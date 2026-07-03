package events

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

// F-21 acceptance: a shuffled burn-event delivery must produce the same
// insertion order on every node. The pipeline orders finalized burns by
// (source rollup id, block height, log index, commitment hash) per paper
// SIV-E, so sorting any permutation must yield one canonical sequence.
func TestPendingBurnDeterministicOrdering(t *testing.T) {
	burns := []PendingBurn{
		{SourceID: 1, BlockNumber: 5, LogIndex: 0, CommitmentHash: [32]byte{0x01}},
		{SourceID: 1, BlockNumber: 5, LogIndex: 1, CommitmentHash: [32]byte{0x02}},
		{SourceID: 1, BlockNumber: 6, LogIndex: 0, CommitmentHash: [32]byte{0x03}},
		{SourceID: 2, BlockNumber: 1, LogIndex: 0, CommitmentHash: [32]byte{0x04}},
		{SourceID: 2, BlockNumber: 1, LogIndex: 0, CommitmentHash: [32]byte{0x05}},
	}

	canonical := append([]PendingBurn(nil), burns...)
	sort.Slice(canonical, func(i, j int) bool { return canonical[i].less(canonical[j]) })

	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 50; trial++ {
		shuffled := append([]PendingBurn(nil), burns...)
		rng.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].less(shuffled[j]) })
		if !reflect.DeepEqual(shuffled, canonical) {
			t.Fatalf("trial %d: order diverged:\n got %v\nwant %v", trial, shuffled, canonical)
		}
	}

	// The canonical order itself follows the paper's sort keys.
	for i := 1; i < len(canonical); i++ {
		if canonical[i].less(canonical[i-1]) {
			t.Fatalf("canonical order violated at %d: %v before %v", i, canonical[i-1], canonical[i])
		}
	}
}
