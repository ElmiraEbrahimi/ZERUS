package oracle

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	votingbatch "l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/config"
	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/memtime"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/merkle"
	"l2alchemy/internal/oracle-repo/util"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	edwards "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Validator = iota
	Aggregator
)

const validatorStateSyncTimeout = 5 * time.Minute

var (
	proofInvalidSelector          = crypto.Keccak256([]byte("ProofInvalid()"))[:4]
	publicInputNotInFieldSelector = crypto.Keccak256([]byte("PublicInputNotInField()"))[:4]
)

type Node struct {
	cfg *config.Config
	wg  *sync.WaitGroup

	ID              uint `json:"id"`
	ecdsaPrivateKey *ecdsa.PrivateKey
	privateKey      *eddsa.PrivateKey
	Account         *gnark.Account
	state           *gnark.State

	Oracle         *Oracle `json:"oracle"`
	OracleMessages chan InternalOraclelMessage
	Role           uint `json:"role"`
	IsRunning      bool `json:"is_running"`

	ethClient *ethclient.Client

	IPFSClient       *db.IPFSClient
	IPFSContent      *IPFSContent
	IncMerkleTree    *merkle.IncrementalMerkleTree
	SparseMerkleTree *merkle.SparseMerkleTree

	incVotesLock sync.Mutex
	wiVotesLock  sync.Mutex
	IncVotes     map[string]map[uint]*IncVote `json:"inc_votes"`
	WiVotes      map[string]map[uint]*WiVote  `json:"wi_votes"`

	// Local view of the committee's nullifier spent-list (paper SIV-E,
	// Theorem 2). The authoritative copy lives in the DFS; this cache is
	// merged from fetched DFS content and from ClaimMinted events.
	spentNullifiersLock sync.Mutex
	spentNullifiers     map[string]bool

	// The node's own per-request claim votes (request ID -> 0/1), recorded
	// when the node verifies a claim. Used to build and sign the node's own
	// batch bitmask (paper SIV-D Step 14); only this node ever signs it.
	myVotesLock sync.Mutex
	myVotes     map[string]*big.Int

	BatchedWiVote *BatchedWiVote
}

// MarkNullifierSpent records a nullifier hash (0x-prefixed hex) in the node's
// local view of the committee spent-list.
func (n *Node) MarkNullifierSpent(nullifierHash string) {
	n.spentNullifiersLock.Lock()
	defer n.spentNullifiersLock.Unlock()
	if n.spentNullifiers == nil {
		n.spentNullifiers = make(map[string]bool)
	}
	n.spentNullifiers[nullifierHash] = true
}

// IsNullifierSpent reports whether the nullifier hash is in the node's local
// view of the committee spent-list.
func (n *Node) IsNullifierSpent(nullifierHash string) bool {
	n.spentNullifiersLock.Lock()
	defer n.spentNullifiersLock.Unlock()
	return n.spentNullifiers[nullifierHash]
}

// MergeSpentNullifiers folds a spent-list fetched from the DFS into the
// node's local view.
func (n *Node) MergeSpentNullifiers(spent map[string]bool) {
	if len(spent) == 0 {
		return
	}
	n.spentNullifiersLock.Lock()
	defer n.spentNullifiersLock.Unlock()
	if n.spentNullifiers == nil {
		n.spentNullifiers = make(map[string]bool)
	}
	for k, v := range spent {
		if v {
			n.spentNullifiers[k] = true
		}
	}
}

// SpentNullifiersSnapshot returns a copy of the node's local spent-list.
func (n *Node) SpentNullifiersSnapshot() map[string]bool {
	n.spentNullifiersLock.Lock()
	defer n.spentNullifiersLock.Unlock()
	out := make(map[string]bool, len(n.spentNullifiers))
	for k, v := range n.spentNullifiers {
		if v {
			out[k] = true
		}
	}
	return out
}

func NewNode(cfg *config.Config, ethClient *ethclient.Client, ipfsClient *db.IPFSClient, oracle *Oracle, id uint, privateKey *eddsa.PrivateKey, offeredAmount uint, role uint, validatorAccounts []*gnark.Account) *Node {
	n := &Node{
		cfg:        cfg,
		ethClient:  ethClient,
		Oracle:     oracle,
		IPFSClient: ipfsClient,
		ID:         id,
		Role:       Validator,
		wg:         &sync.WaitGroup{},
	}

	fmt.Printf("setting account for node (=%v)...\n", n.ID)
	n.privateKey = privateKey
	for i := range validatorAccounts {
		if validatorAccounts[i].Index.Uint64() == uint64(n.ID) {
			n.Account = validatorAccounts[i]
			break
		}
	}

	fmt.Printf("setting state for node (=%v)...\n", n.ID)
	state, err := gnark.NewState(hash.MIMC_BN254.New(), validatorAccounts)
	if err != nil {
		panic(err)
	}
	n.state = state

	pk := strings.TrimPrefix(n.cfg.NodePK, "0x")
	pk = strings.TrimPrefix(pk, "0X")
	ecdsaPrivateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatalf("failed to parseecdsaPrivateKey: %v", err)
	}
	n.ecdsaPrivateKey = ecdsaPrivateKey

	n.OracleMessages = make(chan InternalOraclelMessage, 100)
	n.IncVotes = make(map[string]map[uint]*IncVote)
	n.WiVotes = make(map[string]map[uint]*WiVote)

	zeroValues, err := n.IPFSClient.GetZeroValues()
	if err != nil {
		panic(fmt.Errorf("failed to get zero values from ipfs: %v", err))
	}
	n.IncMerkleTree, err = merkle.NewIncrementalMerkleTree(zeroValues, cfg.IncTreeDepth)
	if err != nil {
		panic(fmt.Errorf("failed to create inc merkle tree: %v", err))
	}
	n.SparseMerkleTree, err = merkle.NewSparseMerkleTree(cfg.SparseTreeDepth, zeroValues)
	if err != nil {
		panic(fmt.Errorf("failed to create sparse merkle tree: %v", err))
	}
	n.IPFSContent = &IPFSContent{
		CommitmentHashIncVote: make(map[string]IncVote),
	}

	n.BatchedWiVote = &BatchedWiVote{
		WithdrawalReqIDs: make([]*big.Int, 0),
		Vote:             make([]*big.Int, 0),
		MajorityOfVotes:  make([]*big.Int, 0),
	}

	return n
}

func (n *Node) IsAggregator() bool {
	return n.Role == Aggregator
}

