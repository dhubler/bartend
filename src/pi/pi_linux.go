package pi

import (
	"time"
	"log"
	"github.com/stianeikeland/go-rpio"
	"os"
	"fmt"
	"os/signal"
)

// SainSmart relay flips off on pin.High and on on pin.Low so things will be
// reversed.
var openedPins map[int]rpio.Pin = make(map[int]rpio.Pin)

func init() {
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nClosing gpio")
			for _, pin := range openedPins {
				pin.High()
			}
			rpio.Close()
			os.Exit(0)
		}
	}()
}

func Pin(id int) rpio.Pin {
	if p, opened := openedPins[id]; opened {
		return p
	}
	p := rpio.Pin(id)
	rpio.PinMode(p, rpio.Output)
	openedPins[id] = p
	return p
}

func TurnOnFor(pin rpio.Pin, howLong time.Duration) {
	log.Printf("pin %d on for %dms", pin, howLong / time.Millisecond)
	pin.Low()
	defer pin.High()
	time.Sleep(howLong)
	log.Printf("pin %d off", pin)
}