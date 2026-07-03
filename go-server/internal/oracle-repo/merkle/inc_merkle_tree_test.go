package merkle

import (
	"fmt"
	"l2alchemy/internal/oracle-repo/util"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
)

func TestValidatorAndWithdrawer(t *testing.T) {
	depth := 5
	zeroValues, _ := GenerateZeroValues(depth)
	tree, err := NewIncrementalMerkleTree(zeroValues, depth)
	if err != nil {
		t.Fatalf("Error initializing tree: %v", err)
	}

	mod := ecc.BN254.ScalarField()
	destinationID := new(big.Int)
	destinationID.SetString("9636219578937187601590327046728695236698322465209974782280717458744997515735", 10)

	// Validator adds 6 nodes
	numLeaves := 6
	for i := 0; i < numLeaves; i++ {
		nullifier, _ := util.GenerateRandomBigInt32Bytes(mod)
		secret, _ := util.GenerateRandomBigInt32Bytes(mod)

		mimcHash := hash.MIMC_BN254.New()
		mimcHash.Reset()
		mimcHash.Write(util.PadTo32Bytes(nullifier))
		mimcHash.Write(util.PadTo32Bytes(secret))
		mimcHash.Write(util.PadTo32Bytes(destinationID))
		commitmentHash := mimcHash.Sum(nil)

		fmt.Printf("Validator: Commitment Hash: %x\n", commitmentHash)
		if _, _, _, err := tree.AddLeafValidator(commitmentHash); err != nil {
			t.Fatalf("Error adding leaf: %v", err)
		}
	}

	// Withdrawer generates proof for leaf 5
	index := uint64(4)
	root, proofPath, err := tree.GetProofPath(index)
	if err != nil {
		t.Fatalf("Error getting proof path: %v", err)
	}

	fmt.Printf("\nWithdrawer: Merkle Root: %x\n", root)
	fmt.Println("Proof Path:")
	for i, hash := range proofPath {
		fmt.Printf("  Path[%d]: %x\n", i, hash)
	}
}

// TestGetProofPathOpensLatestRootForEveryLeaf asserts the property the claim
// flow depends on: after any number of insertions, the proof path of EVERY
// leaf (not just the most recently inserted one) recomputes to the tree's
// latest root, which is the root the aggregator publishes on the Gateway.
func TestGetProofPathOpensLatestRootForEveryLeaf(t *testing.T) {
	depth := 5
	zeroValues, err := GenerateZeroValues(depth)
	if err != nil {
		t.Fatalf("generate zero values: %v", err)
	}
	tree, err := NewIncrementalMerkleTree(zeroValues, depth)
	if err != nil {
		t.Fatalf("initialize tree: %v", err)
	}

	mod := ecc.BN254.ScalarField()
	numLeaves := 6
	for i := 0; i < numLeaves; i++ {
		leaf, _ := util.GenerateRandomBigInt32Bytes(mod)
		if _, _, _, err := tree.AddLeafValidator(util.PadTo32Bytes(leaf)); err != nil {
			t.Fatalf("add leaf %d: %v", i, err)
		}
	}

	latest := tree.LatestRoot()
	for index := uint64(0); index < uint64(numLeaves); index++ {
		root, path, err := tree.GetProofPath(index)
		if err != nil {
			t.Fatalf("proof path for leaf %d: %v", index, err)
		}
		if fmt.Sprintf("%x", root) != fmt.Sprintf("%x", latest) {
			t.Fatalf("leaf %d: proof root %x != latest root %x", index, root, latest)
		}
		if len(path) != depth+1 {
			t.Fatalf("leaf %d: path length %d, want %d", index, len(path), depth+1)
		}

		// Replay the circuit's fold: start at the leaf value, combine with
		// each sibling according to the index bits.
		mimcHash := hash.MIMC_BN254.New()
		sum := path[0]
		pos := index
		for level := 1; level <= depth; level++ {
			mimcHash.Reset()
			if pos%2 == 0 {
				mimcHash.Write(sum)
				mimcHash.Write(path[level])
			} else {
				mimcHash.Write(path[level])
				mimcHash.Write(sum)
			}
			sum = mimcHash.Sum(nil)
			pos /= 2
		}
		if fmt.Sprintf("%x", sum) != fmt.Sprintf("%x", latest) {
			t.Fatalf("leaf %d: recomputed root %x != latest root %x", index, sum, latest)
		}
	}
}
