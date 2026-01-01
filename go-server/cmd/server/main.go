package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"l2alchemy/internal/config"
	"l2alchemy/internal/eth"
	"l2alchemy/internal/logging"
	"l2alchemy/internal/memtime"
	"l2alchemy/internal/oracle_runtime"
	"l2alchemy/internal/oracle_runtime/events"
	"l2alchemy/internal/server/handlers"
	"l2alchemy/internal/zkkeys"

	servers "l2alchemy/internal/server"
)

var (
	txCSVPathMu    sync.Mutex
	txCSVPathCache = map[string]string{}
)

func main() {
	logFile, err := logging.Setup(resolveLogFile())
	if err != nil {
		log.Fatalf("failed to set up logging: %v", err)
	}
	defer logFile.Close()

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

	txCSVPath := resolveTxCSVLogFile(cfg.NodeCount, cfg.BatchSize)
	if _, err := eth.SetupTxReceiptCSV(txCSVPath); err != nil {
		log.Fatalf("failed to set up tx receipt csv logging: %v", err)
	}
	defer eth.CloseTxReceiptCSV()

	memTimePath := memtime.ResolveMemTimeCSVPath(cfg.NodeCount, cfg.BatchSize)
	if _, err := memtime.SetupMemTimeCSV(memTimePath); err != nil {
		log.Fatalf("failed to set up memtime csv logging: %v", err)
	}
	defer memtime.CloseMemTimeCSV()

	ctx := context.Background()
	counterClient, err := eth.NewCounterClient(ctx, cfg.RPCURL, cfg.ChainID, cfg.PrivateKey, cfg.CounterContractAddress)
	if err != nil {
		log.Fatalf("failed to initialize counter client: %v", err)
	}

	keyDir := resolveKeyDir()
	zkMgr := zkkeys.New(keyDir)
	engine, err := oracle_runtime.Init(cfg, keyDir)
	if err != nil {
		log.Fatalf("failed to initialize oracle runtime: %v", err)
	}
	log.Printf("oracle runtime ready!")
	if engine.Oracle != nil {
		engine.Oracle.Start()
	}

	l2Subscriber, err := events.NewL2ContractEventSubscriber(engine)
	if err != nil {
		log.Fatalf("failed to initialize l2 contract event subscriber: %v", err)
	}
	l2Subscriber.Start(context.Background())

	counterHandler := handlers.NewCounterHandler(counterClient)
	zkHandler := handlers.NewZKHandler(zkMgr)
	oracleHandler := handlers.NewOracleHandler(engine)
	mux := servers.NewRouter(counterHandler, zkHandler, oracleHandler)
	rootHandler := servers.WithRequestLogging(mux)

	srv := &http.Server{
		Addr:    cfg.HTTPBindAddr,
		Handler: rootHandler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()
	const serverStartLogWait = 1 * time.Second
	go func() {
		select {
		case <-l2Subscriber.Ready():
		case <-time.After(serverStartLogWait):
		}
		log.Printf("starting server on %s", cfg.HTTPBindAddr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down server")
	l2Subscriber.Stop()
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

func resolveLogFile() string {
	return filepath.Join(resolveLogDir(), "server.log")
}

func resolveTxCSVLogFile(nodeCount, batchSize int) string {
	key := fmt.Sprintf("n%d_b%d", nodeCount, batchSize)

	txCSVPathMu.Lock()
	if path, ok := txCSVPathCache[key]; ok {
		txCSVPathMu.Unlock()
		return path
	}
	txCSVPathMu.Unlock()

	dir := resolveLogDir()
	filename := fmt.Sprintf("n%d_b%d.csv", nodeCount, batchSize)
	path := filepath.Join(dir, filename)

	txCSVPathMu.Lock()
	txCSVPathCache[key] = path
	txCSVPathMu.Unlock()

	return path
}

func resolveLogDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return filepath.Join("logs")
	}
	if st, err := os.Stat(filepath.Join(wd, "go-server")); err == nil && st.IsDir() {
		return filepath.Join(wd, "go-server", "logs")
	}
	return filepath.Join(wd, "logs")
}
