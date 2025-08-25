package main

import "github.com/faiface/beep"

type EQStreamWrapper struct {
	beep.Streamer
	eq *EQ
	d  Display

	buf  []float64
	bufi int
	res  []float64

	err error
}

var _ beep.Streamer = (*EQStreamWrapper)(nil)

func (sw *EQStreamWrapper) Stream(samples [][2]float64) (n int, ok bool) {
	if sw.buf == nil {
		sw.buf = make([]float64, 0, sw.eq.N())
		sw.bufi = 0
	}

	n, ok = sw.Streamer.Stream(samples)
	if !ok {
		return n, ok
	}
	// make sure there's enough room in the buffer
	rem := len(sw.buf) - sw.bufi
	if rem < n {
		sw.buf = append(sw.buf, make([]float64, n-rem)...)
	}
	// mono data and copy into buffer
	ToMono(samples[:n], sw.buf[sw.bufi:(sw.bufi+n)])
	sw.bufi += n

	// if a full N is available
	if sw.bufi >= sw.eq.N() {
		if sw.res == nil {
			sw.res = make([]float64, sw.eq.NumBins)
		}
		// compute and render
		if sw.d != nil {
			sw.eq.Compute(sw.buf[:sw.eq.N()], sw.res)
			err := sw.d.Render(sw.res)
			if err != nil {
				sw.err = err
				return n, false
			}
		}
		sw.buf = sw.buf[sw.eq.N():] // advance buffer
		sw.bufi -= sw.eq.N()
	}
	return n, ok
}

func (sw *EQStreamWrapper) Err() error {
	if sw.err != nil {
		return sw.err
	}
	return sw.Streamer.Err()
}
