package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
	"github.com/kidoman/embed"
	_ "github.com/kidoman/embd/host/rpi" // this load pi driver
)

func main() {
	// set GPIO25 to output mode
	if embd.InitGPIO(); err != nil {
	   panic(err)
	}
	embd.NewPWMPin("
	drv
	pin, err := gpio.OpenPin(rpi.GPIO25, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}

	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nClearing and unexporting the pin.\n")
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	for {
		pin.Set()
		time.Sleep(100 * time.Millisecond)
		pin.Clear()
		time.Sleep(100 * time.Millisecond)
	}
}