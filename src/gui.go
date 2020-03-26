/*
	Author: Steven Ge
	Date: 2020-03-20
*/

// TODO: smooth navigation
// - Find proper resolution for navigation
// - Handle long pressed keys/drags properly (no queueing the key presses)
// TODO: Sometimes displays in a weird format. Need fix
// TODO: Implement scroll to zoom

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
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
var dragWidth int
var dragHeight int
var released = true

// Color factors
var redFactorEntry = widget.NewEntry()
var greenFactorEntry = widget.NewEntry()
var blueFactorEntry = widget.NewEntry()

// Scale level
var scaleEntry = widget.NewEntry()
var zoomSpeedEntry = widget.NewEntry()

// Center coordinates
var centerXEntry = widget.NewEntry()
var centerYEntry = widget.NewEntry()

// Image
var img fyne.Resource

type dragableScrollableIcon struct {
	widget.Icon
}

func main() {
	// Load settings
	byteValue, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Define actions of controls
	win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case "Right":
			move("Right")
		case "Left":
			move("Left")
		case "Up":
			move("Up")
		case "Down":
			move("Down")
		case "A":
			move("Left")
		case "D":
			move("Right")
		case "W":
			move("Up")
		case "S":
			move("Down")

		case "I":
			move("I")
		case "O":
			move("O")
		}
	})

	// Set initial values and placeholder of resolution settings fields
	widthEntry.SetPlaceHolder("Width")
	heightEntry.SetPlaceHolder("Height")

	// Set initial values and placeholder of color factor settings fields
	redFactorEntry.SetPlaceHolder("red factor")
	greenFactorEntry.SetPlaceHolder("green factor")
	blueFactorEntry.SetPlaceHolder("blue factor")

	// Set initial values and placeholder of settings.Scale field
	scaleEntry.SetPlaceHolder("Scale")
	zoomSpeedEntry.SetPlaceHolder("Zoom speed")

	// Set initial values and placeholder of center coordinates fields
	centerXEntry.SetPlaceHolder("Real")
	centerYEntry.SetPlaceHolder("Imaginary")

	// Put the elements in the window
	setInterface("")

	// Start window
	win.Resize(fyne.NewSize(1200, 700))
	win.ShowAndRun()
}

// Actions for when the button is pressed
// These are setting the global variables and start the computation
func buttonAction(w string, h string, rf string, gf string, bf string, sc string, zs string, cx string, cy string) {
	// Convert values of input fields to numerical variable type
	w1, _ := strconv.ParseFloat(w, 64)
	h1, _ := strconv.ParseFloat(h, 64)
	rf1, _ := strconv.ParseFloat(rf, 64)
	gf1, _ := strconv.ParseFloat(gf, 64)
	bf1, _ := strconv.ParseFloat(bf, 64)
	sc1, _ := strconv.ParseFloat(sc, 64)
	zs1, _ := strconv.ParseFloat(zs, 64)
	cx1, _ := strconv.ParseFloat(cx, 64)
	cy1, _ := strconv.ParseFloat(cy, 64)

	// Convert values to correct variable type and set the global variables
	settings.Width = int(w1)
	settings.Height = int(h1)
	settings.RedFactor = int(rf1)
	settings.GreenFactor = int(gf1)
	settings.BlueFactor = int(bf1)
	settings.Scale = sc1
	settings.MovingSpeed = zs1
	settings.CenterX = cx1
	settings.CenterY = cy1

	// Print loading in window
	setInterface("Loading")
	// Recompute and refresh image
	compute()
	setInterface("Computation time:\n" + strconv.FormatFloat(computationTime, 'E', 5, 64) + "\nSeconds")
}

