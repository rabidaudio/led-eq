package eq

import (
	"math"
)

type Bins struct {
	len           int
	bounds        []float64
	isExponential bool
}

func (b *Bins) Len() int {
	return b.len
}

func (b *Bins) Bounds(i int) (float64, float64) {
	if i < 0 || i >= b.len {
		return -1, -1
	}
	return b.bounds[i], b.bounds[i+1]
}

func LinearBins(start, stop float64, len int) Bins {
	step := (stop - start) / float64(len)
	b := make([]float64, len+1)
	for i := range len + 1 {
		b[i] = start + float64(i)*step
	}
	return Bins{
		len:           len,
		bounds:        b,
		isExponential: false,
	}
}

func ExponentialBins(start, stop float64, len int) Bins {
	xstart := log(start)
	xend := log(stop)
	xstep := (xend - xstart) / float64(len)

	b := make([]float64, len+1)
	for i := range len + 1 {
		x := xstart + float64(i)*xstep
		b[i] = math.Pow(10, x)
	}
	return Bins{
		len:           len,
		bounds:        b,
		isExponential: true,
	}
}

func ArbitraryBins(bounds ...float64) Bins {
	return Bins{
		len:           len(bounds) - 1,
		bounds:        bounds,
		isExponential: false,
	}
}

func (b Bins) toExponential() Bins {
	bb := make([]float64, b.len+1)
	for i := range b.len + 1 {
		bb[i] = log(b.bounds[i])
	}
	return Bins{
		len:    b.len,
		bounds: bb,
	}
}

func log(v float64) float64 {
	if v == 0 {
		return 0
	}
	return math.Log10(v) // TODO: log2? logE?
}

func weight(i, j int, src, dest Bins) float64 {
	if dest.isExponential {
		// operate in exponential mode
		src = src.toExponential()
		dest = dest.toExponential()
	}
	slo, shi := src.Bounds(i)
	dlo, dhi := dest.Bounds(j)
	return overlapRatio(slo, shi, dlo, dhi)
}

func resample(fft []float64, src Bins, out []float64, dest Bins) {
	// TODO: matrix multiplication
	for i, v := range fft {
		for j := range out {
			w := weight(i, j, src, dest)
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
