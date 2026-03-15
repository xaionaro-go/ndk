package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/camera"
	"github.com/AndroidGoLab/ndk/looper"
	"github.com/AndroidGoLab/ndk/media"
	"github.com/AndroidGoLab/ndk/window"
)

// cameraTagLensFacing is ACAMERA_LENS_FACING from the NDK header.
const cameraTagLensFacing = 524293

// cameraTagInfoSupportedHardwareLevel is ACAMERA_INFO_SUPPORTED_HARDWARE_LEVEL.
const cameraTagInfoSupportedHardwareLevel = 1376256

func lensFacingString(value int32) string {
	switch value {
	case 0:
		return "front"
	case 1:
		return "back"
	case 2:
		return "external"
	default:
		return fmt.Sprintf("unknown(%d)", value)
	}
}

func hardwareLevelString(value int32) string {
	switch value {
	case 0:
		return "limited"
	case 1:
		return "full"
	case 2:
		return "legacy"
	case 3:
		return "level_3"
	case 4:
		return "external"
	default:
		return fmt.Sprintf("unknown(%d)", value)
	}
}

// imageFormatYUV420888 is AIMAGE_FORMAT_YUV_420_888 from the NDK header.
const imageFormatYUV420888 = 35

var cameraCaptureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture raw frames from a camera using ImageReader",
	Long: `Captures N raw image frames from the specified camera device.
Frames are written sequentially to the output file as raw pixel data.
Uses the camera2 NDK pipeline: ImageReader -> OutputTarget -> CaptureSession.
Runs on a looper thread so Camera2 callbacks are dispatched correctly.`,
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		cameraID, _ := cmd.Flags().GetString("id")
		width, _ := cmd.Flags().GetInt32("width")
		height, _ := cmd.Flags().GetInt32("height")
		format, _ := cmd.Flags().GetInt32("format")
		count, _ := cmd.Flags().GetInt32("count")
		output, _ := cmd.Flags().GetString("output")

		// Run the entire capture pipeline on a looper thread.
		// Camera2 NDK dispatches session callbacks (OnReady, OnActive)
		// on the calling thread's looper — without one, callbacks never
		// fire and the session never produces frames.
		looper.Run(func(lp *looper.Looper) {
			_err = runCameraCapture(lp, cameraID, width, height, format, count, output)
		})
		return _err
	},
}

