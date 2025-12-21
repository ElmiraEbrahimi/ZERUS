package gnark

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

const (
	voteSize = 96 // 32 bytes each for Index, Request, and Vote
)

// Vote represents a validator's vote.
type Vote struct {
	Index     uint64          // Validator's index
	Request   *big.Int        // Request ID
	Vote      *big.Int        // 1 for valid, 0 for invalid
	Sender    eddsa.PublicKey // Validator's public key
	Signature eddsa.Signature // Validator's signature of the vote
}

// Serialize converts the Vote into a fixed-size byte array.
func (v *Vote) Serialize() []byte {
	var b [voteSize]byte

	// Serialize the Index (padded to 32 bytes)
	copy(b[:32], PadOrTrim(big.NewInt(int64(v.Index)).Bytes(), 32))

	// Serialize the Request ID (padded to 32 bytes)
	copy(b[32:64], PadOrTrim(v.Request.Bytes(), 32))

	// Serialize the Vote value (padded to 32 bytes)
	copy(b[64:voteSize], PadOrTrim(v.Vote.Bytes(), 32))

	return b[:]
}
