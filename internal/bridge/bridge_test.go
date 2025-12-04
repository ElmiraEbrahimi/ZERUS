package bridge

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
)

// stubEthClient is a no-op implementation of the eth.Client interface used
// in tests. All methods return zero values.
type stubEthClient struct{}

func (stubEthClient) GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
    return nil, nil
}

func (stubEthClient) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
    return nil, false, nil
}

func (stubEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
    return nil
}

func (stubEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
    return nil, nil
}

// stubZkClient is a no-op implementation of the zksync.Client interface used
// in tests. All methods return zero values.
type stubZkClient struct{}

func (stubZkClient) SendRawTransaction(ctx context.Context, txData []byte) (common.Hash, error) {
    return common.Hash{}, nil
}

func (stubZkClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (map[string]interface{}, error) {
    return nil, nil
}

func (stubZkClient) Call(ctx context.Context, result interface{}, method string, args ...interface{}) error {
    return nil
}

// TestSendAndTrack ensures that sending a message records its status and
// transitions from pending to sent.
func TestSendAndTrack(t *testing.T) {
    svc := NewService(stubEthClient{}, stubZkClient{})
    id, err := svc.SendL1ToL2Message(context.Background(), L1ToL2Request{From: "0x", To: "0x", Data: ""})
    if err != nil {
        t.Fatalf("SendL1ToL2Message returned error: %v", err)
    }
    // Immediately after sending, status should be pending
    status, err := svc.TrackStatus(context.Background(), id)
    if err != nil {
        t.Fatalf("TrackStatus returned error: %v", err)
    }
    if status != StatusPending {
        t.Fatalf("expected status %s, got %s", StatusPending, status)
    }
    // Wait a short time for the goroutine to mark as sent
    time.Sleep(20 * time.Millisecond)
    status, err = svc.TrackStatus(context.Background(), id)
    if err != nil {
        t.Fatalf("TrackStatus returned error after sleep: %v", err)
    }
    if status != StatusSent {
        t.Fatalf("expected status %s after send, got %s", StatusSent, status)
    }
}

// TestUnknownID verifies that querying an unknown message ID returns an error.
func TestUnknownID(t *testing.T) {
    svc := NewService(stubEthClient{}, stubZkClient{})
    _, err := svc.TrackStatus(context.Background(), "nonexistent")
    if err == nil {
        t.Fatal("expected error for unknown id, got nil")
    }
    if !errors.Is(err, errors.New("unknown message id")) {
        // We compare error strings because the service returns a new error each time
        if err.Error() != "unknown message id" {
            t.Fatalf("unexpected error: %v", err)
        }
    }
}