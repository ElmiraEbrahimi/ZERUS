package oracle_runtime

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"

	merkleproof "l2alchemy/circuits/merkle_proof"
	votingbatch "l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/config"
	"l2alchemy/internal/oracle-repo"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/user"
	"l2alchemy/internal/oracle-repo/util"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/ethclient"
)

// OracleEngine is a config-like runtime container that owns the oracle-related state machine
// and data (including users). It is intentionally decoupled from the HTTP server handler layer.
type OracleEngine struct {
	Users map[string]*user.User

	Cfg       *config.Config
	EthClient *ethclient.Client
	Ipfs      *db.IPFSClient

	MerkleCircuit *merkleproof.MerkleProofCircuit
	MerkleR1cs    constraint.ConstraintSystem
	MerklePK      groth16.ProvingKey
	MerkleVK      groth16.VerifyingKey

	VotingCircuit *votingbatch.BatchingVotingCircuit
	VotingR1cs    constraint.ConstraintSystem
	VotingPK      groth16.ProvingKey
	VotingVK      groth16.VerifyingKey

	Oracle *oracle.Oracle

	mu sync.RWMutex
}

var (
	globalMu sync.RWMutex
	global   *OracleEngine
)

// Global returns the initialized OracleEngine instance (if any).
func Global() *OracleEngine {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return global
}

// Init constructs the OracleEngine and performs bootstrap actions that must happen before the
// server starts (including creating an initial user via user.NewUser).
func Init(cfg *config.Config, keyDir string) (*OracleEngine, error) {
	log.Println("oracle runtime initializing...")

	if cfg == nil {
		return nil, fmt.Errorf("oracle runtime: cfg is nil")
	}
	if keyDir == "" {
		return nil, fmt.Errorf("oracle runtime: keyDir is empty")
	}

	// Avoid double-init in case main is invoked twice in tests.
	globalMu.Lock()
	defer globalMu.Unlock()
	if global != nil {
		return global, nil
	}

	// setup and compile circuits:
	log.Println("oracle runtime: setting up circuits...")
	merkleCircuit, merkleR1cs, merklePK, merkleVK := setupMerkleCircuit(cfg, keyDir)
	votingCircuit, votingR1cs, votingPK, votingVK := setupVotingCircuit(cfg, keyDir)
	log.Println("oracle runtime: completed setup circuits")

	ctx := context.Background()
	ethCl, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: dial eth client: %w", err)
	}

	// create ipfs:
	ipfsCl, err := db.NewIPFSClient(true, cfg.IncTreeDepth, cfg.IPFSSimDataPath)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: create ipfs client: %w", err)
	}
	log.Println("oracle runtime: created ipfs client")

	// create oracle and validators:
	log.Println("oracle runtime init: generating validator nodes...")
	var privateKeys []*eddsa.PrivateKey
	privateKeys, err = util.GenerateKeys(cfg.NodeCount)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: generate eddsa keys for validators: %w", err)
	}

	validatorAccounts, err := gnark.CreateAccounts(privateKeys)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: create accounts for validators: %w", err)
	}

	log.Println("oracle runtime init: setting up oracle...")
	oracle := oracle.NewOracle(cfg, ethCl, ipfsCl, cfg.NodeCount, merkleCircuit, merkleR1cs, merkleVK, votingCircuit, votingR1cs, votingPK, votingVK, privateKeys, validatorAccounts)

	// create the default user:
	userPK, err := eddsa.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: generate eddsa key for default user: %w", err)
	}
	log.Println("oracle runtime: successfully generated private key for default user")

	userAcct, err := gnark.CreateAccount(userPK, 0)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: create account for default user: %w", err)
	}
	log.Println("oracle runtime: successfully created account for default user")

	usr := user.NewUser(
		cfg,
		ethCl,
		ipfsCl,
		merkleCircuit,
		merkleR1cs,
		merklePK,
		merkleVK,
		"default-user",
		userPK,
		userAcct,
	)
	usrPtr := &usr

	engine := &OracleEngine{
		Users:         map[string]*user.User{usrPtr.Name: usrPtr},
		Cfg:           cfg,
		EthClient:     ethCl,
		Ipfs:          ipfsCl,
		MerkleCircuit: merkleCircuit,
		MerkleR1cs:    merkleR1cs,
		MerklePK:      merklePK,
		MerkleVK:      merkleVK,
		VotingCircuit: votingCircuit,
		VotingR1cs:    votingR1cs,
		VotingPK:      votingPK,
		VotingVK:      votingVK,
		Oracle:        oracle,
	}

	global = engine
	return engine, nil
}
