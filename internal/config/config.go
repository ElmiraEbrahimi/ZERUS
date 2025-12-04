package config

import (
    "os"
    "strconv"
)

// DBConfig holds database connection settings.
// Values are sourced from environment variables and sensible defaults.
type DBConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
    SSLMode  string
}

// EthConfig contains configuration for the Ethereum L1 client.
type EthConfig struct {
    RPCURL  string
    ChainID int64
    Timeout int // seconds
}

// ZkSyncConfig contains configuration for the zkSync L2 client.
type ZkSyncConfig struct {
    RPCURL  string
    ChainID int64
    Timeout int // seconds
}

// AppConfig encapsulates all application configuration sections.
type AppConfig struct {
    Env      string
    HTTPPort string
    DB       DBConfig
    Eth      EthConfig
    ZkSync   ZkSyncConfig
}

// Load reads environment variables to construct an AppConfig.
// Sensible defaults are applied when variables are not present, making
// local development setup straightforward.
func Load() (*AppConfig, error) {
    cfg := &AppConfig{}
    cfg.Env = getEnv("APP_ENV", "local")
    cfg.HTTPPort = getEnv("APP_PORT", "8080")
    // Database configuration
    cfg.DB.Host = getEnv("DB_HOST", "db")
    cfg.DB.Port = getEnvAsInt("DB_PORT", 5432)
    cfg.DB.User = getEnv("DB_USER", "postgres")
    cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
    cfg.DB.Name = getEnv("DB_NAME", "postgres")
    cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")
    // Ethereum L1 configuration
    cfg.Eth.RPCURL = getEnv("ETH_RPC_URL", "http://eth:8545")
    cfg.Eth.ChainID = getEnvAsInt64("ETH_CHAIN_ID", 1)
    cfg.Eth.Timeout = getEnvAsInt("ETH_TIMEOUT", 10)
    // zkSync L2 configuration
    cfg.ZkSync.RPCURL = getEnv("ZKSYNC_RPC_URL", "http://zksync:3050")
    cfg.ZkSync.ChainID = getEnvAsInt64("ZKSYNC_CHAIN_ID", 280)
    cfg.ZkSync.Timeout = getEnvAsInt("ZKSYNC_TIMEOUT", 10)
    return cfg, nil
}

// getEnv returns the value of the given environment variable or the default
// value if the variable is not set.
func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}

// getEnvAsInt parses the environment variable into an int. If parsing fails
// or the variable is not set, the default value is returned.
func getEnvAsInt(key string, defaultVal int) int {
    if valStr := os.Getenv(key); valStr != "" {
        if val, err := strconv.Atoi(valStr); err == nil {
            return val
        }
    }
    return defaultVal
}

// getEnvAsInt64 parses the environment variable into an int64. If parsing
// fails or the variable is not set, the default value is returned.
func getEnvAsInt64(key string, defaultVal int64) int64 {
    if valStr := os.Getenv(key); valStr != "" {
        if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
            return val
        }
    }
    return defaultVal
}