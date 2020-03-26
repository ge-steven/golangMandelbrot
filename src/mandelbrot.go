/*
	Author: Steven Ge
	Date: 2020-03-20
*/

package main

import (
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
type Settings struct {
	Width       int     `json:"Width"`
	Height      int     `json:"Height"`
	Scale       float64 `json:"Scale"`
	CenterX     float64 `json:"CenterX"`
	CenterY     float64 `json:"CenterY"`
	RedFactor   int     `json:"RedFactor"`
	GreenFactor int     `json:"GreenFactor"`
	BlueFactor  int     `json:"BlueFactor"`
	MovingSpeed float64 `json:"MovingSpeed"`
}

var settings Settings

var upLeft image.Point
var lowRight image.Point
var min image.Point
var max image.Point

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
	min.X = -int(settings.Width / 2)
	max.X = settings.Width + min.X
	min.Y = -int(settings.Height / 2)
	max.Y = settings.Height + min.Y

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
	x := float64(coordinates.X)/float64(-min.X)*settings.Scale + settings.CenterX
	y := float64(coordinates.Y)/float64(-min.Y)*settings.Scale + settings.CenterY

	c := complex(x, y)
	z = (z * z) + c
	computation(z, coordinates, iterations-1, img)
}

// Draws the mandelbrot from the value
func draw(coordinates image.Point, iterations int, img image.RGBA) {
	// draw the value given
	colorValue := uint8(iterations)
	pixelColor := color.RGBA{colorValue * uint8(settings.RedFactor), colorValue * uint8(settings.BlueFactor), colorValue * uint8(settings.GreenFactor), 0xff}
	x := coordinates.X - min.X
	y := coordinates.Y - min.Y

	drawLock.Lock()
	img.Set(x, y, pixelColor)
	drawLock.Unlock()
}
