package main

import (
	"io"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
)

type WavReader struct {
	fmt beep.Format
	r   io.Reader
	s   beep.Streamer
}

func OpenWav(r io.Reader) (wv *WavReader, err error) {
	wv = &WavReader{}
	wv.r = r
	wv.s, wv.fmt, err = wav.Decode(wv.r)
	return
}

func (wv *WavReader) SampleRate() int {
	return int(wv.fmt.SampleRate)
}

func (wv *WavReader) NumChannels() int {
	return wv.fmt.NumChannels
}

// func (wv *WavReader) LenSamples() int {
// 	return wv.s.Len()
// }

func (wv *WavReader) Read(p [][2]float64) (n int, err error) {
	n, ok := wv.s.Stream(p)
	if !ok {
		err = wv.s.Err()
		if err != nil {
			return n, err
		}
		return n, io.EOF
	}
	return n, nil
}

func (wv *WavReader) ReadMono(p []float64) (n int, err error) {
	pp := make([][2]float64, len(p)) // TODO: avoid additional alloc
	n, err = wv.Read(pp)
	for i := range n {
		p[i] = (pp[i][0] + pp[i][1]) / 2
	}
	return n, err
}

// func (wv *WavReader) Close() error {
// 	defer wv.f.Close()

// 	return wv.s.Close()
// }
