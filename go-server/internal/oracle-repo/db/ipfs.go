package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"l2alchemy/internal/oracle-repo/merkle"
	"os"
	"path/filepath"
	"sync"

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
	storagePath string
	mu          sync.RWMutex
}

type ipfsSimState struct {
	Data       map[string][]byte
	LatestHash string
}

func NewIPFSClient(simulate bool, zeroValueTreeDepth int, storagePath string) (*IPFSClient, error) {
	if simulate {
		client := &IPFSClient{
			IPFS: &IPFS{
				Data: make(map[string][]byte),
			},
			IsSimulation: true,
			storagePath:  storagePath,
		}

		if storagePath != "" {
			if err := client.loadState(); err != nil {
				return nil, err
			}
		}
		if _, ok := client.IPFS.Data[ZERO_VALUES_HASH]; !ok {
			if err := client.saveZeroValues(zeroValueTreeDepth); err != nil {
				return nil, err
			}
		}

		return client, nil
	}
	return nil, errors.New("not implemented")
}

// region core

func (i *IPFSClient) Upload(content []byte) (string, error) {
	if i.IsSimulation {
		i.mu.Lock()
		defer i.mu.Unlock()
		hash := generateRandomHash()
		i.IPFS.Data[hash] = content
		i.LatestHash = hash
		if err := i.saveState(); err != nil {
			return "", err
		}
		return hash, nil
	}
	return "", errors.New("not implemented")
}

func (i *IPFSClient) Download(hash string) ([]byte, error) {
	if i.IsSimulation {
		i.mu.RLock()
		defer i.mu.RUnlock()
		if hash == DEFAULT_HASH {
			hash = i.LatestHash
		}
		if hash == "" {
			return nil, errors.New("no latest hash available in simulated IPFS")
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
	i.mu.RLock()
	defer i.mu.RUnlock()
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
	i.mu.Lock()
	defer i.mu.Unlock()
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
	return i.saveState()
}

func (i *IPFSClient) loadState() error {
	if i.storagePath == "" {
		return nil
	}
	info, err := os.Stat(i.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return errors.New("ipfs storage path is a directory")
	}

	file, err := os.Open(i.storagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	var state ipfsSimState
	if err := dec.Decode(&state); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	if state.Data == nil {
		state.Data = make(map[string][]byte)
	}

	i.IPFS.Data = state.Data
	i.LatestHash = state.LatestHash
	return nil
}

func (i *IPFSClient) saveState() error {
	if i.storagePath == "" {
		return nil
	}
	dir := filepath.Dir(i.storagePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := os.Create(i.storagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	state := ipfsSimState{
		Data:       i.IPFS.Data,
		LatestHash: i.LatestHash,
	}
	if err := enc.Encode(state); err != nil {
		return err
	}

	return nil
}

// endregion

func generateRandomHash() string {
	id := uuid.New()
	return id.String()
}
