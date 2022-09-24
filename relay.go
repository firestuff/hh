package main

import "github.com/stianeikeland/go-rpio/v4"

type Relay struct {
	pin rpio.Pin
}

func NewRelay(pin int) *Relay {
	r := &Relay{
		pin: rpio.Pin(pin),
	}

	r.pin.Output()

	return r
}

func (r *Relay) On() {
	r.pin.High()
}

func (r *Relay) Off() {
	r.pin.Low()
}
