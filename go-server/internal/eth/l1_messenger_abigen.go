// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// L1MessengerMetaData contains all meta data concerning the L1Messenger contract.
var L1MessengerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mailboxAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"l2MessengerAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"l2ChainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"useBridgehub_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"L2_MESSENGER_SYSTEM_CONTRACT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"consumedMessages\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2ChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2Messenger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastMessageFromL2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastMessageToL2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mailbox\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIZkSyncMailbox\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveFromL2\",\"inputs\":[{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2MessageIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2TxNumberInBlock\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"proof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendToL2\",\"inputs\":[{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"l2GasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2GasPerPubdataByteLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendToL2Direct\",\"inputs\":[{\"name\":\"l2ChainIdParam\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"l2GasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2GasPerPubdataByteLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setL2Messenger\",\"inputs\":[{\"name\":\"l2MessengerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"useBridgehub\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"L1ToL2MessageRequested\",\"inputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L2MessengerUpdated\",\"inputs\":[{\"name\":\"newL2Messenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L2ToL1MessageReceived\",\"inputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// L1MessengerABI is the input ABI used to generate the binding from.
// Deprecated: Use L1MessengerMetaData.ABI instead.
var L1MessengerABI = L1MessengerMetaData.ABI

// L1Messenger is an auto generated Go binding around an Ethereum contract.
type L1Messenger struct {
	L1MessengerCaller     // Read-only binding to the contract
	L1MessengerTransactor // Write-only binding to the contract
	L1MessengerFilterer   // Log filterer for contract events
}

// L1MessengerCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1MessengerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1MessengerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1MessengerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1MessengerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1MessengerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1MessengerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1MessengerSession struct {
	Contract     *L1Messenger      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1MessengerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1MessengerCallerSession struct {
	Contract *L1MessengerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// L1MessengerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1MessengerTransactorSession struct {
	Contract     *L1MessengerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// L1MessengerRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1MessengerRaw struct {
	Contract *L1Messenger // Generic contract binding to access the raw methods on
}

// L1MessengerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1MessengerCallerRaw struct {
	Contract *L1MessengerCaller // Generic read-only contract binding to access the raw methods on
}

// L1MessengerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1MessengerTransactorRaw struct {
	Contract *L1MessengerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1Messenger creates a new instance of L1Messenger, bound to a specific deployed contract.
func NewL1Messenger(address common.Address, backend bind.ContractBackend) (*L1Messenger, error) {
	contract, err := bindL1Messenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1Messenger{L1MessengerCaller: L1MessengerCaller{contract: contract}, L1MessengerTransactor: L1MessengerTransactor{contract: contract}, L1MessengerFilterer: L1MessengerFilterer{contract: contract}}, nil
}

// NewL1MessengerCaller creates a new read-only instance of L1Messenger, bound to a specific deployed contract.
func NewL1MessengerCaller(address common.Address, caller bind.ContractCaller) (*L1MessengerCaller, error) {
	contract, err := bindL1Messenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1MessengerCaller{contract: contract}, nil
}

// NewL1MessengerTransactor creates a new write-only instance of L1Messenger, bound to a specific deployed contract.
func NewL1MessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*L1MessengerTransactor, error) {
	contract, err := bindL1Messenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1MessengerTransactor{contract: contract}, nil
}

// NewL1MessengerFilterer creates a new log filterer instance of L1Messenger, bound to a specific deployed contract.
func NewL1MessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*L1MessengerFilterer, error) {
	contract, err := bindL1Messenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1MessengerFilterer{contract: contract}, nil
}

// bindL1Messenger binds a generic wrapper to an already deployed contract.
func bindL1Messenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1MessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Messenger *L1MessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Messenger.Contract.L1MessengerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Messenger *L1MessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Messenger.Contract.L1MessengerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Messenger *L1MessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Messenger.Contract.L1MessengerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Messenger *L1MessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Messenger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Messenger *L1MessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Messenger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Messenger *L1MessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Messenger.Contract.contract.Transact(opts, method, params...)
}

