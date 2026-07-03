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

// OracleAccount is an auto generated low-level Go binding around an user-defined struct.
type OracleAccount struct {
	Index   *big.Int
	PubKey  OraclePublicKey
	Balance *big.Int
}

// OraclePublicKey is an auto generated low-level Go binding around an user-defined struct.
type OraclePublicKey struct {
	X *big.Int
	Y *big.Int
}

// OracleReplacementRequest is an auto generated low-level Go binding around an user-defined struct.
type OracleReplacementRequest struct {
	RequestId            *big.Int
	TargetValidatorID    *big.Int
	CandidateValidatorID *big.Int
	CandidateAddr        common.Address
	CandidateStake       *big.Int
	CandidatePubKey      OraclePublicKey
	TargetAccount        OracleAccount
	Path                 []*big.Int
	LeafIndex            *big.Int
	Depth                *big.Int
}

// OracleValidatorInput is an auto generated low-level Go binding around an user-defined struct.
type OracleValidatorInput struct {
	ValidatorID   *big.Int
	ValidatorAddr common.Address
	Stake         *big.Int
	PubKey        OraclePublicKey
}

// OracleMetaData contains all meta data concerning the Oracle contract.
var OracleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_levels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_seedX\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_seedY\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"votingVerifierAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_batchSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_aggregatorTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"AGGREGATOR_REWARD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_LEVELS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VALIDATOR_REWARD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ZERO_VALUE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"aggregatorTimeout\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"chooseNewAggregator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"publicWitness\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nullifierHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitmentRootByEpoch\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"computeRootFromPath\",\"inputs\":[{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"leafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"exit\",\"inputs\":[{\"name\":\"account\",\"type\":\"tuple\",\"internalType\":\"structOracle.Account\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"leafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAggregator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLevels\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextLeafIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReward\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSeed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"tuple\",\"internalType\":\"structOracle.Account\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashLeftRight\",\"inputs\":[{\"name\":\"left\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"right\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"importValidatorsFromL1\",\"inputs\":[{\"name\":\"inputs\",\"type\":\"tuple[]\",\"internalType\":\"structOracle.ValidatorInput[]\",\"components\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isPublishedCommitmentRoot\",\"inputs\":[{\"name\":\"root\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l1Hub\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"leafSum\",\"inputs\":[{\"name\":\"leaf\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"levels\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"publishCommitmentRoot\",\"inputs\":[{\"name\":\"root\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"dfsRef\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUser\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publicKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"replace\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"toReplace\",\"type\":\"tuple\",\"internalType\":\"structOracle.Account\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"leafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"replaceValidatorFromL1\",\"inputs\":[{\"name\":\"request\",\"type\":\"tuple\",\"internalType\":\"structOracle.ReplacementRequest\",\"components\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"targetValidatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidateValidatorID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidateAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"candidateStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"candidatePubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"targetAccount\",\"type\":\"tuple\",\"internalType\":\"structOracle.Account\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"leafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setL1Hub\",\"inputs\":[{\"name\":\"l1HubAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"spentNullifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"submitWiVote\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"batchCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"honestBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vote\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"postStateRoot\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"postSeedX\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"postSeedY\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verify\",\"inputs\":[{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"leafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"viewBalance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"viewLatestCommitmentRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"viewLatestIPFSHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"account\",\"type\":\"tuple\",\"internalType\":\"structOracle.Account\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"tuple\",\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"path\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"leafIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"zeros\",\"inputs\":[{\"name\":\"i\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"BurnSubmitted\",\"inputs\":[{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClaimMinted\",\"inputs\":[{\"name\":\"uniqueID\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"nullifierHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClaimSubmitted\",\"inputs\":[{\"name\":\"uniqueID\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"publicWitness\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"nullifierHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitmentRootPublished\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"root\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"dfsRef\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Exiting\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L1HubUpdated\",\"inputs\":[{\"name\":\"l1Hub\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L2ToL1MessageSent\",\"inputs\":[{\"name\":\"msgType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumOracle.L2ToL1MsgType\"},{\"name\":\"payload\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewAggregator\",\"inputs\":[{\"name\":\"validatorID\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Registered\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Replaced\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"replaced\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"stake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplacementFromL1Processed\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"newRoot\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UserRegistered\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRegistered\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorID\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOracle.PublicKey\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"balance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorsImportedFromL1\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newRoot\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WiVoteSubmitted\",\"inputs\":[{\"name\":\"submitter\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"validators\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"honestBits\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"request\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"majorityVote\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleMetaData.ABI instead.
var OracleABI = OracleMetaData.ABI

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// AGGREGATORREWARD is a free data retrieval call binding the contract method 0xd97c0155.
//
// Solidity: function AGGREGATOR_REWARD() view returns(uint256)
func (_Oracle *OracleCaller) AGGREGATORREWARD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "AGGREGATOR_REWARD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AGGREGATORREWARD is a free data retrieval call binding the contract method 0xd97c0155.
//
// Solidity: function AGGREGATOR_REWARD() view returns(uint256)
func (_Oracle *OracleSession) AGGREGATORREWARD() (*big.Int, error) {
	return _Oracle.Contract.AGGREGATORREWARD(&_Oracle.CallOpts)
}

// AGGREGATORREWARD is a free data retrieval call binding the contract method 0xd97c0155.
//
// Solidity: function AGGREGATOR_REWARD() view returns(uint256)
func (_Oracle *OracleCallerSession) AGGREGATORREWARD() (*big.Int, error) {
	return _Oracle.Contract.AGGREGATORREWARD(&_Oracle.CallOpts)
}

// MAXLEVELS is a free data retrieval call binding the contract method 0x6702c7bf.
//
// Solidity: function MAX_LEVELS() view returns(uint256)
func (_Oracle *OracleCaller) MAXLEVELS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "MAX_LEVELS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXLEVELS is a free data retrieval call binding the contract method 0x6702c7bf.
//
// Solidity: function MAX_LEVELS() view returns(uint256)
func (_Oracle *OracleSession) MAXLEVELS() (*big.Int, error) {
	return _Oracle.Contract.MAXLEVELS(&_Oracle.CallOpts)
}

// MAXLEVELS is a free data retrieval call binding the contract method 0x6702c7bf.
//
// Solidity: function MAX_LEVELS() view returns(uint256)
func (_Oracle *OracleCallerSession) MAXLEVELS() (*big.Int, error) {
	return _Oracle.Contract.MAXLEVELS(&_Oracle.CallOpts)
}

// VALIDATORREWARD is a free data retrieval call binding the contract method 0xa3c3cda0.
//
// Solidity: function VALIDATOR_REWARD() view returns(uint256)
func (_Oracle *OracleCaller) VALIDATORREWARD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "VALIDATOR_REWARD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VALIDATORREWARD is a free data retrieval call binding the contract method 0xa3c3cda0.
//
// Solidity: function VALIDATOR_REWARD() view returns(uint256)
func (_Oracle *OracleSession) VALIDATORREWARD() (*big.Int, error) {
	return _Oracle.Contract.VALIDATORREWARD(&_Oracle.CallOpts)
}

// VALIDATORREWARD is a free data retrieval call binding the contract method 0xa3c3cda0.
//
// Solidity: function VALIDATOR_REWARD() view returns(uint256)
func (_Oracle *OracleCallerSession) VALIDATORREWARD() (*big.Int, error) {
	return _Oracle.Contract.VALIDATORREWARD(&_Oracle.CallOpts)
}

// ZEROVALUE is a free data retrieval call binding the contract method 0xec732959.
//
// Solidity: function ZERO_VALUE() view returns(uint256)
func (_Oracle *OracleCaller) ZEROVALUE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "ZERO_VALUE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ZEROVALUE is a free data retrieval call binding the contract method 0xec732959.
//
// Solidity: function ZERO_VALUE() view returns(uint256)
func (_Oracle *OracleSession) ZEROVALUE() (*big.Int, error) {
	return _Oracle.Contract.ZEROVALUE(&_Oracle.CallOpts)
}

// ZEROVALUE is a free data retrieval call binding the contract method 0xec732959.
//
// Solidity: function ZERO_VALUE() view returns(uint256)
func (_Oracle *OracleCallerSession) ZEROVALUE() (*big.Int, error) {
	return _Oracle.Contract.ZEROVALUE(&_Oracle.CallOpts)
}

// AggregatorTimeout is a free data retrieval call binding the contract method 0x0983fe4a.
//
// Solidity: function aggregatorTimeout() view returns(uint256)
func (_Oracle *OracleCaller) AggregatorTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "aggregatorTimeout")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AggregatorTimeout is a free data retrieval call binding the contract method 0x0983fe4a.
//
// Solidity: function aggregatorTimeout() view returns(uint256)
func (_Oracle *OracleSession) AggregatorTimeout() (*big.Int, error) {
	return _Oracle.Contract.AggregatorTimeout(&_Oracle.CallOpts)
}

// AggregatorTimeout is a free data retrieval call binding the contract method 0x0983fe4a.
//
// Solidity: function aggregatorTimeout() view returns(uint256)
func (_Oracle *OracleCallerSession) AggregatorTimeout() (*big.Int, error) {
	return _Oracle.Contract.AggregatorTimeout(&_Oracle.CallOpts)
}

// BatchSize is a free data retrieval call binding the contract method 0xf4daaba1.
//
// Solidity: function batchSize() view returns(uint256)
func (_Oracle *OracleCaller) BatchSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "batchSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BatchSize is a free data retrieval call binding the contract method 0xf4daaba1.
//
// Solidity: function batchSize() view returns(uint256)
func (_Oracle *OracleSession) BatchSize() (*big.Int, error) {
	return _Oracle.Contract.BatchSize(&_Oracle.CallOpts)
}

// BatchSize is a free data retrieval call binding the contract method 0xf4daaba1.
//
// Solidity: function batchSize() view returns(uint256)
func (_Oracle *OracleCallerSession) BatchSize() (*big.Int, error) {
	return _Oracle.Contract.BatchSize(&_Oracle.CallOpts)
}

