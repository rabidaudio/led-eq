package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	wv, err := OpenWav(os.Stdin)
	if err != nil {
		panic(err)
	}

	// TODO: func N(timestep)
	eq := EQ{SampleRate: wv.SampleRate(), NumBins: 16, Timestep: 50 * time.Millisecond}

	td := NewTerminalDisplay(&eq)
	wrap := EQStreamWrapper{Streamer: wv, eq: &eq, d: td}

	speaker.Init(wv.fmt.SampleRate, eq.N())

	go func() {
		speaker.Play(beep.Seq(&wrap, beep.Callback(func() {
			td.Done()
		})))
	}()

	td.Run()
}
