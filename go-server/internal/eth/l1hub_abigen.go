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

// L1HubPublicKey is an auto generated low-level Go binding around an user-defined struct.
type L1HubPublicKey struct {
	X *big.Int
	Y *big.Int
}

// L1HubReplacementParams is an auto generated low-level Go binding around an user-defined struct.
type L1HubReplacementParams struct {
	TargetValidatorID    *big.Int
	CandidateValidatorID *big.Int
	CandidatePubKey      L1HubPublicKey
	CandidateStake       *big.Int
	TargetLeafIndex      *big.Int
	Path                 []*big.Int
	Depth                *big.Int
}

// L1HubMetaData contains all meta data concerning the L1Hub contract.
var L1HubMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mailboxAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"l2OracleAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"l2ChainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"useBridgehub_\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawDelay_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"L2_MESSENGER_SYSTEM_CONTRACT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchImportValidatorsToL2\",\"inputs\":[{\"name\":\"validatorIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"l2GasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2GasPerPubdataByteLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refundRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"checkpointRootByRound\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"consumedMessages\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeFromL2\",\"inputs\":[{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2LogIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2TxNumberInBlock\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"proof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"l2ChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2Oracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRound\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mailbox\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIZkSyncMailbox\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextReplacementRequestId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingValidatorIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerValidatorL1\",\"inputs\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structL1Hub.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"replacements\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"targetValidatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidateValidatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidateAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"candidateStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidatePubKey\",\"type\":\"tuple\",\"internalType\":\"structL1Hub.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestReplacementL1\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structL1Hub.ReplacementParams\",\"components\":[{\"name\":\"targetValidatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidateValidatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidatePubKey\",\"type\":\"tuple\",\"internalType\":\"structL1Hub.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"candidateStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"targetLeafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"l2GasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"l2GasPerPubdataByteLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refundRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setL2Oracle\",\"inputs\":[{\"name\":\"l2OracleAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUseBridgehub\",\"inputs\":[{\"name\":\"useBridgehub_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawDelay\",\"inputs\":[{\"name\":\"withdrawDelay_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"useBridgehub\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structL1Hub.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exitFinalized\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawn\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"importRequested\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exitFinalizedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"CheckpointFinalized\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"root\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExitFinalized\",\"inputs\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L2MessageConsumed\",\"inputs\":[{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"msgType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumL1Hub.L2ToL1MsgType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L2OracleUpdated\",\"inputs\":[{\"name\":\"l2Oracle\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplacementFinalized\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"newRoot\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplacementRequested\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"targetValidatorID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"candidateAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"stake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRegisteredL1\",\"inputs\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"stake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorsImportFinalized\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newRoot\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorsImportRequested\",\"inputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"count\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalFinalized\",\"inputs\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// L1HubABI is the input ABI used to generate the binding from.
// Deprecated: Use L1HubMetaData.ABI instead.
var L1HubABI = L1HubMetaData.ABI

// L1Hub is an auto generated Go binding around an Ethereum contract.
type L1Hub struct {
	L1HubCaller     // Read-only binding to the contract
	L1HubTransactor // Write-only binding to the contract
	L1HubFilterer   // Log filterer for contract events
}

// L1HubCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1HubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1HubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1HubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1HubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1HubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1HubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1HubSession struct {
	Contract     *L1Hub            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1HubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1HubCallerSession struct {
	Contract *L1HubCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// L1HubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1HubTransactorSession struct {
	Contract     *L1HubTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1HubRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1HubRaw struct {
	Contract *L1Hub // Generic contract binding to access the raw methods on
}

// L1HubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1HubCallerRaw struct {
	Contract *L1HubCaller // Generic read-only contract binding to access the raw methods on
}

// L1HubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1HubTransactorRaw struct {
	Contract *L1HubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1Hub creates a new instance of L1Hub, bound to a specific deployed contract.
func NewL1Hub(address common.Address, backend bind.ContractBackend) (*L1Hub, error) {
	contract, err := bindL1Hub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1Hub{L1HubCaller: L1HubCaller{contract: contract}, L1HubTransactor: L1HubTransactor{contract: contract}, L1HubFilterer: L1HubFilterer{contract: contract}}, nil
}

// NewL1HubCaller creates a new read-only instance of L1Hub, bound to a specific deployed contract.
func NewL1HubCaller(address common.Address, caller bind.ContractCaller) (*L1HubCaller, error) {
	contract, err := bindL1Hub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1HubCaller{contract: contract}, nil
}

// NewL1HubTransactor creates a new write-only instance of L1Hub, bound to a specific deployed contract.
func NewL1HubTransactor(address common.Address, transactor bind.ContractTransactor) (*L1HubTransactor, error) {
	contract, err := bindL1Hub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1HubTransactor{contract: contract}, nil
}

// NewL1HubFilterer creates a new log filterer instance of L1Hub, bound to a specific deployed contract.
func NewL1HubFilterer(address common.Address, filterer bind.ContractFilterer) (*L1HubFilterer, error) {
	contract, err := bindL1Hub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1HubFilterer{contract: contract}, nil
}

// bindL1Hub binds a generic wrapper to an already deployed contract.
func bindL1Hub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1HubMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Hub *L1HubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Hub.Contract.L1HubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Hub *L1HubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Hub.Contract.L1HubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Hub *L1HubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Hub.Contract.L1HubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Hub *L1HubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Hub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Hub *L1HubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Hub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Hub *L1HubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Hub.Contract.contract.Transact(opts, method, params...)
}

// L2MESSENGERSYSTEMCONTRACT is a free data retrieval call binding the contract method 0xed2364c7.
//
// Solidity: function L2_MESSENGER_SYSTEM_CONTRACT() view returns(address)
func (_L1Hub *L1HubCaller) L2MESSENGERSYSTEMCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "L2_MESSENGER_SYSTEM_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2MESSENGERSYSTEMCONTRACT is a free data retrieval call binding the contract method 0xed2364c7.
//
// Solidity: function L2_MESSENGER_SYSTEM_CONTRACT() view returns(address)
func (_L1Hub *L1HubSession) L2MESSENGERSYSTEMCONTRACT() (common.Address, error) {
	return _L1Hub.Contract.L2MESSENGERSYSTEMCONTRACT(&_L1Hub.CallOpts)
}

// L2MESSENGERSYSTEMCONTRACT is a free data retrieval call binding the contract method 0xed2364c7.
//
// Solidity: function L2_MESSENGER_SYSTEM_CONTRACT() view returns(address)
func (_L1Hub *L1HubCallerSession) L2MESSENGERSYSTEMCONTRACT() (common.Address, error) {
	return _L1Hub.Contract.L2MESSENGERSYSTEMCONTRACT(&_L1Hub.CallOpts)
}