// L2MESSENGERSYSTEMCONTRACT is a free data retrieval call binding the contract method 0xed2364c7.
//
// Solidity: function L2_MESSENGER_SYSTEM_CONTRACT() view returns(address)
func (_L1Messenger *L1MessengerCaller) L2MESSENGERSYSTEMCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "L2_MESSENGER_SYSTEM_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2MESSENGERSYSTEMCONTRACT is a free data retrieval call binding the contract method 0xed2364c7.
//
// Solidity: function L2_MESSENGER_SYSTEM_CONTRACT() view returns(address)
func (_L1Messenger *L1MessengerSession) L2MESSENGERSYSTEMCONTRACT() (common.Address, error) {
	return _L1Messenger.Contract.L2MESSENGERSYSTEMCONTRACT(&_L1Messenger.CallOpts)
}

// L2MESSENGERSYSTEMCONTRACT is a free data retrieval call binding the contract method 0xed2364c7.
//
// Solidity: function L2_MESSENGER_SYSTEM_CONTRACT() view returns(address)
func (_L1Messenger *L1MessengerCallerSession) L2MESSENGERSYSTEMCONTRACT() (common.Address, error) {
	return _L1Messenger.Contract.L2MESSENGERSYSTEMCONTRACT(&_L1Messenger.CallOpts)
}

// ConsumedMessages is a free data retrieval call binding the contract method 0x39bc3c81.
//
// Solidity: function consumedMessages(bytes32 ) view returns(bool)
func (_L1Messenger *L1MessengerCaller) ConsumedMessages(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "consumedMessages", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ConsumedMessages is a free data retrieval call binding the contract method 0x39bc3c81.
//
// Solidity: function consumedMessages(bytes32 ) view returns(bool)
func (_L1Messenger *L1MessengerSession) ConsumedMessages(arg0 [32]byte) (bool, error) {
	return _L1Messenger.Contract.ConsumedMessages(&_L1Messenger.CallOpts, arg0)
}

// ConsumedMessages is a free data retrieval call binding the contract method 0x39bc3c81.
//
// Solidity: function consumedMessages(bytes32 ) view returns(bool)
func (_L1Messenger *L1MessengerCallerSession) ConsumedMessages(arg0 [32]byte) (bool, error) {
	return _L1Messenger.Contract.ConsumedMessages(&_L1Messenger.CallOpts, arg0)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_L1Messenger *L1MessengerCaller) L2ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "l2ChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_L1Messenger *L1MessengerSession) L2ChainId() (*big.Int, error) {
	return _L1Messenger.Contract.L2ChainId(&_L1Messenger.CallOpts)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_L1Messenger *L1MessengerCallerSession) L2ChainId() (*big.Int, error) {
	return _L1Messenger.Contract.L2ChainId(&_L1Messenger.CallOpts)
}

// L2Messenger is a free data retrieval call binding the contract method 0xf5730a72.
//
// Solidity: function l2Messenger() view returns(address)
func (_L1Messenger *L1MessengerCaller) L2Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "l2Messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Messenger is a free data retrieval call binding the contract method 0xf5730a72.
//
// Solidity: function l2Messenger() view returns(address)
func (_L1Messenger *L1MessengerSession) L2Messenger() (common.Address, error) {
	return _L1Messenger.Contract.L2Messenger(&_L1Messenger.CallOpts)
}

// L2Messenger is a free data retrieval call binding the contract method 0xf5730a72.
//
// Solidity: function l2Messenger() view returns(address)
func (_L1Messenger *L1MessengerCallerSession) L2Messenger() (common.Address, error) {
	return _L1Messenger.Contract.L2Messenger(&_L1Messenger.CallOpts)
}

// LastMessageFromL2 is a free data retrieval call binding the contract method 0x537d5421.
//
// Solidity: function lastMessageFromL2() view returns(string)
func (_L1Messenger *L1MessengerCaller) LastMessageFromL2(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "lastMessageFromL2")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// LastMessageFromL2 is a free data retrieval call binding the contract method 0x537d5421.
//
// Solidity: function lastMessageFromL2() view returns(string)
func (_L1Messenger *L1MessengerSession) LastMessageFromL2() (string, error) {
	return _L1Messenger.Contract.LastMessageFromL2(&_L1Messenger.CallOpts)
}

