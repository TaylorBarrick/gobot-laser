// +build go1.9

package main

import (
	"math/bits"

	"github.com/pkg/errors"
)

type Encoder struct {
	c    chan bool
	mask uint16
	mode *Mode
}

func NewEncoder(out chan bool, mode *Mode) *Encoder {
	var e Encoder
	e.c = out
	e.mode = mode
	trailingBits := 0
	switch mode.StopBits {
	case OneStopBit:
		trailingBits = 1
	case TwoStopBits:
		trailingBits = 2
	}
	switch mode.Parity {
	case NoParity:
	default:
		trailingBits++
	}
	e.mask = 1 << uint16(mode.DataBits+trailingBits)
	return &e
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	for i, b := range p {
		u16 := uint16(b)
		if pb := parity(uint(b), e.mode.Parity); *pb {
			u16 = u16<<1 | 1
		}
		u16 = e.mask | u16<<uint(e.mode.StopBits)
		feedChannel(u16, 4+e.mode.DataBits, e.c)
		n = i
	}
	return
}

func parity(u uint, p Parity) *bool {
	var parity bool
	switch p {
	case NoParity:
		return nil
	case OddParity:
		parity = bits.OnesCount(u)%2 != 0
	case EvenParity:
		parity = bits.OnesCount(u)%2 == 0
	case MarkParity:
		parity = true
	case SpaceParity:
		parity = false
	}
	return &parity
}

func feedChannel(ub uint16, packetSize int, c chan bool) {
	for i := packetSize; i >= 0; i-- { // Start at fire bit
		mask := uint16(1 << byte(i))
		on := ub&mask == mask
		c <- on
	}
}

type Decoder struct {
	c       chan byte
	buffer  byte
	reading bool
	pos     int
	mode    *Mode
}

func NewDecoder(out chan byte, mode *Mode) *Decoder {
	return &Decoder{c: out, mode: mode}
}

func (d *Decoder) Read(bit bool) (err error) {
	if !d.reading {
		d.reading = bit
		return
	}

	if d.pos < d.mode.DataBits {
		var b byte
		if bit {
			b = 1
		}
		d.buffer = d.buffer<<1 | b
	}

	if d.pos == d.mode.DataBits {
		if pb := parity(uint(d.buffer), d.mode.Parity); pb != nil && *pb != bit {
			return errors.New("Parity Error")
		}
		d.c <- d.buffer
	}

	d.pos++

	if d.pos == d.mode.DataBits+3 {
		d.reading = false
		d.pos = 0
		d.buffer = 0x00
		return nil
	}

	return nil
}
