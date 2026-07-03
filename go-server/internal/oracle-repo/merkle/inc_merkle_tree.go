package merkle

import (
	"fmt"
	"l2alchemy/internal/oracle-repo/util"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	_ "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/hash"
)

// IncrementalMerkleTree represents a Merkle Tree structure.
type IncrementalMerkleTree struct {
	Depth          int            // Depth of the tree
	ZeroValues     [][]byte       // Precomputed zero values, padded to 32 bytes
	Leaves         [][]byte       // Leaves of the tree
	FilledSubtrees map[int][]byte // Per-level filled subtrees (similar to Tornado Cash's filledSubtrees)
	Roots          [][]byte       // History of all roots (one for each leaf insertion)
}

// NewIncrementalMerkleTree initializes a new IncrementalMerkleTree with the given depth.
func NewIncrementalMerkleTree(zeroValues [][]byte, depth int) (*IncrementalMerkleTree, error) {
	return &IncrementalMerkleTree{
		Depth:          depth,
		ZeroValues:     zeroValues,
		Leaves:         [][]byte{},
		FilledSubtrees: make(map[int][]byte),
		Roots:          [][]byte{}, // Start as empty slice
	}, nil
}

func GenerateZeroValues(depth int) ([][]byte, error) {
	mod := ecc.BN254.ScalarField()
	mimcHash := hash.MIMC_BN254.New()

	zeroValues := make([][]byte, depth)
	initialZeroValue := new(big.Int)
	initialZeroValue.SetString("20019762671335178393512154978075455201849332419823879662510519485824706883752", 10)
	initialZeroValue.Mod(initialZeroValue, mod)

	zeroValues[0] = util.PadTo32Bytes(initialZeroValue)
	for i := 1; i < depth; i++ {
		mimcHash.Reset()
		mimcHash.Write(zeroValues[i-1])
		mimcHash.Write(zeroValues[i-1])
		zeroValues[i] = mimcHash.Sum(nil)
	}

	return zeroValues, nil
}

// Add a new leaf to the tree and update the root.
func (tree *IncrementalMerkleTree) AddLeafValidator(commitmentHash []byte) (uint64, []byte, []byte, error) {
	mimcHash := hash.MIMC_BN254.New()

	// Append the new leaf
	tree.Leaves = append(tree.Leaves, commitmentHash)
	currentHash := commitmentHash
	currentIndex := len(tree.Leaves) - 1

	// Compute hashes level by level
	for i := 0; i < tree.Depth; i++ {
		var sibling []byte

		// Determine sibling
		if currentIndex%2 == 0 { // Left child
			sibling = tree.ZeroValues[i]
			tree.FilledSubtrees[i] = currentHash // Update filled subtree for this level
		} else { // Right child
			sibling = tree.FilledSubtrees[i]
		}

		// Compute parent hash
		mimcHash.Reset()
		if currentIndex%2 == 0 { // Left child
			mimcHash.Write(currentHash) // Left child comes first
			mimcHash.Write(sibling)     // Sibling (right child)
		} else { // Right child
			mimcHash.Write(sibling)     // Left sibling comes first
			mimcHash.Write(currentHash) // Right child
		}
		currentHash = mimcHash.Sum(nil)

		currentIndex /= 2
	}

	// Append the root to history
	tree.Roots = append(tree.Roots, currentHash)

	return uint64(len(tree.Leaves) - 1), commitmentHash, currentHash, nil
}

// LatestRoot returns the most recent root, or nil if the tree is empty.
func (tree *IncrementalMerkleTree) LatestRoot() []byte {
	if len(tree.Roots) == 0 {
		return nil
	}
	return tree.Roots[len(tree.Roots)-1]
}

// GetProofPath generates the Merkle proof path for a given leaf index.
//
// The proof must open the tree's LATEST root — the root the aggregator
// published on the Gateway — regardless of which position the target leaf
// occupies. The tree is therefore rebuilt level by level over ALL current
// leaves (zero-padded on the right), and the siblings of `index` are read
// off each level. (A previous version rebuilt only leaves 0..index, so
// proofs for any leaf other than the most recently inserted one opened a
// stale root and were rejected against the published commitment root.)
func (tree *IncrementalMerkleTree) GetProofPath(index uint64) ([]byte, [][]byte, error) {
	if index >= uint64(len(tree.Leaves)) {
		return nil, nil, fmt.Errorf("invalid index")
	}

	mimcHash := hash.MIMC_BN254.New()

	level := make([][]byte, len(tree.Leaves))
	copy(level, tree.Leaves)

	// proofPath[0] is the leaf value itself; proofPath[1..Depth] are the
	// sibling hashes, matching the circuit's VerifyProofIncremental layout.
	proofPath := [][]byte{tree.Leaves[index]}
	pos := index

	for d := 0; d < tree.Depth; d++ {
		sibPos := pos ^ 1
		var sibling []byte
		if sibPos < uint64(len(level)) {
			sibling = level[sibPos]
		} else {
			sibling = tree.ZeroValues[d]
		}
		proofPath = append(proofPath, sibling)

		next := make([][]byte, (len(level)+1)/2)
		for i := 0; i < len(level); i += 2 {
			left := level[i]
			var right []byte
			if i+1 < len(level) {
				right = level[i+1]
			} else {
				right = tree.ZeroValues[d]
			}
			mimcHash.Reset()
			mimcHash.Write(left)
			mimcHash.Write(right)
			next[i/2] = mimcHash.Sum(nil)
		}
		level = next
		pos /= 2
	}

	merkleRoot := level[0]
	return merkleRoot, proofPath, nil
}
