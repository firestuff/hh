package hhio

import "sort"

func NewMedianFilter(in chan float64, num int) chan float64 {
	out := make(chan float64)

	go func() {
		defer close(out)

		buf := make([]float64, num, num)
		srt := make([]float64, num, num)
		next := 0

		for {
			v, ok := <-in
			if !ok {
				return
			}

			buf[next%num] = v
			next++

			if next < num {
				continue
			}

			copy(srt, buf)
			sort.Float64s(srt)

			out <- srt[num/2+1]
		}
	}()

	return out
}
