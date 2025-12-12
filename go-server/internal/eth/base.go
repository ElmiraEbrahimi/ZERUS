package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainClient struct {
	Eth     *ethclient.Client
	ChainID *big.Int
	PrivKey *ecdsa.PrivateKey
	From    common.Address
}

func NewChainClient(ctx context.Context, rpcURL string, chainID int64, privKeyHex string) (*ChainClient, error) {
	// strip optional 0x
	if len(privKeyHex) >= 2 && privKeyHex[:2] == "0x" {
		privKeyHex = privKeyHex[2:]
	}

	b, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}

	privKey, err := crypto.ToECDSA(b)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	cli, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	from := crypto.PubkeyToAddress(privKey.PublicKey)

	return &ChainClient{
		Eth:     cli,
		ChainID: big.NewInt(chainID),
		PrivKey: privKey,
		From:    from,
	}, nil
}

// Shared helper; useful for health checks or generic endpoints.
func (c *ChainClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.Eth.BlockNumber(ctx)
}
