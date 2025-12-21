package db

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPFS(t *testing.T) {
	t.Parallel()

	const TREE_DEPTH = 10

	t.Run("init simulation ipfs client", func(t *testing.T) {
		t.Parallel()

		ipfsClient, err := NewIPFSClient(true, TREE_DEPTH)

		assert.Equal(t, err, nil)
		assert.NotEqual(t, ipfsClient, nil)
		assert.NotEqual(t, ipfsClient.IPFS, nil)
		assert.NotEqual(t, ipfsClient.IPFS.Data, nil)
	})

	t.Run("get/set in simulation ipfs data", func(t *testing.T) {
		t.Parallel()

		k, v := "sample hash", []byte("sample bytes")
		ipfsClient, _ := NewIPFSClient(true, TREE_DEPTH)
		_, ok := ipfsClient.IPFS.Data[k]
		assert.Equal(t, ok, false)

		ipfsClient.IPFS.Data[k] = v
		fetchedValue, ok := ipfsClient.IPFS.Data[k]
		assert.Equal(t, ok, true)
		assert.Equal(t, fetchedValue, v)
	})

	t.Run("upload/download in simulation ipfs client", func(t *testing.T) {
		t.Parallel()

		k, v := "sample hash", []byte("sample bytes")
		ipfsClient, _ := NewIPFSClient(true, TREE_DEPTH)
		_, ok := ipfsClient.IPFS.Data[k]
		assert.Equal(t, ok, false)

		hash, err := ipfsClient.Upload(v)

		assert.Equal(t, err, nil)
		assert.NotEqual(t, hash, "")

		content, err := ipfsClient.Download(hash)

		assert.Equal(t, err, nil)
		assert.NotEqual(t, content, nil)
		assert.Equal(t, content, v)
	})

	t.Run("generate and save zero values", func(t *testing.T) {
		t.Parallel()

		ipfsClient, _ := NewIPFSClient(true, TREE_DEPTH)
		v, ok := ipfsClient.IPFS.Data[ZERO_VALUES_HASH]
		assert.True(t, ok)
		assert.NotEmpty(t, v)

	})

	t.Run("get zero values", func(t *testing.T) {
		t.Parallel()

		ipfsClient, _ := NewIPFSClient(true, TREE_DEPTH)
		v, err := ipfsClient.GetZeroValues()
		log.Println(v)
		assert.Equal(t, err, nil)
		assert.NotEqual(t, v, nil)
		assert.NotEmpty(t, v)
	})
}
