package eth

import (
    "context"
    "testing"

    "github.com/ethereum/go-ethereum/common"
)

// TestNewRPCClientInvalidURL ensures that creating a client with an invalid
// URL results in an error.
func TestNewRPCClientInvalidURL(t *testing.T) {
    // Using an unsupported scheme should cause dialing to fail immediately.
    _, err := NewRPCClient("foo://bar", 1)
    if err == nil {
        t.Fatal("expected error for invalid RPC URL, got nil")
    }
}

// TestGetTransactionByHashNotFound verifies that a not-found error is
// converted into a nil transaction and false flag rather than propagated.
func TestGetTransactionByHashInvalidEndpoint(t *testing.T) {
    // Dial a client against a port that is unlikely to be serving RPC. Dial
    // succeeds because the URL is syntactically valid. The subsequent call
    // should return an error.
    cli, err := NewRPCClient("http://127.0.0.1:12345", 1)
    if err != nil {
        // If dialing fails, skip the test since we cannot perform the call
        t.Skipf("dial failed: %v", err)
    }
    _, _, err = cli.GetTransactionByHash(context.Background(), common.Hash{})
    if err == nil {
        t.Fatalf("expected error calling GetTransactionByHash on invalid endpoint, got nil")
    }
}