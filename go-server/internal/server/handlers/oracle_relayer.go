package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	"l2alchemy/internal/config"
	"l2alchemy/internal/eth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// l1MessageSentTopic is the topic of the zkSync system messenger's
// L1MessageSent(address indexed sender, bytes32 indexed hash, bytes message)
// event, from which the raw L2->L1 message bytes are recovered.
var l1MessageSentTopic = crypto.Keccak256Hash([]byte("L1MessageSent(address,bytes32,bytes)"))

const l2ToL1SystemMessenger = "0x0000000000000000000000000000000000008008"

// OracleL2ToL1Relayer forwards the oracle's typed L2->L1 messages
// (checkpoints, exit/withdraw requests, import/replacement results) to the
// L1 Hub: it recovers the message bytes from the L2 receipt, waits for the
// batch to become executable, fetches the zkSync log-inclusion proof, and
// calls L1Hub.finalizeFromL2 (paper SIV-B Step 5; F-02).
type OracleL2ToL1Relayer struct {
	l2         *eth.ChainClient
	l1         *eth.ChainClient
	hub        *eth.L1Hub
	oracleAddr common.Address
	useDirect  bool
}

// NewOracleL2ToL1Relayer builds a relayer from configuration; it returns
// (nil, nil) when the L1 Hub is not configured, in which case the runtime
// operates without L1 anchoring (single-chain simulation).
func NewOracleL2ToL1Relayer(ctx context.Context, cfg *config.Config) (*OracleL2ToL1Relayer, error) {
	if cfg == nil || cfg.L1HubContractAddress == "" || cfg.L1RPCURL == "" {
		return nil, nil
	}

	l2, err := eth.NewChainClient(ctx, cfg.RPCURL, cfg.ChainID, cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("relayer: init L2 client: %w", err)
	}
	l1, err := eth.NewChainClient(ctx, cfg.L1RPCURL, cfg.L1ChainID, cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("relayer: init L1 client: %w", err)
	}
	if cfg.L1GasPriceWei > 0 {
		l1.GasPriceOverride = big.NewInt(cfg.L1GasPriceWei)
	}

	hubAddr := common.HexToAddress(cfg.L1HubContractAddress)
	hub, err := eth.NewL1Hub(hubAddr, l1.Eth)
	if err != nil {
		return nil, fmt.Errorf("relayer: bind L1Hub: %w", err)
	}

	return &OracleL2ToL1Relayer{
		l2:         l2,
		l1:         l1,
		hub:        hub,
		oracleAddr: common.HexToAddress(cfg.OracleContractAddress),
		useDirect:  cfg.L1UseDirectMessaging,
	}, nil
}

