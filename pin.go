package bartend

import (
	"log"
	"time"

	"github.com/kidoman/embd"
)

type Pin interface {
	Write(v int) error
}

func TurnOnFor(pin Pin, howLong time.Duration) {
	log.Printf("pin %d on for %dms", pin, howLong/time.Millisecond)
	pin.Write(embd.Low)
	defer pin.Write(embd.High)
	time.Sleep(howLong)
	log.Printf("pin %d off", pin)
}
