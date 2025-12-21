package gnark

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
)

// PadOrTrim adjusts a byte slice to a specified size by either padding with zeros or truncating the leading bytes.
// PadTo32Bytes pads a big.Int to a fixed 32-byte length (no truncation) and is specific to cryptographic contexts.

// PadOrTrim adjusts the byte slice to a fixed length by padding or truncating.
func PadOrTrim(bb []byte, size int) []byte {
	l := len(bb)
	if l == size {
		return bb
	}
	if l > size {
		return bb[l-size:]
	}
	tmp := make([]byte, size)
	copy(tmp[size-l:], bb)
	return tmp
}

// Helper function to pad values to 32 bytes
func PadTo32Bytes(value *big.Int) []byte {
	const byteLength = 32
	valueBytes := value.Bytes()
	padded := make([]byte, byteLength)
	copy(padded[byteLength-len(valueBytes):], valueBytes)
	return padded
}

// GenerateZeroValues calculates zero values for each level of the tree
func GenerateZeroValuesIncremental(depth int) ([][]byte, error) {
	mod := ecc.BN254.ScalarField()
	mimcHash := hash.MIMC_BN254.New()

	zeroValues := make([][]byte, depth)

	// Start with the zero value
	initialZeroValue := new(big.Int)
	initialZeroValue.SetString("4555114089170143013007615382799372902997177870479602349537353593038812875418", 10)
	initialZeroValue.Mod(initialZeroValue, mod)

	zeroValues[0] = PadTo32Bytes(initialZeroValue)
	for i := 1; i < depth; i++ {
		mimcHash.Reset()
		mimcHash.Write(zeroValues[i-1])
		mimcHash.Write(zeroValues[i-1])
		zeroValues[i] = mimcHash.Sum(nil)
	}

	return zeroValues, nil
}
