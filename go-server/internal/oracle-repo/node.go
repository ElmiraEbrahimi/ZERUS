package oracle

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"

	"l2alchemy/internal/config"
	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/merkle"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	stateSync       *gnark.StateSync

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
	n := &Node{cfg: cfg, ethClient: ethClient, Oracle: oracle, IPFSClient: ipfsClient, ID: id, Role: Validator}

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

	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}

	fmt.Printf("setting stateSync for node (=%v)...\n", n.ID)
	n.stateSync = gnark.NewStateSync(cfg, uint64(id), state, ethClient, bcClient)

	ecdsaPrivateKey, err := crypto.HexToECDSA(n.cfg.NodesPK[n.ID])
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

func (n *Node) listenToBlockchainEvents(quit chan struct{}) {
	contractAddr := common.HexToAddress(n.cfg.OracleContractAddress)

	// filterers:
	oracleFilterer, err := bc.NewOracleFilterer(contractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create oracle filterer: %v", err)
	}
	newAggEventChan := make(chan *bc.OracleNewAggregator)
	accountRegEventChan := make(chan *bc.OracleValidatorRegistered)
	burnEventChan := make(chan *bc.OracleBurnSubmitted)
	claimEventChan := make(chan *bc.OracleClaimSubmitted)

	// event subs:
	ctx := context.Background()
	newAggSub, err := oracleFilterer.WatchNewAggregator(&bind.WatchOpts{Context: ctx}, newAggEventChan)
	if err != nil {
		log.Fatalf("set up new aggregator event subscription: %v", err)
	}
	accountRegSub, err := oracleFilterer.WatchValidatorRegistered(&bind.WatchOpts{Context: ctx}, accountRegEventChan)
	if err != nil {
		log.Fatalf("set up account register event subscription: %v", err)
	}
	burnSub, err := oracleFilterer.WatchBurnSubmitted(&bind.WatchOpts{Context: ctx}, burnEventChan)
	if err != nil {
		log.Fatalf("set up burn event subscription: %v", err)
	}
	claimSub, err := oracleFilterer.WatchClaimSubmitted(&bind.WatchOpts{Context: ctx}, claimEventChan)
	if err != nil {
		log.Fatalf("set up claim event subscription: %v", err)
	}

	n.wg.Add(1)
	go func() {
		for {
			select {
			// new aggregator event:
			case newAggEvent := <-newAggEventChan:
				// fmt.Printf("node %d received new aggregator event, new aggregator ID=%v\n", n.ID, newAggEvent.Arg0)
				newAggID := newAggEvent.ValidatorID.Uint64()
				if uint(newAggID) == n.ID {
					fmt.Printf("node %d says: I am the new aggregator\n", n.ID)
					n.aggregatorResetIncVotes()
					n.aggregatorResetWiVotes()
					n.resetBatchedWiVotes()
					n.Role = Aggregator
					n.Oracle.AggregatorID = n.ID
				} else {
					n.Role = Validator
				}

			case accountRegEvent := <-accountRegEventChan:
				if accountRegEvent.ValidatorID.Uint64() == uint64(n.ID) {
					n.Account.Index = accountRegEvent.Index
					n.Account.Balance = accountRegEvent.Balance
					n.Account.Reputation = accountRegEvent.Reputation
					n.Account.SeverityCount = accountRegEvent.SeverityCount
					fmt.Printf("node=%v registered: %v\n", n.ID, n.Account)
				}

			// burn event:
			case burnEvent := <-burnEventChan:
				// fmt.Printf("node %d received burn event (commitmentHash=%v)...\n", n.ID, burnEvent.CommitmentHash)
				fmt.Printf("node %d received burn event\n", n.ID)
				eventCommitmentHash := burnEvent.CommitmentHash[:]
				// incTreeIndex, commitmentHash, rootHash, err := n.MerkleTree.AddLeafValidator(eventCommitmentHash)
				incTreeIndex, _, rootHash, err := n.IncMerkleTree.AddLeafValidator(eventCommitmentHash)
				if err != nil {
					fmt.Printf("node %d failed to add leaf to merkle tree: %v\n", n.ID, err)
					break
				}
				incVote := &IncVote{
					CommitmentHash: eventCommitmentHash,
					NodeID:         n.ID,
					IncTreeIndex:   incTreeIndex,
					RootHash:       rootHash,
				}
				fmt.Printf("node %d publishing incVote for burn event (incTreeIndex=%v)...\n", n.ID, incVote.IncTreeIndex)
				n.Oracle.PublishIncVote(incVote)

			// claim event:
			case claimEvent := <-claimEventChan:
				fmt.Printf("node %d received claim event\n", n.ID)

				wiVote, err := n.VerifyClaim(claimEvent)
				if err != nil {
					log.Fatalf("failed verifying claim: %v\n", err)
				}
				n.Oracle.PublishWiVote(wiVote)

			// sub errors:
			case err := <-newAggSub.Err():
				log.Fatalf("burn sub error: %v\n", err)
			case err := <-accountRegSub.Err():
				log.Fatalf("account reg sub error: %v\n", err)
			case err := <-burnSub.Err():
				log.Fatalf("burn sub error: %v\n", err)
			case err := <-claimSub.Err():
				log.Fatalf("claim sub error: %v\n", err)
			// quit:
			case <-quit:
				fmt.Printf("blockchain listener of node %v quitting...\n", n.ID)
				newAggSub.Unsubscribe()
				accountRegSub.Unsubscribe()
				burnSub.Unsubscribe()
				claimSub.Unsubscribe()
				close(newAggEventChan)
				close(accountRegEventChan)
				close(burnEventChan)
				close(claimEventChan)
				n.wg.Done()
				return
			}
		}
	}()
}