// CommitmentRootByEpoch is a free data retrieval call binding the contract method 0xaa89f870.
//
// Solidity: function commitmentRootByEpoch(uint256 ) view returns(uint256)
func (_Oracle *OracleCaller) CommitmentRootByEpoch(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "commitmentRootByEpoch", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentRootByEpoch is a free data retrieval call binding the contract method 0xaa89f870.
//
// Solidity: function commitmentRootByEpoch(uint256 ) view returns(uint256)
func (_Oracle *OracleSession) CommitmentRootByEpoch(arg0 *big.Int) (*big.Int, error) {
	return _Oracle.Contract.CommitmentRootByEpoch(&_Oracle.CallOpts, arg0)
}

// CommitmentRootByEpoch is a free data retrieval call binding the contract method 0xaa89f870.
//
// Solidity: function commitmentRootByEpoch(uint256 ) view returns(uint256)
func (_Oracle *OracleCallerSession) CommitmentRootByEpoch(arg0 *big.Int) (*big.Int, error) {
	return _Oracle.Contract.CommitmentRootByEpoch(&_Oracle.CallOpts, arg0)
}

// ComputeRootFromPath is a free data retrieval call binding the contract method 0x8cfc5a3c.
//
// Solidity: function computeRootFromPath(uint256[] path, uint256 leafIndex, uint256 depth) pure returns(uint256)
func (_Oracle *OracleCaller) ComputeRootFromPath(opts *bind.CallOpts, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "computeRootFromPath", path, leafIndex, depth)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ComputeRootFromPath is a free data retrieval call binding the contract method 0x8cfc5a3c.
//
// Solidity: function computeRootFromPath(uint256[] path, uint256 leafIndex, uint256 depth) pure returns(uint256)
func (_Oracle *OracleSession) ComputeRootFromPath(path []*big.Int, leafIndex *big.Int, depth *big.Int) (*big.Int, error) {
	return _Oracle.Contract.ComputeRootFromPath(&_Oracle.CallOpts, path, leafIndex, depth)
}

// ComputeRootFromPath is a free data retrieval call binding the contract method 0x8cfc5a3c.
//
// Solidity: function computeRootFromPath(uint256[] path, uint256 leafIndex, uint256 depth) pure returns(uint256)
func (_Oracle *OracleCallerSession) ComputeRootFromPath(path []*big.Int, leafIndex *big.Int, depth *big.Int) (*big.Int, error) {
	return _Oracle.Contract.ComputeRootFromPath(&_Oracle.CallOpts, path, leafIndex, depth)
}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(uint256)
func (_Oracle *OracleCaller) GetAggregator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getAggregator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(uint256)
func (_Oracle *OracleSession) GetAggregator() (*big.Int, error) {
	return _Oracle.Contract.GetAggregator(&_Oracle.CallOpts)
}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(uint256)
func (_Oracle *OracleCallerSession) GetAggregator() (*big.Int, error) {
	return _Oracle.Contract.GetAggregator(&_Oracle.CallOpts)
}

// GetLevels is a free data retrieval call binding the contract method 0x0c394a60.
//
// Solidity: function getLevels() view returns(uint256)
func (_Oracle *OracleCaller) GetLevels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getLevels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLevels is a free data retrieval call binding the contract method 0x0c394a60.
//
// Solidity: function getLevels() view returns(uint256)
func (_Oracle *OracleSession) GetLevels() (*big.Int, error) {
	return _Oracle.Contract.GetLevels(&_Oracle.CallOpts)
}

// GetLevels is a free data retrieval call binding the contract method 0x0c394a60.
//
// Solidity: function getLevels() view returns(uint256)
func (_Oracle *OracleCallerSession) GetLevels() (*big.Int, error) {
	return _Oracle.Contract.GetLevels(&_Oracle.CallOpts)
}

// GetNextLeafIndex is a free data retrieval call binding the contract method 0x50e9b925.
//
// Solidity: function getNextLeafIndex() view returns(uint256)
func (_Oracle *OracleCaller) GetNextLeafIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getNextLeafIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNextLeafIndex is a free data retrieval call binding the contract method 0x50e9b925.
//
// Solidity: function getNextLeafIndex() view returns(uint256)
func (_Oracle *OracleSession) GetNextLeafIndex() (*big.Int, error) {
	return _Oracle.Contract.GetNextLeafIndex(&_Oracle.CallOpts)
}

// GetNextLeafIndex is a free data retrieval call binding the contract method 0x50e9b925.
//
// Solidity: function getNextLeafIndex() view returns(uint256)
func (_Oracle *OracleCallerSession) GetNextLeafIndex() (*big.Int, error) {
	return _Oracle.Contract.GetNextLeafIndex(&_Oracle.CallOpts)
}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_Oracle *OracleCaller) GetReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_Oracle *OracleSession) GetReward() (*big.Int, error) {
	return _Oracle.Contract.GetReward(&_Oracle.CallOpts)
}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_Oracle *OracleCallerSession) GetReward() (*big.Int, error) {
	return _Oracle.Contract.GetReward(&_Oracle.CallOpts)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(uint256)
func (_Oracle *OracleCaller) GetRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(uint256)
func (_Oracle *OracleSession) GetRoot() (*big.Int, error) {
	return _Oracle.Contract.GetRoot(&_Oracle.CallOpts)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(uint256)
func (_Oracle *OracleCallerSession) GetRoot() (*big.Int, error) {
	return _Oracle.Contract.GetRoot(&_Oracle.CallOpts)
}

// GetSeed is a free data retrieval call binding the contract method 0x39e7357c.
//
// Solidity: function getSeed() view returns(uint256, uint256)
func (_Oracle *OracleCaller) GetSeed(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getSeed")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetSeed is a free data retrieval call binding the contract method 0x39e7357c.
//
// Solidity: function getSeed() view returns(uint256, uint256)
func (_Oracle *OracleSession) GetSeed() (*big.Int, *big.Int, error) {
	return _Oracle.Contract.GetSeed(&_Oracle.CallOpts)
}

// GetSeed is a free data retrieval call binding the contract method 0x39e7357c.
//
// Solidity: function getSeed() view returns(uint256, uint256)
func (_Oracle *OracleCallerSession) GetSeed() (*big.Int, *big.Int, error) {
	return _Oracle.Contract.GetSeed(&_Oracle.CallOpts)
}

// HashAccount is a free data retrieval call binding the contract method 0xea368cff.
//
// Solidity: function hashAccount((uint256,(uint256,uint256),uint256) account) pure returns(uint256)
func (_Oracle *OracleCaller) HashAccount(opts *bind.CallOpts, account OracleAccount) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "hashAccount", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HashAccount is a free data retrieval call binding the contract method 0xea368cff.
//
// Solidity: function hashAccount((uint256,(uint256,uint256),uint256) account) pure returns(uint256)
func (_Oracle *OracleSession) HashAccount(account OracleAccount) (*big.Int, error) {
	return _Oracle.Contract.HashAccount(&_Oracle.CallOpts, account)
}

// HashAccount is a free data retrieval call binding the contract method 0xea368cff.
//
// Solidity: function hashAccount((uint256,(uint256,uint256),uint256) account) pure returns(uint256)
func (_Oracle *OracleCallerSession) HashAccount(account OracleAccount) (*big.Int, error) {
	return _Oracle.Contract.HashAccount(&_Oracle.CallOpts, account)
}

// HashLeftRight is a free data retrieval call binding the contract method 0x5bb93995.
//
// Solidity: function hashLeftRight(uint256 left, uint256 right) pure returns(uint256)
func (_Oracle *OracleCaller) HashLeftRight(opts *bind.CallOpts, left *big.Int, right *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "hashLeftRight", left, right)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HashLeftRight is a free data retrieval call binding the contract method 0x5bb93995.
//
// Solidity: function hashLeftRight(uint256 left, uint256 right) pure returns(uint256)
func (_Oracle *OracleSession) HashLeftRight(left *big.Int, right *big.Int) (*big.Int, error) {
	return _Oracle.Contract.HashLeftRight(&_Oracle.CallOpts, left, right)
}

// HashLeftRight is a free data retrieval call binding the contract method 0x5bb93995.
//
// Solidity: function hashLeftRight(uint256 left, uint256 right) pure returns(uint256)
func (_Oracle *OracleCallerSession) HashLeftRight(left *big.Int, right *big.Int) (*big.Int, error) {
	return _Oracle.Contract.HashLeftRight(&_Oracle.CallOpts, left, right)
}

// IsPublishedCommitmentRoot is a free data retrieval call binding the contract method 0x911c410e.
//
// Solidity: function isPublishedCommitmentRoot(uint256 root) view returns(bool)
func (_Oracle *OracleCaller) IsPublishedCommitmentRoot(opts *bind.CallOpts, root *big.Int) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "isPublishedCommitmentRoot", root)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPublishedCommitmentRoot is a free data retrieval call binding the contract method 0x911c410e.
//
// Solidity: function isPublishedCommitmentRoot(uint256 root) view returns(bool)
func (_Oracle *OracleSession) IsPublishedCommitmentRoot(root *big.Int) (bool, error) {
	return _Oracle.Contract.IsPublishedCommitmentRoot(&_Oracle.CallOpts, root)
}

// IsPublishedCommitmentRoot is a free data retrieval call binding the contract method 0x911c410e.
//
// Solidity: function isPublishedCommitmentRoot(uint256 root) view returns(bool)
func (_Oracle *OracleCallerSession) IsPublishedCommitmentRoot(root *big.Int) (bool, error) {
	return _Oracle.Contract.IsPublishedCommitmentRoot(&_Oracle.CallOpts, root)
}

// L1Hub is a free data retrieval call binding the contract method 0x43752ce9.
//
// Solidity: function l1Hub() view returns(address)
func (_Oracle *OracleCaller) L1Hub(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "l1Hub")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Hub is a free data retrieval call binding the contract method 0x43752ce9.
//
// Solidity: function l1Hub() view returns(address)
func (_Oracle *OracleSession) L1Hub() (common.Address, error) {
	return _Oracle.Contract.L1Hub(&_Oracle.CallOpts)
}

// L1Hub is a free data retrieval call binding the contract method 0x43752ce9.
//
// Solidity: function l1Hub() view returns(address)
func (_Oracle *OracleCallerSession) L1Hub() (common.Address, error) {
	return _Oracle.Contract.L1Hub(&_Oracle.CallOpts)
}

// LeafSum is a free data retrieval call binding the contract method 0x78014d33.
//
// Solidity: function leafSum(uint256 leaf) pure returns(uint256)
func (_Oracle *OracleCaller) LeafSum(opts *bind.CallOpts, leaf *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "leafSum", leaf)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafSum is a free data retrieval call binding the contract method 0x78014d33.
//
// Solidity: function leafSum(uint256 leaf) pure returns(uint256)
func (_Oracle *OracleSession) LeafSum(leaf *big.Int) (*big.Int, error) {
	return _Oracle.Contract.LeafSum(&_Oracle.CallOpts, leaf)
}

// LeafSum is a free data retrieval call binding the contract method 0x78014d33.
//
// Solidity: function leafSum(uint256 leaf) pure returns(uint256)
func (_Oracle *OracleCallerSession) LeafSum(leaf *big.Int) (*big.Int, error) {
	return _Oracle.Contract.LeafSum(&_Oracle.CallOpts, leaf)
}

// Levels is a free data retrieval call binding the contract method 0x4ecf518b.
//
// Solidity: function levels() view returns(uint256)
func (_Oracle *OracleCaller) Levels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "levels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Levels is a free data retrieval call binding the contract method 0x4ecf518b.
//
// Solidity: function levels() view returns(uint256)
func (_Oracle *OracleSession) Levels() (*big.Int, error) {
	return _Oracle.Contract.Levels(&_Oracle.CallOpts)
}

// Levels is a free data retrieval call binding the contract method 0x4ecf518b.
//
// Solidity: function levels() view returns(uint256)
func (_Oracle *OracleCallerSession) Levels() (*big.Int, error) {
	return _Oracle.Contract.Levels(&_Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleCallerSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

// SpentNullifiers is a free data retrieval call binding the contract method 0x1c70406a.
//
// Solidity: function spentNullifiers(bytes32 ) view returns(bool)
func (_Oracle *OracleCaller) SpentNullifiers(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "spentNullifiers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SpentNullifiers is a free data retrieval call binding the contract method 0x1c70406a.
//
// Solidity: function spentNullifiers(bytes32 ) view returns(bool)
func (_Oracle *OracleSession) SpentNullifiers(arg0 [32]byte) (bool, error) {
	return _Oracle.Contract.SpentNullifiers(&_Oracle.CallOpts, arg0)
}

// SpentNullifiers is a free data retrieval call binding the contract method 0x1c70406a.
//
// Solidity: function spentNullifiers(bytes32 ) view returns(bool)
func (_Oracle *OracleCallerSession) SpentNullifiers(arg0 [32]byte) (bool, error) {
	return _Oracle.Contract.SpentNullifiers(&_Oracle.CallOpts, arg0)
}

// Verify is a free data retrieval call binding the contract method 0x6c0642b0.
//
// Solidity: function verify(uint256[] path, uint256 leafIndex, uint256 depth) view returns(bool)
func (_Oracle *OracleCaller) Verify(opts *bind.CallOpts, path []*big.Int, leafIndex *big.Int, depth *big.Int) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "verify", path, leafIndex, depth)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0x6c0642b0.
//
// Solidity: function verify(uint256[] path, uint256 leafIndex, uint256 depth) view returns(bool)
func (_Oracle *OracleSession) Verify(path []*big.Int, leafIndex *big.Int, depth *big.Int) (bool, error) {
	return _Oracle.Contract.Verify(&_Oracle.CallOpts, path, leafIndex, depth)
}

// Verify is a free data retrieval call binding the contract method 0x6c0642b0.
//
// Solidity: function verify(uint256[] path, uint256 leafIndex, uint256 depth) view returns(bool)
func (_Oracle *OracleCallerSession) Verify(path []*big.Int, leafIndex *big.Int, depth *big.Int) (bool, error) {
	return _Oracle.Contract.Verify(&_Oracle.CallOpts, path, leafIndex, depth)
}

// ViewBalance is a free data retrieval call binding the contract method 0x3ff1e05b.
//
// Solidity: function viewBalance() view returns(uint256, uint256)
func (_Oracle *OracleCaller) ViewBalance(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "viewBalance")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// ViewBalance is a free data retrieval call binding the contract method 0x3ff1e05b.
//
// Solidity: function viewBalance() view returns(uint256, uint256)
func (_Oracle *OracleSession) ViewBalance() (*big.Int, *big.Int, error) {
	return _Oracle.Contract.ViewBalance(&_Oracle.CallOpts)
}

// ViewBalance is a free data retrieval call binding the contract method 0x3ff1e05b.
//
// Solidity: function viewBalance() view returns(uint256, uint256)
func (_Oracle *OracleCallerSession) ViewBalance() (*big.Int, *big.Int, error) {
	return _Oracle.Contract.ViewBalance(&_Oracle.CallOpts)
}

// ViewLatestCommitmentRoot is a free data retrieval call binding the contract method 0x2b66c943.
//
// Solidity: function viewLatestCommitmentRoot() view returns(uint256, uint256)
func (_Oracle *OracleCaller) ViewLatestCommitmentRoot(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "viewLatestCommitmentRoot")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// ViewLatestCommitmentRoot is a free data retrieval call binding the contract method 0x2b66c943.
//
// Solidity: function viewLatestCommitmentRoot() view returns(uint256, uint256)
func (_Oracle *OracleSession) ViewLatestCommitmentRoot() (*big.Int, *big.Int, error) {
	return _Oracle.Contract.ViewLatestCommitmentRoot(&_Oracle.CallOpts)
}

// ViewLatestCommitmentRoot is a free data retrieval call binding the contract method 0x2b66c943.
//
// Solidity: function viewLatestCommitmentRoot() view returns(uint256, uint256)
func (_Oracle *OracleCallerSession) ViewLatestCommitmentRoot() (*big.Int, *big.Int, error) {
	return _Oracle.Contract.ViewLatestCommitmentRoot(&_Oracle.CallOpts)
}

// ViewLatestIPFSHash is a free data retrieval call binding the contract method 0x1bede7cc.
//
// Solidity: function viewLatestIPFSHash() view returns(string)
func (_Oracle *OracleCaller) ViewLatestIPFSHash(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "viewLatestIPFSHash")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ViewLatestIPFSHash is a free data retrieval call binding the contract method 0x1bede7cc.
//
// Solidity: function viewLatestIPFSHash() view returns(string)
func (_Oracle *OracleSession) ViewLatestIPFSHash() (string, error) {
	return _Oracle.Contract.ViewLatestIPFSHash(&_Oracle.CallOpts)
}

// ViewLatestIPFSHash is a free data retrieval call binding the contract method 0x1bede7cc.
//
// Solidity: function viewLatestIPFSHash() view returns(string)
func (_Oracle *OracleCallerSession) ViewLatestIPFSHash() (string, error) {
	return _Oracle.Contract.ViewLatestIPFSHash(&_Oracle.CallOpts)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 i) pure returns(uint256)
func (_Oracle *OracleCaller) Zeros(opts *bind.CallOpts, i *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "zeros", i)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 i) pure returns(uint256)
func (_Oracle *OracleSession) Zeros(i *big.Int) (*big.Int, error) {
	return _Oracle.Contract.Zeros(&_Oracle.CallOpts, i)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 i) pure returns(uint256)
func (_Oracle *OracleCallerSession) Zeros(i *big.Int) (*big.Int, error) {
	return _Oracle.Contract.Zeros(&_Oracle.CallOpts, i)
}

// Burn is a paid mutator transaction binding the contract method 0x08a1eee1.
//
// Solidity: function burn(bytes32 commitmentHash) returns()
func (_Oracle *OracleTransactor) Burn(opts *bind.TransactOpts, commitmentHash [32]byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "burn", commitmentHash)
}

// Burn is a paid mutator transaction binding the contract method 0x08a1eee1.
//
// Solidity: function burn(bytes32 commitmentHash) returns()
func (_Oracle *OracleSession) Burn(commitmentHash [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.Burn(&_Oracle.TransactOpts, commitmentHash)
}

// Burn is a paid mutator transaction binding the contract method 0x08a1eee1.
//
// Solidity: function burn(bytes32 commitmentHash) returns()
func (_Oracle *OracleTransactorSession) Burn(commitmentHash [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.Burn(&_Oracle.TransactOpts, commitmentHash)
}

// ChooseNewAggregator is a paid mutator transaction binding the contract method 0xf9a86b3b.
//
// Solidity: function chooseNewAggregator() returns()
func (_Oracle *OracleTransactor) ChooseNewAggregator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "chooseNewAggregator")
}

// ChooseNewAggregator is a paid mutator transaction binding the contract method 0xf9a86b3b.
//
// Solidity: function chooseNewAggregator() returns()
func (_Oracle *OracleSession) ChooseNewAggregator() (*types.Transaction, error) {
	return _Oracle.Contract.ChooseNewAggregator(&_Oracle.TransactOpts)
}

// ChooseNewAggregator is a paid mutator transaction binding the contract method 0xf9a86b3b.
//
// Solidity: function chooseNewAggregator() returns()
func (_Oracle *OracleTransactorSession) ChooseNewAggregator() (*types.Transaction, error) {
	return _Oracle.Contract.ChooseNewAggregator(&_Oracle.TransactOpts)
}

// Claim is a paid mutator transaction binding the contract method 0xbfc4441a.
//
// Solidity: function claim(bytes proof, bytes publicWitness, bytes32 nullifierHash) returns()
func (_Oracle *OracleTransactor) Claim(opts *bind.TransactOpts, proof []byte, publicWitness []byte, nullifierHash [32]byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "claim", proof, publicWitness, nullifierHash)
}

// Claim is a paid mutator transaction binding the contract method 0xbfc4441a.
//
// Solidity: function claim(bytes proof, bytes publicWitness, bytes32 nullifierHash) returns()
func (_Oracle *OracleSession) Claim(proof []byte, publicWitness []byte, nullifierHash [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.Claim(&_Oracle.TransactOpts, proof, publicWitness, nullifierHash)
}

// Claim is a paid mutator transaction binding the contract method 0xbfc4441a.
//
// Solidity: function claim(bytes proof, bytes publicWitness, bytes32 nullifierHash) returns()
func (_Oracle *OracleTransactorSession) Claim(proof []byte, publicWitness []byte, nullifierHash [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.Claim(&_Oracle.TransactOpts, proof, publicWitness, nullifierHash)
}

// Exit is a paid mutator transaction binding the contract method 0x4c5a04f4.
//
// Solidity: function exit((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256 leafIndex, uint256 depth) returns()
func (_Oracle *OracleTransactor) Exit(opts *bind.TransactOpts, account OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "exit", account, path, leafIndex, depth)
}

// Exit is a paid mutator transaction binding the contract method 0x4c5a04f4.
//
// Solidity: function exit((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256 leafIndex, uint256 depth) returns()
func (_Oracle *OracleSession) Exit(account OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Exit(&_Oracle.TransactOpts, account, path, leafIndex, depth)
}

// Exit is a paid mutator transaction binding the contract method 0x4c5a04f4.
//
// Solidity: function exit((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256 leafIndex, uint256 depth) returns()
func (_Oracle *OracleTransactorSession) Exit(account OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Exit(&_Oracle.TransactOpts, account, path, leafIndex, depth)
}

// ImportValidatorsFromL1 is a paid mutator transaction binding the contract method 0x958ecd48.
//
// Solidity: function importValidatorsFromL1((uint256,address,uint256,(uint256,uint256))[] inputs) returns()
func (_Oracle *OracleTransactor) ImportValidatorsFromL1(opts *bind.TransactOpts, inputs []OracleValidatorInput) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "importValidatorsFromL1", inputs)
}