func (n *Node) Start() {
	if n.wg == nil {
		n.wg = &sync.WaitGroup{}
	}

	n.wg.Add(1)
	go func() {
		for e := range n.OracleMessages {
			// ignore messages from self:
			// if e.From == n.ID {
			// 	continue
			// }

			// message from oracle:
			if e.From == OracleID {
				if e.Title == MessageTitleTerminate {
					fmt.Printf("node %d terminating...\n", n.ID)
					/////////////////////////////////////////////////////
					// if n.ID == 0 { // or if n.IsAggregator() or totalValidators-1
					// //  or with setting in config n.ShouldWithdrawOnExit

					// 	err := n.WithDrawAccounts([]*gnark.Account{n.Account}, n.state)
					// 	if err != nil {
					// 		fmt.Printf("node %d failed to withdraw: %v\n", n.ID, err)
					// 	} else {
					// 		fmt.Printf("node %d successfully withdrew before termination.\n", n.ID)
					// 	}
					// }
					// /////////////////////////////////////////////////////
					n.wg.Done()
					return
				}
				if e.Title == MessageTittleSelectWiVote && n.IsAggregator() {
					fmt.Printf("aggregator (node=%d) is performing wiVote selection...\n", n.ID)
					fmt.Printf("wivotes: %v\n", n.WiVotes)
					if len(n.WiVotes) != 0 {
						// Process claim requests in ascending identifier
						// order so batched votes line up with the
						// deterministic round windows of SIV-D.
						reqIDs := make([]*big.Int, 0, len(n.WiVotes))
						for uniqueReqID := range n.WiVotes {
							id, ok := new(big.Int).SetString(uniqueReqID, 10)
							if !ok {
								// Skip malformed ids instead of killing the
								// committee process (F-28).
								fmt.Printf("node %d: invalid wiVote request id %q, skipping\n", n.ID, uniqueReqID)
								continue
							}
							reqIDs = append(reqIDs, id)
						}
						sort.Slice(reqIDs, func(i, j int) bool {
							return reqIDs[i].Cmp(reqIDs[j]) < 0
						})
						for _, id := range reqIDs {
							// aggregator wiVote process:
							err := n.AggregatorProcessWiVote(id.String())
							if err != nil {
								// Log and continue; the request stays queued
								// for a later round (F-28).
								fmt.Printf("node %d: failed to process wiVote %s: %v\n", n.ID, id, err)
								continue
							}
						}
					}

					// reset wivotes:
					n.aggregatorResetWiVotes()

					if n.batchedWiVoteCount() >= n.cfg.BatchSize {
						n.Oracle.PublishBatchedWiVoteRes()
					}
				}
				//$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
				if e.Title == MessageTittleBatchedWiVote && n.IsAggregator() {
					fmt.Printf("aggregator (node=%d) is performing batch wiVote... (len of batched votes=%d)\n", n.ID, len(n.BatchedWiVote.WithdrawalReqIDs))

					if n.batchedWiVoteCount() < n.cfg.BatchSize {
						fmt.Printf("aggregator (node=%d) skipping batch wiVote: need %d requests, have %d\n", n.ID, n.cfg.BatchSize, len(n.BatchedWiVote.WithdrawalReqIDs))
						continue
					}
					n.processPendingBatchedWiVotes()
				}
				//$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
				if e.Title == MessageTittleSelectIncVote && n.IsAggregator() {
					fmt.Printf("aggregator (node %d) is performing incVote selection...\n", n.ID)
					fmt.Printf("incvotes: %v\n", n.IncVotes)
					if len(n.IncVotes) < n.cfg.BatchSize {
						fmt.Printf("aggregator (node=%d) skipping incVote batch: need %d requests, have %d\n", n.ID, n.cfg.BatchSize, len(n.IncVotes))
						continue
					}

					// fetch and update ipfs:
					if n.IPFSClient.LatestHash != "" {
						fmt.Printf("fetching ipfs content... (node_id=%v)\n", n.ID)
						n.IPFSContent = n.FetchIPFS()
					}

					for len(n.IncVotes) >= n.cfg.BatchSize {
						selected, keys, err := n.selectIncVoteBatch(n.cfg.BatchSize)
						if err != nil {
							fmt.Printf("failed to select incVote batch: %v\n", err)
							break
						}
						if len(selected) != n.cfg.BatchSize {
							fmt.Printf("aggregator (node=%d) waiting for full incVote batch: need %d requests, have %d ready\n", n.ID, n.cfg.BatchSize, len(selected))
							break
						}
						for _, incVote := range selected {
							fmt.Printf("selected tree index is: %v\n", incVote.IncTreeIndex)
							n.IPFSContent.CommitmentHashIncVote[string(incVote.CommitmentHash)] = *incVote
						}
						n.IPFSContent.IncMerkleTree = *n.IncMerkleTree
						fmt.Printf("updating ipfs content... (node_id=%v)\n", n.ID)
						latestIPFSHash, err := n.UpdateIPFS()
						if err != nil {
							log.Printf("failed to update IPFS content (node=%d): %v", n.ID, err)
							break
						}
						// Anchor the finalized commitment root with the DFS
						// reference on-chain (paper SIV-E Steps 6-7).
						root := new(big.Int)
						if r := n.IncMerkleTree.LatestRoot(); r != nil {
							root.SetBytes(r)
						}
						if err := n.publishCommitmentRootTx(root, latestIPFSHash, "L2 publish IPFS burn root"); err != nil {
							log.Printf("failed to publish IPFS burn root (node=%d): %v", n.ID, err)
							break
						}
						for _, key := range keys {
							delete(n.IncVotes, key)
						}
					}
				}
			}

			// message from other nodes:
			if e.Title == MessageTittleWiVote && n.IsAggregator() {
				fmt.Printf("aggregator (node %d) is collecting wiVote...\n", n.ID)
				wiVote := e.Message.(*WiVote)
				n.aggregatorCollectWiVote(wiVote)
			}
			if e.Title == MessageTittleIncVote && n.IsAggregator() {
				fmt.Printf("aggregator (node %d) is collecting incVote...\n", n.ID)
				incVote := e.Message.(*IncVote)
				n.aggregatorCollectIncVote(incVote)
			}
		}
	}()
}

func (n *Node) VerifyClaim(claimEvent *bc.OracleClaimSubmitted) (*WiVote, error) {
	proof := groth16.NewProof(ecc.BN254)
	proofReader := bytes.NewReader(claimEvent.Proof)
	if _, err := proof.ReadFrom(proofReader); err != nil {
		return nil, fmt.Errorf("failed to read proof (node=%v): %v", n.ID, err)
	}

	publicWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize public witness (node=%v): %v", n.ID, err)
	}
	err = publicWitness.UnmarshalBinary(claimEvent.PublicWitness)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public witness (node=%v): %v", n.ID, err)
	}

	// Consult the committee's DFS spent-list before verifying the Redeeming
	// proof (paper SIV-E): a claim whose nullifier hash is already spent is
	// rejected regardless of proof validity (Theorem 2).
	nullifierKey := common.BytesToHash(claimEvent.NullifierHash[:]).Hex()
	if n.IPFSClient != nil && n.IPFSClient.LatestHash != "" {
		if content := n.FetchIPFS(); content != nil {
			n.MergeSpentNullifiers(content.SpentNullifiers)
		}
	}

	// Decode the exact public layout: d_dst, R_comm, h_n, a_dst.
	publicInputs, rootErr := redeemingPublicInputs(publicWitness)
	rootPublished := false
	if rootErr == nil {
		rootPublished, rootErr = n.isPublishedCommitmentRoot(publicInputs.CommitmentRoot)
	}
	configuredDestination, destinationOK := new(big.Int).SetString(n.cfg.DestinationID, 10)

	var isApproved *big.Int
	if n.IsNullifierSpent(nullifierKey) {
		fmt.Printf("rejecting claim: nullifier already spent (node=%v nullifier=%s)\n", n.ID, nullifierKey)
		isApproved = big.NewInt(0)
	} else if rootErr != nil {
		fmt.Printf("rejecting claim: cannot validate commitment root (node=%v): %v\n", n.ID, rootErr)
		isApproved = big.NewInt(0)
	} else if !rootPublished {
		fmt.Printf("rejecting claim: commitment root not recorded by the gateway (node=%v root=%s)\n", n.ID, publicInputs.CommitmentRoot)
		isApproved = big.NewInt(0)
	} else if !destinationOK || publicInputs.DestinationID.Cmp(configuredDestination) != 0 || publicInputs.DestinationID.Cmp(claimEvent.DestinationID) != 0 {
		fmt.Printf("rejecting claim: destination mismatch (node=%v)\n", n.ID)
		isApproved = big.NewInt(0)
	} else if publicInputs.NullifierHash.Cmp(new(big.Int).SetBytes(claimEvent.NullifierHash[:])) != 0 {
		fmt.Printf("rejecting claim: nullifier public input mismatch (node=%v)\n", n.ID)
		isApproved = big.NewInt(0)
	} else if publicInputs.RecipientAddress.Cmp(new(big.Int).SetBytes(claimEvent.Recipient.Bytes())) != 0 {
		fmt.Printf("rejecting claim: recipient public input mismatch (node=%v)\n", n.ID)
		isApproved = big.NewInt(0)
	} else if err := groth16.Verify(proof, n.Oracle.IncVK, publicWitness); err != nil {
		fmt.Printf("failed to verify claim proof (node=%v): %v\n", n.ID, err)
		isApproved = big.NewInt(0)
	} else {
		fmt.Printf("successfully verified claim proof (node=%v)\n", n.ID)
		isApproved = big.NewInt(1)
	}

	// Keep the node's own verdict so it can sign its own batch bitmask
	// later (paper SIV-D Step 14; F-11).
	n.recordOwnVote(claimEvent.UniqueID.String(), isApproved)

	wivote := WiVote{
		Index:      n.Account.Index,
		RequestID:  claimEvent.UniqueID,
		IsApproved: isApproved,
	}

	// Sign the hashed message
	msg := hashWiVoteFieldwise(&wivote)
	hfunc := hash.MIMC_BN254.New()
	hfunc.Reset()
	sigBytes, err := n.privateKey.Sign(msg, hfunc)
	if err != nil {
		return nil, fmt.Errorf("sign the vote: %w", err)
	}

	wivote = WiVote{
		Index:      n.Account.Index,
		RequestID:  claimEvent.UniqueID,
		IsApproved: isApproved,
		SenderPK:   n.privateKey.PublicKey,
		Signature:  sigBytes,
	}

	return &wivote, nil
}

// recordOwnVote stores this node's verdict for a claim request so the node
// can later build and sign its own batch bitmask (paper SIV-D Step 14).
func (n *Node) recordOwnVote(requestID string, isApproved *big.Int) {
	n.myVotesLock.Lock()
	defer n.myVotesLock.Unlock()
	if n.myVotes == nil {
		n.myVotes = make(map[string]*big.Int)
	}
	n.myVotes[requestID] = new(big.Int).Set(isApproved)
}

// SignBatchVote builds this node's b-bit vote bitmask for the round window
// batchIDs (ascending) from its own recorded per-request votes and signs
// H(index, C_batch, vote, round) with its own key (paper SIV-D Step 14,
// Alg. 2 lines 12-13). The aggregator must consume the artifact unmodified.
func (n *Node) SignBatchVote(batchIDs []*big.Int, roundID *big.Int, batchCommitment []byte) (*SignedBatchVote, error) {
	mask := big.NewInt(0)
	n.myVotesLock.Lock()
	for pos, id := range batchIDs {
		vote, ok := n.myVotes[id.String()]
		if !ok {
			n.myVotesLock.Unlock()
			return nil, fmt.Errorf("node %d has no recorded vote for request %s", n.ID, id.String())
		}
		if vote != nil && vote.Sign() != 0 {
			mask.SetBit(mask, pos, 1)
		}
	}
	n.myVotesLock.Unlock()

	msg := hashBatchVoteFieldwise(
		new(big.Int).Set(n.Account.Index),
		new(big.Int).SetBytes(batchCommitment[:32]),
		new(big.Int).Set(mask),
		new(big.Int).Set(roundID),
	)
	hfunc := hash.MIMC_BN254.New()
	sig, err := n.privateKey.Sign(msg, hfunc)
	if err != nil {
		return nil, fmt.Errorf("sign batch vote (node=%d): %w", n.ID, err)
	}
	return &SignedBatchVote{
		Index:     new(big.Int).Set(n.Account.Index),
		Vote:      mask,
		Signature: sig,
	}, nil
}

