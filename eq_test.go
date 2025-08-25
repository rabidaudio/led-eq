package main

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func failIfErr(t *testing.T, err error) {
	assert.NoError(t, err)
	if err != nil {
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	eq := EQ{SampleRate: 48_000, NumBins: 3, Timestep: 1 * time.Millisecond}
	eq.initEQ()

	assert.True(t, eq.initialized)

	// timstamp target = 48
	assert.Equal(t, 32, eq.N)
	assert.Equal(t, float64(1500), eq.linearBinSizeHz)

	assert.InDeltaSlice(t, []float64{20, 200, 2000, 20_000}, eq.exBinBounds, 0.1)

	start, end := eq.BinBounds(1)
	assert.InDelta(t, 200, start, 0.1)
	assert.InDelta(t, 2000, end, 0.1)
}

func TestDefaultInit(t *testing.T) {
	eq := EQ{}
	eq.initEQ()

	assert.Equal(t, 44100, eq.SampleRate, "default")
	assert.Equal(t, 20, eq.LowerBoundHz, "default to audible range (20 Hz)")
	assert.Equal(t, 20_000, eq.UpperBoundHz, "default to audible range (20 KHz)")

	assert.Equal(t, 16, eq.NumBins, "default to 16 EQ bars")
	assert.Equal(t, 1*time.Second/60, eq.Timestep, "default 60Hz measurements (to match screen framerates)")

	assert.Equal(t, 512, eq.N, "the largest power of two that fits within Timestep")
}

func TestIsOverlapping(t *testing.T) {
	assert.True(t, isOverlapping(0, 43, 20, 30.8))   // b is completely within a
	assert.False(t, isOverlapping(0, 43, 974, 1500)) // completely disjoint
	assert.True(t, isOverlapping(15, 25, 10, 20))    // a starts within b
	assert.True(t, isOverlapping(15, 25, 20, 30))    // a ends within b
	assert.True(t, isOverlapping(5, 15, 10, 20))     // b starts within a
	assert.True(t, isOverlapping(15, 25, 10, 20))    // b ends within a
	assert.True(t, isOverlapping(10, 20, 10, 20))    // a == b
}

func TestSine(t *testing.T) {
	wv := WavReader{Path: "testdata/440sin.wav"}
	err := wv.Open()
	failIfErr(t, err)

	defer wv.Close()

	eq := EQ{SampleRate: wv.SampleRate(), NumBins: 16, Timestep: 50 * time.Millisecond}

	p := make([]float64, eq.N())
	n, err := wv.ReadMono(p)
	failIfErr(t, err)
	assert.Equal(t, n, len(p))

	out := make([]float64, eq.NumBins)
	eq.Compute(p, out)

	assert.Equal(t, eq.NumBins, len(out))

	for i, v := range out {
		s, e := eq.BinBounds(i)
		if s <= 440 && e > 440 {
			assert.InDelta(t, 0.8, v, 0.01)
		} else {
			assert.InDelta(t, 0, v, 0.01)
		}
	}
}

func TestEnergyConservation(t *testing.T) {
	wv := WavReader{Path: "testdata/440sin.wav"}
	err := wv.Open()
	failIfErr(t, err)

	defer wv.Close()

	eq := EQ{SampleRate: wv.SampleRate(), NumBins: 16, Timestep: 50 * time.Millisecond}

	p := make([]float64, eq.N())
	n, err := wv.ReadMono(p)
	failIfErr(t, err)
	assert.Equal(t, n, len(p))

	// wave file is 0.8 pp
	for _, v := range p {
		assert.InDelta(t, 0, v, 0.401)
	}

	var rms float64
	for _, v := range p {
		rms += (v * v)
	}
	rms /= float64(len(p))
	rms = math.Sqrt(rms)

	// RMS of a sin wave is srt(2)* peak value
	assert.InDelta(t, rms, 0.707*0.4, 0.01)

	out := make([]float64, eq.N())
	realFFT(p, out)

	var avg float64
	for _, v := range out {
		avg += v * v
	}
	avg /= float64(eq.N())
	avg = math.Sqrt(avg)
	assert.InDelta(t, 0.707*0.4, avg, 0.01)
}