// LastMessageFromL2 is a free data retrieval call binding the contract method 0x537d5421.
//
// Solidity: function lastMessageFromL2() view returns(string)
func (_L1Messenger *L1MessengerCallerSession) LastMessageFromL2() (string, error) {
	return _L1Messenger.Contract.LastMessageFromL2(&_L1Messenger.CallOpts)
}

// LastMessageToL2 is a free data retrieval call binding the contract method 0x2f795bb1.
//
// Solidity: function lastMessageToL2() view returns(string)
func (_L1Messenger *L1MessengerCaller) LastMessageToL2(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "lastMessageToL2")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// LastMessageToL2 is a free data retrieval call binding the contract method 0x2f795bb1.
//
// Solidity: function lastMessageToL2() view returns(string)
func (_L1Messenger *L1MessengerSession) LastMessageToL2() (string, error) {
	return _L1Messenger.Contract.LastMessageToL2(&_L1Messenger.CallOpts)
}

// LastMessageToL2 is a free data retrieval call binding the contract method 0x2f795bb1.
//
// Solidity: function lastMessageToL2() view returns(string)
func (_L1Messenger *L1MessengerCallerSession) LastMessageToL2() (string, error) {
	return _L1Messenger.Contract.LastMessageToL2(&_L1Messenger.CallOpts)
}

// Mailbox is a free data retrieval call binding the contract method 0xd5438eae.
//
// Solidity: function mailbox() view returns(address)
func (_L1Messenger *L1MessengerCaller) Mailbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "mailbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Mailbox is a free data retrieval call binding the contract method 0xd5438eae.
//
// Solidity: function mailbox() view returns(address)
func (_L1Messenger *L1MessengerSession) Mailbox() (common.Address, error) {
	return _L1Messenger.Contract.Mailbox(&_L1Messenger.CallOpts)
}

