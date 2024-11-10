package main

import (
	"fmt"

	camera "github.com/xaionaro-go/ndk/camera24"
	media "github.com/xaionaro-go/ndk/media24"
)

type Camera struct {
	*camera.Device
	*media.ImageReader
	ImageChan chan *media.Image
}

var _ camera.CaptureSessionStateCallbacks = (*Camera)(nil)

func (cam *Camera) OnClosed(sess *camera.CaptureSession) {

}
func (cam *Camera) OnReady(sess *camera.CaptureSession) {
}
func (cam *Camera) OnActive(sess *camera.CaptureSession) {

}
func (cam *Camera) OnImage(r *media.ImageReader) {
	img, err := r.AcquireNextImage()
	if err != nil {
		panic(fmt.Errorf("AcquireNextImage: %w", err))
	}

	select {
	case cam.ImageChan <- img:
	default:
		img.Delete()
	}
}

func (cam *Camera) Close() error {
	return nil
}

func (cam *Camera) StartStreaming(
	width uint,
	height uint,
	format media.Formats,
) error {
	var err error
	cam.ImageReader, err = media.NewImageReader(int(width), int(height), format, 4)
	if err != nil {
		return fmt.Errorf("NewImageReader: %w", err)
	}

	err = cam.ImageReader.SetImageListener(cam.OnImage)
	if err != nil {
		return fmt.Errorf("SetImageListener: %w", err)
	}

	win, err := cam.ImageReader.GetWindow()
	if err != nil {
		return fmt.Errorf("ImageReader.GetWindow: %w", err)
	}

	win.Acquire()
	sessionOutput, err := camera.CaptureSessionOutputCreate(win)
	if err != nil {
		return fmt.Errorf("CaptureSessionOutputCreate: %w", err)
	}

	target, err := camera.CameraOutputTargetCreate(win)
	if err != nil {
		return fmt.Errorf("CameraOutputTargetCreate: %w", err)
	}

	request, err := cam.Device.CreateCaptureRequest(camera.TEMPLATE_ZERO_SHUTTER_LAG)
	if err != nil {
		return fmt.Errorf("CreateCaptureRequest: %w", err)
	}

	err = request.AddTarget(target)
	if err != nil {
		return fmt.Errorf("AddTarget: %w", err)
	}

	outputContainer, err := camera.CaptureSessionOutputContainerCreate()
	if err != nil {
		return fmt.Errorf("CaptureSessionOutputContainerCreate: %w", err)
	}
	outputContainer.Add(sessionOutput)

	captureSession, err := cam.Device.CreateCaptureSession(outputContainer, cam)
	if err != nil {
		return fmt.Errorf("CreateCaptureSession: %w", err)
	}

	err = captureSession.SetRepeatingRequest([]*camera.CaptureRequest{request})
	if err != nil {
		return fmt.Errorf("SetRepeatingRequest: %w", err)
	}

	return nil
}
