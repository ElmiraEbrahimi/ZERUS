package oracle_runtime

import (
	"log"
	"os"
	"path/filepath"

	merkleproof "l2alchemy/circuits/merkle_proof"
	votingbatch "l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/config"
	"l2alchemy/internal/zkkeys"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

func setupMerkleCircuit(cfg *config.Config, keyDir string) (*merkleproof.MerkleProofCircuit, constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey) {
	log.Println("setting up merkle circuit...")
	var circuit merkleproof.MerkleProofCircuit
	for i := 0; i < len(circuit.M.Path); i++ {
		circuit.M.Path[i] = 0
	}

	paths := zkkeys.PathsFor(keyDir, zkkeys.CircuitMerkleProof)
	pkPath := paths.PK
	vkPath := paths.VK
	if err := os.MkdirAll(filepath.Dir(pkPath), 0o755); err != nil {
		panic(err)
	}
	_r1cs, pk, vk, _, _, err := zkkeys.GenerateKeysToFiles(zkkeys.CircuitMerkleProof, pkPath, vkPath, false)
	if err != nil {
		panic(err)
	}
	log.Println("completed set up merkle circuit")

	return &circuit, _r1cs, pk, vk

}

func setupVotingCircuit(cfg *config.Config, keyDir string) (*votingbatch.BatchingVotingCircuit, constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey) {
	log.Println("setting up voting circuit...")
	var circuit votingbatch.BatchingVotingCircuit
	for i := 0; i < cfg.NodeCount; i++ {
		circuit.Validators[i] = votingbatch.BatchingValidatorConstraints{}
		path := make([]frontend.Variable, cfg.SparseTreeDepth+1)
		copy(circuit.Validators[i].MerkleProof.Path[:], path)
	}
	aggregatorPath := make([]frontend.Variable, cfg.SparseTreeDepth+1)
	copy(circuit.Aggregator.MerkleProof.Path[:], aggregatorPath)

	withdrawalReqIDs := make([]frontend.Variable, cfg.BatchSize)
	copy(circuit.WithdrawalReqIDs[:], withdrawalReqIDs)

	paths := zkkeys.PathsFor(keyDir, zkkeys.CircuitVotingBatch)
	pkPath := paths.PK
	vkPath := paths.VK
	if err := os.MkdirAll(filepath.Dir(pkPath), 0o755); err != nil {
		panic(err)
	}
	_r1cs, pk, vk, _, _, err := zkkeys.GenerateKeysToFiles(zkkeys.CircuitVotingBatch, pkPath, vkPath, false)
	if err != nil {
		panic(err)
	}
	log.Println("completed set up voting circuit")

	return &circuit, _r1cs, pk, vk
}
