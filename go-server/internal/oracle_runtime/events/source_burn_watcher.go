package events

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"l2alchemy/internal/config"
	"l2alchemy/internal/eth"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SourceBurnWatcher watches BurnSubmitted events on the source rollup's
// Gateway (paper SIV-E / SV: burns execute on source rollup L2A while
// claims execute on destination rollup L2B with separate states; F-20) and
// feeds them into the destination subscriber's burn-finality pipeline
// (F-21). One watcher corresponds to one source Gateway instance and hence
// one commitment tree, per the paper's per-instance trees.
type SourceBurnWatcher struct {
	client       *ethclient.Client
	filterer     *eth.OracleFilterer
	contractAddr common.Address
	sourceID     uint64
	dest         *L2ContractEventSubscriber

	mu     sync.Mutex
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSourceBurnWatcher builds a watcher from configuration; it returns
// (nil, nil) when no source chain is configured, in which case burns are
// taken from the destination Gateway (single-rollup simulation).
func NewSourceBurnWatcher(ctx context.Context, cfg *config.Config, dest *L2ContractEventSubscriber) (*SourceBurnWatcher, error) {
	if cfg == nil || cfg.SourceRPCURL == "" || cfg.SourceOracleContractAddress == "" {
		return nil, nil
	}
	if dest == nil {
		return nil, fmt.Errorf("source watcher: destination subscriber is nil")
	}

	client, err := ethclient.DialContext(ctx, cfg.SourceRPCURL)
	if err != nil {
		return nil, fmt.Errorf("source watcher: dial source rpc: %w", err)
	}

	sourceID := uint64(cfg.SourceChainID)
	if sourceID == 0 {
		chainID, err := client.ChainID(ctx)
		if err != nil {
			return nil, fmt.Errorf("source watcher: resolve source chain id: %w", err)
		}
		sourceID = chainID.Uint64()
	}

	contractAddr := common.HexToAddress(cfg.SourceOracleContractAddress)
	filterer, err := eth.NewOracleFilterer(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("source watcher: create oracle filterer: %w", err)
	}

	// Route burns through the finality pipeline with this source's own
	// head-height provider, and stop treating destination-side burns as
	// commitment-tree inputs.
	dest.SetSourceHeadFn(sourceID, client.BlockNumber)
	dest.MarkSourceConfigured()

	return &SourceBurnWatcher{
		client:       client,
		filterer:     filterer,
		contractAddr: contractAddr,
		sourceID:     sourceID,
		dest:         dest,
	}, nil
}

// Start begins polling the source chain for BurnSubmitted events. Safe to
// call once; subsequent calls are no-ops.
func (w *SourceBurnWatcher) Start(parent context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return
	}
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	w.cancel = cancel
	log.Printf("source watcher: start (source=%d gateway=%s)", w.sourceID, w.contractAddr.Hex())
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.run(ctx)
	}()
}

// Stop stops the watcher and waits for it to exit.
func (w *SourceBurnWatcher) Stop() {
	w.mu.Lock()
	cancel := w.cancel
	w.cancel = nil
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	w.wg.Wait()
}

func (w *SourceBurnWatcher) run(ctx context.Context) {
	var lastBlock uint64
	if head, err := w.client.BlockNumber(ctx); err == nil {
		lastBlock = head
	}

	ticker := time.NewTicker(l2PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		head, err := w.client.BlockNumber(ctx)
		if err != nil {
			log.Printf("source watcher: head lookup failed: %v", err)
			continue
		}
		if head <= lastBlock {
			continue
		}

		query := ethereum.FilterQuery{
			Addresses: []common.Address{w.contractAddr},
			FromBlock: new(big.Int).SetUint64(lastBlock + 1),
			ToBlock:   new(big.Int).SetUint64(head),
		}
		logs, err := w.client.FilterLogs(ctx, query)
		if err != nil {
			log.Printf("source watcher: filter logs failed: %v", err)
			continue
		}
		lastBlock = head

		for _, evlog := range logs {
			evt, err := w.filterer.ParseBurnSubmitted(evlog)
			if err != nil {
				// Not a BurnSubmitted event; the source Gateway emits
				// other events we do not consume here.
				continue
			}
			log.Printf(
				"source burn observed: source=%d block=%d logIndex=%d commitment=%s",
				w.sourceID, evlog.BlockNumber, evlog.Index,
				common.BytesToHash(evt.CommitmentHash[:]).Hex(),
			)
			w.dest.EnqueueBurn(PendingBurn{
				SourceID:       w.sourceID,
				BlockNumber:    evlog.BlockNumber,
				LogIndex:       uint(evlog.Index),
				CommitmentHash: evt.CommitmentHash,
			})
		}
	}
}
