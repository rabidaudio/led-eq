package main

import (
	"fmt"
	"math"

	"github.com/crazy3lf/colorconv"
	"github.com/stianeikeland/go-rpio/v4"
)

type RGB struct{ R, G, B uint8 }

var white = RGB{255, 255, 255}

type SpiCommand uint8

const (
	WriteDisplay SpiCommand = 0xA0
)

type Surface struct {
	W, H uint8
	data []RGB
}

func NewSurface(w, h uint8) Surface {
	return Surface{W: w, H: h, data: make([]RGB, w*h)}
}

func (s *Surface) Clear() {
	for i := range s.data {
		s.data[i].R = 0
		s.data[i].G = 0
		s.data[i].B = 0
	}
}

func (s *Surface) Set(x, y uint8, v RGB) {
	if x >= s.W || y >= s.H {
		panic(fmt.Errorf("surface: %v,%v out of bounds for surface %v by %v", x, y, s.W, s.H))
	}
	s.data[y*s.W+x].R = v.R
	s.data[y*s.W+x].G = v.G
	s.data[y*s.W+x].B = v.B
}

func (s *Surface) DrawRect(x, y, w, h uint8, v RGB) {
	for i := x; i < x+w && i < s.W; i += 1 {
		for j := y; j < y+h && j < s.H; j += 1 {
			s.Set(i, j, v)
		}
	}
}

type LED struct {
	BinWidth   uint8
	Surface    Surface
	ColorSpeed int
	step       int
}

var _ Display = (*LED)(nil)

func (l *LED) Open(dev rpio.SpiDev, cs uint8) error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	err = rpio.SpiBegin(dev)
	if err != nil {
		return err
	}
	rpio.SpiChipSelect(cs)
	rpio.SpiSpeed(1_000_000 /* 1 MHz */)

	return nil
}

func (l *LED) Render(data []float64) error {
	if len(data) != int(l.BinWidth) {
		return fmt.Errorf("led: invalid width, expected %v but was %v", l.BinWidth, len(data))
	}

	// compute screen
	l.Surface.Clear()
	for i, v := range data {
		vv := vscale(v, l.Surface.H)
		// l.Surface.DrawRect(i*int(l.BinWidth), 0, l.BinWidth, vv, white)
		var x, y uint8
		for x = uint8(i) * l.BinWidth; x < uint8(i+1)*l.BinWidth && x < l.Surface.W; x += 1 {
			for y = 0; y < vv && y < l.Surface.H; y += 1 {
				l.Surface.Set(x, y, l.color(x, y))
			}
		}
	}

	l.step += 1
	if l.step > 360 {
		l.step = 0
	}

	// write screen
	npix := int(l.Surface.W) * int(l.Surface.H)
	buf := make([]byte, npix*3+1)
	buf[0] = byte(WriteDisplay)
	for i := range npix {
		buf[(i*3)+1] = l.Surface.data[i].R
		buf[(i*3)+2] = l.Surface.data[i].G
		buf[(i*3)+3] = l.Surface.data[i].B
	}
	rpio.SpiTransmit(buf...)

	return nil
}

func (l *LED) color(x, y uint8) RGB {
	baseHue := float64(l.step) / float64(l.ColorSpeed)
	// adjust lightness from 40 to 75
	// 75% lightness at top right [l.W, l.H]
	radius := math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y))
	ratio := (radius / math.Sqrt(float64(l.Surface.W)*float64(l.Surface.W)+float64(l.Surface.H)*float64(l.Surface.H)))
	lightness := .4 + (.74-.4)*ratio
	hue := baseHue + ratio*15.0
	if hue > 360 {
		hue -= 360
	}
	r, g, b, err := colorconv.HSLToRGB(hue, 1.0, lightness)
	if err != nil {
		panic(err)
	}
	return RGB{r, g, b}
}

func vscale(v float64, h uint8) uint8 {
	return uint8(float64(h) * v)
}

func (*LED) Close() error {
	return rpio.Close()
}
