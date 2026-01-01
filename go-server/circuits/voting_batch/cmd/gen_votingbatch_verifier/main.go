package main

import (
	"log"
	"os"
	"strconv"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/joho/godotenv"

	votingbatch "l2alchemy/circuits/voting_batch"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	depthStr := os.Getenv("SPARSE_TREE_DEPTH")
	if depthStr == "" {
		log.Fatal("SPARSE_TREE_DEPTH is required to size the Merkle proof path")
	}
	depth, err := strconv.Atoi(depthStr)
	if err != nil || depth < 1 {
		log.Fatalf("invalid SPARSE_TREE_DEPTH %q", depthStr)
	}

	batchStr := os.Getenv("BATCH_SIZE")
	if batchStr == "" {
		log.Fatal("BATCH_SIZE is required to size withdrawal request IDs")
	}
	batchSize, err := strconv.Atoi(batchStr)
	if err != nil || batchSize < 1 {
		log.Fatalf("invalid BATCH_SIZE %q", batchStr)
	}

	var circuit votingbatch.BatchingVotingCircuit
	circuit.WithdrawalReqIDs = make([]frontend.Variable, batchSize)
	circuit.Aggregator.MerkleProof.Path = make([]frontend.Variable, depth+1)
	for i := 0; i < len(circuit.Validators); i++ {
		circuit.Validators[i].MerkleProof.Path = make([]frontend.Variable, depth+1)
	}

	cs, err := frontend.Compile(
		ecc.BN254.ScalarField(),
		r1cs.NewBuilder,
		&circuit,
	)
	if err != nil {
		log.Fatalf("compile circuit: %v", err)
	}

	_, vk, err := groth16.Setup(cs)
	if err != nil {
		log.Fatalf("setup: %v", err)
	}

	outDir := "../contracts/src"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}

	f, err := os.Create(outDir + "/VotingBatchVerifier.sol")
	if err != nil {
		log.Fatalf("create verifier file: %v", err)
	}
	defer f.Close()

	if err := vk.ExportSolidity(f); err != nil {
		log.Fatalf("export solidity: %v", err)
	}

	log.Println("Generated contracts/src/VotingBatchVerifier.sol")
}
