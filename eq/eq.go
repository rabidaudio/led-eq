package eq

import (
	"fmt"
	"math"
	"math/cmplx"
	"time"

	"github.com/mjibson/go-dsp/fft"
)

// TODO: window
// https://www.modalshop.com/rental/learn/basics/how-to-choose-fft-window

type EQ struct {
	SampleRate int
	N          int
	OutBins    Bins
	Normalize  float64
	OutputDB   bool
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
		n /= 2 // went 1 step too far
	}
	return n
}

func Default() EQ {
	return EQ{
		SampleRate: 44100,
		N:          2048,
		OutBins:    ExponentialBins(20, 20_000, 16),
	}
}

func New(sampleRate, n, numbins int) EQ {
	eq := Default()
	eq.SampleRate = sampleRate
	eq.N = n
	eq.OutBins = ExponentialBins(20, 20_000, numbins)
	return eq
}

// realFFT runs the fast fourier transform on real values samples, takes the
// magnitude of the results, and scales such that energy is conserved.
func realFFT(samples []float64, out []float64) {
	transformed := fft.FFTReal(samples)
	for i, v := range transformed {
		// https://dsp.stackexchange.com/questions/90327/how-to-normalize-the-fft
		out[i] = cmplx.Abs(v) / float64(len(samples)) // normalized magnitude
	}
}

// Compute takes in a slice of N mono samples, computes the FFT, determines the magnitude
// of each band, and averages into the format specified by [EQ.OutBins]
func (eq *EQ) Compute(samples []float64, out []float64) {
	if len(samples) < eq.N {
		panic(fmt.Errorf("eq: expected N=%v samples but was %v", eq.N, len(samples)))
	}
	if len(out) < eq.OutBins.Len() {
		panic(fmt.Errorf("eq: out must be at least len %v but was %v", eq.OutBins.Len(), len(out)))
	}

	fft := make([]float64, eq.N)
	realFFT(samples, fft)

	// re-bin
	src := LinearBins(0, float64(eq.SampleRate), eq.N)
	resample(src, eq.OutBins, fft, out)
	if eq.Normalize == 0 {
		eq.Normalize = 1
	}
	for i := range out {
		out[i] *= eq.Normalize
	}
	if eq.OutputDB {
		ToDB(out)
	}
}

func Resample(src, dest Bins, in, out []float64) {
	counts := make([]float64, len(out))
	// TODO: matrix multiplication
	for i, v := range in {
		for j := range out {
			w := weights(i, j, src, dest)
			counts[j] += w
			out[j] += w * v
		}
	}
	for i := range out {
		out[i] /= counts[i]
	}
}

func RMS(data []float64) float64 {
	var sum float64
	for _, d := range data {
		sum += d * d
	}
	avg := sum / float64(len(data))
	return math.Sqrt(avg)
}

func db(v float64) float64 {
	return 20 * math.Log10(v)
}

func ToDB(samples []float64) {
	for i := range samples {
		samples[i] = db(samples[i])
	}
}
