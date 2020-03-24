/*
	Author: Steven Ge
	Date: 2020-03-20
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
	"time"
)

// Settings for the image
var width = 3000
var height = 2000
var scale = 2.0
var centerX = -0.5
var centerY = 0.0
var upLeft image.Point
var lowRight image.Point
var min image.Point
var max image.Point
var redFactor = 1
var greenFactor = 1
var blueFactor = 1

// Amount of threads to be used
var maxComputingThreads = runtime.NumCPU()

// Constanst for maximum iterations and modulusSize of the complex number
var maxIterations = 200
var maxModulusSize = float64(6)

// Mutex lock for drawing threads
var drawLock sync.Mutex

// Making a channel for synchronization of the worker/computing threads
var coordinatesValueChannel chan image.Point
var wgComputation sync.WaitGroup

// Keep track of the computation time
var computationTime float64

// Starts the computation of mandelbrot
func compute() {
	// Initialize channel
	coordinatesValueChannel = make(chan image.Point)

	// Set the number of threads to wait for
	wgComputation.Add(maxComputingThreads)

	// Initialize image
	min.X = -int(width / 2)
	max.X = width + min.X
	min.Y = -int(height / 2)
	max.Y = height + min.Y

	upLeft := image.Point{0, 0}
	lowRight := image.Point{max.X - min.X, max.Y - min.Y}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Keep track of the computation time
	start := time.Now()

	// Creating the threads for drawing and computing
	for i := 0; i < maxComputingThreads; i++ {
		go computingThread(*img)
	}

	// maxWaitingTime := -1.0
	// put all the coordinates in the channel
	for i := min.X; i < max.X; i++ {
		for j := min.Y; j < max.Y; j++ {
			coordinates := image.Point{i, j}
			coordinatesValueChannel <- coordinates
		}
	}

	// Wait for threads and close all channels
	close(coordinatesValueChannel)
	wgComputation.Wait()

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	t := time.Now()
	computationTime = t.Sub(start).Seconds()
	fmt.Println("Total time: ", t.Sub(start))
}

// Compute whether the coordinates belong in the set or not.
// This calls the recursion
func computingThread(img image.RGBA) {
	for true {
		coordinates, stopSignal := <-coordinatesValueChannel
		if !stopSignal {
			wgComputation.Done()
			return
		}
		computation(0, coordinates, maxIterations, img)
	}
}

// The recursion for determining whether the value belongs to the set or not
func computation(z complex128, coordinates image.Point, iterations int, img image.RGBA) {
	modulus := cmplx.Abs(z)
	if iterations == 0 || modulus > maxModulusSize {
		draw(coordinates, iterations, img)
		return
	}
	x := float64(coordinates.X)/float64(-min.X)*scale + centerX
	y := float64(coordinates.Y)/float64(-min.Y)*scale + centerY

	c := complex(x, y)
	z = (z * z) + c
	computation(z, coordinates, iterations-1, img)
}

// Draws the mandelbrot from the value
func draw(coordinates image.Point, iterations int, img image.RGBA) {
	// draw the value given
	colorValue := uint8(iterations)
	pixelColor := color.RGBA{colorValue * uint8(redFactor), colorValue * uint8(blueFactor), colorValue * uint8(greenFactor), 0xff}
	x := coordinates.X - min.X
	y := coordinates.Y - min.Y

	drawLock.Lock()
	img.Set(x, y, pixelColor)
	drawLock.Unlock()
}
