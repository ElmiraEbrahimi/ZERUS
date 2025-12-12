package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type CounterClient struct {
	*ChainClient          // embed base client
	contract     *Counter // abigen-generated type from counter_abigen.go
	address      common.Address
}

func NewCounterClient(ctx context.Context, rpcURL string, chainID int64, privKeyHex, contractAddr string) (*CounterClient, error) {
	base, err := NewChainClient(ctx, rpcURL, chainID, privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("init chain client: %w", err)
	}

	addr := common.HexToAddress(contractAddr)

	counter, err := NewCounter(addr, base.Eth) // from abigen
	if err != nil {
		return nil, fmt.Errorf("bind counter: %w", err)
	}

	return &CounterClient{
		ChainClient: base,
		contract:    counter,
		address:     addr,
	}, nil
}

func (c *CounterClient) Get(ctx context.Context) (*big.Int, error) {
	return c.contract.Get(&bind.CallOpts{
		Context: ctx,
	})
}

func (c *CounterClient) Increment(ctx context.Context) (*types.Transaction, error) {
	nonce, err := c.Eth.PendingNonceAt(ctx, c.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice, err := c.Eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivKey, c.ChainID)
	if err != nil {
		return nil, fmt.Errorf("new transactor: %w", err)
	}

	auth.From = c.From
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice
	auth.GasLimit = 200_000 // tweak or later swap to EstimateGas

	tx, err := c.contract.Increment(auth)
	if err != nil {
		return nil, fmt.Errorf("increment counter: %w", err)
	}

	return tx, nil
}

func (c *CounterClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.ChainClient.BlockNumber(ctx)
}
