package main

import (
	"errors"

	serial "go.bug.st/serial.v1"
)

type Sender struct {
	c    chan bool
	mask uint16
	mode *serial.Mode
}

func NewSender(c chan bool, mode *serial.Mode) *Sender {
	var s Sender
	s.c = c
	s.mode = mode
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
	s.mask = 1 << uint16(mode.DataBits+trailingBits)
	return &s
}

func (s *Sender) Write(p []byte) (n int, err error) {
	for i, b := range p {
		ub := uint16(b)
		if pb := parity(b, s.mode.Parity); *pb {
			ub = ub<<1 | 1
		}
		ub = s.mask | ub<<uint(s.mode.StopBits)
		feedChannel(ub, 4+s.mode.DataBits, s.c)
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

type Reader struct {
	c       chan byte
	buffer  byte
	reading bool
	pos     int
	mode    *serial.Mode
}

func NewReader(c chan byte, mode *serial.Mode) *Reader {
	return &Reader{c: c, mode: mode}
}

func (r *Reader) Read(bit bool) (err error) {
	if !r.reading {
		r.reading = bit
		return
	}

	if r.pos < r.mode.DataBits {
		var b byte
		if bit {
			b = 1
		}
		r.buffer = r.buffer<<1 | b
	}

	if r.pos == r.mode.DataBits {
		if pb := parity(r.buffer, r.mode.Parity); pb != nil && *pb != bit {
			return errors.New("Parity Error")
		}
		r.c <- r.buffer
	}

	r.pos++

	if r.pos == r.mode.DataBits+3 {
		r.reading = false
		r.pos = 0
		r.buffer = 0x00
		return nil
	}

	return nil
}
