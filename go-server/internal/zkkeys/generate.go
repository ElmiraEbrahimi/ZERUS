package zkkeys

import (
    "fmt"
    "os"
    "runtime"
    "time"

    "l2alchemy/circuits/merkle_proof"
    "l2alchemy/circuits/voting_batch"

    "github.com/consensys/gnark/backend/groth16"
    "github.com/consensys/gnark/constraint"
    "github.com/consensys/gnark/frontend"
    "github.com/consensys/gnark/frontend/cs/r1cs"
    "github.com/consensys/gnark-crypto/ecc"
)

// GenerateKeysToFiles compiles the circuit and runs Groth16 Setup, writing pk/vk to pkPath and vkPath.
// It returns a coarse estimate of memory usage (MB) and elapsed time (ms) measured around compile+setup.
func GenerateKeysToFiles(c CircuitName, pkPath, vkPath string) (int, int, error) {
    var circuit frontend.Circuit
    switch c {
    case CircuitMerkleProof:
        circuit = &merkleproof.MerkleProofCircuit{}
    case CircuitVotingBatch:
        circuit = &votingbatch.BatchingVotingCircuit{}
    default:
        return 0, 0, fmt.Errorf("unknown circuit: %s", c)
    }

    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    started := time.Now()

    cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
    if err != nil {
        return 0, 0, fmt.Errorf("compile %s: %w", c, err)
    }

    pk, vk, err := groth16.Setup(cs)
    if err != nil {
        return 0, 0, fmt.Errorf("groth16 setup %s: %w", c, err)
    }

    // Write keys atomically (write to temp then rename) to avoid partial files.
    if err := writeKeyAtomic(pkPath, func(f *os.File) error {
        _, werr := pk.WriteRawTo(f)
        return werr
    }); err != nil {
        return 0, 0, fmt.Errorf("write pk %s: %w", c, err)
    }

    if err := writeKeyAtomic(vkPath, func(f *os.File) error {
        _, werr := vk.WriteRawTo(f)
        return werr
    }); err != nil {
        return 0, 0, fmt.Errorf("write vk %s: %w", c, err)
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
    return memMB, timeMS, nil
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
        return frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &merkleproof.MerkleProofCircuit{})
    case CircuitVotingBatch:
        return frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &votingbatch.BatchingVotingCircuit{})
    default:
        return nil, fmt.Errorf("unknown circuit: %s", c)
    }
}
