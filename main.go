package main

import (
	"io"
	"os"
	"time"
)

func main() {
	wv, err := OpenWav(os.Stdin)
	if err != nil {
		panic(err)
	}

	eq := EQ{SampleRate: wv.SampleRate(), NumBins: 16, Timestep: 50 * time.Millisecond}

	// duplicate streamer so we can also play back
	// a, b := beep.Dup(wv.s)
	// wv.s = &a

	// speaker.Init(wv.fmt.SampleRate, eq.N())

	td := NewTerminalDisplay(&eq)

	go func() {
		p := make([]float64, eq.N())
		res := make([]float64, eq.NumBins)
		for {
			_, err := wv.ReadMono(p)
			if err == io.EOF {
				td.Done()
				return
			}
			if err != nil {
				panic(err)
			}
			eq.Compute(p, res)
			td.Render(res)
		}
	}()

	td.Run()

	// done := make(chan bool)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	done <- true
	// })))
}
