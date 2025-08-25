package eq

import (
	"math"
	"testing"
	"time"

	"github.com/rabidaudio/led-eq/wav"
	"github.com/stretchr/testify/assert"
)

func failIfErr(t *testing.T, err error) {
	assert.NoError(t, err)
	if err != nil {
		t.Fail()
	}
}

func rms(data []float64) float64 {
	var sum float64
	for _, d := range data {
		sum += d * d
	}
	avg := sum / float64(len(data))
	return math.Sqrt(avg)
}

func TestNForTimeStep(t *testing.T) {
	assert.Equal(t, 32, NForTimeStep(48_000, 1*time.Millisecond, AtMost) /* target 48*/)
	assert.Equal(t, 64, NForTimeStep(48_000, 1*time.Millisecond, AtLeast) /* target 48*/)
}

func TestNew(t *testing.T) {
	eq := New(48_000, 32, 3)

	assert.Equal(t, 48_000, eq.SampleRate)
	assert.Equal(t, 32, eq.N)

	start, end := eq.OutBins.Bounds(1)
	assert.InDelta(t, 200, start, 0.1)
	assert.InDelta(t, 2000, end, 0.1)
}

func TestDefaultInit(t *testing.T) {
	eq := Default()

	assert.Equal(t, 44100, eq.SampleRate, "default")
	lo, _ := eq.OutBins.Bounds(0)
	assert.InEpsilon(t, 20, lo, 0.001, "default to audible range (20 Hz)")
	_, hi := eq.OutBins.Bounds(eq.OutBins.Len() - 1)
	assert.InEpsilon(t, 20_000, hi, 0.001, "default to audible range (20 KHz)")

	assert.Equal(t, 16, eq.OutBins.Len(), "default to 16 EQ bars")

	assert.Equal(t, 2048, eq.N, "default")
}

func TestSine(t *testing.T) {
	wv, err := wav.OpenWavFile("testdata/440sin.wav")
	failIfErr(t, err)

	defer wv.Close()

	eq := New(wv.SampleRate(), NForTimeStep(wv.SampleRate(), 50*time.Millisecond, AtMost), 16)

	p := make([]float64, eq.N)
	n, err := wv.ReadMono(p)
	failIfErr(t, err)
	assert.Equal(t, n, len(p))

	out := make([]float64, eq.OutBins.Len())
	eq.Compute(p, out)

	assert.Equal(t, eq.OutBins.Len(), len(out))

	t.Log(out)
	for i, v := range out {
		s, e := eq.OutBins.Bounds(i)
		// TODO: how to compute these expected magnitudes?
		if s <= 440 && e > 440 {
			assert.Greater(t, v, 0.8)
		} else {
			assert.InDelta(t, 0, v, 0.3)
		}
	}
}

func TestEnergyConservation(t *testing.T) {
	wv, err := wav.OpenWavFile("testdata/440sin.wav")
	failIfErr(t, err)

	defer wv.Close()

	eq := New(wv.SampleRate(), NForTimeStep(wv.SampleRate(), 50*time.Millisecond, AtMost), 16)

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
