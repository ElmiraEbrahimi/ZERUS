package eth

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"l2alchemy/internal/runindex"
)

var (
	txCSVHeader = []string{"index", "datetime", "source", "gas_used", "receipt_json"}
	txCSVMu     sync.Mutex
	txCSVWriter *csv.Writer
	txCSVFile   *os.File
	txCSVIndex  int
)

// SetupTxReceiptCSV configures CSV logging for transaction receipts.
// The caller is responsible for calling CloseTxReceiptCSV.
func SetupTxReceiptCSV(path string) (*os.File, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}
	txCSVMu.Lock()
	defer txCSVMu.Unlock()
	if txCSVWriter != nil {
		return txCSVFile, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if err := runindex.EnsureCSVHeaderWithIndex(path, txCSVHeader); err != nil {
		log.Printf("tx receipt csv: header update failed: %v", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	info, statErr := file.Stat()

	txCSVWriter = csv.NewWriter(file)
	txCSVFile = file
	if idx, err := runindex.Current(path); err != nil {
		log.Printf("tx receipt csv: index init failed: %v", err)
		txCSVIndex = 0
	} else {
		txCSVIndex = idx
	}
	if statErr == nil && info.Size() == 0 {
		if err := txCSVWriter.Write(txCSVHeader); err != nil {
			log.Printf("tx receipt csv: header write failed: %v", err)
		}
		txCSVWriter.Flush()
	}

	return file, nil
}

// CloseTxReceiptCSV flushes and closes the CSV logger if configured.
func CloseTxReceiptCSV() {
	txCSVMu.Lock()
	if txCSVWriter != nil {
		txCSVWriter.Flush()
	}
	if txCSVFile != nil {
		if err := txCSVFile.Close(); err != nil {
			log.Printf("tx receipt csv: close failed: %v", err)
		}
	}
	txCSVWriter = nil
	txCSVFile = nil
	txCSVMu.Unlock()
}

// LogTxSubmitted logs basic metadata when a transaction is sent.
func LogTxSubmitted(label string, tx *types.Transaction) {
	if tx == nil {
		return
	}
	to := "<contract creation>"
	if tx.To() != nil {
		to = tx.To().Hex()
	}
	log.Printf(
		"%s tx submitted: hash=%s nonce=%d to=%s",
		label,
		tx.Hash().Hex(),
		tx.Nonce(),
		to,
	)
}

// LogTxReceipt logs the transaction hash and gas cost details.
func LogTxReceipt(label string, tx *types.Transaction, receipt *types.Receipt) {
	if tx == nil || receipt == nil {
		return
	}

	effectiveGasPrice := receipt.EffectiveGasPrice
	if effectiveGasPrice == nil {
		effectiveGasPrice = tx.GasPrice()
	}

	gasCostWei := new(big.Int)
	if effectiveGasPrice != nil {
		gasCostWei.Mul(new(big.Int).SetUint64(receipt.GasUsed), effectiveGasPrice)
	} else {
		effectiveGasPrice = big.NewInt(0)
	}

	log.Printf(
		"%s tx: hash=%s gasUsed=%d gasPrice=%s gasCost=%s wei",
		label,
		tx.Hash().Hex(),
		receipt.GasUsed,
		effectiveGasPrice.String(),
		gasCostWei.String(),
	)

	writeTxReceiptCSV(label, receipt)
}

// WaitAndLogTxReceipt waits for the transaction receipt and logs gas cost details.
func WaitAndLogTxReceipt(ctx context.Context, backend bind.DeployBackend, label string, tx *types.Transaction) {
	if tx == nil || backend == nil {
		return
	}

	LogTxSubmitted(label, tx)
	receipt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		log.Printf("%s tx: wait for receipt failed (hash=%s): %v", label, tx.Hash().Hex(), err)
		return
	}

	LogTxReceipt(label, tx, receipt)
}

// WaitMinedAndLogTxReceipt waits for the transaction receipt and logs submission + gas details.
func WaitMinedAndLogTxReceipt(ctx context.Context, backend bind.DeployBackend, label string, tx *types.Transaction) (*types.Receipt, error) {
	if tx == nil || backend == nil {
		return nil, errors.New("tx receipt wait requires transaction and backend")
	}

	LogTxSubmitted(label, tx)
	receipt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		log.Printf("%s tx: wait for receipt failed (hash=%s): %v", label, tx.Hash().Hex(), err)
		return nil, err
	}

	LogTxReceipt(label, tx, receipt)
	return receipt, nil
}

func writeTxReceiptCSV(label string, receipt *types.Receipt) {
	if receipt == nil {
		return
	}

	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		log.Printf("tx receipt csv: marshal failed: %v", err)
		receiptJSON = nil
	}

	row := []string{
		strconv.Itoa(txCSVIndex),
		time.Now().Format("20060102_150405"),
		label,
		strconv.FormatUint(receipt.GasUsed, 10),
		string(receiptJSON),
	}

	txCSVMu.Lock()
	defer txCSVMu.Unlock()
	if txCSVWriter == nil {
		return
	}
	if err := txCSVWriter.Write(row); err != nil {
		log.Printf("tx receipt csv: write failed: %v", err)
		return
	}
	txCSVWriter.Flush()
	if err := txCSVWriter.Error(); err != nil {
		log.Printf("tx receipt csv: flush failed: %v", err)
	}
}
