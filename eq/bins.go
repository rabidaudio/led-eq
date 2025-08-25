package eq

import (
	"fmt"
	"math"
)

type Bins interface {
	Len() int
	Bounds(binidx int) (float64, float64)
	weight(srclo, srchi, destlo, desthi float64) float64
}

type Linear struct {
	StartHz, StopHz, StepHz float64
}

func (lb Linear) Len() int {
	return int((lb.StopHz - lb.StartHz) / lb.StepHz)
}

func (lb Linear) Bounds(binidx int) (lo, hi float64) {
	if binidx >= lb.Len() || binidx < 0 {
		return -1, -1
	}
	return lb.StartHz + (lb.StepHz * float64(binidx)), lb.StartHz + (lb.StepHz * float64(binidx+1))
}

func (lb Linear) weight(srclo, srchi, destlo, desthi float64) float64 {
	return overlapRatio(srclo, srchi, destlo, desthi)
}

type Exponential struct {
	StartHz, StopHz float64
	NumBins         int
}

func (eb Exponential) Len() int {
	return eb.NumBins
}

func (eb Exponential) Bounds(binidx int) (lo, hi float64) {
	if binidx < 0 || binidx > eb.NumBins {
		return -1, -1
	}
	// TODO: cache?
	xstart := log(float64(eb.StartHz))
	xend := log(float64(eb.StopHz))
	xstep := (xend - xstart) / float64(eb.NumBins)

	xlo := xstart + float64(binidx)*xstep
	xhi := xstart + float64(binidx+1)*xstep

	// TODO: pre-compute for arbitrary bins

	// TODO: try power 2 or power e instead?
	return math.Pow(10, xlo), math.Pow(10, xhi)
}

func log(v float64) float64 {
	if v == 0 {
		return 0
	}
	return math.Log10(v)
}

func (eb Exponential) weight(srclo, srchi, destlo, desthi float64) float64 {
	// TODO: cache?
	return overlapRatio(log(srclo), log(srchi), log(destlo), log(desthi))
}

var _ Bins = Linear{}
var _ Bins = Exponential{}

func resample(fft []float64, src Bins, out []float64, dest Bins) {
	N := len(fft)
	B := len(out)
	if N < B {
		panic(fmt.Sprintf("bins: upsampling is not supported. N must be >= num bins. N: %v num bins: %v", N, B))
	}
	if _, ok := src.(Linear); !ok {
		panic("bins: src must be linear")
	}
	for i, v := range fft {
		srclo, srchi := src.Bounds(i)
		for j := range B {
			destlo, desthi := dest.Bounds(j)
			w := dest.weight(srclo, srchi, destlo, desthi)
			out[j] += w * v
		}
	}
}

func overlapRatio(srclo, srchi, destlo, desthi float64) float64 {
	// if S within D
	// D[ S[  ]S  ]D
	if srclo >= destlo && srchi <= desthi {
		return 1
	}
	// if D within S
	// S[ D[  ]D  ]S
	if destlo >= srclo && desthi <= srchi {
		return (desthi - destlo) / (srchi - srclo)
	}
	// if D starts in S
	// S[  D[   S]   D]
	if destlo >= srclo && destlo < srchi {
		return (srchi - destlo) / (srchi - srclo)
	}
	// if D ends within S
	// D[  S[   D]   S]
	if desthi >= srclo && desthi < srchi {
		return (desthi - srclo) / (srchi - srclo)
	}
	return 0 // disjoint
}
