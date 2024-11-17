//go:build android
// +build android

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/pflag"
	"github.com/xaionaro-go/ndk/camera"
	"github.com/xaionaro-go/ndk/media"
)

import "C"

func main() {
	netPprofAddr := pflag.String("net-pprof-addr", "", "")
	pflag.Parse()

	a := app.New()
	w := a.NewWindow("Display camera")

	if *netPprofAddr == "" && runtime.GOOS == "android" {
		*netPprofAddr = "0.0.0.0:0"
	}
	if *netPprofAddr != "" {
		go func() {
			defer func() { processRecover(w, recover()) }()
			log.Println(http.ListenAndServe(*netPprofAddr, nil))
		}()
	}

	defer func() { processRecover(w, recover()) }()

	camID := "0"

	camMan := CameraManagerInstance()
	cam, err := camMan.OpenCamera(camID)
	if err != nil {
		panicInUI(w, fmt.Errorf("OpenCamera: %w", err))
	}
	defer cam.Close()

	imageListener, err := NewImageListener(
		1920, 1080,
		media.FORMAT_YUV_420_888,
	)
	if err != nil {
		panicInUI(w, fmt.Errorf("NewImageListener: %w", err))
	}

	err = cam.StartStreaming(
		camera.TEMPLATE_ZERO_SHUTTER_LAG,
		imageListener.Window,
	)
	if err != nil {
		panicInUI(w, fmt.Errorf("StartStreaming: %w", err))
	}

	mImg := <-imageListener.ImageChan
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
	fyneImg.FillMode = canvas.ImageFillContain

	c := w.Canvas()
	c.SetContent(fyneImg)

	go func() {
		defer func() { processRecover(w, recover()) }()
		for newMImg := range imageListener.ImageChan {
			_, err := ConvertImage(newMImg)
			if err != nil {
				newMImg.Delete()
				panicInUI(w, fmt.Errorf("ConvertImage: %w", err))
			}
			fyneImg.Image = img
			fyneImg.Refresh()
			mImg.Delete()
			mImg = newMImg
		}
	}()

	w.ShowAndRun()
	<-context.Background().Done()
}

func panicInUI(
	w fyne.Window,
	err error,
) {
	log.Println(err.Error())
	text := widget.NewLabel(err.Error())
	text.Wrapping = fyne.TextWrapWord
	w.SetContent(text)
	w.ShowAndRun()
	<-context.Background().Done()
}

func processRecover(
	w fyne.Window,
	r any,
) {
	if r == nil {
		return
	}

	debug.PrintStack()
	panicInUI(w, fmt.Errorf("%v\n\n%s", r, debug.Stack()))
}
