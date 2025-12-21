package oracle

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"l2alchemy/internal/oracle-repo/merkle"
)

type IPFSContent struct {
	CommitmentHashIncVote map[string]IncVote           `json:"commitment_hash_inc_vote"`
	IncMerkleTree         merkle.IncrementalMerkleTree `json:"inc_merkle_tree"`
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
		}
	}
	ipfsContent, err := DeserializeIPFSContent(content)
	if err != nil {
		panic(fmt.Errorf("failed to deserialize ipfs content (node_id=%v): %v", err, n.ID))
	}

	return ipfsContent
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