// RelayTx relays the oracle L2->L1 message contained in the given L2
// transaction to the L1 Hub. Intended to run in its own goroutine; failures
// are logged, not fatal (the message can be relayed again later since the
// Hub deduplicates by message hash).
func (r *OracleL2ToL1Relayer) RelayTx(txHash common.Hash) {
	if r == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), l2ToL1RelayTimeout)
	defer cancel()

	log.Printf("oracle relayer: L2->L1 relay start tx=%s", txHash.Hex())

	message, err := r.messageBytes(ctx, txHash)
	if err != nil {
		log.Printf("oracle relayer: recover message failed tx=%s: %v", txHash.Hex(), err)
		return
	}

	logCandidate, err := waitForL2ToL1Log(ctx, r.l2.Eth.Client(), txHash)
	if err != nil {
		log.Printf("oracle relayer: wait for L2->L1 log failed tx=%s: %v", txHash.Hex(), err)
		return
	}
	if expected := padAddress(r.oracleAddr.Hex()); logCandidate.LogKey != "" && !strings.EqualFold(logCandidate.LogKey, expected) {
		log.Printf("oracle relayer: log key mismatch expected=%s actual=%s", expected, logCandidate.LogKey)
	}

	interopMode := ""
	if r.useDirect {
		interopMode = "proof_based_gw"
	}

	proofBlockNumber := logCandidate.L1BatchNumber
	if proofBlockNumber == 0 {
		proofBlockNumber = logCandidate.L2BlockNumber
	}

	for attempt := 1; ; attempt++ {
		if ctx.Err() != nil {
			log.Printf("oracle relayer: relay timed out tx=%s", txHash.Hex())
			return
		}

		if logCandidate.L1BatchNumber != 0 {
			ready, err := isL1BatchExecutable(ctx, r.l2.Eth.Client(), logCandidate.L1BatchNumber)
			if err != nil || !ready {
				time.Sleep(l2ToL1RelayPollInterval)
				continue
			}
		}

		proof, err := fetchL2ToL1LogProof(ctx, r.l2.Eth.Client(), txHash, uint64(logCandidate.LogIndex), interopMode)
		if err != nil || proof == nil || len(proof.Proof) == 0 {
			time.Sleep(l2ToL1RelayPollInterval)
			continue
		}

		txNumber := logCandidate.L2TxNumberBatch
		if proof.Log != nil {
			if proofTxNum, ok := parseUint64(proof.Log["txNumberInBlock"]); ok {
				txNumber = proofTxNum
			}
		}
		if txNumber > math.MaxUint16 {
			log.Printf("oracle relayer: invalid tx index=%d", txNumber)
			return
		}

		proofBytes, err := parseProof(proof.Proof)
		if err != nil {
			log.Printf("oracle relayer: invalid proof tx=%s: %v", txHash.Hex(), err)
			return
		}

		auth, err := r.l1Transactor(ctx)
		if err != nil {
			log.Printf("oracle relayer: build L1 transactor failed: %v", err)
			return
		}

		tx, err := r.hub.FinalizeFromL2(
			auth,
			message,
			new(big.Int).SetUint64(proofBlockNumber),
			new(big.Int).SetUint64(uint64(logCandidate.LogIndex)),
			uint16(txNumber),
			proofBytes,
		)
		if err == nil {
			log.Printf("oracle relayer: finalizeFromL2 submitted l2_tx=%s l1_tx=%s (attempt=%d)", txHash.Hex(), tx.Hash().Hex(), attempt)
			receipt, err := bind.WaitMined(ctx, r.l1.Eth, tx)
			if err != nil {
				log.Printf("oracle relayer: wait finalizeFromL2 receipt failed: %v", err)
				return
			}
			if receipt.Status == 1 {
				log.Printf("oracle relayer: message finalized on L1 (l2_tx=%s)", txHash.Hex())
			} else {
				log.Printf("oracle relayer: finalizeFromL2 reverted (l2_tx=%s)", txHash.Hex())
			}
			return
		}

		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "already consumed") {
			log.Printf("oracle relayer: message already consumed (l2_tx=%s)", txHash.Hex())
			return
		}
		// A deterministic on-chain revert (e.g. "unknown root", "validator
		// mismatch", "withdraw delay not elapsed") cannot succeed by
		// retrying the same message against the same state: give up once
		// and leave the message for a later replay (the Hub deduplicates by
		// message hash, so re-relaying after the state changes is safe).
		if strings.Contains(errMsg, "execution reverted") || strings.Contains(errMsg, "revert") {
			log.Printf("oracle relayer: finalizeFromL2 reverted, not retrying (l2_tx=%s): %v", txHash.Hex(), err)
			return
		}
		log.Printf("oracle relayer: finalizeFromL2 failed (attempt=%d), retrying: %v", attempt, err)
		time.Sleep(l2ToL1RelayPollInterval)
	}
}

// messageBytes recovers the raw L2->L1 message from the transaction's
// L1MessageSent system-messenger event where the sender is the oracle.
func (r *OracleL2ToL1Relayer) messageBytes(ctx context.Context, txHash common.Hash) ([]byte, error) {
	receipt, err := r.l2.Eth.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("fetch receipt: %w", err)
	}

	messenger := common.HexToAddress(l2ToL1SystemMessenger)
	for _, lg := range receipt.Logs {
		if lg.Address != messenger || len(lg.Topics) < 2 {
			continue
		}
		if lg.Topics[0] != l1MessageSentTopic {
			continue
		}
		sender := common.BytesToAddress(lg.Topics[1].Bytes())
		if sender != r.oracleAddr {
			continue
		}
		// Data is one ABI-encoded dynamic `bytes` argument.
		if len(lg.Data) < 64 {
			return nil, fmt.Errorf("malformed L1MessageSent data (len=%d)", len(lg.Data))
		}
		length := new(big.Int).SetBytes(lg.Data[32:64]).Uint64()
		if uint64(len(lg.Data)) < 64+length {
			return nil, fmt.Errorf("truncated L1MessageSent data (len=%d want=%d)", len(lg.Data), 64+length)
		}
		return lg.Data[64 : 64+length], nil
	}
	return nil, fmt.Errorf("no oracle L1MessageSent event in tx %s", txHash.Hex())
}

func (r *OracleL2ToL1Relayer) l1Transactor(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := r.l1.Eth.PendingNonceAt(ctx, r.l1.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}
	gasPrice := r.l1.GasPriceOverride
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice, err = r.l1.Eth.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(r.l1.PrivKey, r.l1.ChainID)
	if err != nil {
		return nil, fmt.Errorf("new transactor: %w", err)
	}
	auth.From = r.l1.From
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	return auth, nil
}
