package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"l2alchemy/internal/eth"
)

// CounterHandler serves endpoints backed by the on-chain counter client.
type CounterHandler struct {
	counter *eth.CounterClient
}

// NewCounterHandler constructs a new CounterHandler with the given dependency.
func NewCounterHandler(counter *eth.CounterClient) *CounterHandler {
	return &CounterHandler{counter: counter}
}

// counterResponse is the shape of responses returned from GET /counter.
type counterResponse struct {
	Value string `json:"value"`
}

// txResponse is the shape of responses returned from transaction endpoints.
type txResponse struct {
	TxHash string `json:"tx_hash"`
}

// getCounter handles GET /counter by reading the current value from the
// on-chain counter and returning it as JSON. It fails with a 500 status code
// if the RPC call cannot be completed.
func (h *CounterHandler) GetCounter(w http.ResponseWriter, r *http.Request) {
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
func (h *CounterHandler) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	tx, err := h.counter.Increment(ctx)
	if err != nil {
		log.Printf("error calling increment(): %v", err)
		http.Error(w, "failed to send transaction", http.StatusInternalServerError)
		return
	}
	go func() {
		waitCtx, waitCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer waitCancel()
		eth.WaitAndLogTxReceipt(waitCtx, h.counter.Eth, "counter increment", tx)
	}()
	resp := txResponse{TxHash: tx.Hash().Hex()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// health handles GET /health by probing the latest block number. If the call
// succeeds it returns 200 and the block number; otherwise it returns 500.
func (h *CounterHandler) Health(w http.ResponseWriter, r *http.Request) {
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
