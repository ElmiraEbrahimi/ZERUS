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

// L2MessengerMetaData contains all meta data concerning the L2Messenger contract.
var L2MessengerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"l1MessengerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"l1Messenger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastMessageFromL1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastMessageToL1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveFromL1\",\"inputs\":[{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendToL1\",\"inputs\":[{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setL1Messenger\",\"inputs\":[{\"name\":\"l1MessengerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"L1MessengerUpdated\",\"inputs\":[{\"name\":\"newL1Messenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L1ToL2MessageReceived\",\"inputs\":[{\"name\":\"message\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L2ToL1MessageSent\",\"inputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// L2MessengerABI is the input ABI used to generate the binding from.
// Deprecated: Use L2MessengerMetaData.ABI instead.
var L2MessengerABI = L2MessengerMetaData.ABI

// L2Messenger is an auto generated Go binding around an Ethereum contract.
type L2Messenger struct {
	L2MessengerCaller     // Read-only binding to the contract
	L2MessengerTransactor // Write-only binding to the contract
	L2MessengerFilterer   // Log filterer for contract events
}

// L2MessengerCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2MessengerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2MessengerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2MessengerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2MessengerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2MessengerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2MessengerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2MessengerSession struct {
	Contract     *L2Messenger      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L2MessengerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2MessengerCallerSession struct {
	Contract *L2MessengerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// L2MessengerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2MessengerTransactorSession struct {
	Contract     *L2MessengerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// L2MessengerRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2MessengerRaw struct {
	Contract *L2Messenger // Generic contract binding to access the raw methods on
}

// L2MessengerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2MessengerCallerRaw struct {
	Contract *L2MessengerCaller // Generic read-only contract binding to access the raw methods on
}

// L2MessengerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2MessengerTransactorRaw struct {
	Contract *L2MessengerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2Messenger creates a new instance of L2Messenger, bound to a specific deployed contract.
func NewL2Messenger(address common.Address, backend bind.ContractBackend) (*L2Messenger, error) {
	contract, err := bindL2Messenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2Messenger{L2MessengerCaller: L2MessengerCaller{contract: contract}, L2MessengerTransactor: L2MessengerTransactor{contract: contract}, L2MessengerFilterer: L2MessengerFilterer{contract: contract}}, nil
}

// NewL2MessengerCaller creates a new read-only instance of L2Messenger, bound to a specific deployed contract.
func NewL2MessengerCaller(address common.Address, caller bind.ContractCaller) (*L2MessengerCaller, error) {
	contract, err := bindL2Messenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2MessengerCaller{contract: contract}, nil
}

// NewL2MessengerTransactor creates a new write-only instance of L2Messenger, bound to a specific deployed contract.
func NewL2MessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*L2MessengerTransactor, error) {
	contract, err := bindL2Messenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2MessengerTransactor{contract: contract}, nil
}

// NewL2MessengerFilterer creates a new log filterer instance of L2Messenger, bound to a specific deployed contract.
func NewL2MessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*L2MessengerFilterer, error) {
	contract, err := bindL2Messenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2MessengerFilterer{contract: contract}, nil
}

// bindL2Messenger binds a generic wrapper to an already deployed contract.
func bindL2Messenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2MessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2Messenger *L2MessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2Messenger.Contract.L2MessengerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2Messenger *L2MessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2Messenger.Contract.L2MessengerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2Messenger *L2MessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2Messenger.Contract.L2MessengerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2Messenger *L2MessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2Messenger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2Messenger *L2MessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2Messenger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2Messenger *L2MessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2Messenger.Contract.contract.Transact(opts, method, params...)
}

