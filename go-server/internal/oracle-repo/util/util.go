package util

import (
	"context"
	crand "crypto/rand"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GenerateRandomBigInt32Bytes(mod *big.Int) (*big.Int, error) {
	const byteLength = 32
	buf := make([]byte, byteLength)
	_, err := crand.Read(buf)
	if err != nil {
		return nil, err
	}
	// Interpret the 32-byte array as a big.Int
	randomValue := new(big.Int).SetBytes(buf)
	// Ensure the value is reduced modulo the curve's scalar field
	randomValue.Mod(randomValue, mod)
	return randomValue, nil
}

func PadTo32Bytes(value *big.Int) []byte {
	const byteLength = 32
	valueBytes := value.Bytes()
	padded := make([]byte, byteLength)
	copy(padded[byteLength-len(valueBytes):], valueBytes)
	return padded
}

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

func GetEthBalance(ethClient *ethclient.Client, address common.Address) (*big.Int, error) {
	balance, err := ethClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

type EthereumProof struct {
	Proof [8]*big.Int
}

func ProofToEthereumProof(p groth16.Proof) (*EthereumProof, error) {
	var proof EthereumProof

	pn, ok := p.(*groth16_bn254.Proof)
	if !ok {
		return nil, fmt.Errorf("expected *groth16_bn254.Proof, got %T", p)
	}

	data := pn.MarshalSolidity()

	// Validate length (8 elements * 32 bytes = 256 bytes)
	if len(data) != 256 {
		return nil, fmt.Errorf("unexpected proof length: %v", len(data))
	}

	// Extract elements in order (no reordering needed)
	for i := 0; i < 8; i++ {
		start := i * 32
		end := start + 32
		proof.Proof[i] = new(big.Int).SetBytes(data[start:end])
	}

	return &proof, nil
}

func ModToBn254(i *big.Int) *big.Int {
	bn254Modulus, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	return new(big.Int).Mod(i, bn254Modulus)
}

func ModToBn254Bytes(data []byte) []byte {
	bn254Modulus, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	i := new(big.Int).SetBytes(data)
	return new(big.Int).Mod(i, bn254Modulus).Bytes()
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// func SetupCSVWriter(filepath string) *csv.Writer {
// 	fmt.Println("setting up csv writer...")
// 	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// 	// defer func(f *os.File) {
// 	// 	err := f.Close()
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }(f)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fileInfo, err := f.Stat()
// 	if err != nil {
// 		panic(fmt.Errorf("could not get csv file info: %v", err))
// 	}

// 	csvWriter := csv.NewWriter(f)

// 	// file is newly created:
// 	if fileInfo.Size() == 0 {
// 		headerRow := []string{
// 			"nodeCount",
// 			"provingTime",
// 			"provingMemory",
// 			"compileMemory",
// 			"compileTime",
// 			"datetime",
// 		}
// 		err = csvWriter.Write(headerRow)
// 		if err != nil {
// 			panic(fmt.Errorf("failed to write header row: %v", err))
// 		}
// 		csvWriter.Flush()
// 	}

//		return csvWriter
//	}
func SetupCSVWriter(filepath string, header []string) *csv.Writer {
	fmt.Println("setting up csv writer...")
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		panic(fmt.Errorf("could not get csv file info: %v", err))
	}

	csvWriter := csv.NewWriter(f)

	// write header if file is new
	if fileInfo.Size() == 0 && header != nil {
		err = csvWriter.Write(header)
		if err != nil {
			panic(fmt.Errorf("failed to write header row: %v", err))
		}
		csvWriter.Flush()
	}

	return csvWriter
}

func ReplaceVerifier(filePath, newName string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	oldStr := "contract Verifier {"
	newStr := "contract " + newName + " {"
	updatedContent := strings.Replace(string(data), oldStr, newStr, 1)

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func PowBatching(x, y *big.Int) *big.Int {
	output := big.NewInt(1)
	binaryY := ToBinary(y, 256)

	for i := 0; i < len(binaryY); i++ {
		if i != 0 {
			output.Mul(output, output)
		}
		multiply := new(big.Int).Mul(output, x)
		if binaryY[len(binaryY)-1-i] == 1 {
			output.Set(multiply)
		}
	}

	return output
}

func ToBinary(n *big.Int, nbDigits int) []int {
	binary := make([]int, nbDigits)
	for i := 0; i < nbDigits; i++ {
		binary[nbDigits-1-i] = int(new(big.Int).And(n, big.NewInt(1)).Int64())
		n.Rsh(n, 1)
	}
	return binary
}
