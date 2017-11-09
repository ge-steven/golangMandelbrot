/*
	Author: Steven Ge
	Date: 2017-11-09
*/

package main

import (
	"fmt"
	"math/cmplx"
	"time"

	"github.com/fogleman/gg"
)

func main() {
	// resolution := 0.000000005
	// maxX := float64(-0.78999)
	// maxY := float64(0.15001)
	// minX := float64(-0.79001)
	// minY := float64(0.14999)

	start := time.Now()

	//---------------------------------------------------------------------------
	// BEGIN CONSTANTS
	//---------------------------------------------------------------------------
	// Settings for the image
	// screenResolution := complex(float64(1080), float64(1920))
	resolution := 0.001
	maxX := float64(1)
	maxY := float64(1)
	minX := float64(-2)
	minY := float64(-1)

	// Amount of threads to be used
	maxComputingThreads := 4
	maxDrawingThreads := 4

	width := getFullCoordinates((maxX - minX), resolution)
	height := getFullCoordinates((maxY - minY), resolution)

	// Constanst for maximum iterations and modulusSize of the complex number
	maxIterations := 200
	maxModulusSize := float64(6)

	//---------------------------------------------------------------------------
	// END CONSTANTS
	//---------------------------------------------------------------------------

	// Making a channel for synchronization of the worker/computing threads
	computedValueChannel := make(chan workerValue)
	coordinatesValueChannel := make(chan CoordinateValue)
	// stopChannelFromMain := make(chan int)
	stopChannelFromPainter := make(chan int)
	stopChannelFromComputing := make(chan int)

	// Create the context for drawing
	dc := gg.NewContext(int(width), int(height))

	// Creating the threads for drawing and computing
	for i := 0; i < maxComputingThreads; i++ {
		go computingThread(maxModulusSize, maxIterations, coordinatesValueChannel, computedValueChannel, stopChannelFromComputing)
	}
	for i := 0; i < maxDrawingThreads; i++ {
		go drawingThread(computedValueChannel, *dc, maxIterations, minX, minY, resolution, stopChannelFromPainter)
	}

	// put all the coordinates in the channel
	for i := minX; i < maxX; i += resolution {
		for j := minY; j < maxY; j += resolution {
			value := new(CoordinateValue)
			value.coordinates = complex(float64(i), float64(j))
			value.stopSignalFromMain = false
			coordinatesValueChannel <- *value
		}
	}

	// When all the coordinates are in the channel, put the stopSignals in the channel
	for i := 0; i < maxComputingThreads; i++ {
		value := new(CoordinateValue)
		value.coordinates = complex(float64(i), float64(i))
		value.stopSignalFromMain = true
		coordinatesValueChannel <- *value
	}

	// Wait until all threads signal that they're done
	for i := 0; i < maxComputingThreads; i++ {
		<-stopChannelFromComputing
	}

	// Set a stop signal so that the drawingThread stops
	for i := 0; i < maxDrawingThreads; i++ {
		stopSignal := new(workerValue)
		stopSignal.stopSignalFromMain = true
		computedValueChannel <- *stopSignal
	}

	// Wait until the drawingThread signals that it's done
	for i := 0; i < maxDrawingThreads; i++ {
		<-stopChannelFromPainter
	}

	// Close all channels
	close(computedValueChannel)
	close(coordinatesValueChannel)
	// close(stopChannelFromMain)
	close(stopChannelFromPainter)
	close(stopChannelFromComputing)
	dc.SavePNG("HopefullyMandelbrotmultithreaded.png")

	t := time.Now()
	fmt.Print(t.Sub(start))
	fmt.Println()
}

// Compute whether the coordinates belong in the set or not.
// This calls the recursion
func computingThread(maxModulusSize float64, iterations int,
	coordinatesValueChannel chan CoordinateValue, computedValueChannel chan workerValue,
	stopChannelFromComputing chan int) {

	coordinates := new(CoordinateValue)

	for true {
		coordinates := <-coordinatesValueChannel

		if coordinates.stopSignalFromMain {
			break
		}

		computation(0, coordinates.coordinates, maxModulusSize, iterations, computedValueChannel)
	}

	stopChannelFromComputing <- int(real(coordinates.coordinates))
	return
}

// The recursion for determining whether the value belongs to the set or not
func computation(z complex128, coordinates complex128, maxModulusSize float64,
	iterations int, computedValueChannel chan workerValue) {
	modulus := cmplx.Abs(z)
	if iterations == 0 || modulus > maxModulusSize {
		result := new(workerValue)
		result.coordinates = coordinates
		result.iteration = iterations
		result.stopSignalFromMain = false

		computedValueChannel <- *result
		// fmt.Print(coordinates)
		// fmt.Print(<-ch)
		return
	}

	z = (z * z) + coordinates
	computation(z, coordinates, maxModulusSize, iterations-1, computedValueChannel)
}

// Draws the mandelbrot
// Gets the values from the computedValueChannel
func drawingThread(computedValueChannel chan workerValue, dc gg.Context, maxIterations int,
	minX float64, minY float64, resolution float64, stopChannelFromPainter chan int) {
	for true {
		value := <-computedValueChannel

		// If it is the stopSignal, break out of the loop
		if value.stopSignalFromMain {
			// Let the main thread know that we're done
			stopChannelFromPainter <- 1
			return
		}

		// draw the value given
		iterations := value.iteration
		coordinates := value.coordinates

		color := (float64(iterations) / float64(2)) / float64(maxIterations)

		dc.SetRGB(color, color, color)
		x := getFullCoordinates(real(coordinates)-float64(minX), resolution)
		y := getFullCoordinates(imag(coordinates)-float64(minY), resolution)

		dc.DrawPoint(x, y, float64(1))
		dc.Fill()
	}
}

// Compute the coordinates on the image (scale the coordinates up)
func getFullCoordinates(x float64, resolution float64) float64 {
	return (x / resolution)
}

// Zoom in at a certain point in the mandelbrot
func zoomInAt(coordinates complex128) {

}

/** The computed value from the computing thread
This is used by the drawing thread to determine the color and draw it
This is also used by the main thread to set the stop signal for the
drawing thread
*/
type workerValue struct {
	coordinates        complex128
	iteration          int
	stopSignalFromMain bool
}

/**
The coordinates that the computing threads use to determine whether
the coordinates belong to the set.
This is also used by the main thread to set a stop signal for the computing
threads.
*/

type CoordinateValue struct {
	coordinates        complex128
	stopSignalFromMain bool
}
