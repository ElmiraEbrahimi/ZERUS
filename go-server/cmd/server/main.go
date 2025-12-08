package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "counterapi/internal/config"
    "counterapi/internal/eth"
    servers "counterapi/internal/server"

    "github.com/joho/godotenv"
)

func main() {
    // Best-effort load of a local .env file.
    if err := godotenv.Load(); err != nil {
        log.Printf("no .env file found or unable to load it: %v", err)
    }

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("unable to load configuration: %v", err)
    }

    ctx := context.Background()
    counterClient, err := eth.NewCounterClient(ctx, cfg.RPCURL, cfg.ChainID, cfg.PrivateKey, cfg.ContractAddress)
    if err != nil {
        log.Fatalf("failed to initialise counter client: %v", err)
    }

    handler := servers.NewHandler(counterClient)
    mux := http.NewServeMux()
    servers.RegisterRoutes(mux, handler)

    srv := &http.Server{
        Addr:    cfg.HTTPBindAddr,
        Handler: mux,
    }

    go func() {
        log.Printf("starting counter API on %s", cfg.HTTPBindAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen error: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop
    log.Println("shutting down counter API")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("error during shutdown: %v", err)
    }
}