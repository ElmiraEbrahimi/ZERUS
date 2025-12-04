package eth

import (
    "context"
    "errors"
    "time"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

// Client defines the behaviour required from an Ethereum L1 client.
// By coding to this interface we can later replace the implementation
// with mocks or alternative clients without changing downstream code.
type Client interface {
    GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
    GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
    SendTransaction(ctx context.Context, tx *types.Transaction) error
    FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

// RPCClient is a minimal wrapper around go-ethereum's ethclient.Client.
// It satisfies the Client interface defined above.
type RPCClient struct {
    client *ethclient.Client
    // timeout used for network calls
    timeout time.Duration
}

// NewRPCClient connects to an Ethereum node via the given RPC URL.
// A dial error is returned if the connection cannot be established.
func NewRPCClient(rpcURL string, timeoutSeconds int) (*RPCClient, error) {
    cli, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, err
    }
    return &RPCClient{client: cli, timeout: time.Duration(timeoutSeconds) * time.Second}, nil
}

// GetBlockByHash retrieves the block with the given hash. If the context
// is cancelled or times out, an error is returned.
func (c *RPCClient) GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    return c.client.BlockByHash(ctx, hash)
}

// GetTransactionByHash fetches a transaction by its hash. The bool return
// value indicates whether the transaction was found in the blockchain.
func (c *RPCClient) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    tx, pending, err := c.client.TransactionByHash(ctx, hash)
    if err != nil {
        // convert not found errors to boolean
        if err.Error() == ethereum.NotFound.Error() {
            return nil, false, nil
        }
        return nil, false, err
    }
    return tx, pending, nil
}

// SendTransaction broadcasts a signed transaction to the network.
func (c *RPCClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    return c.client.SendTransaction(ctx, tx)
}

// FilterLogs executes a log filter query against the node and returns the
// matching logs. This method is useful for event subscriptions or indexing.
func (c *RPCClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    return c.client.FilterLogs(ctx, q)
}

// Ensure RPCClient implements the Client interface at compile time.
var _ Client = (*RPCClient)(nil)

// ErrNotImplemented is returned by unimplemented functions.
var ErrNotImplemented = errors.New("not implemented")