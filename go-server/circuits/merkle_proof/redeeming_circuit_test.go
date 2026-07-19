package merkleproof

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func hashFields(values ...*big.Int) *big.Int {
	h := hash.MIMC_BN254.New()
	for _, value := range values {
		var element fr.Element
		element.SetBigInt(value)
		encoded := element.Bytes()
		_, _ = h.Write(encoded[:])
	}
	return new(big.Int).SetBytes(h.Sum(nil))
}

func TestRedeemingProofBindsDestinationAndRecipient(t *testing.T) {
	nullifier := big.NewInt(11)
	secret := big.NewInt(22)
	destination := big.NewInt(33)
	recipient := big.NewInt(44)
	sibling := big.NewInt(55)
	commitment := hashFields(nullifier, secret, destination, recipient)
	root := hashFields(commitment, sibling)
	nullifierHash := hashFields(nullifier)

	circuit := &MerkleProofCircuit{M: MerkleProof{Path: make([]frontend.Variable, 2)}}
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		t.Fatal(err)
	}
	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		t.Fatal(err)
	}
	assignment := &MerkleProofCircuit{
		DestinationID: destination,
		NullifierHash: nullifierHash, RecipientAddress: recipient,
		M:    MerkleProof{RootHash: root, Path: []frontend.Variable{commitment, sibling}},
		Leaf: 0, Nullifier: nullifier, Secret: secret,
	}
	full, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
	proof, err := groth16.Prove(cs, pk, full)
	if err != nil {
		t.Fatal(err)
	}
	public, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		t.Fatal(err)
	}
	if err := groth16.Verify(proof, vk, public); err != nil {
		t.Fatalf("correct public inputs must verify: %v", err)
	}
	vector, ok := public.Vector().(fr.Vector)
	if !ok || len(vector) != 4 {
		t.Fatalf("unexpected public witness vector: %T len=%d", public.Vector(), len(vector))
	}
	wantOrder := []*big.Int{destination, root, nullifierHash, recipient}
	for i := range wantOrder {
		if vector[i].BigInt(new(big.Int)).Cmp(wantOrder[i]) != 0 {
			t.Fatalf("public input %d is out of order", i)
		}
	}

	assertTamperedFails := func(name string, tampered *MerkleProofCircuit) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			w, err := frontend.NewWitness(tampered, ecc.BN254.ScalarField(), frontend.PublicOnly())
			if err != nil {
				t.Fatal(err)
			}
			if err := groth16.Verify(proof, vk, w); err == nil {
				t.Fatal("copied proof verified with changed public input")
			}
		})
	}
	changedDestination := *assignment
	changedDestination.DestinationID = big.NewInt(34)
	assertTamperedFails("destination", &changedDestination)
	changedRecipient := *assignment
	changedRecipient.RecipientAddress = big.NewInt(45)
	assertTamperedFails("recipient", &changedRecipient)
	copiedProof := *assignment
	copiedProof.RecipientAddress = big.NewInt(99)
	assertTamperedFails("copied_proof_recipient_replacement", &copiedProof)
}