func runCameraCapture(
	lp *looper.Looper,
	cameraID string,
	width int32,
	height int32,
	format int32,
	count int32,
	output string,
) (_err error) {
	// Create ImageReader with CPU_READ_OFTEN usage for frame access.
	const usageCPUReadOften = uint64(0x3) // AHARDWAREBUFFER_USAGE_CPU_READ_OFTEN
	var maxImages int32 = 4
	reader, err := media.NewImageReaderWithUsage(width, height, format, usageCPUReadOften, maxImages)
	if err != nil {
		return fmt.Errorf("creating image reader: %w", err)
	}
	defer reader.Close()

	// Get ANativeWindow from the ImageReader.
	win, err := reader.Window()
	if err != nil {
		return fmt.Errorf("getting window from image reader: %w", err)
	}

	// Acquire a reference to the window (required before creating session outputs).
	// The media.Window and window.Window are separate wrapper types for the same
	// underlying ANativeWindow. Convert via unsafe.Pointer for the camera API.
	nw := window.NewWindowFromPointer(win.Pointer())
	nw.Acquire()
	camWindow := (*camera.ANativeWindow)(win.Pointer())

	mgr := camera.NewManager()
	defer mgr.Close()

	outputTarget, err := camera.NewOutputTarget(camWindow)
	if err != nil {
		return fmt.Errorf("creating output target: %w", err)
	}
	defer outputTarget.Close()

	sessionOutput, err := camera.NewSessionOutput(camWindow)
	if err != nil {
		return fmt.Errorf("creating session output: %w", err)
	}
	defer sessionOutput.Close()

	container, err := camera.NewSessionOutputContainer()
	if err != nil {
		return fmt.Errorf("creating session output container: %w", err)
	}
	defer container.Close()

	if err := container.Add(sessionOutput); err != nil {
		return fmt.Errorf("adding session output to container: %w", err)
	}

	fmt.Println("opening camera...")
	device, err := mgr.OpenCamera(cameraID, camera.DeviceStateCallbacks{
		OnDisconnected: func() { fmt.Println("callback: camera disconnected") },
		OnError:        func(code int) { fmt.Printf("callback: camera error %d\n", code) },
	})
	if err != nil {
		return fmt.Errorf("opening camera %q: %w", cameraID, err)
	}
	defer device.Close()

	request, err := device.CreateCaptureRequest(camera.Record)
	if err != nil {
		return fmt.Errorf("creating capture request: %w", err)
	}
	defer request.Close()

	request.AddTarget(outputTarget)

	// Use a WaitGroup to block until the session reports ready.
	var wg sync.WaitGroup
	wg.Add(1)
	var readyOnce sync.Once

	fmt.Println("creating capture session...")
	session, err := device.CreateCaptureSession(container, camera.SessionStateCallbacks{
		OnClosed: func() { fmt.Println("callback: session closed") },
		OnReady: func() {
			fmt.Println("callback: session ready")
			readyOnce.Do(wg.Done)
		},
		OnActive: func() { fmt.Println("callback: session active") },
	})
	if err != nil {
		return fmt.Errorf("creating capture session: %w", err)
	}
	defer session.Close()

	// Wait for the session to become ready. Camera2 dispatches callbacks
	// on its own internal thread — we just need to wait for OnReady.
	fmt.Println("waiting for session ready...")
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Session is ready.
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for camera session to become ready")
	}

	fmt.Println("camera session ready")

	if err := session.SetRepeatingRequest(request); err != nil {
		return fmt.Errorf("setting repeating request: %w", err)
	}
	defer func() {
		if stopErr := session.StopRepeating(); stopErr != nil && _err == nil {
			_err = fmt.Errorf("stopping repeating request: %w", stopErr)
		}
	}()

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	// Acquire and write frames. Poll the looper between acquires
	// to keep Camera2 callbacks flowing.
	var totalBytes int64
	for i := int32(0); i < count; i++ {
		var img *media.Image
		var acquireErr error

		for retries := 0; retries < 200; retries++ {
			img, acquireErr = reader.AcquireNextImage()
			if acquireErr == nil {
				break
			}
			// Poll looper to let Camera2 process internally.
			looper.PollOnce(16*time.Millisecond, nil, nil, nil)
		}
		if acquireErr != nil {
			return fmt.Errorf("acquiring image %d: %w", i, acquireErr)
		}

		dataPtr, dataLen, err := img.PlaneData(0)
		if err != nil {
			img.Close()
			return fmt.Errorf("getting plane data for image %d: %w", i, err)
		}

		data := unsafe.Slice(dataPtr, dataLen)
		if _, err := f.Write(data); err != nil {
			img.Close()
			return fmt.Errorf("writing image %d to file: %w", i, err)
		}
		totalBytes += int64(dataLen)

		img.Close()
		fmt.Printf("captured frame %d/%d (%d bytes)\n", i+1, count, dataLen)
	}

	fmt.Printf("captured %d frames (%d bytes total) to %s\n", count, totalBytes, output)
	return nil
}

var cameraListDetailsCmd = &cobra.Command{
	Use:   "list-details",
	Short: "List all cameras with full characteristics",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		mgr := camera.NewManager()
		defer mgr.Close()

		ids, err := mgr.CameraIDList()
		if err != nil {
			return fmt.Errorf("getting camera ID list: %w", err)
		}

		if len(ids) == 0 {
			fmt.Println("no cameras found")
			return nil
		}

		for _, id := range ids {
			meta, err := mgr.GetCameraCharacteristics(id)
			if err != nil {
				fmt.Printf("camera %s: error getting characteristics: %v\n", id, err)
				continue
			}

			facing := meta.I32At(uint32(cameraTagLensFacing), 0)
			orientation := meta.I32At(uint32(camera.SensorOrientation), 0)
			hwLevel := meta.I32At(uint32(cameraTagInfoSupportedHardwareLevel), 0)

			fmt.Printf("camera %s:\n", id)
			fmt.Printf("  lens facing:     %s\n", lensFacingString(facing))
			fmt.Printf("  orientation:     %d degrees\n", orientation)
			fmt.Printf("  hardware level:  %s\n", hardwareLevelString(hwLevel))

			meta.Close()
		}

		return nil
	},
}

func init() {
	cameraCaptureCmd.Flags().String("id", "0", "camera device ID")
	cameraCaptureCmd.Flags().Int32("width", 640, "image width in pixels")
	cameraCaptureCmd.Flags().Int32("height", 480, "image height in pixels")
	cameraCaptureCmd.Flags().Int32("format", imageFormatYUV420888, "image format (35=YUV_420_888, 256=JPEG)")
	cameraCaptureCmd.Flags().Int32("count", 1, "number of frames to capture")
	cameraCaptureCmd.Flags().String("output", "capture.raw", "output file path")

	cameraCmd.AddCommand(cameraCaptureCmd)
	cameraCmd.AddCommand(cameraListDetailsCmd)
}
