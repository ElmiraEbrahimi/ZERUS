package handlers

import "testing"

func TestIsReadyBatchStatus(t *testing.T) {
	cases := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "empty", status: "", want: false},
		{name: "committed", status: "committed", want: false},
		{name: "fast finalized camel", status: "fastFinalized", want: true},
		{name: "fast finalized snake", status: "fast_finalized", want: true},
		{name: "verified", status: "verified", want: true},
		{name: "executed", status: "executed", want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isReadyBatchStatus(tc.status); got != tc.want {
				t.Fatalf("isReadyBatchStatus(%q) = %t, want %t", tc.status, got, tc.want)
			}
		})
	}
}

func TestHasNonZeroHash(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "empty", value: "", want: false},
		{name: "zero hex prefix", value: "0x", want: false},
		{name: "zero single digit", value: "0x0", want: false},
		{name: "zero padded", value: "0x000000", want: false},
		{name: "leading zeros", value: "0x000001", want: true},
		{name: "non hex", value: "abcd", want: true},
		{name: "non zero zero prefix", value: "0xabcdef", want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasNonZeroHash(tc.value); got != tc.want {
				t.Fatalf("hasNonZeroHash(%q) = %t, want %t", tc.value, got, tc.want)
			}
		})
	}
}
