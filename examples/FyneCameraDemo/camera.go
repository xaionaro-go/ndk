package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/camera"
	"github.com/xaionaro-go/ndk/media"
	"github.com/xaionaro-go/ndk/ndk"
)

type Camera struct {
	*camera.Device
	*media.ImageReader
}

var _ camera.CaptureSessionStateCallbacks = (*Camera)(nil)

func (cam *Camera) OnClosed(sess *camera.CaptureSession) {

}
func (cam *Camera) OnReady(sess *camera.CaptureSession) {
}
func (cam *Camera) OnActive(sess *camera.CaptureSession) {

}
func (cam *Camera) Close() error {
	return fmt.Errorf("not implemented, yet")
}

func (cam *Camera) StartStreaming(
	requestTemplate camera.DeviceRequestTemplate,
	targetWindows ...*ndk.Window,
) error {
	request, err := cam.Device.CreateCaptureRequest(requestTemplate)
	if err != nil {
		return fmt.Errorf("CreateCaptureRequest: %w", err)
	}

	outputContainer, err := camera.CaptureSessionOutputContainerCreate()
	if err != nil {
		return fmt.Errorf("CaptureSessionOutputContainerCreate: %w", err)
	}

	for idx, win := range targetWindows {
		sessionOutput, err := camera.CaptureSessionOutputCreate(win)
		if err != nil {
			return fmt.Errorf("%d: CaptureSessionOutputCreate: %w", idx, err)
		}
		err = outputContainer.Add(sessionOutput)
		if err != nil {
			return fmt.Errorf("%d, outputContainer.Add: %w", idx, err)
		}

		target, err := camera.CameraOutputTargetCreate(win)
		if err != nil {
			return fmt.Errorf("%d: CameraOutputTargetCreate: %w", idx, err)
		}

		err = request.AddTarget(target)
		if err != nil {
			return fmt.Errorf("%d: AddTarget: %w", idx, err)
		}
	}

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
