package eth

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

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
}

// WaitAndLogTxReceipt waits for the transaction receipt and logs gas cost details.
func WaitAndLogTxReceipt(ctx context.Context, backend bind.DeployBackend, label string, tx *types.Transaction) {
	if tx == nil || backend == nil {
		return
	}

	receipt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		log.Printf("%s tx: wait for receipt failed (hash=%s): %v", label, tx.Hash().Hex(), err)
		return
	}

	LogTxReceipt(label, tx, receipt)
}
