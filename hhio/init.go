package hhio

import "github.com/stianeikeland/go-rpio/v4"

func Open() error {
	return rpio.Open()
}

func Close() {
	rpio.Close()
}
