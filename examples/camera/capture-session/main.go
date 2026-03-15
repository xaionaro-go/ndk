// Camera capture session lifecycle example.
//
// Demonstrates the complete Camera2 capture pipeline: creating a manager,
// listing cameras, reading characteristics, opening a device with callbacks,
// creating a capture request and session, starting a repeating request, and
// orderly cleanup. A sync.WaitGroup gates the repeating-request call until
// the session reports ready.
//
// Because the camera pipeline requires an ANativeWindow from an Android
// Surface, this example documents the window-dependent steps in comments
// and executes only the steps that work without a live Surface.
//
// Prerequisites:
//   - Android device with API level 24+.
//   - android.permission.CAMERA must be granted at runtime.
package main

import (
	"fmt"
	"log"

	"github.com/AndroidGoLab/ndk/camera"
)

func main() {
	// --- Stage 1: Create manager ---
	mgr := camera.NewManager()
	if mgr == nil {
		log.Fatal("camera.NewManager returned nil")
	}
	defer func() {
		if err := mgr.Close(); err != nil {
			log.Printf("manager close: %v", err)
		}
		log.Println("manager closed")
	}()
	log.Println("manager created")

	// --- Stage 2: List cameras ---
	ids, err := mgr.CameraIDList()
	if err != nil {
		log.Fatalf("listing cameras: %v", err)
	}
	if len(ids) == 0 {
		log.Fatal("no cameras available on this device")
	}
	fmt.Printf("found %d camera(s):\n", len(ids))
	for i, id := range ids {
		fmt.Printf("  [%d] %s\n", i, id)
	}

	// --- Stage 3: Read characteristics ---
	camID := ids[0]
	chars, err := mgr.GetCameraCharacteristics(camID)
	if err != nil {
		log.Fatalf("get characteristics for %s: %v", camID, err)
	}
	defer chars.Close()

	// Query sensor orientation.
	if chars.I32Count(uint32(camera.SensorOrientation)) > 0 {
		orient := chars.I32At(uint32(camera.SensorOrientation), 0)
		fmt.Printf("camera %s sensor orientation: %d degrees\n", camID, orient)
	}

	// Query available stream configurations (format, width, height, isInput).
	scTag := uint32(camera.ScalerAvailableStreamConfigurations)
	count := chars.I32Count(scTag)
	outputs := 0
	for i := int32(0); i+3 < count; i += 4 {
		isInput := chars.I32At(scTag, i+3)
		if isInput == 0 {
			outputs++
		}
	}
	fmt.Printf("camera %s has %d output stream configurations\n", camID, outputs)

	// --- Stage 4: Open camera device ---
	//
	// Opening the camera requires android.permission.CAMERA. The device
	// state callbacks receive disconnect and error notifications.
	dev, err := mgr.OpenCamera(camID, camera.DeviceStateCallbacks{
		OnDisconnected: func() {
			log.Println("callback: camera disconnected")
		},
		OnError: func(code int) {
			log.Printf("callback: camera error %d", code)
		},
	})
	if err != nil {
		log.Fatalf("opening camera %s: %v", camID, err)
	}
	defer func() {
		if err := dev.Close(); err != nil {
			log.Printf("device close: %v", err)
		}
		log.Println("device closed")
	}()
	log.Printf("camera %s opened", camID)

	// --- Stage 5: Create capture request ---
	//
	// The template type provides sensible defaults. Available templates:
	//   camera.Preview, camera.StillCapture, camera.Record,
	//   camera.VideoSnapshot, camera.ZeroShutterLag, camera.Manual
	req, err := dev.CreateCaptureRequest(camera.Preview)
	if err != nil {
		log.Fatalf("creating capture request: %v", err)
	}
	defer req.Close()
	log.Println("capture request created (Preview template)")

	// --- Stage 6: Output target, session output, session ---
	//
	// The remaining steps require an ANativeWindow, which comes from an
	// Android Surface (SurfaceView, TextureView, ImageReader, or
	// SurfaceTexture). The pattern is shown below. Without a window, we
	// print the pipeline steps and exit cleanly.
	//
	//   // Obtain an ANativeWindow from a Surface.
	//   var nativeWindow *camera.ANativeWindow
	//
	//   // Create output target and add to request.
	//   outTarget, err := camera.NewOutputTarget(nativeWindow)
	//   if err != nil {
	//       log.Fatalf("creating output target: %v", err)
	//   }
	//   defer outTarget.Close()
	//   req.AddTarget(outTarget)
	//
	//   // Create session output and container.
	//   sessOut, err := camera.NewSessionOutput(nativeWindow)
	//   if err != nil {
	//       log.Fatalf("creating session output: %v", err)
	//   }
	//   defer sessOut.Close()
	//
	//   container, err := camera.NewSessionOutputContainer()
	//   if err != nil {
	//       log.Fatalf("creating output container: %v", err)
	//   }
	//   defer container.Close()
	//
	//   if err := container.Add(sessOut); err != nil {
	//       log.Fatalf("adding session output: %v", err)
	//   }
	//
	//   // Create capture session with state callbacks.
	//   // Use a WaitGroup to block until the session is ready.
	//   var wg sync.WaitGroup
	//   wg.Add(1)
	//   var readyOnce sync.Once
	//
	//   session, err := dev.CreateCaptureSession(
	//       container,
	//       camera.SessionStateCallbacks{
	//           OnClosed: func() {
	//               log.Println("callback: session closed")
	//           },
	//           OnReady: func() {
	//               log.Println("callback: session ready")
	//               readyOnce.Do(wg.Done)
	//           },
	//           OnActive: func() {
	//               log.Println("callback: session active")
	//           },
	//       },
	//   )
	//   if err != nil {
	//       log.Fatalf("creating capture session: %v", err)
	//   }
	//   defer session.Close()
	//
	//   // Wait for the session to become ready.
	//   wg.Wait()
	//
	//   // Start continuous preview capture.
	//   if err := session.SetRepeatingRequest(req); err != nil {
	//       log.Fatalf("set repeating request: %v", err)
	//   }
	//   log.Println("repeating request started")
	//
	//   // ... render frames ...
	//
	//   // Stop the repeating request before cleanup.
	//   if err := session.StopRepeating(); err != nil {
	//       log.Printf("stop repeating: %v", err)
	//   }
	//   log.Println("repeating request stopped")

	// --- Cleanup order ---
	//
	// Resources are released in reverse creation order via defers:
	//   1. session.Close()
	//   2. container.Close()
	//   3. sessOut.Close()
	//   4. outTarget.Close()
	//   5. req.Close()
	//   6. dev.Close()
	//   7. mgr.Close()

	log.Println("capture session pipeline complete -- see source for full pattern")
}
