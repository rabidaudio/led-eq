package eq

import (
	"fmt"
	"math"
	"math/cmplx"
	"time"

	"github.com/mjibson/go-dsp/fft"
)

type EQ struct {
	SampleRate int
	N          int
	OutBins    Bins
}

type StepMode int

const (
	AtLeast StepMode = iota
	AtMost
)

// Compute the optimal FFT size N based on the desired
// frequency at which you want measurements (e.g. 60Hz display).
// If [StepMode] is [AtLeast], will return the next largest power of two
// (meaning) time steps may be larger (and thus frequency may be slower).
// If [AtMost], will return the largest power of two within timestep, meaning
// steps may be smaller (and thus frequency may be faster).
func NForTimeStep(sampleRate int, step time.Duration, mode StepMode) int {
	targetN := int(step.Seconds() * float64(sampleRate))
	n := 1
	for n < targetN {
		n *= 2
	}
	if mode == AtMost {
		n /= 2 // went too far
	}
	return n
}

func DefaultEQ() EQ {
	return EQ{
		SampleRate: 44100,
		N:          2048,
		OutBins: Exponential{
			NumBins: 16,
			StartHz: 20,
			StopHz:  20_000,
		},
	}
}

func NewEQ(sampleRate, n, numbins int) EQ {
	eq := DefaultEQ()
	eq.SampleRate = sampleRate
	eq.N = n
	eq.OutBins = Exponential{
		NumBins: numbins,
		StartHz: 20,
		StopHz:  20_000,
	}
	return eq
}

func (eq *EQ) SrcBin() Bins {
	return Linear{StartHz: 0, StopHz: float64(eq.SampleRate), StepHz: float64(eq.SampleRate) / float64(eq.N)}
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
	if len(samples) < eq.N {
		panic(fmt.Errorf("eq: expected N=%v samples but was %v", eq.N, len(samples)))
	}
	if len(out) < eq.OutBins.Len() {
		panic(fmt.Errorf("eq: out must be at least len %v but was %v", eq.OutBins.Len(), len(out)))
	}

	// for i := range eq.OutBins.Len() {
	// 	out[i] = 0
	// }

	transformed := make([]float64, eq.N)
	realFFT(samples, transformed)

	// re-bin
	resample(transformed, eq.SrcBin(), out, eq.OutBins)

	for i := range eq.OutBins.Len() {
		out[i] /= float64(eq.OutBins.Len()) // normalize energy
	}
}
