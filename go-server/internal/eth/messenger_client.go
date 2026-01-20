package eth

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const mailboxBaseCostABI = `[{"inputs":[{"internalType":"uint256","name":"gasPrice","type":"uint256"},{"internalType":"uint256","name":"l2GasLimit","type":"uint256"},{"internalType":"uint256","name":"l2GasPerPubdataByteLimit","type":"uint256"}],"name":"l2TransactionBaseCost","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

const mailboxBaseCostBridgehubABI = `[{"inputs":[{"internalType":"uint256","name":"chainId","type":"uint256"},{"internalType":"uint256","name":"gasPrice","type":"uint256"},{"internalType":"uint256","name":"l2GasLimit","type":"uint256"},{"internalType":"uint256","name":"l2GasPerPubdataByteLimit","type":"uint256"}],"name":"l2TransactionBaseCost","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

type L1MessengerClient struct {
	*ChainClient
	contract *L1Messenger
	address  common.Address
}

func NewL1MessengerClient(ctx context.Context, rpcURL string, chainID int64, privKeyHex, contractAddr string, gasPriceWei int64) (*L1MessengerClient, error) {
	base, err := NewChainClient(ctx, rpcURL, chainID, privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("init chain client: %w", err)
	}
	if gasPriceWei > 0 {
		base.GasPriceOverride = big.NewInt(gasPriceWei)
	}

	addr := common.HexToAddress(contractAddr)
	messenger, err := NewL1Messenger(addr, base.Eth)
	if err != nil {
		return nil, fmt.Errorf("bind L1 messenger: %w", err)
	}

	return &L1MessengerClient{
		ChainClient: base,
		contract:    messenger,
		address:     addr,
	}, nil
}

func (c *L1MessengerClient) Address() common.Address {
	return c.address
}

func (c *L1MessengerClient) LastMessageFromL2(ctx context.Context) (string, error) {
	return c.contract.LastMessageFromL2(&bind.CallOpts{Context: ctx})
}

func (c *L1MessengerClient) LastMessageToL2(ctx context.Context) (string, error) {
	return c.contract.LastMessageToL2(&bind.CallOpts{Context: ctx})
}

func (c *L1MessengerClient) SendToL2(ctx context.Context, message string, l2GasLimit, l2GasPerPubdata uint64, valueWei *big.Int) (*types.Transaction, error) {
	if valueWei == nil || valueWei.Sign() == 0 {
		baseCost, err := c.estimateBaseCost(ctx, nil, l2GasLimit, l2GasPerPubdata)
		if err != nil {
			return nil, err
		}
		valueWei = baseCost
	}
	auth, err := newMessengerTransactor(ctx, c.ChainClient, valueWei)
	if err != nil {
		return nil, err
	}
	return c.contract.SendToL2(
		auth,
		message,
		new(big.Int).SetUint64(l2GasLimit),
		new(big.Int).SetUint64(l2GasPerPubdata),
	)
}

func (c *L1MessengerClient) SendToL2Direct(ctx context.Context, l2ChainID uint64, message string, l2GasLimit, l2GasPerPubdata uint64, valueWei *big.Int) (*types.Transaction, error) {
	if valueWei == nil || valueWei.Sign() == 0 {
		chainID := new(big.Int).SetUint64(l2ChainID)
		baseCost, err := c.estimateBaseCost(ctx, chainID, l2GasLimit, l2GasPerPubdata)
		if err != nil {
			return nil, err
		}
		valueWei = baseCost
	}
	auth, err := newMessengerTransactor(ctx, c.ChainClient, valueWei)
	if err != nil {
		return nil, err
	}
	return c.contract.SendToL2Direct(
		auth,
		new(big.Int).SetUint64(l2ChainID),
		message,
		new(big.Int).SetUint64(l2GasLimit),
		new(big.Int).SetUint64(l2GasPerPubdata),
	)
}

func (c *L1MessengerClient) ReceiveFromL2(ctx context.Context, message string, l2BlockNumber, l2MessageIndex uint64, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	auth, err := newMessengerTransactor(ctx, c.ChainClient, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	return c.contract.ReceiveFromL2(
		auth,
		message,
		new(big.Int).SetUint64(l2BlockNumber),
		new(big.Int).SetUint64(l2MessageIndex),
		l2TxNumberInBlock,
		proof,
	)
}

type L2MessengerClient struct {
	*ChainClient
	contract *L2Messenger
	address  common.Address
}

func NewL2MessengerClient(ctx context.Context, rpcURL string, chainID int64, privKeyHex, contractAddr string) (*L2MessengerClient, error) {
	base, err := NewChainClient(ctx, rpcURL, chainID, privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("init chain client: %w", err)
	}

	addr := common.HexToAddress(contractAddr)
	messenger, err := NewL2Messenger(addr, base.Eth)
	if err != nil {
		return nil, fmt.Errorf("bind L2 messenger: %w", err)
	}

	return &L2MessengerClient{
		ChainClient: base,
		contract:    messenger,
		address:     addr,
	}, nil
}

func (c *L2MessengerClient) Address() common.Address {
	return c.address
}

func (c *L2MessengerClient) LastMessageFromL1(ctx context.Context) (string, error) {
	return c.contract.LastMessageFromL1(&bind.CallOpts{Context: ctx})
}

func (c *L2MessengerClient) LastMessageToL1(ctx context.Context) (string, error) {
	return c.contract.LastMessageToL1(&bind.CallOpts{Context: ctx})
}

func (c *L2MessengerClient) SendToL1(ctx context.Context, message string) (*types.Transaction, error) {
	auth, err := newMessengerTransactor(ctx, c.ChainClient, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	return c.contract.SendToL1(auth, message)
}

func newMessengerTransactor(ctx context.Context, chain *ChainClient, valueWei *big.Int) (*bind.TransactOpts, error) {
	nonce, err := chain.Eth.PendingNonceAt(ctx, chain.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice := chain.GasPriceOverride
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice, err = chain.Eth.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}

	auth, err := bind.NewKeyedTransactorWithChainID(chain.PrivKey, chain.ChainID)
	if err != nil {
		return nil, fmt.Errorf("new transactor: %w", err)
	}

	if valueWei == nil {
		valueWei = big.NewInt(0)
	}

	auth.From = chain.From
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = valueWei
	auth.GasPrice = gasPrice
	return auth, nil
}

func (c *L1MessengerClient) estimateBaseCost(ctx context.Context, chainID *big.Int, l2GasLimit, l2GasPerPubdata uint64) (*big.Int, error) {
	mailbox, err := c.contract.Mailbox(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("resolve mailbox address: %w", err)
	}

	gasPrice := c.GasPriceOverride
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice, err = c.Eth.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}

	l2GasLimitBI := new(big.Int).SetUint64(l2GasLimit)
	l2GasPerPubdataBI := new(big.Int).SetUint64(l2GasPerPubdata)

	if chainID != nil {
		if cost, err := callBaseCost(ctx, c.Eth, mailbox, mailboxBaseCostBridgehubABI, chainID, gasPrice, l2GasLimitBI, l2GasPerPubdataBI); err == nil {
			return cost, nil
		}
	}

	cost, err := callBaseCost(ctx, c.Eth, mailbox, mailboxBaseCostABI, gasPrice, l2GasLimitBI, l2GasPerPubdataBI)
	if err != nil {
		return nil, fmt.Errorf("estimate base cost: %w", err)
	}
	return cost, nil
}

func callBaseCost(ctx context.Context, backend bind.ContractBackend, mailbox common.Address, abiJSON string, args ...any) (*big.Int, error) {
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse mailbox ABI: %w", err)
	}
	data, err := parsed.Pack("l2TransactionBaseCost", args...)
	if err != nil {
		return nil, fmt.Errorf("pack base cost call: %w", err)
	}
	msg := ethereum.CallMsg{To: &mailbox, Data: data}
	out, err := backend.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call base cost: %w", err)
	}
	res, err := parsed.Unpack("l2TransactionBaseCost", out)
	if err != nil {
		return nil, fmt.Errorf("unpack base cost: %w", err)
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("unexpected base cost result len=%d", len(res))
	}
	cost, ok := res[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected base cost type %T", res[0])
	}
	return cost, nil
}
