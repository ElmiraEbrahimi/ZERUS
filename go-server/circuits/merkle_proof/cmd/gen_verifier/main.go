package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"l2alchemy/internal/zkkeys"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	keyDir := resolveKeyDir()
	paths := zkkeys.PathsFor(keyDir, zkkeys.CircuitMerkleProof)
	if err := os.MkdirAll(filepath.Dir(paths.PK), 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", filepath.Dir(paths.PK), err)
	}
	force := strings.EqualFold(os.Getenv("FORCE_ZK_KEYGEN"), "1") || strings.EqualFold(os.Getenv("FORCE_ZK_KEYGEN"), "true")
	_, _, vk, _, _, err := zkkeys.GenerateKeysToFiles(zkkeys.CircuitMerkleProof, paths.PK, paths.VK, force)
	if err != nil {
		log.Fatalf("generate keys: %v", err)
	}
	outDir := "../contracts/src"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}

	f, err := os.Create(outDir + "/MerkleProofVerifier.sol")
	if err != nil {
		log.Fatalf("create verifier file: %v", err)
	}
	defer f.Close()

	if err := vk.ExportSolidity(f); err != nil {
		log.Fatalf("export solidity: %v", err)
	}

	log.Println("Generated contracts/src/MerkleProofVerifier.sol")
}

func resolveKeyDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return filepath.Join("circuits", "build", "keys")
	}
	if st, err := os.Stat(filepath.Join(wd, "go-server")); err == nil && st.IsDir() {
		return filepath.Join(wd, "go-server", "circuits", "build", "keys")
	}
	return filepath.Join(wd, "circuits", "build", "keys")
}
