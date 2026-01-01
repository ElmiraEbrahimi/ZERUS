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

type burnResponse struct {
	User           string `json:"user"`
	TxHash         string `json:"tx_hash"`
	CommitmentHash string `json:"commitment_hash"`
}

type withdrawResponse struct {
	User   string `json:"user"`
	TxHash string `json:"tx_hash"`
}

type balanceResponse struct {
	User     string `json:"user"`
	TokenOne uint   `json:"token_one"`
	TokenTwo uint   `json:"token_two"`
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

// BurnDefaultUser handles POST /users/default/burn by calling BurnTx for the
// user named "default-user".
func (h *OracleHandler) BurnDefaultUser(w http.ResponseWriter, r *http.Request) {
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
	_ = ctx // reserved for future wiring when BurnTx accepts a context.

	txHash, commitmentHash, err := usr.BurnTx()
	if err != nil {
		log.Printf("burn default-user failed: %v", err)
		http.Error(w, "failed to burn user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(burnResponse{
		User:           usr.Name,
		TxHash:         txHash,
		CommitmentHash: commitmentHash,
	})
}

// WithdrawDefaultUser handles POST /users/default/withdraw by calling WithdrawTx
// for the user named "default-user".
func (h *OracleHandler) WithdrawDefaultUser(w http.ResponseWriter, r *http.Request) {
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
	_ = ctx // reserved for future wiring when WithdrawTx accepts a context.

	txHash, err := usr.WithdrawTx()
	if err != nil {
		log.Printf("withdraw default-user failed: %v", err)
		http.Error(w, "failed to withdraw user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawResponse{User: usr.Name, TxHash: txHash})
}

// GetDefaultUserBalance handles GET /users/default/balance by calling GetBalance.
func (h *OracleHandler) GetDefaultUserBalance(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.engine == nil {
		http.Error(w, "oracle engine not initialized", http.StatusInternalServerError)
		return
	}

	usr := h.engine.Users["default-user"]
	if usr == nil {
		http.Error(w, "default-user not found", http.StatusNotFound)
		return
	}

	tokenOne, tokenTwo := usr.GetBalance()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balanceResponse{
		User:     usr.Name,
		TokenOne: tokenOne,
		TokenTwo: tokenTwo,
	})
}
