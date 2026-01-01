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
						for uniqueReqID := range n.WiVotes {
							// aggregator wiVote process:
							err := n.AggregatorProcessWiVote(uniqueReqID)
							if err != nil {
								panic(fmt.Errorf("failed to process wiVote: %v", err))
							}

						}
					}

					// reset wivotes:
					n.aggregatorResetWiVotes()
				}
				//$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
				if e.Title == MessageTittleBatchedWiVote && n.IsAggregator() {
					fmt.Printf("aggregator (node=%d) is performing batch wiVote... (len of batched votes=%d)\n", n.ID, len(n.BatchedWiVote.WithdrawalReqIDs))

					if len(n.BatchedWiVote.WithdrawalReqIDs) < votingbatch.BatchSize {
						fmt.Printf("aggregator (node=%d) skipping batch wiVote: need %d requests, have %d\n", n.ID, votingbatch.BatchSize, len(n.BatchedWiVote.WithdrawalReqIDs))
						continue
					}
					majorityVote, err := n.processBatchedWiVotes()
					if err != nil {
						fmt.Printf("aggregator (node=%d) failed to process batched wiVotes: %v\n", n.ID, err)
						continue
					}
					fmt.Printf("majority vote result ", majorityVote.String())

					// fmt.Printf("aggregator (node=%d) is performing batch slashing...\n", n.ID)
					// err = n.processBatchedSlashing(majorityVote)
					// if err != nil {
					// 	panic(fmt.Errorf("failed to slash: %v", err))
					// }

					n.resetBatchedWiVotes()
				}
				//$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
				if e.Title == MessageTittleSelectIncVote && n.IsAggregator() {
					fmt.Printf("aggregator (node %d) is performing incVote selection...\n", n.ID)
					fmt.Printf("incvotes: %v\n", n.IncVotes)
					selected := make([]*IncVote, 0)
					for commitmentHashStr := range n.IncVotes {
						selectedIncVote, err := n.AggregatorSelectVote(commitmentHashStr)
						if err != nil {
							fmt.Printf("failed to select incVote: %v\n", err)
							continue
						}
						fmt.Printf("selected tree index is: %v\n", selectedIncVote.IncTreeIndex)
						selected = append(selected, selectedIncVote)
					}
					if len(selected) != 0 {
						// fetch and update ipfs:
						if n.IPFSClient.LatestHash != "" {
							fmt.Printf("fetching ipfs content... (node_id=%v)\n", n.ID)
							n.IPFSContent = n.FetchIPFS()
						}
						for _, incVote := range selected {
							n.IPFSContent.CommitmentHashIncVote[string(incVote.CommitmentHash)] = *incVote
						}
						n.IPFSContent.IncMerkleTree = *n.IncMerkleTree
						fmt.Printf("updating ipfs content... (node_id=%v)\n", n.ID)
						latestIPFSHash, _ := n.UpdateIPFS()
						n.updateLatestIPFSHashTx(latestIPFSHash)
					}
					n.aggregatorResetIncVotes()
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
	err = groth16.Verify(proof, n.Oracle.IncVK, publicWitness)

	if err != nil {
		fmt.Printf("failed to verify claim proof (node=%v): %v\n", n.ID, err)
		isApproved = big.NewInt(0)
	} else {
		fmt.Printf("successfully verified claim merkle proof and MiMC hash! (node=%v)\n", n.ID)
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

func (n *Node) resetBatchedWiVotes() {
	for _, node := range n.Oracle.Nodes {
		node.BatchedWiVote.WithdrawalReqIDs = make([]*big.Int, 0)
		node.BatchedWiVote.Vote = make([]*big.Int, 0)
		node.BatchedWiVote.MajorityOfVotes = make([]*big.Int, 0)
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
		if selection[incTreeIndexAndRootHash] > len(n.Oracle.Nodes)/2 {
			selectedIncVote = incVote
			break
		}
	}

	return selectedIncVote, nil
}

func (n *Node) AggregatorSelectWiVote(uniqueReqID string) (isApproved *big.Int, err error) {
	isApproved = big.NewInt(0)
	approvedCount := 0
	for _, wiVote := range n.WiVotes[uniqueReqID] {
		if wiVote.IsApproved.Uint64() == 1 {
			approvedCount++
			if approvedCount > n.cfg.NodeCount/2 {
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
		node := n.Oracle.Nodes[uint(v.Index.Uint64())]
		fmt.Println("adding wivote to the bached wivote of the node:", node.ID)
		node.BatchedWiVote.Index = v.Index
		node.BatchedWiVote.WithdrawalReqIDs = append(node.BatchedWiVote.WithdrawalReqIDs, v.RequestID)
		node.BatchedWiVote.Vote = append(node.BatchedWiVote.Vote, v.IsApproved)
		node.BatchedWiVote.MajorityOfVotes = append(node.BatchedWiVote.MajorityOfVotes, majorityWiVote)
		node.BatchedWiVote.SenderPK = n.privateKey.PublicKey
	}

	return nil
}

func (n *Node) processBatchedWiVotes() (*big.Int, error) {

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

	majorityVote := big.NewInt(31) // TODO: fixme take majority vote from the batch wivotes

	fmt.Printf("selected wivotes bathced: %v\n", voteSlice)

	preStateData := make([]byte, len(n.state.Data))
	copy(preStateData, n.state.Data)

	preStateHData := make([]byte, len(n.state.HData))
	copy(preStateHData, n.state.HData)

	hfunc := hash.MIMC_BN254.New()

	sort.Slice(n.BatchedWiVote.WithdrawalReqIDs, func(i, j int) bool {
		return n.BatchedWiVote.WithdrawalReqIDs[i].Cmp(n.BatchedWiVote.WithdrawalReqIDs[j]) < 0
	})

	hfunc.Reset()
	for i := 0; i < len(n.BatchedWiVote.WithdrawalReqIDs); i++ {
		hfunc.Write(n.BatchedWiVote.WithdrawalReqIDs[i].Bytes())
	}
	batchCommitment := hfunc.Sum(nil)
	batchCommitment = util.ModToBn254Bytes(batchCommitment)[:32]

	hfunc.Reset()
	fmt.Printf("Length of WithdrawalReqIDs: %d\n", len(n.BatchedWiVote.WithdrawalReqIDs))
	if len(n.BatchedWiVote.WithdrawalReqIDs) != votingbatch.BatchSize {
		return nil, fmt.Errorf("batch size mismatch: expected %d got %d", votingbatch.BatchSize, len(n.BatchedWiVote.WithdrawalReqIDs))
	}

	fmt.Println("Sorted WithdrawalReqIDs:")
	for _, id := range n.BatchedWiVote.WithdrawalReqIDs {
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

	var validatorConstraints [votingbatch.NumValidators]votingbatch.BatchingValidatorConstraints

	validatorBits := new(big.Int)
	honestBits := new(big.Int) //tracks validators matching majority

	fmt.Printf("Votes length--------- %x\n", len(voteSlice))

	// var tempPostState []byte
	for i, vote := range voteSlice {

		// 1. Find the node who signed this vote
		node := n.Oracle.Nodes[uint(vote.Index.Uint64())]

		// 2. Rebuild the original message that was signed
		// hfunc.Reset()
		// hfunc.Write(vote.Index.Bytes())
		// hfunc.Write(batchCommitment)
		// hfunc.Write(big.NewInt(31).Bytes()) // TODO: fixme
		// hfunc.Write(new(big.Int).SetInt64(int64(n.Oracle.RoundID)).Bytes())
		// msg := hfunc.Sum(nil)

		msg := hashBatchVoteFieldwise(
			new(big.Int).Set(vote.Index),
			new(big.Int).SetBytes(batchCommitment[:32]),
			big.NewInt(31), // or majorityVote //TODO: fixme
			big.NewInt(int64(n.Oracle.RoundID)),
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

		// make sure use the node public key
		publicKey.Assign(tedwards.BN254, node.privateKey.PublicKey.Bytes())
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
			//Vote:      big.NewInt(31), // TODO: fixme
			Vote: big.NewInt(31), // TODO: fixme based on the majority vote it should be 1

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

		// ------------------------------------------------------------------
		// TEMP FIX (batch voting):
		//   - BatchSize = 5
		//   - Correct vote  = 31 (0b11111)
		//   - Bad vote      = 30 (0b11110)
		//   - MajorityVote  = 31
		// TODO: Replace with real per-validator bitmask comparison later
		// ------------------------------------------------------------------

		// validatorVote := validatorConstraints[i].Vote // this is the value the circuit uses (31)
		validatorVote := big.NewInt(31) // TEMP: validator vote (31 = honest, 30 = bad)
		majorityVote := big.NewInt(31)  // TEMP: majority vote (31 = honest, 30 = bad)

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
			// dishonest validator → FULL SLASH
			validatorAccount.Balance.SetInt64(0)
		}

		err = n.state.WriteAccount(validatorAccount)
		if err != nil {
			return nil, fmt.Errorf("write account: %w", err)
		}

		// *****************************************************************************
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

	withdrawalReqIDs := make([]frontend.Variable, len(n.BatchedWiVote.WithdrawalReqIDs))
	for i := range n.BatchedWiVote.WithdrawalReqIDs {
		withdrawalReqIDs[i] = frontend.Variable(n.BatchedWiVote.WithdrawalReqIDs[i])
	}

	assignment := votingbatch.BatchingVotingCircuit{
		ResultingStateRoot: postStateRoot,
		//RoundID:            n.Oracle.RoundID,
		RoundID:          new(big.Int).SetInt64(int64(n.Oracle.RoundID)),
		BatchCommitment:  batchCommitment[:32],
		MajorityVote:     new(big.Int).Set(majorityVote),
		ValidatorBits:    new(big.Int).Set(validatorBits),
		HonestBits:       new(big.Int).Set(honestBits), // NEW
		WithdrawalReqIDs: withdrawalReqIDs,
		Aggregator:       aggregatorConstraints,
		Validators:       validatorConstraints,
	}

	// Print all variables used in the assignment
	fmt.Printf("ResultingStateRoot: %v\n", postStateRoot)
	fmt.Printf("RoundID: %v\n", new(big.Int).SetInt64(int64(n.Oracle.RoundID)))
	fmt.Printf("BatchCommitment: %v\n", batchCommitment[:32])
	fmt.Printf("MajorityVote: %v\n", new(big.Int).Set(majorityVote))
	fmt.Printf("ValidatorBits: %v\n", new(big.Int).Set(validatorBits))
	fmt.Printf("WithdrawalReqIDs: %v\n", withdrawalReqIDs)
	fmt.Printf("Aggregator: %+v\n", aggregatorConstraints)
	fmt.Printf("Validators: %+v\n", validatorConstraints)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("create witness: %w", err)
	}

	p, err := groth16.Prove(n.Oracle.SparseR1CS, n.Oracle.SparsePK, witness)
	if err != nil {
		return nil, fmt.Errorf("prove: %v", err)
	}

	pw, err := witness.Public()
	if err != nil {
		return nil, fmt.Errorf("public witness: %w", err)
	}

	err = groth16.Verify(p, n.Oracle.SparseVK, pw) // TODO: delete maybe later
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
	check("slashedValIndex", new(big.Int).SetInt64(int64(n.Oracle.RoundID)))
	check("uniqueReqID", new(big.Int).SetInt64(int64(n.Oracle.RoundID)))
	check("postStateRoot", new(big.Int).SetBytes(postStateRoot))
	check("postSeedX", postSeedX)
	check("postSeedY", postSeedY)
	check("batchCommitment", new(big.Int).SetBytes(batchCommitment))
	check("majorityVote", majorityVote)
	check("validatorBits", validatorBits)

	n.aggregatorSubmitWiVoteTx(
		n.Account.Index,
		new(big.Int).SetInt64(int64(n.Oracle.RoundID)), // Pass the first element of the slice as an example
		new(big.Int).SetBytes(batchCommitment[:32]),
		validatorBits,
		honestBits, // NEW
		new(big.Int).Set(majorityVote),
		new(big.Int).SetBytes(postStateRoot),
		postSeedX,
		postSeedY,
		proof.Proof,
	)

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
	reason, err := abi.UnpackRevert(data)
	if err != nil {
		return "", ""
	}
	return reason, ""
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
		fmt.Printf("Transaction failed (node=%v\n)", n.ID)
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
	// trxOpts.GasPrice = big.NewInt(20000000000)
	trxOpts.Value = big.NewInt(account.Balance.Int64())

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))

	_, path, err := n.state.MerkleProofBytes(account.Index.Uint64())
	if err != nil {
		log.Fatalf("merkle proof failed (index=%d): %v", account.Index.Uint64(), err)
	}
	fmt.Printf("*** About to replace account:\n")
	fmt.Printf("  Index         : %x\n", replaceAccount.Index.String())
	fmt.Printf("  PublicKey.X   : %x\n", replaceAccount.PublicKey.A.X.String())
	fmt.Printf("  PublicKey.Y   : %x\n", replaceAccount.PublicKey.A.Y.String())
	fmt.Printf("  Balance       : %x\n", replaceAccount.Balance.String())

	tx, err := bcClient.Replace(
		trxOpts,
		*gnark.PublicKeyToOraclePublicKey(account.PublicKey),
		*gnark.AccountToOracleAccount(&replaceAccount),
		path[:], // Convert fixed-size array to slice
		new(big.Int).Set(account.Index),
		big.NewInt(int64(n.cfg.SparseTreeDepth)),
	)
	if err != nil {
		log.Fatalf("replace failed (index=%d): %v", account.Index.Uint64(), err)
	}

	// 5. Wait for confirmation
	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for replace tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status == 1 {
		bc.LogTxReceipt(fmt.Sprintf("replace account index=%d", account.Index.Uint64()), tx, receipt)
	} else {
		log.Fatalf("replace tx reverted (index=%d)\n", account.Index.Uint64())
	}

	// 6. Update local state
	account.Balance = big.NewInt(0)
	if err := n.state.WriteAccount(account); err != nil {
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
	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for exit tx (index=%d): %v", account.Index.Uint64(), err)
	}

	if receipt.Status == 1 {
		bc.LogTxReceipt(fmt.Sprintf("exit account index=%d", account.Index.Uint64()), tx, receipt)
	} else {
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
	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for withdraw tx (index=%d): %v", account.Index.Uint64(), err)
	}
	bc.LogTxReceipt(fmt.Sprintf("withdraw account index=%d", account.Index.Uint64()), tx, receipt)

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
