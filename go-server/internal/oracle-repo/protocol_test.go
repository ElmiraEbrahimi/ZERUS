package oracle

import (
	"math/big"
	"testing"
)

// F-09: the finalization threshold is the paper's f+1 with
// f = floor((n-1)/3) (SIV-A/SIV-D, Alg. 2 line 19). The table pins the
// quorum boundary for the committee sizes evaluated in SVI.
func TestBFTThreshold(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{0, 0},
		{-1, 0},
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 2},  // f=1
		{7, 3},  // f=2
		{8, 3},  // f=2
		{16, 6}, // f=5
		{32, 11},
		{64, 22},
		{128, 43},
	}
	for _, c := range cases {
		if got := bftThreshold(c.n); got != c.want {
			t.Errorf("bftThreshold(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}

// F-07: round r contains exactly the claim identifiers [r*b, (r+1)*b - 1]
// (paper SIV-D); anything else must be rejected.
func TestValidateRoundWindow(t *testing.T) {
	ids := func(vals ...int64) []*big.Int {
		out := make([]*big.Int, len(vals))
		for i, v := range vals {
			out[i] = big.NewInt(v)
		}
		return out
	}

	t.Run("valid windows", func(t *testing.T) {
		for _, c := range []struct {
			ids   []*big.Int
			b     int
			round int64
		}{
			{ids(0), 1, 0},
			{ids(7), 1, 7},
			{ids(0, 1, 2, 3), 4, 0},
			{ids(20, 21, 22, 23, 24), 5, 4},
		} {
			round, err := validateRoundWindow(c.ids, c.b)
			if err != nil {
				t.Fatalf("validateRoundWindow(%v, %d): %v", c.ids, c.b, err)
			}
			if round.Int64() != c.round {
				t.Errorf("round = %s, want %d", round, c.round)
			}
		}
	})

	t.Run("rejects non-window batches", func(t *testing.T) {
		for _, c := range []struct {
			ids []*big.Int
			b   int
		}{
			{ids(1, 2, 3, 4), 4}, // straddles two windows
			{ids(0, 1, 3, 4), 4}, // gap
			{ids(0, 1, 2), 4},    // short batch
			{ids(0, 0, 1, 2), 4}, // duplicate
		} {
			if _, err := validateRoundWindow(c.ids, c.b); err == nil {
				t.Errorf("validateRoundWindow(%v, %d): expected error", c.ids, c.b)
			}
		}
	})
}
