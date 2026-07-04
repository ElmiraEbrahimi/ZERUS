package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"l2alchemy/internal/eth"
)

const (
	defaultL2GasLimit           = 2_000_000
	defaultL2GasPerPubdataLimit = 800
	l1MessengerSystemAddress    = "0x0000000000000000000000000000000000008008"
	l2ToL1RelayTimeout          = 4 * time.Minute
	l2ToL1RelayPollInterval     = 2 * time.Second
)

type MessengerHandler struct {
	l1               *eth.L1MessengerClient
	l2               *eth.L2MessengerClient
	defaultL2ChainID uint64
	useDirect        bool
}

func NewMessengerHandler(l1 *eth.L1MessengerClient, l2 *eth.L2MessengerClient, defaultL2ChainID uint64, useDirect bool) *MessengerHandler {
	return &MessengerHandler{
		l1:               l1,
		l2:               l2,
		defaultL2ChainID: defaultL2ChainID,
		useDirect:        useDirect,
	}
}

type messageRequest struct {
	Message string `json:"message"`
}

type sendToL2Request struct {
	Message         string `json:"message"`
	L2GasLimit      uint64 `json:"l2_gas_limit"`
	L2GasPerPubdata uint64 `json:"l2_gas_per_pubdata"`
	ValueWei        string `json:"value_wei"`
	L2ChainID       uint64 `json:"l2_chain_id"`
}

