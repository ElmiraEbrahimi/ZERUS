package oracle_runtime

import (
	"fmt"
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
	circuit.M.Path = make([]frontend.Variable, cfg.IncTreeDepth+1)
	for i := 0; i < len(circuit.M.Path); i++ {
		circuit.M.Path[i] = 0
	}

	paths := zkkeys.PathsFor(keyDir, zkkeys.CircuitMerkleProof)
	pkPath := paths.PK
	vkPath := paths.VK
	if err := os.MkdirAll(filepath.Dir(pkPath), 0o755); err != nil {
		panic(err)
	}
	if !fileExists(pkPath) || !fileExists(vkPath) {
		panic(fmt.Errorf("merkle keys missing (pk=%s vk=%s); run `make generate` or `make deploy`", pkPath, vkPath))
	}
	_r1cs, pk, vk, _, _, _, _, err := zkkeys.GenerateKeysToFiles(zkkeys.CircuitMerkleProof, pkPath, vkPath, false)
	if err != nil {
		panic(err)
	}
	log.Println("completed set up merkle circuit")

	return &circuit, _r1cs, pk, vk

}

func setupVotingCircuit(cfg *config.Config, keyDir string) (*votingbatch.BatchingVotingCircuit, constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey) {
	log.Println("setting up voting circuit...")
	var circuit votingbatch.BatchingVotingCircuit
	quorumSize := votingQuorumWitnessSize(cfg.NodeCount)
	circuit.Validators = make([]votingbatch.BatchingValidatorConstraints, quorumSize)
	for i := 0; i < quorumSize; i++ {
		circuit.Validators[i] = votingbatch.BatchingValidatorConstraints{}
		circuit.Validators[i].MerkleProof.Path = make([]frontend.Variable, cfg.SparseTreeDepth+1)
	}
	circuit.Aggregator.MerkleProof.Path = make([]frontend.Variable, cfg.SparseTreeDepth+1)

	circuit.WithdrawalReqIDs = make([]frontend.Variable, cfg.BatchSize)

	paths := zkkeys.PathsFor(keyDir, zkkeys.CircuitVotingBatch)
	pkPath := paths.PK
	vkPath := paths.VK
	if err := os.MkdirAll(filepath.Dir(pkPath), 0o755); err != nil {
		panic(err)
	}
	if !fileExists(pkPath) || !fileExists(vkPath) {
		panic(fmt.Errorf("voting batch keys missing (pk=%s vk=%s); run `make generate` or `make deploy`", pkPath, vkPath))
	}
	_r1cs, pk, vk, _, _, _, _, err := zkkeys.GenerateKeysToFiles(zkkeys.CircuitVotingBatch, pkPath, vkPath, false)
	if err != nil {
		panic(err)
	}
	log.Println("completed set up voting circuit")

	return &circuit, _r1cs, pk, vk
}

func votingQuorumWitnessSize(nValidators int) int {
	if nValidators <= 0 {
		return 0
	}
	f := (nValidators - 1) / 3
	return 2*f + 1
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
