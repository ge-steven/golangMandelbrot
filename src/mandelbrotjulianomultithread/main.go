package main

import (
	"fmt"
	"math/cmplx"
	"time"

	"github.com/fogleman/gg"
)

func main() {

	start := time.Now()
	// resolution := 0.000000005
	// maxX := float64(-0.78999)
	// maxY := float64(0.15001)
	// minX := float64(-0.79001)
	// minY := float64(0.14999)

	resolution := 0.001
	maxX := float64(1)
	maxY := float64(1)
	minX := float64(-2)
	minY := float64(-1)

	width := getFullCoordinates((maxX - minX), resolution)
	height := getFullCoordinates((maxY - minY), resolution)

	// Constanst for maximum iterations and modulusSize of the complex number
	maxIterations := 200
	maxModulusSize := float64(6)

	// Create the context for drawing
	dc := gg.NewContext(int(width), int(height))

	for i := minX; i < maxX; i += resolution {
		for j := minY; j < maxY; j += resolution {
			// Create the list of coordinates for one thread
			coordinates := complex(i, j)

			iterations := computation(complex(0, 0), coordinates, maxModulusSize, maxIterations)

			color := (float64(iterations) / float64(2)) / float64(maxIterations)

			dc.SetRGB(color, color, color)
			x := getFullCoordinates(real(coordinates)-float64(minX), resolution)
			y := getFullCoordinates(imag(coordinates)-float64(minY), resolution)

			dc.DrawPoint(x, y, float64(1))
			dc.Fill()
		}
	}
	dc.SavePNG("HopefullyMandelbrot.png")

	t := time.Now()
	fmt.Print(t.Sub(start))
	fmt.Println()
}

func computation(z complex128, coordinates complex128, maxModulusSize float64,
	iterations int) int {
	modulus := cmplx.Abs(z)
	if iterations == 0 || modulus > maxModulusSize {
		// fmt.Print(coordinates)
		// fmt.Print(<-ch)
		return iterations
	}

	z = (z * z) + coordinates
	return computation(z, coordinates, maxModulusSize, iterations-1)
}

// Compute the coordinates on the image (scale the coordinates up)
func getFullCoordinates(x float64, resolution float64) float64 {
	return (x / resolution)
}

// Zoom in at a certain point in the mandelbrot
func zoomInAt(coordinates complex128) {

}
