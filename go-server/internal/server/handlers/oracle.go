package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

type validatorActionRequest struct {
	NodeID *uint `json:"node_id"`
}

type replaceAccountRequest struct {
	NodeID               *uint   `json:"node_id"`
	ReplaceWithAccountID *uint64 `json:"replace_with_account_id"`
}

type validatorActionResponse struct {
	NodeID uint   `json:"node_id"`
	Status string `json:"status"`
}

type replaceAccountResponse struct {
	NodeID               uint   `json:"node_id"`
	ReplaceWithAccountID uint64 `json:"replace_with_account_id"`
	Status               string `json:"status"`
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
	// Register in ascending node-ID order: the contract assigns leaf indices
	// by insertion order, and the off-chain state tree places account i at
	// position i, so the on-chain and off-chain trees only agree when
	// registration is deterministic (F-19).
	ids := make([]uint, 0, len(oracle.Nodes))
	for id := range oracle.Nodes {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	nodeIDs := make([]uint, 0, len(ids))
	for _, id := range ids {
		node := oracle.Nodes[id]
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

	tokenOne, tokenTwo, err := usr.GetBalance()
	if err != nil {
		log.Printf("get balance default-user failed: %v", err)
		http.Error(w, "failed to get balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balanceResponse{
		User:     usr.Name,
		TokenOne: tokenOne,
		TokenTwo: tokenTwo,
	})
}

// ReplaceValidatorAccount handles POST /validators/replace by calling ReplaceAccountTx
// for the selected validator node.
func (h *OracleHandler) ReplaceValidatorAccount(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.engine == nil || h.engine.Oracle == nil {
		http.Error(w, "oracle engine not initialized", http.StatusInternalServerError)
		return
	}

	req := replaceAccountRequest{}
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	nodeID, err := parseNodeID(r, req.NodeID, h.engine.Oracle.AggregatorID)
	if err != nil {
		http.Error(w, "invalid node_id", http.StatusBadRequest)
		return
	}

	replaceWithID, err := parseReplaceAccountID(r, req.ReplaceWithAccountID)
	if err != nil {
		http.Error(w, "replace_with_account_id is required", http.StatusBadRequest)
		return
	}

	node := h.engine.Oracle.Nodes[nodeID]
	if node == nil {
		http.Error(w, "validator node not found", http.StatusNotFound)
		return
	}

	if err := node.ReplaceAccountTx(replaceWithID); err != nil {
		log.Printf("replace account node=%d failed: %v", nodeID, err)
		http.Error(w, "failed to replace account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replaceAccountResponse{
		NodeID:               nodeID,
		ReplaceWithAccountID: replaceWithID,
		Status:               "ok",
	})
}

// ExitValidatorAccount handles POST /validators/exit by calling ExitTx
// for the selected validator node.
func (h *OracleHandler) ExitValidatorAccount(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.engine == nil || h.engine.Oracle == nil {
		http.Error(w, "oracle engine not initialized", http.StatusInternalServerError)
		return
	}

	req := validatorActionRequest{}
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	nodeID, err := parseNodeID(r, req.NodeID, h.engine.Oracle.AggregatorID)
	if err != nil {
		http.Error(w, "invalid node_id", http.StatusBadRequest)
		return
	}

	node := h.engine.Oracle.Nodes[nodeID]
	if node == nil {
		http.Error(w, "validator node not found", http.StatusNotFound)
		return
	}

	if err := node.ExitTx(); err != nil {
		log.Printf("exit account node=%d failed: %v", nodeID, err)
		http.Error(w, "failed to exit account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validatorActionResponse{
		NodeID: nodeID,
		Status: "ok",
	})
}

// WithdrawValidatorAccount handles POST /validators/withdraw by calling WithdrawAccountTx
// for the selected validator node.
func (h *OracleHandler) WithdrawValidatorAccount(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.engine == nil || h.engine.Oracle == nil {
		http.Error(w, "oracle engine not initialized", http.StatusInternalServerError)
		return
	}

	req := validatorActionRequest{}
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	nodeID, err := parseNodeID(r, req.NodeID, h.engine.Oracle.AggregatorID)
	if err != nil {
		http.Error(w, "invalid node_id", http.StatusBadRequest)
		return
	}

	node := h.engine.Oracle.Nodes[nodeID]
	if node == nil {
		http.Error(w, "validator node not found", http.StatusNotFound)
		return
	}

	if err := node.WithdrawAccountTx(); err != nil {
		log.Printf("withdraw account node=%d failed: %v", nodeID, err)
		http.Error(w, "failed to withdraw account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validatorActionResponse{
		NodeID: nodeID,
		Status: "ok",
	})
}

func decodeJSONBody(r *http.Request, dst any) error {
	if r == nil || r.Body == nil {
		return nil
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func parseNodeID(r *http.Request, bodyNodeID *uint, fallback uint) (uint, error) {
	if bodyNodeID != nil {
		return *bodyNodeID, nil
	}
	val, ok, err := parseUintQuery(r, "node_id")
	if err != nil {
		return 0, err
	}
	if ok {
		return uint(val), nil
	}
	return fallback, nil
}

func parseReplaceAccountID(r *http.Request, bodyReplaceID *uint64) (uint64, error) {
	if bodyReplaceID != nil {
		return *bodyReplaceID, nil
	}
	val, ok, err := parseUintQuery(r, "replace_with_account_id")
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("replace_with_account_id is required")
	}
	return val, nil
}

func parseUintQuery(r *http.Request, key string) (uint64, bool, error) {
	if r == nil {
		return 0, false, nil
	}
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, false, nil
	}
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, true, err
	}
	return val, true, nil
}
