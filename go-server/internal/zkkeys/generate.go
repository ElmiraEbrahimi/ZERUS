package zkkeys

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"l2alchemy/circuits/merkle_proof"
	"l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/memtime"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark-crypto/ecc"
)

var supportedCircuits = []CircuitName{
	CircuitMerkleProof,
	CircuitVotingBatch,
}

// SupportedCircuits returns all circuit identifiers that can be compiled here.
func SupportedCircuits() []CircuitName {
	out := make([]CircuitName, len(supportedCircuits))
	copy(out, supportedCircuits)
	return out
}

// GenerateKeysToFiles compiles the circuit and runs Groth16 Setup, writing pk/vk to pkPath and vkPath.
// It returns the compiled constraint system, keys, and a coarse estimate of memory usage (MB) and elapsed time (ms).
// If force is false and both pk/vk files exist, it skips setup and loads the keys from disk.
func GenerateKeysToFiles(c CircuitName, pkPath, vkPath string, force bool) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, int, int, error) {
	var circuit frontend.Circuit
	switch c {
	case CircuitMerkleProof:
		merkleCircuit, err := merkleCircuitFromEnv()
		if err != nil {
			return nil, nil, nil, 0, 0, err
		}
		circuit = merkleCircuit
	case CircuitVotingBatch:
		votingCircuit, err := votingBatchCircuitFromEnv()
		if err != nil {
			return nil, nil, nil, 0, 0, err
		}
		circuit = votingCircuit
	default:
		return nil, nil, nil, 0, 0, fmt.Errorf("unknown circuit: %s", c)
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	started := time.Now()

	compileLabel := fmt.Sprintf("frontend.Compile %s", c)
	compileSample := memtime.Start(compileLabel)
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	compileSample.End()
	if err != nil {
		return nil, nil, nil, 0, 0, fmt.Errorf("compile %s: %w", c, err)
	}

	if !force && fileExists(pkPath) && fileExists(vkPath) {
		pk, err := readProvingKey(pkPath)
		if err != nil {
			return nil, nil, nil, 0, 0, fmt.Errorf("read pk %s: %w", c, err)
		}
		vk, err := readVerifyingKey(vkPath)
		if err != nil {
			return nil, nil, nil, 0, 0, fmt.Errorf("read vk %s: %w", c, err)
		}

		runtime.ReadMemStats(&m2)
		elapsed := time.Since(started)

		// TotalAlloc is monotonic for the process; a delta around compile+load is a pragmatic signal.
		usedBytes := int64(0)
		if m2.TotalAlloc > m1.TotalAlloc {
			usedBytes = int64(m2.TotalAlloc - m1.TotalAlloc)
		}
		memMB := int(usedBytes / (1024 * 1024))
		timeMS := int(elapsed.Milliseconds())
		return cs, pk, vk, memMB, timeMS, nil
	}

	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return nil, nil, nil, 0, 0, fmt.Errorf("groth16 setup %s: %w", c, err)
	}

	// Write keys atomically (write to temp then rename) to avoid partial files.
	if err := writeKeyAtomic(pkPath, func(f *os.File) error {
		_, werr := pk.WriteRawTo(f)
		return werr
	}); err != nil {
		return nil, nil, nil, 0, 0, fmt.Errorf("write pk %s: %w", c, err)
	}

	if err := writeKeyAtomic(vkPath, func(f *os.File) error {
		_, werr := vk.WriteRawTo(f)
		return werr
	}); err != nil {
		return nil, nil, nil, 0, 0, fmt.Errorf("write vk %s: %w", c, err)
	}

	runtime.ReadMemStats(&m2)
	elapsed := time.Since(started)

	// TotalAlloc is monotonic for the process; a delta around compile+setup is a pragmatic signal.
	usedBytes := int64(0)
	if m2.TotalAlloc > m1.TotalAlloc {
		usedBytes = int64(m2.TotalAlloc - m1.TotalAlloc)
	}
	memMB := int(usedBytes / (1024 * 1024))
	timeMS := int(elapsed.Milliseconds())
	return cs, pk, vk, memMB, timeMS, nil
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func readProvingKey(path string) (groth16.ProvingKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pk := groth16.NewProvingKey(ecc.BN254)
	if _, err := pk.ReadFrom(f); err != nil {
		return nil, err
	}
	return pk, nil
}

func readVerifyingKey(path string) (groth16.VerifyingKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vk := groth16.NewVerifyingKey(ecc.BN254)
	if _, err := vk.ReadFrom(f); err != nil {
		return nil, err
	}
	return vk, nil
}

func writeKeyAtomic(path string, writeFn func(*os.File) error) error {
	tmp := path + ".tmp"
	// Best-effort cleanup.
	_ = os.Remove(tmp)

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := writeFn(f); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

// CompileOnly is a helper for tests or future enhancements.
func CompileOnly(c CircuitName) (constraint.ConstraintSystem, error) {
	switch c {
	case CircuitMerkleProof:
		merkleCircuit, err := merkleCircuitFromEnv()
		if err != nil {
			return nil, err
		}
		sample := memtime.Start("frontend.Compile merkle_proof")
		cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, merkleCircuit)
		sample.End()
		return cs, err
	case CircuitVotingBatch:
		votingCircuit, err := votingBatchCircuitFromEnv()
		if err != nil {
			return nil, err
		}
		sample := memtime.Start("frontend.Compile voting_batch")
		cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, votingCircuit)
		sample.End()
		return cs, err
	default:
		return nil, fmt.Errorf("unknown circuit: %s", c)
	}
}