type receiveFromL2Request struct {
	Message           string   `json:"message"`
	L2BlockNumber     uint64   `json:"l2_block_number"`
	L2LogIndex        *uint64  `json:"l2_log_index"`
	L2MessageIndex    *uint64  `json:"l2_message_index"`
	L2TxNumberInBlock uint16   `json:"l2_tx_number_in_block"`
	Proof             []string `json:"proof"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func (h *MessengerHandler) SendL2ToL1(w http.ResponseWriter, r *http.Request) {
	if h.l2 == nil {
		http.Error(w, "L2 messenger not configured", http.StatusServiceUnavailable)
		return
	}

	var req messageRequest
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	log.Printf("messaging: L2 -> L1 send message=%q", req.Message)
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	tx, err := h.l2.SendToL1(ctx, req.Message)
	if err != nil {
		log.Printf("messaging: L2 sendToL1 failed: %v", err)
		http.Error(w, "failed to send L2 -> L1 message", http.StatusInternalServerError)
		return
	}

	go waitForTxReceipt("L2->L1 send", h.l2.Eth, tx)
	if h.l1 != nil {
		go h.autoRelayL2ToL1(tx.Hash(), req.Message)
	}
	writeJSON(w, txResponse{TxHash: tx.Hash().Hex()})
}

func (h *MessengerHandler) SendL1ToL2(w http.ResponseWriter, r *http.Request) {
	if h.l1 == nil {
		http.Error(w, "L1 messenger not configured", http.StatusServiceUnavailable)
		return
	}

	var req sendToL2Request
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	if req.L2GasLimit == 0 {
		req.L2GasLimit = defaultL2GasLimit
	}
	if req.L2GasPerPubdata == 0 {
		req.L2GasPerPubdata = defaultL2GasPerPubdataLimit
	}
	valueWei, err := parseBigInt(req.ValueWei)
	if err != nil {
		http.Error(w, "invalid value_wei", http.StatusBadRequest)
		return
	}
	autoValue := valueWei == nil || valueWei.Sign() == 0

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	directChainID := req.L2ChainID
	if directChainID == 0 {
		directChainID = h.defaultL2ChainID
	}

	if req.L2ChainID != 0 || h.useDirect {
		if directChainID == 0 {
			http.Error(w, "l2_chain_id is required", http.StatusBadRequest)
			return
		}
		log.Printf(
			"messaging: L1 -> L2 direct send chain_id=%d message=%q l2_gas_limit=%d l2_gas_per_pubdata=%d value_wei=%s",
			directChainID,
			req.Message,
			req.L2GasLimit,
			req.L2GasPerPubdata,
			valueWei.String(),
		)
		directValue := valueWei
		if autoValue {
			directValue = nil
		}
		tx, err := h.l1.SendToL2Direct(ctx, directChainID, req.Message, req.L2GasLimit, req.L2GasPerPubdata, directValue)
		if err != nil {
			log.Printf("messaging: L1 sendToL2Direct failed: %v", err)
			http.Error(w, "failed to send L1 -> L2 direct message", http.StatusInternalServerError)
			return
		}

		go waitForTxReceipt("L1->L2 direct send", h.l1.Eth, tx)
		writeJSON(w, txResponse{TxHash: tx.Hash().Hex()})
		return
	}

	log.Printf("messaging: L1 -> L2 send message=%q l2_gas_limit=%d l2_gas_per_pubdata=%d value_wei=%s", req.Message, req.L2GasLimit, req.L2GasPerPubdata, valueWei.String())
	tx, err := h.l1.SendToL2(ctx, req.Message, req.L2GasLimit, req.L2GasPerPubdata, valueWei)
	if err != nil {
		if isLegacyMailboxErr(err) && directChainID != 0 && !h.useDirect {
			log.Printf("messaging: legacy L1->L2 failed, retrying direct chain_id=%d: %v", directChainID, err)
			var directValue *big.Int
			if !autoValue {
				directValue = valueWei
			}
			tx, err = h.l1.SendToL2Direct(ctx, directChainID, req.Message, req.L2GasLimit, req.L2GasPerPubdata, directValue)
		}
		if err != nil {
			log.Printf("messaging: L1 sendToL2 failed: %v", err)
			http.Error(w, "failed to send L1 -> L2 message", http.StatusInternalServerError)
			return
		}
	}

	go waitForTxReceipt("L1->L2 send", h.l1.Eth, tx)
	writeJSON(w, txResponse{TxHash: tx.Hash().Hex()})
}

func isLegacyMailboxErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "legacy interface") || strings.Contains(msg, "only available for era")
}

type l2ToL1LogProof struct {
	ID    uint64
	Proof []string
	Log   map[string]any
}

func (h *MessengerHandler) autoRelayL2ToL1(txHash common.Hash, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), l2ToL1RelayTimeout)
	defer cancel()

	log.Printf("messaging: L2->L1 relay start tx=%s", txHash.Hex())
	logCandidate, err := waitForL2ToL1Log(ctx, h.l2.Eth.Client(), txHash)
	if err != nil {
		log.Printf("messaging: L2->L1 relay wait log failed: %v", err)
		return
	}

	l1Addr := ""
	l2Addr := ""
	if h.l1 != nil {
		l1Addr = h.l1.Address().Hex()
	}
	if h.l2 != nil {
		l2Addr = h.l2.Address().Hex()
	}
	log.Printf("messaging: L2->L1 relay context l1_messenger=%s l2_messenger=%s use_direct=%t", l1Addr, l2Addr, h.useDirect)

	logIndex := logCandidate.LogIndex
	l1BatchNumber := logCandidate.L1BatchNumber
	l2BlockNumber := logCandidate.L2BlockNumber
	l2TxNumberInBlock := logCandidate.L2TxNumberBatch

	interopMode := ""
	if h.useDirect {
		interopMode = "proof_based_gw"
	}
	expectedKey := ""
	expectedValue := ""
	if l2Addr != "" {
		expectedKey = padAddress(l2Addr)
	}
	if l1Addr != "" {
		payload, err := encodeL2ToL1Payload(common.HexToAddress(l1Addr), message)
		if err != nil {
			log.Printf("messaging: L2->L1 relay payload encode failed: %v", err)
		} else {
			expectedValue = crypto.Keccak256Hash(payload).Hex()
		}
	}
	if expectedKey != "" || expectedValue != "" {
		log.Printf("messaging: L2->L1 relay expected key=%s value=%s", expectedKey, expectedValue)
	}
	if expectedKey != "" && logCandidate.LogKey != "" && !strings.EqualFold(logCandidate.LogKey, expectedKey) {
		log.Printf("messaging: L2->L1 relay key mismatch expected=%s actual=%s", expectedKey, logCandidate.LogKey)
	}
	if expectedValue != "" && logCandidate.Value != "" && !strings.EqualFold(logCandidate.Value, expectedValue) {
		log.Printf("messaging: L2->L1 relay value mismatch expected=%s actual=%s", expectedValue, logCandidate.Value)
	}

	proofBlockNumber := l1BatchNumber
	if proofBlockNumber == 0 {
		proofBlockNumber = l2BlockNumber
	}

	for attempt := 1; ; attempt++ {
		if l1BatchNumber != 0 {
			ready, err := isL1BatchExecutable(ctx, h.l2.Eth.Client(), l1BatchNumber)
			if err != nil {
				if ctx.Err() != nil {
					log.Printf("messaging: L2->L1 relay batch check failed (attempt=%d): %v", attempt, err)
					return
				}
				log.Printf("messaging: L2->L1 relay batch check failed (attempt=%d), retrying: %v", attempt, err)
				time.Sleep(l2ToL1RelayPollInterval)
				continue
			}
			if !ready {
				if ctx.Err() != nil {
					log.Printf("messaging: L2->L1 relay batch not executable (attempt=%d)", attempt)
					return
				}
				log.Printf("messaging: L2->L1 relay batch not executable (attempt=%d), waiting", attempt)
				time.Sleep(l2ToL1RelayPollInterval)
				continue
			}
		}

		proof, err := fetchL2ToL1LogProof(ctx, h.l2.Eth.Client(), txHash, uint64(logIndex), interopMode)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("messaging: L2->L1 relay proof failed (attempt=%d): %v", attempt, err)
				return
			}
			log.Printf("messaging: L2->L1 relay proof failed (attempt=%d), retrying: %v", attempt, err)
			time.Sleep(l2ToL1RelayPollInterval)
			continue
		}
		if proof == nil || len(proof.Proof) == 0 {
			if ctx.Err() != nil {
				log.Printf("messaging: L2->L1 relay proof missing (attempt=%d)", attempt)
				return
			}
			log.Printf("messaging: L2->L1 relay proof missing (attempt=%d), retrying", attempt)
			time.Sleep(l2ToL1RelayPollInterval)
			continue
		}
		l2LogIndex := proof.ID
		if proof.ID != uint64(logIndex) {
			log.Printf("messaging: L2->L1 relay log index mismatch proof_id=%d log_index=%d", proof.ID, logIndex)
		}
		if attempt == 1 {
			log.Printf(
				"messaging: L2->L1 relay params l2_block=%d l1_batch=%d proof_block=%d l2_log_index=%d l2_tx_in_block=%d log_index=%d proof_len=%d",
				l2BlockNumber,
				l1BatchNumber,
				proofBlockNumber,
				uint64(logIndex),
				l2TxNumberInBlock,
				logIndex,
				len(proof.Proof),
			)
		}

		attemptTxNumber := l2TxNumberInBlock
		if proof.Log != nil {
			if proofTxNum, ok := parseUint64(proof.Log["txNumberInBlock"]); ok {
				if proofTxNum > math.MaxUint16 {
					log.Printf("messaging: L2->L1 relay proof txNumberInBlock too large=%d", proofTxNum)
					return
				}
				if proofTxNum != attemptTxNumber {
					log.Printf("messaging: L2->L1 relay txNumberInBlock mismatch proof=%d receipt=%d", proofTxNum, attemptTxNumber)
				}
				attemptTxNumber = proofTxNum
			}
			if shard, ok := parseUint64(proof.Log["l2ShardId"]); ok {
				isService := formatRaw(proof.Log["isService"])
				log.Printf("messaging: L2->L1 relay proof log l2ShardId=%d isService=%s", shard, isService)
			}
		}

		if attemptTxNumber > math.MaxUint16 {
			log.Printf("messaging: L2->L1 relay invalid tx index=%d", attemptTxNumber)
			return
		}

		proofBytes, err := parseProof(proof.Proof)
		if err != nil {
			log.Printf("messaging: L2->L1 relay invalid proof: %v", err)
			return
		}

		tx, err := h.l1.ReceiveFromL2(ctx, message, proofBlockNumber, l2LogIndex, uint16(attemptTxNumber), proofBytes)
		if err == nil {
			log.Printf("messaging: L2->L1 relay receive tx submitted (attempt=%d) hash=%s", attempt, tx.Hash().Hex())
			go waitForTxReceipt("L2->L1 receive", h.l1.Eth, tx)
			return
		}

		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "message already consumed") || strings.Contains(errMsg, "already consumed") {
			log.Printf("messaging: L2->L1 relay already received (attempt=%d): %v", attempt, err)
			return
		}
		if ctx.Err() != nil {
			log.Printf("messaging: L2->L1 relay receive failed (attempt=%d): %v", attempt, err)
			return
		}
		log.Printf("messaging: L2->L1 relay receive failed (attempt=%d), retrying: %v", attempt, err)
		time.Sleep(l2ToL1RelayPollInterval)
	}
}

func waitForL2ToL1Log(ctx context.Context, rpcClient *rpc.Client, txHash common.Hash) (*l2ToL1LogCandidate, error) {
	ticker := time.NewTicker(l2ToL1RelayPollInterval)
	defer ticker.Stop()

	for {
		receipt, err := fetchReceipt(ctx, rpcClient, txHash)
		if err != nil {
			return nil, err
		}
		if receipt != nil {
			candidate, ok := extractL2ToL1Log(receipt)
			if ok && candidate != nil {
				return candidate, nil
			}
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting for L2->L1 log")
		case <-ticker.C:
		}
	}
}

func waitForL2ToL1LogValue(ctx context.Context, rpcClient *rpc.Client, txHash common.Hash, expectedValue string) (*l2ToL1LogCandidate, error) {
	ticker := time.NewTicker(l2ToL1RelayPollInterval)
	defer ticker.Stop()

	expectedValue = strings.ToLower(strings.TrimSpace(expectedValue))
	for {
		receipt, err := fetchReceipt(ctx, rpcClient, txHash)
		if err != nil {
			return nil, err
		}
		if receipt != nil {
			candidate, ok := extractL2ToL1LogValue(receipt, expectedValue)
			if ok && candidate != nil {
				return candidate, nil
			}
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting for L2->L1 log value %s", expectedValue)
		case <-ticker.C:
		}
	}
}

func fetchReceipt(ctx context.Context, rpcClient *rpc.Client, txHash common.Hash) (map[string]any, error) {
	var out map[string]any
	if err := rpcClient.CallContext(ctx, &out, "eth_getTransactionReceipt", txHash.Hex()); err != nil {
		return nil, err
	}
	return out, nil
}

func isL1BatchExecutable(ctx context.Context, rpcClient *rpc.Client, batchNumber uint64) (bool, error) {
	var out map[string]any
	if err := rpcClient.CallContext(ctx, &out, "zks_getL1BatchDetails", batchNumber); err != nil {
		// If the method is not supported, treat as ready.
		if strings.Contains(strings.ToLower(err.Error()), "method not found") {
			return true, nil
		}
		return false, err
	}
	if out == nil {
		return false, nil
	}
	status := strings.ToLower(parseString(out["status"]))
	if isReadyBatchStatus(status) {
		return true, nil
	}
	if hasNonZeroHash(parseString(out["executeTxHash"])) {
		return true, nil
	}
	if hasNonZeroHash(parseString(out["proveTxHash"])) {
		return true, nil
	}
	return false, nil
}

func isReadyBatchStatus(status string) bool {
	switch strings.ToLower(status) {
	case "executed", "verified", "fastfinalized", "fast_finalized":
		return true
	default:
		return false
	}
}

func hasNonZeroHash(value string) bool {
	if value == "" {
		return false
	}
	if strings.HasPrefix(value, "0x") {
		trimmed := strings.TrimPrefix(value, "0x")
		trimmed = strings.TrimLeft(trimmed, "0")
		return trimmed != ""
	}
	return true
}

func extractL2ToL1Log(receipt map[string]any) (*l2ToL1LogCandidate, bool) {
	if receipt == nil {
		return nil, false
	}

	l1BatchNumberFallback, _ := parseUint64(receipt["l1BatchNumber"])
	l1BatchTxIndexFallback, _ := parseUint64(receipt["l1BatchTxIndex"])
	l2BlockNumberFallback, _ := parseUint64(receipt["blockNumber"])

	logsRaw, ok := receipt["l2ToL1Logs"].([]any)
	if !ok || len(logsRaw) == 0 {
		return nil, false
	}

	receiptFrom := parseString(receipt["from"])
	receiptTo := parseString(receipt["to"])
	paddedFrom := padAddress(receiptFrom)
	paddedTo := padAddress(receiptTo)
	var fallback *l2ToL1LogCandidate

	for idx, raw := range logsRaw {
		logMap, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		sender := parseString(logMap["sender"])
		if sender == "" || !strings.EqualFold(sender, l1MessengerSystemAddress) {
			continue
		}

		logKey := strings.ToLower(parseString(logMap["key"]))
		logValue := strings.ToLower(parseString(logMap["value"]))
		keyMatches := false
		if paddedFrom != "" && strings.EqualFold(logKey, paddedFrom) {
			keyMatches = true
		}
		if paddedTo != "" && strings.EqualFold(logKey, paddedTo) {
			keyMatches = true
		}

		l1BatchNumber, ok := parseUint64(logMap["l1BatchNumber"])
		if !ok {
			l1BatchNumber = l1BatchNumberFallback
		}
		if l1BatchNumber == 0 {
			continue
		}

		l2BlockNumber, ok := parseUint64(logMap["blockNumber"])
		if !ok {
			l2BlockNumber = l2BlockNumberFallback
		}

		l2TxNumberInBatch, ok := parseUint64(logMap["txIndexInL1Batch"])
		if !ok {
			l2TxNumberInBatch, ok = parseUint64(logMap["txNumberInBatch"])
		}
		if !ok {
			l2TxNumberInBatch = l1BatchTxIndexFallback
		}
		if !ok || l2TxNumberInBatch == 0 {
			l2TxNumberInBatch, ok = parseUint64(logMap["transactionIndex"])
		}
		if !ok || l2TxNumberInBatch == 0 {
			l2TxNumberInBatch, _ = parseUint64(logMap["txNumberInBlock"])
		}

		receiptLogIndex, ok := parseUint64(logMap["transactionLogIndex"])
		if !ok {
			receiptLogIndex, ok = parseUint64(logMap["logIndex"])
		}
		if !ok {
			receiptLogIndex = uint64(idx)
		}
		// zks_getL2ToL1LogProof expects the index inside receipt.l2ToL1Logs,
		// not the EVM receipt logIndex. Using the receipt logIndex can prove a
		// different system log and makes L1Hub.finalizeFromL2 revert.
		l2ToL1LogIndex := uint64(idx)

		txIndexInL1BatchRaw := formatRaw(logMap["txIndexInL1Batch"])
		txNumberInBatchRaw := formatRaw(logMap["txNumberInBatch"])
		transactionIndexRaw := formatRaw(logMap["transactionIndex"])
		txNumberInBlockRaw := formatRaw(logMap["txNumberInBlock"])
		l2ShardRaw := formatRaw(logMap["l2ShardId"])
		if l2ShardRaw == "" {
			l2ShardRaw = formatRaw(logMap["shardId"])
		}
		isServiceRaw := formatRaw(logMap["isService"])

		candidate := &l2ToL1LogCandidate{
			LogIndex:        int(l2ToL1LogIndex),
			L2BlockNumber:   l2BlockNumber,
			L1BatchNumber:   l1BatchNumber,
			L2TxNumberBatch: l2TxNumberInBatch,
			Sender:          sender,
			LogKey:          logKey,
			Value:           logValue,
			Raw:             logMap,
		}

		if keyMatches {
			log.Printf(
				"messaging: L2->L1 relay log sender=%s key=%s value=%s l1_batch=%d tx_in_batch=%d l2_to_l1_log_index=%d receipt_log_index=%d (l2ShardId=%s isService=%s txIndexInL1Batch=%s txNumberInBatch=%s transactionIndex=%s txNumberInBlock=%s from=%s to=%s)",
				sender,
				logKey,
				logValue,
				l1BatchNumber,
				l2TxNumberInBatch,
				l2ToL1LogIndex,
				receiptLogIndex,
				l2ShardRaw,
				isServiceRaw,
				txIndexInL1BatchRaw,
				txNumberInBatchRaw,
				transactionIndexRaw,
				txNumberInBlockRaw,
				paddedFrom,
				paddedTo,
			)
			return candidate, true
		}

		if fallback == nil {
			fallback = candidate
		}
	}

	if fallback != nil {
		log.Printf(
			"messaging: L2->L1 relay log fallback sender=%s key=%s value=%s l1_batch=%d tx_in_batch=%d log_index=%d (from=%s to=%s)",
			fallback.Sender,
			fallback.LogKey,
			fallback.Value,
			fallback.L1BatchNumber,
			fallback.L2TxNumberBatch,
			fallback.LogIndex,
			paddedFrom,
			paddedTo,
		)
		return fallback, true
	}

	return nil, false
}

func extractL2ToL1LogValue(receipt map[string]any, expectedValue string) (*l2ToL1LogCandidate, bool) {
	if receipt == nil || expectedValue == "" {
		return nil, false
	}

	l1BatchNumberFallback, _ := parseUint64(receipt["l1BatchNumber"])
	l1BatchTxIndexFallback, _ := parseUint64(receipt["l1BatchTxIndex"])
	l2BlockNumberFallback, _ := parseUint64(receipt["blockNumber"])

	logsRaw, ok := receipt["l2ToL1Logs"].([]any)
	if !ok || len(logsRaw) == 0 {
		return nil, false
	}

	for idx, raw := range logsRaw {
		logMap, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		sender := parseString(logMap["sender"])
		if sender == "" || !strings.EqualFold(sender, l1MessengerSystemAddress) {
			continue
		}

		logValue := strings.ToLower(parseString(logMap["value"]))
		if !strings.EqualFold(logValue, expectedValue) {
			continue
		}

		l1BatchNumber, ok := parseUint64(logMap["l1BatchNumber"])
		if !ok {
			l1BatchNumber = l1BatchNumberFallback
		}
		if l1BatchNumber == 0 {
			continue
		}

		l2BlockNumber, ok := parseUint64(logMap["blockNumber"])
		if !ok {
			l2BlockNumber = l2BlockNumberFallback
		}

		l2TxNumberInBatch, ok := parseUint64(logMap["txIndexInL1Batch"])
		if !ok {
			l2TxNumberInBatch, ok = parseUint64(logMap["txNumberInBatch"])
		}
		if !ok {
			l2TxNumberInBatch = l1BatchTxIndexFallback
		}
		if !ok || l2TxNumberInBatch == 0 {
			l2TxNumberInBatch, ok = parseUint64(logMap["transactionIndex"])
		}
		if !ok || l2TxNumberInBatch == 0 {
			l2TxNumberInBatch, _ = parseUint64(logMap["txNumberInBlock"])
		}

		receiptLogIndex, ok := parseUint64(logMap["transactionLogIndex"])
		if !ok {
			receiptLogIndex, ok = parseUint64(logMap["logIndex"])
		}
		if !ok {
			receiptLogIndex = uint64(idx)
		}
		// zks_getL2ToL1LogProof expects the index inside receipt.l2ToL1Logs.
		// The EVM receipt log index may point to another system log and makes
		// L1Hub.finalizeFromL2 fail proof verification.
		l2ToL1LogIndex := uint64(idx)

		candidate := &l2ToL1LogCandidate{
			LogIndex:        int(l2ToL1LogIndex),
			L2BlockNumber:   l2BlockNumber,
			L1BatchNumber:   l1BatchNumber,
			L2TxNumberBatch: l2TxNumberInBatch,
			Sender:          sender,
			LogKey:          strings.ToLower(parseString(logMap["key"])),
			Value:           logValue,
			Raw:             logMap,
		}
		log.Printf(
			"messaging: L2->L1 relay matched value=%s l1_batch=%d tx_in_batch=%d l2_to_l1_log_index=%d receipt_log_index=%d",
			logValue,
			l1BatchNumber,
			l2TxNumberInBatch,
			l2ToL1LogIndex,
			receiptLogIndex,
		)
		return candidate, true
	}

	return nil, false
}

func fetchL2ToL1LogProof(ctx context.Context, rpcClient *rpc.Client, txHash common.Hash, logIndex uint64, interopMode string) (*l2ToL1LogProof, error) {
	var out map[string]any
	var err error
	if interopMode != "" {
		err = rpcClient.CallContext(ctx, &out, "zks_getL2ToL1LogProof", txHash.Hex(), logIndex, interopMode)
	} else {
		err = rpcClient.CallContext(ctx, &out, "zks_getL2ToL1LogProof", txHash.Hex(), logIndex)
	}
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	id, ok := parseUint64(out["id"])
	if !ok {
		return nil, fmt.Errorf("missing proof id")
	}
	proofItems := parseStringSlice(out["proof"])
	if len(proofItems) == 0 {
		return nil, fmt.Errorf("missing proof items")
	}

	logMap, _ := out["log"].(map[string]any)
	return &l2ToL1LogProof{
		ID:    id,
		Proof: proofItems,
		Log:   logMap,
	}, nil
}

func parseUint64(raw any) (uint64, bool) {
	switch v := raw.(type) {
	case nil:
		return 0, false
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return 0, false
		}
		if strings.HasPrefix(s, "0x") {
			n, err := hexutil.DecodeUint64(s)
			return n, err == nil
		}
		n, err := strconv.ParseUint(s, 10, 64)
		return n, err == nil
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case json.Number:
		n, err := v.Int64()
		if err != nil || n < 0 {
			return 0, false
		}
		return uint64(n), true
	default:
		return 0, false
	}
}

func parseString(raw any) string {
	if raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}

func parseStringSlice(raw any) []string {
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		val := parseString(item)
		if val == "" {
			continue
		}
		out = append(out, val)
	}
	return out
}

type l2ToL1LogCandidate struct {
	LogIndex        int
	L2BlockNumber   uint64
	L1BatchNumber   uint64
	L2TxNumberBatch uint64
	Sender          string
	LogKey          string
	Value           string
	Raw             map[string]any
}

func padAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ""
	}
	if !strings.HasPrefix(addr, "0x") {
		return ""
	}
	hex := strings.TrimPrefix(addr, "0x")
	if len(hex) > 40 {
		hex = hex[len(hex)-40:]
	}
	return "0x" + strings.Repeat("0", 64-len(hex)) + strings.ToLower(hex)
}

func formatRaw(raw any) string {
	if raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatUint(uint64(v), 10)
	default:
		return fmt.Sprintf("%v", raw)
	}
}

func (h *MessengerHandler) SendL1ToL2Direct(w http.ResponseWriter, r *http.Request) {
	if h.l1 == nil {
		http.Error(w, "L1 messenger not configured", http.StatusServiceUnavailable)
		return
	}

	var req sendToL2Request
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	if req.L2ChainID == 0 {
		http.Error(w, "l2_chain_id is required", http.StatusBadRequest)
		return
	}

	if req.L2GasLimit == 0 {
		req.L2GasLimit = defaultL2GasLimit
	}
	if req.L2GasPerPubdata == 0 {
		req.L2GasPerPubdata = defaultL2GasPerPubdataLimit
	}
	valueWei, err := parseBigInt(req.ValueWei)
	if err != nil {
		http.Error(w, "invalid value_wei", http.StatusBadRequest)
		return
	}

	log.Printf("messaging: L1 -> L2 direct send chain_id=%d message=%q l2_gas_limit=%d l2_gas_per_pubdata=%d value_wei=%s", req.L2ChainID, req.Message, req.L2GasLimit, req.L2GasPerPubdata, valueWei.String())
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	tx, err := h.l1.SendToL2Direct(ctx, req.L2ChainID, req.Message, req.L2GasLimit, req.L2GasPerPubdata, valueWei)
	if err != nil {
		log.Printf("messaging: L1 sendToL2Direct failed: %v", err)
		http.Error(w, "failed to send L1 -> L2 direct message", http.StatusInternalServerError)
		return
	}

	go waitForTxReceipt("L1->L2 direct send", h.l1.Eth, tx)
	writeJSON(w, txResponse{TxHash: tx.Hash().Hex()})
}

func (h *MessengerHandler) ReceiveFromL2(w http.ResponseWriter, r *http.Request) {
	if h.l1 == nil {
		http.Error(w, "L1 messenger not configured", http.StatusServiceUnavailable)
		return
	}

	var req receiveFromL2Request
	if err := decodeJSONBody(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	if len(req.Proof) == 0 {
		http.Error(w, "proof is required", http.StatusBadRequest)
		return
	}

	proof, err := parseProof(req.Proof)
	if err != nil {
		http.Error(w, "invalid proof", http.StatusBadRequest)
		return
	}
	logIndex := uint64(0)
	if req.L2LogIndex != nil {
		logIndex = *req.L2LogIndex
	} else if req.L2MessageIndex != nil {
		logIndex = *req.L2MessageIndex
	}

	log.Printf(
		"messaging: L2 -> L1 receive message=%q block=%d log_index=%d tx_number=%d proof_len=%d",
		req.Message,
		req.L2BlockNumber,
		logIndex,
		req.L2TxNumberInBlock,
		len(proof),
	)
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	tx, err := h.l1.ReceiveFromL2(ctx, req.Message, req.L2BlockNumber, logIndex, req.L2TxNumberInBlock, proof)
	if err != nil {
		log.Printf("messaging: L1 receiveFromL2 failed: %v", err)
		http.Error(w, "failed to receive L2 -> L1 message", http.StatusInternalServerError)
		return
	}

	go waitForTxReceipt("L2->L1 receive", h.l1.Eth, tx)
	writeJSON(w, txResponse{TxHash: tx.Hash().Hex()})
}

func (h *MessengerHandler) GetL1LastFromL2(w http.ResponseWriter, r *http.Request) {
	if h.l1 == nil {
		http.Error(w, "L1 messenger not configured", http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	msg, err := h.l1.LastMessageFromL2(ctx)
	if err != nil {
		log.Printf("messaging: L1 lastMessageFromL2 failed: %v", err)
		http.Error(w, "failed to query L1 message", http.StatusInternalServerError)
		return
	}
	log.Printf("messaging: L1 lastMessageFromL2=%q", msg)
	writeJSON(w, messageResponse{Message: msg})
}

func (h *MessengerHandler) GetL2LastFromL1(w http.ResponseWriter, r *http.Request) {
	if h.l2 == nil {
		http.Error(w, "L2 messenger not configured", http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	msg, err := h.l2.LastMessageFromL1(ctx)
	if err != nil {
		log.Printf("messaging: L2 lastMessageFromL1 failed: %v", err)
		http.Error(w, "failed to query L2 message", http.StatusInternalServerError)
		return
	}
	log.Printf("messaging: L2 lastMessageFromL1=%q", msg)
	writeJSON(w, messageResponse{Message: msg})
}

func parseBigInt(raw string) (*big.Int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return big.NewInt(0), nil
	}
	n, ok := new(big.Int).SetString(raw, 0)
	if !ok {
		return nil, fmt.Errorf("invalid big int %q", raw)
	}
	return n, nil
}

func parseProof(items []string) ([][32]byte, error) {
	proof := make([][32]byte, 0, len(items))
	for _, item := range items {
		val := strings.TrimSpace(item)
		if val == "" {
			continue
		}
		if !strings.HasPrefix(val, "0x") {
			return nil, fmt.Errorf("missing 0x prefix")
		}
		hash := common.HexToHash(val)
		proof = append(proof, [32]byte(hash))
	}
	if len(proof) == 0 {
		return nil, fmt.Errorf("empty proof")
	}
	return proof, nil
}

func encodeL2ToL1Payload(l1Messenger common.Address, message string) ([]byte, error) {
	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	args := abi.Arguments{{Type: addressType}, {Type: stringType}}
	return args.Pack(l1Messenger, message)
}

func waitForTxReceipt(label string, backend bind.DeployBackend, tx *types.Transaction) {
	waitCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	eth.WaitAndLogTxReceipt(waitCtx, backend, label, tx)
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
