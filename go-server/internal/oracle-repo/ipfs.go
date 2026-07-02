package oracle

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"l2alchemy/internal/oracle-repo/merkle"
	"math/big"
)

type IPFSContent struct {
	CommitmentHashIncVote map[string]IncVote           `json:"commitment_hash_inc_vote"`
	IncMerkleTree         merkle.IncrementalMerkleTree `json:"inc_merkle_tree"`
	// SpentNullifiers is the committee's nullifier spent-list (paper SIV-E):
	// validators reject a claim whose nullifier hash is already recorded here.
	// Keys are 0x-prefixed 32-byte nullifier hashes.
	SpentNullifiers map[string]bool `json:"spent_nullifiers"`
}

func (c *IPFSContent) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeserializeIPFSContent(data []byte) (*IPFSContent, error) {
	var content IPFSContent
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (n *Node) FetchIPFS() *IPFSContent {
	content, err := n.IPFSClient.Download("")
	if err != nil {
		panic(fmt.Errorf("failed to fetch ipfs content (node_id=%v): %v", err, n.ID))
	}
	if content == nil {
		return &IPFSContent{
			CommitmentHashIncVote: make(map[string]IncVote),
			SpentNullifiers:       make(map[string]bool),
		}
	}
	ipfsContent, err := DeserializeIPFSContent(content)
	if err != nil {
		panic(fmt.Errorf("failed to deserialize ipfs content (node_id=%v): %v", err, n.ID))
	}
	if ipfsContent.SpentNullifiers == nil {
		ipfsContent.SpentNullifiers = make(map[string]bool)
	}

	return ipfsContent
}

// PersistSpentNullifiers publishes the aggregator's spent-list to the DFS so
// validators can consult it when verifying claims (paper SIV-E).
func (n *Node) PersistSpentNullifiers() error {
	if !n.IsAggregator() {
		return fmt.Errorf("node is not an aggregator (node_id=%v)", n.ID)
	}
	content := n.FetchIPFS()
	n.MergeSpentNullifiers(content.SpentNullifiers)
	content.SpentNullifiers = n.SpentNullifiersSnapshot()
	n.IPFSContent = content
	latestHash, err := n.UpdateIPFS()
	if err != nil {
		return err
	}
	// Re-anchor the (unchanged) commitment root with the new DFS reference
	// so users always resolve the pointer the Gateway has recorded (F-22).
	root := new(big.Int)
	if r := content.IncMerkleTree.LatestRoot(); r != nil {
		root.SetBytes(r)
	}
	return n.publishCommitmentRootTx(root, latestHash)
}

func (n *Node) UpdateIPFS() (string, error) {
	if !n.IsAggregator() {
		return "", fmt.Errorf("node is not an aggregator (node_id=%v)", n.ID)
	}
	content, err := n.IPFSContent.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize node ipfs content(node_id=%v): %v", n.ID, err)
	}
	latestHash, err := n.IPFSClient.Upload(content)
	if err != nil {
		return "", fmt.Errorf("failed to update ipfs content (node_id=%v): %v", n.ID, err)
	}

	return latestHash, nil
}
