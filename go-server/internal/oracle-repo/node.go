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

	BatchedWiVote *BatchedWiVote
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
								panic(fmt.Errorf("invalid wiVote request id: %q", uniqueReqID))
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
								panic(fmt.Errorf("failed to process wiVote: %v", err))
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
						latestIPFSHash, _ := n.UpdateIPFS()
						n.updateLatestIPFSHashTx(latestIPFSHash)
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

	var isApproved *big.Int
	if err := groth16.Verify(proof, n.Oracle.IncVK, publicWitness); err != nil {
		fmt.Printf("failed to verify claim proof (node=%v): %v\n", n.ID, err)
		isApproved = big.NewInt(0)
	} else {
		fmt.Printf("successfully verified claim proof (node=%v)\n", n.ID)
		isApproved = big.NewInt(1)
	}

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
	b := big.NewInt(int64(n.cfg.BatchSize))
	roundID := new(big.Int).Div(batchIDs[0], b)
	windowBase := new(big.Int).Mul(roundID, b)
	for i, id := range batchIDs {
		expected := new(big.Int).Add(windowBase, big.NewInt(int64(i)))
		if id.Cmp(expected) != 0 {
			return nil, fmt.Errorf("batch does not form round window [%s..%s]: position %d has id %s, want %s",
				windowBase, new(big.Int).Add(windowBase, big.NewInt(int64(n.cfg.BatchSize-1))), i, id, expected)
		}
	}

	positions := make(map[string]int, len(batchIDs))
	for i, id := range batchIDs {
		positions[id.String()] = i
	}

	buildVoteMask := func(vote *BatchedWiVote) (*big.Int, error) {
		if len(vote.WithdrawalReqIDs) < n.cfg.BatchSize || len(vote.Vote) < n.cfg.BatchSize {
			return nil, fmt.Errorf("validator %d batch vote missing: need %d requests", vote.Index.Uint64(), n.cfg.BatchSize)
		}
		mask := big.NewInt(0)
		for j := 0; j < n.cfg.BatchSize; j++ {
			reqID := vote.WithdrawalReqIDs[j]
			pos, ok := positions[reqID.String()]
			if !ok {
				return nil, fmt.Errorf("validator %d request mismatch: requestID=%s", vote.Index.Uint64(), reqID.String())
			}
			if vote.Vote[j] != nil && vote.Vote[j].Sign() != 0 {
				mask.SetBit(mask, pos, 1)
			}
		}
		return mask, nil
	}

	voteMasks := make([]*big.Int, len(voteSlice))
	voteCounts := make(map[string]int, len(voteSlice))
	for i, vote := range voteSlice {
		mask, err := buildVoteMask(vote)
		if err != nil {
			return nil, err
		}
		voteMasks[i] = mask
		voteCounts[mask.String()]++
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

	if majorityCount < threshold {
		return nil, fmt.Errorf("no BFT quorum for vote mask: max=%d threshold=%d", majorityCount, threshold)
	}

	hfunc := hash.MIMC_BN254.New()
	hfunc.Reset()
	var commitmentElem fr.Element
	for i := 0; i < len(batchIDs); i++ {
		commitmentElem.SetBigInt(batchIDs[i])
		hfunc.Write(commitmentElem.Marshal()[:])
	}
	batchCommitment := hfunc.Sum(nil)
	batchCommitment = util.PadOrTrim(util.ModToBn254Bytes(batchCommitment), 32)

	hfunc.Reset()
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

	validatorConstraints := make([]votingbatch.BatchingValidatorConstraints, n.cfg.NodeCount)

	validatorBits := new(big.Int)
	honestBits := new(big.Int) //tracks validators matching majority

	fmt.Printf("Votes length--------- %x\n", len(voteSlice))

	// var tempPostState []byte
	for i, vote := range voteSlice {

		// 1. Find the node who signed this vote
		node, err := n.nodeByAccountIndex(vote.Index)
		if err != nil {
			return nil, err
		}

		// 2. Rebuild the original message that was signed
		// hfunc.Reset()
		// hfunc.Write(vote.Index.Bytes())
		// hfunc.Write(batchCommitment)
		// hfunc.Write(big.NewInt(31).Bytes()) // TODO: fixme
		// hfunc.Write(new(big.Int).SetInt64(int64(n.Oracle.RoundID)).Bytes())
		// msg := hfunc.Sum(nil)

		validatorVoteMask := new(big.Int).Set(voteMasks[i])

		msg := hashBatchVoteFieldwise(
			new(big.Int).Set(vote.Index),
			new(big.Int).SetBytes(batchCommitment[:32]),
			new(big.Int).Set(validatorVoteMask),
			new(big.Int).Set(roundID),
		)

		// 3. Re-sign the vote (this can be for validation or to force resync)
		vote.Signature, err = node.privateKey.Sign(msg, hfunc)
		if err != nil {
			return nil, fmt.Errorf("sign the vote: %w", err)
		}

		// 4. Read validator account from state
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

func (n *Node) updateLatestIPFSHashTx(latestIPFSHash string) error {
	fmt.Printf("updating latest ipfs hash (validator=%v) ...\n", n.ID)
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

	tx, err := bcClient.UpdateLatestIPFSHash(trxOpts, latestIPFSHash)
	if err != nil {
		log.Fatalf("call updateLatestIPFSHash() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	bc.LogTxReceipt("update latest ipfs hash", tx, receipt)
	if receipt.Status == 1 {
		fmt.Printf("successfully updated ipfs hash (node=%v)\n", n.ID)
	} else {
		fmt.Printf("Transaction failed (node=%v\n)", n.ID)
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
	bc.LogTxReceipt("submit wivote", tx, receipt)

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

func (n *Node) ReplaceAccountTx(replaceWithAccountID uint64) error {
	fmt.Printf("starting to replace account (node=%v)...\n", n.ID)

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
		log.Fatalf("replace tx reverted (index=%d)\n", account.Index.Uint64())
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

func (n *Node) ExitTx() error {
	fmt.Printf("exiting account (node=%v)...\n", n.ID)

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
		fmt.Sprintf("exit account index=%d", account.Index.Uint64()),
		tx,
	)
	if err != nil {
		log.Fatalf("failed to wait for exit tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status != 1 {
		log.Fatalf("exit tx reverted (index=%d)\n", account.Index.Uint64())
	}

	return nil
}

func (n *Node) WithdrawAccountTx() error {
	fmt.Printf("starting withdraw account (node=%v)...\n", n.ID)

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
		fmt.Sprintf("withdraw account index=%d", account.Index.Uint64()),
		tx,
	)
	if err != nil {
		log.Fatalf("failed to wait for withdraw tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status == 1 {
		log.Printf("withdrawn account index=%d", account.Index.Uint64())
	} else {
		log.Fatalf("withdraw tx reverted (index=%d)\n", account.Index.Uint64())
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
func bftThreshold(nValidators int) int {
	if nValidators <= 0 {
		return 0
	}
	f := (nValidators - 1) / 3
	return 2*f + 1
}
