package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// requestIDKey is used to store a request id in the request context.
type requestIDKey struct{}

// RequestIDFromContext returns the request id stored in ctx, if present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(requestIDKey{})
	s, ok := v.(string)
	return s, ok
}

// WithRequestLogging wraps next with request logging + panic recovery.
//
// A single log line is emitted for every request that hits the server, including
// 404s and method-not-allowed responses. In case of panic, the panic is logged
// and a 500 is returned.
func WithRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := requestIDFromRequest(r)
		ctx := context.WithValue(r.Context(), requestIDKey{}, reqID)
		r = r.WithContext(ctx)

		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		remoteIP := clientIP(r)

		// Ensure we always log exactly once.
		defer func() {
			dur := time.Since(start)
			// Emit the access log line.
			log.Printf(
				"http_request request_id=%s method=%s path=%q status=%d bytes=%d duration_ms=%d remote_ip=%s ua=%q",
				reqID,
				r.Method,
				r.URL.RequestURI(),
				lrw.status,
				lrw.bytes,
				dur.Milliseconds(),
				remoteIP,
				r.UserAgent(),
			)
		}()

		// Panic recovery (still allows the deferred access log above to run).
		defer func() {
			if rec := recover(); rec != nil {
				// Best-effort 500 if nothing has been written yet.
				if !lrw.wroteHeader {
					lrw.WriteHeader(http.StatusInternalServerError)
				}
				log.Printf(
					"http_panic request_id=%s method=%s path=%q remote_ip=%s panic=%v stack=%q",
					reqID,
					r.Method,
					r.URL.RequestURI(),
					remoteIP,
					rec,
					string(debug.Stack()),
				)
				// Only write an error body if we still can.
				if !lrw.wroteBody {
					http.Error(lrw, "internal server error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(lrw, r)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
	wroteBody   bool
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.status = statusCode
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	if n > 0 {
		w.wroteBody = true
	}
	return n, err
}

// Preserve optional interfaces (e.g. Flusher) when present.
func (w *loggingResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Preserve Hijacker if present.
func (w *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

// Preserve Pusher if present (HTTP/2).
func (w *loggingResponseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Preserve ReadFrom when supported for io.Copy optimizations.
func (w *loggingResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	if rf, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		if !w.wroteHeader {
			w.WriteHeader(http.StatusOK)
		}
		n, err := rf.ReadFrom(r)
		w.bytes += int(n)
		if n > 0 {
			w.wroteBody = true
		}
		return n, err
	}
	return 0, http.ErrNotSupported
}

func requestIDFromRequest(r *http.Request) string {
	// Honor standard/request-id headers when provided by a reverse proxy.
	if v := strings.TrimSpace(r.Header.Get("X-Request-Id")); v != "" {
		return v
	}
	if v := strings.TrimSpace(r.Header.Get("X-Correlation-Id")); v != "" {
		return v
	}
	// 16 random bytes -> 32 hex chars.
	b := make([]byte, 16)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	// Fallback: timestamp based.
	return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
}

func clientIP(r *http.Request) string {
	// If behind a proxy, use X-Forwarded-For / X-Real-IP (best effort).
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		// XFF can be a comma-separated list: client, proxy1, proxy2
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}
	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
