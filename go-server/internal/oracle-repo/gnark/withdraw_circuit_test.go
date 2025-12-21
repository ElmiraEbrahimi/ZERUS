package gnark

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
)

// IncrementalMerkleTree represents a Merkle Tree structure.
type IncrementalMerkleTree struct {
	Depth          int            // Depth of the tree
	ZeroValues     [][]byte       // Precomputed zero values, padded to 32 bytes
	Leaves         [][]byte       // Leaves of the tree
	FilledSubtrees map[int][]byte // Per-level filled subtrees
	Roots          [][]byte       // History of all roots
}

// Initialize a new IncrementalMerkleTree with the given depth.
func NewIncrementalMerkleTree(depth int) (*IncrementalMerkleTree, error) {
	zeroValues, err := GenerateZeroValuesIncremental(depth)
	if err != nil {
		return nil, err
	}
	for i := 0; i < depth; i++ {
		fmt.Printf("Zero Value at Level %d: %x\n", i, zeroValues[i])
	}

	return &IncrementalMerkleTree{
		Depth:          depth,
		ZeroValues:     zeroValues,
		Leaves:         [][]byte{},
		FilledSubtrees: make(map[int][]byte),
		Roots:          [][]byte{},
	}, nil
}

// Add a new leaf to the tree and update the root.
func (tree *IncrementalMerkleTree) AddLeafValidator(commitmentHash []byte) (int, []byte, []byte, error) {
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

	return len(tree.Leaves) - 1, commitmentHash, currentHash, nil
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

// Helper function to generate fixed-length 32-byte random big.Int values
func generateRandomBigInt32Bytes(mod *big.Int) (*big.Int, error) {
	const byteLength = 32
	buf := make([]byte, byteLength)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	// Interpret the 32-byte array as a big.Int
	randomValue := new(big.Int).SetBytes(buf)
	// Ensure the value is reduced modulo the curve's scalar field
	randomValue.Mod(randomValue, mod)
	return randomValue, nil
}

func VerifyProofWithZeroValues(merkleRoot []byte, proofSet [][]byte, proofIndex uint64) bool {
	if merkleRoot == nil || len(proofSet) == 0 {
		return false
	}
	fmt.Println(" VerifyProofWithZeroValues Proof Path Length:", len(proofSet))
	h := hash.MIMC_BN254.New()
	// Start with the leaf (proofSet[0] is the commitment hash itself)
	currentHash := proofSet[0]
	var sibling []byte
	currentIndex := proofIndex
	for i := 1; i < len(proofSet); i++ {
		sibling = proofSet[i]

		h.Reset()
		if currentIndex%2 == 0 { // Left child
			h.Write(currentHash)
			h.Write(sibling)
		} else { // Right child
			h.Write(sibling)
			h.Write(currentHash)
		}
		currentHash = h.Sum(nil)
		fmt.Printf("VerifyProofWithZeroValues Level %d: currentHash=%x, sibling=%x\n", i, currentHash, sibling)

		currentIndex = currentIndex / 2
	}
	fmt.Println(" VerifyProofWithZeroValues Proof Path:")
	for i, hash := range proofSet {
		fmt.Printf("  Path[%d]: %x\n", i, hash)
	}
	fmt.Printf("\nVerifyProofWithZeroValues: Merkle Root: %x\n", currentHash)
	// Compare the calculated Merkle root to the expected root
	return bytes.Equal(currentHash, merkleRoot)
}

func TestMerkleProofCircuit(t *testing.T) {

	assert := test.NewAssert(t)
	hGo := hash.MIMC_BN254.New()
	mod := ecc.BN254.ScalarField()

	// Initialize the Incremental Merkle Tree
	depth := 10
	tree, err := NewIncrementalMerkleTree(depth)
	if err != nil {
		t.Fatalf("Error initializing tree: %v", err)
	}

	numLeaves := uint64(5)
	proofIndex := uint64(4) // Test the inclusion proof for leaf 4

	var nullifier, secret, destinationID *big.Int
	var commitmentHashIndexed, commitmentHash, nullifierHashIndexed []byte
	var nullifierIndexed *big.Int
	var secretIndexed *big.Int

	// Step 1: Generate random leaves
	for i := uint64(0); i < numLeaves; i++ {
		nullifier, _ = generateRandomBigInt32Bytes(mod)
		assert.NoError(err)
		secret, _ = generateRandomBigInt32Bytes(mod)
		assert.NoError(err)

		mimcHash := hash.MIMC_BN254.New()
		mimcHash.Reset()
		mimcHash.Write(PadTo32Bytes(nullifier))
		mimcHash.Write(PadTo32Bytes(secret))
		destinationID = new(big.Int)
		destinationID.SetString("9636219578937187601590327046728695236698322465209974782280717458744997515735", 10)
		mimcHash.Write(PadTo32Bytes(destinationID))
		commitmentHash = mimcHash.Sum(nil)

		// Step 2: Validator adds leaves to the tree
		index, commitmentHash, root, _ := tree.AddLeafValidator(commitmentHash)
		fmt.Printf("Validator: Leaf %d added, Commitment Hash: %x, Root: %x\n", index, commitmentHash, root)

		if i == proofIndex {
			commitmentHashIndexed = commitmentHash
			nullifierIndexed = nullifier
			secretIndexed = secret
			hGo.Reset()
			hGo.Write(nullifier.Bytes())
			nullifierHashIndexed = hGo.Sum(nil)

			fmt.Printf("\nProof Index Leaf Hash: %x\n", commitmentHashIndexed)
			fmt.Printf("\nGenerated Nullifier Hash:%x\n", nullifierHashIndexed)
		}
	}

	// Step 3: Build the Merkle proof
	merkleRoot, proofPath, err := tree.GetProofPath(proofIndex)
	if err != nil {
		t.Fatalf("Failed to generate Merkle proof: %v", err)
	}
	fmt.Println("--- Generated Merkle Proof Path ---")
	fmt.Printf("\nWithdrawer: Merkle Root: %x\n", merkleRoot)
	fmt.Println("Proof Path:")
	for i, hash := range proofPath {
		fmt.Printf("  Path[%d]: %x\n", i, hash)
	}

	// // Step 4: Verify proof in plain Go
	// verified := VerifyProofWithZeroValues(merkleRoot, proofPath, proofIndex)
	// fmt.Println("\n--- Proof Verification ---")
	// fmt.Printf("\nMerkle Root for Verification:%x\n", merkleRoot)
	// fmt.Println("Proof Path Length:", len(proofPath))
	// fmt.Println("Proof Verification Status:", verified)
	// if !verified {
	// 	t.Fatal("Merkle proof verification failed in plain Go")
	// }

	/////////////////////////////////////////CIRCUIT/////////////////////////////////////////

	// Step 1: Create the witness
	var witness MerkleProofCircuit
	var circuit MerkleProofCircuit
	circuit.M.Path = make([]frontend.Variable, depth+1)

	witness.Nullifier = nullifierIndexed
	witness.Secret = secretIndexed
	witness.DestinationID = destinationID
	witness.Leaf = proofIndex // Index of the leaf (= IncTreeIndex)
	// witness.CommitmentHash = commitmentHashIndexed // (= commitmentHash)
	witness.NullifierHash = nullifierHashIndexed
	witness.M.RootHash = merkleRoot //

	witness.M.Path = make([]frontend.Variable, depth+1)
	for i := 0; i < depth+1; i++ {
		witness.M.Path[i] = frontend.Variable(new(big.Int).SetBytes(proofPath[i]))
	}

	// Step 2: Compile the circuit
	// Initialize all empty slice with correct length
	fmt.Printf("Circuit Path Length: %d\n", len(circuit.M.Path))
	fmt.Printf("Witness Path Length: %d\n", len(witness.M.Path))
	fmt.Printf("Expected Depth: %d\n", depth)
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		t.Fatalf("Failed to compile circuit: %v", err)
	}

	// Step 3: Setup the proving/verifying keys
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		t.Fatalf("Failed to set up Groth16: %v", err)
	}

	// Step 4: Generate the proof
	fullWitness, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatalf("Failed to create full witness: %v", err)
	}

	proof, err := groth16.Prove(r1cs, pk, fullWitness)
	if err != nil {
		t.Fatalf("Failed to generate Groth16 proof: %v", err)
	}

	// Step 5: Verify the proof
	publicWitness, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		t.Fatalf("Failed to create public witness: %v", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		t.Fatalf("Failed to verify Groth16 proof: %v", err)
	}

	fmt.Println("Merkle proof and MiMC hash verified successfully!")

	//
	//
	//  SECOND TIME VERIFICATION
	//
	//

	numLeaves_secondRound := uint64(6)
	proofIndex_secondRound := uint64(5)
	// Step 1: Generate random leaves
	var nullifier_secondRound, secret_secondRound, destinationID_secondRound *big.Int
	var commitmentHashIndexed_secondRound, commitmentHash_secondRound, nullifierHashIndexed_secondRound []byte
	var nullifierIndexed_secondRound *big.Int
	var secretIndexed_secondRound *big.Int
	startIndex := numLeaves // Start second round indexing from last round's last index + 1

	for i := startIndex; i < startIndex+numLeaves_secondRound; i++ {
		//for i := uint64(0); i < numLeaves_secondRound; i++ {
		nullifier_secondRound, _ = generateRandomBigInt32Bytes(mod)
		assert.NoError(err)
		secret_secondRound, _ = generateRandomBigInt32Bytes(mod)
		assert.NoError(err)

		mimcHash := hash.MIMC_BN254.New()
		mimcHash.Reset()
		mimcHash.Write(PadTo32Bytes(nullifier_secondRound))
		mimcHash.Write(PadTo32Bytes(secret_secondRound))
		destinationID_secondRound = new(big.Int)
		destinationID_secondRound.SetString("9636219578937187601590327046728695236698322465209974782280717458744997515735", 10)
		mimcHash.Write(PadTo32Bytes(destinationID_secondRound))
		commitmentHash_secondRound = mimcHash.Sum(nil)

		// Step 2: Validator adds leaves to the tree
		index, commitmentHash, root, _ := tree.AddLeafValidator(commitmentHash_secondRound)
		fmt.Printf("Validator: Leaf %d added, Commitment Hash: %x, Root: %x\n", index, commitmentHash, root)

		if i == proofIndex_secondRound {
			commitmentHashIndexed_secondRound = commitmentHash_secondRound
			nullifierIndexed_secondRound = nullifier_secondRound
			secretIndexed_secondRound = secret_secondRound
			hGo.Reset()
			hGo.Write(nullifier_secondRound.Bytes())
			nullifierHashIndexed_secondRound = hGo.Sum(nil)

			fmt.Printf("\nProof Index Leaf Hash secondRound: %x\n", commitmentHashIndexed_secondRound)
			fmt.Printf("\nGenerated Nullifier Hash secondRound:%x\n", nullifierHashIndexed_secondRound)
		}
	}

	// Step 3: Build the Merkle proof
	merkleRoot_secondRound, proofPath_secondRound, err := tree.GetProofPath(proofIndex_secondRound)
	if err != nil {
		t.Fatalf("Failed to generate Merkle proof: %v", err)
	}
	fmt.Println("--- Generated Merkle Proof Path ---")
	fmt.Printf("\nWithdrawer: Merkle Root: %x\n", merkleRoot_secondRound)
	fmt.Println("Proof Path:")
	for i, hash := range proofPath_secondRound {
		fmt.Printf("  Path[%d]: %x\n", i, hash)
	}

	// // Step 4: Verify proof in plain Go
	// verified := VerifyProofWithZeroValues(merkleRoot, proofPath, proofIndex)
	// fmt.Println("\n--- Proof Verification ---")
	// fmt.Printf("\nMerkle Root for Verification:%x\n", merkleRoot)
	// fmt.Println("Proof Path Length:", len(proofPath))
	// fmt.Println("Proof Verification Status:", verified)
	// if !verified {
	// 	t.Fatal("Merkle proof verification failed in plain Go")
	// }

	/////////////////////////////////////////CIRCUIT/////////////////////////////////////////

	// Step 1: Create the witness
	var witness_secondRound MerkleProofCircuit // NEW INSTANCE for second round

	witness_secondRound.Nullifier = nullifierIndexed_secondRound
	witness_secondRound.Secret = secretIndexed_secondRound
	witness_secondRound.DestinationID = destinationID_secondRound
	witness_secondRound.Leaf = proofIndex_secondRound
	// witness_secondRound.CommitmentHash = commitmentHashIndexed_secondRound
	witness_secondRound.NullifierHash = nullifierHashIndexed_secondRound
	witness_secondRound.M.RootHash = merkleRoot_secondRound

	witness_secondRound.M.Path = make([]frontend.Variable, depth+1)
	for i := 0; i < depth+1; i++ {
		witness_secondRound.M.Path[i] = frontend.Variable(new(big.Int).SetBytes(proofPath_secondRound[i]))
	}

	// Step 2: Compile the circuit
	// Initialize all empty slice with correct length
	fmt.Printf("Circuit Path Length: %d\n", len(circuit.M.Path))
	fmt.Printf("Witness Path Length: %d\n", len(witness_secondRound.M.Path))
	fmt.Printf("Expected Depth: %d\n", depth)
	// r1cs, err = frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	// if err != nil {
	// 	t.Fatalf("Failed to compile circuit: %v", err)
	// }

	// Step 3: Setup the proving/verifying keys

	// Step 4: Generate the proof
	fullWitness_secondRound, err := frontend.NewWitness(&witness_secondRound, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatalf("Failed to create full witness: %v", err)
	}
	fmt.Println(" Debugging Second Proof:")
	fmt.Printf("Nullifier: %v\n", witness_secondRound.Nullifier)
	fmt.Printf("Secret: %v\n", witness_secondRound.Secret)
	fmt.Printf("DestinationID: %v\n", witness_secondRound.DestinationID)
	// fmt.Printf("CommitmentHash: %v\n", witness_secondRound.CommitmentHash)
	fmt.Printf("NullifierHash: %v\n", witness_secondRound.NullifierHash)
	fmt.Printf("Merkle Root: %v\n", witness_secondRound.M.RootHash)
	proof_secondRound, err := groth16.Prove(r1cs, pk, fullWitness_secondRound)
	if err != nil {
		t.Fatalf("Failed to generate Groth16 proof: %v", err)
	}

	// Step 5: Verify the proof
	publicWitness_secondRound, err := frontend.NewWitness(&witness_secondRound, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		t.Fatalf("Failed to create public witness: %v", err)
	}

	err = groth16.Verify(proof_secondRound, vk, publicWitness_secondRound)
	if err != nil {
		t.Fatalf("Failed to verify Groth16 proof: %v", err)
	}

	fmt.Println("Merkle proof and MiMC hash verified successfully!")

}
