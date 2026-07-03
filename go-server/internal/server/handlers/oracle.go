package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/oracle_runtime"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	Mode    string `json:"mode,omitempty"`
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
	mode := "l2-direct"
	if h.engine.Cfg != nil && h.engine.Cfg.L1RPCURL != "" && h.engine.Cfg.L1HubContractAddress != "" {
		registered, err := h.registerValidatorsViaL1Hub(r.Context(), ids)
		if err != nil {
			log.Printf("register validators via L1Hub failed: %v", err)
			http.Error(w, "failed to register validators through L1 hub", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(registerValidatorsResponse{NodeIDs: registered, Count: len(registered), Mode: "l1-hub"})
		return
	}

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
	json.NewEncoder(w).Encode(registerValidatorsResponse{NodeIDs: nodeIDs, Count: len(nodeIDs), Mode: mode})
}

func (h *OracleHandler) registerValidatorsViaL1Hub(ctx context.Context, ids []uint) ([]uint, error) {
	cfg := h.engine.Cfg
	oracle := h.engine.Oracle
	if cfg == nil || oracle == nil {
		return nil, fmt.Errorf("oracle engine not initialized")
	}

	l1Client, err := bc.NewChainClient(ctx, cfg.L1RPCURL, cfg.L1ChainID, cfg.NodePK)
	if err != nil {
		return nil, fmt.Errorf("init L1 validator client: %w", err)
	}
	if cfg.L1GasPriceWei > 0 {
		l1Client.GasPriceOverride = big.NewInt(cfg.L1GasPriceWei)
	}

	hub, err := bc.NewL1Hub(common.HexToAddress(cfg.L1HubContractAddress), l1Client.Eth)
	if err != nil {
		return nil, fmt.Errorf("bind L1Hub: %w", err)
	}

	registered := make([]uint, 0, len(ids))
	importIDs := make([]*big.Int, 0, len(ids))
	for _, id := range ids {
		node := oracle.Nodes[id]
		if node == nil || node.Account == nil || node.Account.PublicKey == nil {
			continue
		}

		auth, err := newL1HubTransactor(ctx, l1Client, new(big.Int).Set(node.Account.Balance))
		if err != nil {
			return nil, fmt.Errorf("node %d L1 auth: %w", id, err)
		}
		pubKey := bc.L1HubPublicKey{
			X: node.Account.PublicKey.A.X.BigInt(new(big.Int)),
			Y: node.Account.PublicKey.A.Y.BigInt(new(big.Int)),
		}
		tx, err := hub.RegisterValidatorL1(auth, big.NewInt(int64(id)), pubKey)
		if err != nil {
			return nil, fmt.Errorf("register validator %d on L1: %w", id, err)
		}
		receipt, err := bc.WaitMinedAndLogTxReceipt(ctx, l1Client.Eth, fmt.Sprintf("L1 validator stake registration node=%d", id), tx)
		if err != nil {
			return nil, fmt.Errorf("wait L1 register validator %d: %w", id, err)
		}
		if receipt.Status != 1 {
			return nil, fmt.Errorf("L1 register validator %d reverted (tx=%s)", id, tx.Hash().Hex())
		}

		registered = append(registered, id)
		importIDs = append(importIDs, big.NewInt(int64(id)))
	}

	if len(importIDs) == 0 {
		return registered, nil
	}

	auth, err := newL1HubTransactor(ctx, l1Client, l1ToL2FeeValue(cfg.L1L2ValueWei))
	if err != nil {
		return nil, fmt.Errorf("L1 import auth: %w", err)
	}
	tx, err := hub.BatchImportValidatorsToL2(
		auth,
		importIDs,
		big.NewInt(30_000_000),
		big.NewInt(800),
		l1Client.From,
	)
	if err != nil {
		return nil, fmt.Errorf("request validator import to L2: %w", err)
	}
	receipt, err := bc.WaitMinedAndLogTxReceipt(ctx, l1Client.Eth, "L1->L2 validator import request", tx)
	if err != nil {
		return nil, fmt.Errorf("wait L1 validator import request: %w", err)
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("L1 validator import request reverted (tx=%s)", tx.Hash().Hex())
	}

	return registered, nil
}

func newL1HubTransactor(ctx context.Context, chain *bc.ChainClient, valueWei *big.Int) (*bind.TransactOpts, error) {
	nonce, err := chain.Eth.PendingNonceAt(ctx, chain.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}
	gasPrice := chain.GasPriceOverride
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice, err = chain.Eth.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(chain.PrivKey, chain.ChainID)
	if err != nil {
		return nil, fmt.Errorf("new transactor: %w", err)
	}
	if valueWei == nil {
		valueWei = big.NewInt(0)
	}
	auth.From = chain.From
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = valueWei
	auth.GasPrice = gasPrice
	return auth, nil
}

func l1ToL2FeeValue(raw string) *big.Int {
	if strings.TrimSpace(raw) == "" {
		// Local zkSync direct mailbox requests are payable. This default affects
		// transferred value only; CSV gas remains the receipt gasUsed.
		return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	}
	v, ok := new(big.Int).SetString(strings.TrimSpace(raw), 10)
	if !ok || v.Sign() < 0 {
		return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	}
	return v
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

	var replaceErr error
	if h.engine.Cfg != nil && h.engine.Cfg.L1RPCURL != "" && h.engine.Cfg.L1HubContractAddress != "" {
		replaceErr = node.RequestReplacementFromL1Tx(replaceWithID)
	} else {
		replaceErr = node.ReplaceAccountTx(replaceWithID)
	}
	if replaceErr != nil {
		log.Printf("replace account node=%d failed: %v", nodeID, replaceErr)
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
