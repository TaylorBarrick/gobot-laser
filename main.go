package main

import (
	"fmt"
	"io"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"

	serial "go.bug.st/serial.v1"
)

const LASERTHRESHOLD = 100
const CLOCK = 500 * time.Millisecond
const LIGHTSENSOR = "5"
const LASERPIN = "10"
const COMPORT = "COM4"

var bitsToTransmit chan bool
var bytesReceived chan byte
var firmataAdaptor *firmata.Adaptor
var laserDriver Emitter
var rdr *Reader
var snd *Sender

func main() {
	bitsToTransmit = make(chan bool, 10000)
	bytesReceived = make(chan byte, 100)

	firmataAdaptor = firmata.NewAdaptor(COMPORT)
	laserDriver = gpio.NewLedDriver(firmataAdaptor, LASERPIN)

	mode := &serial.Mode{DataBits: 8, Parity: serial.OddParity, StopBits: serial.TwoStopBits}

	rdr = NewReader(bytesReceived, mode)
	snd = NewSender(bitsToTransmit, mode)

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{},
		func() {
			go detectAndSendInput(snd)
			go func() {
				for {
					select {
					case b := <-bytesReceived:
						fmt.Printf("%c", b)
					}
				}
			}()
			tick := time.Tick(CLOCK)
			for {
				select {
				case <-tick:
					checkLightSensor(firmataAdaptor, LIGHTSENSOR, LASERTHRESHOLD, rdr)
					checkLaserChannel(laserDriver)
				}
			}
		},
	)

	robot.Start()
}

type Emitter interface {
	On() error
	Off() error
}

type AnalogSensor interface {
	AnalogRead(pin string) (val int, err error)
}

type BitReader interface {
	Read(bit bool) error
}

func checkLaserChannel(e Emitter) (done bool) {
	select {
	case b := <-bitsToTransmit:
		if b {
			e.On()
		} else {
			e.Off()
		}
	default:
		return true
	}
	return
}

func checkLightSensor(s AnalogSensor, pin string, threshold int, r BitReader) {
	i, _ := s.AnalogRead(pin)
	r.Read(i < threshold)
}

func detectAndSendInput(w io.Writer) {
	in := ""
	for {
		fmt.Scanln(&in)
		if _, err := w.Write([]byte(in)); err != nil {
			panic(err)
		}
	}
}
