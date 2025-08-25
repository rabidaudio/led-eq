package main

import (
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

func TestNForTimeStep(t *testing.T) {
	assert.Equal(t, 32, NForTimeStep(48_000, 1*time.Millisecond, AtMost) /* target 48*/)
	assert.Equal(t, 64, NForTimeStep(48_000, 1*time.Millisecond, AtLeast) /* target 48*/)
}

func TestInit(t *testing.T) {
	eq := DefaultEQ()
	eq.SampleRate = 48_000
	eq.NumBins = 3
	eq.N = 32
	eq.initEQ()

	assert.True(t, eq.initialized)

	assert.Equal(t, 32, eq.N)
	assert.Equal(t, float64(1500), eq.linearBinSizeHz)

	assert.InDeltaSlice(t, []float64{20, 200, 2000, 20_000}, eq.exBinBounds, 0.1)

	start, end := eq.BinBounds(1)
	assert.InDelta(t, 200, start, 0.1)
	assert.InDelta(t, 2000, end, 0.1)
}

func TestDefaultInit(t *testing.T) {
	eq := DefaultEQ()
	eq.initEQ()

	assert.Equal(t, 44100, eq.SampleRate, "default")
	assert.Equal(t, 20, eq.LowerBoundHz, "default to audible range (20 Hz)")
	assert.Equal(t, 20_000, eq.UpperBoundHz, "default to audible range (20 KHz)")

	assert.Equal(t, 16, eq.NumBins, "default to 16 EQ bars")

	assert.Equal(t, 2048, eq.N, "default")
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
	wv, err := OpenWavFile("testdata/440sin.wav")
	failIfErr(t, err)

	defer wv.Close()

	eq := NewEQ(wv.SampleRate(), 16, NForTimeStep(wv.SampleRate(), 50*time.Millisecond, AtMost))

	p := make([]float64, eq.N)
	n, err := wv.ReadMono(p)
	failIfErr(t, err)
	assert.Equal(t, n, len(p))

	out := make([]float64, eq.NumBins)
	eq.Compute(p, out)

	assert.Equal(t, eq.NumBins, len(out))

	t.Log(out)
	for i, v := range out {
		s, e := eq.BinBounds(i)
		// TODO: how to compute these expected magnitudes?
		if s <= 440 && e > 440 {
			assert.Greater(t, v, 0.8)
		} else {
			assert.InDelta(t, 0, v, 0.3)
		}
	}
}

func TestEnergyConservation(t *testing.T) {
	wv, err := OpenWavFile("testdata/440sin.wav")
	failIfErr(t, err)

	defer wv.Close()

	eq := NewEQ(wv.SampleRate(), 16, NForTimeStep(wv.SampleRate(), 50*time.Millisecond, AtMost))

	p := make([]float64, eq.N)
	n, err := wv.ReadMono(p)
	failIfErr(t, err)
	assert.Equal(t, n, len(p))

	// wave file is 0.8 pp
	for _, v := range p {
		assert.InDelta(t, 0, v, 0.401)
	}

	// RMS of a sin wave is srt(2)* peak value
	assert.InDelta(t, rms(p), 0.707*0.4, 0.01)

	out := make([]float64, eq.N)
	realFFT(p, out)

	assert.InDelta(t, 0.707*0.4, rms(out), 0.01)
}

func TestGetWeight(t *testing.T) {
	// l:[21.533203125  43.06640625]	e:[3556.558820077839. 5476.839268528712]	w: -6.367775273051676	v: -2.1567005054830806
	assert.GreaterOrEqual(t, getWeight(21.5, 43, 3556.6, 5476.8), float64(0))
	assert.LessOrEqual(t, getWeight(21.5, 43, 3556.6, 5476.8), float64(1))

	lbins := []float64{0, 3000, 6000, 9000, 12000, 15000, 18000, 21000, 24000, 27000, 30000, 33000, 36000, 39000, 42000, 45000, 48000}
	ebins := []float64{20, 200, 2000, 20000}

	results := make([]float64, 0)
	for l := range len(lbins) - 1 {
		for e := range len(ebins) - 1 {
			w := getWeight(lbins[l], lbins[l+1], ebins[e], ebins[e+1])
			results = append(results, w)
		}
	}

	t.Log(results)
	assert.InDeltaSlice(t, []float64{
		0.2875942272, 0.2875942272, 0.05064282956,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 0.6834904379,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}, results, 0.01)
}
