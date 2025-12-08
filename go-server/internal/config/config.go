package config

import (
    "fmt"
    "os"
    "strconv"
)

// Config holds runtime configuration for the counter API server. All fields
// correspond to environment variables that must be set before startup. A
// zero-value Config is invalid.
type Config struct {
    RPCURL          string
    ChainID         int64
    PrivateKey      string
    ContractAddress string
    HTTPBindAddr    string
}

// Load reads environment variables and constructs a Config. If any required
// variable is missing or malformed it returns an error.
func Load() (*Config, error) {
    cfg := &Config{
        RPCURL:          os.Getenv("ZKSYNC_RPC_URL"),
        PrivateKey:      os.Getenv("ZKSYNC_PRIVATE_KEY"),
        ContractAddress: os.Getenv("COUNTER_CONTRACT_ADDRESS"),
        HTTPBindAddr:    os.Getenv("HTTP_BIND_ADDR"),
    }
    if cfg.RPCURL == "" {
        return nil, fmt.Errorf("ZKSYNC_RPC_URL must be set")
    }
    if cfg.PrivateKey == "" {
        return nil, fmt.Errorf("ZKSYNC_PRIVATE_KEY must be set")
    }
    if cfg.ContractAddress == "" {
        return nil, fmt.Errorf("COUNTER_CONTRACT_ADDRESS must be set")
    }
    if cfg.HTTPBindAddr == "" {
        cfg.HTTPBindAddr = ":8080"
    }
    chainIDEnv := os.Getenv("ZKSYNC_CHAIN_ID")
    if chainIDEnv == "" {
        return nil, fmt.Errorf("ZKSYNC_CHAIN_ID must be set")
    }
    chainID, err := strconv.ParseInt(chainIDEnv, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("unable to parse ZKSYNC_CHAIN_ID (%s): %w", chainIDEnv, err)
    }
    cfg.ChainID = chainID
    return cfg, nil
}