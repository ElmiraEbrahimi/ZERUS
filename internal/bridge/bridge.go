package bridge

import (
    "context"
    "errors"
    "sync"

    "github.com/google/uuid"

    "github.com/example/bridge-service/internal/eth"
    "github.com/example/bridge-service/internal/zksync"
)

// MessageStatus represents the lifecycle state of a bridge message.
type MessageStatus string

const (
    // StatusPending indicates the message has been accepted by the service
    // and will be processed shortly.
    StatusPending MessageStatus = "pending"
    // StatusSent marks a message that has been forwarded to the target chain.
    StatusSent MessageStatus = "sent"
    // StatusFailed marks a message whose processing failed.
    StatusFailed MessageStatus = "failed"
)

// L1ToL2Request captures the minimum information needed to send a message
// from Ethereum L1 to zkSync L2. Additional fields such as value or gas
// limits can be added later as required.
type L1ToL2Request struct {
    From string `json:"from"`
    To   string `json:"to"`
    Data string `json:"data"`
}

// L2ToL1Request defines a message from L2 to L1. Fields can be extended
// when integrating with actual bridge contracts.
type L2ToL1Request struct {
    From string `json:"from"`
    To   string `json:"to"`
    Data string `json:"data"`
}

// Service ties together Ethereum and zkSync clients and provides a simple
// abstraction over cross-layer messaging. Messages and their statuses are
// tracked in-memory for now; persistent storage can be added later via the
// database layer.
type Service struct {
    ethClient   eth.Client
    zksyncClient zksync.Client
    mu          sync.RWMutex
    statuses    map[string]MessageStatus
}

// NewService constructs a bridge service using the provided clients. The
// returned service is ready to process messages.
func NewService(ethClient eth.Client, zClient zksync.Client) *Service {
    return &Service{
        ethClient:    ethClient,
        zksyncClient: zClient,
        statuses:     make(map[string]MessageStatus),
    }
}

// SendL1ToL2Message accepts a request to bridge data from L1 to L2. It
// generates a unique identifier for the request and records its status. The
// underlying transaction submission is currently simulated; integration with
// actual contracts will be added later.
func (s *Service) SendL1ToL2Message(ctx context.Context, req L1ToL2Request) (string, error) {
    id := uuid.New().String()
    s.mu.Lock()
    s.statuses[id] = StatusPending
    s.mu.Unlock()
    // TODO: implement actual call to ethClient to send message to L2 bridge
    // For now we simply mark the message as sent after a short delay.
    go func() {
        // In a real implementation, this goroutine would wait on transaction
        // confirmations and update status accordingly. Here we mark as sent
        // immediately to satisfy the API contract.
        s.mu.Lock()
        s.statuses[id] = StatusSent
        s.mu.Unlock()
    }()
    return id, nil
}

// SendL2ToL1Message bridges a message from L2 back to L1. Similar to
// SendL1ToL2Message, this is a placeholder for real contract interactions.
func (s *Service) SendL2ToL1Message(ctx context.Context, req L2ToL1Request) (string, error) {
    id := uuid.New().String()
    s.mu.Lock()
    s.statuses[id] = StatusPending
    s.mu.Unlock()
    // TODO: call zksyncClient to send message to L1
    go func() {
        s.mu.Lock()
        s.statuses[id] = StatusSent
        s.mu.Unlock()
    }()
    return id, nil
}

// TrackStatus retrieves the status of a previously submitted message. If the
// identifier is unknown, an error is returned.
func (s *Service) TrackStatus(ctx context.Context, id string) (MessageStatus, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    status, ok := s.statuses[id]
    if !ok {
        return "", errors.New("unknown message id")
    }
    return status, nil
}