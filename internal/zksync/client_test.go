package zksync

import "testing"

// TestNewRPCClientInvalidURL ensures that constructing a zksync client with
// an invalid endpoint returns an error. Dialling a non-existent host
// typically results in a network error.
func TestNewRPCClientInvalidURL(t *testing.T) {
    // Using an unsupported scheme should cause dialing to fail immediately.
    _, err := NewRPCClient("foo://bar", 1)
    if err == nil {
        t.Fatal("expected error for invalid RPC URL, got nil")
    }
}