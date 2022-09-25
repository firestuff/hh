package hhio

import "github.com/stianeikeland/go-rpio/v4"

type Relay struct {
	pin rpio.Pin
	on  bool
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
	r.on = true
}

func (r *Relay) Off() {
	r.pin.Low()
	r.on = false
}

func (r *Relay) IsOn() bool {
	return r.on
}
