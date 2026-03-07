// Camera2 pipeline overview.
//
// Documents the complete Android Camera2 capture pipeline using the ndk
// camera package types. Each stage of the pipeline is shown in sequence with
// the types that participate. Because a real capture requires an
// ANativeWindow (obtained from a Surface), runtime camera permissions, and
// a running Android activity, this example cannot execute end-to-end on its
// own. It serves as a reference for the object graph and cleanup order.
//
// Pipeline summary:
//
//	Manager
//	  -> open Device (via manager.OpenCamera)
//	       -> create CaptureRequest (via device.CreateCaptureRequest)
//	       -> create SessionOutputContainer
//	            -> add SessionOutput (wrapping an ANativeWindow)
//	       -> create CaptureSession (via device.CreateCaptureSession)
//	            -> set repeating request or single capture
//	            -> stop repeating
//	       <- close CaptureSession
//	  <- close Device
//	<- close Manager
//
// Prerequisites:
//   - Android device with API level 24+.
//   - android.permission.CAMERA must be granted at runtime before opening
//     any camera device. Without it the NDK returns camera.ErrPermissionDenied.
//   - A valid ANativeWindow handle is needed to create OutputTarget and
//     SessionOutput. Typically this comes from an Android Surface provided
//     by SurfaceView, TextureView, or an ImageReader.
package main

import (
	"log"

	"github.com/xaionaro-go/ndk/camera"
)

func main() {
	// ---------------------------------------------------------------
	// Stage 1: Manager
	//
	// The Manager is the root of the camera subsystem. It provides
	// camera enumeration and device-open operations.
	// ---------------------------------------------------------------
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

	// In a real application, you would now call manager methods to
	// enumerate camera IDs, then open a camera device with a
	// device-state callback.
	//
	// The ACameraDevice_StateCallbacks struct receives disconnect and
	// error notifications from the camera service.

	// ---------------------------------------------------------------
	// Stage 2: Device
	//
	// A Device represents a single opened camera. It is obtained by
	// calling the manager's open-camera method. Devices create capture
	// requests and capture sessions.
	// ---------------------------------------------------------------
	//
	// var device *camera.Device  // obtained from manager.OpenCamera
	//
	// defer func() {
	//     if err := device.Close(); err != nil {
	//         log.Printf("device close: %v", err)
	//     }
	// }()

	// ---------------------------------------------------------------
	// Stage 3: CaptureRequest
	//
	// A CaptureRequest describes a single frame capture configuration.
	// It is created from a Device using a TemplateType that supplies
	// reasonable defaults:
	//
	//   camera.Preview         - viewfinder / preview
	//   camera.StillCapture    - high-quality still photo
	//   camera.Record          - video recording
	//   camera.VideoSnapshot   - still image during video
	//   camera.ZeroShutterLag  - near-zero-lag capture
	//   camera.Manual          - fully manual control
	//
	// Created via:
	//   request, err := device.CreateCaptureRequest(camera.Preview)
	// ---------------------------------------------------------------
	//
	// var request *camera.CaptureRequest
	// defer request.Close()

	// ---------------------------------------------------------------
	// Stage 4: OutputTarget
	//
	// An OutputTarget wraps an ANativeWindow and is added to a
	// CaptureRequest to specify where captured frames are sent.
	// Created via:
	//   target, err := camera.NewOutputTarget(nativeWindow)
	// ---------------------------------------------------------------
	//
	// var target *camera.OutputTarget
	// request.AddTarget(target)
	// defer target.Close()

	// ---------------------------------------------------------------
	// Stage 5: SessionOutput + SessionOutputContainer
	//
	// A SessionOutput also wraps an ANativeWindow and is registered in
	// a SessionOutputContainer. The container is passed when creating a
	// CaptureSession so the camera HAL knows which outputs to prepare.
	// ---------------------------------------------------------------
	//
	// var sessionOutput *camera.SessionOutput
	//   // created via camera.NewSessionOutput(nativeWindow)
	//
	// var container *camera.SessionOutputContainer
	//   // created via camera.NewSessionOutputContainer()
	//
	// container.Add(sessionOutput)
	//
	// defer container.Close()
	// defer sessionOutput.Close()

	// ---------------------------------------------------------------
	// Stage 6: CaptureSession
	//
	// The CaptureSession is created from a Device with the output
	// container and a state-callback struct. It manages the stream
	// pipeline between the camera and the output surfaces.
	//
	// The capture session state callbacks receive ready, active, and
	// closed notifications.
	//
	// Once created, submit capture requests:
	//   - setRepeatingRequest: continuous preview / recording
	//   - capture:             single-shot capture
	// ---------------------------------------------------------------
	//
	// var session *camera.CaptureSession
	//   // created via device.CreateCaptureSession(container, &callbacks)
	//
	// // Start continuous preview:
	// //   session.SetRepeatingRequest(&captureCallbacks, 1, &request)
	//
	// // Later, stop it:
	// //   session.StopRepeating()
	//
	// defer session.Close()

	// ---------------------------------------------------------------
	// Cleanup order
	//
	// Resources must be released in reverse creation order:
	//   1. session.Close()        - stop capture pipeline
	//   2. container.Close()      - free output container
	//   3. sessionOutput.Close()  - free session output
	//   4. target.Close()         - free output target
	//   5. request.Close()        - free capture request
	//   6. device.Close()         - close camera device
	//   7. mgr.Close()            - delete camera manager
	//
	// All Close methods are idempotent and nil-safe.
	// ---------------------------------------------------------------

	// ---------------------------------------------------------------
	// Error handling
	//
	// Camera operations return camera.Error values. Key codes:
	//   camera.ErrPermissionDenied   (-10013) - CAMERA permission not granted
	//   camera.ErrCameraInUse        (-10010) - another app holds the camera
	//   camera.ErrCameraDisconnected (-10002) - camera physically removed
	//   camera.ErrCameraDevice       (-10005) - fatal camera hardware error
	//   camera.ErrSessionClosed      (-10007) - session invalidated
	//
	// Use errors.Is to check for specific failures:
	//   if errors.Is(err, camera.ErrPermissionDenied) { ... }
	// ---------------------------------------------------------------

	log.Println("pipeline overview complete -- see source comments for details")
}