type redeemingInputs struct {
	DestinationID    *big.Int
	CommitmentRoot   *big.Int
	NullifierHash    *big.Int
	RecipientAddress *big.Int
}

// redeemingPublicInputs extracts Algorithm 1's exact public-input order:
// d_dst, R_comm, h_n, a_dst.
func redeemingPublicInputs(w witness.Witness) (*redeemingInputs, error) {
	vec, ok := w.Vector().(fr.Vector)
	if !ok || len(vec) != 4 {
		return nil, fmt.Errorf("unexpected public witness layout")
	}
	return &redeemingInputs{
		DestinationID:    vec[0].BigInt(new(big.Int)),
		CommitmentRoot:   vec[1].BigInt(new(big.Int)),
		NullifierHash:    vec[2].BigInt(new(big.Int)),
		RecipientAddress: vec[3].BigInt(new(big.Int)),
	}, nil
}

// isPublishedCommitmentRoot checks the Gateway's record of commitment roots
// (paper SIV-E Steps 6-7): claims are only accepted against recorded roots.
func (n *Node) isPublishedCommitmentRoot(root *big.Int) (bool, error) {
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		return false, fmt.Errorf("create contract client instance: %w", err)
	}
	callOpts := &bind.CallOpts{Context: context.Background()}
	return bcClient.IsPublishedCommitmentRoot(callOpts, root)
}

// endregion

// region aggregator

func (n *Node) aggregatorResetIncVotes() {
	n.IncVotes = make(map[string]map[uint]*IncVote)
}

func (n *Node) aggregatorResetWiVotes() {
	n.WiVotes = make(map[string]map[uint]*WiVote)
}

func (n *Node) nodeByAccountIndex(index *big.Int) (*Node, error) {
	if n == nil || n.Oracle == nil {
		return nil, fmt.Errorf("oracle not initialized")
	}
	if index == nil {
		return nil, fmt.Errorf("account index is nil")
	}
	for _, node := range n.Oracle.Nodes {
		if node == nil || node.Account == nil || node.Account.Index == nil {
			continue
		}
		if node.Account.Index.Cmp(index) == 0 {
			return node, nil
		}
	}
	return nil, fmt.Errorf("node not found for account index %s", index.String())
}

func (n *Node) resetBatchedWiVotes() {
	for _, node := range n.Oracle.Nodes {
		node.BatchedWiVote.WithdrawalReqIDs = make([]*big.Int, 0)
		node.BatchedWiVote.Vote = make([]*big.Int, 0)
		node.BatchedWiVote.MajorityOfVotes = make([]*big.Int, 0)
	}
}

func (n *Node) batchedWiVoteCount() int {
	if n == nil || n.BatchedWiVote == nil {
		return 0
	}
	return len(n.BatchedWiVote.WithdrawalReqIDs)
}

func (n *Node) peekBatchedWiVoteIDs(batchSize int) []*big.Int {
	if n == nil || n.BatchedWiVote == nil {
		return nil
	}
	if len(n.BatchedWiVote.WithdrawalReqIDs) < batchSize {
		return nil
	}
	batch := make([]*big.Int, batchSize)
	copy(batch, n.BatchedWiVote.WithdrawalReqIDs[:batchSize])
	return batch
}

func (n *Node) consumeBatchedWiVotes(batchSize int) {
	if n == nil || n.Oracle == nil {
		return
	}
	for _, node := range n.Oracle.Nodes {
		if node == nil || node.BatchedWiVote == nil {
			continue
		}
		if len(node.BatchedWiVote.WithdrawalReqIDs) >= batchSize {
			node.BatchedWiVote.WithdrawalReqIDs = node.BatchedWiVote.WithdrawalReqIDs[batchSize:]
		} else {
			node.BatchedWiVote.WithdrawalReqIDs = make([]*big.Int, 0)
		}
		if len(node.BatchedWiVote.Vote) >= batchSize {
			node.BatchedWiVote.Vote = node.BatchedWiVote.Vote[batchSize:]
		} else {
			node.BatchedWiVote.Vote = make([]*big.Int, 0)
		}
		if len(node.BatchedWiVote.MajorityOfVotes) >= batchSize {
			node.BatchedWiVote.MajorityOfVotes = node.BatchedWiVote.MajorityOfVotes[batchSize:]
		} else {
			node.BatchedWiVote.MajorityOfVotes = make([]*big.Int, 0)
		}
	}
}

func (n *Node) processPendingBatchedWiVotes() {
	if n == nil || n.BatchedWiVote == nil {
		return
	}
	for n.batchedWiVoteCount() >= n.cfg.BatchSize {
		batchIDs := n.peekBatchedWiVoteIDs(n.cfg.BatchSize)
		if len(batchIDs) != n.cfg.BatchSize {
			return
		}
		majorityVote, err := n.processBatchedWiVotes(batchIDs)
		if err != nil {
			fmt.Printf("aggregator (node=%d) failed to process batched wiVotes: %v\n", n.ID, err)
			return
		}
		fmt.Printf("majority vote result %s\n", majorityVote.String())
		n.consumeBatchedWiVotes(n.cfg.BatchSize)
	}
}

func (n *Node) aggregatorCollectIncVote(incVote *IncVote) {
	n.incVotesLock.Lock()
	defer n.incVotesLock.Unlock()
	_, ok := n.IncVotes[string(incVote.CommitmentHash)]
	if !ok {
		n.IncVotes[string(incVote.CommitmentHash)] = make(map[uint]*IncVote)
	}
	n.IncVotes[string(incVote.CommitmentHash)][incVote.NodeID] = incVote
}

func (n *Node) aggregatorCollectWiVote(wiVote *WiVote) {
	n.wiVotesLock.Lock()
	defer n.wiVotesLock.Unlock()
	_, ok := n.WiVotes[wiVote.RequestID.String()]
	if !ok {
		n.WiVotes[wiVote.RequestID.String()] = make(map[uint]*WiVote)
	}
	n.WiVotes[wiVote.RequestID.String()][uint(wiVote.Index.Uint64())] = wiVote
}

func (n *Node) AggregatorSelectVote(commitmentHashStr string) (*IncVote, error) {
	n.incVotesLock.Lock()
	defer n.incVotesLock.Unlock()

	commitmentHashVotes, ok := n.IncVotes[commitmentHashStr]
	if !ok {
		return nil, fmt.Errorf("commitment hash not found in votes: commitmentHash=%v", commitmentHashStr)
	}

	var selectedIncVote *IncVote
	selection := make(map[string]int, len(commitmentHashVotes))
	for _, incVote := range commitmentHashVotes {
		incTreeIndexAndRootHash := fmt.Sprintf("%v-%v", incVote.IncTreeIndex, string(incVote.RootHash))
		selection[incTreeIndexAndRootHash] += 1
		nValidators := len(n.Oracle.Nodes)
		threshold := bftThreshold(nValidators)
		if selection[incTreeIndexAndRootHash] >= threshold {
			selectedIncVote = incVote
			break
		}

	}

	return selectedIncVote, nil
}

func (n *Node) selectIncVoteBatch(batchSize int) ([]*IncVote, []string, error) {
	if len(n.IncVotes) < batchSize {
		return nil, nil, nil
	}
	keys := make([]string, 0, len(n.IncVotes))
	for key := range n.IncVotes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	keys = keys[:batchSize]

	selected := make([]*IncVote, 0, batchSize)
	for _, key := range keys {
		selectedIncVote, err := n.AggregatorSelectVote(key)
		if err != nil {
			return nil, nil, err
		}
		if selectedIncVote == nil {
			return nil, nil, nil
		}
		selected = append(selected, selectedIncVote)
	}

	return selected, keys, nil
}

func (n *Node) AggregatorSelectWiVote(uniqueReqID string) (isApproved *big.Int, err error) {
	isApproved = big.NewInt(0)
	approvedCount := 0
	for _, wiVote := range n.WiVotes[uniqueReqID] {
		if wiVote.IsApproved.Uint64() == 1 {
			approvedCount++
			nValidators := n.cfg.NodeCount
			threshold := bftThreshold(nValidators)

			if approvedCount >= threshold {
				isApproved = big.NewInt(1)
				break
			}

		}
	}

	return isApproved, nil
}