// Initialize the items in the window
func setInterface(loadingString string) {
	// Save settings
	settingsJson, _ := json.Marshal(settings)
	err := ioutil.WriteFile("settings.json", settingsJson, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Load image
	file, err := os.Open(imagePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	img = fyne.NewStaticResource("image", bytes)

	// Set the values of the button (Needed when using mouse/keyboard controls)
	// Set initial values and placeholder of resolution settings fields
	widthEntry.SetText(strconv.Itoa(settings.Width))
	heightEntry.SetText(strconv.Itoa(settings.Height))

	// Set initial values and placeholder of color factor settings fields
	redFactorEntry.SetText(strconv.Itoa(settings.RedFactor))
	greenFactorEntry.SetText(strconv.Itoa(settings.GreenFactor))
	blueFactorEntry.SetText(strconv.Itoa(settings.BlueFactor))

	// Set initial values and placeholder of settings.Scale field
	scaleEntry.SetText(strconv.FormatFloat(settings.Scale, 'f', -1, 64))
	zoomSpeedEntry.SetText(strconv.FormatFloat(settings.MovingSpeed, 'f', -1, 64))

	// Set initial values and placeholder of center coordinates fields
	centerXEntry.SetText(strconv.FormatFloat(settings.CenterX, 'f', -1, 64))
	centerYEntry.SetText(strconv.FormatFloat(settings.CenterY, 'f', -1, 64))

	// Set settings button
	setButton := widget.NewButton("Set", func() {
		buttonAction(
			widthEntry.Text,
			heightEntry.Text,
			redFactorEntry.Text,
			greenFactorEntry.Text,
			blueFactorEntry.Text,
			scaleEntry.Text,
			zoomSpeedEntry.Text,
			centerXEntry.Text,
			centerYEntry.Text)
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
		widget.NewLabel("Scale"),
		scaleEntry,
		widget.NewLabel("Move speed"),
		zoomSpeedEntry,
		widget.NewLabel("Center coordinates"),
		centerXEntry,
		centerYEntry,
		setButton,
		widget.NewLabel(loadingString))

	// Defining the layout of the window
	r := newdragableScrollableIcon(img)
	container := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, settingsList, nil),
		settingsList, r,
	)
	win.SetContent(container)
}

// Keyboard navigation
func move(movement string) {
	// Print loading in window
	setInterface("Loading")
	switch movement {
	case "Right":
		settings.CenterX = settings.CenterX + (settings.Scale * settings.MovingSpeed)
	case "Left":
		settings.CenterX = settings.CenterX - (settings.Scale * settings.MovingSpeed)
	case "Up":
		settings.CenterY = settings.CenterY - (settings.Scale * settings.MovingSpeed)
	case "Down":
		settings.CenterY = settings.CenterY + (settings.Scale * settings.MovingSpeed)

	case "I":
		settings.Scale = settings.Scale - (settings.Scale * settings.MovingSpeed)
	case "O":
		settings.Scale = settings.Scale + (settings.Scale * settings.MovingSpeed)
	}
	// Recompute and refresh image
	compute()
	setInterface("Computation time:\n" + strconv.FormatFloat(computationTime, 'E', 5, 64) + "\nSeconds")
}

// Mouse navigation functions
func drag(x int, y int) {
	// Print loading in window
	setInterface("Loading")
	settings.CenterX = settings.CenterX - (settings.Scale * float64(x) * settings.MovingSpeed)
	settings.CenterY = settings.CenterY - (settings.Scale * float64(y) * settings.MovingSpeed)
	// Recompute and refresh image
	compute()
	setInterface("Computation time:\n" + strconv.FormatFloat(computationTime, 'E', 5, 64) + "\nSeconds")
}

func scroll(s int) {
	// Print loading in window
	// setInterface("Loading")

	settings.Scale = settings.Scale - (settings.Scale * settings.MovingSpeed * float64(s))

	fmt.Println(s)
	// Recompute and refresh image
	compute()
	// setInterface("Computation time:\n" + strconv.FormatFloat(computationTime, 'E', 5, 64) + "\nSeconds")
}

func (t *dragableScrollableIcon) Dragged(d *fyne.DragEvent) {
	if released {
		dragWidth = settings.Width
		dragHeight = settings.Height
		released = false
	}

	settings.Width = 50
	settings.Height = 25

	drag(d.DraggedX, d.DraggedY)
}

func (t *dragableScrollableIcon) DragEnd() {
	// Print loading in window
	setInterface("Loading")
	settings.Width = dragWidth
	settings.Height = dragHeight
	// Recompute and refresh image
	compute()
	setInterface("Computation time:\n" + strconv.FormatFloat(computationTime, 'E', 5, 64) + "\nSeconds")
	released = true
}

func (t *dragableScrollableIcon) Scrolled(s *fyne.ScrollEvent) {
	if released {
		dragWidth = settings.Width
		dragHeight = settings.Height
		released = false
	}
	settings.Width = 50
	settings.Height = 25

	scroll(s.DeltaY)

	settings.Width = dragWidth
	settings.Height = dragHeight
	released = true
}

func newdragableScrollableIcon(res fyne.Resource) *dragableScrollableIcon {
	icon := &dragableScrollableIcon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)

	return icon
}
