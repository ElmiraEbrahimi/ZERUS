package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"l2alchemy/internal/oracle-repo/merkle"

	"github.com/google/uuid"
)

const (
	DEFAULT_HASH     = ""
	ZERO_VALUES_HASH = "zv"
)

type IPFS struct {
	Data map[string][]byte
}

type IPFSClient struct {
	IPFS         *IPFS
	IsSimulation bool   `json:"is_simulation"`
	LatestHash   string `json:"latest_hash"`
}

func NewIPFSClient(simulate bool, zeroValueTreeDepth int) (*IPFSClient, error) {
	if simulate {
		data := make(map[string][]byte)
		client := &IPFSClient{
			IPFS: &IPFS{
				Data: data,
			},
			IsSimulation: true,
		}

		client.saveZeroValues(zeroValueTreeDepth)

		return client, nil
	}
	return nil, errors.New("not implemented")
}

// region core

func (i *IPFSClient) Upload(content []byte) (string, error) {
	if i.IsSimulation {
		hash := generateRandomHash()
		i.IPFS.Data[hash] = content
		i.LatestHash = hash
		return hash, nil
	}
	return "", errors.New("not implemented")
}

func (i *IPFSClient) Download(hash string) ([]byte, error) {
	if i.IsSimulation {
		if hash == DEFAULT_HASH {
			hash = i.LatestHash
		}
		content, ok := i.IPFS.Data[hash]
		if !ok {
			return nil, errors.New("hash does not exist in IPFS")
		}
		return content, nil
	}
	return nil, errors.New("not implemented")
}

// endregion

// region merkle tree

func (i *IPFSClient) GetZeroValues() ([][]byte, error) {
	var zeroValues [][]byte
	content, err := i.Download(ZERO_VALUES_HASH)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(content)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&zeroValues)
	if err != nil {
		return nil, err
	}

	return zeroValues, nil
}

func (i *IPFSClient) saveZeroValues(depth int) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	zeroValues, err := merkle.GenerateZeroValues(depth)
	if err != nil {
		return err
	}

	err = enc.Encode(zeroValues)
	if err != nil {
		return err
	}
	content := buf.Bytes()
	i.IPFS.Data[ZERO_VALUES_HASH] = content

	return nil
}

// endregion

func generateRandomHash() string {
	id := uuid.New()
	return id.String()
}
