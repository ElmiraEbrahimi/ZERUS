package events

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"l2alchemy/internal/eth"
	"l2alchemy/internal/oracle-repo"
	"l2alchemy/internal/oracle_runtime"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const l2LogsChanBuffer = 256
const l2PollInterval = 1 * time.Second

// L2ContractEventSubscriber subscribes to L2 contract logs and keeps a handle to the
// oracle engine so events can be applied to the runtime state.
//
// This is intentionally L2-scoped so a future L1 subscriber can coexist
// independently.
type L2ContractEventSubscriber struct {
	engine *oracle_runtime.OracleEngine

	contractAddr common.Address
	contractABI  abi.ABI
	filterer     *eth.OracleFilterer

	mu        sync.Mutex
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	readyOnce sync.Once
	readyCh   chan struct{}
}

// NewL2ContractEventSubscriber constructs a subscriber bound to the configured L2 Oracle
// contract address.
func NewL2ContractEventSubscriber(engine *oracle_runtime.OracleEngine) (*L2ContractEventSubscriber, error) {
	if engine == nil {
		return nil, fmt.Errorf("l2 subscriber: engine is nil")
	}
	if engine.Cfg == nil {
		return nil, fmt.Errorf("l2 subscriber: engine cfg is nil")
	}
	if engine.EthClient == nil {
		return nil, fmt.Errorf("l2 subscriber: engine eth client is nil")
	}
	if strings.TrimSpace(engine.Cfg.OracleContractAddress) == "" {
		return nil, fmt.Errorf("l2 subscriber: ORACLE_CONTRACT_ADDRESS is empty")
	}

	contractAddr := common.HexToAddress(engine.Cfg.OracleContractAddress)
	parsedABI, err := abi.JSON(strings.NewReader(eth.OracleMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("l2 subscriber: parse oracle ABI: %w", err)
	}
	filterer, err := eth.NewOracleFilterer(contractAddr, engine.EthClient)
	if err != nil {
		return nil, fmt.Errorf("l2 subscriber: create oracle filterer: %w", err)
	}

	return &L2ContractEventSubscriber{
		engine:       engine,
		contractAddr: contractAddr,
		contractABI:  parsedABI,
		filterer:     filterer,
		readyCh:      make(chan struct{}),
	}, nil
}

// Engine returns the oracle engine instance owned by this subscriber.
func (l *L2ContractEventSubscriber) Engine() *oracle_runtime.OracleEngine {
	return l.engine
}

// Ready returns a channel that closes once the subscriber has settled on
// either a live subscription or polling.
func (l *L2ContractEventSubscriber) Ready() <-chan struct{} {
	return l.readyCh
}

func (l *L2ContractEventSubscriber) markReady() {
	l.readyOnce.Do(func() {
		close(l.readyCh)
	})
}

// Start begins the background subscription loop. It is safe to call Start
// multiple times; subsequent calls are no-ops.
func (l *L2ContractEventSubscriber) Start(parent context.Context) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.cancel != nil {
		return
	}
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	l.cancel = cancel
	log.Printf("l2 subscriber: start (contract=%s)", l.contractAddr.Hex())
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.run(ctx)
	}()
}

// Stop stops the subscriber and waits for it to exit.
func (l *L2ContractEventSubscriber) Stop() {
	l.mu.Lock()
	cancel := l.cancel
	l.cancel = nil
	l.mu.Unlock()

	if cancel != nil {
		log.Printf("l2 subscriber: stop")
		cancel()
	}
	l.wg.Wait()
}

