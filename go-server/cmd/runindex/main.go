package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"l2alchemy/internal/runindex"
)

func main() {
	logDir := resolveLogDir()
	logPath := filepath.Join(logDir, "run_index")
	idx, err := runindex.Advance(logPath)
	if err != nil {
		log.Fatalf("run index advance failed: %v", err)
	}
	fmt.Printf("run index: %d\n", idx)
}

func resolveLogDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "logs"
	}
	if st, err := os.Stat(filepath.Join(wd, "go-server")); err == nil && st.IsDir() {
		return filepath.Join(wd, "go-server", "logs")
	}
	return filepath.Join(wd, "logs")
}
