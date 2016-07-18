package pi

import (
	"time"
	"github.com/c2g/c2"
)

func LightLed(pin int, howLong time.Duration) {
	c2.Debug.Printf("Turn on gpio %d for %dms", pin, howLong)
	time.Sleep(howLong)
}
