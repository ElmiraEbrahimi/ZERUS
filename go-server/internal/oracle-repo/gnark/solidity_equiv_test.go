package gnark

// Cross-implementation equivalence tests (plan F-18).
//
// Three MiMC/Merkle implementations must agree for on-chain verification to
// succeed: the Solidity MiMC library (contracts/src/mimc.sol) used by
// MerkleTree.verify/insert/update, the gnark in-circuit MiMC used by the
// Redeeming and Aggregating circuits, and gnark-crypto's MIMC_BN254 that
// produces the off-chain validator-state root. The gnark in-circuit MiMC is
// generated from gnark-crypto's constants, so proving Solidity == gnark-crypto
// on shared inputs covers all three. The Solidity logic is ported to Go here
// verbatim (keccak256 constant chain, 110 rounds, exponent 5,
// Miyaguchi-Preneel chaining) so any drift in mimc.sol parameters breaks this
// test.

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/crypto"
)

const solRounds = 110

var frModulus = fr.Modulus()

// solMiMCpe5 ports MiMC.MiMCpe5 from contracts/src/mimc.sol.
func solMiMCpe5(x, k *big.Int) *big.Int {
	c := crypto.Keccak256([]byte("seed"))
	out := new(big.Int).Set(x)
	t := new(big.Int)
	a := new(big.Int)
	for i := 0; i < solRounds; i++ {
		c = crypto.Keccak256(c)
		t.SetBytes(c)
		t.Add(t, out).Add(t, k).Mod(t, frModulus) // x + c_i + k
		a.Mul(t, t).Mod(a, frModulus)             // t^2
		a.Mul(a, a).Mod(a, frModulus)             // t^4
		out.Mul(a, t).Mod(out, frModulus)         // t^5
	}
	return out.Add(out, k).Mod(out, frModulus)
}

// solMiMCHash ports MiMC.hash (MiMCpe5_mp with key 0) from mimc.sol.
func solMiMCHash(msgs ...*big.Int) *big.Int {
	r := big.NewInt(0)
	for _, m := range msgs {
		e := solMiMCpe5(m, r)
		r = new(big.Int).Add(r, m)
		r.Add(r, e).Mod(r, frModulus)
	}
	return r
}

func gnarkMiMCHash(t *testing.T, msgs ...*big.Int) *big.Int {
	t.Helper()
	h := mimc.NewMiMC()
	for _, m := range msgs {
		buf := make([]byte, 32)
		m.FillBytes(buf)
		if _, err := h.Write(buf); err != nil {
			t.Fatalf("gnark mimc write: %v", err)
		}
	}
	return new(big.Int).SetBytes(h.Sum(nil))
}

func randomFieldElement(t *testing.T) *big.Int {
	t.Helper()
	v, err := rand.Int(rand.Reader, frModulus)
	if err != nil {
		t.Fatalf("rand: %v", err)
	}
	return v
}

// TestSolidityMiMCMatchesGnarkCrypto proves mimc.sol's hash equals
// gnark-crypto's MIMC_BN254 for 1-, 2- and 4-element inputs (leafSum,
// hashLeftRight, and hashAccount shapes respectively).
func TestSolidityMiMCMatchesGnarkCrypto(t *testing.T) {
	for _, n := range []int{1, 2, 4} {
		for trial := 0; trial < 8; trial++ {
			msgs := make([]*big.Int, n)
			for i := range msgs {
				msgs[i] = randomFieldElement(t)
			}
			sol := solMiMCHash(msgs...)
			ref := gnarkMiMCHash(t, msgs...)
			if sol.Cmp(ref) != 0 {
				t.Fatalf("n=%d trial=%d mismatch:\n solidity     = %s\n gnark-crypto = %s", n, trial, sol, ref)
			}
		}
	}
}

