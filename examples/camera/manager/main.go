// Camera manager lifecycle example.
//
// Demonstrates the simplest possible use of the camera package: creating a
// camera manager and closing it. The manager is the entry point to the
// Android Camera2 NDK API -- all camera discovery and device access starts
// here.
//
// Prerequisites:
//   - This program must run on an Android device (API level 24+).
//   - The android.permission.CAMERA runtime permission must be granted before
//     any camera device can be opened. Manager creation itself does not
//     require the permission, but subsequent operations (listing cameras,
//     opening a device) will fail with ErrPermissionDenied (-10013) without it.
//   - On Android 10+ the permission must be requested at runtime through the
//     Java/Kotlin activity layer; it cannot be granted from native code alone.
package main

import (
	"log"

	"github.com/AndroidGoLab/ndk/camera"
)

func main() {
	// Create a camera manager. This allocates the underlying
	// ACameraManager and connects to the camera service.
	mgr := camera.NewManager()
	if mgr == nil {
		log.Fatal("camera.NewManager returned nil")
	}
	log.Println("camera manager created")

	// The manager is ready. In a real application you would now:
	//   1. Query available cameras  (ACameraManager_getCameraIdList)
	//   2. Retrieve camera metadata (ACameraManager_getCameraCharacteristics)
	//   3. Open a camera device     (ACameraManager_openCamera)
	//
	// All of these require android.permission.CAMERA to be granted.
	// Without it the NDK functions return ErrPermissionDenied.

	// Close the manager to release native resources.
	// Close is idempotent -- calling it again is harmless.
	if err := mgr.Close(); err != nil {
		log.Fatalf("manager close: %v", err)
	}
	log.Println("camera manager closed")

	// Demonstrate idempotent close.
	if err := mgr.Close(); err != nil {
		log.Fatalf("second close: %v", err)
	}
	log.Println("second close succeeded (no-op)")
}
