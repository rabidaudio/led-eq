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

	// var wtot float64
	for i, v := range transformed {
		lbinhzlow := eq.linearBinSizeHz * float64(i)
		lbinhzhigh := eq.linearBinSizeHz * float64(i+1)

		for j := range eq.NumBins {
			exbinhzlow, exbinhzhigh := eq.BinBounds(j)

			w := getWeight(lbinhzlow, lbinhzhigh, exbinhzlow, exbinhzhigh)
			// wtot += w
			out[j] += w * v
		}

		// wl := getWeight(lbinhzlow, lbinhzhigh, 0, float64(eq.LowerBoundHz))
		// wh := getWeight(lbinhzlow, lbinhzhigh, float64(eq.UpperBoundHz), float64(eq.SampleRate))
		// fmt.Printf("unused bin weights for [%v %v]:\t %v\t%v\n", lbinhzlow, lbinhzhigh, wl, wh)
	}

	// scale := wtot / float64(eq.N())
	// fmt.Printf("used weights: %v of %v (%v)\n", wtot, eq.N(), scale)

	for i := range eq.NumBins {
		out[i] /= float64(eq.NumBins) // normalize energy
		// out[i] /= (1 - scale)
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

func getWeight(linstart, linend, exstart, exend float64) float64 {
	if linstart >= exstart && linend <= exend {
		// case 1: linear contained within exponential
		return 1
	}
	// TODO: pre-compute during init
	la := math.Log10(linstart)
	if linstart == 0 {
		la = 0
	}
	lb := math.Log10(linend)
	ea := math.Log10(exstart)
	if exstart == 0 {
		ea = 0
	}
	eb := math.Log10(exend)
	if exstart >= linstart && exend <= linend {
		// case 2: exponential contained within linear
		return (eb - ea) / (lb - la)
	}
	if exstart >= linstart && exstart < linend {
		// case 3: exponential crosses end of linear
		return (lb - ea) / (lb - la)
	}
	if linstart >= exstart && linstart < exend {
		// case 4: exponential crosses start of linear
		return (eb - la) / (lb - la)
	}
	return 0 // disjoint
}

func rms(data []float64) float64 {
	var sum float64
	for _, d := range data {
		sum += d * d
	}
	avg := sum / float64(len(data))
	return math.Sqrt(avg)
}
