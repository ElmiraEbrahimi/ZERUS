package main

import (
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	merkleproof "l2alchemy/circuits/merkle_proof"
)

func main() {
	var circuit merkleproof.MerkleProofCircuit

	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
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