func (n *Node) AggregatorProcessWiVote(uniqueReqID string) error {
	n.wiVotesLock.Lock()
	defer n.wiVotesLock.Unlock()

	votes, ok := n.WiVotes[uniqueReqID]
	if !ok {
		return fmt.Errorf("uniqueID not found in wiVotes: uniqueReqID=%v", uniqueReqID)
	}

	var voteSlice []*WiVote
	for _, v := range votes {
		voteSlice = append(voteSlice, v)
	}
	majorityWiVote, err := n.AggregatorSelectWiVote(uniqueReqID)
	if err != nil {
		return fmt.Errorf("failed to select wiVote: %v", err)
	}
	for _, v := range voteSlice {
		node, err := n.nodeByAccountIndex(v.Index)
		if err != nil {
			return err
		}
		fmt.Println("adding wivote to the bached wivote of the node:", node.ID)
		node.BatchedWiVote.Index = v.Index
		node.BatchedWiVote.WithdrawalReqIDs = append(node.BatchedWiVote.WithdrawalReqIDs, v.RequestID)
		node.BatchedWiVote.Vote = append(node.BatchedWiVote.Vote, v.IsApproved)
		node.BatchedWiVote.MajorityOfVotes = append(node.BatchedWiVote.MajorityOfVotes, majorityWiVote)
		node.BatchedWiVote.SenderPK = n.privateKey.PublicKey
	}

	return nil
}