// ImportValidatorsFromL1 is a paid mutator transaction binding the contract method 0x958ecd48.
//
// Solidity: function importValidatorsFromL1((uint256,address,uint256,(uint256,uint256))[] inputs) returns()
func (_Oracle *OracleSession) ImportValidatorsFromL1(inputs []OracleValidatorInput) (*types.Transaction, error) {
	return _Oracle.Contract.ImportValidatorsFromL1(&_Oracle.TransactOpts, inputs)
}

// ImportValidatorsFromL1 is a paid mutator transaction binding the contract method 0x958ecd48.
//
// Solidity: function importValidatorsFromL1((uint256,address,uint256,(uint256,uint256))[] inputs) returns()
func (_Oracle *OracleTransactorSession) ImportValidatorsFromL1(inputs []OracleValidatorInput) (*types.Transaction, error) {
	return _Oracle.Contract.ImportValidatorsFromL1(&_Oracle.TransactOpts, inputs)
}

// PublishCommitmentRoot is a paid mutator transaction binding the contract method 0x12dfaac4.
//
// Solidity: function publishCommitmentRoot(uint256 root, string dfsRef) returns()
func (_Oracle *OracleTransactor) PublishCommitmentRoot(opts *bind.TransactOpts, root *big.Int, dfsRef string) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "publishCommitmentRoot", root, dfsRef)
}

