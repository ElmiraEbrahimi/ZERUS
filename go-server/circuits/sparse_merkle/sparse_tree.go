package merkle

import "github.com/consensys/gnark-crypto/hash"

type SparseMerkleTree struct {
	Depth      int
	Leaves     [][]byte
	ZeroValues [][]byte
}

func NewSparseMerkleTree(depth int, zeroValues [][]byte) (*SparseMerkleTree, error) {
	leafCount := 1 << depth
	leaves := make([][]byte, leafCount)
	for i := range leaves {
		leaves[i] = zeroValues[0] // Initialize leaves with zero values
	}
	return &SparseMerkleTree{
		Depth:      depth,
		Leaves:     leaves,
		ZeroValues: zeroValues,
	}, nil
}

func (tree *SparseMerkleTree) AddLeaf(index int, leaf []byte) {
	tree.Leaves[index] = leaf
}

func (tree *SparseMerkleTree) ComputeRoot() ([]byte, error) {
	mimcHash := hash.MIMC_BN254.New()
	currentLevel := tree.Leaves

	for level := 0; level < tree.Depth; level++ {
		nextLevel := make([][]byte, len(currentLevel)/2)
		for i := 0; i < len(currentLevel); i += 2 {
			mimcHash.Reset()
			mimcHash.Write(currentLevel[i])
			if i+1 < len(currentLevel) {
				mimcHash.Write(currentLevel[i+1])
			} else {
				mimcHash.Write(tree.ZeroValues[level])
			}
			nextLevel[i/2] = mimcHash.Sum(nil)
		}
		currentLevel = nextLevel
	}

	return currentLevel[0], nil
}

func (tree *SparseMerkleTree) GenerateProof(index int) ([][]byte, error) {
	proof := make([][]byte, tree.Depth+1)
	proof[tree.Depth] = tree.Leaves[index]
	currentIndex := index
	for level := 0; level < tree.Depth; level++ {
		siblingIndex := currentIndex ^ 1
		if siblingIndex < len(tree.Leaves) {
			proof[level] = tree.Leaves[siblingIndex]
		} else {
			proof[level] = tree.ZeroValues[level]
		}
		currentIndex /= 2
	}
	return proof, nil
}
