package bartend

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func Open() {
}

func Pin(id int) rpio.Pin {
	pin := rpio.Pin(id)
	return pin
}

func PinOn(pin rpio.Pin, on bool) {
	log.Printf("pin %d %v", pin, on)
}

func TurnOnFor(pin rpio.Pin, howLong time.Duration) {
	log.Printf("pin %d on for %dms", pin, howLong/time.Millisecond)
	time.Sleep(howLong)
	log.Printf("pin %d off", pin)
}