// PublishCommitmentRoot is a paid mutator transaction binding the contract method 0x12dfaac4.
//
// Solidity: function publishCommitmentRoot(uint256 root, string dfsRef) returns()
func (_Oracle *OracleSession) PublishCommitmentRoot(root *big.Int, dfsRef string) (*types.Transaction, error) {
	return _Oracle.Contract.PublishCommitmentRoot(&_Oracle.TransactOpts, root, dfsRef)
}

// PublishCommitmentRoot is a paid mutator transaction binding the contract method 0x12dfaac4.
//
// Solidity: function publishCommitmentRoot(uint256 root, string dfsRef) returns()
func (_Oracle *OracleTransactorSession) PublishCommitmentRoot(root *big.Int, dfsRef string) (*types.Transaction, error) {
	return _Oracle.Contract.PublishCommitmentRoot(&_Oracle.TransactOpts, root, dfsRef)
}

// RegisterUser is a paid mutator transaction binding the contract method 0xe90c3796.
//
// Solidity: function registerUser((uint256,uint256) publicKey) payable returns()
func (_Oracle *OracleTransactor) RegisterUser(opts *bind.TransactOpts, publicKey OraclePublicKey) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "registerUser", publicKey)
}

// RegisterUser is a paid mutator transaction binding the contract method 0xe90c3796.
//
// Solidity: function registerUser((uint256,uint256) publicKey) payable returns()
func (_Oracle *OracleSession) RegisterUser(publicKey OraclePublicKey) (*types.Transaction, error) {
	return _Oracle.Contract.RegisterUser(&_Oracle.TransactOpts, publicKey)
}

// RegisterUser is a paid mutator transaction binding the contract method 0xe90c3796.
//
// Solidity: function registerUser((uint256,uint256) publicKey) payable returns()
func (_Oracle *OracleTransactorSession) RegisterUser(publicKey OraclePublicKey) (*types.Transaction, error) {
	return _Oracle.Contract.RegisterUser(&_Oracle.TransactOpts, publicKey)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x7e56f05f.
//
// Solidity: function registerValidator(uint256 validatorID, (uint256,uint256) publicKey) payable returns()
func (_Oracle *OracleTransactor) RegisterValidator(opts *bind.TransactOpts, validatorID *big.Int, publicKey OraclePublicKey) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "registerValidator", validatorID, publicKey)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x7e56f05f.
//
// Solidity: function registerValidator(uint256 validatorID, (uint256,uint256) publicKey) payable returns()
func (_Oracle *OracleSession) RegisterValidator(validatorID *big.Int, publicKey OraclePublicKey) (*types.Transaction, error) {
	return _Oracle.Contract.RegisterValidator(&_Oracle.TransactOpts, validatorID, publicKey)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x7e56f05f.
//
// Solidity: function registerValidator(uint256 validatorID, (uint256,uint256) publicKey) payable returns()
func (_Oracle *OracleTransactorSession) RegisterValidator(validatorID *big.Int, publicKey OraclePublicKey) (*types.Transaction, error) {
	return _Oracle.Contract.RegisterValidator(&_Oracle.TransactOpts, validatorID, publicKey)
}

// Replace is a paid mutator transaction binding the contract method 0x622f2400.
//
// Solidity: function replace((uint256,uint256) publicKey, (uint256,(uint256,uint256),uint256) toReplace, uint256[] path, uint256 leafIndex, uint256 depth) payable returns()
func (_Oracle *OracleTransactor) Replace(opts *bind.TransactOpts, publicKey OraclePublicKey, toReplace OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "replace", publicKey, toReplace, path, leafIndex, depth)
}

// Replace is a paid mutator transaction binding the contract method 0x622f2400.
//
// Solidity: function replace((uint256,uint256) publicKey, (uint256,(uint256,uint256),uint256) toReplace, uint256[] path, uint256 leafIndex, uint256 depth) payable returns()
func (_Oracle *OracleSession) Replace(publicKey OraclePublicKey, toReplace OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Replace(&_Oracle.TransactOpts, publicKey, toReplace, path, leafIndex, depth)
}

// Replace is a paid mutator transaction binding the contract method 0x622f2400.
//
// Solidity: function replace((uint256,uint256) publicKey, (uint256,(uint256,uint256),uint256) toReplace, uint256[] path, uint256 leafIndex, uint256 depth) payable returns()
func (_Oracle *OracleTransactorSession) Replace(publicKey OraclePublicKey, toReplace OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Replace(&_Oracle.TransactOpts, publicKey, toReplace, path, leafIndex, depth)
}

// ReplaceValidatorFromL1 is a paid mutator transaction binding the contract method 0x2f886a4a.
//
// Solidity: function replaceValidatorFromL1((uint256,uint256,uint256,address,uint256,(uint256,uint256),(uint256,(uint256,uint256),uint256),uint256[],uint256,uint256) request) returns()
func (_Oracle *OracleTransactor) ReplaceValidatorFromL1(opts *bind.TransactOpts, request OracleReplacementRequest) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "replaceValidatorFromL1", request)
}

// ReplaceValidatorFromL1 is a paid mutator transaction binding the contract method 0x2f886a4a.
//
// Solidity: function replaceValidatorFromL1((uint256,uint256,uint256,address,uint256,(uint256,uint256),(uint256,(uint256,uint256),uint256),uint256[],uint256,uint256) request) returns()
func (_Oracle *OracleSession) ReplaceValidatorFromL1(request OracleReplacementRequest) (*types.Transaction, error) {
	return _Oracle.Contract.ReplaceValidatorFromL1(&_Oracle.TransactOpts, request)
}

// ReplaceValidatorFromL1 is a paid mutator transaction binding the contract method 0x2f886a4a.
//
// Solidity: function replaceValidatorFromL1((uint256,uint256,uint256,address,uint256,(uint256,uint256),(uint256,(uint256,uint256),uint256),uint256[],uint256,uint256) request) returns()
func (_Oracle *OracleTransactorSession) ReplaceValidatorFromL1(request OracleReplacementRequest) (*types.Transaction, error) {
	return _Oracle.Contract.ReplaceValidatorFromL1(&_Oracle.TransactOpts, request)
}

// SetL1Hub is a paid mutator transaction binding the contract method 0xd97178f2.
//
// Solidity: function setL1Hub(address l1HubAddress) returns()
func (_Oracle *OracleTransactor) SetL1Hub(opts *bind.TransactOpts, l1HubAddress common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setL1Hub", l1HubAddress)
}

// SetL1Hub is a paid mutator transaction binding the contract method 0xd97178f2.
//
// Solidity: function setL1Hub(address l1HubAddress) returns()
func (_Oracle *OracleSession) SetL1Hub(l1HubAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetL1Hub(&_Oracle.TransactOpts, l1HubAddress)
}

// SetL1Hub is a paid mutator transaction binding the contract method 0xd97178f2.
//
// Solidity: function setL1Hub(address l1HubAddress) returns()
func (_Oracle *OracleTransactorSession) SetL1Hub(l1HubAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetL1Hub(&_Oracle.TransactOpts, l1HubAddress)
}

