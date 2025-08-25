package main

import (
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
)

type WavReader struct {
	Path string

	fmt beep.Format
	f   *os.File
	s   beep.StreamSeekCloser
}

func (wv *WavReader) Open() error {
	f, err := os.Open(wv.Path)
	if err != nil {
		return err
	}
	wv.f = f
	wv.s, wv.fmt, err = wav.Decode(wv.f)
	return err
}

func (wv *WavReader) SampleRate() int {
	return int(wv.fmt.SampleRate)
}

func (wv *WavReader) NumChannels() int {
	return wv.fmt.NumChannels
}

func (wv *WavReader) LenSamples() int {
	return wv.s.Len()
}

func (wv *WavReader) Read(p [][2]float64) (n int, err error) {
	n, ok := wv.s.Stream(p)
	if !ok {
		return n, wv.s.Err()
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

func (wv *WavReader) Close() error {
	defer wv.f.Close()

	return wv.s.Close()
}
