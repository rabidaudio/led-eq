package eq

import (
	"math"
)

type Bins []float64

func (b Bins) Len() int {
	return len(b) - 1
}

func (b Bins) Bounds(i int) (float64, float64) {
	if i < 0 || i >= b.Len() {
		return -1, -1
	}
	return b[i], b[i+1]
}

func LinearBins(start, stop float64, len int) Bins {
	step := (stop - start) / float64(len)
	b := make([]float64, len+1)
	for i := range len + 1 {
		b[i] = start + float64(i)*step
	}
	return b
}

func ExponentialBins(start, stop float64, len int) Bins {
	xstart := log(start)
	xend := log(stop)
	xstep := (xend - xstart) / float64(len)

	b := make([]float64, len+1)
	for i := range len + 1 {
		x := xstart + float64(i)*xstep
		b[i] = math.Exp(x)
	}
	return b
}

func ArbitraryBins(bounds ...float64) Bins {
	return bounds
}

func log(v float64) float64 {
	if v == 0 {
		return 0
	}
	return math.Log(v)
}

func weights(i, j int, src, dest Bins) float64 {
	slo, shi := src.Bounds(i)
	dlo, dhi := dest.Bounds(j)
	return overlapRatio(slo, shi, dlo, dhi) // TODO: normalize transfer?  / (dhi - dlo)
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