func (n *Node) processBatchedWiVotes(withdrawalReqIDs []*big.Int) (*big.Int, error) {

	voteSlice := make([]*BatchedWiVote, 0)
	for i, node := range n.Oracle.Nodes {
		voteSlice = append(voteSlice, node.BatchedWiVote)
		fmt.Printf("Before sort - voteSlice[%d]: Index = %v\n", i, node.BatchedWiVote.Index)
	}

	// FIXED: Sort by Index
	sort.Slice(voteSlice, func(i, j int) bool {
		return voteSlice[i].Index.Uint64() < voteSlice[j].Index.Uint64()
	})

	// ✅ Print after sorting
	for i, vote := range voteSlice {
		fmt.Printf("After sort - voteSlice[%d]: Index = %v\n", i, vote.Index)
	}

	fmt.Printf("selected wivotes bathced: %v\n", voteSlice)

	preStateData := make([]byte, len(n.state.Data))
	copy(preStateData, n.state.Data)

	preStateHData := make([]byte, len(n.state.HData))
	copy(preStateHData, n.state.HData)

	if len(withdrawalReqIDs) != n.cfg.BatchSize {
		return nil, fmt.Errorf("batch size mismatch: expected %d got %d", n.cfg.BatchSize, len(withdrawalReqIDs))
	}

	batchIDs := append([]*big.Int(nil), withdrawalReqIDs...)
	sort.Slice(batchIDs, func(i, j int) bool {
		return batchIDs[i].Cmp(batchIDs[j]) < 0
	})
	if len(batchIDs) == 0 {
		return nil, fmt.Errorf("empty withdrawal request batch")
	}
	// Deterministic batching rule (paper SIV-D): round r contains exactly the
	// claim identifiers [r*b, (r+1)*b - 1]. The Gateway assigns sequential
	// IDs, so the sorted batch must form that window; its round identifier is
	// id/b.
	roundID, err := validateRoundWindow(batchIDs, n.cfg.BatchSize)
	if err != nil {
		return nil, err
	}

	// Batch commitment C_batch = H(id_1 .. id_b), binding votes to this
	// specific batch (paper SIV-D).
	hfunc := hash.MIMC_BN254.New()
	hfunc.Reset()
	var commitmentElem fr.Element
	for i := 0; i < len(batchIDs); i++ {
		commitmentElem.SetBigInt(batchIDs[i])
		hfunc.Write(commitmentElem.Marshal()[:])
	}
	batchCommitment := hfunc.Sum(nil)
	batchCommitment = util.PadOrTrim(util.ModToBn254Bytes(batchCommitment), 32)

	// Each validator computes and signs its own batch bitmask from its own
	// recorded votes (paper SIV-D Step 14, Alg. 2 lines 12-13). The
	// aggregator consumes the signed artifacts unmodified and never touches
	// another validator's key (F-11).
	signedVotes := make([]*SignedBatchVote, len(voteSlice))
	voteMasks := make([]*big.Int, len(voteSlice))
	voteCounts := make(map[string]int, len(voteSlice))
	for i, vote := range voteSlice {
		node, err := n.nodeByAccountIndex(vote.Index)
		if err != nil {
			return nil, err
		}
		signed, err := node.SignBatchVote(batchIDs, roundID, batchCommitment)
		if err != nil {
			return nil, err
		}
		signedVotes[i] = signed
		voteMasks[i] = new(big.Int).Set(signed.Vote)
		voteCounts[signed.Vote.String()]++
	}

	if len(voteMasks) == 0 {
		return nil, fmt.Errorf("no validator votes in batch")
	}

	majorityVote := new(big.Int).Set(voteMasks[0])
	majorityCount := voteCounts[majorityVote.String()]
	for _, mask := range voteMasks[1:] {
		count := voteCounts[mask.String()]
		if count > majorityCount {
			majorityVote = new(big.Int).Set(mask)
			majorityCount = count
		}
	}
	nValidators := n.cfg.NodeCount
	threshold := bftThreshold(nValidators)
	quorumSize := bftQuorumSize(nValidators)

	if majorityCount < threshold {
		return nil, fmt.Errorf("no BFT quorum for vote mask: max=%d threshold=%d", majorityCount, threshold)
	}
	if len(voteSlice) < quorumSize {
		return nil, fmt.Errorf("not enough votes for BFT quorum witness: have=%d need=%d", len(voteSlice), quorumSize)
	}

	voteSlice, signedVotes, voteMasks, err = selectQuorumBatchVotes(voteSlice, signedVotes, voteMasks, majorityVote, threshold, quorumSize)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Length of WithdrawalReqIDs: %d\n", len(batchIDs))

	fmt.Println("Sorted WithdrawalReqIDs:")
	for _, id := range batchIDs {
		fmt.Println(id)
	}

	fmt.Println("Sorted validator indexes in voteSlice:")
	for _, v := range voteSlice {
		fmt.Println(v.Index.Uint64())
	}

	aggregatorRootBytes, aggregatorProofBytes, err := n.state.MerkleProof(n.Account.Index.Uint64())
	if err != nil {
		return nil, fmt.Errorf("aggregator merkle proof: %w", err)
	}

	fmt.Printf("aggregatorRootBytes Aggregator state.MerleProof--------- %x\n", aggregatorRootBytes)
	aggregatorProof := make([]frontend.Variable, len(aggregatorProofBytes))
	for j := range aggregatorProofBytes {
		aggregatorProof[j] = frontend.Variable(new(big.Int).SetBytes(aggregatorProofBytes[j]))
	}

	preSeedX, preSeedY, err := n.getSeed()
	if err != nil {
		return nil, fmt.Errorf("get seed: %w", err)
	}

	preSeed := edwards.NewPointAffine(*new(fr.Element).SetBigInt(preSeedX), *new(fr.Element).SetBigInt(preSeedY))
	modulus := edwards.GetEdwardsCurve().Order
	sk := big.NewInt(0).SetBytes(n.privateKey.Bytes()[fp.Bytes : 2*fp.Bytes])
	sk.Mod(sk, &modulus)

	var postSeed edwards.PointAffine
	postSeed.ScalarMultiplication(&preSeed, sk)

	postSeedX := new(big.Int)
	postSeedY := new(big.Int)

	postSeed.X.BigInt(postSeedX)
	postSeed.Y.BigInt(postSeedY)
	aggregatorAccount, err := n.state.ReadAccount(n.Account.Index.Uint64())
	if err != nil {
		return nil, fmt.Errorf("read aggregator account: %w", err)
	}

	aggregatorConstraints := votingbatch.BatchingAggregatorConstraints{
		Index:     new(big.Int).Set(n.Account.Index),
		PreSeed:   twistededwards.Point{X: new(big.Int).Set(preSeedX), Y: new(big.Int).Set(preSeedY)},
		PostSeed:  twistededwards.Point{X: new(big.Int).Set(postSeedX), Y: new(big.Int).Set(postSeedY)},
		SecretKey: new(big.Int).Set(sk),
		Balance:   new(big.Int).Set(aggregatorAccount.Balance),
		MerkleProof: votingbatch.MerkleProofW{
			RootHash: aggregatorRootBytes,
			Path:     aggregatorProof,
		},
	}

	aggregatorAccount.Balance.Add(new(big.Int).Set(aggregatorAccount.Balance), big.NewInt(votingbatch.RewardAggregator))

	err = n.state.WriteAccount(aggregatorAccount)
	if err != nil {
		return nil, fmt.Errorf("write account: %w", err)
	}

	validatorConstraints := make([]votingbatch.BatchingValidatorConstraints, quorumSize)

	validatorBits := new(big.Int)
	honestBits := new(big.Int) //tracks validators matching majority

	fmt.Printf("Votes length--------- %x\n", len(voteSlice))

	// var tempPostState []byte
	for i, vote := range voteSlice {

		validatorVoteMask := new(big.Int).Set(voteMasks[i])

		// The validator's own signed batch bitmask (F-11): the aggregator
		// attaches it to the witness unmodified.
		vote.Signature = signedVotes[i].Signature

		// Read validator account from state
		validatorAccount, err := n.state.ReadAccount(vote.Index.Uint64())
		if err != nil {
			return nil, fmt.Errorf("read validator account: %w", err)
		}

		// 5. Use the correct public key & assign the signature
		var publicKey eddsa2.PublicKey
		var signature eddsa2.Signature

		publicKey.Assign(tedwards.BN254, validatorAccount.PublicKey.Bytes())
		signature.Assign(tedwards.BN254, vote.Signature)

		fmt.Printf("**************Validator [%d] index: %d **********\n", i, validatorAccount.Index.Uint64())
		rootBytes, proofBytes, err := n.state.MerkleProofBytes(validatorAccount.Index.Uint64())
		if err != nil {
			return nil, fmt.Errorf("validator merkle proof: %w", err)
		}
		fmt.Printf("rootBytes Validator state.MerleProof--------- %x\n", rootBytes)
		validatorProof := make([]frontend.Variable, len(proofBytes))
		for j := range proofBytes {
			validatorProof[j] = frontend.Variable(proofBytes[j])
		}

		validatorConstraints[i] = votingbatch.BatchingValidatorConstraints{
			//Index:     new(big.Int).Set(validatorAccount.Index),
			Index:     new(big.Int).Set(validatorAccount.Index),
			PublicKey: publicKey,
			Balance:   new(big.Int).Set(validatorAccount.Balance), //passed by reference
			MerkleProof: votingbatch.MerkleProofW{
				RootHash: rootBytes,
				Path:     validatorProof,
			},
			Signature: signature,
			Vote:      new(big.Int).Set(validatorVoteMask),
		}

		fmt.Printf("Validator[%d] LeafHash Inputs:\n", i)
		fmt.Printf("  Index     =  %x\n", validatorAccount.Index.String())
		fmt.Printf("  PubKey.X  =  %x\n", publicKey.A.X)
		fmt.Printf("  PubKey.Y  =  %x\n", publicKey.A.Y)
		fmt.Printf("  Balance   =  %x\n", validatorAccount.Balance.String())
		fmt.Printf("  Merkle Root = %x\n", rootBytes)

		validatorBit := new(big.Int)
		validatorBit.Exp(big.NewInt(2), new(big.Int).Set(vote.Index), nil)

		validatorBits = validatorBits.Add(new(big.Int).Set(validatorBits), new(big.Int).Set(validatorBit))
		// *****************************************************************************
		// //mark validator as honest if vote matches majority
		// if vote.IsApproved.Cmp(majorityWiVote) == 0 {
		// 	honestBits = honestBits.Add(
		// 		new(big.Int).Set(honestBits),
		// 		new(big.Int).Set(validatorBit),
		// 	)
		// }
		// // slash
		// if vote.IsApproved.Cmp(majorityWiVote) == 0 {
		// 	// honest validator → fixed reward
		// 	validatorAccount.Balance.Add(
		// 		new(big.Int).Set(validatorAccount.Balance),
		// 		big.NewInt(gnark.RewardValidator),
		// 	)
		// } else {
		// 	// dishonest validator → FULL SLASH
		// 	validatorAccount.Balance.SetInt64(0)
		// }

		validatorVote := new(big.Int).Set(validatorVoteMask)

		fmt.Printf("DEBUG honest calc: idx=%d validatorVote=%s majorityVote=%s\n",
			vote.Index.Uint64(), validatorVote.String(), majorityVote.String(),
		)

		// mark validator as honest if vote matches majority
		if validatorVote.Cmp(majorityVote) == 0 {
			honestBits = honestBits.Add(
				new(big.Int).Set(honestBits),
				new(big.Int).Set(validatorBit),
			)

			// honest validator → fixed reward
			validatorAccount.Balance.Add(
				new(big.Int).Set(validatorAccount.Balance),
				big.NewInt(votingbatch.RewardValidator),
			)
		} else {
			penalty := big.NewInt(votingbatch.PenaltyValidator)
			if validatorAccount.Balance.Cmp(penalty) <= 0 {
				validatorAccount.Balance.SetInt64(0)
			} else {
				validatorAccount.Balance.Sub(validatorAccount.Balance, penalty)
			}

		}

		err = n.state.WriteAccount(validatorAccount)
		if err != nil {
			return nil, fmt.Errorf("write account: %w", err)
		}
	}

	postStateRoot, err := n.state.Root()
	if err != nil {
		return nil, fmt.Errorf("state root: %w", err)
	}

	// uniqueReqIdInt, ok := new(big.Int).SetString(uniqueReqID, 10)
	// if !ok {
	// 	panic("failed to convert uniqueReqID to big.Int")
	// }

	// fmt.Printf("POST STATE POST: %x  --- TEMP POST STATE: %x ", postStateRoot, tempPostState)

	withdrawalReqVars := make([]frontend.Variable, len(batchIDs))
	for i := range batchIDs {
		withdrawalReqVars[i] = frontend.Variable(batchIDs[i])
	}

	assignment := votingbatch.BatchingVotingCircuit{
		ResultingStateRoot: postStateRoot,
		//RoundID:            n.Oracle.RoundID,
		RoundID:          new(big.Int).Set(roundID),
		BatchCommitment:  batchCommitment[:32],
		MajorityVote:     new(big.Int).Set(majorityVote),
		ValidatorBits:    new(big.Int).Set(validatorBits),
		HonestBits:       new(big.Int).Set(honestBits), // NEW
		WithdrawalReqIDs: withdrawalReqVars,
		Aggregator:       aggregatorConstraints,
		Validators:       validatorConstraints,
	}

	// Print all variables used in the assignment
	fmt.Printf("ResultingStateRoot: %v\n", postStateRoot)
	fmt.Printf("RoundID: %v\n", new(big.Int).Set(roundID))
	fmt.Printf("BatchCommitment: %v\n", batchCommitment[:32])
	fmt.Printf("MajorityVote: %v\n", new(big.Int).Set(majorityVote))
	fmt.Printf("ValidatorBits: %v\n", new(big.Int).Set(validatorBits))
	fmt.Printf("WithdrawalReqIDs: %v\n", withdrawalReqVars)
	fmt.Printf("Aggregator: %+v\n", aggregatorConstraints)
	fmt.Printf("Validators: %+v\n", validatorConstraints)

	var (
		w witness.Witness
		p groth16.Proof
	)

	_, err = memtime.MeasurePeak(
		"groth16.Prove voting_batch",
		memtime.PeakSampleInterval,
		func() error {
			var e error
			w, e = frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
			if e != nil {
				return fmt.Errorf("create witness: %w", e)
			}

			p, e = groth16.Prove(n.Oracle.SparseR1CS, n.Oracle.SparsePK, w)
			return e
		},
	)
	if err != nil {
		return nil, fmt.Errorf("prove: %w", err)
	}

	pw, err := w.Public()
	if err != nil {
		return nil, fmt.Errorf("public witness: %w", err)
	}

	_, err = memtime.MeasurePeak(
		"groth16.Verify voting_batch",
		memtime.PeakSampleInterval,
		func() error {
			return groth16.Verify(p, n.Oracle.SparseVK, pw)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("verify proof: %w", err)
	}

	proof, err := util.ProofToEthereumProof(p)
	if err != nil {
		return nil, fmt.Errorf("proof to ethereum proof: %w", err)
	}

	//Reset state
	n.state.SetData(preStateData)
	n.state.SetHData(preStateHData)

	fmt.Println(" --- Scalar Field Check Complete --- ")

	check("aggregator account index", n.Account.Index)
	check("slashedValIndex", roundID)
	check("uniqueReqID", roundID)
	check("postStateRoot", new(big.Int).SetBytes(postStateRoot))
	check("postSeedX", postSeedX)
	check("postSeedY", postSeedY)
	check("batchCommitment", new(big.Int).SetBytes(batchCommitment))
	check("majorityVote", majorityVote)
	check("validatorBits", validatorBits)

	if err := n.aggregatorSubmitWiVoteTx(
		n.Account.Index,
		roundID,
		new(big.Int).SetBytes(batchCommitment[:32]),
		validatorBits,
		honestBits, // NEW
		new(big.Int).Set(majorityVote),
		new(big.Int).SetBytes(postStateRoot),
		postSeedX,
		postSeedY,
		proof.Proof,
	); err != nil {
		return nil, err
	}

	fmt.Println("aggregatorSubmitWiVoteTx submitted successfully")

	return majorityVote, nil
}

func hashWiVoteFieldwise(wivote *WiVote) []byte {
	hfunc := hash.MIMC_BN254.New()
	hfunc.Reset()

	var fe fr.Element

	fe.SetBigInt(wivote.Index)
	hfunc.Write(fe.Marshal()[:]) // feeds the field element in canonical form

	fe.SetBigInt(wivote.RequestID)
	hfunc.Write(fe.Marshal()[:])

	fe.SetBigInt(wivote.IsApproved)
	hfunc.Write(fe.Marshal()[:])

	return hfunc.Sum(nil)
}

func hashBatchVoteFieldwise(index, batchCommitment, vote, roundID *big.Int) []byte {
	hFunc := hash.MIMC_BN254.New()
	hFunc.Reset()

	var fe fr.Element

	fe.SetBigInt(index)
	hFunc.Write(fe.Marshal()[:])

	fe.SetBigInt(batchCommitment)
	hFunc.Write(fe.Marshal()[:])

	fe.SetBigInt(vote)
	hFunc.Write(fe.Marshal()[:])

	fe.SetBigInt(roundID)
	hFunc.Write(fe.Marshal()[:])

	return hFunc.Sum(nil)
}

// func marshalHash(Index *big.Int, RequestID *big.Int, IsApproved *big.Int) []byte {
// 	hfunc := hash.MIMC_BN254.New()
// 	hfunc.Reset()

// 	var fe fr.Element

// 	fe.SetBigInt(Index)
// 	hfunc.Write(fe.Marshal()[:]) // feeds the field element in canonical form

// 	fe.SetBigInt(RequestID)
// 	hfunc.Write(fe.Marshal()[:])

// 	fe.SetBigInt(IsApproved)
// 	hfunc.Write(fe.Marshal()[:])

// 	return hfunc.Sum(nil)
// }

// endregion

// region contract

func (n *Node) newTransactOpts() (*bind.TransactOpts, error) {
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		return nil, err
	}
	if n.cfg.TxGasLimit > 0 {
		// Avoid RPC gas estimation when the endpoint doesn't support eth_estimateGas.
		trxOpts.GasLimit = uint64(n.cfg.TxGasLimit)
	}
	if n.cfg.TxGasPriceWei > 0 {
		// Avoid RPC gas price discovery when eth_gasPrice isn't available.
		trxOpts.GasPrice = big.NewInt(n.cfg.TxGasPriceWei)
	} else if n.cfg.TxGasFeeCapWei > 0 || n.cfg.TxGasTipCapWei > 0 {
		// Allow explicit EIP-1559 values without RPC lookups.
		if n.cfg.TxGasFeeCapWei > 0 {
			trxOpts.GasFeeCap = big.NewInt(n.cfg.TxGasFeeCapWei)
		}
		if n.cfg.TxGasTipCapWei > 0 {
			trxOpts.GasTipCap = big.NewInt(n.cfg.TxGasTipCapWei)
		}
	}
	return trxOpts, nil
}

