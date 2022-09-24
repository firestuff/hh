package hhio

import "context"
import "errors"
import "time"

import "github.com/stianeikeland/go-rpio/v4"

type Ultrasonic struct {
	C chan float64

	trig rpio.Pin
	echo rpio.Pin
}

var timeout = errors.New("timed out waiting for device")

func NewUltrasonic(ctx context.Context, trig, echo int) *Ultrasonic {
	us := &Ultrasonic{
		C:    make(chan float64),
		trig: rpio.Pin(trig),
		echo: rpio.Pin(echo),
	}

	us.trig.Output()
	us.trig.Low()

	us.echo.Input()

	go us.loop(ctx)

	return us
}

func (us *Ultrasonic) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(us.C)
			return

		default:
			us.measure()
		}
	}
}

func (us *Ultrasonic) measure() {
	deadline := time.NewTimer(200 * time.Millisecond)

	us.trigger()

	start := us.wait(rpio.High, deadline)
	end := us.wait(rpio.Low, deadline)

	if start.IsZero() || end.IsZero() {
		return
	}

	duration := end.Sub(start)
	cm := duration.Seconds() * 17150

	if cm > 400 {
		return
	}

	us.C <- cm
}

func (us *Ultrasonic) trigger() {
	us.trig.High()
	time.Sleep(10 * time.Microsecond)
	us.trig.Low()
}

func (us *Ultrasonic) wait(goal rpio.State, deadline *time.Timer) time.Time {
	for {
		select {
		case <-deadline.C:
			return time.Time{}

		default:
			if us.echo.Read() == goal {
				return time.Now()
			}
		}
	}
}
