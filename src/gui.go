/*
	Author: Steven Ge
	Date: 2020-03-20
*/

package main

import (
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// Window
var application = app.New()
var win = application.NewWindow("Mandelbrot")

// Path to mandelbrot image
var imagePath = "image.png"

// Initialize fields of the settings
// Resolution
var widthEntry = widget.NewEntry()
var heightEntry = widget.NewEntry()

// Color factors
var redFactorEntry = widget.NewEntry()
var greenFactorEntry = widget.NewEntry()
var blueFactorEntry = widget.NewEntry()

// Zoom level
var zoom = 2.0

// Load image
var img = canvas.NewImageFromFile(imagePath)

func main() {

	// Set initial values and placeholder of resolution settings fields
	widthEntry.SetText(strconv.Itoa(width))
	heightEntry.SetText(strconv.Itoa(height))
	widthEntry.SetPlaceHolder("width")
	heightEntry.SetPlaceHolder("Height")

	// Set initial values and placeholder of color factor settings fields
	redFactorEntry.SetText(strconv.Itoa(redFactor))
	greenFactorEntry.SetText(strconv.Itoa(greenFactor))
	blueFactorEntry.SetText(strconv.Itoa(blueFactor))
	redFactorEntry.SetPlaceHolder("red factor")
	greenFactorEntry.SetPlaceHolder("green factor")
	blueFactorEntry.SetPlaceHolder("blue factor")

	// Put the elements in the window
	initializeInterface("")

	// Start window
	win.Resize(fyne.NewSize(640, 480))
	win.ShowAndRun()
}

// Actions for when the button is pressed
// These are setting the global variables and start the computation
func buttonAction(w string, h string, s float64, cx float64, cy float64, rf string, gf string, bf string) {
	// Print loading in window
	initializeInterface("Loading")

	// Convert values of input fields to numerical variable type
	w1, _ := strconv.ParseFloat(w, 64)
	h1, _ := strconv.ParseFloat(h, 64)
	rf1, _ := strconv.ParseFloat(rf, 64)
	gf1, _ := strconv.ParseFloat(gf, 64)
	bf1, _ := strconv.ParseFloat(bf, 64)

	// Convert values to correct variable type and set the global variables
	width = int(w1)
	height = int(h1)
	redFactor = int(rf1)
	greenFactor = int(gf1)
	blueFactor = int(bf1)

	// Recompute and refresh image
	compute()
	initializeInterface("Computation time:\n" + strconv.FormatFloat(computationTime, 'E', 5, 64) + "\nSeconds")
}

// Initialize the items in the window
func initializeInterface(loadingString string) {
	// Set settings button
	setButton := widget.NewButton("Set", func() {
		buttonAction(
			widthEntry.Text,
			heightEntry.Text,
			zoom,
			-0.5,
			0,
			redFactorEntry.Text,
			greenFactorEntry.Text,
			blueFactorEntry.Text)
	})

	// Grouping the input fields and labels of settings
	settingsList := widget.NewGroup("          Settings          ",
		widget.NewLabel("Resolution"),
		widthEntry,
		heightEntry,
		widget.NewLabel("Color Factor"),
		redFactorEntry,
		greenFactorEntry,
		blueFactorEntry,
		setButton,
		widget.NewLabel(loadingString))

	// Defining the layout of the window
	container := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, settingsList, nil),
		settingsList, img,
	)

	win.SetContent(container)
}
