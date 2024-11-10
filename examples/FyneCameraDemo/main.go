//go:build android
// +build android

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/pflag"
	media "github.com/xaionaro-go/ndk/media24"
)

import "C"

func main() {
	netPprofAddr := pflag.String("net-pprof-addr", "", "")
	pflag.Parse()

	a := app.New()
	w := a.NewWindow("Display camera")

	if *netPprofAddr != "" {
		go func() {
			defer func() { displayPanic(w, recover()) }()
			log.Println(http.ListenAndServe(*netPprofAddr, nil))
		}()
	}

	defer func() { displayPanic(w, recover()) }()

	camID := "0"

	camMan := CameraManagerInstance()
	cam, err := camMan.OpenCamera(camID)
	if err != nil {
		panicInUI(w, fmt.Errorf("OpenCamera: %w", err))
	}
	defer cam.Close()

	err = cam.StartStreaming(1920, 1080, media.FORMAT_YUV_420_888)
	if err != nil {
		panicInUI(w, fmt.Errorf("StartStreaming: %w", err))
	}

	mImg := <-cam.ImageChan
	if mImg == nil {
		panicInUI(w, fmt.Errorf("mImg == nil"))
	}
	width, err := mImg.GetWidth()
	if err != nil {
		panicInUI(w, fmt.Errorf("GetWidth: %w", err))
	}
	if width == 0 {
		panicInUI(w, fmt.Errorf("width == 0"))
	}

	height, err := mImg.GetHeight()
	if err != nil {
		panicInUI(w, fmt.Errorf("GetHeight: %w", err))
	}
	if height == 0 {
		panicInUI(w, fmt.Errorf("height == 0"))
	}

	img, err := ConvertImage(mImg)
	if err != nil {
		panicInUI(w, fmt.Errorf("ConvertImage: %w", err))
	}

	log.Println("width: ", width, "; height: ", height)

	// checking we are not panicking while reading the image
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.At(x, y)
		}
	}

	fyneImg := canvas.NewImageFromImage(img)
	fyneImg.ScaleMode = canvas.ImageScaleFastest
	fyneImg.FillMode = canvas.ImageFillOriginal

	c := w.Canvas()
	c.SetContent(fyneImg)

	if false {
		go func() {
			defer func() { displayPanic(w, recover()) }()
			for mImg := range cam.ImageChan {
				img, err := ConvertImage(mImg)
				if err != nil {
					panicInUI(w, fmt.Errorf("ConvertImage: %w", err))
				}
				fyneImg.Image = img
				fyneImg.Refresh()
			}
		}()
	}

	w.ShowAndRun()
	<-context.Background().Done()
}

func panicInUI(
	w fyne.Window,
	err error,
) {
	log.Panicln(err.Error())
	text := widget.NewLabel(err.Error())
	text.Wrapping = fyne.TextWrapWord
	w.SetContent(text)
	w.ShowAndRun()
	<-context.Background().Done()
}

func displayPanic(
	w fyne.Window,
	r any,
) {
	if r == nil {
		return
	}

	debug.PrintStack()
	panicInUI(w, fmt.Errorf("%v\n\n%s", r, debug.Stack()))
}
