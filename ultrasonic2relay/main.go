package main

import "context"
import "fmt"

import "github.com/firestuff/hh/hhio"

func main() {
	err := hhio.Open()
	if err != nil {
		panic(err)
	}
	defer hhio.Close()

	us := hhio.NewUltrasonic(context.Background(), 6, 5)
	mf := hhio.NewMedianFilter(us.C, 47)

  r := hhio.NewRelay(21)

	for dist := range mf {
		fmt.Printf("%f\n", dist)

    if dist < 50 {
      r.On()
    } else if dist > 70 {
      r.Off()
    }
	}
}
