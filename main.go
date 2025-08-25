package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/rabidaudio/led-eq/eq"
	"github.com/rabidaudio/led-eq/wav"
)

func main() {
	wv, err := wav.OpenWav(os.Stdin)
	// wv, err := wav.OpenWavFile("eq/testdata/440sin.wav")
	if err != nil {
		panic(err)
	}

	N := eq.NForTimeStep(wv.SampleRate(), 1*time.Second/60.0 /*60Hz*/, eq.AtLeast)
	eq := eq.NewEQ(wv.SampleRate(), N, 20)

	speaker.Init(beep.SampleRate(wv.SampleRate()), eq.N)

	// wrap := EQStreamWrapper{Streamer: wv, eq: &eq}
	// done := make(chan struct{})
	// speaker.Play(beep.Seq(&wrap, beep.Callback(func() {
	// 	done <- struct{}{}
	// })))
	// <-done

	td := NewTerminalDisplay(&eq)
	wrap := EQStreamWrapper{Streamer: wv, eq: &eq, d: td}

	go func() {
		speaker.Play(beep.Seq(&wrap, beep.Callback(func() {
			td.Done()
		})))
	}()

	td.Run()
}
