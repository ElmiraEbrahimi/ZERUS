package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"l2alchemy/internal/eth"
	"l2alchemy/internal/zkkeys"
)

// Handler aggregates dependencies for HTTP endpoints. It holds a reference to
// the CounterClient used to service requests.
type Handler struct {
	counter *eth.CounterClient
	zk      *zkkeys.Manager
}

// NewHandler constructs a new Handler with the given dependencies.
func NewHandler(counter *eth.CounterClient, zk *zkkeys.Manager) *Handler {
	return &Handler{counter: counter, zk: zk}
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

	// ZK key generation and status endpoints.
	// POST /circuits/keygen triggers Groth16 setup and writes pk/vk for both circuits.
	mux.HandleFunc("/circuits/keygen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.generateCircuitKeys(w, r)
	})

	// GET /circuits/keys returns whether pk/vk exist on disk for each circuit.
	mux.HandleFunc("/circuits/keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.circuitKeysStatus(w, r)
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

type circuitKeysStatusResponse struct {
	Circuits []zkkeys.KeyStatus `json:"circuits"`
}

type circuitKeyGenResponse struct {
	Results []zkkeys.KeyGenResult `json:"results"`
}

// generateCircuitKeys handles POST /circuits/keygen.
func (h *Handler) generateCircuitKeys(w http.ResponseWriter, r *http.Request) {
	if h.zk == nil {
		http.Error(w, "zk key manager not configured", http.StatusInternalServerError)
		return
	}
	// These operations can be heavy; set a generous request timeout upstream if needed.
	results, err := h.zk.GenerateAll()
	if err != nil {
		log.Printf("key generation failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(circuitKeyGenResponse{Results: results})
}

// circuitKeysStatus handles GET /circuits/keys.
func (h *Handler) circuitKeysStatus(w http.ResponseWriter, r *http.Request) {
	if h.zk == nil {
		http.Error(w, "zk key manager not configured", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(circuitKeysStatusResponse{Circuits: h.zk.Status()})
}
