package gnark

import (
	"context"
	"encoding/hex"
	"fmt"
	"l2alchemy/internal/config"
	bc "l2alchemy/internal/eth"
	"math/big"
	"strings"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type StateSync struct {
	cfg       *config.Config
	wg        *sync.WaitGroup
	index     uint64
	state     *State
	ethClient *ethclient.Client
	contract  *bc.Oracle
}

func NewStateSync(cfg *config.Config, wg *sync.WaitGroup, index uint64, state *State, ethClient *ethclient.Client, contract *bc.Oracle) *StateSync {
	return &StateSync{cfg: cfg, wg: wg, index: index, state: state, ethClient: ethClient, contract: contract}
}

func (s *StateSync) Synchronize() error {

	if err := s.Update(context.Background()); err != nil {
		return fmt.Errorf("update state: %w", err)
	}

	go func() {
		if err := WatchEvent(context.Background(), s.contract.WatchRegistered, s.HandleRegisteredEvent); err != nil {
			fmt.Println(fmt.Errorf("watch registered event: %w", err))
		}
	}()

	go func() {
		if err := WatchEvent(context.Background(), s.contract.WatchWiVoteSubmitted, s.HandleWatchWiVoteSubmittedEvent); err != nil {
			fmt.Println(fmt.Errorf("watch wivote submitted event: %w", err))
		}
	}()

	return nil
}

func (s *StateSync) Update(ctx context.Context) error {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(50),
		Addresses: []common.Address{
			common.HexToAddress(s.cfg.OracleContractAddress),
		},
	}

	logs, err := s.ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("filter logs: %w", err)
	}
	contractABI, err := abi.JSON(strings.NewReader(bc.OracleABI))
	if err != nil {
		return fmt.Errorf("reading contract abi: %w", err)
	}

	for _, log := range logs {
		event, err := contractABI.EventByID(log.Topics[0])
		if err != nil {
			return fmt.Errorf("event by id: %w", err)
		}
		switch event.Name {
		case "Registered":
			e, err := s.contract.ParseRegistered(log)
			if err != nil {
				return fmt.Errorf("parse registered: %w", err)
			}
			err = s.HandleRegisteredEvent(ctx, e)
			if err != nil {
				return fmt.Errorf("handle registered: %w", err)
			}
		case "WiVoteSubmitted":
			e, err := s.contract.ParseWiVoteSubmitted(log)
			if err != nil {
				return fmt.Errorf("parse wivote submitted: %w", err)
			}
			err = s.HandleWatchWiVoteSubmittedEvent(ctx, e)
			if err != nil {
				return fmt.Errorf("handle wivote submitted: %w", err)
			}
		}
	}

	return nil
}

func (s *StateSync) HandleRegisteredEvent(ctx context.Context, event *bc.OracleRegistered) error {
	// fmt.Printf("handle registered event: index=%d, pubKeyX=%s, pubKeyY=%s, balance=%s\n",
	// 	event.Index.Uint64(),
	// 	event.Pubkey.X.String(),
	// 	event.Pubkey.Y.String(),
	// 	event.Value.String())

	x := fr.NewElement(0)
	y := fr.NewElement(0)

	x.SetBigInt(event.Pubkey.X)
	y.SetBigInt(event.Pubkey.Y)

	publicKey := eddsa.PublicKey{A: twistededwards.NewPointAffine(x, y)}

	account := Account{
		Index:         event.Index,
		PublicKey:     &publicKey,
		Balance:       event.Value,
		Reputation:    event.Reputation,
		SeverityCount: event.SeverityCount,
	}

	err := s.state.WriteAccount(account)
	if err != nil {
		return fmt.Errorf("write account: %w", err)
	}

	return nil
}

func (s *StateSync) HandleWatchWiVoteSubmittedEvent(ctx context.Context, event *bc.OracleWiVoteSubmitted) error {
	fmt.Println("handle wivote submitted event: request:%w, submitter%w", event.Request.Uint64(), "submitter", event.Submitter.Uint64())

	aggregatorAccount, err := s.state.ReadAccount(event.Submitter.Uint64())
	if err != nil {
		return fmt.Errorf("read aggregator account: %w", err)
	}
	aggregatorAccount.Reputation.Add(aggregatorAccount.Reputation, big.NewInt(2))
	aggregatorAccount.SeverityCount.Add(aggregatorAccount.SeverityCount, big.NewInt(1))
	aggregatorAccount.Balance.Add(aggregatorAccount.Balance, big.NewInt(RewardAggregator))
	err = s.state.WriteAccount(aggregatorAccount)
	if err != nil {
		return fmt.Errorf("write account: %w", err)
	}

	for i := 0; i < s.cfg.NodeCount; i++ {
		if event.Validators.Bit(i) == 0 {
			continue
		}

		account, err := s.state.ReadAccount(uint64(i))
		if err != nil {
			return fmt.Errorf("read aggregator account: %w", err)
		}
		account.Reputation.Add(account.Reputation, big.NewInt(2))
		account.SeverityCount.Add(account.SeverityCount, big.NewInt(1))
		account.Balance.Add(account.Balance, big.NewInt(RewardValidator))
		err = s.state.WriteAccount(account)
		if err != nil {
			return fmt.Errorf("write account: %w", err)
		}
	}
	root, _ := s.state.Root()
	fmt.Println("root after update:", hex.EncodeToString(root))
	return nil
}
