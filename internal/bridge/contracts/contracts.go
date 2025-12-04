package contracts

// BridgeContract defines the methods required to interact with a bridge
// smart contract. When integrating your actual Solidity contracts,
// generate Go bindings (e.g. via abigen) and ensure they satisfy this
// interface. Keeping the interface small and focused allows you to swap
// implementations or mock it in tests.
type BridgeContract interface {
    // Deposit initiates a transfer from L1 to L2. Parameters such as
    // recipient, token address and amount should be captured in a
    // struct defined in your implementation.
    Deposit(opts interface{}, args ...interface{}) (interface{}, error)

    // Withdraw initiates a transfer from L2 back to L1.
    Withdraw(opts interface{}, args ...interface{}) (interface{}, error)

    // Other bridge-specific functions can be added here.
}