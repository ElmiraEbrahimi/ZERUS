package oracle

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	votingbatch "l2alchemy/circuits/voting_batch"
	"l2alchemy/internal/eth"
	"l2alchemy/internal/oracle-repo/gnark"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func (o *Oracle) ApplyRegisteredEvent(event *eth.OracleRegistered) error {
	if o == nil {
		return fmt.Errorf("oracle: nil")
	}
	if event == nil {
		return fmt.Errorf("oracle: registered event is nil")
	}

	x := fr.NewElement(0)
	y := fr.NewElement(0)
	x.SetBigInt(event.Pubkey.X)
	y.SetBigInt(event.Pubkey.Y)
	publicKey := eddsa.PublicKey{A: twistededwards.NewPointAffine(x, y)}

	account := gnark.Account{
		Index:         event.Index,
		PublicKey:     &publicKey,
		Balance:       event.Value,
		Reputation:    big.NewInt(50),
		SeverityCount: big.NewInt(0),
	}

	return o.applyToStates(func(state *gnark.State) error {
		if err := state.WriteAccount(account); err != nil {
			return fmt.Errorf("write account: %w", err)
		}
		return nil
	})
}

func (o *Oracle) ApplyWiVoteSubmittedEvent(event *eth.OracleWiVoteSubmitted) error {
	if o == nil {
		return fmt.Errorf("oracle: nil")
	}
	if event == nil {
		return fmt.Errorf("oracle: wivote event is nil")
	}
	if o.cfg == nil {
		return fmt.Errorf("oracle: cfg is nil")
	}

	nodeCount := o.cfg.NodeCount
	err := o.applyToStates(func(state *gnark.State) error {
		aggregatorAccount, err := state.ReadAccount(event.Submitter.Uint64())
		if err != nil {
			return fmt.Errorf("read aggregator account: %w", err)
		}
		aggregatorAccount.Balance.Add(aggregatorAccount.Balance, big.NewInt(votingbatch.RewardAggregator))
		if err := state.WriteAccount(aggregatorAccount); err != nil {
			return fmt.Errorf("write aggregator account: %w", err)
		}

		for i := 0; i < nodeCount; i++ {
			if event.Validators.Bit(i) == 0 {
				continue
			}

			account, err := state.ReadAccount(uint64(i))
			if err != nil {
				return fmt.Errorf("read validator account (i=%d): %w", i, err)
			}

			if event.HonestBits.Bit(i) == 1 {
				account.Balance.Add(account.Balance, big.NewInt(votingbatch.RewardValidator))
			} else {
				// dishonest validator → subtract penalty, clamp at 0
				penalty := big.NewInt(votingbatch.PenaltyValidator)
				if account.Balance.Cmp(penalty) <= 0 {
					account.Balance.SetInt64(0)
				} else {
					account.Balance.Sub(account.Balance, penalty)
				}
			}

			if err := state.WriteAccount(account); err != nil {
				return fmt.Errorf("write validator account (i=%d): %w", i, err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	o.logStateRoot("wivote submitted")
	return nil
}

func (o *Oracle) applyToStates(apply func(state *gnark.State) error) error {
	if apply == nil {
		return nil
	}
	for id, node := range o.Nodes {
		if node == nil || node.state == nil {
			continue
		}
		if err := apply(node.state); err != nil {
			return fmt.Errorf("node %d: %w", id, err)
		}
	}
	return nil
}

func (o *Oracle) logStateRoot(reason string) {
	var state *gnark.State
	for _, node := range o.Nodes {
		if node != nil && node.state != nil {
			state = node.state
			break
		}
	}
	if state == nil {
		return
	}
	root, err := state.Root()
	if err != nil {
		log.Printf("oracle state root (%s): %v", reason, err)
		return
	}
	log.Printf("oracle state root after %s: %s", reason, hex.EncodeToString(root))
}