func (l *L2ContractEventSubscriber) run(ctx context.Context) {
	if err := l.syncState(ctx); err != nil {
		log.Printf("l2 subscriber: initial sync failed: %v", err)
	}
	backoff := time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		logsCh := make(chan types.Log, l2LogsChanBuffer)
		query := ethereum.FilterQuery{Addresses: []common.Address{l.contractAddr}}
		sub, err := l.engine.EthClient.SubscribeFilterLogs(ctx, query, logsCh)
		if err != nil {
			if isNotificationsUnsupported(err) {
				log.Printf("l2 subscriber: subscribe logs unsupported, switching to polling (err=%v)", err)
				l.pollLogs(ctx)
				return
			}
			log.Printf("l2 subscriber: subscribe logs failed: %v (backoff=%s)", err, backoff)
			sleep(ctx, backoff)
			backoff = nextBackoff(backoff)
			continue
		}
		backoff = time.Second

		log.Printf("l2 subscriber: subscribed (contract=%s)", l.contractAddr.Hex())
		l.markReady()

		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return

			case err, ok := <-sub.Err():
				// If the error channel is closed, the subscription is no longer usable.
				if !ok {
					sub.Unsubscribe()
					log.Printf("l2 subscriber: subscription error channel closed (resubscribing)")
					goto RESUBSCRIBE
				}
				if err != nil {
					log.Printf("l2 subscriber: subscription error: %v (resubscribing)", err)
				}
				sub.Unsubscribe()
				goto RESUBSCRIBE

			case evlog, ok := <-logsCh:
				if !ok {
					sub.Unsubscribe()
					log.Printf("l2 subscriber: logs channel closed (resubscribing)")
					goto RESUBSCRIBE
				}
				l.logEvent(evlog)
			}
		}

	RESUBSCRIBE:
		continue
	}
}

func (l *L2ContractEventSubscriber) syncState(ctx context.Context) error {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(50),
		Addresses: []common.Address{l.contractAddr},
	}
	logs, err := l.engine.EthClient.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("filter logs: %w", err)
	}
	for _, evlog := range logs {
		l.logEvent(evlog)
	}
	return nil
}

func (l *L2ContractEventSubscriber) pollLogs(ctx context.Context) {
	log.Printf("l2 subscriber: polling enabled (contract=%s interval=%s)", l.contractAddr.Hex(), l2PollInterval)
	l.markReady()
	ticker := time.NewTicker(l2PollInterval)
	defer ticker.Stop()

	lastBlock, ok := l.fetchLatestBlock(ctx)
	if !ok {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			latest, ok := l.fetchLatestBlock(ctx)
			if !ok {
				continue
			}
			if latest <= lastBlock {
				continue
			}

			query := ethereum.FilterQuery{
				Addresses: []common.Address{l.contractAddr},
				FromBlock: new(big.Int).SetUint64(lastBlock + 1),
				ToBlock:   new(big.Int).SetUint64(latest),
			}
			logs, err := l.engine.EthClient.FilterLogs(ctx, query)
			if err != nil {
				log.Printf("l2 subscriber: poll logs failed: %v", err)
				continue
			}
			for _, evlog := range logs {
				l.logEvent(evlog)
			}
			lastBlock = latest
		}
	}
}

func (l *L2ContractEventSubscriber) fetchLatestBlock(ctx context.Context) (uint64, bool) {
	latest, err := l.engine.EthClient.BlockNumber(ctx)
	if err != nil {
		log.Printf("l2 subscriber: fetch latest block failed: %v", err)
		sleep(ctx, nextBackoff(time.Second))
		return 0, false
	}
	return latest, true
}

func (l *L2ContractEventSubscriber) logEvent(evlog types.Log) {
	name, indexed, nonIndexed, decodeErr := l.decodeEvent(evlog)
	if decodeErr != "" {
		log.Printf(
			"l2 event decode error: name=%s %s topics=%d data=%d indexed=%v args=%v err=%s",
			name,
			eventMeta(evlog),
			len(evlog.Topics),
			len(evlog.Data),
			indexed,
			nonIndexed,
			decodeErr,
		)
	}
	l.handleTypedEvent(name, evlog)
}

func (l *L2ContractEventSubscriber) decodeEvent(evlog types.Log) (string, map[string]any, map[string]any, string) {
	name := "unknown"
	nonIndexed := map[string]any{}
	indexed := map[string]any{}
	decodeErr := ""

	if len(evlog.Topics) == 0 {
		return name, indexed, nonIndexed, decodeErr
	}

	evt, err := l.contractABI.EventByID(evlog.Topics[0])
	if err != nil {
		decodeErr = "event_lookup_err=" + err.Error()
		return name, indexed, nonIndexed, decodeErr
	}
	name = evt.Name

	// Decode indexed fields from topics (excluding the 0th topic, which is the event signature).
	// This allows downstream consumers to work with a complete view of the event.
	if len(evlog.Topics) > 1 {
		idxArgs := abi.Arguments{}
		for _, in := range evt.Inputs {
			if in.Indexed {
				idxArgs = append(idxArgs, in)
			}
		}
		if len(idxArgs) > 0 {
			if err := abi.ParseTopicsIntoMap(indexed, idxArgs, evlog.Topics[1:]); err != nil {
				decodeErr = appendDecodeErr(decodeErr, "parse_topics_err="+err.Error())
			}
		}
	}

	// Decode non-indexed fields from the data payload.
	if len(evlog.Data) > 0 {
		if err := l.contractABI.UnpackIntoMap(nonIndexed, evt.Name, evlog.Data); err != nil {
			decodeErr = appendDecodeErr(decodeErr, "unpack_err="+err.Error())
		}
	}

	return name, indexed, nonIndexed, decodeErr
}

