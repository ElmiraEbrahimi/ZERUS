package util

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func GenerateKey() (*eddsa.PrivateKey, error) {
	privateKey := &eddsa.PrivateKey{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	privateKey, err := eddsa.GenerateKey(r)
	if err != nil {
		return nil, fmt.Errorf("eddsa generate key: %w", err)
	}

	return privateKey, nil
}

func GenerateKeys(count int) ([]*eddsa.PrivateKey, error) {
	if count <= 0 {
		return nil, errors.New("count must be above zero")
	}
	privateKeys := make([]*eddsa.PrivateKey, count)
	for i := 0; i < count; i++ {
		sk, err := GenerateKey()
		if err != nil {
			panic(err)
		}
		privateKeys[i] = sk
	}
	return privateKeys, nil
}
