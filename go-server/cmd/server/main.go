package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"l2alchemy/internal/config"
	"l2alchemy/internal/eth"

	// "l2alchemy/internal/oracle_runtime"
	servers "l2alchemy/internal/server"
	"l2alchemy/internal/zkkeys"
)

func main() {
	// Best-effort load of a .env file (repo root).
	//
	// Supported run modes:
	//   - from repo root:       go run ./go-server/cmd/server
	//   - from go-server/:      go run ./cmd/server

	printEnvValues := true
	cfg, err := config.Load(printEnvValues)
	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	ctx := context.Background()
	counterClient, err := eth.NewCounterClient(ctx, cfg.RPCURL, cfg.ChainID, cfg.PrivateKey, cfg.CounterContractAddress)
	if err != nil {
		log.Fatalf("failed to initialise counter client: %v", err)
	}

	// oracleEngine, err := oracle_runtime.Init(cfg)
	// if err != nil {
	// 	log.Fatalf("failed to initialise oracle runtime: %v", err)
	// }
	// log.Printf("oracle runtime ready: state=%s users=%d", len(oracleEngine.Users))

	keyDir := resolveKeyDir()
	zkMgr := zkkeys.New(keyDir)
	handler := servers.NewHandler(counterClient, zkMgr)
	mux := http.NewServeMux()
	servers.RegisterRoutes(mux, handler)
	rootHandler := servers.WithRequestLogging(mux)

	srv := &http.Server{
		Addr:    cfg.HTTPBindAddr,
		Handler: rootHandler,
	}

	go func() {
		log.Printf("starting server on %s", cfg.HTTPBindAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during shutdown: %v", err)
	}
}

// resolveKeyDir returns a stable folder for pk/vk files, regardless of whether
// the binary is started from repo root or from within go-server/.
func resolveKeyDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return filepath.Join("circuits", "build", "keys")
	}
	// If started from repo root, prefer go-server/build/keys.
	if st, err := os.Stat(filepath.Join(wd, "go-server")); err == nil && st.IsDir() {
		return filepath.Join(wd, "go-server", "circuits", "build", "keys")
	}
	// Otherwise assume current working dir is go-server/.
	return filepath.Join(wd, "circuits", "build", "keys")
}
