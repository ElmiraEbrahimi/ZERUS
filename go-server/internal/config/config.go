package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config is loaded from environment variables (optionally from a .env file).
// Bindings are declared via struct tags to avoid repetitive os.Getenv calls.
//
// Tag format:
//
//	env:"ENV_NAME[,required]"   - environment variable name, optionally required
//	default:"value"             - default value if env var is unset/empty
//
// Types supported: string, bool, int, int64.
type Config struct {
	// Strings
	RPCURL     string `env:"ZKSYNC_RPC_URL,required"`
	PrivateKey string `env:"ZKSYNC_PRIVATE_KEY,required"`

	CounterContractAddress     string `env:"COUNTER_CONTRACT_ADDRESS,required"`
	MerkleVerifierAddress      string `env:"MERKLE_VERIFIER_ADDRESS,required"`
	VotingBatchVerifierAddress string `env:"VOTING_BATCH_VERIFIER_ADDRESS,required"`
	MerkleTreeContractAddress  string `env:"MERKLE_TREE_CONTRACT_ADDRESS,required"`
	OracleContractAddress      string `env:"ORACLE_CONTRACT_ADDRESS,required"`

	L1RPCURL                     string `env:"L1_RPC_URL"`
	L1GasPriceWei                int64  `env:"L1_GAS_PRICE_WEI" default:"0"`
	L1L2ValueWei                 string `env:"L1_L2_VALUE_WEI"`
	L1MessengerContractAddress   string `env:"L1_MESSENGER_CONTRACT_ADDRESS"`
	L2MessengerContractAddress   string `env:"L2_MESSENGER_CONTRACT_ADDRESS"`
	L1MailboxAddress             string `env:"L1_MAILBOX_ADDRESS"`
	L1UseDirectMessaging         bool   `env:"L1_USE_DIRECT_MESSAGING" default:"false"`
	// L1HubContractAddress enables L1 anchoring (paper SIV-B/SIV-C): when
	// set together with L1_RPC_URL, the runtime relays the oracle's L2->L1
	// messages (checkpoints, exits, withdrawals, import/replacement
	// results) to the L1 Hub via finalizeFromL2 (F-01/F-02).
	L1HubContractAddress string `env:"L1_HUB_CONTRACT_ADDRESS"`

	BlockchainURL string `env:"ZKSYNC_RPC_URL,required"`

	OracleSeedX string `env:"ORACLE_SEED_X,required"`
	OracleSeedY string `env:"ORACLE_SEED_Y,required"`
	UserPK      string `env:"USER_PK,required"`
	NodePK      string `env:"NODE_PK,required"`

	ChainID                 int64 `env:"ZKSYNC_CHAIN_ID,required"`
	L1ChainID               int64 `env:"L1_CHAIN_ID" default:"0"`
	NodeCount               int   `env:"NODE_COUNT,required"`
	SparseTreeDepth         int   `env:"SPARSE_TREE_DEPTH,required"`
	BatchSize               int   `env:"BATCH_SIZE,required"`
	IncTreeDepth            int   `env:"INC_TREE_DEPTH,required"`
	// DestinationID is the per-deployment destination-rollup identifier
	// d_dst (paper SIV-E; F-24): one value per Gateway instance, used in the
	// burn commitment C = H(n_rd || s_rd || d_dst). Decimal field element.
	DestinationID string `env:"DESTINATION_ID" default:"9636219578937187601590327046728695236698322465209974782280717458744997515735"`
	// BurnConfirmationDepth is the source-rollup finality rule used before
	// inserting burns into the commitment tree (paper SIV-E / SVII-C;
	// F-21), expressed as a block-depth confirmation for the local stack.
	BurnConfirmationDepth int `env:"BURN_CONFIRMATION_DEPTH" default:"0"`
	// InitialValidatorStake is each simulated validator's starting stake /
	// state-tree balance, configurable for the evaluation sweeps (F-28).
	InitialValidatorStake int `env:"INITIAL_VALIDATOR_STAKE" default:"1000"`
	// Source rollup profile (paper SV: two independent zkSync Era
	// instances; F-20). When set, burns are observed on the source chain's
	// Gateway while claims run against the local (destination) Gateway.
	SourceRPCURL                string `env:"SOURCE_RPC_URL"`
	SourceOracleContractAddress string `env:"SOURCE_ORACLE_CONTRACT_ADDRESS"`
	SourceChainID               int64  `env:"SOURCE_CHAIN_ID" default:"0"`
	TxGasLimit              int64 `env:"TX_GAS_LIMIT" default:"3000000"`
	TxGasPriceWei           int64 `env:"TX_GAS_PRICE_WEI" default:"20000000000"`
	TxGasFeeCapWei          int64 `env:"TX_GAS_FEE_CAP_WEI" default:"0"`
	TxGasTipCapWei          int64 `env:"TX_GAS_TIP_CAP_WEI" default:"0"`

	HTTPBindAddr    string `env:"HTTP_BIND_ADDR" default:":8080"`
	IPFSSimDataPath string `env:"IPFS_SIM_DATA_PATH" default:".ipfs_sim_data.gob"`
	// IPFSAPIURL, when set (e.g. http://127.0.0.1:5001), stores DFS content
	// through a real IPFS daemon's HTTP API (paper SV; F-26). When empty,
	// the local simulation backed by IPFS_SIM_DATA_PATH is used.
	IPFSAPIURL string `env:"IPFS_API_URL"`
	UserStatePath   string `env:"USER_STATE_PATH" default:".user_state.gob"`
}

