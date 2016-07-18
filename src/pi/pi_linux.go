package pi

import (
	"github.com/davecheney/gpio"
	"time"
)

func LightLed(pinId int, howLong time.Duration) {
	pin, err := gpio.OpenPin(pinId, gpio.ModeOutput)
	if err != nil {
		panic(err)
	}
	pin.Set()
	defer pin.Clear()
	time.Sleep(howLong)
}