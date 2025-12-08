package server

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "counterapi/internal/eth"
)

// Handler aggregates dependencies for HTTP endpoints. It holds a reference to
// the CounterClient used to service requests.
type Handler struct {
    counter *eth.CounterClient
}

// NewHandler constructs a new Handler with the given CounterClient.
func NewHandler(counter *eth.CounterClient) *Handler {
    return &Handler{counter: counter}
}

// counterResponse is the shape of responses returned from GET /counter.
type counterResponse struct {
    Value string `json:"value"`
}

// txResponse is the shape of responses returned from POST /counter/increment.
type txResponse struct {
    TxHash string `json:"tx_hash"`
}

// RegisterRoutes binds HTTP paths to handler methods. It can be called from
// main() to configure the default ServeMux.
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
    mux.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        h.getCounter(w, r)
    })

    mux.HandleFunc("/counter/increment", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        h.incrementCounter(w, r)
    })

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        h.health(w, r)
    })
}

// getCounter handles GET /counter by reading the current value from the
// on‑chain counter and returning it as JSON. It fails with a 500 status code
// if the RPC call cannot be completed.
func (h *Handler) getCounter(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    val, err := h.counter.Get(ctx)
    if err != nil {
        log.Printf("error calling get(): %v", err)
        http.Error(w, "failed to query counter", http.StatusInternalServerError)
        return
    }
    resp := counterResponse{Value: val.String()}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// incrementCounter handles POST /counter/increment by sending a transaction to
// the counter contract. On success it returns the resulting transaction hash.
// Errors are reported with a 500 status code.
func (h *Handler) incrementCounter(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
    defer cancel()
    tx, err := h.counter.Increment(ctx)
    if err != nil {
        log.Printf("error calling increment(): %v", err)
        http.Error(w, "failed to send transaction", http.StatusInternalServerError)
        return
    }
    resp := txResponse{TxHash: tx.Hash().Hex()}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// health handles GET /health by probing the latest block number. If the call
// succeeds it returns 200 and the block number; otherwise it returns 500.
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    block, err := h.counter.BlockNumber(ctx)
    if err != nil {
        log.Printf("health check failed: %v", err)
        http.Error(w, "unhealthy", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "block_number": block,
    })
}