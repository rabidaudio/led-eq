package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-audio/wav"
	"github.com/mjibson/go-dsp/fft"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const FreqBins = 16
const SampleRate = 44100

const LowBoundHz = 20
const HighBoundHz = 20_000

var Timestep = time.Second / 60 // 60Hz
var N int
var binSizeHz float64
var audibleBins int
var binSize int

func init() {
	targetN := Timestep.Seconds() * SampleRate
	N = 1
	for N < int(targetN) {
		N *= 2
	}

	binSizeHz = float64(SampleRate) / float64(N)
	audibleBins = int(math.Round(20_000 /*kHz*/ / binSizeHz))
	binSize = audibleBins / FreqBins
}

func main() {
	fmt.Printf("Sampling over %v sec time steps (N=%v) to %v audible bins, and averaging across %v to %v bins\n",
		Timestep.Seconds(), N, audibleBins, binSize, FreqBins)
	wg := sync.WaitGroup{}

	samples := make(chan float64, 0)
	lines := make(chan []float64, 0)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		readSamples("/Users/personal/projects/led-eq/lower_wolves.wav", samples)
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		x := make([]float64, N)

		xstart := math.Log10(LowBoundHz)
		xend := math.Log10(HighBoundHz)
		xstep := (xend - xstart) / FreqBins

		// lower frequency bounds (upper is x+1 or 20kHz)
		binBoundsHz := make([]float64, FreqBins+1)
		for i := range FreqBins {
			binBoundsHz[i] = math.Pow(10, xstart+float64(i)*xstep)
		}
		binBoundsHz[FreqBins] = HighBoundHz

		for {
			// load the next FreqBins samples
			for i := range N {
				sample, ok := <-samples
				if !ok {
					fmt.Printf("no more samples\n")
					close(lines)
					return // closed, finished
				}
				x[i] = sample
			}
			transformed := fft.FFTReal(x)

			l := make([]float64, FreqBins)
			counts := make([]int, FreqBins)

			for i := range transformed {
				v := cmplx.Abs(transformed[i]) / math.Sqrt(float64(N)) // normalized magnitude

				lbinhz := binSizeHz * float64(i)

				for j := range FreqBins {
					// TODO: could weighted average bins that fall across bounds, but for now we're just going
					// to assign completely to both
					if (lbinhz >= binBoundsHz[j] && lbinhz < binBoundsHz[j+1]) || // lower is within bin
						((lbinhz+binSizeHz) >= binBoundsHz[j] && (lbinhz+binSizeHz) < binBoundsHz[j+1]) || // upper is within
						(lbinhz <= binBoundsHz[j] && (lbinhz+binSizeHz) >= binBoundsHz[j+1]) { // covers the whole

						l[j] += v
						counts[j] += 1
					}
				}
			}

			// convert sums to averages
			for i := range FreqBins {
				if counts[i] != 0 {
					l[i] = l[i] / float64(counts[i])
				}
			}

			lines <- l
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// draw(fileSizeSamples, lines)
		f, err := os.Create("out.tsv")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		i := 0
		for {
			line, ok := <-lines
			if !ok {
				fmt.Printf("done writing tsv after %v lines\n", i)
				return
			}
			var sbuf []byte

			for _, v := range line {
				sbuf = strconv.AppendFloat(sbuf, v, 'f', -1, 64)
				sbuf = append(sbuf, '\t')
			}
			sbuf = append(sbuf, '\n')

			_, err = f.WriteString(string(sbuf))
			if err != nil {
				panic(err)
			}
			i += 1
		}
	}(&wg)

	wg.Wait()
}

func draw(fileSizeSamples int, lines chan []float64) {
	c := canvas.New(float64(fileSizeSamples/FreqBins), FreqBins)

	r := canvas.Rectangle(1, 1)

	x := 0
	for {
		line, ok := <-lines
		if !ok {
			break
		}

		for y, sample := range line {
			v := uint8(sample * 256)
			color := canvas.RGB(v, v, v)
			c.RenderPath(r, canvas.Style{Fill: canvas.Paint{Color: color}}, canvas.Identity.Translate(float64(x), float64(y)))
		}
		x += 1
	}

	// save file to png
	// img : = raster c.WriteFile("out.png", canvas)/

	err := c.WriteFile("out.png", renderers.PNG())
	if err != nil {
		panic(err)
	}
	// img := rasterizer.Draw(c, canvas.DefaultResolution, canvas.DefaultColorSpace)
	// img.
}

func readSamples(path string, samples chan float64) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// b := make([]byte, 4)
	// f.Read(b)
	// fmt.Printf("header: %v\n", string(b))
	// f.Seek(0, io.SeekStart)

	fmt.Printf("reading file...\n")
	wd := wav.NewDecoder(f)
	buf, err := wd.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	d, err := wd.Duration()
	if err != nil {
		panic(err)
	}
	nsamp := d.Seconds() * SampleRate
	fmt.Printf("expecting %v lines from %v samples\n", nsamp/float64(N), nsamp)

	fb := buf.AsFloatBuffer()
	for i := range len(fb.Data) / fb.Format.NumChannels {
		l := fb.Data[i*2] / float64(1<<15)
		r := fb.Data[(i*2)+1] / float64(1<<15)
		mono := (l + r) / 2
		samples <- mono
	}
	fmt.Printf("file complete\n")
	close(samples)

	// chunkSize := 512
	// p := make([]byte, chunkSize*2*2)
	// s := 0
	// for {
	// 	n := 0
	// 	done := false
	// 	for n < len(p) {
	// 		nn, err := f.Read(p[n:])
	// 		n += nn
	// 		if err == io.EOF {
	// 			fmt.Printf("file complete\n")
	// 			done = true
	// 			break
	// 		}
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// 	for i := range chunkSize {
	// 		var sample [2]float64
	// 		l := int16(binary.LittleEndian.Uint16(p[i*4 : i*4+2]))
	// 		r := int16(binary.LittleEndian.Uint16(p[i*4+2 : (i+1)*4]))
	// 		sample[0] = float64(l) / float64(1<<15)
	// 		sample[1] = float64(r) / float64(1<<15)
	// 		samples <- sample
	// 	}
	// 	s += n
	// 	if done {
	// 		fmt.Printf("closing input file after %v bytes\n", s)
	// 		close(samples)
	// 		return
	// 	}
	// }
}