func (l *L2ContractEventSubscriber) handleTypedEvent(name string, evlog types.Log) {
	if l.filterer == nil {
		return
	}

	switch name {
	case "BurnSubmitted":
		l.handleBurnSubmitted(evlog)
	case "ClaimSubmitted":
		l.handleClaimSubmitted(evlog)
	case "Exiting":
		l.handleExiting(evlog)
	case "NewAggregator":
		l.handleNewAggregator(evlog)
	case "Registered":
		l.handleRegistered(evlog)
	case "Replaced":
		l.handleReplaced(evlog)
	case "UserRegistered":
		l.handleUserRegistered(evlog)
	case "ValidatorRegistered":
		l.handleValidatorRegistered(evlog)
	case "WiVoteSubmitted":
		l.handleWiVoteSubmitted(evlog)
	case "Withdrawn":
		l.handleWithdrawn(evlog)
	}
}

func (l *L2ContractEventSubscriber) handleBurnSubmitted(evlog types.Log) {
	evt, err := l.filterer.ParseBurnSubmitted(evlog)
	if err != nil {
		log.Printf("l2 event BurnSubmitted: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event BurnSubmitted: commitmentHash=%s %s",
		common.BytesToHash(evt.CommitmentHash[:]).Hex(),
		eventMeta(evlog),
	)
	if err := l.applyBurnSubmittedEvent(evt); err != nil {
		log.Printf("l2 event BurnSubmitted: apply error: %v", err)
	}
}

func (l *L2ContractEventSubscriber) handleValidatorRegistered(evlog types.Log) {
	evt, err := l.filterer.ParseValidatorRegistered(evlog)
	if err != nil {
		log.Printf("l2 event ValidatorRegistered: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event ValidatorRegistered: addr=%s validatorID=%s index=%s pubkey=(%s,%s) balance=%s %s",
		evt.Addr.Hex(),
		evt.ValidatorID,
		evt.Index,
		evt.Pubkey.X,
		evt.Pubkey.Y,
		evt.Balance,
		eventMeta(evlog),
	)
	if err := l.applyValidatorRegisteredEvent(evt); err != nil {
		log.Printf("l2 event ValidatorRegistered: apply error: %v", err)
	}
}

func (l *L2ContractEventSubscriber) handleClaimSubmitted(evlog types.Log) {
	evt, err := l.filterer.ParseClaimSubmitted(evlog)
	if err != nil {
		log.Printf("l2 event ClaimSubmitted: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event ClaimSubmitted: uniqueID=%s nullifierHash=%s proof=0x%x publicWitness=0x%x proofLen=%d publicWitnessLen=%d %s",
		evt.UniqueID,
		common.BytesToHash(evt.NullifierHash[:]).Hex(),
		evt.Proof,
		evt.PublicWitness,
		len(evt.Proof),
		len(evt.PublicWitness),
		eventMeta(evlog),
	)
	if err := l.applyClaimSubmittedEvent(evt); err != nil {
		log.Printf("l2 event ClaimSubmitted: apply error: %v", err)
	}
}

func (l *L2ContractEventSubscriber) handleExiting(evlog types.Log) {
	evt, err := l.filterer.ParseExiting(evlog)
	if err != nil {
		log.Printf("l2 event Exiting: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event Exiting: sender=%s %s",
		evt.Sender.Hex(),
		eventMeta(evlog),
	)
}

func (l *L2ContractEventSubscriber) handleNewAggregator(evlog types.Log) {
	evt, err := l.filterer.ParseNewAggregator(evlog)
	if err != nil {
		log.Printf("l2 event NewAggregator: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event NewAggregator: validatorID=%s %s",
		evt.ValidatorID,
		eventMeta(evlog),
	)
	if err := l.applyNewAggregatorEvent(evt); err != nil {
		log.Printf("l2 event NewAggregator: apply error: %v", err)
	}
}

func (l *L2ContractEventSubscriber) handleRegistered(evlog types.Log) {
	evt, err := l.filterer.ParseRegistered(evlog)
	if err != nil {
		log.Printf("l2 event Registered: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event Registered: sender=%s index=%s pubkey=(%s,%s) value=%s %s",
		evt.Sender.Hex(),
		evt.Index,
		evt.Pubkey.X,
		evt.Pubkey.Y,
		evt.Value,
		eventMeta(evlog),
	)
	if err := l.applyRegisteredEvent(evt); err != nil {
		log.Printf("l2 event Registered: apply error: %v", err)
	}
}

func (l *L2ContractEventSubscriber) handleReplaced(evlog types.Log) {
	evt, err := l.filterer.ParseReplaced(evlog)
	if err != nil {
		log.Printf("l2 event Replaced: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event Replaced: sender=%s replaced=%s %s",
		evt.Sender.Hex(),
		evt.Replaced.Hex(),
		eventMeta(evlog),
	)
}

func (l *L2ContractEventSubscriber) handleUserRegistered(evlog types.Log) {
	evt, err := l.filterer.ParseUserRegistered(evlog)
	if err != nil {
		log.Printf("l2 event UserRegistered: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event UserRegistered: addr=%s index=%s pubkey=(%s,%s) balance=%s %s",
		evt.Addr.Hex(),
		evt.Index,
		evt.Pubkey.X,
		evt.Pubkey.Y,
		evt.Balance,
		eventMeta(evlog),
	)
}

func (l *L2ContractEventSubscriber) handleWiVoteSubmitted(evlog types.Log) {
	evt, err := l.filterer.ParseWiVoteSubmitted(evlog)
	if err != nil {
		log.Printf("l2 event WiVoteSubmitted: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event WiVoteSubmitted: submitter=%s validators=%s honestBits=%s request=%s majorityVote=%s %s",
		evt.Submitter,
		evt.Validators,
		evt.HonestBits,
		evt.Request,
		evt.MajorityVote,
		eventMeta(evlog),
	)
	if err := l.applyWiVoteSubmittedEvent(evt); err != nil {
		log.Printf("l2 event WiVoteSubmitted: apply error: %v", err)
	}
}

func (l *L2ContractEventSubscriber) applyRegisteredEvent(evt *eth.OracleRegistered) error {
	if l.engine == nil || l.engine.Oracle == nil {
		return fmt.Errorf("oracle engine not initialized")
	}
	return l.engine.Oracle.ApplyRegisteredEvent(evt)
}

func (l *L2ContractEventSubscriber) applyWiVoteSubmittedEvent(evt *eth.OracleWiVoteSubmitted) error {
	if l.engine == nil || l.engine.Oracle == nil {
		return fmt.Errorf("oracle engine not initialized")
	}
	return l.engine.Oracle.ApplyWiVoteSubmittedEvent(evt)
}

func (l *L2ContractEventSubscriber) applyNewAggregatorEvent(evt *eth.OracleNewAggregator) error {
	if l.engine == nil || l.engine.Oracle == nil {
		return fmt.Errorf("oracle engine not initialized")
	}
	newAggID := uint(evt.ValidatorID.Uint64())
	var found bool
	for _, node := range l.engine.Oracle.Nodes {
		if node == nil {
			continue
		}
		if node.ID == newAggID {
			log.Printf("node %d says: I am the new aggregator", node.ID)
			node.IncVotes = make(map[string]map[uint]*oracle.IncVote)
			node.WiVotes = make(map[string]map[uint]*oracle.WiVote)
			l.resetBatchedWiVotes()
			node.Role = oracle.Aggregator
			found = true
		} else {
			node.Role = oracle.Validator
		}
	}
	if found {
		l.engine.Oracle.AggregatorID = newAggID
	}
	return nil
}

func (l *L2ContractEventSubscriber) applyValidatorRegisteredEvent(evt *eth.OracleValidatorRegistered) error {
	if l.engine == nil || l.engine.Oracle == nil {
		return fmt.Errorf("oracle engine not initialized")
	}
	for _, node := range l.engine.Oracle.Nodes {
		if node == nil || node.Account == nil {
			continue
		}
		if evt.ValidatorID.Uint64() != uint64(node.ID) {
			continue
		}
		node.Account.Index = evt.Index
		node.Account.Balance = evt.Balance
		log.Printf("node=%v registered: %v", node.ID, node.Account)
	}
	return nil
}

func (l *L2ContractEventSubscriber) applyBurnSubmittedEvent(evt *eth.OracleBurnSubmitted) error {
	if l.engine == nil || l.engine.Oracle == nil {
		return fmt.Errorf("oracle engine not initialized")
	}
	for _, node := range l.engine.Oracle.Nodes {
		if node == nil || node.IncMerkleTree == nil {
			continue
		}
		log.Printf("node %d received burn event", node.ID)
		eventCommitmentHash := evt.CommitmentHash[:]
		incTreeIndex, _, rootHash, err := node.IncMerkleTree.AddLeafValidator(eventCommitmentHash)
		if err != nil {
			log.Printf("node %d failed to add leaf to merkle tree: %v", node.ID, err)
			continue
		}
		incVote := &oracle.IncVote{
			CommitmentHash: eventCommitmentHash,
			NodeID:         node.ID,
			IncTreeIndex:   incTreeIndex,
			RootHash:       rootHash,
		}
		log.Printf("node %d publishing incVote for burn event (incTreeIndex=%v)...", node.ID, incVote.IncTreeIndex)
		node.Oracle.PublishIncVote(incVote)
	}
	l.engine.Oracle.PublishIncVoteSelectionRes()
	return nil
}

func (l *L2ContractEventSubscriber) applyClaimSubmittedEvent(evt *eth.OracleClaimSubmitted) error {
	if l.engine == nil || l.engine.Oracle == nil {
		return fmt.Errorf("oracle engine not initialized")
	}
	for _, node := range l.engine.Oracle.Nodes {
		if node == nil {
			continue
		}
		log.Printf("node %d received claim event", node.ID)
		wiVote, err := node.VerifyClaim(evt)
		if err != nil {
			log.Fatalf("failed verifying claim: %v", err)
		}
		node.Oracle.PublishWiVote(wiVote)
	}
	l.engine.Oracle.PublishWiVoteSelectionRes()
	return nil
}

func (l *L2ContractEventSubscriber) resetBatchedWiVotes() {
	if l.engine == nil || l.engine.Oracle == nil {
		return
	}
	for _, node := range l.engine.Oracle.Nodes {
		if node == nil || node.BatchedWiVote == nil {
			continue
		}
		node.BatchedWiVote.WithdrawalReqIDs = make([]*big.Int, 0)
		node.BatchedWiVote.Vote = make([]*big.Int, 0)
		node.BatchedWiVote.MajorityOfVotes = make([]*big.Int, 0)
	}
}

func (l *L2ContractEventSubscriber) handleWithdrawn(evlog types.Log) {
	evt, err := l.filterer.ParseWithdrawn(evlog)
	if err != nil {
		log.Printf("l2 event Withdrawn: parse error: %v", err)
		return
	}
	log.Printf(
		"l2 event Withdrawn: sender=%s %s",
		evt.Sender.Hex(),
		eventMeta(evlog),
	)
}

func appendDecodeErr(current, next string) string {
	if current == "" {
		return next
	}
	return current + " " + next
}

func eventMeta(evlog types.Log) string {
	return fmt.Sprintf("block=%d tx=%s idx=%d", evlog.BlockNumber, evlog.TxHash.Hex(), evlog.Index)
}

func formatDecodeErr(err string) string {
	if err == "" {
		return ""
	}
	return " err=" + err
}

func isNotificationsUnsupported(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notifications not supported") ||
		strings.Contains(msg, "notifications are not supported") ||
		strings.Contains(msg, "unsupported subscription")
}

func nextBackoff(cur time.Duration) time.Duration {
	next := cur * 2
	if next > 30*time.Second {
		return 30 * time.Second
	}
	return next
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
