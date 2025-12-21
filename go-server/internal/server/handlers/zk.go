package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"l2alchemy/internal/zkkeys"
)

// ZKHandler serves endpoints for ZK key management.
type ZKHandler struct {
	zk *zkkeys.Manager
}

// NewZKHandler constructs a new ZKHandler with the given dependency.
func NewZKHandler(zk *zkkeys.Manager) *ZKHandler {
	return &ZKHandler{zk: zk}
}

type circuitKeysStatusResponse struct {
	Circuits []zkkeys.KeyStatus `json:"circuits"`
}

type circuitKeyGenResponse struct {
	Results []zkkeys.KeyGenResult `json:"results"`
}

// generateCircuitKeys handles POST /circuits/keygen.
func (h *ZKHandler) GenerateCircuitKeys(w http.ResponseWriter, r *http.Request) {
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
func (h *ZKHandler) CircuitKeysStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(circuitKeysStatusResponse{Circuits: h.zk.Status()})
}
