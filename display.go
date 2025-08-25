package main

type Display interface {
	Render(values []float64) error
}
