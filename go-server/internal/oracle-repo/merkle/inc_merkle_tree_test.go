package merkle

import (
	"fmt"
	"math/big"
	"l2alchemy/internal/oracle-repo/util"
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
