package main

import (
	"errors"

	serial "go.bug.st/serial.v1"
)

type Encoder struct {
	c    chan bool
	mask uint16
	mode *serial.Mode
}

func NewEncoder(out chan bool, mode *serial.Mode) *Encoder {
	var e Encoder
	e.c = out
	e.mode = mode
	trailingBits := 0
	switch mode.StopBits {
	case serial.OneStopBit:
		trailingBits = 1
	case serial.OnePointFiveStopBits:
	//???
	case serial.TwoStopBits:
		trailingBits = 2
	}
	switch mode.Parity {
	case serial.NoParity:
	default:
		trailingBits++
	}
	e.mask = 1 << uint16(mode.DataBits+trailingBits)
	return &e
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	for i, b := range p {
		ub := uint16(b)
		if pb := parity(b, e.mode.Parity); *pb {
			ub = ub<<1 | 1
		}
		ub = e.mask | ub<<uint(e.mode.StopBits)
		feedChannel(ub, 4+e.mode.DataBits, e.c)
		n = i
	}
	return
}

func parity(b byte, p serial.Parity) *bool {
	var parity bool
	switch p {
	case serial.NoParity:
		return nil
	case serial.OddParity:
		parity = false
	case serial.EvenParity:
		parity = true
	case serial.MarkParity:
		parity = true
		return &parity
	case serial.SpaceParity:
		parity = false
		return &parity
	}
	for x := 0; x < 8; x++ {
		mask := byte(1 << byte(7-x))
		on := b&mask == mask
		if on {
			parity = !parity
		}
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
	mode    *serial.Mode
}

func NewDecoder(out chan byte, mode *serial.Mode) *Decoder {
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
		if pb := parity(d.buffer, d.mode.Parity); pb != nil && *pb != bit {
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
