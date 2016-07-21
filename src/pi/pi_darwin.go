package pi

import (
	"time"
	"log"
	"github.com/stianeikeland/go-rpio"
)

func Open() {
}

func Pin(id int) rpio.Pin {
	pin := rpio.Pin(id)
	return pin
}

func TurnOnFor(pin rpio.Pin, howLong time.Duration) {
	log.Printf("pin %d on for %dms", pin, howLong / time.Millisecond)
	time.Sleep(howLong)
	log.Printf("pin %d off", pin)
}
