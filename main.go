package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/rabidaudio/led-eq/eq"
	"github.com/rabidaudio/led-eq/wav"
)

var debug = false

func init() {
	if v, ok := os.LookupEnv("DEBUG"); ok && v != "0" {
		debug = true
	}
}

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

// TODO: normalize: y axis log, clip above 0db

func main() {
	var wv *wav.WavReader
	if debug {
		wv = must(wav.OpenWavFile("eq/testdata/440sin.wav")).WavReader
	} else {
		wv = must(wav.OpenWav(os.Stdin))
	}

	N := eq.NForTimeStep(wv.SampleRate(), 1*time.Second/60.0 /*60Hz*/, eq.AtLeast)
	e := eq.EQ{
		SampleRate: wv.SampleRate(),
		N:          N,
		// OutBins:    eq.ExponentialBins(20, 20_000, 32),
		// OutBins: eq.LinearBins(0, float64(wv.SampleRate()), N),
		OutBins: eq.ArbitraryBins(
			50, 100, 200, 400, 800, 1600, 3200, 6400, 20_000,
			// 25, 50, 75, 100, 150, 200, 300, 400, 600, 800, 1200, 1600, 2400, 3200, 4800, 6400, 9600, 20_000,
		),
		Normalize: 2,
		OutputDB:  false,
	}

	speaker.Init(beep.SampleRate(wv.SampleRate()), e.N)

	var td *TerminalDisplay
	if !debug {
		td = NewTerminalDisplay(&e)
	}

	wrap := EQStreamWrapper{Streamer: wv, eq: &e, d: td}

	done := make(chan struct{})
	go func() {
		speaker.Play(beep.Seq(&wrap, beep.Callback(func() {
			if td != nil {
				td.Done()
			}
			done <- struct{}{}
		})))
	}()

	if td != nil {
		td.Run()
	} else {
		<-done
	}
}
