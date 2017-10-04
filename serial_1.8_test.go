// +build !go1.9

package main

import (
	"testing"
)

func BenchmarkParity_Even_18(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parity(0xFF, OddParity)
	}
}

func BenchmarkParity_Odd_18(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parity(0xFF, EvenParity)
	}
}
