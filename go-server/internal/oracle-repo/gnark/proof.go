// Package proof provides utilities for handling zkSNARK proofs and converting them into formats
// compatible with Ethereum smart contracts.
//
// The ProofToEthereumProof function takes a zkSNARK proof and converts it into the EthereumProof format.
// This format splits the proof into its components (A, B, and C) as required by Ethereum's zkSNARK verification contracts.
package gnark

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark/backend/groth16"
)

// EthereumProof represents the zkSNARK proof in the format required by Ethereum's smart contract verifier.
type EthereumProof struct {
	A [2]*big.Int    // Component A of the proof (two field elements).
	B [2][2]*big.Int // Component B of the proof (two pairs of field elements).
	C [2]*big.Int    // Component C of the proof (two field elements).
}

// ProofToEthereumProof converts a groth16 zkSNARK proof into EthereumProof format.
// The function reads the raw byte representation of the proof and extracts each field element (A, B, C)
// based on their predefined positions in the proof structure.
func ProofToEthereumProof(p groth16.Proof) (*EthereumProof, error) {
	// Initialize an EthereumProof structure to hold the converted proof.
	var proof EthereumProof

	// Serialize the proof into a byte buffer using the WriteRawTo method.
	var buf bytes.Buffer
	_, err := p.WriteRawTo(&buf)
	if err != nil {
		return nil, fmt.Errorf("write raw proof to: %w", err)
	}
	proofBytes := buf.Bytes() // Get the raw bytes of the proof.

	// Extract components from the serialized proofBytes. Each component corresponds to a field element in the zkSNARK proof.

	// Extract A (two field elements).
	proof.A[0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*0 : fp.Bytes*1]) // First field element of A.
	proof.A[1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*1 : fp.Bytes*2]) // Second field element of A.

	// Extract B (two pairs of field elements).
	proof.B[0][0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*2 : fp.Bytes*3]) // First pair: First element.
	proof.B[0][1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*3 : fp.Bytes*4]) // First pair: Second element.
	proof.B[1][0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*4 : fp.Bytes*5]) // Second pair: First element.
	proof.B[1][1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*5 : fp.Bytes*6]) // Second pair: Second element.

	// Extract C (two field elements).
	proof.C[0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*6 : fp.Bytes*7]) // First field element of C.
	proof.C[1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*7 : fp.Bytes*8]) // Second field element of C.

	// Return the converted EthereumProof.
	return &proof, nil
}