// SubmitWiVote is a paid mutator transaction binding the contract method 0xdbc1dd69.
//
// Solidity: function submitWiVote(uint256 index, uint256 roundId, uint256 batchCommitment, uint256 validatorBits, uint256 honestBits, uint256 vote, uint256 postStateRoot, uint256 postSeedX, uint256 postSeedY, uint256[8] proof) returns()
func (_Oracle *OracleTransactor) SubmitWiVote(opts *bind.TransactOpts, index *big.Int, roundId *big.Int, batchCommitment *big.Int, validatorBits *big.Int, honestBits *big.Int, vote *big.Int, postStateRoot *big.Int, postSeedX *big.Int, postSeedY *big.Int, proof [8]*big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "submitWiVote", index, roundId, batchCommitment, validatorBits, honestBits, vote, postStateRoot, postSeedX, postSeedY, proof)
}

// SubmitWiVote is a paid mutator transaction binding the contract method 0xdbc1dd69.
//
// Solidity: function submitWiVote(uint256 index, uint256 roundId, uint256 batchCommitment, uint256 validatorBits, uint256 honestBits, uint256 vote, uint256 postStateRoot, uint256 postSeedX, uint256 postSeedY, uint256[8] proof) returns()
func (_Oracle *OracleSession) SubmitWiVote(index *big.Int, roundId *big.Int, batchCommitment *big.Int, validatorBits *big.Int, honestBits *big.Int, vote *big.Int, postStateRoot *big.Int, postSeedX *big.Int, postSeedY *big.Int, proof [8]*big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.SubmitWiVote(&_Oracle.TransactOpts, index, roundId, batchCommitment, validatorBits, honestBits, vote, postStateRoot, postSeedX, postSeedY, proof)
}

