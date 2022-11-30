//go:build !(linux && arm)

package bartend

import "log"

type fakePin int

func (self fakePin) Write(v int) error {
	log.Printf("pin #%d=%d", self, v)
	return nil
}

func GetPin(id int) (Pin, error) {
	return fakePin(id), nil
}
