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

func main() {
	var wv *wav.WavReader
	if debug {
		wv = must(wav.OpenWavFile("eq/testdata/440sin.wav")).WavReader
	} else {
		wv = must(wav.OpenWav(os.Stdin))
	}

	N := eq.NForTimeStep(wv.SampleRate(), 1*time.Second/60.0 /*60Hz*/, eq.AtLeast)
	eq := eq.New(wv.SampleRate(), N, 32)

	speaker.Init(beep.SampleRate(wv.SampleRate()), eq.N)

	var td *TerminalDisplay
	if !debug {
		td = NewTerminalDisplay(&eq)
	}

	wrap := EQStreamWrapper{Streamer: wv, eq: &eq, d: td}

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
