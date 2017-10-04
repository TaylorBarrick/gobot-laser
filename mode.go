package main

type Mode struct {
	DataBits int
	Parity   Parity
	StopBits StopBits
}

type Parity int

const (
	NoParity Parity = iota
	MarkParity
	SpaceParity
	EvenParity
	OddParity
)

type StopBits int

const (
	OneStopBit StopBits = iota
	TwoStopBits
)
