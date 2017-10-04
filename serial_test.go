package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parity(t *testing.T) {
	bOne := true
	bZero := false
	tests := []struct {
		name string
		u    uint
		p    Parity
		want *bool
	}{
		{"even_match", 0xFF, EvenParity, &bOne},
		{"odd_match", 0xFE, OddParity, &bOne},
		{"even_bad", 0xFE, EvenParity, &bZero},
		{"odd_bad", 0xFF, OddParity, &bZero},
		{"mark", 0x00, MarkParity, &bOne},
		{"spaces", 0x00, SpaceParity, &bZero},
		{"none", 0x00, NoParity, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parity(tt.u, tt.p); !cmp.Equal(got, tt.want) {
				t.Errorf("parity() got and want mismatch\n%v", cmp.Diff(got, tt.want))
			}
		})
	}
}
