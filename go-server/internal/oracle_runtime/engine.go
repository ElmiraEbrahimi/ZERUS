package oracle_runtime

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"

	"l2alchemy/internal/config"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/user"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/ethclient"
)

// OracleEngine is a config-like runtime container that owns the oracle-related state machine
// and data (including users). It is intentionally decoupled from the HTTP server handler layer.
type OracleEngine struct {
	Users map[string]*user.User

	cfg       *config.Config
	wg        *sync.WaitGroup
	ethClient *ethclient.Client
	ipfs      *db.IPFSClient

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
func Init(cfg *config.Config) (*OracleEngine, error) {
	if cfg == nil {
		return nil, fmt.Errorf("oracle runtime init: cfg is nil")
	}

	// Avoid double-init in case main is invoked twice in tests.
	globalMu.Lock()
	defer globalMu.Unlock()
	if global != nil {
		return global, nil
	}

	ctx := context.Background()
	ethCl, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: dial eth client: %w", err)
	}

	// create ipfs:
	ipfsCl, err := db.NewIPFSClient(true, cfg.IncTreeDepth)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: create ipfs client: %w", err)
	}
	log.Println("oracle runtime: created ipfs client")

	wg := &sync.WaitGroup{}

	// create the default user:
	userPK, err := eddsa.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: generate eddsa key for default user: %w", err)
	}
	log.Println("oracle runtime: generated private key for default user")

	userAcct, err := gnark.CreateAccount(userPK, 0)
	if err != nil {
		return nil, fmt.Errorf("oracle runtime init: create account for default user: %w", err)
	}
	log.Println("oracle runtime: created account for default user")

	usr := user.NewUser(
		wg,
		cfg,
		ethCl,
		ipfsCl,
		nil, // TODO
		nil,
		nil,
		nil,
		"default-user",
		userPK,
		userAcct,
	)
	usrPtr := &usr

	engine := &OracleEngine{
		Users:     map[string]*user.User{usrPtr.Name: usrPtr},
		cfg:       cfg,
		wg:        wg,
		ethClient: ethCl,
		ipfs:      ipfsCl,
	}

	global = engine
	return engine, nil
}