// Mailbox is a free data retrieval call binding the contract method 0xd5438eae.
//
// Solidity: function mailbox() view returns(address)
func (_L1Messenger *L1MessengerCallerSession) Mailbox() (common.Address, error) {
	return _L1Messenger.Contract.Mailbox(&_L1Messenger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Messenger *L1MessengerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Messenger *L1MessengerSession) Owner() (common.Address, error) {
	return _L1Messenger.Contract.Owner(&_L1Messenger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Messenger *L1MessengerCallerSession) Owner() (common.Address, error) {
	return _L1Messenger.Contract.Owner(&_L1Messenger.CallOpts)
}

// UseBridgehub is a free data retrieval call binding the contract method 0x3cc53b3c.
//
// Solidity: function useBridgehub() view returns(bool)
func (_L1Messenger *L1MessengerCaller) UseBridgehub(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1Messenger.contract.Call(opts, &out, "useBridgehub")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UseBridgehub is a free data retrieval call binding the contract method 0x3cc53b3c.
//
// Solidity: function useBridgehub() view returns(bool)
func (_L1Messenger *L1MessengerSession) UseBridgehub() (bool, error) {
	return _L1Messenger.Contract.UseBridgehub(&_L1Messenger.CallOpts)
}

// UseBridgehub is a free data retrieval call binding the contract method 0x3cc53b3c.
//
// Solidity: function useBridgehub() view returns(bool)
func (_L1Messenger *L1MessengerCallerSession) UseBridgehub() (bool, error) {
	return _L1Messenger.Contract.UseBridgehub(&_L1Messenger.CallOpts)
}

// ReceiveFromL2 is a paid mutator transaction binding the contract method 0x61daaf17.
//
// Solidity: function receiveFromL2(string message, uint256 l2BlockNumber, uint256 l2MessageIndex, uint16 l2TxNumberInBlock, bytes32[] proof) returns()
func (_L1Messenger *L1MessengerTransactor) ReceiveFromL2(opts *bind.TransactOpts, message string, l2BlockNumber *big.Int, l2MessageIndex *big.Int, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	return _L1Messenger.contract.Transact(opts, "receiveFromL2", message, l2BlockNumber, l2MessageIndex, l2TxNumberInBlock, proof)
}

// ReceiveFromL2 is a paid mutator transaction binding the contract method 0x61daaf17.
//
// Solidity: function receiveFromL2(string message, uint256 l2BlockNumber, uint256 l2MessageIndex, uint16 l2TxNumberInBlock, bytes32[] proof) returns()
func (_L1Messenger *L1MessengerSession) ReceiveFromL2(message string, l2BlockNumber *big.Int, l2MessageIndex *big.Int, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	return _L1Messenger.Contract.ReceiveFromL2(&_L1Messenger.TransactOpts, message, l2BlockNumber, l2MessageIndex, l2TxNumberInBlock, proof)
}

// ReceiveFromL2 is a paid mutator transaction binding the contract method 0x61daaf17.
//
// Solidity: function receiveFromL2(string message, uint256 l2BlockNumber, uint256 l2MessageIndex, uint16 l2TxNumberInBlock, bytes32[] proof) returns()
func (_L1Messenger *L1MessengerTransactorSession) ReceiveFromL2(message string, l2BlockNumber *big.Int, l2MessageIndex *big.Int, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	return _L1Messenger.Contract.ReceiveFromL2(&_L1Messenger.TransactOpts, message, l2BlockNumber, l2MessageIndex, l2TxNumberInBlock, proof)
}

// SendToL2 is a paid mutator transaction binding the contract method 0x7b4cc8e9.
//
// Solidity: function sendToL2(string message, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit) payable returns(bytes32 txHash)
func (_L1Messenger *L1MessengerTransactor) SendToL2(opts *bind.TransactOpts, message string, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int) (*types.Transaction, error) {
	return _L1Messenger.contract.Transact(opts, "sendToL2", message, l2GasLimit, l2GasPerPubdataByteLimit)
}

// SendToL2 is a paid mutator transaction binding the contract method 0x7b4cc8e9.
//
// Solidity: function sendToL2(string message, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit) payable returns(bytes32 txHash)
func (_L1Messenger *L1MessengerSession) SendToL2(message string, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int) (*types.Transaction, error) {
	return _L1Messenger.Contract.SendToL2(&_L1Messenger.TransactOpts, message, l2GasLimit, l2GasPerPubdataByteLimit)
}

// SendToL2 is a paid mutator transaction binding the contract method 0x7b4cc8e9.
//
// Solidity: function sendToL2(string message, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit) payable returns(bytes32 txHash)
func (_L1Messenger *L1MessengerTransactorSession) SendToL2(message string, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int) (*types.Transaction, error) {
	return _L1Messenger.Contract.SendToL2(&_L1Messenger.TransactOpts, message, l2GasLimit, l2GasPerPubdataByteLimit)
}

// SendToL2Direct is a paid mutator transaction binding the contract method 0x12587039.
//
// Solidity: function sendToL2Direct(uint256 l2ChainIdParam, string message, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit) payable returns(bytes32 txHash)
func (_L1Messenger *L1MessengerTransactor) SendToL2Direct(opts *bind.TransactOpts, l2ChainIdParam *big.Int, message string, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int) (*types.Transaction, error) {
	return _L1Messenger.contract.Transact(opts, "sendToL2Direct", l2ChainIdParam, message, l2GasLimit, l2GasPerPubdataByteLimit)
}

// SendToL2Direct is a paid mutator transaction binding the contract method 0x12587039.
//
// Solidity: function sendToL2Direct(uint256 l2ChainIdParam, string message, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit) payable returns(bytes32 txHash)
func (_L1Messenger *L1MessengerSession) SendToL2Direct(l2ChainIdParam *big.Int, message string, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int) (*types.Transaction, error) {
	return _L1Messenger.Contract.SendToL2Direct(&_L1Messenger.TransactOpts, l2ChainIdParam, message, l2GasLimit, l2GasPerPubdataByteLimit)
}

// SendToL2Direct is a paid mutator transaction binding the contract method 0x12587039.
//
// Solidity: function sendToL2Direct(uint256 l2ChainIdParam, string message, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit) payable returns(bytes32 txHash)
func (_L1Messenger *L1MessengerTransactorSession) SendToL2Direct(l2ChainIdParam *big.Int, message string, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int) (*types.Transaction, error) {
	return _L1Messenger.Contract.SendToL2Direct(&_L1Messenger.TransactOpts, l2ChainIdParam, message, l2GasLimit, l2GasPerPubdataByteLimit)
}

// SetL2Messenger is a paid mutator transaction binding the contract method 0x17cb75cb.
//
// Solidity: function setL2Messenger(address l2MessengerAddress) returns()
func (_L1Messenger *L1MessengerTransactor) SetL2Messenger(opts *bind.TransactOpts, l2MessengerAddress common.Address) (*types.Transaction, error) {
	return _L1Messenger.contract.Transact(opts, "setL2Messenger", l2MessengerAddress)
}

// SetL2Messenger is a paid mutator transaction binding the contract method 0x17cb75cb.
//
// Solidity: function setL2Messenger(address l2MessengerAddress) returns()
func (_L1Messenger *L1MessengerSession) SetL2Messenger(l2MessengerAddress common.Address) (*types.Transaction, error) {
	return _L1Messenger.Contract.SetL2Messenger(&_L1Messenger.TransactOpts, l2MessengerAddress)
}

// SetL2Messenger is a paid mutator transaction binding the contract method 0x17cb75cb.
//
// Solidity: function setL2Messenger(address l2MessengerAddress) returns()
func (_L1Messenger *L1MessengerTransactorSession) SetL2Messenger(l2MessengerAddress common.Address) (*types.Transaction, error) {
	return _L1Messenger.Contract.SetL2Messenger(&_L1Messenger.TransactOpts, l2MessengerAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Messenger *L1MessengerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L1Messenger.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Messenger *L1MessengerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Messenger.Contract.TransferOwnership(&_L1Messenger.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Messenger *L1MessengerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Messenger.Contract.TransferOwnership(&_L1Messenger.TransactOpts, newOwner)
}

// L1MessengerL1ToL2MessageRequestedIterator is returned from FilterL1ToL2MessageRequested and is used to iterate over the raw logs and unpacked data for L1ToL2MessageRequested events raised by the L1Messenger contract.
type L1MessengerL1ToL2MessageRequestedIterator struct {
	Event *L1MessengerL1ToL2MessageRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *L1MessengerL1ToL2MessageRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1MessengerL1ToL2MessageRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(L1MessengerL1ToL2MessageRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *L1MessengerL1ToL2MessageRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1MessengerL1ToL2MessageRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1MessengerL1ToL2MessageRequested represents a L1ToL2MessageRequested event raised by the L1Messenger contract.
type L1MessengerL1ToL2MessageRequested struct {
	TxHash  [32]byte
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterL1ToL2MessageRequested is a free log retrieval operation binding the contract event 0x8c4680c5ba24908229e7b8e78e02cd8ba9fdb737edb1551fc52c6c6c61000953.
//
// Solidity: event L1ToL2MessageRequested(bytes32 indexed txHash, string message)
func (_L1Messenger *L1MessengerFilterer) FilterL1ToL2MessageRequested(opts *bind.FilterOpts, txHash [][32]byte) (*L1MessengerL1ToL2MessageRequestedIterator, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}

	logs, sub, err := _L1Messenger.contract.FilterLogs(opts, "L1ToL2MessageRequested", txHashRule)
	if err != nil {
		return nil, err
	}
	return &L1MessengerL1ToL2MessageRequestedIterator{contract: _L1Messenger.contract, event: "L1ToL2MessageRequested", logs: logs, sub: sub}, nil
}

// WatchL1ToL2MessageRequested is a free log subscription operation binding the contract event 0x8c4680c5ba24908229e7b8e78e02cd8ba9fdb737edb1551fc52c6c6c61000953.
//
// Solidity: event L1ToL2MessageRequested(bytes32 indexed txHash, string message)
func (_L1Messenger *L1MessengerFilterer) WatchL1ToL2MessageRequested(opts *bind.WatchOpts, sink chan<- *L1MessengerL1ToL2MessageRequested, txHash [][32]byte) (event.Subscription, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}

	logs, sub, err := _L1Messenger.contract.WatchLogs(opts, "L1ToL2MessageRequested", txHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1MessengerL1ToL2MessageRequested)
				if err := _L1Messenger.contract.UnpackLog(event, "L1ToL2MessageRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseL1ToL2MessageRequested is a log parse operation binding the contract event 0x8c4680c5ba24908229e7b8e78e02cd8ba9fdb737edb1551fc52c6c6c61000953.
//
// Solidity: event L1ToL2MessageRequested(bytes32 indexed txHash, string message)
func (_L1Messenger *L1MessengerFilterer) ParseL1ToL2MessageRequested(log types.Log) (*L1MessengerL1ToL2MessageRequested, error) {
	event := new(L1MessengerL1ToL2MessageRequested)
	if err := _L1Messenger.contract.UnpackLog(event, "L1ToL2MessageRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1MessengerL2MessengerUpdatedIterator is returned from FilterL2MessengerUpdated and is used to iterate over the raw logs and unpacked data for L2MessengerUpdated events raised by the L1Messenger contract.
type L1MessengerL2MessengerUpdatedIterator struct {
	Event *L1MessengerL2MessengerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *L1MessengerL2MessengerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1MessengerL2MessengerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(L1MessengerL2MessengerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *L1MessengerL2MessengerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1MessengerL2MessengerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1MessengerL2MessengerUpdated represents a L2MessengerUpdated event raised by the L1Messenger contract.
type L1MessengerL2MessengerUpdated struct {
	NewL2Messenger common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterL2MessengerUpdated is a free log retrieval operation binding the contract event 0xcc03cd3ae70f712db44504c94bdbf0b3b14fc41d55bd236a8240ae85073ef57f.
//
// Solidity: event L2MessengerUpdated(address indexed newL2Messenger)
func (_L1Messenger *L1MessengerFilterer) FilterL2MessengerUpdated(opts *bind.FilterOpts, newL2Messenger []common.Address) (*L1MessengerL2MessengerUpdatedIterator, error) {

	var newL2MessengerRule []interface{}
	for _, newL2MessengerItem := range newL2Messenger {
		newL2MessengerRule = append(newL2MessengerRule, newL2MessengerItem)
	}

	logs, sub, err := _L1Messenger.contract.FilterLogs(opts, "L2MessengerUpdated", newL2MessengerRule)
	if err != nil {
		return nil, err
	}
	return &L1MessengerL2MessengerUpdatedIterator{contract: _L1Messenger.contract, event: "L2MessengerUpdated", logs: logs, sub: sub}, nil
}

// WatchL2MessengerUpdated is a free log subscription operation binding the contract event 0xcc03cd3ae70f712db44504c94bdbf0b3b14fc41d55bd236a8240ae85073ef57f.
//
// Solidity: event L2MessengerUpdated(address indexed newL2Messenger)
func (_L1Messenger *L1MessengerFilterer) WatchL2MessengerUpdated(opts *bind.WatchOpts, sink chan<- *L1MessengerL2MessengerUpdated, newL2Messenger []common.Address) (event.Subscription, error) {

	var newL2MessengerRule []interface{}
	for _, newL2MessengerItem := range newL2Messenger {
		newL2MessengerRule = append(newL2MessengerRule, newL2MessengerItem)
	}

	logs, sub, err := _L1Messenger.contract.WatchLogs(opts, "L2MessengerUpdated", newL2MessengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1MessengerL2MessengerUpdated)
				if err := _L1Messenger.contract.UnpackLog(event, "L2MessengerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseL2MessengerUpdated is a log parse operation binding the contract event 0xcc03cd3ae70f712db44504c94bdbf0b3b14fc41d55bd236a8240ae85073ef57f.
//
// Solidity: event L2MessengerUpdated(address indexed newL2Messenger)
func (_L1Messenger *L1MessengerFilterer) ParseL2MessengerUpdated(log types.Log) (*L1MessengerL2MessengerUpdated, error) {
	event := new(L1MessengerL2MessengerUpdated)
	if err := _L1Messenger.contract.UnpackLog(event, "L2MessengerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1MessengerL2ToL1MessageReceivedIterator is returned from FilterL2ToL1MessageReceived and is used to iterate over the raw logs and unpacked data for L2ToL1MessageReceived events raised by the L1Messenger contract.
type L1MessengerL2ToL1MessageReceivedIterator struct {
	Event *L1MessengerL2ToL1MessageReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *L1MessengerL2ToL1MessageReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1MessengerL2ToL1MessageReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(L1MessengerL2ToL1MessageReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *L1MessengerL2ToL1MessageReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1MessengerL2ToL1MessageReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1MessengerL2ToL1MessageReceived represents a L2ToL1MessageReceived event raised by the L1Messenger contract.
type L1MessengerL2ToL1MessageReceived struct {
	MessageHash [32]byte
	Message     string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterL2ToL1MessageReceived is a free log retrieval operation binding the contract event 0xb3c20656c6722a2b39018d45c4b0887cd1ef0a41088739c83b6d9093f8219daf.
//
// Solidity: event L2ToL1MessageReceived(bytes32 indexed messageHash, string message)
func (_L1Messenger *L1MessengerFilterer) FilterL2ToL1MessageReceived(opts *bind.FilterOpts, messageHash [][32]byte) (*L1MessengerL2ToL1MessageReceivedIterator, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _L1Messenger.contract.FilterLogs(opts, "L2ToL1MessageReceived", messageHashRule)
	if err != nil {
		return nil, err
	}
	return &L1MessengerL2ToL1MessageReceivedIterator{contract: _L1Messenger.contract, event: "L2ToL1MessageReceived", logs: logs, sub: sub}, nil
}

// WatchL2ToL1MessageReceived is a free log subscription operation binding the contract event 0xb3c20656c6722a2b39018d45c4b0887cd1ef0a41088739c83b6d9093f8219daf.
//
// Solidity: event L2ToL1MessageReceived(bytes32 indexed messageHash, string message)
func (_L1Messenger *L1MessengerFilterer) WatchL2ToL1MessageReceived(opts *bind.WatchOpts, sink chan<- *L1MessengerL2ToL1MessageReceived, messageHash [][32]byte) (event.Subscription, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _L1Messenger.contract.WatchLogs(opts, "L2ToL1MessageReceived", messageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1MessengerL2ToL1MessageReceived)
				if err := _L1Messenger.contract.UnpackLog(event, "L2ToL1MessageReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseL2ToL1MessageReceived is a log parse operation binding the contract event 0xb3c20656c6722a2b39018d45c4b0887cd1ef0a41088739c83b6d9093f8219daf.
//
// Solidity: event L2ToL1MessageReceived(bytes32 indexed messageHash, string message)
func (_L1Messenger *L1MessengerFilterer) ParseL2ToL1MessageReceived(log types.Log) (*L1MessengerL2ToL1MessageReceived, error) {
	event := new(L1MessengerL2ToL1MessageReceived)
	if err := _L1Messenger.contract.UnpackLog(event, "L2ToL1MessageReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1MessengerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L1Messenger contract.
type L1MessengerOwnershipTransferredIterator struct {
	Event *L1MessengerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *L1MessengerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1MessengerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(L1MessengerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *L1MessengerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1MessengerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1MessengerOwnershipTransferred represents a OwnershipTransferred event raised by the L1Messenger contract.
type L1MessengerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Messenger *L1MessengerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1MessengerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Messenger.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1MessengerOwnershipTransferredIterator{contract: _L1Messenger.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Messenger *L1MessengerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L1MessengerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Messenger.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1MessengerOwnershipTransferred)
				if err := _L1Messenger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Messenger *L1MessengerFilterer) ParseOwnershipTransferred(log types.Log) (*L1MessengerOwnershipTransferred, error) {
	event := new(L1MessengerOwnershipTransferred)
	if err := _L1Messenger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