// CheckpointRootByRound is a free data retrieval call binding the contract method 0x41907e64.
//
// Solidity: function checkpointRootByRound(uint256 ) view returns(uint256)
func (_L1Hub *L1HubCaller) CheckpointRootByRound(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "checkpointRootByRound", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CheckpointRootByRound is a free data retrieval call binding the contract method 0x41907e64.
//
// Solidity: function checkpointRootByRound(uint256 ) view returns(uint256)
func (_L1Hub *L1HubSession) CheckpointRootByRound(arg0 *big.Int) (*big.Int, error) {
	return _L1Hub.Contract.CheckpointRootByRound(&_L1Hub.CallOpts, arg0)
}

// CheckpointRootByRound is a free data retrieval call binding the contract method 0x41907e64.
//
// Solidity: function checkpointRootByRound(uint256 ) view returns(uint256)
func (_L1Hub *L1HubCallerSession) CheckpointRootByRound(arg0 *big.Int) (*big.Int, error) {
	return _L1Hub.Contract.CheckpointRootByRound(&_L1Hub.CallOpts, arg0)
}

// ConsumedMessages is a free data retrieval call binding the contract method 0x39bc3c81.
//
// Solidity: function consumedMessages(bytes32 ) view returns(bool)
func (_L1Hub *L1HubCaller) ConsumedMessages(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "consumedMessages", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ConsumedMessages is a free data retrieval call binding the contract method 0x39bc3c81.
//
// Solidity: function consumedMessages(bytes32 ) view returns(bool)
func (_L1Hub *L1HubSession) ConsumedMessages(arg0 [32]byte) (bool, error) {
	return _L1Hub.Contract.ConsumedMessages(&_L1Hub.CallOpts, arg0)
}

// ConsumedMessages is a free data retrieval call binding the contract method 0x39bc3c81.
//
// Solidity: function consumedMessages(bytes32 ) view returns(bool)
func (_L1Hub *L1HubCallerSession) ConsumedMessages(arg0 [32]byte) (bool, error) {
	return _L1Hub.Contract.ConsumedMessages(&_L1Hub.CallOpts, arg0)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_L1Hub *L1HubCaller) L2ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "l2ChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_L1Hub *L1HubSession) L2ChainId() (*big.Int, error) {
	return _L1Hub.Contract.L2ChainId(&_L1Hub.CallOpts)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_L1Hub *L1HubCallerSession) L2ChainId() (*big.Int, error) {
	return _L1Hub.Contract.L2ChainId(&_L1Hub.CallOpts)
}

// L2Oracle is a free data retrieval call binding the contract method 0x9b5f694a.
//
// Solidity: function l2Oracle() view returns(address)
func (_L1Hub *L1HubCaller) L2Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "l2Oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Oracle is a free data retrieval call binding the contract method 0x9b5f694a.
//
// Solidity: function l2Oracle() view returns(address)
func (_L1Hub *L1HubSession) L2Oracle() (common.Address, error) {
	return _L1Hub.Contract.L2Oracle(&_L1Hub.CallOpts)
}

// L2Oracle is a free data retrieval call binding the contract method 0x9b5f694a.
//
// Solidity: function l2Oracle() view returns(address)
func (_L1Hub *L1HubCallerSession) L2Oracle() (common.Address, error) {
	return _L1Hub.Contract.L2Oracle(&_L1Hub.CallOpts)
}

// LatestRoot is a free data retrieval call binding the contract method 0xd7b0fef1.
//
// Solidity: function latestRoot() view returns(uint256)
func (_L1Hub *L1HubCaller) LatestRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "latestRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRoot is a free data retrieval call binding the contract method 0xd7b0fef1.
//
// Solidity: function latestRoot() view returns(uint256)
func (_L1Hub *L1HubSession) LatestRoot() (*big.Int, error) {
	return _L1Hub.Contract.LatestRoot(&_L1Hub.CallOpts)
}

// LatestRoot is a free data retrieval call binding the contract method 0xd7b0fef1.
//
// Solidity: function latestRoot() view returns(uint256)
func (_L1Hub *L1HubCallerSession) LatestRoot() (*big.Int, error) {
	return _L1Hub.Contract.LatestRoot(&_L1Hub.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_L1Hub *L1HubCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_L1Hub *L1HubSession) LatestRound() (*big.Int, error) {
	return _L1Hub.Contract.LatestRound(&_L1Hub.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_L1Hub *L1HubCallerSession) LatestRound() (*big.Int, error) {
	return _L1Hub.Contract.LatestRound(&_L1Hub.CallOpts)
}

// Mailbox is a free data retrieval call binding the contract method 0xd5438eae.
//
// Solidity: function mailbox() view returns(address)
func (_L1Hub *L1HubCaller) Mailbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "mailbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Mailbox is a free data retrieval call binding the contract method 0xd5438eae.
//
// Solidity: function mailbox() view returns(address)
func (_L1Hub *L1HubSession) Mailbox() (common.Address, error) {
	return _L1Hub.Contract.Mailbox(&_L1Hub.CallOpts)
}

// Mailbox is a free data retrieval call binding the contract method 0xd5438eae.
//
// Solidity: function mailbox() view returns(address)
func (_L1Hub *L1HubCallerSession) Mailbox() (common.Address, error) {
	return _L1Hub.Contract.Mailbox(&_L1Hub.CallOpts)
}

// NextReplacementRequestId is a free data retrieval call binding the contract method 0xde230c0b.
//
// Solidity: function nextReplacementRequestId() view returns(uint256)
func (_L1Hub *L1HubCaller) NextReplacementRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "nextReplacementRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextReplacementRequestId is a free data retrieval call binding the contract method 0xde230c0b.
//
// Solidity: function nextReplacementRequestId() view returns(uint256)
func (_L1Hub *L1HubSession) NextReplacementRequestId() (*big.Int, error) {
	return _L1Hub.Contract.NextReplacementRequestId(&_L1Hub.CallOpts)
}

// NextReplacementRequestId is a free data retrieval call binding the contract method 0xde230c0b.
//
// Solidity: function nextReplacementRequestId() view returns(uint256)
func (_L1Hub *L1HubCallerSession) NextReplacementRequestId() (*big.Int, error) {
	return _L1Hub.Contract.NextReplacementRequestId(&_L1Hub.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Hub *L1HubCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Hub *L1HubSession) Owner() (common.Address, error) {
	return _L1Hub.Contract.Owner(&_L1Hub.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Hub *L1HubCallerSession) Owner() (common.Address, error) {
	return _L1Hub.Contract.Owner(&_L1Hub.CallOpts)
}

// PendingValidatorIds is a free data retrieval call binding the contract method 0x43bfdb53.
//
// Solidity: function pendingValidatorIds(uint256 ) view returns(uint256)
func (_L1Hub *L1HubCaller) PendingValidatorIds(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "pendingValidatorIds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingValidatorIds is a free data retrieval call binding the contract method 0x43bfdb53.
//
// Solidity: function pendingValidatorIds(uint256 ) view returns(uint256)
func (_L1Hub *L1HubSession) PendingValidatorIds(arg0 *big.Int) (*big.Int, error) {
	return _L1Hub.Contract.PendingValidatorIds(&_L1Hub.CallOpts, arg0)
}

// PendingValidatorIds is a free data retrieval call binding the contract method 0x43bfdb53.
//
// Solidity: function pendingValidatorIds(uint256 ) view returns(uint256)
func (_L1Hub *L1HubCallerSession) PendingValidatorIds(arg0 *big.Int) (*big.Int, error) {
	return _L1Hub.Contract.PendingValidatorIds(&_L1Hub.CallOpts, arg0)
}

// Replacements is a free data retrieval call binding the contract method 0x44688429.
//
// Solidity: function replacements(uint256 ) view returns(uint256 requestId, uint256 targetValidatorID, uint256 candidateValidatorID, address candidateAddr, uint256 candidateStake, (uint256,uint256) candidatePubKey, bool finalized, bool success)
func (_L1Hub *L1HubCaller) Replacements(opts *bind.CallOpts, arg0 *big.Int) (struct {
	RequestId            *big.Int
	TargetValidatorID    *big.Int
	CandidateValidatorID *big.Int
	CandidateAddr        common.Address
	CandidateStake       *big.Int
	CandidatePubKey      L1HubPublicKey
	Finalized            bool
	Success              bool
}, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "replacements", arg0)

	outstruct := new(struct {
		RequestId            *big.Int
		TargetValidatorID    *big.Int
		CandidateValidatorID *big.Int
		CandidateAddr        common.Address
		CandidateStake       *big.Int
		CandidatePubKey      L1HubPublicKey
		Finalized            bool
		Success              bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TargetValidatorID = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.CandidateValidatorID = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CandidateAddr = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.CandidateStake = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.CandidatePubKey = *abi.ConvertType(out[5], new(L1HubPublicKey)).(*L1HubPublicKey)
	outstruct.Finalized = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.Success = *abi.ConvertType(out[7], new(bool)).(*bool)

	return *outstruct, err

}

// Replacements is a free data retrieval call binding the contract method 0x44688429.
//
// Solidity: function replacements(uint256 ) view returns(uint256 requestId, uint256 targetValidatorID, uint256 candidateValidatorID, address candidateAddr, uint256 candidateStake, (uint256,uint256) candidatePubKey, bool finalized, bool success)
func (_L1Hub *L1HubSession) Replacements(arg0 *big.Int) (struct {
	RequestId            *big.Int
	TargetValidatorID    *big.Int
	CandidateValidatorID *big.Int
	CandidateAddr        common.Address
	CandidateStake       *big.Int
	CandidatePubKey      L1HubPublicKey
	Finalized            bool
	Success              bool
}, error) {
	return _L1Hub.Contract.Replacements(&_L1Hub.CallOpts, arg0)
}

// Replacements is a free data retrieval call binding the contract method 0x44688429.
//
// Solidity: function replacements(uint256 ) view returns(uint256 requestId, uint256 targetValidatorID, uint256 candidateValidatorID, address candidateAddr, uint256 candidateStake, (uint256,uint256) candidatePubKey, bool finalized, bool success)
func (_L1Hub *L1HubCallerSession) Replacements(arg0 *big.Int) (struct {
	RequestId            *big.Int
	TargetValidatorID    *big.Int
	CandidateValidatorID *big.Int
	CandidateAddr        common.Address
	CandidateStake       *big.Int
	CandidatePubKey      L1HubPublicKey
	Finalized            bool
	Success              bool
}, error) {
	return _L1Hub.Contract.Replacements(&_L1Hub.CallOpts, arg0)
}

// UseBridgehub is a free data retrieval call binding the contract method 0x3cc53b3c.
//
// Solidity: function useBridgehub() view returns(bool)
func (_L1Hub *L1HubCaller) UseBridgehub(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "useBridgehub")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UseBridgehub is a free data retrieval call binding the contract method 0x3cc53b3c.
//
// Solidity: function useBridgehub() view returns(bool)
func (_L1Hub *L1HubSession) UseBridgehub() (bool, error) {
	return _L1Hub.Contract.UseBridgehub(&_L1Hub.CallOpts)
}

// UseBridgehub is a free data retrieval call binding the contract method 0x3cc53b3c.
//
// Solidity: function useBridgehub() view returns(bool)
func (_L1Hub *L1HubCallerSession) UseBridgehub() (bool, error) {
	return _L1Hub.Contract.UseBridgehub(&_L1Hub.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) view returns(address validatorAddr, uint256 stake, (uint256,uint256) pubKey, bool active, bool exitFinalized, bool withdrawn, bool importRequested, uint256 exitFinalizedAt)
func (_L1Hub *L1HubCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (struct {
	ValidatorAddr   common.Address
	Stake           *big.Int
	PubKey          L1HubPublicKey
	Active          bool
	ExitFinalized   bool
	Withdrawn       bool
	ImportRequested bool
	ExitFinalizedAt *big.Int
}, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "validators", arg0)

	outstruct := new(struct {
		ValidatorAddr   common.Address
		Stake           *big.Int
		PubKey          L1HubPublicKey
		Active          bool
		ExitFinalized   bool
		Withdrawn       bool
		ImportRequested bool
		ExitFinalizedAt *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ValidatorAddr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Stake = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PubKey = *abi.ConvertType(out[2], new(L1HubPublicKey)).(*L1HubPublicKey)
	outstruct.Active = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.ExitFinalized = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.Withdrawn = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.ImportRequested = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.ExitFinalizedAt = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) view returns(address validatorAddr, uint256 stake, (uint256,uint256) pubKey, bool active, bool exitFinalized, bool withdrawn, bool importRequested, uint256 exitFinalizedAt)
func (_L1Hub *L1HubSession) Validators(arg0 *big.Int) (struct {
	ValidatorAddr   common.Address
	Stake           *big.Int
	PubKey          L1HubPublicKey
	Active          bool
	ExitFinalized   bool
	Withdrawn       bool
	ImportRequested bool
	ExitFinalizedAt *big.Int
}, error) {
	return _L1Hub.Contract.Validators(&_L1Hub.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) view returns(address validatorAddr, uint256 stake, (uint256,uint256) pubKey, bool active, bool exitFinalized, bool withdrawn, bool importRequested, uint256 exitFinalizedAt)
func (_L1Hub *L1HubCallerSession) Validators(arg0 *big.Int) (struct {
	ValidatorAddr   common.Address
	Stake           *big.Int
	PubKey          L1HubPublicKey
	Active          bool
	ExitFinalized   bool
	Withdrawn       bool
	ImportRequested bool
	ExitFinalizedAt *big.Int
}, error) {
	return _L1Hub.Contract.Validators(&_L1Hub.CallOpts, arg0)
}

// WithdrawDelay is a free data retrieval call binding the contract method 0x0288a39c.
//
// Solidity: function withdrawDelay() view returns(uint256)
func (_L1Hub *L1HubCaller) WithdrawDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Hub.contract.Call(opts, &out, "withdrawDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawDelay is a free data retrieval call binding the contract method 0x0288a39c.
//
// Solidity: function withdrawDelay() view returns(uint256)
func (_L1Hub *L1HubSession) WithdrawDelay() (*big.Int, error) {
	return _L1Hub.Contract.WithdrawDelay(&_L1Hub.CallOpts)
}

// WithdrawDelay is a free data retrieval call binding the contract method 0x0288a39c.
//
// Solidity: function withdrawDelay() view returns(uint256)
func (_L1Hub *L1HubCallerSession) WithdrawDelay() (*big.Int, error) {
	return _L1Hub.Contract.WithdrawDelay(&_L1Hub.CallOpts)
}

// BatchImportValidatorsToL2 is a paid mutator transaction binding the contract method 0x5efcf9ff.
//
// Solidity: function batchImportValidatorsToL2(uint256[] validatorIds, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit, address refundRecipient) payable returns(bytes32 txHash)
func (_L1Hub *L1HubTransactor) BatchImportValidatorsToL2(opts *bind.TransactOpts, validatorIds []*big.Int, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int, refundRecipient common.Address) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "batchImportValidatorsToL2", validatorIds, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient)
}

// BatchImportValidatorsToL2 is a paid mutator transaction binding the contract method 0x5efcf9ff.
//
// Solidity: function batchImportValidatorsToL2(uint256[] validatorIds, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit, address refundRecipient) payable returns(bytes32 txHash)
func (_L1Hub *L1HubSession) BatchImportValidatorsToL2(validatorIds []*big.Int, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int, refundRecipient common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.BatchImportValidatorsToL2(&_L1Hub.TransactOpts, validatorIds, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient)
}

// BatchImportValidatorsToL2 is a paid mutator transaction binding the contract method 0x5efcf9ff.
//
// Solidity: function batchImportValidatorsToL2(uint256[] validatorIds, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit, address refundRecipient) payable returns(bytes32 txHash)
func (_L1Hub *L1HubTransactorSession) BatchImportValidatorsToL2(validatorIds []*big.Int, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int, refundRecipient common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.BatchImportValidatorsToL2(&_L1Hub.TransactOpts, validatorIds, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient)
}

// FinalizeFromL2 is a paid mutator transaction binding the contract method 0x300a93e5.
//
// Solidity: function finalizeFromL2(bytes message, uint256 l2BlockNumber, uint256 l2LogIndex, uint16 l2TxNumberInBlock, bytes32[] proof) returns()
func (_L1Hub *L1HubTransactor) FinalizeFromL2(opts *bind.TransactOpts, message []byte, l2BlockNumber *big.Int, l2LogIndex *big.Int, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "finalizeFromL2", message, l2BlockNumber, l2LogIndex, l2TxNumberInBlock, proof)
}

// FinalizeFromL2 is a paid mutator transaction binding the contract method 0x300a93e5.
//
// Solidity: function finalizeFromL2(bytes message, uint256 l2BlockNumber, uint256 l2LogIndex, uint16 l2TxNumberInBlock, bytes32[] proof) returns()
func (_L1Hub *L1HubSession) FinalizeFromL2(message []byte, l2BlockNumber *big.Int, l2LogIndex *big.Int, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	return _L1Hub.Contract.FinalizeFromL2(&_L1Hub.TransactOpts, message, l2BlockNumber, l2LogIndex, l2TxNumberInBlock, proof)
}

// FinalizeFromL2 is a paid mutator transaction binding the contract method 0x300a93e5.
//
// Solidity: function finalizeFromL2(bytes message, uint256 l2BlockNumber, uint256 l2LogIndex, uint16 l2TxNumberInBlock, bytes32[] proof) returns()
func (_L1Hub *L1HubTransactorSession) FinalizeFromL2(message []byte, l2BlockNumber *big.Int, l2LogIndex *big.Int, l2TxNumberInBlock uint16, proof [][32]byte) (*types.Transaction, error) {
	return _L1Hub.Contract.FinalizeFromL2(&_L1Hub.TransactOpts, message, l2BlockNumber, l2LogIndex, l2TxNumberInBlock, proof)
}

// RegisterValidatorL1 is a paid mutator transaction binding the contract method 0x8372b086.
//
// Solidity: function registerValidatorL1(uint256 validatorID, (uint256,uint256) pubKey) payable returns()
func (_L1Hub *L1HubTransactor) RegisterValidatorL1(opts *bind.TransactOpts, validatorID *big.Int, pubKey L1HubPublicKey) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "registerValidatorL1", validatorID, pubKey)
}

// RegisterValidatorL1 is a paid mutator transaction binding the contract method 0x8372b086.
//
// Solidity: function registerValidatorL1(uint256 validatorID, (uint256,uint256) pubKey) payable returns()
func (_L1Hub *L1HubSession) RegisterValidatorL1(validatorID *big.Int, pubKey L1HubPublicKey) (*types.Transaction, error) {
	return _L1Hub.Contract.RegisterValidatorL1(&_L1Hub.TransactOpts, validatorID, pubKey)
}

// RegisterValidatorL1 is a paid mutator transaction binding the contract method 0x8372b086.
//
// Solidity: function registerValidatorL1(uint256 validatorID, (uint256,uint256) pubKey) payable returns()
func (_L1Hub *L1HubTransactorSession) RegisterValidatorL1(validatorID *big.Int, pubKey L1HubPublicKey) (*types.Transaction, error) {
	return _L1Hub.Contract.RegisterValidatorL1(&_L1Hub.TransactOpts, validatorID, pubKey)
}

// RequestReplacementL1 is a paid mutator transaction binding the contract method 0x98bbf3f2.
//
// Solidity: function requestReplacementL1((uint256,uint256,(uint256,uint256),uint256,uint256,uint256[],uint256) params, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit, address refundRecipient) payable returns(uint256 requestId, bytes32 txHash)
func (_L1Hub *L1HubTransactor) RequestReplacementL1(opts *bind.TransactOpts, params L1HubReplacementParams, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int, refundRecipient common.Address) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "requestReplacementL1", params, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient)
}

// RequestReplacementL1 is a paid mutator transaction binding the contract method 0x98bbf3f2.
//
// Solidity: function requestReplacementL1((uint256,uint256,(uint256,uint256),uint256,uint256,uint256[],uint256) params, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit, address refundRecipient) payable returns(uint256 requestId, bytes32 txHash)
func (_L1Hub *L1HubSession) RequestReplacementL1(params L1HubReplacementParams, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int, refundRecipient common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.RequestReplacementL1(&_L1Hub.TransactOpts, params, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient)
}

// RequestReplacementL1 is a paid mutator transaction binding the contract method 0x98bbf3f2.
//
// Solidity: function requestReplacementL1((uint256,uint256,(uint256,uint256),uint256,uint256,uint256[],uint256) params, uint256 l2GasLimit, uint256 l2GasPerPubdataByteLimit, address refundRecipient) payable returns(uint256 requestId, bytes32 txHash)
func (_L1Hub *L1HubTransactorSession) RequestReplacementL1(params L1HubReplacementParams, l2GasLimit *big.Int, l2GasPerPubdataByteLimit *big.Int, refundRecipient common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.RequestReplacementL1(&_L1Hub.TransactOpts, params, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient)
}

// SetL2Oracle is a paid mutator transaction binding the contract method 0x7ab4189e.
//
// Solidity: function setL2Oracle(address l2OracleAddress) returns()
func (_L1Hub *L1HubTransactor) SetL2Oracle(opts *bind.TransactOpts, l2OracleAddress common.Address) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "setL2Oracle", l2OracleAddress)
}

// SetL2Oracle is a paid mutator transaction binding the contract method 0x7ab4189e.
//
// Solidity: function setL2Oracle(address l2OracleAddress) returns()
func (_L1Hub *L1HubSession) SetL2Oracle(l2OracleAddress common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.SetL2Oracle(&_L1Hub.TransactOpts, l2OracleAddress)
}

// SetL2Oracle is a paid mutator transaction binding the contract method 0x7ab4189e.
//
// Solidity: function setL2Oracle(address l2OracleAddress) returns()
func (_L1Hub *L1HubTransactorSession) SetL2Oracle(l2OracleAddress common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.SetL2Oracle(&_L1Hub.TransactOpts, l2OracleAddress)
}

// SetUseBridgehub is a paid mutator transaction binding the contract method 0x346b7f00.
//
// Solidity: function setUseBridgehub(bool useBridgehub_) returns()
func (_L1Hub *L1HubTransactor) SetUseBridgehub(opts *bind.TransactOpts, useBridgehub_ bool) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "setUseBridgehub", useBridgehub_)
}

// SetUseBridgehub is a paid mutator transaction binding the contract method 0x346b7f00.
//
// Solidity: function setUseBridgehub(bool useBridgehub_) returns()
func (_L1Hub *L1HubSession) SetUseBridgehub(useBridgehub_ bool) (*types.Transaction, error) {
	return _L1Hub.Contract.SetUseBridgehub(&_L1Hub.TransactOpts, useBridgehub_)
}

// SetUseBridgehub is a paid mutator transaction binding the contract method 0x346b7f00.
//
// Solidity: function setUseBridgehub(bool useBridgehub_) returns()
func (_L1Hub *L1HubTransactorSession) SetUseBridgehub(useBridgehub_ bool) (*types.Transaction, error) {
	return _L1Hub.Contract.SetUseBridgehub(&_L1Hub.TransactOpts, useBridgehub_)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x72f0cb30.
//
// Solidity: function setWithdrawDelay(uint256 withdrawDelay_) returns()
func (_L1Hub *L1HubTransactor) SetWithdrawDelay(opts *bind.TransactOpts, withdrawDelay_ *big.Int) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "setWithdrawDelay", withdrawDelay_)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x72f0cb30.
//
// Solidity: function setWithdrawDelay(uint256 withdrawDelay_) returns()
func (_L1Hub *L1HubSession) SetWithdrawDelay(withdrawDelay_ *big.Int) (*types.Transaction, error) {
	return _L1Hub.Contract.SetWithdrawDelay(&_L1Hub.TransactOpts, withdrawDelay_)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x72f0cb30.
//
// Solidity: function setWithdrawDelay(uint256 withdrawDelay_) returns()
func (_L1Hub *L1HubTransactorSession) SetWithdrawDelay(withdrawDelay_ *big.Int) (*types.Transaction, error) {
	return _L1Hub.Contract.SetWithdrawDelay(&_L1Hub.TransactOpts, withdrawDelay_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Hub *L1HubTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L1Hub.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Hub *L1HubSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.TransferOwnership(&_L1Hub.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Hub *L1HubTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Hub.Contract.TransferOwnership(&_L1Hub.TransactOpts, newOwner)
}

// L1HubCheckpointFinalizedIterator is returned from FilterCheckpointFinalized and is used to iterate over the raw logs and unpacked data for CheckpointFinalized events raised by the L1Hub contract.
type L1HubCheckpointFinalizedIterator struct {
	Event *L1HubCheckpointFinalized // Event containing the contract specifics and raw log

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
func (it *L1HubCheckpointFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubCheckpointFinalized)
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
		it.Event = new(L1HubCheckpointFinalized)
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
func (it *L1HubCheckpointFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubCheckpointFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubCheckpointFinalized represents a CheckpointFinalized event raised by the L1Hub contract.
type L1HubCheckpointFinalized struct {
	Round *big.Int
	Root  *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCheckpointFinalized is a free log retrieval operation binding the contract event 0x85e167eeb69b4ba558856559b63e15f189b339a4db8fec09f7adababb7216f1f.
//
// Solidity: event CheckpointFinalized(uint256 indexed round, uint256 root)
func (_L1Hub *L1HubFilterer) FilterCheckpointFinalized(opts *bind.FilterOpts, round []*big.Int) (*L1HubCheckpointFinalizedIterator, error) {

	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "CheckpointFinalized", roundRule)
	if err != nil {
		return nil, err
	}
	return &L1HubCheckpointFinalizedIterator{contract: _L1Hub.contract, event: "CheckpointFinalized", logs: logs, sub: sub}, nil
}

// WatchCheckpointFinalized is a free log subscription operation binding the contract event 0x85e167eeb69b4ba558856559b63e15f189b339a4db8fec09f7adababb7216f1f.
//
// Solidity: event CheckpointFinalized(uint256 indexed round, uint256 root)
func (_L1Hub *L1HubFilterer) WatchCheckpointFinalized(opts *bind.WatchOpts, sink chan<- *L1HubCheckpointFinalized, round []*big.Int) (event.Subscription, error) {

	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "CheckpointFinalized", roundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubCheckpointFinalized)
				if err := _L1Hub.contract.UnpackLog(event, "CheckpointFinalized", log); err != nil {
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

// ParseCheckpointFinalized is a log parse operation binding the contract event 0x85e167eeb69b4ba558856559b63e15f189b339a4db8fec09f7adababb7216f1f.
//
// Solidity: event CheckpointFinalized(uint256 indexed round, uint256 root)
func (_L1Hub *L1HubFilterer) ParseCheckpointFinalized(log types.Log) (*L1HubCheckpointFinalized, error) {
	event := new(L1HubCheckpointFinalized)
	if err := _L1Hub.contract.UnpackLog(event, "CheckpointFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubExitFinalizedIterator is returned from FilterExitFinalized and is used to iterate over the raw logs and unpacked data for ExitFinalized events raised by the L1Hub contract.
type L1HubExitFinalizedIterator struct {
	Event *L1HubExitFinalized // Event containing the contract specifics and raw log

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
func (it *L1HubExitFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubExitFinalized)
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
		it.Event = new(L1HubExitFinalized)
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
func (it *L1HubExitFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubExitFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubExitFinalized represents a ExitFinalized event raised by the L1Hub contract.
type L1HubExitFinalized struct {
	ValidatorID   *big.Int
	ValidatorAddr common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterExitFinalized is a free log retrieval operation binding the contract event 0xbf865cf80bb20ce94cc64a00c758ae8f3a59e33abbe4a625131ae692d6e997c6.
//
// Solidity: event ExitFinalized(uint256 indexed validatorID, address indexed validatorAddr)
func (_L1Hub *L1HubFilterer) FilterExitFinalized(opts *bind.FilterOpts, validatorID []*big.Int, validatorAddr []common.Address) (*L1HubExitFinalizedIterator, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}
	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "ExitFinalized", validatorIDRule, validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return &L1HubExitFinalizedIterator{contract: _L1Hub.contract, event: "ExitFinalized", logs: logs, sub: sub}, nil
}

// WatchExitFinalized is a free log subscription operation binding the contract event 0xbf865cf80bb20ce94cc64a00c758ae8f3a59e33abbe4a625131ae692d6e997c6.
//
// Solidity: event ExitFinalized(uint256 indexed validatorID, address indexed validatorAddr)
func (_L1Hub *L1HubFilterer) WatchExitFinalized(opts *bind.WatchOpts, sink chan<- *L1HubExitFinalized, validatorID []*big.Int, validatorAddr []common.Address) (event.Subscription, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}
	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "ExitFinalized", validatorIDRule, validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubExitFinalized)
				if err := _L1Hub.contract.UnpackLog(event, "ExitFinalized", log); err != nil {
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

// ParseExitFinalized is a log parse operation binding the contract event 0xbf865cf80bb20ce94cc64a00c758ae8f3a59e33abbe4a625131ae692d6e997c6.
//
// Solidity: event ExitFinalized(uint256 indexed validatorID, address indexed validatorAddr)
func (_L1Hub *L1HubFilterer) ParseExitFinalized(log types.Log) (*L1HubExitFinalized, error) {
	event := new(L1HubExitFinalized)
	if err := _L1Hub.contract.UnpackLog(event, "ExitFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubL2MessageConsumedIterator is returned from FilterL2MessageConsumed and is used to iterate over the raw logs and unpacked data for L2MessageConsumed events raised by the L1Hub contract.
type L1HubL2MessageConsumedIterator struct {
	Event *L1HubL2MessageConsumed // Event containing the contract specifics and raw log

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
func (it *L1HubL2MessageConsumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubL2MessageConsumed)
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
		it.Event = new(L1HubL2MessageConsumed)
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
func (it *L1HubL2MessageConsumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubL2MessageConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubL2MessageConsumed represents a L2MessageConsumed event raised by the L1Hub contract.
type L1HubL2MessageConsumed struct {
	MessageHash [32]byte
	MsgType     uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterL2MessageConsumed is a free log retrieval operation binding the contract event 0x7853b350a8b8e5355f9b088e7f20139dd39ecc8ddc8ebebb0b5f88956f40e706.
//
// Solidity: event L2MessageConsumed(bytes32 indexed messageHash, uint8 msgType)
func (_L1Hub *L1HubFilterer) FilterL2MessageConsumed(opts *bind.FilterOpts, messageHash [][32]byte) (*L1HubL2MessageConsumedIterator, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "L2MessageConsumed", messageHashRule)
	if err != nil {
		return nil, err
	}
	return &L1HubL2MessageConsumedIterator{contract: _L1Hub.contract, event: "L2MessageConsumed", logs: logs, sub: sub}, nil
}

// WatchL2MessageConsumed is a free log subscription operation binding the contract event 0x7853b350a8b8e5355f9b088e7f20139dd39ecc8ddc8ebebb0b5f88956f40e706.
//
// Solidity: event L2MessageConsumed(bytes32 indexed messageHash, uint8 msgType)
func (_L1Hub *L1HubFilterer) WatchL2MessageConsumed(opts *bind.WatchOpts, sink chan<- *L1HubL2MessageConsumed, messageHash [][32]byte) (event.Subscription, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "L2MessageConsumed", messageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubL2MessageConsumed)
				if err := _L1Hub.contract.UnpackLog(event, "L2MessageConsumed", log); err != nil {
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

// ParseL2MessageConsumed is a log parse operation binding the contract event 0x7853b350a8b8e5355f9b088e7f20139dd39ecc8ddc8ebebb0b5f88956f40e706.
//
// Solidity: event L2MessageConsumed(bytes32 indexed messageHash, uint8 msgType)
func (_L1Hub *L1HubFilterer) ParseL2MessageConsumed(log types.Log) (*L1HubL2MessageConsumed, error) {
	event := new(L1HubL2MessageConsumed)
	if err := _L1Hub.contract.UnpackLog(event, "L2MessageConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubL2OracleUpdatedIterator is returned from FilterL2OracleUpdated and is used to iterate over the raw logs and unpacked data for L2OracleUpdated events raised by the L1Hub contract.
type L1HubL2OracleUpdatedIterator struct {
	Event *L1HubL2OracleUpdated // Event containing the contract specifics and raw log

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
func (it *L1HubL2OracleUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubL2OracleUpdated)
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
		it.Event = new(L1HubL2OracleUpdated)
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
func (it *L1HubL2OracleUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubL2OracleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubL2OracleUpdated represents a L2OracleUpdated event raised by the L1Hub contract.
type L1HubL2OracleUpdated struct {
	L2Oracle common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterL2OracleUpdated is a free log retrieval operation binding the contract event 0x506054661a1c4e95e9ab185811eb2dc08e25bd05116fced9a45c25d25e9ab3c9.
//
// Solidity: event L2OracleUpdated(address indexed l2Oracle)
func (_L1Hub *L1HubFilterer) FilterL2OracleUpdated(opts *bind.FilterOpts, l2Oracle []common.Address) (*L1HubL2OracleUpdatedIterator, error) {

	var l2OracleRule []interface{}
	for _, l2OracleItem := range l2Oracle {
		l2OracleRule = append(l2OracleRule, l2OracleItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "L2OracleUpdated", l2OracleRule)
	if err != nil {
		return nil, err
	}
	return &L1HubL2OracleUpdatedIterator{contract: _L1Hub.contract, event: "L2OracleUpdated", logs: logs, sub: sub}, nil
}

// WatchL2OracleUpdated is a free log subscription operation binding the contract event 0x506054661a1c4e95e9ab185811eb2dc08e25bd05116fced9a45c25d25e9ab3c9.
//
// Solidity: event L2OracleUpdated(address indexed l2Oracle)
func (_L1Hub *L1HubFilterer) WatchL2OracleUpdated(opts *bind.WatchOpts, sink chan<- *L1HubL2OracleUpdated, l2Oracle []common.Address) (event.Subscription, error) {

	var l2OracleRule []interface{}
	for _, l2OracleItem := range l2Oracle {
		l2OracleRule = append(l2OracleRule, l2OracleItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "L2OracleUpdated", l2OracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubL2OracleUpdated)
				if err := _L1Hub.contract.UnpackLog(event, "L2OracleUpdated", log); err != nil {
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

// ParseL2OracleUpdated is a log parse operation binding the contract event 0x506054661a1c4e95e9ab185811eb2dc08e25bd05116fced9a45c25d25e9ab3c9.
//
// Solidity: event L2OracleUpdated(address indexed l2Oracle)
func (_L1Hub *L1HubFilterer) ParseL2OracleUpdated(log types.Log) (*L1HubL2OracleUpdated, error) {
	event := new(L1HubL2OracleUpdated)
	if err := _L1Hub.contract.UnpackLog(event, "L2OracleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L1Hub contract.
type L1HubOwnershipTransferredIterator struct {
	Event *L1HubOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L1HubOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubOwnershipTransferred)
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
		it.Event = new(L1HubOwnershipTransferred)
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
func (it *L1HubOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubOwnershipTransferred represents a OwnershipTransferred event raised by the L1Hub contract.
type L1HubOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Hub *L1HubFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1HubOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1HubOwnershipTransferredIterator{contract: _L1Hub.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Hub *L1HubFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L1HubOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubOwnershipTransferred)
				if err := _L1Hub.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L1Hub *L1HubFilterer) ParseOwnershipTransferred(log types.Log) (*L1HubOwnershipTransferred, error) {
	event := new(L1HubOwnershipTransferred)
	if err := _L1Hub.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubReplacementFinalizedIterator is returned from FilterReplacementFinalized and is used to iterate over the raw logs and unpacked data for ReplacementFinalized events raised by the L1Hub contract.
type L1HubReplacementFinalizedIterator struct {
	Event *L1HubReplacementFinalized // Event containing the contract specifics and raw log

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
func (it *L1HubReplacementFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubReplacementFinalized)
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
		it.Event = new(L1HubReplacementFinalized)
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
func (it *L1HubReplacementFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubReplacementFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubReplacementFinalized represents a ReplacementFinalized event raised by the L1Hub contract.
type L1HubReplacementFinalized struct {
	RequestId *big.Int
	Success   bool
	NewRoot   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReplacementFinalized is a free log retrieval operation binding the contract event 0x9e95d26ac127c801cbc56a49eeb180606a50d94b4306b1e24fcc8725c8effec7.
//
// Solidity: event ReplacementFinalized(uint256 indexed requestId, bool success, uint256 newRoot)
func (_L1Hub *L1HubFilterer) FilterReplacementFinalized(opts *bind.FilterOpts, requestId []*big.Int) (*L1HubReplacementFinalizedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "ReplacementFinalized", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &L1HubReplacementFinalizedIterator{contract: _L1Hub.contract, event: "ReplacementFinalized", logs: logs, sub: sub}, nil
}

// WatchReplacementFinalized is a free log subscription operation binding the contract event 0x9e95d26ac127c801cbc56a49eeb180606a50d94b4306b1e24fcc8725c8effec7.
//
// Solidity: event ReplacementFinalized(uint256 indexed requestId, bool success, uint256 newRoot)
func (_L1Hub *L1HubFilterer) WatchReplacementFinalized(opts *bind.WatchOpts, sink chan<- *L1HubReplacementFinalized, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "ReplacementFinalized", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubReplacementFinalized)
				if err := _L1Hub.contract.UnpackLog(event, "ReplacementFinalized", log); err != nil {
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

// ParseReplacementFinalized is a log parse operation binding the contract event 0x9e95d26ac127c801cbc56a49eeb180606a50d94b4306b1e24fcc8725c8effec7.
//
// Solidity: event ReplacementFinalized(uint256 indexed requestId, bool success, uint256 newRoot)
func (_L1Hub *L1HubFilterer) ParseReplacementFinalized(log types.Log) (*L1HubReplacementFinalized, error) {
	event := new(L1HubReplacementFinalized)
	if err := _L1Hub.contract.UnpackLog(event, "ReplacementFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubReplacementRequestedIterator is returned from FilterReplacementRequested and is used to iterate over the raw logs and unpacked data for ReplacementRequested events raised by the L1Hub contract.
type L1HubReplacementRequestedIterator struct {
	Event *L1HubReplacementRequested // Event containing the contract specifics and raw log

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
func (it *L1HubReplacementRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubReplacementRequested)
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
		it.Event = new(L1HubReplacementRequested)
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
func (it *L1HubReplacementRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubReplacementRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubReplacementRequested represents a ReplacementRequested event raised by the L1Hub contract.
type L1HubReplacementRequested struct {
	RequestId         *big.Int
	TargetValidatorID *big.Int
	CandidateAddr     common.Address
	Stake             *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterReplacementRequested is a free log retrieval operation binding the contract event 0x6b73393a7efe044782f5280f8d882c747547fccc671295192546a63ecdf47920.
//
// Solidity: event ReplacementRequested(uint256 indexed requestId, uint256 indexed targetValidatorID, address indexed candidateAddr, uint256 stake)
func (_L1Hub *L1HubFilterer) FilterReplacementRequested(opts *bind.FilterOpts, requestId []*big.Int, targetValidatorID []*big.Int, candidateAddr []common.Address) (*L1HubReplacementRequestedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var targetValidatorIDRule []interface{}
	for _, targetValidatorIDItem := range targetValidatorID {
		targetValidatorIDRule = append(targetValidatorIDRule, targetValidatorIDItem)
	}
	var candidateAddrRule []interface{}
	for _, candidateAddrItem := range candidateAddr {
		candidateAddrRule = append(candidateAddrRule, candidateAddrItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "ReplacementRequested", requestIdRule, targetValidatorIDRule, candidateAddrRule)
	if err != nil {
		return nil, err
	}
	return &L1HubReplacementRequestedIterator{contract: _L1Hub.contract, event: "ReplacementRequested", logs: logs, sub: sub}, nil
}

// WatchReplacementRequested is a free log subscription operation binding the contract event 0x6b73393a7efe044782f5280f8d882c747547fccc671295192546a63ecdf47920.
//
// Solidity: event ReplacementRequested(uint256 indexed requestId, uint256 indexed targetValidatorID, address indexed candidateAddr, uint256 stake)
func (_L1Hub *L1HubFilterer) WatchReplacementRequested(opts *bind.WatchOpts, sink chan<- *L1HubReplacementRequested, requestId []*big.Int, targetValidatorID []*big.Int, candidateAddr []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var targetValidatorIDRule []interface{}
	for _, targetValidatorIDItem := range targetValidatorID {
		targetValidatorIDRule = append(targetValidatorIDRule, targetValidatorIDItem)
	}
	var candidateAddrRule []interface{}
	for _, candidateAddrItem := range candidateAddr {
		candidateAddrRule = append(candidateAddrRule, candidateAddrItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "ReplacementRequested", requestIdRule, targetValidatorIDRule, candidateAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubReplacementRequested)
				if err := _L1Hub.contract.UnpackLog(event, "ReplacementRequested", log); err != nil {
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

// ParseReplacementRequested is a log parse operation binding the contract event 0x6b73393a7efe044782f5280f8d882c747547fccc671295192546a63ecdf47920.
//
// Solidity: event ReplacementRequested(uint256 indexed requestId, uint256 indexed targetValidatorID, address indexed candidateAddr, uint256 stake)
func (_L1Hub *L1HubFilterer) ParseReplacementRequested(log types.Log) (*L1HubReplacementRequested, error) {
	event := new(L1HubReplacementRequested)
	if err := _L1Hub.contract.UnpackLog(event, "ReplacementRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubValidatorRegisteredL1Iterator is returned from FilterValidatorRegisteredL1 and is used to iterate over the raw logs and unpacked data for ValidatorRegisteredL1 events raised by the L1Hub contract.
type L1HubValidatorRegisteredL1Iterator struct {
	Event *L1HubValidatorRegisteredL1 // Event containing the contract specifics and raw log

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
func (it *L1HubValidatorRegisteredL1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubValidatorRegisteredL1)
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
		it.Event = new(L1HubValidatorRegisteredL1)
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
func (it *L1HubValidatorRegisteredL1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubValidatorRegisteredL1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubValidatorRegisteredL1 represents a ValidatorRegisteredL1 event raised by the L1Hub contract.
type L1HubValidatorRegisteredL1 struct {
	ValidatorID   *big.Int
	ValidatorAddr common.Address
	Stake         *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterValidatorRegisteredL1 is a free log retrieval operation binding the contract event 0x4ddb0569442b90a927479e7e5e9822b98a1cd6d459edaee6dcba1814524f0ae1.
//
// Solidity: event ValidatorRegisteredL1(uint256 indexed validatorID, address indexed validatorAddr, uint256 stake)
func (_L1Hub *L1HubFilterer) FilterValidatorRegisteredL1(opts *bind.FilterOpts, validatorID []*big.Int, validatorAddr []common.Address) (*L1HubValidatorRegisteredL1Iterator, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}
	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "ValidatorRegisteredL1", validatorIDRule, validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return &L1HubValidatorRegisteredL1Iterator{contract: _L1Hub.contract, event: "ValidatorRegisteredL1", logs: logs, sub: sub}, nil
}

// WatchValidatorRegisteredL1 is a free log subscription operation binding the contract event 0x4ddb0569442b90a927479e7e5e9822b98a1cd6d459edaee6dcba1814524f0ae1.
//
// Solidity: event ValidatorRegisteredL1(uint256 indexed validatorID, address indexed validatorAddr, uint256 stake)
func (_L1Hub *L1HubFilterer) WatchValidatorRegisteredL1(opts *bind.WatchOpts, sink chan<- *L1HubValidatorRegisteredL1, validatorID []*big.Int, validatorAddr []common.Address) (event.Subscription, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}
	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "ValidatorRegisteredL1", validatorIDRule, validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubValidatorRegisteredL1)
				if err := _L1Hub.contract.UnpackLog(event, "ValidatorRegisteredL1", log); err != nil {
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

// ParseValidatorRegisteredL1 is a log parse operation binding the contract event 0x4ddb0569442b90a927479e7e5e9822b98a1cd6d459edaee6dcba1814524f0ae1.
//
// Solidity: event ValidatorRegisteredL1(uint256 indexed validatorID, address indexed validatorAddr, uint256 stake)
func (_L1Hub *L1HubFilterer) ParseValidatorRegisteredL1(log types.Log) (*L1HubValidatorRegisteredL1, error) {
	event := new(L1HubValidatorRegisteredL1)
	if err := _L1Hub.contract.UnpackLog(event, "ValidatorRegisteredL1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubValidatorsImportFinalizedIterator is returned from FilterValidatorsImportFinalized and is used to iterate over the raw logs and unpacked data for ValidatorsImportFinalized events raised by the L1Hub contract.
type L1HubValidatorsImportFinalizedIterator struct {
	Event *L1HubValidatorsImportFinalized // Event containing the contract specifics and raw log

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
func (it *L1HubValidatorsImportFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubValidatorsImportFinalized)
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
		it.Event = new(L1HubValidatorsImportFinalized)
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
func (it *L1HubValidatorsImportFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubValidatorsImportFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubValidatorsImportFinalized represents a ValidatorsImportFinalized event raised by the L1Hub contract.
type L1HubValidatorsImportFinalized struct {
	Count   *big.Int
	NewRoot *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterValidatorsImportFinalized is a free log retrieval operation binding the contract event 0x0f4345ceb28ddbe3e09d7dad3def28fd3b5c0b91216a3057b57c430d88da4b73.
//
// Solidity: event ValidatorsImportFinalized(uint256 count, uint256 newRoot)
func (_L1Hub *L1HubFilterer) FilterValidatorsImportFinalized(opts *bind.FilterOpts) (*L1HubValidatorsImportFinalizedIterator, error) {

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "ValidatorsImportFinalized")
	if err != nil {
		return nil, err
	}
	return &L1HubValidatorsImportFinalizedIterator{contract: _L1Hub.contract, event: "ValidatorsImportFinalized", logs: logs, sub: sub}, nil
}

// WatchValidatorsImportFinalized is a free log subscription operation binding the contract event 0x0f4345ceb28ddbe3e09d7dad3def28fd3b5c0b91216a3057b57c430d88da4b73.
//
// Solidity: event ValidatorsImportFinalized(uint256 count, uint256 newRoot)
func (_L1Hub *L1HubFilterer) WatchValidatorsImportFinalized(opts *bind.WatchOpts, sink chan<- *L1HubValidatorsImportFinalized) (event.Subscription, error) {

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "ValidatorsImportFinalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubValidatorsImportFinalized)
				if err := _L1Hub.contract.UnpackLog(event, "ValidatorsImportFinalized", log); err != nil {
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

// ParseValidatorsImportFinalized is a log parse operation binding the contract event 0x0f4345ceb28ddbe3e09d7dad3def28fd3b5c0b91216a3057b57c430d88da4b73.
//
// Solidity: event ValidatorsImportFinalized(uint256 count, uint256 newRoot)
func (_L1Hub *L1HubFilterer) ParseValidatorsImportFinalized(log types.Log) (*L1HubValidatorsImportFinalized, error) {
	event := new(L1HubValidatorsImportFinalized)
	if err := _L1Hub.contract.UnpackLog(event, "ValidatorsImportFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubValidatorsImportRequestedIterator is returned from FilterValidatorsImportRequested and is used to iterate over the raw logs and unpacked data for ValidatorsImportRequested events raised by the L1Hub contract.
type L1HubValidatorsImportRequestedIterator struct {
	Event *L1HubValidatorsImportRequested // Event containing the contract specifics and raw log

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
func (it *L1HubValidatorsImportRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubValidatorsImportRequested)
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
		it.Event = new(L1HubValidatorsImportRequested)
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
func (it *L1HubValidatorsImportRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubValidatorsImportRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubValidatorsImportRequested represents a ValidatorsImportRequested event raised by the L1Hub contract.
type L1HubValidatorsImportRequested struct {
	TxHash [32]byte
	Count  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterValidatorsImportRequested is a free log retrieval operation binding the contract event 0xa6a406bba8edd318c95802e11c6fa6c696172b7bbf1aa90e28baa11a306bd3ca.
//
// Solidity: event ValidatorsImportRequested(bytes32 indexed txHash, uint256 count)
func (_L1Hub *L1HubFilterer) FilterValidatorsImportRequested(opts *bind.FilterOpts, txHash [][32]byte) (*L1HubValidatorsImportRequestedIterator, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "ValidatorsImportRequested", txHashRule)
	if err != nil {
		return nil, err
	}
	return &L1HubValidatorsImportRequestedIterator{contract: _L1Hub.contract, event: "ValidatorsImportRequested", logs: logs, sub: sub}, nil
}

// WatchValidatorsImportRequested is a free log subscription operation binding the contract event 0xa6a406bba8edd318c95802e11c6fa6c696172b7bbf1aa90e28baa11a306bd3ca.
//
// Solidity: event ValidatorsImportRequested(bytes32 indexed txHash, uint256 count)
func (_L1Hub *L1HubFilterer) WatchValidatorsImportRequested(opts *bind.WatchOpts, sink chan<- *L1HubValidatorsImportRequested, txHash [][32]byte) (event.Subscription, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "ValidatorsImportRequested", txHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubValidatorsImportRequested)
				if err := _L1Hub.contract.UnpackLog(event, "ValidatorsImportRequested", log); err != nil {
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

// ParseValidatorsImportRequested is a log parse operation binding the contract event 0xa6a406bba8edd318c95802e11c6fa6c696172b7bbf1aa90e28baa11a306bd3ca.
//
// Solidity: event ValidatorsImportRequested(bytes32 indexed txHash, uint256 count)
func (_L1Hub *L1HubFilterer) ParseValidatorsImportRequested(log types.Log) (*L1HubValidatorsImportRequested, error) {
	event := new(L1HubValidatorsImportRequested)
	if err := _L1Hub.contract.UnpackLog(event, "ValidatorsImportRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1HubWithdrawalFinalizedIterator is returned from FilterWithdrawalFinalized and is used to iterate over the raw logs and unpacked data for WithdrawalFinalized events raised by the L1Hub contract.
type L1HubWithdrawalFinalizedIterator struct {
	Event *L1HubWithdrawalFinalized // Event containing the contract specifics and raw log

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
func (it *L1HubWithdrawalFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1HubWithdrawalFinalized)
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
		it.Event = new(L1HubWithdrawalFinalized)
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
func (it *L1HubWithdrawalFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1HubWithdrawalFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1HubWithdrawalFinalized represents a WithdrawalFinalized event raised by the L1Hub contract.
type L1HubWithdrawalFinalized struct {
	ValidatorID   *big.Int
	ValidatorAddr common.Address
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalFinalized is a free log retrieval operation binding the contract event 0x2144673b243cbac0f18cd0de7b99b0d8e12caf36ff781e39524425d75ba4b273.
//
// Solidity: event WithdrawalFinalized(uint256 indexed validatorID, address indexed validatorAddr, uint256 amount)
func (_L1Hub *L1HubFilterer) FilterWithdrawalFinalized(opts *bind.FilterOpts, validatorID []*big.Int, validatorAddr []common.Address) (*L1HubWithdrawalFinalizedIterator, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}
	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _L1Hub.contract.FilterLogs(opts, "WithdrawalFinalized", validatorIDRule, validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return &L1HubWithdrawalFinalizedIterator{contract: _L1Hub.contract, event: "WithdrawalFinalized", logs: logs, sub: sub}, nil
}

// WatchWithdrawalFinalized is a free log subscription operation binding the contract event 0x2144673b243cbac0f18cd0de7b99b0d8e12caf36ff781e39524425d75ba4b273.
//
// Solidity: event WithdrawalFinalized(uint256 indexed validatorID, address indexed validatorAddr, uint256 amount)
func (_L1Hub *L1HubFilterer) WatchWithdrawalFinalized(opts *bind.WatchOpts, sink chan<- *L1HubWithdrawalFinalized, validatorID []*big.Int, validatorAddr []common.Address) (event.Subscription, error) {

	var validatorIDRule []interface{}
	for _, validatorIDItem := range validatorID {
		validatorIDRule = append(validatorIDRule, validatorIDItem)
	}
	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _L1Hub.contract.WatchLogs(opts, "WithdrawalFinalized", validatorIDRule, validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1HubWithdrawalFinalized)
				if err := _L1Hub.contract.UnpackLog(event, "WithdrawalFinalized", log); err != nil {
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

// ParseWithdrawalFinalized is a log parse operation binding the contract event 0x2144673b243cbac0f18cd0de7b99b0d8e12caf36ff781e39524425d75ba4b273.
//
// Solidity: event WithdrawalFinalized(uint256 indexed validatorID, address indexed validatorAddr, uint256 amount)
func (_L1Hub *L1HubFilterer) ParseWithdrawalFinalized(log types.Log) (*L1HubWithdrawalFinalized, error) {
	event := new(L1HubWithdrawalFinalized)
	if err := _L1Hub.contract.UnpackLog(event, "WithdrawalFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
