package wav

import (
	"io"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
)

func ToMono(p [][2]float64, out []float64) {
	for i := range p {
		out[i] = (p[i][0] + p[i][1]) / 2
	}
}

type WavReader struct {
	beep.Streamer

	fmt beep.Format
	r   io.Reader
}

func OpenWav(r io.Reader) (wv *WavReader, err error) {
	wv = &WavReader{}
	wv.r = r
	wv.Streamer, wv.fmt, err = wav.Decode(wv.r)
	return wv, err
}

func (wv *WavReader) SampleRate() int {
	return int(wv.fmt.SampleRate)
}

func (wv *WavReader) NumChannels() int {
	return wv.fmt.NumChannels
}

func (wv *WavReader) Read(p [][2]float64) (n int, err error) {
	n, ok := wv.Stream(p)
	if !ok {
		err = wv.Streamer.Err()
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
	ToMono(pp, p)
	return n, err
}

type WaveFileReader struct {
	*WavReader
	f    *os.File
	Path string
}

func OpenWavFile(path string) (*WaveFileReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	wv, err := OpenWav(f)
	if err != nil {
		return nil, err
	}
	return &WaveFileReader{WavReader: wv, Path: path, f: f}, nil
}

func (wv *WaveFileReader) LenSamples() int {
	s := wv.Streamer.(beep.StreamSeekCloser)
	return s.Len()
}

func (wv *WaveFileReader) Close() error {
	s := wv.Streamer.(beep.StreamSeekCloser)
	defer s.Close()

	return wv.f.Close()
}
