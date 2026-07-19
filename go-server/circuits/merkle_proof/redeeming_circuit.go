package merkleproof

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// // MerkleProof represents the Merkle proof structure used in the circuit.
type MerkleProof struct {
	RootHash frontend.Variable   `gnark:",public"` // Public: Merkle root
	Path     []frontend.Variable // Proof path (sibling hashes the first element is the leaf itself)
}

// MerkleProofCircuit represents the ZK circuit for verifying a Merkle proof.
type MerkleProofCircuit struct {
	// Public inputs. Keep this declaration order aligned with Algorithm 1:
	// d_dst, R_comm, h_n, a_dst.
	DestinationID    frontend.Variable `gnark:",public"`
	M                MerkleProof
	NullifierHash    frontend.Variable `gnark:",public"`
	RecipientAddress frontend.Variable `gnark:",public"`
	Leaf             frontend.Variable

	Nullifier frontend.Variable
	Secret    frontend.Variable
}

// LeafSum calculates the hash for a leaf node.
func LeafSumW(api frontend.API, h mimc.MiMC, data frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(data)
	return h.Sum()
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
	depth := len(mp.Path) - 1
	sum := mp.Path[0]

	//  Binary decomposition of the leaf index determines the path
	binLeaf := api.ToBinary(leaf, depth)
	api.Println("[Circuit] Binary Leaf:", binLeaf)

	for i := 1; i < len(mp.Path); i++ {
		h.Reset()
		d1 := api.Select(binLeaf[i-1], mp.Path[i], sum)
		d2 := api.Select(binLeaf[i-1], sum, mp.Path[i])
		sum = NodeSumW(api, h, d1, d2)
	}

	// Check if the calculated root matches the provided root
	api.Println("[Circuit] calculated rootHash", sum)
	api.Println("[Withdrawer] expected rootHash", mp.RootHash)
	api.AssertIsEqual(sum, mp.RootHash)
}

// Define defines the constraints for the Merkle proof circuit.
func (circuit *MerkleProofCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	h.Reset()
	h.Write(circuit.Nullifier)
	nullifierHash := h.Sum()
	api.Println("[Circuit] Nullifier:", nullifierHash)
	api.Println("[out] NullifierHash:", circuit.NullifierHash)
	// Enforce that the circuit's computed nullifierHash matches the public input
	api.AssertIsEqual(nullifierHash, circuit.NullifierHash)

	// Reset the MiMC hash for the next computation

	// C = H(n_rd, s_rd, d_dst, a_dst). Do not change this order.
	h.Reset()
	h.Write(circuit.Nullifier)
	h.Write(circuit.Secret)
	h.Write(circuit.DestinationID)
	h.Write(circuit.RecipientAddress)

	hash := h.Sum()

	// Enforce that the circuit's computed hash matches the public hash
	api.Println("[Circuit] hash:", hash)

	circuit.M.VerifyProofIncremental(api, h, circuit.Leaf)
	api.AssertIsEqual(hash, circuit.M.Path[0])

	return nil

}
