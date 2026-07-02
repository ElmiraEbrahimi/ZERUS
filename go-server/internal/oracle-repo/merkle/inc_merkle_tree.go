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

// Generate the Merkle proof path for a given leaf index.
func (tree *IncrementalMerkleTree) GetProofPath(index uint64) ([]byte, [][]byte, error) {
	if index >= uint64(len(tree.Leaves)) {
		return nil, nil, fmt.Errorf("invalid index")
	}

	mimcHash := hash.MIMC_BN254.New()
	filledSubtrees := make(map[int][]byte) // Cache filled left nodes
	proofPath := [][]byte{}
	var merkleRoot []byte // Final Merkle root

	// Step 1: Rebuild the tree from leaf 0 to the proof index
	for i := uint64(0); i <= index; i++ {
		currentLeafHash := tree.Leaves[i]
		currentLeafIndex := i

		for level := 0; level < tree.Depth; level++ {
			var sibling []byte
			siblingIndex := currentLeafIndex ^ 1 // Flip last bit to find sibling

			// Determine sibling value
			if currentLeafIndex%2 == 0 { // Left child
				// Right sibling exists at level 0
				if level == 0 && siblingIndex < uint64(len(tree.Leaves)) && siblingIndex <= index {
					sibling = tree.Leaves[siblingIndex]
				} else {
					sibling = tree.ZeroValues[level] // Use zero value for right siblings
				}
			} else { // Right child
				// Fetch left sibling from cache **only if the sibling index is even**
				if siblingIndex%2 == 0 {
					if val, exists := filledSubtrees[level]; exists {
						sibling = val // Use cached intermediate left node
					} else {
						sibling = tree.ZeroValues[level] // Default to zero if no cache
					}
				} else {
					sibling = tree.ZeroValues[level] // Force zero for odd siblings
				}
			}

			// If processing the target leaf, capture the sibling for the proof path
			if i == index {
				proofPath = append(proofPath, sibling)
			}

			// Cache the left child at this level
			if currentLeafIndex%2 == 0 {
				filledSubtrees[level] = currentLeafHash
			}

			// Compute parent hash
			mimcHash.Reset()
			if currentLeafIndex%2 == 0 { // Left child
				mimcHash.Write(currentLeafHash)
				mimcHash.Write(sibling)
			} else { // Right child
				mimcHash.Write(sibling)
				mimcHash.Write(currentLeafHash)
			}

			currentLeafHash = mimcHash.Sum(nil)
			currentLeafIndex /= 2 // Move to parent level

			// Capture the final root when at the last level
			if i == index && level == tree.Depth-1 {
				merkleRoot = currentLeafHash
			}
		}
	}
	// Append the leaf value itself to the proof path
	proofPath = append([][]byte{tree.Leaves[index]}, proofPath...)

	return merkleRoot, proofPath, nil
}