// Load reads .env (if present) and environment variables, binds them to Config,
// and validates required fields and type conversions.
func Load(printVals bool) (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../.env"); err2 != nil {
			log.Printf("no .env file found or unable to load it: %v", err2)
		}
	}
	cfg := &Config{}
	if err := bindEnv(cfg); err != nil {
		return nil, err
	}
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	if printVals {
		printEnvValues(cfg)
	}

	return cfg, nil
}

func bindEnv(out any) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("bindEnv expects pointer to struct, got %T", out)
	}

	v := rv.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		envTag := strings.TrimSpace(ft.Tag.Get("env"))
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envKey := strings.TrimSpace(parts[0])
		required := false
		for _, p := range parts[1:] {
			if strings.TrimSpace(p) == "required" {
				required = true
			}
		}

		val, ok := os.LookupEnv(envKey)
		if !ok || strings.TrimSpace(val) == "" {
			if def := ft.Tag.Get("default"); def != "" {
				val = def
			} else if required {
				return fmt.Errorf("missing required env var: %s", envKey)
			} else {
				continue
			}
		}

		if err := setValue(fv, envKey, val); err != nil {
			return err
		}
	}

	return nil
}

func setValue(field reflect.Value, key, raw string) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field for env var %s (field must be exported)", key)
	}

	raw = strings.TrimSpace(raw)

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid bool for %s: %q", key, raw)
		}
		field.SetBool(b)
		return nil

	case reflect.Int:
		n, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return fmt.Errorf("invalid int for %s: %q", key, raw)
		}
		field.SetInt(n)
		return nil

	case reflect.Int64:
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int64 for %s: %q", key, raw)
		}
		field.SetInt(n)
		return nil

	case reflect.Slice:
		// Only support []string for now
		if field.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice element type %s for env var %s", field.Type().Elem().Kind(), key)
		}

		// Parse comma-separated list
		// Example: LIST_ENV=val1,val2,val3
		parts := strings.Split(raw, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			s := strings.TrimSpace(p)
			if s == "" {
				continue
			}
			out = append(out, s)
		}

		field.Set(reflect.ValueOf(out))
		return nil

	default:
		return fmt.Errorf("unsupported field kind %s for env var %s", field.Kind(), key)
	}
}

func validateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if cfg.NodeCount < 1 {
		return fmt.Errorf("NODE_COUNT must be >= 1")
	}
	// The Aggregating circuit compiles Merkle paths of fixed length
	// SPARSE_TREE_DEPTH+1, and the off-chain gnark-crypto tree over
	// NODE_COUNT leaves only produces paths of that length (and roots that
	// match the fixed-depth on-chain tree) when NODE_COUNT is exactly
	// 2^SPARSE_TREE_DEPTH.
	if cfg.NodeCount != (1 << cfg.SparseTreeDepth) {
		return fmt.Errorf("NODE_COUNT=%d must equal 2^SPARSE_TREE_DEPTH (=%d)", cfg.NodeCount, 1<<cfg.SparseTreeDepth)
	}
	// The on-chain MerkleTree contract precomputes zero-subtree hashes only up
	// to depth 8 (see contracts/src/merkle_tree.sol MAX_LEVELS); it is deployed
	// with SPARSE_TREE_DEPTH levels, so reject configs it cannot support.
	const maxOnChainTreeDepth = 8
	if cfg.SparseTreeDepth > maxOnChainTreeDepth {
		return fmt.Errorf("SPARSE_TREE_DEPTH=%d exceeds the on-chain MerkleTree maximum of %d", cfg.SparseTreeDepth, maxOnChainTreeDepth)
	}
	return nil
}
func printEnvValues(cfg *Config) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	log.Println("=== Loaded environment configuration ===")

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		envTag := strings.TrimSpace(ft.Tag.Get("env"))
		if envTag == "" {
			continue
		}

		envKey := strings.Split(envTag, ",")[0]
		fieldType := ft.Type.String()

		log.Printf(
			"ENV %-35s | TYPE %-8s | VALUE %v",
			envKey,
			fieldType,
			fv.Interface(),
		)
	}

	log.Printf("=== End configuration dump ===\n\n")
}