func (n *Node) selectNewAggregatorTx() error {
	fmt.Printf("selecting new aggregator from contract...\n")
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	tx, err := bcClient.ChooseNewAggregator(trxOpts)
	if err != nil {
		log.Fatalf("call SelectNewAggregatorTx() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	bc.LogTxReceipt("select new aggregator", tx, receipt)
	if receipt.Status == 1 {
		fmt.Printf("successfully called SelectNewAggregatorTx (by oracle)\n")
	} else {
		fmt.Printf("Transaction failed (by oracle)\n")
	}

	return nil
}

func (n *Node) RegisterValidatorTx() error {
	fmt.Printf("registering validator=%v ...\n", n.ID)
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	// Stake the validator's collateral: the on-chain leaf is
	// MiMC(index, pk.x, pk.y, msg.value), so the registered value must equal
	// the off-chain account balance for Merkle proofs to verify (F-19).
	trxOpts.Value = new(big.Int).Set(n.Account.Balance)

	pk := gnark.PublicKeyToOraclePublicKey(n.Account.PublicKey)
	tx, err := bcClient.RegisterValidator(trxOpts, big.NewInt(int64(n.ID)), *pk)
	if err != nil {
		log.Fatalf("call RegisterValidator() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	bc.LogTxReceipt(fmt.Sprintf("register validator node=%v", n.ID), tx, receipt)
	if receipt.Status == 1 {
		fmt.Printf("successfully registered validator node=%v\n", n.ID)
	} else {
		fmt.Printf("Transaction failed (node=%v\n)", n.ID)
		revertReason, callErr := n.revertReason(context.Background(), tx, receipt.BlockNumber)
		if revertReason != "" {
			return fmt.Errorf("register validator failed (node=%v tx=%s revert=%s)", n.ID, tx.Hash().Hex(), revertReason)
		}
		if callErr != "" {
			return fmt.Errorf("register validator failed (node=%v tx=%s call_err=%s)", n.ID, tx.Hash().Hex(), callErr)
		}
		return fmt.Errorf("register validator failed (node=%v tx=%s)", n.ID, tx.Hash().Hex())
	}

	return nil
}

func (n *Node) revertReason(ctx context.Context, tx *types.Transaction, blockNumber *big.Int) (string, string) {
	if n == nil || n.ethClient == nil || tx == nil {
		return "", ""
	}
	if ctx == nil {
		ctx = context.Background()
	}
	to := tx.To()
	if to == nil {
		return "", ""
	}

	if n.ecdsaPrivateKey == nil {
		return "", ""
	}

	msg := ethereum.CallMsg{
		From:  crypto.PubkeyToAddress(n.ecdsaPrivateKey.PublicKey),
		To:    to,
		Gas:   tx.Gas(),
		Value: tx.Value(),
		Data:  tx.Data(),
	}

	switch tx.Type() {
	case types.LegacyTxType, types.AccessListTxType:
		msg.GasPrice = tx.GasPrice()
	case types.DynamicFeeTxType:
		msg.GasFeeCap = tx.GasFeeCap()
		msg.GasTipCap = tx.GasTipCap()
	}

	data, err := n.ethClient.CallContract(ctx, msg, blockNumber)
	if err != nil && len(data) == 0 {
		return "", err.Error()
	}
	if len(data) == 0 {
		return "", ""
	}
	if reason := decodeRevertData(data); reason != "" {
		return reason, ""
	}
	return fmt.Sprintf("reverted (data=0x%x)", data), ""
}

func decodeRevertData(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	reason, err := abi.UnpackRevert(data)
	if err == nil && strings.TrimSpace(reason) != "" {
		return reason
	}
	if len(data) < 4 {
		return ""
	}
	selector := data[:4]
	switch {
	case bytes.Equal(selector, proofInvalidSelector):
		return "ProofInvalid"
	case bytes.Equal(selector, publicInputNotInFieldSelector):
		return "PublicInputNotInField"
	default:
		return fmt.Sprintf("custom error 0x%x (data=0x%x)", selector, data)
	}
}

// publishCommitmentRootTx reports the finalized commitment-tree root and the
// DFS object identifier to the Gateway (paper SIV-E Steps 6-7). Only the
// current round aggregator's transaction is accepted on-chain.
func (n *Node) publishCommitmentRootTx(root *big.Int, dfsRef string, label string) error {
	if label == "" {
		label = "L2 publish IPFS root"
	}
	fmt.Printf("publishing commitment root (validator=%v root=%s) ...\n", n.ID, root)
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	tx, err := bcClient.PublishCommitmentRoot(trxOpts, root, dfsRef)
	if err != nil {
		log.Fatalf("call publishCommitmentRoot() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	bc.LogTxReceipt(label, tx, receipt)
	if receipt.Status == 1 {
		fmt.Printf("successfully published commitment root (node=%v)\n", n.ID)
	} else {
		revertReason, callErr := n.revertReason(context.Background(), tx, receipt.BlockNumber)
		if revertReason != "" {
			return fmt.Errorf("%s reverted (node=%v tx=%s revert=%s)", label, n.ID, tx.Hash().Hex(), revertReason)
		}
		if callErr != "" {
			return fmt.Errorf("%s failed (node=%v tx=%s call_err=%s)", label, n.ID, tx.Hash().Hex(), callErr)
		}
		return fmt.Errorf("%s failed (node=%v tx=%s)", label, n.ID, tx.Hash().Hex())
	}

	return nil
}

func (n *Node) aggregatorSubmitWiVoteTx(index *big.Int, uniqueReqID *big.Int, batchCommitment *big.Int, validatorBits *big.Int, honestBits *big.Int, vote *big.Int, postStateRoot *big.Int, postSeedX *big.Int, postSeedY *big.Int, proof [8]*big.Int) error {
	fmt.Printf("submitting wivote (node=%v) ...\n", n.ID)
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	tx, err := bcClient.SubmitWiVote(
		trxOpts,
		index,
		uniqueReqID,
		batchCommitment,
		toScalarField(validatorBits),
		toScalarField(honestBits),
		vote,
		postStateRoot,
		postSeedX,
		postSeedY,
		proof,
	)
	if err != nil {
		log.Fatalf("call SubmitWiVote() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	bc.LogTxReceipt("L2 submitWiVote batch finalization", tx, receipt)

	if receipt.Status == 1 {
		fmt.Printf("successfully submitted wivote (node=%v)\n", n.ID)
	} else {
		fmt.Printf("Transaction failed (node=%v)\n", n.ID)
		revertReason, callErr := n.revertReason(context.Background(), tx, receipt.BlockNumber)
		if revertReason != "" {
			log.Printf("submit wivote reverted (node=%v tx=%s): %s", n.ID, tx.Hash().Hex(), revertReason)
			return fmt.Errorf("submit wivote failed (node=%v tx=%s revert=%s)", n.ID, tx.Hash().Hex(), revertReason)
		}
		if callErr != "" {
			log.Printf("submit wivote call failed (node=%v tx=%s): %s", n.ID, tx.Hash().Hex(), callErr)
			return fmt.Errorf("submit wivote failed (node=%v tx=%s call_err=%s)", n.ID, tx.Hash().Hex(), callErr)
		}
		return fmt.Errorf("submit wivote failed (node=%v tx=%s)", n.ID, tx.Hash().Hex())
	}

	return nil
}

// awaitStateSync blocks until this node's local validator-state root equals
// the Gateway's on-chain root. Round finalization (WiVoteSubmitted) and leaf
// replacements are applied to local state asynchronously by the event
// subscriber; a membership proof built before those events land would open a
// stale root and revert on-chain (F-19 enforcement).
func (n *Node) awaitStateSync(timeout time.Duration) error {
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		return fmt.Errorf("create contract client instance: %w", err)
	}
	deadline := time.Now().Add(timeout)
	for {
		onChainRoot, err := bcClient.GetRoot(&bind.CallOpts{})
		if err != nil {
			return fmt.Errorf("read on-chain validator-state root: %w", err)
		}
		localRoot, err := n.state.Root()
		if err != nil {
			return fmt.Errorf("compute local state root: %w", err)
		}
		if onChainRoot.Cmp(new(big.Int).SetBytes(localRoot)) == 0 {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("local validator state not synced with on-chain root within %s (node=%d local=%x onchain=%s)",
				timeout, n.ID, localRoot, onChainRoot.String())
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// AwaitStateSync exposes the validator-tree sync barrier to HTTP/setup code.
func (n *Node) AwaitStateSync(timeout time.Duration) error {
	return n.awaitStateSync(timeout)
}

// AwaitAccountIndex waits until asynchronous L1->L2 lifecycle events have
// updated the live node handle to the expected validator-tree leaf.
func (n *Node) AwaitAccountIndex(index uint64, timeout time.Duration) error {
	if n == nil || n.Account == nil || n.Account.Index == nil {
		return fmt.Errorf("node account not initialized")
	}
	deadline := time.Now().Add(timeout)
	for {
		if n.Account.Index.Uint64() == index {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("node account index not synced within %s (node=%d current=%s expected=%d)", timeout, n.ID, n.Account.Index.String(), index)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (n *Node) ReplaceAccountTx(replaceWithAccountID uint64) error {
	fmt.Printf("starting to replace account (node=%v)...\n", n.ID)

	if err := n.awaitStateSync(validatorStateSyncTimeout); err != nil {
		return fmt.Errorf("replace: %w", err)
	}

	account, err := n.state.ReadAccount(n.Account.Index.Uint64())
	if err != nil {
		log.Fatalf("read account at index %d: %v", account.Index.Uint64(), err)
	}

	replaceAccount, err := n.state.ReadAccount(replaceWithAccountID)
	if err != nil {
		log.Fatalf("read replace account at index %d: %v", replaceAccount.Index.Uint64(), err)
	}

	fmt.Printf("Account index=%d with balance:%s is replacing with account inedx=%d balance=%s\n", account.Index.Uint64(), account.Balance.String(), replaceAccount.Index.Uint64(), replaceAccount.Balance.String())

	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	// Paper SIV-C: the replacement stake must strictly exceed the incumbent's.
	replacementStake := new(big.Int).Add(replaceAccount.Balance, big.NewInt(1))
	trxOpts.Value = replacementStake

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	// The Merkle proof must open the incumbent's leaf: the contract binds
	// path[0] to hashAccount(toReplace) and verifies it at the target index.
	_, path, err := n.state.MerkleProofBytes(replaceAccount.Index.Uint64())
	if err != nil {
		log.Fatalf("merkle proof failed (index=%d): %v", replaceAccount.Index.Uint64(), err)
	}
	fmt.Printf("*** About to replace account:\n")
	fmt.Printf("  Index         : %s\n", replaceAccount.Index.String())
	fmt.Printf("  PublicKey.X   : %s\n", replaceAccount.PublicKey.A.X.String())
	fmt.Printf("  PublicKey.Y   : %s\n", replaceAccount.PublicKey.A.Y.String())
	fmt.Printf("  Balance       : %s\n", replaceAccount.Balance.String())

	tx, err := bcClient.Replace(
		trxOpts,
		*gnark.PublicKeyToOraclePublicKey(account.PublicKey),
		*gnark.AccountToOracleAccount(&replaceAccount),
		path[:], // Convert fixed-size array to slice
		new(big.Int).Set(replaceAccount.Index),
		big.NewInt(int64(n.cfg.SparseTreeDepth)),
	)
	if err != nil {
		log.Fatalf("replace failed (index=%d): %v", account.Index.Uint64(), err)
	}

	// 5. Wait for confirmation
	receipt, err := bc.WaitMinedAndLogTxReceipt(
		context.Background(),
		n.ethClient,
		fmt.Sprintf("replace account index=%d", account.Index.Uint64()),
		tx,
	)
	if err != nil {
		log.Fatalf("failed to wait for replace tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("replace tx reverted (index=%d tx=%s)", account.Index.Uint64(), tx.Hash().Hex())
	}

	// 6. Update local state to mirror the on-chain leaf replacement: the
	// target slot now holds the caller's key with the new (higher) stake.
	replaceAccount.PublicKey = account.PublicKey
	replaceAccount.Balance = replacementStake
	if err := n.state.WriteAccount(replaceAccount); err != nil {
		log.Fatalf("failed to update account after replace: %v", err)
	}

	return nil
}

func (n *Node) RequestReplacementFromL1Tx(replaceWithAccountID uint64) error {
	fmt.Printf("requesting L1 replacement (node=%v)...\n", n.ID)

	if n.cfg == nil || n.cfg.L1RPCURL == "" || n.cfg.L1HubContractAddress == "" {
		return n.ReplaceAccountTx(replaceWithAccountID)
	}
	if err := n.awaitStateSync(validatorStateSyncTimeout); err != nil {
		return fmt.Errorf("l1 replacement: %w", err)
	}

	targetAccount, err := n.state.ReadAccount(replaceWithAccountID)
	if err != nil {
		return fmt.Errorf("read replacement target account at index %d: %w", replaceWithAccountID, err)
	}
	_, path, err := n.state.MerkleProofBytes(targetAccount.Index.Uint64())
	if err != nil {
		return fmt.Errorf("replacement target merkle proof (index=%d): %w", targetAccount.Index.Uint64(), err)
	}

	candidateStake := new(big.Int).Add(targetAccount.Balance, big.NewInt(1))
	l1Client, err := bc.NewChainClient(context.Background(), n.cfg.L1RPCURL, n.cfg.L1ChainID, n.cfg.NodePK)
	if err != nil {
		return fmt.Errorf("init L1 replacement client: %w", err)
	}
	if n.cfg.L1GasPriceWei > 0 {
		l1Client.GasPriceOverride = big.NewInt(n.cfg.L1GasPriceWei)
	}
	hub, err := bc.NewL1Hub(common.HexToAddress(n.cfg.L1HubContractAddress), l1Client.Eth)
	if err != nil {
		return fmt.Errorf("bind L1Hub: %w", err)
	}
	if err := waitL1ValidatorImported(context.Background(), hub, targetAccount.Index, 90*time.Second); err != nil {
		return fmt.Errorf("wait replacement target imported on L1: %w", err)
	}

	feeValue := l1ToL2RequestValue(n.cfg.L1L2ValueWei)
	totalValue := new(big.Int).Add(candidateStake, feeValue)
	auth, err := newL1HubNodeTransactor(context.Background(), l1Client, totalValue)
	if err != nil {
		return fmt.Errorf("L1 replacement auth: %w", err)
	}

	params := bc.L1HubReplacementParams{
		TargetValidatorID:    new(big.Int).Set(targetAccount.Index),
		CandidateValidatorID: big.NewInt(int64(n.ID)),
		CandidatePubKey: bc.L1HubPublicKey{
			X: n.Account.PublicKey.A.X.BigInt(new(big.Int)),
			Y: n.Account.PublicKey.A.Y.BigInt(new(big.Int)),
		},
		CandidateStake:  candidateStake,
		TargetLeafIndex: new(big.Int).Set(targetAccount.Index),
		TargetPubKey: bc.L1HubPublicKey{
			X: targetAccount.PublicKey.A.X.BigInt(new(big.Int)),
			Y: targetAccount.PublicKey.A.Y.BigInt(new(big.Int)),
		},
		TargetBalance: new(big.Int).Set(targetAccount.Balance),
		Path:          path[:],
		Depth:         big.NewInt(int64(n.cfg.SparseTreeDepth)),
	}

	tx, err := hub.RequestReplacementL1(
		auth,
		params,
		big.NewInt(30_000_000),
		big.NewInt(800),
		l1Client.From,
	)
	if err != nil {
		return fmt.Errorf("request replacement on L1: %w", err)
	}
	receipt, err := bc.WaitMinedAndLogTxReceipt(
		context.Background(),
		l1Client.Eth,
		fmt.Sprintf("L1->L2 validator replacement request index=%d", targetAccount.Index.Uint64()),
		tx,
	)
	if err != nil {
		return fmt.Errorf("wait L1 replacement request: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("L1 replacement request reverted (tx=%s)", tx.Hash().Hex())
	}
	return nil
}

func waitL1ValidatorImported(ctx context.Context, hub *bc.L1Hub, validatorID *big.Int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		record, err := hub.Validators(&bind.CallOpts{Context: ctx}, validatorID)
		if err != nil {
			return err
		}
		if record.ImportRequested {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("validator %s is not imported after %s", validatorID.String(), timeout)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func newL1HubNodeTransactor(ctx context.Context, chain *bc.ChainClient, valueWei *big.Int) (*bind.TransactOpts, error) {
	nonce, err := chain.Eth.PendingNonceAt(ctx, chain.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}
	gasPrice := chain.GasPriceOverride
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice, err = chain.Eth.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(chain.PrivKey, chain.ChainID)
	if err != nil {
		return nil, fmt.Errorf("new transactor: %w", err)
	}
	if valueWei == nil {
		valueWei = big.NewInt(0)
	}
	auth.From = chain.From
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = valueWei
	auth.GasPrice = gasPrice
	return auth, nil
}

func l1ToL2RequestValue(raw string) *big.Int {
	if strings.TrimSpace(raw) == "" {
		return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	}
	v, ok := new(big.Int).SetString(strings.TrimSpace(raw), 10)
	if !ok || v.Sign() < 0 {
		return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	}
	return v
}

func (n *Node) ExitTx() error {
	fmt.Printf("exiting account (node=%v)...\n", n.ID)

	if err := n.awaitStateSync(validatorStateSyncTimeout); err != nil {
		return fmt.Errorf("exit: %w", err)
	}

	account, err := n.state.ReadAccount(n.Account.Index.Uint64())
	if err != nil {
		log.Fatalf("read account at index %d: %v", account.Index.Uint64(), err)
	}

	fmt.Printf("Account index exit=%d balance=%s\n", account.Index.Uint64(), account.Balance.String())

	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	_, path, err := n.state.MerkleProofBytes(account.Index.Uint64())
	if err != nil {
		log.Fatalf("merkle proof failed (index=%d): %v", account.Index.Uint64(), err)
	}

	// 4. Call withdraw
	tx, err := bcClient.Exit(
		trxOpts,
		*gnark.AccountToOracleAccount(&account),
		path[:], // Convert fixed-size array to slice
		new(big.Int).Set(account.Index),
		big.NewInt(int64(n.cfg.SparseTreeDepth)),
	)
	if err != nil {
		log.Fatalf("exit failed (index=%d): %v", account.Index.Uint64(), err)
	}

	// 5. Wait for confirmation
	receipt, err := bc.WaitMinedAndLogTxReceipt(
		context.Background(),
		n.ethClient,
		fmt.Sprintf("L2 validator exit request index=%d", account.Index.Uint64()),
		tx,
	)
	if err != nil {
		log.Fatalf("failed to wait for exit tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("exit tx reverted (index=%d tx=%s)", account.Index.Uint64(), tx.Hash().Hex())
	}

	return nil
}

func (n *Node) WithdrawAccountTx() error {
	fmt.Printf("starting withdraw account (node=%v)...\n", n.ID)

	if err := n.awaitStateSync(validatorStateSyncTimeout); err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}

	account, err := n.state.ReadAccount(n.Account.Index.Uint64())
	if err != nil {
		log.Fatalf("read account at index %d: %v", account.Index.Uint64(), err)
	}

	// if account.Balance.Sign() == 0 {
	// 	return nil // skip zero balance accounts
	// }

	fmt.Printf("Account index withdraw=%d balance=%s\n", account.Index.Uint64(), account.Balance.String())

	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}

	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	_, path, err := n.state.MerkleProofBytes(account.Index.Uint64())
	if err != nil {
		log.Fatalf("merkle proof failed (index=%d): %v", account.Index.Uint64(), err)
	}

	// 4. Call withdraw
	tx, err := bcClient.Withdraw(
		trxOpts,
		*gnark.AccountToOracleAccount(&account),
		path[:], // Convert fixed-size array to slice
		new(big.Int).Set(account.Index),
		big.NewInt(int64(n.cfg.SparseTreeDepth)),
	)
	if err != nil {
		log.Fatalf("withdraw failed (index=%d): %v", account.Index.Uint64(), err)
	}

	// 5. Wait for confirmation
	receipt, err := bc.WaitMinedAndLogTxReceipt(
		context.Background(),
		n.ethClient,
		fmt.Sprintf("L2 validator withdraw request index=%d", account.Index.Uint64()),
		tx,
	)
	if err != nil {
		log.Fatalf("failed to wait for withdraw tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status == 1 {
		log.Printf("withdrawn account index=%d", account.Index.Uint64())
	} else {
		return fmt.Errorf("withdraw tx reverted (index=%d tx=%s)", account.Index.Uint64(), tx.Hash().Hex())
	}
	//averageCost = averageCost / uint64(len(accounts))

	// 6. Update local state
	account.Balance = big.NewInt(0)
	if err := n.state.WriteAccount(account); err != nil {
		log.Fatalf("failed to update account after withdraw: %v", err)
	}

	return nil
}

// endregion

func toScalarField(value *big.Int) *big.Int {
	mod := new(big.Int)
	mod.SetString("30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001", 16)
	return new(big.Int).Mod(value, mod)
}

func (n *Node) getSeed() (*big.Int, *big.Int, error) {
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	trxOpts, err := n.newTransactOpts()
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	callOpts := &bind.CallOpts{
		Context: context.Background(),
		From:    trxOpts.From,
	}
	preSeedX, preSeedY, err := bcClient.GetSeed(callOpts)
	if err != nil {
		log.Fatalf("call ViewLatestIPFSHash() function: %v", err)
	}

	return preSeedX, preSeedY, err
}

func check(name string, val *big.Int) {
	R := new(big.Int)
	R.SetString("30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001", 16)

	if val.Cmp(R) >= 0 {
		fmt.Printf("❌ %s exceeds scalar field R\n  → value: %s\n", name, val.String())
	} else {
		fmt.Printf("✅ %s is within scalar field\n  → value: %s\n", name, val.String())
	}
}

// bftThreshold returns the paper's finalization threshold (SIV-A/SIV-D,
// Alg. 2 line 19): f + 1 with f = floor((n-1)/3), i.e. at least one honest
// supporting validator under the BFT assumption n = 3f + 1.
func bftThreshold(nValidators int) int {
	if nValidators <= 0 {
		return 0
	}
	f := (nValidators - 1) / 3
	return f + 1
}

// bftQuorumSize returns the number of validator responses represented in the
// voting proof. The full committee has n=3f+1 validators, but the aggregator
// only needs a quorum of 2f+1 responses to finalize.
func bftQuorumSize(nValidators int) int {
	if nValidators <= 0 {
		return 0
	}
	f := (nValidators - 1) / 3
	return 2*f + 1
}

func selectQuorumBatchVotes(votes []*BatchedWiVote, signed []*SignedBatchVote, masks []*big.Int, majorityVote *big.Int, agreementThreshold int, quorumSize int) ([]*BatchedWiVote, []*SignedBatchVote, []*big.Int, error) {
	if len(votes) != len(signed) || len(votes) != len(masks) {
		return nil, nil, nil, fmt.Errorf("vote slice length mismatch")
	}
	if quorumSize <= 0 || agreementThreshold <= 0 {
		return nil, nil, nil, fmt.Errorf("invalid quorum parameters: quorum=%d agreement=%d", quorumSize, agreementThreshold)
	}
	if len(votes) < quorumSize {
		return nil, nil, nil, fmt.Errorf("not enough votes for quorum: have=%d need=%d", len(votes), quorumSize)
	}

	selected := make([]bool, len(votes))
	selectedVotes := make([]*BatchedWiVote, 0, quorumSize)
	selectedSigned := make([]*SignedBatchVote, 0, quorumSize)
	selectedMasks := make([]*big.Int, 0, quorumSize)
	agreementCount := 0

	add := func(i int) {
		selected[i] = true
		selectedVotes = append(selectedVotes, votes[i])
		selectedSigned = append(selectedSigned, signed[i])
		selectedMasks = append(selectedMasks, masks[i])
	}

	for i, mask := range masks {
		if agreementCount >= agreementThreshold {
			break
		}
		if mask != nil && mask.Cmp(majorityVote) == 0 {
			add(i)
			agreementCount++
		}
	}
	if agreementCount < agreementThreshold {
		return nil, nil, nil, fmt.Errorf("not enough agreeing votes for quorum witness: have=%d need=%d", agreementCount, agreementThreshold)
	}

	for i := range votes {
		if len(selectedVotes) >= quorumSize {
			break
		}
		if !selected[i] {
			add(i)
		}
	}
	if len(selectedVotes) != quorumSize {
		return nil, nil, nil, fmt.Errorf("failed to build quorum witness: have=%d need=%d", len(selectedVotes), quorumSize)
	}

	return selectedVotes, selectedSigned, selectedMasks, nil
}

// validateRoundWindow enforces the deterministic batching rule (paper
// SIV-D): round r contains exactly the claim identifiers
// [r*b, (r+1)*b - 1]. The ascending batch must form that window; the round
// identifier is id/b.
func validateRoundWindow(batchIDs []*big.Int, batchSize int) (*big.Int, error) {
	if batchSize <= 0 {
		return nil, fmt.Errorf("invalid batch size %d", batchSize)
	}
	if len(batchIDs) != batchSize {
		return nil, fmt.Errorf("batch size mismatch: expected %d got %d", batchSize, len(batchIDs))
	}
	b := big.NewInt(int64(batchSize))
	roundID := new(big.Int).Div(batchIDs[0], b)
	windowBase := new(big.Int).Mul(roundID, b)
	for i, id := range batchIDs {
		expected := new(big.Int).Add(windowBase, big.NewInt(int64(i)))
		if id.Cmp(expected) != 0 {
			return nil, fmt.Errorf("batch does not form round window [%s..%s]: position %d has id %s, want %s",
				windowBase, new(big.Int).Add(windowBase, big.NewInt(int64(batchSize-1))), i, id, expected)
		}
	}
	return roundID, nil
}
