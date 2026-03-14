package main

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/camera"
	mediacapi "github.com/xaionaro-go/ndk/capi/media"
	"github.com/xaionaro-go/ndk/media"
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

var cameraCaptureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture raw frames from a camera using ImageReader",
	Long: `Captures N raw image frames from the specified camera device.
Frames are written sequentially to the output file as raw pixel data.
Uses the camera2 NDK pipeline: ImageReader -> OutputTarget -> CaptureSession.`,
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		cameraID, _ := cmd.Flags().GetString("id")
		width, _ := cmd.Flags().GetInt32("width")
		height, _ := cmd.Flags().GetInt32("height")
		format, _ := cmd.Flags().GetInt32("format")
		count, _ := cmd.Flags().GetInt32("count")
		output, _ := cmd.Flags().GetString("output")

		// Create ImageReader via capi (the idiomatic NewImageReader has a broken
		// signature that casts **ImageReader to **capi.AImageReader incorrectly).
		// maxImages is the max number of images the reader can hold simultaneously.
		// Keep it small (4) — we acquire and release frames in a loop.
		var maxImages int32 = 4
		var readerPtr *mediacapi.AImageReader
		status := mediacapi.AImageReader_new(width, height, format, maxImages, &readerPtr)
		if status < 0 {
			return fmt.Errorf("creating image reader: media error %d", status)
		}
		reader := media.NewImageReaderFromPointer(unsafe.Pointer(readerPtr))
		defer reader.Close()

		// Get ANativeWindow from the ImageReader.
		var nativeWindow *mediacapi.ANativeWindow
		status = mediacapi.AImageReader_getWindow(readerPtr, &nativeWindow)
		if status < 0 {
			return fmt.Errorf("getting window from image reader: media error %d", status)
		}

		// The camera package has its own ANativeWindow type alias; cast through
		// unsafe.Pointer since both are C.ANativeWindow underneath.
		camWindow := (*camera.ANativeWindow)(unsafe.Pointer(nativeWindow))

		mgr := camera.NewManager()
		defer mgr.Close()

		// Create OutputTarget for the capture request.
		outputTarget, err := camera.NewOutputTarget(camWindow)
		if err != nil {
			return fmt.Errorf("creating output target: %w", err)
		}
		defer outputTarget.Close()

		// Create SessionOutput + container for the capture session.
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

		// Open the camera device with no-op state callbacks.
		device, err := mgr.OpenCamera(cameraID, camera.DeviceStateCallbacks{
			OnDisconnected: func() {},
			OnError:        func(int) {},
		})
		if err != nil {
			return fmt.Errorf("opening camera %q: %w", cameraID, err)
		}
		defer device.Close()

		// Create a capture request using the StillCapture template.
		request, err := device.CreateCaptureRequest(camera.StillCapture)
		if err != nil {
			return fmt.Errorf("creating capture request: %w", err)
		}
		defer request.Close()

		request.AddTarget(outputTarget)

		// Create a capture session with no-op state callbacks.
		session, err := device.CreateCaptureSession(container, camera.SessionStateCallbacks{
			OnClosed: func() {},
			OnReady:  func() {},
			OnActive: func() {},
		})
		if err != nil {
			return fmt.Errorf("creating capture session: %w", err)
		}
		defer session.Close()

		// Start repeating request so the camera produces frames.
		if err := session.SetRepeatingRequest(request); err != nil {
			return fmt.Errorf("setting repeating request: %w", err)
		}
		defer func() {
			if stopErr := session.StopRepeating(); stopErr != nil && _err == nil {
				_err = fmt.Errorf("stopping repeating request: %w", stopErr)
			}
		}()

		// Open output file.
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()

		// Wait for the camera pipeline to start producing frames.
		time.Sleep(500 * time.Millisecond)

		// Acquire and write frames.
		var totalBytes int64
		for i := int32(0); i < count; i++ {
			var imagePtr *mediacapi.AImage

			// Retry acquiring — the camera may not have a frame ready yet.
			for retries := 0; retries < 100; retries++ {
				status = mediacapi.AImageReader_acquireNextImage(readerPtr, &imagePtr)
				if status >= 0 {
					break
				}
				time.Sleep(33 * time.Millisecond) // ~30fps polling
			}
			if status < 0 {
				return fmt.Errorf("acquiring image %d: media error %d", i, status)
			}

			// Read plane 0 data.
			var dataPtr *uint8
			var dataLen int32
			status = mediacapi.AImage_getPlaneData(imagePtr, 0, &dataPtr, &dataLen)
			if status < 0 {
				mediacapi.AImage_delete(imagePtr)
				return fmt.Errorf("getting plane data for image %d: media error %d", i, status)
			}

			// Copy the pixel data from the C buffer to a Go slice and write it.
			data := unsafe.Slice(dataPtr, dataLen)
			if _, err := f.Write(data); err != nil {
				mediacapi.AImage_delete(imagePtr)
				return fmt.Errorf("writing image %d to file: %w", i, err)
			}
			totalBytes += int64(dataLen)

			mediacapi.AImage_delete(imagePtr)
			fmt.Printf("captured frame %d/%d (%d bytes)\n", i+1, count, dataLen)
		}

		fmt.Printf("captured %d frames (%d bytes total) to %s\n", count, totalBytes, output)
		return nil
	},
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
	cameraCaptureCmd.Flags().Int32("format", mediacapi.AIMAGE_FORMAT_YUV_420_888, "image format (35=YUV_420_888, 256=JPEG)")
	cameraCaptureCmd.Flags().Int32("count", 1, "number of frames to capture")
	cameraCaptureCmd.Flags().String("output", "capture.raw", "output file path")

	cameraCmd.AddCommand(cameraCaptureCmd)
	cameraCmd.AddCommand(cameraListDetailsCmd)
}
