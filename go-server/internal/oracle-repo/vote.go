package oracle

import (
	"bytes"
	"l2alchemy/internal/oracle-repo/util"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

const wiVoteSize = 96

// region incvote

type IncVote struct {
	NodeID         uint
	CommitmentHash []byte
	IncTreeIndex   uint64
	RootHash       []byte
}

type ByCommitmentHash []*IncVote

func (s ByCommitmentHash) Len() int      { return len(s) }
func (s ByCommitmentHash) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByCommitmentHash) Less(i, j int) bool {
	a, b := s[i].CommitmentHash, s[j].CommitmentHash
	if len(a) != len(b) {
		return len(a) > len(b)
	}
	return bytes.Compare(a, b) > 0
}

// endregion

// region wivote

type WiVote struct {
	NodeID     uint
	Index      *big.Int
	RequestID  *big.Int
	IsApproved *big.Int
	SenderPK   eddsa.PublicKey
	Signature  []byte
}

func (v *WiVote) Serialize() []byte {
	var b [wiVoteSize]byte

	copy(b[:32], util.PadOrTrim(v.Index.Bytes(), 32))
	copy(b[32:64], util.PadOrTrim(v.RequestID.Bytes(), 32))
	copy(b[64:wiVoteSize], util.PadOrTrim(v.IsApproved.Bytes(), 32))

	return b[:]
}

type BatchedWiVote struct {
	Index            *big.Int
	WithdrawalReqIDs []*big.Int
	Vote             []*big.Int
	MajorityOfVotes  []*big.Int
	SenderPK         eddsa.PublicKey
	Signature        []byte
}

// endregion
