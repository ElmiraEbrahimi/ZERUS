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

// VotingBatchVerifierMetaData contains all meta data concerning the VotingBatchVerifier contract.
var VotingBatchVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"compressProof\",\"inputs\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"}],\"outputs\":[{\"name\":\"compressed\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyCompressedProof\",\"inputs\":[{\"name\":\"compressedProof\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"},{\"name\":\"input\",\"type\":\"uint256[12]\",\"internalType\":\"uint256[12]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyProof\",\"inputs\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"input\",\"type\":\"uint256[12]\",\"internalType\":\"uint256[12]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"ProofInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicInputNotInField\",\"inputs\":[]}]",
}

// VotingBatchVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use VotingBatchVerifierMetaData.ABI instead.
var VotingBatchVerifierABI = VotingBatchVerifierMetaData.ABI

// VotingBatchVerifier is an auto generated Go binding around an Ethereum contract.
type VotingBatchVerifier struct {
	VotingBatchVerifierCaller     // Read-only binding to the contract
	VotingBatchVerifierTransactor // Write-only binding to the contract
	VotingBatchVerifierFilterer   // Log filterer for contract events
}

// VotingBatchVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VotingBatchVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingBatchVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VotingBatchVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingBatchVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VotingBatchVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingBatchVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VotingBatchVerifierSession struct {
	Contract     *VotingBatchVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// VotingBatchVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VotingBatchVerifierCallerSession struct {
	Contract *VotingBatchVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// VotingBatchVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VotingBatchVerifierTransactorSession struct {
	Contract     *VotingBatchVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// VotingBatchVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VotingBatchVerifierRaw struct {
	Contract *VotingBatchVerifier // Generic contract binding to access the raw methods on
}

// VotingBatchVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VotingBatchVerifierCallerRaw struct {
	Contract *VotingBatchVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VotingBatchVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VotingBatchVerifierTransactorRaw struct {
	Contract *VotingBatchVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVotingBatchVerifier creates a new instance of VotingBatchVerifier, bound to a specific deployed contract.
func NewVotingBatchVerifier(address common.Address, backend bind.ContractBackend) (*VotingBatchVerifier, error) {
	contract, err := bindVotingBatchVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VotingBatchVerifier{VotingBatchVerifierCaller: VotingBatchVerifierCaller{contract: contract}, VotingBatchVerifierTransactor: VotingBatchVerifierTransactor{contract: contract}, VotingBatchVerifierFilterer: VotingBatchVerifierFilterer{contract: contract}}, nil
}

// NewVotingBatchVerifierCaller creates a new read-only instance of VotingBatchVerifier, bound to a specific deployed contract.
func NewVotingBatchVerifierCaller(address common.Address, caller bind.ContractCaller) (*VotingBatchVerifierCaller, error) {
	contract, err := bindVotingBatchVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VotingBatchVerifierCaller{contract: contract}, nil
}

// NewVotingBatchVerifierTransactor creates a new write-only instance of VotingBatchVerifier, bound to a specific deployed contract.
func NewVotingBatchVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VotingBatchVerifierTransactor, error) {
	contract, err := bindVotingBatchVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VotingBatchVerifierTransactor{contract: contract}, nil
}

// NewVotingBatchVerifierFilterer creates a new log filterer instance of VotingBatchVerifier, bound to a specific deployed contract.
func NewVotingBatchVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VotingBatchVerifierFilterer, error) {
	contract, err := bindVotingBatchVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VotingBatchVerifierFilterer{contract: contract}, nil
}

// bindVotingBatchVerifier binds a generic wrapper to an already deployed contract.
func bindVotingBatchVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VotingBatchVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VotingBatchVerifier *VotingBatchVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VotingBatchVerifier.Contract.VotingBatchVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VotingBatchVerifier *VotingBatchVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VotingBatchVerifier.Contract.VotingBatchVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VotingBatchVerifier *VotingBatchVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VotingBatchVerifier.Contract.VotingBatchVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VotingBatchVerifier *VotingBatchVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VotingBatchVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VotingBatchVerifier *VotingBatchVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VotingBatchVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VotingBatchVerifier *VotingBatchVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VotingBatchVerifier.Contract.contract.Transact(opts, method, params...)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_VotingBatchVerifier *VotingBatchVerifierCaller) CompressProof(opts *bind.CallOpts, proof [8]*big.Int) ([4]*big.Int, error) {
	var out []interface{}
	err := _VotingBatchVerifier.contract.Call(opts, &out, "compressProof", proof)

	if err != nil {
		return *new([4]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([4]*big.Int)).(*[4]*big.Int)

	return out0, err

}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_VotingBatchVerifier *VotingBatchVerifierSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _VotingBatchVerifier.Contract.CompressProof(&_VotingBatchVerifier.CallOpts, proof)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_VotingBatchVerifier *VotingBatchVerifierCallerSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _VotingBatchVerifier.Contract.CompressProof(&_VotingBatchVerifier.CallOpts, proof)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xb28408ae.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[12] input) view returns()
func (_VotingBatchVerifier *VotingBatchVerifierCaller) VerifyCompressedProof(opts *bind.CallOpts, compressedProof [4]*big.Int, input [12]*big.Int) error {
	var out []interface{}
	err := _VotingBatchVerifier.contract.Call(opts, &out, "verifyCompressedProof", compressedProof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xb28408ae.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[12] input) view returns()
func (_VotingBatchVerifier *VotingBatchVerifierSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [12]*big.Int) error {
	return _VotingBatchVerifier.Contract.VerifyCompressedProof(&_VotingBatchVerifier.CallOpts, compressedProof, input)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xb28408ae.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[12] input) view returns()
func (_VotingBatchVerifier *VotingBatchVerifierCallerSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [12]*big.Int) error {
	return _VotingBatchVerifier.Contract.VerifyCompressedProof(&_VotingBatchVerifier.CallOpts, compressedProof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x8aa330f1.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[12] input) view returns()
func (_VotingBatchVerifier *VotingBatchVerifierCaller) VerifyProof(opts *bind.CallOpts, proof [8]*big.Int, input [12]*big.Int) error {
	var out []interface{}
	err := _VotingBatchVerifier.contract.Call(opts, &out, "verifyProof", proof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyProof is a free data retrieval call binding the contract method 0x8aa330f1.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[12] input) view returns()
func (_VotingBatchVerifier *VotingBatchVerifierSession) VerifyProof(proof [8]*big.Int, input [12]*big.Int) error {
	return _VotingBatchVerifier.Contract.VerifyProof(&_VotingBatchVerifier.CallOpts, proof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x8aa330f1.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[12] input) view returns()
func (_VotingBatchVerifier *VotingBatchVerifierCallerSession) VerifyProof(proof [8]*big.Int, input [12]*big.Int) error {
	return _VotingBatchVerifier.Contract.VerifyProof(&_VotingBatchVerifier.CallOpts, proof, input)
}
