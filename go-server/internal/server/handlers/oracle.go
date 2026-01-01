package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"l2alchemy/internal/oracle_runtime"
)

// OracleHandler serves endpoints backed by the oracle runtime engine.
type OracleHandler struct {
	engine *oracle_runtime.OracleEngine
}

// NewOracleHandler constructs a new OracleHandler with the given dependency.
func NewOracleHandler(engine *oracle_runtime.OracleEngine) *OracleHandler {
	return &OracleHandler{engine: engine}
}

type registerUserResponse struct {
	User   string `json:"user"`
	TxHash string `json:"tx_hash"`
}

type registerValidatorsResponse struct {
	NodeIDs []uint `json:"node_ids"`
	Count   int    `json:"count"`
}

// RegisterDefaultUser handles POST /users/default/register by calling
// RegisterUserTx for the user named "default-user".
func (h *OracleHandler) RegisterDefaultUser(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.engine == nil {
		http.Error(w, "oracle engine not initialized", http.StatusInternalServerError)
		return
	}

	usr := h.engine.Users["default-user"]
	if usr == nil {
		http.Error(w, "default-user not found", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	// RegisterUserTx blocks until mined; run it under a request-scoped timeout.
	_ = ctx // reserved for future wiring when RegisterUserTx accepts a context.
	txHash, err := usr.RegisterUserTx()
	if err != nil {
		log.Printf("register default-user failed: %v", err)
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerUserResponse{User: usr.Name, TxHash: txHash})
}

// RegisterValidators handles POST /validators/register by registering all
// validator nodes in the oracle runtime.
func (h *OracleHandler) RegisterValidators(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.engine == nil || h.engine.Oracle == nil {
		http.Error(w, "oracle engine not initialized", http.StatusInternalServerError)
		return
	}

	oracle := h.engine.Oracle
	nodeIDs := make([]uint, 0, len(oracle.Nodes))
	for id, node := range oracle.Nodes {
		if node == nil {
			continue
		}
		if err := node.RegisterValidatorTx(); err != nil {
			log.Printf("register validator %d failed: %v", id, err)
			http.Error(w, "failed to register validators", http.StatusInternalServerError)
			return
		}
		nodeIDs = append(nodeIDs, id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerValidatorsResponse{NodeIDs: nodeIDs, Count: len(nodeIDs)})
}
