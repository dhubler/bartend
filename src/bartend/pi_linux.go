package bartend

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

// SainSmart relay flips off on pin.High and on on pin.Low so things will be
// reversed.
var openedPins map[int]embd.DigitalPin = make(map[int]embd.DigitalPin)

func init() {
	// On exit, turn off all pumps
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\nClosing gpio")
		for _, pin := range openedPins {
			pin.Write(embd.High)
		}
		embd.CloseGPIO()
		os.Exit(1)
	}()
}

func GetPin(id int) (Pin, error) {
	if p, opened := openedPins[id]; opened {
		return p, nil
	}
	dp, err := embd.NewDigitalPin(id)
	if err != nil {
		return nil, err
	}
	if err = dp.SetDirection(embd.Out); err != nil {
		return nil, err
	}
	if err = dp.Write(embd.High); err != nil {
		return nil, err
	}
	openedPins[id] = dp
	return dp, nil
}