// L1Messenger is a free data retrieval call binding the contract method 0x6140e0e6.
//
// Solidity: function l1Messenger() view returns(address)
func (_L2Messenger *L2MessengerCaller) L1Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2Messenger.contract.Call(opts, &out, "l1Messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Messenger is a free data retrieval call binding the contract method 0x6140e0e6.
//
// Solidity: function l1Messenger() view returns(address)
func (_L2Messenger *L2MessengerSession) L1Messenger() (common.Address, error) {
	return _L2Messenger.Contract.L1Messenger(&_L2Messenger.CallOpts)
}

// L1Messenger is a free data retrieval call binding the contract method 0x6140e0e6.
//
// Solidity: function l1Messenger() view returns(address)
func (_L2Messenger *L2MessengerCallerSession) L1Messenger() (common.Address, error) {
	return _L2Messenger.Contract.L1Messenger(&_L2Messenger.CallOpts)
}

// LastMessageFromL1 is a free data retrieval call binding the contract method 0x1e0617d7.
//
// Solidity: function lastMessageFromL1() view returns(string)
func (_L2Messenger *L2MessengerCaller) LastMessageFromL1(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L2Messenger.contract.Call(opts, &out, "lastMessageFromL1")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// LastMessageFromL1 is a free data retrieval call binding the contract method 0x1e0617d7.
//
// Solidity: function lastMessageFromL1() view returns(string)
func (_L2Messenger *L2MessengerSession) LastMessageFromL1() (string, error) {
	return _L2Messenger.Contract.LastMessageFromL1(&_L2Messenger.CallOpts)
}

// LastMessageFromL1 is a free data retrieval call binding the contract method 0x1e0617d7.
//
// Solidity: function lastMessageFromL1() view returns(string)
func (_L2Messenger *L2MessengerCallerSession) LastMessageFromL1() (string, error) {
	return _L2Messenger.Contract.LastMessageFromL1(&_L2Messenger.CallOpts)
}

// LastMessageToL1 is a free data retrieval call binding the contract method 0x5fa2534c.
//
// Solidity: function lastMessageToL1() view returns(string)
func (_L2Messenger *L2MessengerCaller) LastMessageToL1(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L2Messenger.contract.Call(opts, &out, "lastMessageToL1")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// LastMessageToL1 is a free data retrieval call binding the contract method 0x5fa2534c.
//
// Solidity: function lastMessageToL1() view returns(string)
func (_L2Messenger *L2MessengerSession) LastMessageToL1() (string, error) {
	return _L2Messenger.Contract.LastMessageToL1(&_L2Messenger.CallOpts)
}

// LastMessageToL1 is a free data retrieval call binding the contract method 0x5fa2534c.
//
// Solidity: function lastMessageToL1() view returns(string)
func (_L2Messenger *L2MessengerCallerSession) LastMessageToL1() (string, error) {
	return _L2Messenger.Contract.LastMessageToL1(&_L2Messenger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2Messenger *L2MessengerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2Messenger.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2Messenger *L2MessengerSession) Owner() (common.Address, error) {
	return _L2Messenger.Contract.Owner(&_L2Messenger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2Messenger *L2MessengerCallerSession) Owner() (common.Address, error) {
	return _L2Messenger.Contract.Owner(&_L2Messenger.CallOpts)
}

// ReceiveFromL1 is a paid mutator transaction binding the contract method 0xeef00d73.
//
// Solidity: function receiveFromL1(string message) returns()
func (_L2Messenger *L2MessengerTransactor) ReceiveFromL1(opts *bind.TransactOpts, message string) (*types.Transaction, error) {
	return _L2Messenger.contract.Transact(opts, "receiveFromL1", message)
}

// ReceiveFromL1 is a paid mutator transaction binding the contract method 0xeef00d73.
//
// Solidity: function receiveFromL1(string message) returns()
func (_L2Messenger *L2MessengerSession) ReceiveFromL1(message string) (*types.Transaction, error) {
	return _L2Messenger.Contract.ReceiveFromL1(&_L2Messenger.TransactOpts, message)
}

// ReceiveFromL1 is a paid mutator transaction binding the contract method 0xeef00d73.
//
// Solidity: function receiveFromL1(string message) returns()
func (_L2Messenger *L2MessengerTransactorSession) ReceiveFromL1(message string) (*types.Transaction, error) {
	return _L2Messenger.Contract.ReceiveFromL1(&_L2Messenger.TransactOpts, message)
}

// SendToL1 is a paid mutator transaction binding the contract method 0xaa46d7de.
//
// Solidity: function sendToL1(string message) returns(bytes32 messageHash)
func (_L2Messenger *L2MessengerTransactor) SendToL1(opts *bind.TransactOpts, message string) (*types.Transaction, error) {
	return _L2Messenger.contract.Transact(opts, "sendToL1", message)
}

// SendToL1 is a paid mutator transaction binding the contract method 0xaa46d7de.
//
// Solidity: function sendToL1(string message) returns(bytes32 messageHash)
func (_L2Messenger *L2MessengerSession) SendToL1(message string) (*types.Transaction, error) {
	return _L2Messenger.Contract.SendToL1(&_L2Messenger.TransactOpts, message)
}

// SendToL1 is a paid mutator transaction binding the contract method 0xaa46d7de.
//
// Solidity: function sendToL1(string message) returns(bytes32 messageHash)
func (_L2Messenger *L2MessengerTransactorSession) SendToL1(message string) (*types.Transaction, error) {
	return _L2Messenger.Contract.SendToL1(&_L2Messenger.TransactOpts, message)
}

// SetL1Messenger is a paid mutator transaction binding the contract method 0x5a860897.
//
// Solidity: function setL1Messenger(address l1MessengerAddress) returns()
func (_L2Messenger *L2MessengerTransactor) SetL1Messenger(opts *bind.TransactOpts, l1MessengerAddress common.Address) (*types.Transaction, error) {
	return _L2Messenger.contract.Transact(opts, "setL1Messenger", l1MessengerAddress)
}

// SetL1Messenger is a paid mutator transaction binding the contract method 0x5a860897.
//
// Solidity: function setL1Messenger(address l1MessengerAddress) returns()
func (_L2Messenger *L2MessengerSession) SetL1Messenger(l1MessengerAddress common.Address) (*types.Transaction, error) {
	return _L2Messenger.Contract.SetL1Messenger(&_L2Messenger.TransactOpts, l1MessengerAddress)
}

// SetL1Messenger is a paid mutator transaction binding the contract method 0x5a860897.
//
// Solidity: function setL1Messenger(address l1MessengerAddress) returns()
func (_L2Messenger *L2MessengerTransactorSession) SetL1Messenger(l1MessengerAddress common.Address) (*types.Transaction, error) {
	return _L2Messenger.Contract.SetL1Messenger(&_L2Messenger.TransactOpts, l1MessengerAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2Messenger *L2MessengerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L2Messenger.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2Messenger *L2MessengerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2Messenger.Contract.TransferOwnership(&_L2Messenger.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2Messenger *L2MessengerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2Messenger.Contract.TransferOwnership(&_L2Messenger.TransactOpts, newOwner)
}

// L2MessengerL1MessengerUpdatedIterator is returned from FilterL1MessengerUpdated and is used to iterate over the raw logs and unpacked data for L1MessengerUpdated events raised by the L2Messenger contract.
type L2MessengerL1MessengerUpdatedIterator struct {
	Event *L2MessengerL1MessengerUpdated // Event containing the contract specifics and raw log

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
func (it *L2MessengerL1MessengerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2MessengerL1MessengerUpdated)
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
		it.Event = new(L2MessengerL1MessengerUpdated)
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
func (it *L2MessengerL1MessengerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2MessengerL1MessengerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2MessengerL1MessengerUpdated represents a L1MessengerUpdated event raised by the L2Messenger contract.
type L2MessengerL1MessengerUpdated struct {
	NewL1Messenger common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterL1MessengerUpdated is a free log retrieval operation binding the contract event 0xa9504df670c53f044496433789d10797c347d80e0d74c0c2492405e5be676175.
//
// Solidity: event L1MessengerUpdated(address indexed newL1Messenger)
func (_L2Messenger *L2MessengerFilterer) FilterL1MessengerUpdated(opts *bind.FilterOpts, newL1Messenger []common.Address) (*L2MessengerL1MessengerUpdatedIterator, error) {

	var newL1MessengerRule []interface{}
	for _, newL1MessengerItem := range newL1Messenger {
		newL1MessengerRule = append(newL1MessengerRule, newL1MessengerItem)
	}

	logs, sub, err := _L2Messenger.contract.FilterLogs(opts, "L1MessengerUpdated", newL1MessengerRule)
	if err != nil {
		return nil, err
	}
	return &L2MessengerL1MessengerUpdatedIterator{contract: _L2Messenger.contract, event: "L1MessengerUpdated", logs: logs, sub: sub}, nil
}

// WatchL1MessengerUpdated is a free log subscription operation binding the contract event 0xa9504df670c53f044496433789d10797c347d80e0d74c0c2492405e5be676175.
//
// Solidity: event L1MessengerUpdated(address indexed newL1Messenger)
func (_L2Messenger *L2MessengerFilterer) WatchL1MessengerUpdated(opts *bind.WatchOpts, sink chan<- *L2MessengerL1MessengerUpdated, newL1Messenger []common.Address) (event.Subscription, error) {

	var newL1MessengerRule []interface{}
	for _, newL1MessengerItem := range newL1Messenger {
		newL1MessengerRule = append(newL1MessengerRule, newL1MessengerItem)
	}

	logs, sub, err := _L2Messenger.contract.WatchLogs(opts, "L1MessengerUpdated", newL1MessengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2MessengerL1MessengerUpdated)
				if err := _L2Messenger.contract.UnpackLog(event, "L1MessengerUpdated", log); err != nil {
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

// ParseL1MessengerUpdated is a log parse operation binding the contract event 0xa9504df670c53f044496433789d10797c347d80e0d74c0c2492405e5be676175.
//
// Solidity: event L1MessengerUpdated(address indexed newL1Messenger)
func (_L2Messenger *L2MessengerFilterer) ParseL1MessengerUpdated(log types.Log) (*L2MessengerL1MessengerUpdated, error) {
	event := new(L2MessengerL1MessengerUpdated)
	if err := _L2Messenger.contract.UnpackLog(event, "L1MessengerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2MessengerL1ToL2MessageReceivedIterator is returned from FilterL1ToL2MessageReceived and is used to iterate over the raw logs and unpacked data for L1ToL2MessageReceived events raised by the L2Messenger contract.
type L2MessengerL1ToL2MessageReceivedIterator struct {
	Event *L2MessengerL1ToL2MessageReceived // Event containing the contract specifics and raw log

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
func (it *L2MessengerL1ToL2MessageReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2MessengerL1ToL2MessageReceived)
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
		it.Event = new(L2MessengerL1ToL2MessageReceived)
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
func (it *L2MessengerL1ToL2MessageReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2MessengerL1ToL2MessageReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2MessengerL1ToL2MessageReceived represents a L1ToL2MessageReceived event raised by the L2Messenger contract.
type L2MessengerL1ToL2MessageReceived struct {
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterL1ToL2MessageReceived is a free log retrieval operation binding the contract event 0xd699719d7f28484e6107073cf4b32ec9f7873b2a62e4cfb81f5030de4c223751.
//
// Solidity: event L1ToL2MessageReceived(string message)
func (_L2Messenger *L2MessengerFilterer) FilterL1ToL2MessageReceived(opts *bind.FilterOpts) (*L2MessengerL1ToL2MessageReceivedIterator, error) {

	logs, sub, err := _L2Messenger.contract.FilterLogs(opts, "L1ToL2MessageReceived")
	if err != nil {
		return nil, err
	}
	return &L2MessengerL1ToL2MessageReceivedIterator{contract: _L2Messenger.contract, event: "L1ToL2MessageReceived", logs: logs, sub: sub}, nil
}

// WatchL1ToL2MessageReceived is a free log subscription operation binding the contract event 0xd699719d7f28484e6107073cf4b32ec9f7873b2a62e4cfb81f5030de4c223751.
//
// Solidity: event L1ToL2MessageReceived(string message)
func (_L2Messenger *L2MessengerFilterer) WatchL1ToL2MessageReceived(opts *bind.WatchOpts, sink chan<- *L2MessengerL1ToL2MessageReceived) (event.Subscription, error) {

	logs, sub, err := _L2Messenger.contract.WatchLogs(opts, "L1ToL2MessageReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2MessengerL1ToL2MessageReceived)
				if err := _L2Messenger.contract.UnpackLog(event, "L1ToL2MessageReceived", log); err != nil {
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

// ParseL1ToL2MessageReceived is a log parse operation binding the contract event 0xd699719d7f28484e6107073cf4b32ec9f7873b2a62e4cfb81f5030de4c223751.
//
// Solidity: event L1ToL2MessageReceived(string message)
func (_L2Messenger *L2MessengerFilterer) ParseL1ToL2MessageReceived(log types.Log) (*L2MessengerL1ToL2MessageReceived, error) {
	event := new(L2MessengerL1ToL2MessageReceived)
	if err := _L2Messenger.contract.UnpackLog(event, "L1ToL2MessageReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2MessengerL2ToL1MessageSentIterator is returned from FilterL2ToL1MessageSent and is used to iterate over the raw logs and unpacked data for L2ToL1MessageSent events raised by the L2Messenger contract.
type L2MessengerL2ToL1MessageSentIterator struct {
	Event *L2MessengerL2ToL1MessageSent // Event containing the contract specifics and raw log

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
func (it *L2MessengerL2ToL1MessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2MessengerL2ToL1MessageSent)
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
		it.Event = new(L2MessengerL2ToL1MessageSent)
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
func (it *L2MessengerL2ToL1MessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2MessengerL2ToL1MessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2MessengerL2ToL1MessageSent represents a L2ToL1MessageSent event raised by the L2Messenger contract.
type L2MessengerL2ToL1MessageSent struct {
	MessageHash [32]byte
	Message     string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterL2ToL1MessageSent is a free log retrieval operation binding the contract event 0x3b92c4ff008d613401d7303fe637e0803a07d143cc0d43cf361183e4d2a5e25b.
//
// Solidity: event L2ToL1MessageSent(bytes32 indexed messageHash, string message)
func (_L2Messenger *L2MessengerFilterer) FilterL2ToL1MessageSent(opts *bind.FilterOpts, messageHash [][32]byte) (*L2MessengerL2ToL1MessageSentIterator, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _L2Messenger.contract.FilterLogs(opts, "L2ToL1MessageSent", messageHashRule)
	if err != nil {
		return nil, err
	}
	return &L2MessengerL2ToL1MessageSentIterator{contract: _L2Messenger.contract, event: "L2ToL1MessageSent", logs: logs, sub: sub}, nil
}

// WatchL2ToL1MessageSent is a free log subscription operation binding the contract event 0x3b92c4ff008d613401d7303fe637e0803a07d143cc0d43cf361183e4d2a5e25b.
//
// Solidity: event L2ToL1MessageSent(bytes32 indexed messageHash, string message)
func (_L2Messenger *L2MessengerFilterer) WatchL2ToL1MessageSent(opts *bind.WatchOpts, sink chan<- *L2MessengerL2ToL1MessageSent, messageHash [][32]byte) (event.Subscription, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _L2Messenger.contract.WatchLogs(opts, "L2ToL1MessageSent", messageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2MessengerL2ToL1MessageSent)
				if err := _L2Messenger.contract.UnpackLog(event, "L2ToL1MessageSent", log); err != nil {
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

// ParseL2ToL1MessageSent is a log parse operation binding the contract event 0x3b92c4ff008d613401d7303fe637e0803a07d143cc0d43cf361183e4d2a5e25b.
//
// Solidity: event L2ToL1MessageSent(bytes32 indexed messageHash, string message)
func (_L2Messenger *L2MessengerFilterer) ParseL2ToL1MessageSent(log types.Log) (*L2MessengerL2ToL1MessageSent, error) {
	event := new(L2MessengerL2ToL1MessageSent)
	if err := _L2Messenger.contract.UnpackLog(event, "L2ToL1MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2MessengerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L2Messenger contract.
type L2MessengerOwnershipTransferredIterator struct {
	Event *L2MessengerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L2MessengerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2MessengerOwnershipTransferred)
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
		it.Event = new(L2MessengerOwnershipTransferred)
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
func (it *L2MessengerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2MessengerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2MessengerOwnershipTransferred represents a OwnershipTransferred event raised by the L2Messenger contract.
type L2MessengerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2Messenger *L2MessengerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L2MessengerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2Messenger.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L2MessengerOwnershipTransferredIterator{contract: _L2Messenger.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2Messenger *L2MessengerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L2MessengerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2Messenger.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2MessengerOwnershipTransferred)
				if err := _L2Messenger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L2Messenger *L2MessengerFilterer) ParseOwnershipTransferred(log types.Log) (*L2MessengerOwnershipTransferred, error) {
	event := new(L2MessengerOwnershipTransferred)
	if err := _L2Messenger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
