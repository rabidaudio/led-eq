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

	N := NForTimeStep(wv.SampleRate(), 1*time.Second/60.0 /*60Hz*/, AtLeast)
	eq := NewEQ(wv.SampleRate(), 32, N)

	td := NewTerminalDisplay(&eq)
	wrap := EQStreamWrapper{Streamer: wv, eq: &eq, d: td}

	speaker.Init(wv.fmt.SampleRate, eq.N)

	go func() {
		speaker.Play(beep.Seq(&wrap, beep.Callback(func() {
			td.Done()
		})))
	}()

	td.Run()
}
