package main

import "time"

func main() {
	wv := WavReader{Path: "testdata/lower_wolves.wav"}
	err := wv.Open()
	if err != nil {
		panic(err)
	}

	// duplicate streamer so we can also play back
	// a, b := beep.Dup(wv.s)
	// wv.s = &a

	// speaker.Init(wv.fmt.SampleRate, wv.fmt.SampleRate.N(time.Second/10))

	defer wv.Close()

	eq := EQ{SampleRate: wv.SampleRate(), NumBins: 16, Timestep: 50 * time.Millisecond}

	p := make([]float64, eq.N())
	td := NewTerminalDisplay(&eq, func(res []float64) {
		_, err := wv.ReadMono(p)
		if err != nil {
			panic(err)
		}
		eq.Compute(p, res)
	})

	td.Run()

	// done := make(chan bool)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	done <- true
	// })))
}
