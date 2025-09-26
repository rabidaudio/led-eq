package main

import (
	"github.com/rabidaudio/led-eq/colorlight"
)

type ColorlightDisplay struct {
	cl     colorlight.Colorlight
	canvas colorlight.Canvas
}

func (cd *ColorlightDisplay) Open() error {
	err := cd.cl.Open("lo0")
	if err != nil {
		return err
	}
	// alloc canvas
	w, h := cd.cl.Dims()
	cd.canvas = make(colorlight.Canvas, h)
	for i := range w {
		cd.canvas[i] = make([]colorlight.Pixel, w)
	}
	return nil
}

func (cd *ColorlightDisplay) clear() {
	for _, row := range cd.canvas {
		for _, pix := range row {
			pix.Red = 0
			pix.Blue = 0
			pix.Green = 0
		}
	}
}

func (cd *ColorlightDisplay) setPixel(x, y int, color colorlight.Pixel) {
	cd.canvas[y][x] = color
}

func (cd *ColorlightDisplay) Render(values []float64) error {
	cd.clear()

	w, h := cd.cl.Dims()
	boxw := int(w) / len(values)
	// TODO: hsv
	for _, v := range values {
		boxh := int(float64(h) * v)
		for x := range boxw {
			for y := range boxh {
				cd.setPixel(x, y, colorlight.Pixel{Red: 255})
			}
		}
	}
	return cd.cl.Draw(cd.canvas)
}
