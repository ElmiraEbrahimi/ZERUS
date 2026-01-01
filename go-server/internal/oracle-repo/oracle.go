package oracle

import (
	"encoding/csv"
	"fmt"
	merkleproof "l2alchemy/circuits/merkle_proof"
	votingbatch "l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/config"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Oracle struct {
	cfg          *config.Config
	wg           *sync.WaitGroup
	ethClient    *ethclient.Client
	Nodes        map[uint]*Node `json:"nodes"`
	AggregatorID uint           `json:"aggregator_id"`

	IncCircuit *merkleproof.MerkleProofCircuit
	IncR1CS    constraint.ConstraintSystem
	IncVK      groth16.VerifyingKey

	SparseCircuit *votingbatch.BatchingVotingCircuit
	SparseR1CS    constraint.ConstraintSystem
	SparsePK      groth16.ProvingKey
	SparseVK      groth16.VerifyingKey

	Oracle *ethclient.Client

	messageLock sync.Mutex

	GasCosts *GasCosts

	CircuitMemTime         *CircuitMemTime
	SparseMemTimeCSVWriter *csv.Writer

	RoundID int
}

type GasCosts struct {
	RegisterValidatorCost uint64
	SubmitWiVoteCost      uint64
	WithdrawCost          uint64

	ReplaceCost          uint64
	ExitCost             uint64
	UpdateLatestIPFSHash uint64
}

type CircuitMemTime struct {
	SparseProvingTime        int
	SparseProvingMemoryUsage int

	SparseCompileMemory int
	SparseCompileTime   int
}

type InternalOraclelMessage struct {
	From    uint   `json:"from"`
	Title   string `json:"title"`
	Message any    `json:"message"`
}

const (
	OracleID                   = 7777
	MessageTitleTerminate      = "terminate"
	MessageTittleIncVote       = "inc_vote"
	MessageTittleWiVote        = "wi_vote"
	MessageTittleSelectIncVote = "inc_vote_selection"
	MessageTittleSelectWiVote  = "wivote_selection"
	MessageTittleBatchedWiVote = "batched_wivote_selection"
)

func NewOracle(
	cfg *config.Config,
	ethClient *ethclient.Client,
	ipfsClient *db.IPFSClient,
	nodesCount int,

	incCircuit *merkleproof.MerkleProofCircuit,
	incR1cs constraint.ConstraintSystem,
	incVk groth16.VerifyingKey,

	sparseCircuit *votingbatch.BatchingVotingCircuit,
	sparseR1cs constraint.ConstraintSystem,
	sparsePK groth16.ProvingKey,
	sparseVK groth16.VerifyingKey,

	privteKeys []*eddsa.PrivateKey,
	accounts []*gnark.Account,

) *Oracle {

	o := new(Oracle)
	o.cfg = cfg
	o.ethClient = ethClient

	o.IncCircuit = incCircuit
	o.IncR1CS = incR1cs
	o.IncVK = incVk

	o.SparseCircuit = sparseCircuit
	o.SparseR1CS = sparseR1cs
	o.SparsePK = sparsePK
	o.SparseVK = sparseVK

	fmt.Println("creating nodes...")
	nodes := make(map[uint]*Node, nodesCount)
	for i := 0; i < nodesCount; i++ {
		offeredAmount := 1000 // TODO: change initial validator balance / move to cfg
		nodes[uint(i)] = NewNode(cfg, ethClient, ipfsClient, o, uint(i), privteKeys[i], uint(offeredAmount), Validator, accounts)
	}
	o.Nodes = nodes
	o.AggregatorID = 0

	return o
}

func (o *Oracle) Start() {
	for _, node := range o.Nodes {
		node.Start()
	}
}

func (o *Oracle) SelectNewAggregator() error {
	aggregator := o.Nodes[o.AggregatorID]
	err := aggregator.selectNewAggregatorTx()

	return err
}

// region internal messages

func (o *Oracle) PublishInternalMessage(msg InternalOraclelMessage) {
	for _, node := range o.Nodes {
		node.OracleMessages <- msg
	}
}

func (o *Oracle) PublishIncVote(incVote *IncVote) {
	o.messageLock.Lock()
	defer o.messageLock.Unlock()
	o.PublishInternalMessage(
		InternalOraclelMessage{From: incVote.NodeID, Title: MessageTittleIncVote, Message: incVote},
	)
}
func (o *Oracle) PublishWiVote(wiVote *WiVote) {
	o.messageLock.Lock()
	defer o.messageLock.Unlock()
	o.PublishInternalMessage(
		InternalOraclelMessage{From: wiVote.NodeID, Title: MessageTittleWiVote, Message: wiVote},
	)
}

func (o *Oracle) PublishIncVoteSelectionRes() {
	o.PublishInternalMessage(
		InternalOraclelMessage{From: OracleID, Title: MessageTittleSelectIncVote, Message: ""},
	)
}

func (o *Oracle) PublishWiVoteSelectionRes() {
	o.PublishInternalMessage(
		InternalOraclelMessage{From: OracleID, Title: MessageTittleSelectWiVote, Message: ""},
	)
}

func (o *Oracle) PublishBatchedWiVoteRes() {
	o.PublishInternalMessage(
		InternalOraclelMessage{From: OracleID, Title: MessageTittleBatchedWiVote, Message: ""},
	)
}

func (o *Oracle) PublishTerminate() {
	o.PublishInternalMessage(
		InternalOraclelMessage{From: OracleID, Title: MessageTitleTerminate, Message: ""},
	)
}

// endregion
