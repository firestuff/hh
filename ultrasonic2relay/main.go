package main

import "context"
import "encoding/json"
import "flag"
import "log"
import "math"
import "os"
import "time"

import "github.com/firestuff/hh/hhio"
import "gopkg.in/yaml.v3"

var conf = flag.String("conf", "", "path to config file")

type Config struct {
	Relay
	Ultrasonics []Ultrasonic

	MedianBuffer int

	OnCM  float64
	OffCM float64

	OnSeconds  float64
	OffSeconds float64
}

type Relay struct {
	Control int
}

type Ultrasonic struct {
	Trigger int
	Echo    int
}

type state int

const (
	onTimer state = iota
	offTimer
	watching
)

func main() {
	flag.Parse()

	cf, err := readConf()
	if err != nil {
		panic(err)
	}

	err = hhio.Open()
	if err != nil {
		panic(err)
	}
	defer hhio.Close()

	uss := []chan float64{}

	for _, uscf := range cf.Ultrasonics {
		us := hhio.NewUltrasonic(context.Background(), uscf.Trigger, uscf.Echo)
		mf := hhio.NewMedianFilter(us.C, cf.MedianBuffer)
		uss = append(uss, mf)
	}

	r := hhio.NewRelay(cf.Relay.Control)

	last := make([]float64, len(uss))
	for i := range uss {
		last[i] = math.MaxFloat64
	}

	var st state = watching
	var onUntil, offUntil time.Time

	for {
		// Fetch new values
		for i, us := range uss {
			select {
			case dist := <-us:
				last[i] = dist

			default:
			}
		}

		// Count votes
		var on, off int
		for _, v := range last {
			if v < cf.OnCM {
				on++
			} else if v > cf.OffCM {
				off++
			}
		}

		switch st {
		case watching:
			if on > 0 {
				log.Printf("on     %s", fmtDists(last))
				r.On()
				onUntil = time.Now().Add(time.Duration(cf.OnSeconds * float64(time.Second)))
				st = onTimer
			}

		case onTimer:
			if time.Now().After(onUntil) {
				log.Printf("off    %s", fmtDists(last))
				r.Off()
				offUntil = time.Now().Add(time.Duration(cf.OffSeconds * float64(time.Second)))
				st = offTimer
			}

		case offTimer:
			if time.Now().After(offUntil) && on == 0 && off == len(uss) {
				log.Printf("reset  %s", fmtDists(last))
				st = watching
			}
		}
	}
}

func readConf() (*Config, error) {
	fh, err := os.Open(*conf)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	dec := yaml.NewDecoder(fh)
	dec.KnownFields(true)

	c := &Config{}

	err = dec.Decode(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func fmtDists(dists []float64) string {
	ints := make([]int, len(dists))

	for i, d := range dists {
		ints[i] = int(d)
	}

	b, err := json.Marshal(ints)
	if err != nil {
		panic(err)
	}

	return string(b)
}
