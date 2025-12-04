package main

import (
    "log"

    "github.com/example/bridge-service/internal/api"
    "github.com/example/bridge-service/internal/bridge"
    "github.com/example/bridge-service/internal/config"
    "github.com/example/bridge-service/internal/db"
    elog "github.com/example/bridge-service/internal/log"
    "github.com/example/bridge-service/internal/eth"
    "github.com/example/bridge-service/internal/zksync"
)

// main is the entrypoint for the bridge service. It loads configuration,
// initializes dependencies and starts the HTTP server. Any fatal errors
// during startup cause the program to exit with a non-zero status.
func main() {
    // Configure logger
    elog.Init()
    // Load configuration from environment
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }
    // Connect to the database
    dbConn, err := db.Connect(cfg.DB)
    if err != nil {
        log.Fatalf("database connection failed: %v", err)
    }
    defer dbConn.Close()
    // Initialize Ethereum L1 client
    ethClient, err := eth.NewRPCClient(cfg.Eth.RPCURL, cfg.Eth.Timeout)
    if err != nil {
        log.Fatalf("failed to create Ethereum client: %v", err)
    }
    // Initialize zkSync L2 client
    zksClient, err := zksync.NewRPCClient(cfg.ZkSync.RPCURL, cfg.ZkSync.Timeout)
    if err != nil {
        log.Fatalf("failed to create zkSync client: %v", err)
    }
    // Create bridge service
    bridgeSvc := bridge.NewService(ethClient, zksClient)
    // Create and start HTTP server
    srv := api.NewServer(bridgeSvc)
    addr := ":" + cfg.HTTPPort
    log.Printf("starting HTTP server on %s", addr)
    if err := srv.Start(addr); err != nil {
        log.Fatalf("HTTP server exited: %v", err)
    }
}