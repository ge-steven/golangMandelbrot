/*
	Author: Steven Ge
	Date: 2020-03-20
*/

package main

import (
	"fmt"
	"runtime"
	"strconv"

	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()
	var img canvas.Image
	img.File := "../mandelbrotjulianomultithread/HopefullyMandelbrot.png"
	fmt.Println(img.File)

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewLabel(strconv.Itoa(runtime.NumCPU())),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}
