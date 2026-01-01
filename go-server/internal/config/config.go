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

	BlockchainURL string `env:"ZKSYNC_RPC_URL,required"`

	OracleSeedX string `env:"ORACLE_SEED_X,required"`
	OracleSeedY string `env:"ORACLE_SEED_Y,required"`
	UserPK      string `env:"USER_PK,required"`
	NodePK      string `env:"NODE_PK,required"`

	ChainID                 int64 `env:"ZKSYNC_CHAIN_ID,required"`
	RoundDurationMS         int   `env:"ROUND_DURATION_MILI_SECONDS,required"`
	NodeCount               int   `env:"NODE_COUNT,required"`
	SparseTreeDepth         int   `env:"SPARSE_TREE_DEPTH,required"`
	BatchSize               int   `env:"BATCH_SIZE,required"`
	RoundToChangeAggregator int   `env:"ROUND_TO_CHANGE_AGGREGATOR,required"`
	RoundLimit              int   `env:"ROUND_LIMIT,required"`
	RoundIncVoteSelection   int   `env:"ROUND_INCVOTE_SELECTION,required"`
	RoundWiVoteSelection    int   `env:"ROUND_WIVOTE_SELECTION,required"`
	IncTreeDepth            int   `env:"INC_TREE_DEPTH,required"`
	TxGasLimit              int64 `env:"TX_GAS_LIMIT" default:"3000000"`
	TxGasPriceWei           int64 `env:"TX_GAS_PRICE_WEI" default:"20000000000"`
	TxGasFeeCapWei          int64 `env:"TX_GAS_FEE_CAP_WEI" default:"0"`
	TxGasTipCapWei          int64 `env:"TX_GAS_TIP_CAP_WEI" default:"0"`

	HTTPBindAddr    string `env:"HTTP_BIND_ADDR" default:":8080"`
	IPFSSimDataPath string `env:"IPFS_SIM_DATA_PATH" default:".ipfs_sim_data.gob"`
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
	minDepth := 0
	for (1 << minDepth) < cfg.NodeCount {
		minDepth++
	}
	if cfg.SparseTreeDepth < minDepth {
		return fmt.Errorf("SPARSE_TREE_DEPTH=%d too small for NODE_COUNT=%d (need at least %d)", cfg.SparseTreeDepth, cfg.NodeCount, minDepth)
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
