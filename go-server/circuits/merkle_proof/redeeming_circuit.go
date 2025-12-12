package merkleproof

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// MerkleDepth is the fixed depth of the Merkle tree (number of levels of siblings).
// Path[0] is the leaf; Path[1..MerkleDepth] are the sibling nodes.
const MerkleDepth = 3

type MerkleProof struct {
	RootHash frontend.Variable                  `gnark:",public"` // Public: Merkle root
	Path     [MerkleDepth + 1]frontend.Variable // leaf + MerkleDepth siblings
}

// MerkleProofCircuit represents the ZK circuit for verifying a Merkle proof.
type MerkleProofCircuit struct {
	M    MerkleProof       // Public: Merkle root
	Leaf frontend.Variable // The leaf value to be verified

	Nullifier     frontend.Variable // secret input
	Secret        frontend.Variable // secret input
	DestinationID frontend.Variable // Receiver's ID (could be a wallet address or identifier)

	NullifierHash frontend.Variable `gnark:",public"`
}


// nodeSum calculates the hash for a parent node from two children.
func NodeSumW(api frontend.API, h mimc.MiMC, a, b frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(a)
	h.Write(b)
	return h.Sum()
}

// VerifyProof defines the logic for verifying a Merkle proof.
func (mp *MerkleProof) VerifyProofIncremental(api frontend.API, h mimc.MiMC, leaf frontend.Variable) {
	depth := MerkleDepth
	sum := mp.Path[0]

	// Binary decomposition of the leaf index determines the path.
	// `leaf` is interpreted as the index of the leaf in the tree.
	binLeaf := api.ToBinary(leaf, depth)
	api.Println("[Circuit] Binary Leaf:", binLeaf)

	for i := 1; i <= depth; i++ { // i up to depth, which matches Path[1..depth]
		h.Reset()
		d1 := api.Select(binLeaf[i-1], mp.Path[i], sum)
		d2 := api.Select(binLeaf[i-1], sum, mp.Path[i])
		sum = NodeSumW(api, h, d1, d2)
	}

	api.Println("[Circuit] calculated rootHash", sum)
	api.Println("[Withdrawer] calculated rootHash", mp.RootHash)
	api.AssertIsEqual(sum, mp.RootHash)
}

// Define defines the constraints for the Merkle proof circuit.
func (circuit *MerkleProofCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	// Compute the nullifierHash
	h.Reset()
	h.Write(circuit.Nullifier)
	nullifierHash := h.Sum()
	api.Println("[Circuit] Nullifier:", nullifierHash)
	api.Println("[out] NullifierHash:", circuit.NullifierHash)
	// Enforce that the circuit's computed nullifierHash matches the public input
	api.AssertIsEqual(nullifierHash, circuit.NullifierHash)

	// Reset the MiMC hash for the next computation

	// Compute the hash using nullifier and secret
	h.Reset()
	h.Write(circuit.Nullifier)
	h.Write(circuit.Secret)
	h.Write(circuit.DestinationID)
	// h.Write(circuit.Credentials)

	hash := h.Sum()

	// Enforce that the circuit's computed hash matches the public hash
	api.Println("[Circuit] hash:", hash)

	circuit.M.VerifyProofIncremental(api, h, circuit.Leaf)
	api.AssertIsEqual(hash, circuit.M.Path[0])

	return nil

}
