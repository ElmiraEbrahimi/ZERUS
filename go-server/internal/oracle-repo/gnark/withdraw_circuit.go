package gnark

import (
	"fmt"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"

	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

// // MerkleProof represents the Merkle proof structure used in the circuit.
type MerkleProof struct {
	RootHash frontend.Variable   `gnark:",public"` // Public: Merkle root
	Path     []frontend.Variable // Proof path (sibling hashes the first element is the leaf itself)
}

// MerkleProofCircuit represents the ZK circuit for verifying a Merkle proof.
type MerkleProofCircuit struct {
	M    MerkleProof       // Public: Merkle root
	Leaf frontend.Variable // The leaf value to be verified

	Nullifier     frontend.Variable // secret input
	Secret        frontend.Variable // secret input
	DestinationID frontend.Variable // Receiver's ID (could be a wallet address or identifier)

	//***
	Credentials   frontend.Variable
	NullifierHash frontend.Variable `gnark:",public"` // public output

	CredentialsHash frontend.Variable `gnark:",public"` // public output
	Issuer          IssuerConstraints

	RevocationProof MerkleProof       // Merkle path for PoNR
	RevocationLeaf  frontend.Variable // Index of the credential hash in the revocation tree
}

// ***
type IssuerConstraints struct {
	PublicKey eddsa.PublicKey `gnark:",public"`
	Signature eddsa.Signature
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
	api.Println("[Withdrawer] calculated rootHash", mp.RootHash)
	api.AssertIsEqual(sum, mp.RootHash)
}

// Define defines the constraints for the Merkle proof circuit.
func (circuit *MerkleProofCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	//***
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return fmt.Errorf("curve initialization: %w", err)
	}

	//***
	h.Reset()
	h.Write(circuit.Credentials)
	CredentialsHash := h.Sum()

	api.Println("[Circuit] CredentialsHash:", CredentialsHash)
	api.Println("[out] CredentialsHash:", circuit.CredentialsHash)
	api.AssertIsEqual(CredentialsHash, circuit.CredentialsHash)

	//***
	//calculate msg
	msg := CredentialsHash
	h.Reset()
	api.Println("[withdraw circuit] Verifying EDDSA Signature for Issuer...")
	if err := eddsa.Verify(curve, circuit.Issuer.Signature, msg, circuit.Issuer.PublicKey, &h); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
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
	h.Write(circuit.Credentials)

	hash := h.Sum()

	// Enforce that the circuit's computed hash matches the public hash
	api.Println("[Circuit] hash:", hash)

	circuit.M.VerifyProofIncremental(api, h, circuit.Leaf)
	api.AssertIsEqual(hash, circuit.M.Path[0])

	// Revocation tree check (PoNR)
	h.Reset()
	api.Println("[Circuit] Verifying PoNR Merkle proof...")
	circuit.RevocationProof.VerifyProofIncremental(api, h, circuit.RevocationLeaf)
	return nil
	
}
