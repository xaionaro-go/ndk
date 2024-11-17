package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/media"
	"github.com/xaionaro-go/ndk/ndk"
)

type ImageListener struct {
	Window      *ndk.Window
	ImageReader *media.ImageReader
	ImageChan   chan *media.Image
}

func NewImageListener(
	width uint,
	height uint,
	format media.Formats,
) (*ImageListener, error) {
	l := &ImageListener{
		ImageChan: make(chan *media.Image, 1),
	}

	var err error
	l.ImageReader, err = media.NewImageReader(int(width), int(height), format, 4)
	if err != nil {
		return nil, fmt.Errorf("NewImageReader: %w", err)
	}

	err = l.ImageReader.SetImageListener(l.OnImage)
	if err != nil {
		return nil, fmt.Errorf("SetImageListener: %w", err)
	}

	l.Window, err = l.ImageReader.GetWindow()
	if err != nil {
		return nil, fmt.Errorf("ImageReader.GetWindow: %w", err)
	}

	l.Window.Acquire()
	return l, nil
}

func (l *ImageListener) OnImage(r *media.ImageReader) {
	img, err := r.AcquireNextImage()
	if err != nil {
		panic(fmt.Errorf("AcquireNextImage: %w", err))
	}

	select {
	case l.ImageChan <- img:
	default:
		img.Delete()
	}
}