func merkleCircuitFromEnv() (*merkleproof.MerkleProofCircuit, error) {
	depthStr := os.Getenv("INC_TREE_DEPTH")
	if depthStr == "" {
		return nil, fmt.Errorf("INC_TREE_DEPTH is required to size the Merkle proof path")
	}
	depth, err := strconv.Atoi(depthStr)
	if err != nil || depth < 1 {
		return nil, fmt.Errorf("invalid INC_TREE_DEPTH %q", depthStr)
	}

	circuit := &merkleproof.MerkleProofCircuit{}
	circuit.M.Path = make([]frontend.Variable, depth+1)
	return circuit, nil
}

func votingBatchCircuitFromEnv() (*votingbatch.BatchingVotingCircuit, error) {
	depthStr := os.Getenv("SPARSE_TREE_DEPTH")
	if depthStr == "" {
		return nil, fmt.Errorf("SPARSE_TREE_DEPTH is required to size the Merkle proof path")
	}
	depth, err := strconv.Atoi(depthStr)
	if err != nil || depth < 1 {
		return nil, fmt.Errorf("invalid SPARSE_TREE_DEPTH %q", depthStr)
	}

	batchStr := os.Getenv("BATCH_SIZE")
	if batchStr == "" {
		return nil, fmt.Errorf("BATCH_SIZE is required to size withdrawal request IDs")
	}
	batchSize, err := strconv.Atoi(batchStr)
	if err != nil || batchSize < 1 {
		return nil, fmt.Errorf("invalid BATCH_SIZE %q", batchStr)
	}

	nodeCountStr := os.Getenv("NODE_COUNT")
	if nodeCountStr == "" {
		return nil, fmt.Errorf("NODE_COUNT is required to size validator list")
	}
	nodeCount, err := strconv.Atoi(nodeCountStr)
	if err != nil || nodeCount < 1 {
		return nil, fmt.Errorf("invalid NODE_COUNT %q", nodeCountStr)
	}

	circuit := &votingbatch.BatchingVotingCircuit{}
	circuit.Validators = make([]votingbatch.BatchingValidatorConstraints, nodeCount)
	circuit.WithdrawalReqIDs = make([]frontend.Variable, batchSize)
	circuit.Aggregator.MerkleProof.Path = make([]frontend.Variable, depth+1)
	for i := range circuit.Validators {
		circuit.Validators[i].MerkleProof.Path = make([]frontend.Variable, depth+1)
	}
	return circuit, nil
}
