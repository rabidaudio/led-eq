package main

import (
	"fmt"
	"time"

	alsa "github.com/cocoonlife/goalsa"
	"github.com/rabidaudio/led-eq/eq"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

// TODO: normalize: y axis log, clip above 0db

func DeInterleaveToMono(in []int16, out []float64) {
	if len(in) != 2*len(out) {
		panic(fmt.Sprintf("wrong size bufs: %v and %v\n", len(in), len(out)))
	}
	for i := range out {
		l := in[2*i]
		r := in[2*i+1]
		mono := (l + r) / 2
		out[i] = float64(mono) / float64(1<<15)
	}
}

func main() {
	// wv := must(wav.OpenWav(os.Stdin))

	dev := must(alsa.NewCaptureDevice("hw:1,1,0", 2, alsa.Format(alsa.FormatS16LE), 44100, alsa.BufferParams{}))

	N := eq.NForTimeStep(dev.Rate, 1*time.Second/60.0 /*60Hz*/, eq.AtLeast)
	e := eq.EQ{
		SampleRate: dev.Rate,
		N:          N,
		// OutBins:    eq.ExponentialBins(20, 20_000, 32),
		// OutBins: eq.LinearBins(0, float64(wv.SampleRate()), N),
		OutBins: eq.ArbitraryBins(
			50, 100, 200, 400, 800, 1600, 3200, 6400, 20_000,
		),
		Normalize: 2,
		OutputDB:  false,
	}

	td := NewTerminalDisplay(&e)

	go func() {
		read := make([]int16, 2*N)
		buf := make([]float64, N)
		out := make([]float64, e.OutBins.Len())
		for {
			// _, err := wv.ReadMono(buf)
			nn := 0
			for nn < 2*N {
				n := must(dev.Read(read[nn:]))
				nn += n
			}
			DeInterleaveToMono(read, buf)

			for i := range out {
				out[i] = 0
			}
			e.Compute(buf, out)
			err := td.Render(out)
			if err != nil {
				panic(err)
			}
		}
	}()

	td.Run()
}
