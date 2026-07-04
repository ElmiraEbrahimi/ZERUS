package zkkeys

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"l2alchemy/circuits/merkle_proof"
	"l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/memtime"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
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
// If an r1cs cache exists alongside the keys, it is reused to avoid recompiling.
func GenerateKeysToFiles(
	c CircuitName,
	pkPath, vkPath string,
	force bool,
) (
	constraint.ConstraintSystem,
	groth16.ProvingKey,
	groth16.VerifyingKey,
	int, int,
	int, int,
	error,
) {

	// -------------------------------
	// 1) Resolve circuit builder
	// -------------------------------
	var buildCircuit func() (frontend.Circuit, error)
	switch c {
	case CircuitMerkleProof:
		buildCircuit = func() (frontend.Circuit, error) {
			return merkleCircuitFromEnv()
		}
	case CircuitVotingBatch:
		buildCircuit = func() (frontend.Circuit, error) {
			return votingBatchCircuitFromEnv()
		}
	default:
		return nil, nil, nil, 0, 0, 0, 0, fmt.Errorf("unknown circuit: %s", c)
	}

	// -------------------------------
	// 2) COMPILE — PEAK MEMORY
	// -------------------------------
	r1csPath := r1csPathFromPK(pkPath)
	var cs constraint.ConstraintSystem
	compilePeakMB := 0
	compileTimeMS := 0

	if !force {
		cached, err := readR1CS(r1csPath)
		if err == nil {
			cs = cached
		}
	}

	if cs == nil {
		circuit, err := buildCircuit()
		if err != nil {
			return nil, nil, nil, 0, 0, 0, 0, err
		}
		compileRes, err := memtime.MeasurePeak(
			fmt.Sprintf("frontend.Compile %s", c),
			memtime.PeakSampleInterval,
			func() error {
				var e error
				cs, e = frontend.Compile(
					ecc.BN254.ScalarField(),
					r1cs.NewBuilder,
					circuit,
				)
				return e
			},
		)
		if err != nil {
			return nil, nil, nil, 0, 0, 0, 0, fmt.Errorf("compile %s: %w", c, err)
		}

		if err := writeR1CSAtomic(r1csPath, cs); err != nil {
			log.Printf("zkkeys: write r1cs %s: %v", r1csPath, err)
		}

		compilePeakMB = memtime.BytesToMB(compileRes.PeakBytes)
		compileTimeMS = int(compileRes.Time.Milliseconds())
	}

	// -------------------------------
	// 3) LOAD KEYS IF EXIST (NO SETUP)
	// -------------------------------
	if !force && fileExists(pkPath) && fileExists(vkPath) {
		pk, err := readProvingKey(pkPath)
		if err != nil {
			return nil, nil, nil, 0, 0, 0, 0, err
		}
		vk, err := readVerifyingKey(vkPath)
		if err != nil {
			return nil, nil, nil, 0, 0, 0, 0, err
		}

		// Setup NOT executed → setup peak = 0
		return cs, pk, vk,
			compilePeakMB, compileTimeMS,
			0, 0,
			nil
	}

	// -------------------------------
	// 4) SETUP — PEAK MEMORY
	// -------------------------------
	var pk groth16.ProvingKey
	var vk groth16.VerifyingKey

	setupRes, err := memtime.MeasurePeak(
		fmt.Sprintf("groth16.Setup %s", c),
		memtime.PeakSampleInterval,
		func() error {
			var e error
			pk, vk, e = groth16.Setup(cs)
			return e
		},
	)
	if err != nil {
		return nil, nil, nil, 0, 0, 0, 0, fmt.Errorf("groth16 setup %s: %w", c, err)
	}

	setupPeakMB := memtime.BytesToMB(setupRes.PeakBytes)
	setupTimeMS := int(setupRes.Time.Milliseconds())

	// -------------------------------
	// 5) WRITE KEYS (not measured)
	// -------------------------------
	if err := writeKeyAtomic(pkPath, func(f *os.File) error {
		_, e := pk.WriteRawTo(f)
		return e
	}); err != nil {
		return nil, nil, nil, 0, 0, 0, 0, err
	}

	if err := writeKeyAtomic(vkPath, func(f *os.File) error {
		_, e := vk.WriteRawTo(f)
		return e
	}); err != nil {
		return nil, nil, nil, 0, 0, 0, 0, err
	}

	return cs, pk, vk,
		compilePeakMB, compileTimeMS,
		setupPeakMB, setupTimeMS,
		nil
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

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
	_, err = pk.ReadFrom(f)
	return pk, err
}

func readVerifyingKey(path string) (groth16.VerifyingKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vk.ReadFrom(f)
	return vk, err
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

func r1csPathFromPK(pkPath string) string {
	return filepath.Join(filepath.Dir(pkPath), "r1cs")
}

func readR1CS(path string) (cs constraint.ConstraintSystem, err error) {
	defer func() {
		if r := recover(); r != nil {
			cs = nil
			err = fmt.Errorf("read r1cs %s: panic: %v", path, r)
		}
	}()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cs = groth16.NewCS(ecc.BN254)
	if _, err := cs.ReadFrom(f); err != nil {
		return nil, err
	}
	return cs, nil
}

func writeR1CSAtomic(path string, cs constraint.ConstraintSystem) error {
	tmp := path + ".tmp"
	_ = os.Remove(tmp)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := cs.WriteTo(f); err != nil {
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

// -----------------------------------------------------------------------------
// Circuit builders (unchanged)
// -----------------------------------------------------------------------------

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
	circuit.Validators = make([]votingbatch.BatchingValidatorConstraints, votingQuorumWitnessSize(nodeCount))
	circuit.WithdrawalReqIDs = make([]frontend.Variable, batchSize)
	circuit.Aggregator.MerkleProof.Path = make([]frontend.Variable, depth+1)
	for i := range circuit.Validators {
		circuit.Validators[i].MerkleProof.Path = make([]frontend.Variable, depth+1)
	}
	return circuit, nil
}

func votingQuorumWitnessSize(nValidators int) int {
	if nValidators <= 0 {
		return 0
	}
	f := (nValidators - 1) / 3
	return 2*f + 1
}
