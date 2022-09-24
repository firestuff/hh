package main

import "context"
import "fmt"

import "github.com/stianeikeland/go-rpio/v4"

func main() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()

	us := NewUltrasonic(context.TODO(), 6, 5)
	mf := NewMedianFilter(us.C, 9)

	for dist := range mf {
		fmt.Printf("%f\n", dist)
	}
}