func (n *Node) Start() {
	// listen to blockchain events:
	go func() {
		if err := n.stateSync.Synchronize(); err != nil {
			panic(fmt.Errorf("failed to sync state: %v", err))
		}
	}()

	// register validator:
	n.registerValidatorTx()

	// n.wg.Add(1)
	// go func() {
	// 	for e := range n.OracleMessages {

	// 		// message from oracle:
	// 		if e.From == OracleID {
	// 			if e.Title == MessageTitleTerminate {
	// 				fmt.Printf("node %d terminating...\n", n.ID)
	// 				/////////////////////////////////////////////////////
	// 				// if n.ID == 0 { // or if n.IsAggregator() or totalValidators-1
	// 				// //  or with setting in config n.ShouldWithdrawOnExit

	// 				// 	err := n.WithDrawAccounts([]*gnark.Account{n.Account}, n.state)
	// 				// 	if err != nil {
	// 				// 		fmt.Printf("node %d failed to withdraw: %v\n", n.ID, err)
	// 				// 	} else {
	// 				// 		fmt.Printf("node %d successfully withdrew before termination.\n", n.ID)
	// 				// 	}
	// 				// }
	// 				// /////////////////////////////////////////////////////
	// 				n.wg.Done()
	// 				return
	// 			}
	// 			if e.Title == MessageTittleSelectWiVote && n.IsAggregator() {
	// 				fmt.Printf("aggregator (node=%d) is performing wiVote selection...\n", n.ID)
	// 				fmt.Printf("wivotes: %v\n", n.WiVotes)
	// 				if len(n.WiVotes) != 0 {
	// 					for uniqueReqID := range n.WiVotes {
	// 						// aggregator wiVote process:
	// 						err := n.AggregatorProcessWiVote(uniqueReqID)
	// 						if err != nil {
	// 							panic(fmt.Errorf("failed to process wiVote: %v", err))
	// 						}

	// 					}
	// 				}

	// 				// reset wivotes:
	// 				n.aggregatorResetWiVotes()
	// 			}

	// 			if e.Title == MessageTittleSelectIncVote && n.IsAggregator() {
	// 				fmt.Printf("aggregator (node %d) is performing incVote selection...\n", n.ID)
	// 				fmt.Printf("incvotes: %v\n", n.IncVotes)
	// 				selected := make([]*IncVote, 0)
	// 				for commitmentHashStr := range n.IncVotes {
	// 					selectedIncVote, err := n.AggregatorSelectVote(commitmentHashStr)
	// 					if err != nil {
	// 						fmt.Printf("failed to select incVote: %v\n", err)
	// 						continue
	// 					}
	// 					fmt.Printf("selected tree index is: %v\n", selectedIncVote.IncTreeIndex)
	// 					selected = append(selected, selectedIncVote)
	// 				}
	// 				if len(selected) != 0 {
	// 					// fetch and update ipfs:
	// 					if n.IPFSClient.LatestHash != "" {
	// 						fmt.Printf("fetching ipfs content... (node_id=%v)\n", n.ID)
	// 						n.IPFSContent = n.FetchIPFS()
	// 					}
	// 					for _, incVote := range selected {
	// 						n.IPFSContent.CommitmentHashIncVote[string(incVote.CommitmentHash)] = *incVote
	// 					}
	// 					n.IPFSContent.IncMerkleTree = *n.IncMerkleTree
	// 					fmt.Printf("updating ipfs content... (node_id=%v)\n", n.ID)
	// 					latestIPFSHash, _ := n.UpdateIPFS()
	// 					n.updateLatestIPFSHashTx(latestIPFSHash)
	// 				}
	// 				n.aggregatorResetIncVotes()
	// 			}
	// 		}

	// 		// message from other nodes:
	// 		if e.Title == MessageTittleWiVote && n.IsAggregator() {
	// 			fmt.Printf("aggregator (node %d) is collecting wiVote...\n", n.ID)
	// 			wiVote := e.Message.(*WiVote)
	// 			n.aggregatorCollectWiVote(wiVote)
	// 		}
	// 		if e.Title == MessageTittleIncVote && n.IsAggregator() {
	// 			fmt.Printf("aggregator (node %d) is collecting incVote...\n", n.ID)
	// 			incVote := e.Message.(*IncVote)
	// 			n.aggregatorCollectIncVote(incVote)
	// 		}
	// 	}
	// }()
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

// endregion

// region contract

func (n *Node) selectNewAggregatorTx() error {
	fmt.Printf("selecting new aggregator from contract...\n")
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	tx, err := bcClient.ChooseNewAggregator(trxOpts)
	if err != nil {
		log.Fatalf("call SelectNewAggregatorTx() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	if receipt.Status == 1 {
		fmt.Printf("successfully called SelectNewAggregatorTx (by oracle)\n")
	} else {
		fmt.Printf("Transaction failed (by oracle)\n")
	}

	return nil
}

func (n *Node) registerValidatorTx() error {
	fmt.Printf("registering validator=%v ...\n", n.ID)
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	pk := gnark.PublicKeyToOraclePublicKey(n.Account.PublicKey)
	tx, err := bcClient.RegisterValidator(trxOpts, big.NewInt(int64(n.ID)), *pk)
	if err != nil {
		log.Fatalf("call RegisterValidator() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	if receipt.Status == 1 {
		fmt.Printf("successfully registered validator node=%v\n", n.ID)
	} else {
		fmt.Printf("Transaction failed (node=%v\n)", n.ID)
	}

	return nil
}

func (n *Node) updateLatestIPFSHashTx(latestIPFSHash string) error {
	fmt.Printf("updating latest ipfs hash (validator=%v) ...\n", n.ID)
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	tx, err := bcClient.UpdateLatestIPFSHash(trxOpts, latestIPFSHash)
	if err != nil {
		log.Fatalf("call updateLatestIPFSHash() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	if receipt.Status == 1 {
		fmt.Printf("successfully updated ipfs hash (node=%v)\n", n.ID)
	} else {
		fmt.Printf("Transaction failed (node=%v\n)", n.ID)
	}

	return nil
}

func (n *Node) aggregatorSubmitWiVoteTx(index *big.Int, uniqueReqID *big.Int, batchCommitment *big.Int, validatorBits *big.Int, vote *big.Int, postStateRoot *big.Int, postSeedX *big.Int, postSeedY *big.Int, proof [8]*big.Int) error {
	fmt.Printf("submitting wivote (node=%v) ...\n", n.ID)
	oracleContractAddr := common.HexToAddress(n.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, n.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	tx, err := bcClient.SubmitWiVote(
		trxOpts,
		index,
		uniqueReqID,
		batchCommitment,
		validatorBits,
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
	chainID := big.NewInt(n.cfg.ChainID)

	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
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
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	_, path, err := n.state.MerkleProofBytes(account.Index.Uint64())
	if err != nil {
		log.Fatalf("merkle proof failed (index=%d): %v", account.Index.Uint64(), err)
	}
	fmt.Printf("*** About to replace account:\n")
	fmt.Printf("  Index         : %x\n", replaceAccount.Index.String())
	fmt.Printf("  PublicKey.X   : %x\n", replaceAccount.PublicKey.A.X.String())
	fmt.Printf("  PublicKey.Y   : %x\n", replaceAccount.PublicKey.A.Y.String())
	fmt.Printf("  Balance       : %x\n", replaceAccount.Balance.String())
	fmt.Printf("  Reputation    : %x\n", replaceAccount.Reputation.String())
	fmt.Printf("  SeverityCount : %x\n", replaceAccount.SeverityCount.String())

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
		log.Printf("replace: account index=%d | gas=%d | tx=%s\n", account.Index.Uint64(), receipt.GasUsed, tx.Hash().Hex())
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
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

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
		log.Printf("exit: account index=%d | gas=%d | tx=%s\n", account.Index.Uint64(), receipt.GasUsed, tx.Hash().Hex())
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

	chainID := big.NewInt(n.cfg.ChainID)

	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}

	pendingNonce, err := n.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := n.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

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
	log.Printf("✅ Account index=%d | balance=%s | Withdraw gas used: %d\n", account.Index.Uint64(), account.Balance.String(), receipt.GasUsed)
	log.Printf(" Saving to Oracle.GasCosts.WithdrawCost = %d\n", receipt.GasUsed)
	log.Printf("DEBUG: Withdrawal receipt status: %d\n", receipt.Status)

	if receipt.Status == 1 {
		log.Printf("Withdrawn: account index=%d | gas=%d | tx=%s\n", account.Index.Uint64(), receipt.GasUsed, tx.Hash().Hex())
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
	chainID := big.NewInt(n.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, chainID)
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
