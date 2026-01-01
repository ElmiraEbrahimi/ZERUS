package votingbatch

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// MerkleProofW represents the Merkle proof structure used across circuits.
type MerkleProofW struct {
	RootHash frontend.Variable
	Path     []frontend.Variable // Proof path (sibling hashes, first element is the leaf itself)
}

// LeafSum calculates the hash for a leaf node.
func LeafSum(api frontend.API, h mimc.MiMC, data frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(data)
	return h.Sum()
}

// NodeSum calculates the hash for a parent node from two children.
func NodeSum(api frontend.API, h mimc.MiMC, a, b frontend.Variable) frontend.Variable {
	h.Reset()      // Reset only once
	h.Write(a)     // First input
	h.Write(b)     // Second input
	return h.Sum() // Compute the hash

}

// MerkleProof Methods

// VerifyProof for MerkleProof
func (mp *MerkleProofW) VerifyProof(api frontend.API, h mimc.MiMC, leaf frontend.Variable) {
	depth := len(mp.Path) - 1
	sum := LeafSum(api, h, mp.Path[0])

	// Binary decomposition of the leaf index determines the path
	binLeaf := api.ToBinary(leaf, depth)
	api.Println("[Circuit] Binary Leaf:", binLeaf)

	for i := 1; i < len(mp.Path); i++ {
		h.Reset()
		d1 := api.Select(binLeaf[i-1], mp.Path[i], sum)
		d2 := api.Select(binLeaf[i-1], sum, mp.Path[i])
		sum = NodeSum(api, h, d1, d2)
	}

	// Check if the calculated root matches the provided root
	api.Println("[Circuit_voting verify proof] calculated rootHash", sum)
	api.Println("[out_voting verify proof] calculated rootHash", mp.RootHash)
	api.AssertIsEqual(sum, mp.RootHash)
}

// ComputeRootFromPath for MerkleProof
func (mp *MerkleProofW) ComputeRootFromPath(api frontend.API, h mimc.MiMC, leaf frontend.Variable) frontend.Variable {
	depth := len(mp.Path) - 1
	sum := LeafSum(api, h, mp.Path[0])

	// Binary decomposition of the leaf index determines the path
	binLeaf := api.ToBinary(leaf, depth)
	api.Println("[Circuit] Binary Leaf:", binLeaf)

	for i := 1; i < len(mp.Path); i++ {
		h.Reset()
		d1 := api.Select(binLeaf[i-1], mp.Path[i], sum)
		d2 := api.Select(binLeaf[i-1], sum, mp.Path[i])
		sum = NodeSum(api, h, d1, d2)
	}

	return sum
}
