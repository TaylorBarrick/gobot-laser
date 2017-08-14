package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"

	"flag"

	serial "go.bug.st/serial.v1"
)

var lightLevel int
var lightSensorPin string
var laserPin string
var device string
var clockSpeed int64

func init() {
	flag.IntVar(&lightLevel, "LightLevel", 100, "--LightLevel=100 //Represents the amount of resistance measured when the laser is detected.  Higher values represent not detecting the laser, while lower values detect the laser as on.")
	flag.StringVar(&lightSensorPin, "LightPin", "5", "--Light=5 //The pin in which the light sensor is connected")
	flag.StringVar(&laserPin, "LaserPin", "10", "--Laser=10 //The pin in which the laser is connected")
	flag.Int64Var(&clockSpeed, "Clock", 50, "--Clock=500 //Interval for sending and receiving data in Milliseconds")
	flag.StringVar(&device, "Device", "COM4", "--Device=COM4 //The serial device used for Firmata communication")
}

func main() {
	flag.Parse()

	bitsToTransmit := make(chan bool, 10000)
	bytesReceived := make(chan byte, 100)

	firmataAdaptor := firmata.NewAdaptor(device)
	laserDriver := gpio.NewLedDriver(firmataAdaptor, laserPin)

	mode := &serial.Mode{DataBits: 8, Parity: serial.OddParity, StopBits: serial.TwoStopBits}

	decoder := NewDecoder(bytesReceived, mode)
	encoder := NewEncoder(bitsToTransmit, mode)

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{},
		func() {
			go detectAndSendInput(os.Stdin, encoder)
			go func() {
				for {
					select {
					case b := <-bytesReceived:
						fmt.Printf("%c", b)
					}
				}
			}()
			tick := time.Tick(time.Duration(clockSpeed) * time.Millisecond)
			for {
				select {
				case <-tick:
					checkLaserChannel(bitsToTransmit, Laser{laserDriver})
					checkLightSensor(firmataAdaptor, lightSensorPin, lightLevel, decoder)
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

type Laser struct {
	Emitter
}

func (l Laser) Send(b bool) error {
	if b {
		return l.On()
	}
	return l.Off()
}
func checkLaserChannel(c chan bool, l Laser) error {
	select {
	case b := <-c:
		return l.Send(b)
	}
}

type BitReader interface {
	Read(bit bool) error
}

func checkLightSensor(r aio.AnalogReader, pin string, threshold int, br BitReader) (err error) {
	i := 0
	if i, err = r.AnalogRead(pin); err != nil {
		return err
	}
	br.Read(i < threshold)
	return
}

func detectAndSendInput(r io.Reader, w io.Writer) {
	in := ""
	for {
		fmt.Fscanln(r, &in)
		if _, err := w.Write([]byte(in)); err != nil {
			panic(err)
		}
	}
}
