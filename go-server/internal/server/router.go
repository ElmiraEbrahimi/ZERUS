package server

import "net/http"

// NewRouter constructs a new http.ServeMux and registers all routes using the
// provided Handler. This helper is useful in tests and when embedding the
// router in a custom HTTP server.
func NewRouter(h *Handler) *http.ServeMux {
    mux := http.NewServeMux()
    RegisterRoutes(mux, h)
    return mux
}