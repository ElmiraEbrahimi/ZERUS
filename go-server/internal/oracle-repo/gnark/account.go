package gnark

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"

	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/oracle-repo/util"
)

const AccountSize = 192 // 32 bytes × 6 fields

type Account struct {
	Index         *big.Int
	PublicKey     *eddsa.PublicKey
	Balance       *big.Int
	Reputation    *big.Int
	SeverityCount *big.Int
}

// Serialize Index ∥ pubkeyX ∥ pubkeyY ∥ balance
func (a *Account) Serialize() []byte {
	var b [AccountSize]byte

	copy(b[:32], util.PadOrTrim(a.Index.Bytes(), 32))

	var buf [32]byte
	buf = a.PublicKey.A.X.Bytes()
	copy(b[32:], buf[:])
	buf = a.PublicKey.A.Y.Bytes()
	copy(b[64:], buf[:])

	copy(b[96:], util.PadOrTrim(a.Balance.Bytes(), 32))
	copy(b[128:], util.PadOrTrim(a.Reputation.Bytes(), 32))
	copy(b[160:], util.PadOrTrim(a.SeverityCount.Bytes(), 32))

	return b[:]
}

func (a *Account) Deserialize(data []byte) {

	a.Index = big.NewInt(0).SetBytes(data[:32])

	a.PublicKey = new(eddsa.PublicKey)

	a.PublicKey.A.X.SetZero()
	a.PublicKey.A.Y.SetOne()

	a.PublicKey.A.X.SetBytes(data[32:64])
	a.PublicKey.A.Y.SetBytes(data[64:96])

	a.Balance = big.NewInt(0).SetBytes(data[96:128])
	a.Reputation = big.NewInt(0).SetBytes(data[128:160])
	a.SeverityCount = big.NewInt(0).SetBytes(data[160:AccountSize])
}

func PublicKeyToOraclePublicKey(publicKey *eddsa.PublicKey) *bc.OraclePublicKey {
	return &bc.OraclePublicKey{
		X: publicKey.A.X.BigInt(big.NewInt(0)),
		Y: publicKey.A.Y.BigInt(big.NewInt(0)),
	}
}

func AccountToOracleAccount(account *Account) *bc.OracleAccount {
	return &bc.OracleAccount{
		Index:         util.ModToBn254(account.Index),
		PubKey:        *PublicKeyToOraclePublicKey(account.PublicKey),
		Balance:       util.ModToBn254(account.Balance),
		Reputation:    util.ModToBn254(account.Reputation),
		SeverityCount: util.ModToBn254(account.SeverityCount),
	}
}

func CreateAccounts(privateKeys []*eddsa.PrivateKey) ([]*Account, error) {
	accounts := make([]*Account, len(privateKeys))
	for i, privateKey := range privateKeys {
		accounts[i] = &Account{
			Index:         big.NewInt(int64(i)),
			PublicKey:     &privateKey.PublicKey,
			Balance:       big.NewInt(0),
			Reputation:    big.NewInt(50),
			SeverityCount: big.NewInt(0),
		}
	}
	return accounts, nil
}

func CreateAccount(privateKey *eddsa.PrivateKey, index int64) (*Account, error) {
	account := &Account{
		Index:         big.NewInt(index),
		PublicKey:     &privateKey.PublicKey,
		Balance:       big.NewInt(0),
		Reputation:    big.NewInt(50),
		SeverityCount: big.NewInt(0),
	}

	return account, nil
}
