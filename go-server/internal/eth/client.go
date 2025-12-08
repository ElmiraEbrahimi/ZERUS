package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// counterABIJSON is the JSON ABI definition for the Counter contract. It allows
// method calls to be encoded and decoded at runtime without relying on
// generated Go bindings. If the Counter contract changes this constant must be
// updated accordingly.
const counterABIJSON = `[
    {
        "inputs": [],
        "stateMutability": "nonpayable",
        "type": "constructor"
    },
    {
        "inputs": [],
        "name": "get",
        "outputs": [
            {
                "internalType": "uint256",
                "name": "",
                "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "increment",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "value",
        "outputs": [
            {
                "internalType": "uint256",
                "name": "",
                "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    }
]`

// CounterClient provides convenience methods for interacting with a deployed
// Counter contract over JSON‑RPC. It wraps a generic BoundContract to avoid
// having to use abigen.
type CounterClient struct {
    ethClient *ethclient.Client
    contract  *bind.BoundContract
    abi       abi.ABI
    address   common.Address
    chainID   *big.Int
    privKey   *ecdsa.PrivateKey
}

// NewCounterClient dials an RPC endpoint, parses the private key and initialises
// a new CounterClient bound to the given contract address.
func NewCounterClient(ctx context.Context, rpcURL string, chainID int64, privKeyHex string, contractAddr string) (*CounterClient, error) {
    // Remove optional 0x prefix from the private key
    if strings.HasPrefix(privKeyHex, "0x") {
        privKeyHex = privKeyHex[2:]
    }
    privKeyBytes, err := hex.DecodeString(privKeyHex)
    if err != nil {
        return nil, err
    }
    pk, err := crypto.ToECDSA(privKeyBytes)
    if err != nil {
        return nil, err
    }
    parsedABI, err := abi.JSON(strings.NewReader(counterABIJSON))
    if err != nil {
        return nil, err
    }
    client, err := ethclient.DialContext(ctx, rpcURL)
    if err != nil {
        return nil, err
    }
    address := common.HexToAddress(contractAddr)
    bound := bind.NewBoundContract(address, parsedABI, client, client, client)

    return &CounterClient{
        ethClient: client,
        contract:  bound,
        abi:       parsedABI,
        address:   address,
        chainID:   big.NewInt(chainID),
        privKey:   pk,
    }, nil
}

// Get retrieves the current value of the counter by performing an eth_call.
func (c *CounterClient) Get(ctx context.Context) (*big.Int, error) {
    // Call expects a pointer to a slice of any for the outputs.
    var outputs []any
    if err := c.contract.Call(&bind.CallOpts{Context: ctx}, &outputs, "get"); err != nil {
        return nil, err
    }

    if len(outputs) != 1 {
        return nil, fmt.Errorf("counter get: expected 1 return value, got %d", len(outputs))
    }

    // The ABI decoder should give us *big.Int or big.Int as the first element.
    switch v := outputs[0].(type) {
    case *big.Int:
        return v, nil
    case big.Int:
        // Copy to heap so caller always gets *big.Int.
        vv := new(big.Int).Set(&v)
        return vv, nil
    default:
        return nil, fmt.Errorf("counter get: unexpected return type %T", v)
    }
}
// Increment sends a transaction to call the increment() function. It returns
// the transaction object but does not wait for the receipt. Consumers can
// choose to monitor the tx separately.
func (c *CounterClient) Increment(ctx context.Context) (*types.Transaction, error) {
    from := crypto.PubkeyToAddress(c.privKey.PublicKey)
    nonce, err := c.ethClient.PendingNonceAt(ctx, from)
    if err != nil {
        return nil, err
    }
    gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
    if err != nil {
        return nil, err
    }
    auth, err := bind.NewKeyedTransactorWithChainID(c.privKey, c.chainID)
    if err != nil {
        return nil, err
    }
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)       // no ETH is sent
    auth.GasPrice = gasPrice
    // Provide a generous gas limit to cover typical zkSync L2 execution. Users
    // can adjust this by estimating gas ahead of time if desired.
    auth.GasLimit = uint64(200_000)
    return c.contract.Transact(auth, "increment")
}

// BlockNumber returns the latest block number from the underlying RPC client.
// It can be used by health checks to ensure the node is responsive.
func (c *CounterClient) BlockNumber(ctx context.Context) (uint64, error) {
    return c.ethClient.BlockNumber(ctx)
}