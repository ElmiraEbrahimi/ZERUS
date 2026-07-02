package user

import (
	"bytes"
	"math/big"
	"testing"

	"l2alchemy/internal/oracle-repo/util"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/hash"
)

// F-25 regression: the off-chain nullifier hash must be computed over the
// canonical 32-byte field encoding of the nullifier, matching the in-circuit
// MiMC hash of the field element, even when the nullifier has leading zero
// bytes (~1/256 of random draws).
func TestNullifierHashLeadingZeroPadding(t *testing.T) {
	// A nullifier whose big-endian encoding is shorter than 32 bytes.
	nullifier := new(big.Int).SetInt64(42)

	var elem fr.Element
	elem.SetBigInt(nullifier)
	canonical := elem.Marshal()

	padded := util.PadTo32Bytes(nullifier)
	if !bytes.Equal(canonical, padded) {
		t.Fatalf("PadTo32Bytes does not match the canonical field encoding: %x vs %x", padded, canonical)
	}

	h := hash.MIMC_BN254.New()
	h.Reset()
	if _, err := h.Write(canonical); err != nil {
		t.Fatalf("write canonical encoding: %v", err)
	}
	inCircuitStyle := h.Sum(nil)

	if got := NullifierHashBytes(nullifier); !bytes.Equal(got, inCircuitStyle) {
		t.Fatalf("NullifierHashBytes diverges from the field-element hash: %x vs %x", got, inCircuitStyle)
	}

	// Note: gnark-crypto v0.19.0's MiMC left-pads a single short Write, so
	// hashing the unpadded nullifier bytes happens to produce the same
	// digest today. The explicit padding (and this test) pins the canonical
	// encoding so the off-chain hash stays aligned with the in-circuit
	// H(n_rd) regardless of library-internal normalization behaviour.
}
