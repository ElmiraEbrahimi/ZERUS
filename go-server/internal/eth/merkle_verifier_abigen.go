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

// MerkleProofVerifierMetaData contains all meta data concerning the MerkleProofVerifier contract.
var MerkleProofVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"compressProof\",\"inputs\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"}],\"outputs\":[{\"name\":\"compressed\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyCompressedProof\",\"inputs\":[{\"name\":\"compressedProof\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"},{\"name\":\"input\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyProof\",\"inputs\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"input\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"ProofInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicInputNotInField\",\"inputs\":[]}]",
}

// MerkleProofVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleProofVerifierMetaData.ABI instead.
var MerkleProofVerifierABI = MerkleProofVerifierMetaData.ABI

// MerkleProofVerifier is an auto generated Go binding around an Ethereum contract.
type MerkleProofVerifier struct {
	MerkleProofVerifierCaller     // Read-only binding to the contract
	MerkleProofVerifierTransactor // Write-only binding to the contract
	MerkleProofVerifierFilterer   // Log filterer for contract events
}

// MerkleProofVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleProofVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleProofVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleProofVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleProofVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleProofVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleProofVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleProofVerifierSession struct {
	Contract     *MerkleProofVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MerkleProofVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleProofVerifierCallerSession struct {
	Contract *MerkleProofVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// MerkleProofVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleProofVerifierTransactorSession struct {
	Contract     *MerkleProofVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// MerkleProofVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleProofVerifierRaw struct {
	Contract *MerkleProofVerifier // Generic contract binding to access the raw methods on
}

// MerkleProofVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleProofVerifierCallerRaw struct {
	Contract *MerkleProofVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleProofVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleProofVerifierTransactorRaw struct {
	Contract *MerkleProofVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleProofVerifier creates a new instance of MerkleProofVerifier, bound to a specific deployed contract.
func NewMerkleProofVerifier(address common.Address, backend bind.ContractBackend) (*MerkleProofVerifier, error) {
	contract, err := bindMerkleProofVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleProofVerifier{MerkleProofVerifierCaller: MerkleProofVerifierCaller{contract: contract}, MerkleProofVerifierTransactor: MerkleProofVerifierTransactor{contract: contract}, MerkleProofVerifierFilterer: MerkleProofVerifierFilterer{contract: contract}}, nil
}

// NewMerkleProofVerifierCaller creates a new read-only instance of MerkleProofVerifier, bound to a specific deployed contract.
func NewMerkleProofVerifierCaller(address common.Address, caller bind.ContractCaller) (*MerkleProofVerifierCaller, error) {
	contract, err := bindMerkleProofVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleProofVerifierCaller{contract: contract}, nil
}

// NewMerkleProofVerifierTransactor creates a new write-only instance of MerkleProofVerifier, bound to a specific deployed contract.
func NewMerkleProofVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleProofVerifierTransactor, error) {
	contract, err := bindMerkleProofVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleProofVerifierTransactor{contract: contract}, nil
}

// NewMerkleProofVerifierFilterer creates a new log filterer instance of MerkleProofVerifier, bound to a specific deployed contract.
func NewMerkleProofVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleProofVerifierFilterer, error) {
	contract, err := bindMerkleProofVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleProofVerifierFilterer{contract: contract}, nil
}

// bindMerkleProofVerifier binds a generic wrapper to an already deployed contract.
func bindMerkleProofVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MerkleProofVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleProofVerifier *MerkleProofVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleProofVerifier.Contract.MerkleProofVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleProofVerifier *MerkleProofVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleProofVerifier.Contract.MerkleProofVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleProofVerifier *MerkleProofVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleProofVerifier.Contract.MerkleProofVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleProofVerifier *MerkleProofVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleProofVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleProofVerifier *MerkleProofVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleProofVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleProofVerifier *MerkleProofVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleProofVerifier.Contract.contract.Transact(opts, method, params...)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_MerkleProofVerifier *MerkleProofVerifierCaller) CompressProof(opts *bind.CallOpts, proof [8]*big.Int) ([4]*big.Int, error) {
	var out []interface{}
	err := _MerkleProofVerifier.contract.Call(opts, &out, "compressProof", proof)

	if err != nil {
		return *new([4]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([4]*big.Int)).(*[4]*big.Int)

	return out0, err

}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_MerkleProofVerifier *MerkleProofVerifierSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _MerkleProofVerifier.Contract.CompressProof(&_MerkleProofVerifier.CallOpts, proof)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_MerkleProofVerifier *MerkleProofVerifierCallerSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _MerkleProofVerifier.Contract.CompressProof(&_MerkleProofVerifier.CallOpts, proof)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xf2457c8d.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[4] input) view returns()
func (_MerkleProofVerifier *MerkleProofVerifierCaller) VerifyCompressedProof(opts *bind.CallOpts, compressedProof [4]*big.Int, input [4]*big.Int) error {
	var out []interface{}
	err := _MerkleProofVerifier.contract.Call(opts, &out, "verifyCompressedProof", compressedProof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xf2457c8d.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[4] input) view returns()
func (_MerkleProofVerifier *MerkleProofVerifierSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [4]*big.Int) error {
	return _MerkleProofVerifier.Contract.VerifyCompressedProof(&_MerkleProofVerifier.CallOpts, compressedProof, input)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xf2457c8d.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[4] input) view returns()
func (_MerkleProofVerifier *MerkleProofVerifierCallerSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [4]*big.Int) error {
	return _MerkleProofVerifier.Contract.VerifyCompressedProof(&_MerkleProofVerifier.CallOpts, compressedProof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x23572511.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[4] input) view returns()
func (_MerkleProofVerifier *MerkleProofVerifierCaller) VerifyProof(opts *bind.CallOpts, proof [8]*big.Int, input [4]*big.Int) error {
	var out []interface{}
	err := _MerkleProofVerifier.contract.Call(opts, &out, "verifyProof", proof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyProof is a free data retrieval call binding the contract method 0x23572511.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[4] input) view returns()
func (_MerkleProofVerifier *MerkleProofVerifierSession) VerifyProof(proof [8]*big.Int, input [4]*big.Int) error {
	return _MerkleProofVerifier.Contract.VerifyProof(&_MerkleProofVerifier.CallOpts, proof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x23572511.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[4] input) view returns()
func (_MerkleProofVerifier *MerkleProofVerifierCallerSession) VerifyProof(proof [8]*big.Int, input [4]*big.Int) error {
	return _MerkleProofVerifier.Contract.VerifyProof(&_MerkleProofVerifier.CallOpts, proof, input)
}
