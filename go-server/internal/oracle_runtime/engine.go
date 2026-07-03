package oracle_runtime

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	merkleproof "l2alchemy/circuits/merkle_proof"
	votingbatch "l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/config"
	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/oracle-repo"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/user"
	"l2alchemy/internal/oracle-repo/util"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

// validateOnChainBatchSize fails fast when the configured BATCH_SIZE does
// not match the deployed Gateway's immutable batch size: the deterministic
// claim windows (paper SIV-D) are computed from both, and a silent mismatch
// would misalign minting against the committee's batches.
//
// The read is retried for a grace window: right after `make deploy` the
// local zkSync node may briefly report "no contract code" for a freshly
// deployed contract (and the local stack is known to roll back recent
// miniblocks when it is in an inconsistent state), so a transient read
// failure must not kill the server instantly. A confirmed mismatch is
// fatal immediately.
func validateOnChainBatchSize(ctx context.Context, ethCl *ethclient.Client, cfg *config.Config) error {
	const (
		attempts = 10
		backoff  = 3 * time.Second
	)

	oracleAddr := common.HexToAddress(cfg.OracleContractAddress)
	oracleClient, err := bc.NewOracle(oracleAddr, ethCl)
	if err != nil {
		return fmt.Errorf("oracle runtime init: bind oracle contract: %w", err)
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		onChain, err := oracleClient.BatchSize(&bind.CallOpts{Context: ctx})
		if err == nil {
			if !onChain.IsInt64() || onChain.Int64() != int64(cfg.BatchSize) {
				return fmt.Errorf(
					"oracle runtime init: BATCH_SIZE=%d does not match the deployed Gateway's batch size %s; redeploy the Oracle (make deploy-oracle) or fix .env",
					cfg.BatchSize, onChain,
				)
			}
			return nil
		}
		lastErr = err
		log.Printf(
			"oracle runtime init: read on-chain batch size failed (attempt %d/%d, oracle=%s): %v",
			attempt, attempts, cfg.OracleContractAddress, err,
		)
		if attempt < attempts {
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf(
		"oracle runtime init: could not read the Gateway's batch size at %s after %d attempts: %w\n"+
			"  - if the address is stale, redeploy with 'make deploy'\n"+
			"  - if the local chain lost recently deployed state (rolled-back miniblocks), reset it with 'make down && make up-deploy'",
		cfg.OracleContractAddress, attempts, lastErr,
	)
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

	// The claim-identifier windows [r*b, (r+1)*b - 1] (paper SIV-D) are
	// enforced with the contract's immutable batch size; a BATCH_SIZE that
	// diverges from the deployed Gateway would silently misalign minting,
	// so fail fast on mismatch.
	if err := validateOnChainBatchSize(ctx, ethCl, cfg); err != nil {
		return nil, err
	}

	// create ipfs: real daemon when IPFS_API_URL is set (paper SV; F-26),
	// local simulation otherwise so tests run without a daemon.
	ipfsCl, err := db.NewIPFSClient(cfg.IPFSAPIURL == "", cfg.IncTreeDepth, cfg.IPFSSimDataPath, cfg.IPFSAPIURL)
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
