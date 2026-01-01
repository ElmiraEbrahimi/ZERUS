package gnark

import (
	"bytes"
	"fmt"
	"hash"
	"math/big"
	"sync"

	"l2alchemy/internal/oracle-repo/util"

	"github.com/consensys/gnark-crypto/accumulator/merkletree"
)

const (
	SegmentSize = 32 // Size of each leaf in the Merkle tree
)

type State struct {
	sync.RWMutex
	hFunc hash.Hash
	Data  []byte // Serialized account data
	HData []byte // Hashed account data
}

// Initializes a new State object with serialized accounts
func NewState(hFunc hash.Hash, accounts []*Account) (*State, error) {
	data := make([]byte, AccountSize*len(accounts))
	hData := make([]byte, hFunc.Size()*len(accounts))

	for i, account := range accounts {
		hFunc.Reset()

		accountData := account.Serialize()
		s, err := hashAccountForMerkle(hFunc, account)
		if err != nil {
			return nil, fmt.Errorf("hash account: %w", err)
		}

		copy(data[i*AccountSize:(i+1)*AccountSize], accountData)
		copy(hData[i*hFunc.Size():(i+1)*hFunc.Size()], s)
	}

	return &State{
		hFunc: hFunc,
		Data:  data,
		HData: hData,
	}, nil
}

// Update the state with a new account and update the corresponding hash.
func (s *State) WriteAccount(account Account) error {
	s.Lock()
	defer s.Unlock()

	i := int(account.Index.Int64())
	accountData := account.Serialize()
	copy(s.Data[i*AccountSize:], accountData)

	sHash, err := hashAccountForMerkle(s.hFunc, &account)
	if err != nil {
		return fmt.Errorf("hash account: %w", err)
	}
	copy(s.HData[i*s.hFunc.Size():(i+1)*s.hFunc.Size()], sHash)

	return nil
}

// retrieve the account data at the specified index.
func (s *State) ReadAccount(i uint64) (Account, error) {
	s.RLock()
	defer s.RUnlock()

	var res Account
	res.Deserialize(s.Data[int(i)*AccountSize : int(i)*AccountSize+AccountSize])
	return res, nil
}

// Compute the Merkle root of the current state using the `ReaderRoot` function.
func (s *State) Root() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	var stateBuf bytes.Buffer
	_, err := stateBuf.Write(s.HData)
	if err != nil {
		return nil, fmt.Errorf("failed to write state data to buffer: %v", err)
	}

	root, err := merkletree.ReaderRoot(&stateBuf, s.hFunc, SegmentSize)
	if err != nil {
		return nil, fmt.Errorf("failed to compute Merkle root: %w", err)
	}

	return root, nil
}

// MerkleProof generates a Merkle proof for an account at the given index using `BuildReaderProof`.
func (s *State) MerkleProof(i uint64) ([]byte, [][]byte, error) {
	s.RLock()
	defer s.RUnlock()

	var path [][]byte

	// Create a buffer for the hashed data
	var stateBuf bytes.Buffer
	_, err := stateBuf.Write(s.HData)
	if err != nil {
		return nil, path, fmt.Errorf("failed to write state data: %v", err)
	}

	// Use BuildReaderProof to compute the Merkle root and proof set
	root, proofSet, _, err := merkletree.BuildReaderProof(&stateBuf, s.hFunc, SegmentSize, i)
	if err != nil {
		return nil, path, fmt.Errorf("failed to build Merkle proof: %v", err)
	}
	p := make([]*big.Int, len(proofSet))
	for i, node := range proofSet {
		p[i] = big.NewInt(0).SetBytes(node)
	}
	for i := 0; i < len(proofSet); i++ {
		path = append(path, proofSet[i])
	}

	return root, path, nil
}

// MerkleProofBytes generates a Merkle proof for an account and returns the proof as `*big.Int` slices.
func (s *State) MerkleProofBytes(i uint64) ([]byte, []*big.Int, error) {
	s.RLock()
	defer s.RUnlock()
	path := make([]*big.Int, 0)
	var stateBuf bytes.Buffer
	_, err := stateBuf.Write(s.HData)
	if err != nil {
		return nil, path, fmt.Errorf("failed to write state data: %v", err)
	}
	root, proofSet, _, err := merkletree.BuildReaderProof(&stateBuf, s.hFunc, SegmentSize, i)
	if err != nil {
		return nil, path, fmt.Errorf("failed to build Merkle proof: %v", err)
	}
	path = make([]*big.Int, len(proofSet))
	for i := 0; i < len(proofSet); i++ {
		path[i] = big.NewInt(0).SetBytes(proofSet[i])
	}
	return root, path, nil
}

// SetData updates the serialized account data.
func (s *State) SetData(data []byte) {
	s.Lock()
	defer s.Unlock()
	s.Data = data
}

// SetHData updates the hashed account data.
func (s *State) SetHData(hData []byte) {
	s.Lock()
	defer s.Unlock()
	s.HData = hData
}

func hashAccountForMerkle(hFunc hash.Hash, account *Account) ([]byte, error) {
	if account == nil || account.PublicKey == nil {
		return nil, fmt.Errorf("account or public key is nil")
	}
	hFunc.Reset()
	if _, err := hFunc.Write(util.PadOrTrim(account.Index.Bytes(), SegmentSize)); err != nil {
		return nil, err
	}
	pkX := account.PublicKey.A.X.Bytes()
	pkY := account.PublicKey.A.Y.Bytes()
	if _, err := hFunc.Write(pkX[:]); err != nil {
		return nil, err
	}
	if _, err := hFunc.Write(pkY[:]); err != nil {
		return nil, err
	}
	if _, err := hFunc.Write(util.PadOrTrim(account.Balance.Bytes(), SegmentSize)); err != nil {
		return nil, err
	}
	return hFunc.Sum(nil), nil
}
