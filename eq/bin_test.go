package eq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearBins(t *testing.T) {
	b := LinearBins(100, 1000, 9)

	assert.Equal(t, 9, b.Len())
	x, y := b.Bounds(0)
	assert.Equal(t, 100.0, x)
	assert.Equal(t, 200.0, y)

	x, y = b.Bounds(8)
	assert.Equal(t, 900.0, x)
	assert.Equal(t, 1000.0, y)

	x, y = b.Bounds(9)
	assert.Equal(t, -1.0, x)
	assert.Equal(t, -1.0, y)

	x, y = b.Bounds(-1)
	assert.Equal(t, -1.0, x)
	assert.Equal(t, -1.0, y)
}

func TestExponentialBins(t *testing.T) {
	b := ExponentialBins(20, 20_000, 3)

	assert.Equal(t, 3, b.Len())
	x, y := b.Bounds(0)
	assert.InEpsilon(t, 20.0, x, 0.0001)
	assert.InEpsilon(t, 200.0, y, 0.0001)

	x, y = b.Bounds(1)
	assert.InEpsilon(t, 200.0, x, 0.0001)
	assert.InEpsilon(t, 2000.0, y, 0.0001)

	x, y = b.Bounds(2)
	assert.InEpsilon(t, 2000.0, x, 0.0001)
	assert.InEpsilon(t, 20_000.0, y, 0.0001)

	x, y = b.Bounds(3)
	assert.InEpsilon(t, -1.0, x, 0.0001)
	assert.InEpsilon(t, -1.0, y, 0.0001)

	x, y = b.Bounds(-1)
	assert.InEpsilon(t, -1.0, x, 0.0001)
	assert.InEpsilon(t, -1.0, y, 0.0001)
}

func TestArbitraryBins(t *testing.T) {
	b := ArbitraryBins(20, 100, 250, 500, 1000, 2500, 20_000)

	assert.Equal(t, 6, b.Len())
	x, y := b.Bounds(0)
	assert.Equal(t, 20.0, x)
	assert.Equal(t, 100.0, y)

	x, y = b.Bounds(1)
	assert.Equal(t, 100.0, x)
	assert.Equal(t, 250.0, y)

	x, y = b.Bounds(5)
	assert.Equal(t, 2500.0, x)
	assert.Equal(t, 20_000.0, y)

	x, y = b.Bounds(6)
	assert.Equal(t, -1.0, x)
	assert.Equal(t, -1.0, y)

	x, y = b.Bounds(-1)
	assert.Equal(t, -1.0, x)
	assert.Equal(t, -1.0, y)
}

func TestOverlapRatio(t *testing.T) {
	// l:[21.533203125  43.06640625]	e:[3556.558820077839. 5476.839268528712]	w: -6.367775273051676	v: -2.1567005054830806
	assert.GreaterOrEqual(t, overlapRatio(21.5, 43, 3556.6, 5476.8), float64(0))
	assert.LessOrEqual(t, overlapRatio(21.5, 43, 3556.6, 5476.8), float64(1))

	assert.True(t, overlapRatio(0, 43, 20, 30.8) > 0)   // b is completely within a
	assert.False(t, overlapRatio(0, 43, 974, 1500) > 0) // completely disjoint
	assert.True(t, overlapRatio(15, 25, 10, 20) > 0)    // a starts within b
	assert.True(t, overlapRatio(15, 25, 20, 30) > 0)    // a ends within b
	assert.True(t, overlapRatio(5, 15, 10, 20) > 0)     // b starts within a
	assert.True(t, overlapRatio(15, 25, 10, 20) > 0)    // b ends within a
	assert.True(t, overlapRatio(10, 20, 10, 20) > 0)    // a == b
}

func TestLinearConversion(t *testing.T) {
	a := LinearBins(0, 24, 3)  // 0 8 16 24
	b := LinearBins(12, 24, 4) // 12 15 18 21 24

	//        00  01  02  03  04  05  06  07  08  09  10  11  12  13  14  15  16  17  18  19  20  21  22  23  24
	//src:   [a                             ][b                             ][c                                 ]
	//dest:                                                  [d         ][e         ][f         ][g             ]
	//
	//     d   e   f   g
	// a
	// b  3/8 1/8
	// c      2/8 3/8 3/8
	results := make([]float64, 0)
	for i := range a.Len() {
		for j := range b.Len() {
			w := weights(i, j, a, b)
			results = append(results, w)
		}
	}
	t.Log(results)
	assert.InDeltaSlice(t, []float64{
		0, 0, 0, 0,
		3.0 / 8, 1.0 / 8, 0, 0,
		0, 2.0 / 8, 3.0 / 8, 3.0 / 8,
	}, results, 0.01)
}

func TestExponentialConversion(t *testing.T) {
	lbins := LinearBins(0, 48_000, 16)
	ebins := ExponentialBins(20, 20_000, 3)

	results := make([]float64, 0)
	for l := range lbins.Len() {
		for e := range ebins.Len() {
			w := weights(l, e, lbins, ebins)
			results = append(results, w)
		}
	}

	t.Log(results)
	assert.InDeltaSlice(t, []float64{
		0.2875942272, 0.2875942272, 0.05064282956,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 0.6834904379,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}, results, 0.01)
}

func TestUpscaleConversion(t *testing.T) {
	a := LinearBins(12, 24, 4) // 12 15 18 21 24
	b := LinearBins(0, 24, 3)  // 0 8 16 24

	//        00  01  02  03  04  05  06  07  08  09  10  11  12  13  14  15  16  17  18  19  20  21  22  23  24
	//src:                                                   [d         ][e         ][f         ][g             ]
	//dest:  [a                             ][b                             ][c                                 ]
	//
	//     a    b    c
	// d        1
	// e       1/3  2/3
	// f             1
	// g             1
	results := make([]float64, 0)
	for i := range a.Len() {
		for j := range b.Len() {
			w := weights(i, j, a, b)
			results = append(results, w)
		}
	}
	t.Log(results)
	assert.InDeltaSlice(t, []float64{
		0, 1.0, 0,
		0, 1.0 / 3, 2.0 / 3,
		0, 0, 1.0,
		0, 0, 1.0,
	}, results, 0.01)
}

func TestArbitraryConversion(t *testing.T) {
	a := LinearBins(0, 1000, 10)
	b := ArbitraryBins(0, 100, 250, 500, 2500)

	results := make([]float64, 0)
	for i := range a.Len() {
		for j := range b.Len() {
			w := weights(i, j, a, b)
			results = append(results, w)
		}
	}
	t.Log(results)
	assert.InDeltaSlice(t, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0.5, 0.5, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
	}, results, 0.01)
}
