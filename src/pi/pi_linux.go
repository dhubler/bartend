package pi

import (
	"github.com/davecheney/gpio"
	"time"
)

func LightLed() {
	pin, err := gpio.OpenPin(7, gpio.ModeOutput)
	if err != nil {
		panic(err)
	}
	for {
		pin.Set()
		time.Sleep(1 * time.Second)
		pin.Clear()
		time.Sleep(1 * time.Second)
	}
}