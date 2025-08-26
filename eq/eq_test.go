package eq

import (
	"fmt"
	"os"
	"slices"
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

func TestEnergyConservationSine(t *testing.T) {
	wv, err := wav.OpenWavFile("testdata/440sin.wav")
	failIfErr(t, err)

	defer wv.Close()

	N := NForTimeStep(wv.SampleRate(), 50*time.Millisecond, AtMost)
	eq := New(wv.SampleRate(), N, 16)

	p := make([]float64, N)
	n, err := wv.ReadMono(p)
	failIfErr(t, err)
	assert.Equal(t, n, len(p))

	// wave file is 0.8 pp
	for _, v := range p {
		assert.InDelta(t, 0, v, 0.401)
	}

	// RMS of a sin wave is srt(2)* peak value
	assert.InDelta(t, RMS(p), 0.707*0.4, 0.01)

	out := make([]float64, N)
	realFFT(p, out)

	assert.InDelta(t, 0.707*0.4, RMS(out), 0.01)

	f, err := os.Create("out.txt")
	failIfErr(t, err)
	defer f.Close()

	// rms should be approximately for binning conversions regardless of bin size
	for i := 1; i < N; i += 1 {
		eq.OutBins = LinearBins(0, float64(wv.SampleRate()), i)
		out = make([]float64, i)
		eq.Compute(p, out)
		r := RMS(out)

		fmt.Fprintf(f, "%v\t%v\n", i, r)
		// t.Log(i, r)
		// assert.InDelta(t, 0.707*0.4, r, 0.1)
	}
}

func TestNormalizedSinValue(t *testing.T) {
	wv, err := wav.OpenWavFile("testdata/440sin_1.wav")
	failIfErr(t, err)

	defer wv.Close()

	// f, err := os.Create("out.txt")
	// failIfErr(t, err)
	// defer f.Close()

	wavdata := make([]float64, 0, 4*wv.SampleRate())
	chunk := make([][2]float64, 512)
	chunkMono := make([]float64, 512)
	for {
		n, ok := wv.Stream(chunk)
		wav.ToMono(chunk[:n], chunkMono[:n])
		wavdata = append(wavdata, chunkMono[:n]...)
		if !ok {
			break
		}
	}

	// t.Log("read file\n")

	// for n := 0; n < 11; n += 1 {
	// N := (1 << n)
	avg := 0.0
	navg := 0.0
	for N := 1; N < 2048; N += 1 {

		// for b := 0; b < 6; b += 1 {
		for b := 5; b < 6; b += 1 {
			B := (1 << b)
			eq := EQ{SampleRate: wv.SampleRate(), N: N, OutBins: ExponentialBins(20, 20_000, B)}

			t.Log(N, B)

			max := 0.0
			c := 0
			n := 0
			out := make([]float64, B)
			for n < len(wavdata) {
				end := n + N
				if end > len(wavdata) {
					break
				}
				eq.Compute(wavdata[n:end], out)
				max = slices.Max(out)
				n += N
				c += 1
			}
			// a /= float64(c)
			// fmt.Fprintf(f, "%d\t%d\t%f\n", N, B, max)
			avg += max
			navg += 1
		}
	}
	avg /= navg
	t.Log(avg)
	assert.InEpsilon(t, Normalization440HzSin, 1/avg, 0.01)
}

func avg(p []float64) float64 {
	a := 0.0
	for i := range p {
		a += p[i]
	}
	return a / float64(len(p))
}

// TODO: sweep frequency range and see how RMS changes
