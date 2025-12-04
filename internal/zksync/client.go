package zksync

import (
    "context"
    "time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/rpc"
)

// Client defines the minimal interface required to interact with a zkSync node.
// The methods here mirror common JSON-RPC calls. Keeping the interface small
// allows easy mocking and replacement.
type Client interface {
    SendRawTransaction(ctx context.Context, txData []byte) (common.Hash, error)
    GetTransactionReceipt(ctx context.Context, txHash common.Hash) (map[string]interface{}, error)
    // Call performs a generic RPC call and stores the result in the provided
    // result pointer. This can be used to call arbitrary zkSync-specific
    // endpoints.
    Call(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// RPCClient is a thin wrapper around go-ethereum's rpc.Client. It manages
// timeouts and provides convenience methods for common zkSync operations.
type RPCClient struct {
    client  *rpc.Client
    timeout time.Duration
}

// NewRPCClient dials the provided RPC endpoint and returns a zkSync client.
// A timeout in seconds can be provided to bound network calls.
func NewRPCClient(rpcURL string, timeoutSeconds int) (*RPCClient, error) {
    cli, err := rpc.Dial(rpcURL)
    if err != nil {
        return nil, err
    }
    return &RPCClient{client: cli, timeout: time.Duration(timeoutSeconds) * time.Second}, nil
}

// SendRawTransaction broadcasts a raw, RLP‐encoded transaction to the zkSync
// network. The returned hash can be used to track the transaction status.
func (c *RPCClient) SendRawTransaction(ctx context.Context, txData []byte) (common.Hash, error) {
    var hash common.Hash
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    // The eth_sendRawTransaction method takes a hex‐encoded payload
    err := c.client.CallContext(ctx, &hash, "eth_sendRawTransaction", hexutil.Encode(txData))
    return hash, err
}

// GetTransactionReceipt retrieves the receipt for a transaction hash. The receipt
// is returned as a generic map because zkSync receipts may contain fields
// beyond the standard Ethereum ones.
func (c *RPCClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (map[string]interface{}, error) {
    var receipt map[string]interface{}
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    err := c.client.CallContext(ctx, &receipt, "eth_getTransactionReceipt", txHash)
    return receipt, err
}

// Call performs a low-level JSON-RPC call. It can be used to invoke any
// zkSync-specific method not covered by the convenience helpers. The result
// parameter must be a pointer to the expected response type.
func (c *RPCClient) Call(ctx context.Context, result interface{}, method string, args ...interface{}) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()
    return c.client.CallContext(ctx, result, method, args...)
}

// Ensure RPCClient implements the Client interface.
var _ Client = (*RPCClient)(nil)