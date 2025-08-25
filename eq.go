package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"time"

	"github.com/mjibson/go-dsp/fft"
)

type EQ struct {
	NumBins      int
	SampleRate   int
	Timestep     time.Duration
	UpperBoundHz int
	LowerBoundHz int

	n               int
	initialized     bool
	linearBinSizeHz float64
	exBinBounds     []float64
	counts          []int
}

func (eq *EQ) initEQ() {
	if eq.NumBins == 0 {
		eq.NumBins = 16
	}
	if eq.LowerBoundHz == 0 {
		eq.LowerBoundHz = 20
	}
	if eq.UpperBoundHz == 0 {
		eq.UpperBoundHz = 20_000
	}
	if eq.SampleRate == 0 {
		eq.SampleRate = 44100
	}
	if eq.Timestep == 0 {
		eq.Timestep = 1 * time.Second / 60
	}
	if eq.n == 0 {
		targetN := eq.Timestep.Seconds() * float64(eq.SampleRate)
		eq.n = 1 // choose the largest power of 2 within the given timestep
		for {
			if eq.n > int(targetN) {
				eq.n /= 2 // went too far
				break
			} else {
				eq.n *= 2
			}
		}
		eq.linearBinSizeHz = float64(eq.SampleRate) / float64(eq.n)

		xstart := math.Log10(float64(eq.LowerBoundHz))
		xend := math.Log10(float64(eq.UpperBoundHz))
		xstep := (xend - xstart) / float64(eq.NumBins)
		// lower frequency bounds (upper is x+1 or 20kHz)
		eq.exBinBounds = make([]float64, eq.NumBins+1)
		for i := range eq.NumBins {
			eq.exBinBounds[i] = math.Pow(10, xstart+float64(i)*xstep)
		}
		eq.exBinBounds[eq.NumBins] = float64(eq.UpperBoundHz)
	}
	eq.counts = make([]int, eq.NumBins)
	eq.initialized = true
}

func (eq *EQ) N() int {
	if !eq.initialized {
		eq.initEQ()
	}
	return eq.n
}

func realFFT(samples []float64, out []float64) {
	transformed := fft.FFTReal(samples)
	for i, v := range transformed {
		out[i] = cmplx.Abs(v) / math.Sqrt(float64(len(samples))) // normalized magnitude
	}
}

// Compute takes in a slice of N mono samples, computes the FFT, determines the magnitude
// of each band, and averages into NumBins exponential buckets (out)
func (eq *EQ) Compute(samples []float64, out []float64) {
	if !eq.initialized {
		eq.initEQ()
	}

	if len(samples) < eq.N() {
		panic(fmt.Errorf("eq: expected N=%v samples but was %v", eq.N(), len(samples)))
	}
	if len(out) < eq.NumBins {
		panic(fmt.Errorf("eq: out must be at least len %v but was %v", eq.NumBins, len(out)))
	}

	for i := range eq.NumBins {
		out[i] = 0
		eq.counts[i] = 0
	}

	transformed := make([]float64, eq.N())
	realFFT(samples, transformed)

	for i, v := range transformed {
		lbinhzlow := eq.linearBinSizeHz * float64(i)
		lbinhzhigh := eq.linearBinSizeHz * float64(i+1)

		for j := range eq.NumBins {
			exbinhzlow, exbinhzhigh := eq.BinBounds(j)

			// for any linear bin that overlaps with an exponential bin, include it's measurement
			if isOverlapping(lbinhzlow, lbinhzhigh, exbinhzlow, exbinhzhigh) {
				out[j] += v
				eq.counts[j] += 1
			}
		}
	}

	// convert sums to averages
	for i := range eq.NumBins {
		if eq.counts[i] != 0 { // avoid divide by zero
			out[i] = out[i] / float64(eq.counts[i])
		}
	}
}

func (eq *EQ) BinBounds(i int) (start, end float64) {
	return eq.exBinBounds[i], eq.exBinBounds[i+1]
}

func isOverlapping(alo, ahi, blo, bhi float64) bool {
	// if a starts within b
	// B[   A[  ]A   ]B
	// B[   A[       ]B   ]A
	if alo >= blo && alo < bhi {
		return true
	}
	// if a ends within b
	// A[   B[  ]A   ]B
	// A[   B[       ]B   ]A
	if blo >= alo && blo < ahi {
		return true
	}
	return false
}