// TestSolidityZeroSubtreeTable pins the zeros() table in merkle_tree.sol:
// zeros(i+1) = MiMC(zeros(i) || zeros(i)).
func TestSolidityZeroSubtreeTable(t *testing.T) {
	want := []string{
		"4555114089170143013007615382799372902997177870479602349537353593038812875418",
		"19836274635794509466014568838483484624600013041370895881939528370985142278954",
		"13617893609837248081008283892047260108796513789343133982990357780506196299573",
		"17019866949954776883112551581972189814819563148075409968100822283266078590629",
		"2138726945387356468363054067854480562066796241025489545311417198275461384100",
		"9530699369832533880639149369063437013450242560813587771525635668574190709561",
		"3103391909757162857233288457092921903629638236252641536170476897663767868208",
		"5180474051210496608687766129670424206209071725010745029464643986723420860574",
	}
	z, ok := new(big.Int).SetString(want[0], 10)
	if !ok {
		t.Fatal("bad ZERO_VALUE literal")
	}
	for i := 1; i < len(want); i++ {
		z = gnarkMiMCHash(t, z, z)
		if z.String() != want[i] {
			t.Fatalf("zeros(%d): got %s want %s", i, z, want[i])
		}
	}
}

func testAccounts(t *testing.T, n int) []*Account {
	t.Helper()
	keys := make([]*eddsa.PrivateKey, n)
	for i := range keys {
		k, err := eddsa.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("eddsa key: %v", err)
		}
		keys[i] = k
	}
	accounts, err := CreateAccounts(keys)
	if err != nil {
		t.Fatalf("create accounts: %v", err)
	}
	for i, a := range accounts {
		a.Balance = big.NewInt(int64(1000 + i))
	}
	return accounts
}

// TestSolidityMerkleTreeMatchesState proves that the fixed-depth tree built
// with the Solidity conventions (tree leaf = leafSum(accountHash), node =
// hashLeftRight(l, r)) yields the same root as the off-chain gnark-crypto
// state tree for a full power-of-two leaf set, and that the proof paths
// emitted by State.MerkleProofBytes replay to that root under the Solidity
// verify() logic (path[0] = raw account hash, hashed once as the leaf).
func TestSolidityMerkleTreeMatchesState(t *testing.T) {
	const nLeaves = 4 // depth 2

	accounts := testAccounts(t, nLeaves)
	state, err := NewState(mimc.NewMiMC(), accounts)
	if err != nil {
		t.Fatalf("new state: %v", err)
	}

	stateRoot, err := state.Root()
	if err != nil {
		t.Fatalf("state root: %v", err)
	}

	// Solidity-convention tree over the account hashes.
	level := make([]*big.Int, nLeaves)
	for i := 0; i < nLeaves; i++ {
		accountHash := new(big.Int).SetBytes(state.HData[i*32 : (i+1)*32])
		level[i] = solMiMCHash(accountHash) // leafSum
	}
	for len(level) > 1 {
		next := make([]*big.Int, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			next = append(next, solMiMCHash(level[i], level[i+1])) // hashLeftRight
		}
		level = next
	}
	if level[0].Cmp(new(big.Int).SetBytes(stateRoot)) != 0 {
		t.Fatalf("root mismatch:\n solidity tree = %s\n state tree    = %s", level[0], new(big.Int).SetBytes(stateRoot))
	}

	// Replay every proof path through the Solidity verify() logic.
	for idx := uint64(0); idx < nLeaves; idx++ {
		root, path, err := state.MerkleProofBytes(idx)
		if err != nil {
			t.Fatalf("merkle proof %d: %v", idx, err)
		}
		computed := solMiMCHash(path[0]) // leafSum(path[0])
		depth := len(path) - 1
		for i := 1; i <= depth; i++ {
			bit := (idx >> (i - 1)) & 1
			if bit == 1 {
				computed = solMiMCHash(path[i], computed)
			} else {
				computed = solMiMCHash(computed, path[i])
			}
		}
		if computed.Cmp(new(big.Int).SetBytes(root)) != 0 {
			t.Fatalf("verify() replay mismatch for leaf %d:\n computed = %s\n root     = %s", idx, computed, new(big.Int).SetBytes(root))
		}
	}
}