// SubmitWiVote is a paid mutator transaction binding the contract method 0xdbc1dd69.
//
// Solidity: function submitWiVote(uint256 index, uint256 roundId, uint256 batchCommitment, uint256 validatorBits, uint256 honestBits, uint256 vote, uint256 postStateRoot, uint256 postSeedX, uint256 postSeedY, uint256[8] proof) returns()
func (_Oracle *OracleTransactorSession) SubmitWiVote(index *big.Int, roundId *big.Int, batchCommitment *big.Int, validatorBits *big.Int, honestBits *big.Int, vote *big.Int, postStateRoot *big.Int, postSeedX *big.Int, postSeedY *big.Int, proof [8]*big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.SubmitWiVote(&_Oracle.TransactOpts, index, roundId, batchCommitment, validatorBits, honestBits, vote, postStateRoot, postSeedX, postSeedY, proof)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0218350a.
//
// Solidity: function withdraw((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256 leafIndex, uint256 depth) returns()
func (_Oracle *OracleTransactor) Withdraw(opts *bind.TransactOpts, account OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "withdraw", account, path, leafIndex, depth)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0218350a.
//
// Solidity: function withdraw((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256 leafIndex, uint256 depth) returns()
func (_Oracle *OracleSession) Withdraw(account OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Withdraw(&_Oracle.TransactOpts, account, path, leafIndex, depth)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0218350a.
//
// Solidity: function withdraw((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256 leafIndex, uint256 depth) returns()
func (_Oracle *OracleTransactorSession) Withdraw(account OracleAccount, path []*big.Int, leafIndex *big.Int, depth *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Withdraw(&_Oracle.TransactOpts, account, path, leafIndex, depth)
}

// OracleBurnSubmittedIterator is returned from FilterBurnSubmitted and is used to iterate over the raw logs and unpacked data for BurnSubmitted events raised by the Oracle contract.
type OracleBurnSubmittedIterator struct {
	Event *OracleBurnSubmitted // Event containing the contract specifics and raw log

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
func (it *OracleBurnSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleBurnSubmitted)
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
		it.Event = new(OracleBurnSubmitted)
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
func (it *OracleBurnSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleBurnSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleBurnSubmitted represents a BurnSubmitted event raised by the Oracle contract.
type OracleBurnSubmitted struct {
	CommitmentHash [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBurnSubmitted is a free log retrieval operation binding the contract event 0xb82baf1380ae33e09102fa90a0334892bb7395766e1a7adc29550d6e70af2b28.
//
// Solidity: event BurnSubmitted(bytes32 commitmentHash)
func (_Oracle *OracleFilterer) FilterBurnSubmitted(opts *bind.FilterOpts) (*OracleBurnSubmittedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "BurnSubmitted")
	if err != nil {
		return nil, err
	}
	return &OracleBurnSubmittedIterator{contract: _Oracle.contract, event: "BurnSubmitted", logs: logs, sub: sub}, nil
}

// WatchBurnSubmitted is a free log subscription operation binding the contract event 0xb82baf1380ae33e09102fa90a0334892bb7395766e1a7adc29550d6e70af2b28.
//
// Solidity: event BurnSubmitted(bytes32 commitmentHash)
func (_Oracle *OracleFilterer) WatchBurnSubmitted(opts *bind.WatchOpts, sink chan<- *OracleBurnSubmitted) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "BurnSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleBurnSubmitted)
				if err := _Oracle.contract.UnpackLog(event, "BurnSubmitted", log); err != nil {
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

// ParseBurnSubmitted is a log parse operation binding the contract event 0xb82baf1380ae33e09102fa90a0334892bb7395766e1a7adc29550d6e70af2b28.
//
// Solidity: event BurnSubmitted(bytes32 commitmentHash)
func (_Oracle *OracleFilterer) ParseBurnSubmitted(log types.Log) (*OracleBurnSubmitted, error) {
	event := new(OracleBurnSubmitted)
	if err := _Oracle.contract.UnpackLog(event, "BurnSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleClaimMintedIterator is returned from FilterClaimMinted and is used to iterate over the raw logs and unpacked data for ClaimMinted events raised by the Oracle contract.
type OracleClaimMintedIterator struct {
	Event *OracleClaimMinted // Event containing the contract specifics and raw log

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
func (it *OracleClaimMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleClaimMinted)
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
		it.Event = new(OracleClaimMinted)
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
func (it *OracleClaimMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleClaimMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleClaimMinted represents a ClaimMinted event raised by the Oracle contract.
type OracleClaimMinted struct {
	UniqueID      *big.Int
	Recipient     common.Address
	NullifierHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterClaimMinted is a free log retrieval operation binding the contract event 0x7572a3ce62c551cbc2492ffef7cef6cec05d63659999bed312a093c1f0492831.
//
// Solidity: event ClaimMinted(uint256 uniqueID, address recipient, bytes32 nullifierHash)
func (_Oracle *OracleFilterer) FilterClaimMinted(opts *bind.FilterOpts) (*OracleClaimMintedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "ClaimMinted")
	if err != nil {
		return nil, err
	}
	return &OracleClaimMintedIterator{contract: _Oracle.contract, event: "ClaimMinted", logs: logs, sub: sub}, nil
}

// WatchClaimMinted is a free log subscription operation binding the contract event 0x7572a3ce62c551cbc2492ffef7cef6cec05d63659999bed312a093c1f0492831.
//
// Solidity: event ClaimMinted(uint256 uniqueID, address recipient, bytes32 nullifierHash)
func (_Oracle *OracleFilterer) WatchClaimMinted(opts *bind.WatchOpts, sink chan<- *OracleClaimMinted) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "ClaimMinted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleClaimMinted)
				if err := _Oracle.contract.UnpackLog(event, "ClaimMinted", log); err != nil {
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

// ParseClaimMinted is a log parse operation binding the contract event 0x7572a3ce62c551cbc2492ffef7cef6cec05d63659999bed312a093c1f0492831.
//
// Solidity: event ClaimMinted(uint256 uniqueID, address recipient, bytes32 nullifierHash)
func (_Oracle *OracleFilterer) ParseClaimMinted(log types.Log) (*OracleClaimMinted, error) {
	event := new(OracleClaimMinted)
	if err := _Oracle.contract.UnpackLog(event, "ClaimMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleClaimSubmittedIterator is returned from FilterClaimSubmitted and is used to iterate over the raw logs and unpacked data for ClaimSubmitted events raised by the Oracle contract.
type OracleClaimSubmittedIterator struct {
	Event *OracleClaimSubmitted // Event containing the contract specifics and raw log

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
func (it *OracleClaimSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleClaimSubmitted)
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
		it.Event = new(OracleClaimSubmitted)
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
func (it *OracleClaimSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleClaimSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleClaimSubmitted represents a ClaimSubmitted event raised by the Oracle contract.
type OracleClaimSubmitted struct {
	UniqueID      *big.Int
	Proof         []byte
	PublicWitness []byte
	NullifierHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterClaimSubmitted is a free log retrieval operation binding the contract event 0x72c99a852ec54aa83e97e78ddb722293e43867a5e36b2446dfd80e0ae266157c.
//
// Solidity: event ClaimSubmitted(uint256 uniqueID, bytes proof, bytes publicWitness, bytes32 nullifierHash)
func (_Oracle *OracleFilterer) FilterClaimSubmitted(opts *bind.FilterOpts) (*OracleClaimSubmittedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "ClaimSubmitted")
	if err != nil {
		return nil, err
	}
	return &OracleClaimSubmittedIterator{contract: _Oracle.contract, event: "ClaimSubmitted", logs: logs, sub: sub}, nil
}

// WatchClaimSubmitted is a free log subscription operation binding the contract event 0x72c99a852ec54aa83e97e78ddb722293e43867a5e36b2446dfd80e0ae266157c.
//
// Solidity: event ClaimSubmitted(uint256 uniqueID, bytes proof, bytes publicWitness, bytes32 nullifierHash)
func (_Oracle *OracleFilterer) WatchClaimSubmitted(opts *bind.WatchOpts, sink chan<- *OracleClaimSubmitted) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "ClaimSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleClaimSubmitted)
				if err := _Oracle.contract.UnpackLog(event, "ClaimSubmitted", log); err != nil {
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

// ParseClaimSubmitted is a log parse operation binding the contract event 0x72c99a852ec54aa83e97e78ddb722293e43867a5e36b2446dfd80e0ae266157c.
//
// Solidity: event ClaimSubmitted(uint256 uniqueID, bytes proof, bytes publicWitness, bytes32 nullifierHash)
func (_Oracle *OracleFilterer) ParseClaimSubmitted(log types.Log) (*OracleClaimSubmitted, error) {
	event := new(OracleClaimSubmitted)
	if err := _Oracle.contract.UnpackLog(event, "ClaimSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleCommitmentRootPublishedIterator is returned from FilterCommitmentRootPublished and is used to iterate over the raw logs and unpacked data for CommitmentRootPublished events raised by the Oracle contract.
type OracleCommitmentRootPublishedIterator struct {
	Event *OracleCommitmentRootPublished // Event containing the contract specifics and raw log

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
func (it *OracleCommitmentRootPublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleCommitmentRootPublished)
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
		it.Event = new(OracleCommitmentRootPublished)
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
func (it *OracleCommitmentRootPublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleCommitmentRootPublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleCommitmentRootPublished represents a CommitmentRootPublished event raised by the Oracle contract.
type OracleCommitmentRootPublished struct {
	Epoch  *big.Int
	Root   *big.Int
	DfsRef string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCommitmentRootPublished is a free log retrieval operation binding the contract event 0x23f1667a6a7b338cbfcb35db53f683b1f3a9342d577bef2e197582dd8fc5d612.
//
// Solidity: event CommitmentRootPublished(uint256 indexed epoch, uint256 root, string dfsRef)
func (_Oracle *OracleFilterer) FilterCommitmentRootPublished(opts *bind.FilterOpts, epoch []*big.Int) (*OracleCommitmentRootPublishedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "CommitmentRootPublished", epochRule)
	if err != nil {
		return nil, err
	}
	return &OracleCommitmentRootPublishedIterator{contract: _Oracle.contract, event: "CommitmentRootPublished", logs: logs, sub: sub}, nil
}

// WatchCommitmentRootPublished is a free log subscription operation binding the contract event 0x23f1667a6a7b338cbfcb35db53f683b1f3a9342d577bef2e197582dd8fc5d612.
//
// Solidity: event CommitmentRootPublished(uint256 indexed epoch, uint256 root, string dfsRef)
func (_Oracle *OracleFilterer) WatchCommitmentRootPublished(opts *bind.WatchOpts, sink chan<- *OracleCommitmentRootPublished, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "CommitmentRootPublished", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleCommitmentRootPublished)
				if err := _Oracle.contract.UnpackLog(event, "CommitmentRootPublished", log); err != nil {
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

// ParseCommitmentRootPublished is a log parse operation binding the contract event 0x23f1667a6a7b338cbfcb35db53f683b1f3a9342d577bef2e197582dd8fc5d612.
//
// Solidity: event CommitmentRootPublished(uint256 indexed epoch, uint256 root, string dfsRef)
func (_Oracle *OracleFilterer) ParseCommitmentRootPublished(log types.Log) (*OracleCommitmentRootPublished, error) {
	event := new(OracleCommitmentRootPublished)
	if err := _Oracle.contract.UnpackLog(event, "CommitmentRootPublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleExitingIterator is returned from FilterExiting and is used to iterate over the raw logs and unpacked data for Exiting events raised by the Oracle contract.
type OracleExitingIterator struct {
	Event *OracleExiting // Event containing the contract specifics and raw log

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
func (it *OracleExitingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleExiting)
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
		it.Event = new(OracleExiting)
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
func (it *OracleExitingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleExitingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleExiting represents a Exiting event raised by the Oracle contract.
type OracleExiting struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExiting is a free log retrieval operation binding the contract event 0x157b1813744b6e99f33bf9153540dd45bd711cbb7322dc3b9c43822687e94180.
//
// Solidity: event Exiting(address indexed sender)
func (_Oracle *OracleFilterer) FilterExiting(opts *bind.FilterOpts, sender []common.Address) (*OracleExitingIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Exiting", senderRule)
	if err != nil {
		return nil, err
	}
	return &OracleExitingIterator{contract: _Oracle.contract, event: "Exiting", logs: logs, sub: sub}, nil
}

// WatchExiting is a free log subscription operation binding the contract event 0x157b1813744b6e99f33bf9153540dd45bd711cbb7322dc3b9c43822687e94180.
//
// Solidity: event Exiting(address indexed sender)
func (_Oracle *OracleFilterer) WatchExiting(opts *bind.WatchOpts, sink chan<- *OracleExiting, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Exiting", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleExiting)
				if err := _Oracle.contract.UnpackLog(event, "Exiting", log); err != nil {
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

// ParseExiting is a log parse operation binding the contract event 0x157b1813744b6e99f33bf9153540dd45bd711cbb7322dc3b9c43822687e94180.
//
// Solidity: event Exiting(address indexed sender)
func (_Oracle *OracleFilterer) ParseExiting(log types.Log) (*OracleExiting, error) {
	event := new(OracleExiting)
	if err := _Oracle.contract.UnpackLog(event, "Exiting", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleL1HubUpdatedIterator is returned from FilterL1HubUpdated and is used to iterate over the raw logs and unpacked data for L1HubUpdated events raised by the Oracle contract.
type OracleL1HubUpdatedIterator struct {
	Event *OracleL1HubUpdated // Event containing the contract specifics and raw log

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
func (it *OracleL1HubUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleL1HubUpdated)
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
		it.Event = new(OracleL1HubUpdated)
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
func (it *OracleL1HubUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleL1HubUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleL1HubUpdated represents a L1HubUpdated event raised by the Oracle contract.
type OracleL1HubUpdated struct {
	L1Hub common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterL1HubUpdated is a free log retrieval operation binding the contract event 0x5c8ac826140120ac2c42eed2cb0cbd58a8ff043284a66015263de4d870bae35a.
//
// Solidity: event L1HubUpdated(address indexed l1Hub)
func (_Oracle *OracleFilterer) FilterL1HubUpdated(opts *bind.FilterOpts, l1Hub []common.Address) (*OracleL1HubUpdatedIterator, error) {

	var l1HubRule []interface{}
	for _, l1HubItem := range l1Hub {
		l1HubRule = append(l1HubRule, l1HubItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "L1HubUpdated", l1HubRule)
	if err != nil {
		return nil, err
	}
	return &OracleL1HubUpdatedIterator{contract: _Oracle.contract, event: "L1HubUpdated", logs: logs, sub: sub}, nil
}

// WatchL1HubUpdated is a free log subscription operation binding the contract event 0x5c8ac826140120ac2c42eed2cb0cbd58a8ff043284a66015263de4d870bae35a.
//
// Solidity: event L1HubUpdated(address indexed l1Hub)
func (_Oracle *OracleFilterer) WatchL1HubUpdated(opts *bind.WatchOpts, sink chan<- *OracleL1HubUpdated, l1Hub []common.Address) (event.Subscription, error) {

	var l1HubRule []interface{}
	for _, l1HubItem := range l1Hub {
		l1HubRule = append(l1HubRule, l1HubItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "L1HubUpdated", l1HubRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleL1HubUpdated)
				if err := _Oracle.contract.UnpackLog(event, "L1HubUpdated", log); err != nil {
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

// ParseL1HubUpdated is a log parse operation binding the contract event 0x5c8ac826140120ac2c42eed2cb0cbd58a8ff043284a66015263de4d870bae35a.
//
// Solidity: event L1HubUpdated(address indexed l1Hub)
func (_Oracle *OracleFilterer) ParseL1HubUpdated(log types.Log) (*OracleL1HubUpdated, error) {
	event := new(OracleL1HubUpdated)
	if err := _Oracle.contract.UnpackLog(event, "L1HubUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleL2ToL1MessageSentIterator is returned from FilterL2ToL1MessageSent and is used to iterate over the raw logs and unpacked data for L2ToL1MessageSent events raised by the Oracle contract.
type OracleL2ToL1MessageSentIterator struct {
	Event *OracleL2ToL1MessageSent // Event containing the contract specifics and raw log

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
func (it *OracleL2ToL1MessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleL2ToL1MessageSent)
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
		it.Event = new(OracleL2ToL1MessageSent)
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
func (it *OracleL2ToL1MessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleL2ToL1MessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleL2ToL1MessageSent represents a L2ToL1MessageSent event raised by the Oracle contract.
type OracleL2ToL1MessageSent struct {
	MsgType uint8
	Payload []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterL2ToL1MessageSent is a free log retrieval operation binding the contract event 0xecdbbfd291a96bc9da716a82a98fa25ddf2c3ce6d2c1d945a8b7861581d585fb.
//
// Solidity: event L2ToL1MessageSent(uint8 msgType, bytes payload)
func (_Oracle *OracleFilterer) FilterL2ToL1MessageSent(opts *bind.FilterOpts) (*OracleL2ToL1MessageSentIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "L2ToL1MessageSent")
	if err != nil {
		return nil, err
	}
	return &OracleL2ToL1MessageSentIterator{contract: _Oracle.contract, event: "L2ToL1MessageSent", logs: logs, sub: sub}, nil
}

// WatchL2ToL1MessageSent is a free log subscription operation binding the contract event 0xecdbbfd291a96bc9da716a82a98fa25ddf2c3ce6d2c1d945a8b7861581d585fb.
//
// Solidity: event L2ToL1MessageSent(uint8 msgType, bytes payload)
func (_Oracle *OracleFilterer) WatchL2ToL1MessageSent(opts *bind.WatchOpts, sink chan<- *OracleL2ToL1MessageSent) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "L2ToL1MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleL2ToL1MessageSent)
				if err := _Oracle.contract.UnpackLog(event, "L2ToL1MessageSent", log); err != nil {
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

// ParseL2ToL1MessageSent is a log parse operation binding the contract event 0xecdbbfd291a96bc9da716a82a98fa25ddf2c3ce6d2c1d945a8b7861581d585fb.
//
// Solidity: event L2ToL1MessageSent(uint8 msgType, bytes payload)
func (_Oracle *OracleFilterer) ParseL2ToL1MessageSent(log types.Log) (*OracleL2ToL1MessageSent, error) {
	event := new(OracleL2ToL1MessageSent)
	if err := _Oracle.contract.UnpackLog(event, "L2ToL1MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleNewAggregatorIterator is returned from FilterNewAggregator and is used to iterate over the raw logs and unpacked data for NewAggregator events raised by the Oracle contract.
type OracleNewAggregatorIterator struct {
	Event *OracleNewAggregator // Event containing the contract specifics and raw log

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
func (it *OracleNewAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleNewAggregator)
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
		it.Event = new(OracleNewAggregator)
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
func (it *OracleNewAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleNewAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleNewAggregator represents a NewAggregator event raised by the Oracle contract.
type OracleNewAggregator struct {
	ValidatorID *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewAggregator is a free log retrieval operation binding the contract event 0xd4f48376eb9019821edbe6122e0dfd423575a8d95aeb6bc42995757f4cfc65bc.
//
// Solidity: event NewAggregator(uint256 validatorID)
func (_Oracle *OracleFilterer) FilterNewAggregator(opts *bind.FilterOpts) (*OracleNewAggregatorIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewAggregator")
	if err != nil {
		return nil, err
	}
	return &OracleNewAggregatorIterator{contract: _Oracle.contract, event: "NewAggregator", logs: logs, sub: sub}, nil
}

// WatchNewAggregator is a free log subscription operation binding the contract event 0xd4f48376eb9019821edbe6122e0dfd423575a8d95aeb6bc42995757f4cfc65bc.
//
// Solidity: event NewAggregator(uint256 validatorID)
func (_Oracle *OracleFilterer) WatchNewAggregator(opts *bind.WatchOpts, sink chan<- *OracleNewAggregator) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "NewAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleNewAggregator)
				if err := _Oracle.contract.UnpackLog(event, "NewAggregator", log); err != nil {
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

// ParseNewAggregator is a log parse operation binding the contract event 0xd4f48376eb9019821edbe6122e0dfd423575a8d95aeb6bc42995757f4cfc65bc.
//
// Solidity: event NewAggregator(uint256 validatorID)
func (_Oracle *OracleFilterer) ParseNewAggregator(log types.Log) (*OracleNewAggregator, error) {
	event := new(OracleNewAggregator)
	if err := _Oracle.contract.UnpackLog(event, "NewAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the Oracle contract.
type OracleRegisteredIterator struct {
	Event *OracleRegistered // Event containing the contract specifics and raw log

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
func (it *OracleRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleRegistered)
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
		it.Event = new(OracleRegistered)
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
func (it *OracleRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleRegistered represents a Registered event raised by the Oracle contract.
type OracleRegistered struct {
	Sender common.Address
	Index  *big.Int
	Pubkey OraclePublicKey
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x195cbf3e686960cbbd84ae23e6bd39c44cf97ba7b9f1fbe6b06a84b2a7d7757c.
//
// Solidity: event Registered(address sender, uint256 index, (uint256,uint256) pubkey, uint256 value)
func (_Oracle *OracleFilterer) FilterRegistered(opts *bind.FilterOpts) (*OracleRegisteredIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &OracleRegisteredIterator{contract: _Oracle.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x195cbf3e686960cbbd84ae23e6bd39c44cf97ba7b9f1fbe6b06a84b2a7d7757c.
//
// Solidity: event Registered(address sender, uint256 index, (uint256,uint256) pubkey, uint256 value)
func (_Oracle *OracleFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *OracleRegistered) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleRegistered)
				if err := _Oracle.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x195cbf3e686960cbbd84ae23e6bd39c44cf97ba7b9f1fbe6b06a84b2a7d7757c.
//
// Solidity: event Registered(address sender, uint256 index, (uint256,uint256) pubkey, uint256 value)
func (_Oracle *OracleFilterer) ParseRegistered(log types.Log) (*OracleRegistered, error) {
	event := new(OracleRegistered)
	if err := _Oracle.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleReplacedIterator is returned from FilterReplaced and is used to iterate over the raw logs and unpacked data for Replaced events raised by the Oracle contract.
type OracleReplacedIterator struct {
	Event *OracleReplaced // Event containing the contract specifics and raw log

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
func (it *OracleReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleReplaced)
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
		it.Event = new(OracleReplaced)
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
func (it *OracleReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleReplaced represents a Replaced event raised by the Oracle contract.
type OracleReplaced struct {
	Sender   common.Address
	Replaced common.Address
	Index    *big.Int
	Pubkey   OraclePublicKey
	Stake    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReplaced is a free log retrieval operation binding the contract event 0x8be3b32745cd34cae7fce370fa47fea5adb649d8805f55d6ed039cfa1b063345.
//
// Solidity: event Replaced(address indexed sender, address indexed replaced, uint256 index, (uint256,uint256) pubkey, uint256 stake)
func (_Oracle *OracleFilterer) FilterReplaced(opts *bind.FilterOpts, sender []common.Address, replaced []common.Address) (*OracleReplacedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var replacedRule []interface{}
	for _, replacedItem := range replaced {
		replacedRule = append(replacedRule, replacedItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Replaced", senderRule, replacedRule)
	if err != nil {
		return nil, err
	}
	return &OracleReplacedIterator{contract: _Oracle.contract, event: "Replaced", logs: logs, sub: sub}, nil
}

// WatchReplaced is a free log subscription operation binding the contract event 0x8be3b32745cd34cae7fce370fa47fea5adb649d8805f55d6ed039cfa1b063345.
//
// Solidity: event Replaced(address indexed sender, address indexed replaced, uint256 index, (uint256,uint256) pubkey, uint256 stake)
func (_Oracle *OracleFilterer) WatchReplaced(opts *bind.WatchOpts, sink chan<- *OracleReplaced, sender []common.Address, replaced []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var replacedRule []interface{}
	for _, replacedItem := range replaced {
		replacedRule = append(replacedRule, replacedItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Replaced", senderRule, replacedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleReplaced)
				if err := _Oracle.contract.UnpackLog(event, "Replaced", log); err != nil {
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

// ParseReplaced is a log parse operation binding the contract event 0x8be3b32745cd34cae7fce370fa47fea5adb649d8805f55d6ed039cfa1b063345.
//
// Solidity: event Replaced(address indexed sender, address indexed replaced, uint256 index, (uint256,uint256) pubkey, uint256 stake)
func (_Oracle *OracleFilterer) ParseReplaced(log types.Log) (*OracleReplaced, error) {
	event := new(OracleReplaced)
	if err := _Oracle.contract.UnpackLog(event, "Replaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleReplacementFromL1ProcessedIterator is returned from FilterReplacementFromL1Processed and is used to iterate over the raw logs and unpacked data for ReplacementFromL1Processed events raised by the Oracle contract.
type OracleReplacementFromL1ProcessedIterator struct {
	Event *OracleReplacementFromL1Processed // Event containing the contract specifics and raw log

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
func (it *OracleReplacementFromL1ProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleReplacementFromL1Processed)
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
		it.Event = new(OracleReplacementFromL1Processed)
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
func (it *OracleReplacementFromL1ProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleReplacementFromL1ProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleReplacementFromL1Processed represents a ReplacementFromL1Processed event raised by the Oracle contract.
type OracleReplacementFromL1Processed struct {
	RequestId *big.Int
	Success   bool
	NewRoot   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReplacementFromL1Processed is a free log retrieval operation binding the contract event 0x1c30d040fa77d7d013fc7ce2f6604b8fb3ebc94d7b070a101468df13e4393d83.
//
// Solidity: event ReplacementFromL1Processed(uint256 indexed requestId, bool success, uint256 newRoot)
func (_Oracle *OracleFilterer) FilterReplacementFromL1Processed(opts *bind.FilterOpts, requestId []*big.Int) (*OracleReplacementFromL1ProcessedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "ReplacementFromL1Processed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleReplacementFromL1ProcessedIterator{contract: _Oracle.contract, event: "ReplacementFromL1Processed", logs: logs, sub: sub}, nil
}

// WatchReplacementFromL1Processed is a free log subscription operation binding the contract event 0x1c30d040fa77d7d013fc7ce2f6604b8fb3ebc94d7b070a101468df13e4393d83.
//
// Solidity: event ReplacementFromL1Processed(uint256 indexed requestId, bool success, uint256 newRoot)
func (_Oracle *OracleFilterer) WatchReplacementFromL1Processed(opts *bind.WatchOpts, sink chan<- *OracleReplacementFromL1Processed, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "ReplacementFromL1Processed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleReplacementFromL1Processed)
				if err := _Oracle.contract.UnpackLog(event, "ReplacementFromL1Processed", log); err != nil {
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

// ParseReplacementFromL1Processed is a log parse operation binding the contract event 0x1c30d040fa77d7d013fc7ce2f6604b8fb3ebc94d7b070a101468df13e4393d83.
//
// Solidity: event ReplacementFromL1Processed(uint256 indexed requestId, bool success, uint256 newRoot)
func (_Oracle *OracleFilterer) ParseReplacementFromL1Processed(log types.Log) (*OracleReplacementFromL1Processed, error) {
	event := new(OracleReplacementFromL1Processed)
	if err := _Oracle.contract.UnpackLog(event, "ReplacementFromL1Processed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleUserRegisteredIterator is returned from FilterUserRegistered and is used to iterate over the raw logs and unpacked data for UserRegistered events raised by the Oracle contract.
type OracleUserRegisteredIterator struct {
	Event *OracleUserRegistered // Event containing the contract specifics and raw log

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
func (it *OracleUserRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleUserRegistered)
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
		it.Event = new(OracleUserRegistered)
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
func (it *OracleUserRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleUserRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleUserRegistered represents a UserRegistered event raised by the Oracle contract.
type OracleUserRegistered struct {
	Addr    common.Address
	Index   *big.Int
	Pubkey  OraclePublicKey
	Balance *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUserRegistered is a free log retrieval operation binding the contract event 0x0140a42f83a6b549a1f9aeb320dd6aa0d41c340639250b5f02a15a7ea7055302.
//
// Solidity: event UserRegistered(address addr, uint256 index, (uint256,uint256) pubkey, uint256 balance)
func (_Oracle *OracleFilterer) FilterUserRegistered(opts *bind.FilterOpts) (*OracleUserRegisteredIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "UserRegistered")
	if err != nil {
		return nil, err
	}
	return &OracleUserRegisteredIterator{contract: _Oracle.contract, event: "UserRegistered", logs: logs, sub: sub}, nil
}

// WatchUserRegistered is a free log subscription operation binding the contract event 0x0140a42f83a6b549a1f9aeb320dd6aa0d41c340639250b5f02a15a7ea7055302.
//
// Solidity: event UserRegistered(address addr, uint256 index, (uint256,uint256) pubkey, uint256 balance)
func (_Oracle *OracleFilterer) WatchUserRegistered(opts *bind.WatchOpts, sink chan<- *OracleUserRegistered) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "UserRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleUserRegistered)
				if err := _Oracle.contract.UnpackLog(event, "UserRegistered", log); err != nil {
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

// ParseUserRegistered is a log parse operation binding the contract event 0x0140a42f83a6b549a1f9aeb320dd6aa0d41c340639250b5f02a15a7ea7055302.
//
// Solidity: event UserRegistered(address addr, uint256 index, (uint256,uint256) pubkey, uint256 balance)
func (_Oracle *OracleFilterer) ParseUserRegistered(log types.Log) (*OracleUserRegistered, error) {
	event := new(OracleUserRegistered)
	if err := _Oracle.contract.UnpackLog(event, "UserRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleValidatorRegisteredIterator is returned from FilterValidatorRegistered and is used to iterate over the raw logs and unpacked data for ValidatorRegistered events raised by the Oracle contract.
type OracleValidatorRegisteredIterator struct {
	Event *OracleValidatorRegistered // Event containing the contract specifics and raw log

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
func (it *OracleValidatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleValidatorRegistered)
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
		it.Event = new(OracleValidatorRegistered)
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
func (it *OracleValidatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleValidatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleValidatorRegistered represents a ValidatorRegistered event raised by the Oracle contract.
type OracleValidatorRegistered struct {
	Addr        common.Address
	ValidatorID *big.Int
	Index       *big.Int
	Pubkey      OraclePublicKey
	Balance     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterValidatorRegistered is a free log retrieval operation binding the contract event 0xab98958ae468add4aa956316226e8db683a814ddc402eede01b06d5786d43011.
//
// Solidity: event ValidatorRegistered(address addr, uint256 validatorID, uint256 index, (uint256,uint256) pubkey, uint256 balance)
func (_Oracle *OracleFilterer) FilterValidatorRegistered(opts *bind.FilterOpts) (*OracleValidatorRegisteredIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "ValidatorRegistered")
	if err != nil {
		return nil, err
	}
	return &OracleValidatorRegisteredIterator{contract: _Oracle.contract, event: "ValidatorRegistered", logs: logs, sub: sub}, nil
}

// WatchValidatorRegistered is a free log subscription operation binding the contract event 0xab98958ae468add4aa956316226e8db683a814ddc402eede01b06d5786d43011.
//
// Solidity: event ValidatorRegistered(address addr, uint256 validatorID, uint256 index, (uint256,uint256) pubkey, uint256 balance)
func (_Oracle *OracleFilterer) WatchValidatorRegistered(opts *bind.WatchOpts, sink chan<- *OracleValidatorRegistered) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "ValidatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleValidatorRegistered)
				if err := _Oracle.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
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

// ParseValidatorRegistered is a log parse operation binding the contract event 0xab98958ae468add4aa956316226e8db683a814ddc402eede01b06d5786d43011.
//
// Solidity: event ValidatorRegistered(address addr, uint256 validatorID, uint256 index, (uint256,uint256) pubkey, uint256 balance)
func (_Oracle *OracleFilterer) ParseValidatorRegistered(log types.Log) (*OracleValidatorRegistered, error) {
	event := new(OracleValidatorRegistered)
	if err := _Oracle.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleValidatorsImportedFromL1Iterator is returned from FilterValidatorsImportedFromL1 and is used to iterate over the raw logs and unpacked data for ValidatorsImportedFromL1 events raised by the Oracle contract.
type OracleValidatorsImportedFromL1Iterator struct {
	Event *OracleValidatorsImportedFromL1 // Event containing the contract specifics and raw log

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
func (it *OracleValidatorsImportedFromL1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleValidatorsImportedFromL1)
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
		it.Event = new(OracleValidatorsImportedFromL1)
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
func (it *OracleValidatorsImportedFromL1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleValidatorsImportedFromL1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleValidatorsImportedFromL1 represents a ValidatorsImportedFromL1 event raised by the Oracle contract.
type OracleValidatorsImportedFromL1 struct {
	Count   *big.Int
	NewRoot *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterValidatorsImportedFromL1 is a free log retrieval operation binding the contract event 0x1179e923b391693c40c7210e742a45dc34fce4496618b8d9f02fd90249e3a9c7.
//
// Solidity: event ValidatorsImportedFromL1(uint256 count, uint256 newRoot)
func (_Oracle *OracleFilterer) FilterValidatorsImportedFromL1(opts *bind.FilterOpts) (*OracleValidatorsImportedFromL1Iterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "ValidatorsImportedFromL1")
	if err != nil {
		return nil, err
	}
	return &OracleValidatorsImportedFromL1Iterator{contract: _Oracle.contract, event: "ValidatorsImportedFromL1", logs: logs, sub: sub}, nil
}

// WatchValidatorsImportedFromL1 is a free log subscription operation binding the contract event 0x1179e923b391693c40c7210e742a45dc34fce4496618b8d9f02fd90249e3a9c7.
//
// Solidity: event ValidatorsImportedFromL1(uint256 count, uint256 newRoot)
func (_Oracle *OracleFilterer) WatchValidatorsImportedFromL1(opts *bind.WatchOpts, sink chan<- *OracleValidatorsImportedFromL1) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "ValidatorsImportedFromL1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleValidatorsImportedFromL1)
				if err := _Oracle.contract.UnpackLog(event, "ValidatorsImportedFromL1", log); err != nil {
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

// ParseValidatorsImportedFromL1 is a log parse operation binding the contract event 0x1179e923b391693c40c7210e742a45dc34fce4496618b8d9f02fd90249e3a9c7.
//
// Solidity: event ValidatorsImportedFromL1(uint256 count, uint256 newRoot)
func (_Oracle *OracleFilterer) ParseValidatorsImportedFromL1(log types.Log) (*OracleValidatorsImportedFromL1, error) {
	event := new(OracleValidatorsImportedFromL1)
	if err := _Oracle.contract.UnpackLog(event, "ValidatorsImportedFromL1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleWiVoteSubmittedIterator is returned from FilterWiVoteSubmitted and is used to iterate over the raw logs and unpacked data for WiVoteSubmitted events raised by the Oracle contract.
type OracleWiVoteSubmittedIterator struct {
	Event *OracleWiVoteSubmitted // Event containing the contract specifics and raw log

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
func (it *OracleWiVoteSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleWiVoteSubmitted)
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
		it.Event = new(OracleWiVoteSubmitted)
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
func (it *OracleWiVoteSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleWiVoteSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleWiVoteSubmitted represents a WiVoteSubmitted event raised by the Oracle contract.
type OracleWiVoteSubmitted struct {
	Submitter    *big.Int
	Validators   *big.Int
	HonestBits   *big.Int
	Request      *big.Int
	MajorityVote *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWiVoteSubmitted is a free log retrieval operation binding the contract event 0x73f336750f8104b46b6917b76a289163d6ceb486db25df9b8b71c5316ee064dd.
//
// Solidity: event WiVoteSubmitted(uint256 submitter, uint256 validators, uint256 honestBits, uint256 request, uint256 majorityVote)
func (_Oracle *OracleFilterer) FilterWiVoteSubmitted(opts *bind.FilterOpts) (*OracleWiVoteSubmittedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "WiVoteSubmitted")
	if err != nil {
		return nil, err
	}
	return &OracleWiVoteSubmittedIterator{contract: _Oracle.contract, event: "WiVoteSubmitted", logs: logs, sub: sub}, nil
}

// WatchWiVoteSubmitted is a free log subscription operation binding the contract event 0x73f336750f8104b46b6917b76a289163d6ceb486db25df9b8b71c5316ee064dd.
//
// Solidity: event WiVoteSubmitted(uint256 submitter, uint256 validators, uint256 honestBits, uint256 request, uint256 majorityVote)
func (_Oracle *OracleFilterer) WatchWiVoteSubmitted(opts *bind.WatchOpts, sink chan<- *OracleWiVoteSubmitted) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "WiVoteSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleWiVoteSubmitted)
				if err := _Oracle.contract.UnpackLog(event, "WiVoteSubmitted", log); err != nil {
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

// ParseWiVoteSubmitted is a log parse operation binding the contract event 0x73f336750f8104b46b6917b76a289163d6ceb486db25df9b8b71c5316ee064dd.
//
// Solidity: event WiVoteSubmitted(uint256 submitter, uint256 validators, uint256 honestBits, uint256 request, uint256 majorityVote)
func (_Oracle *OracleFilterer) ParseWiVoteSubmitted(log types.Log) (*OracleWiVoteSubmitted, error) {
	event := new(OracleWiVoteSubmitted)
	if err := _Oracle.contract.UnpackLog(event, "WiVoteSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the Oracle contract.
type OracleWithdrawnIterator struct {
	Event *OracleWithdrawn // Event containing the contract specifics and raw log

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
func (it *OracleWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleWithdrawn)
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
		it.Event = new(OracleWithdrawn)
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
func (it *OracleWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleWithdrawn represents a Withdrawn event raised by the Oracle contract.
type OracleWithdrawn struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xf45a04d08a70caa7eb4b747571305559ad9fdf4a093afd41506b35c8a306fa94.
//
// Solidity: event Withdrawn(address indexed sender)
func (_Oracle *OracleFilterer) FilterWithdrawn(opts *bind.FilterOpts, sender []common.Address) (*OracleWithdrawnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Withdrawn", senderRule)
	if err != nil {
		return nil, err
	}
	return &OracleWithdrawnIterator{contract: _Oracle.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xf45a04d08a70caa7eb4b747571305559ad9fdf4a093afd41506b35c8a306fa94.
//
// Solidity: event Withdrawn(address indexed sender)
func (_Oracle *OracleFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *OracleWithdrawn, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Withdrawn", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleWithdrawn)
				if err := _Oracle.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xf45a04d08a70caa7eb4b747571305559ad9fdf4a093afd41506b35c8a306fa94.
//
// Solidity: event Withdrawn(address indexed sender)
func (_Oracle *OracleFilterer) ParseWithdrawn(log types.Log) (*OracleWithdrawn, error) {
	event := new(OracleWithdrawn)
	if err := _Oracle.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
