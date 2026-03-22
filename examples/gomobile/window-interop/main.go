// gomobile bind handle interop: round-trip multiple NDK handle types through int64.
//
// Demonstrates passing different NDK handle types (looper, config, sensor)
// across the gomobile bind boundary using int64 transport. Each handle is
// created, converted to int64, restored, and used — proving the round-trip
// works for the full type hierarchy.
//
// The int64 pattern applies to ALL NDK handle types, including window and
// EGL handles that are only available inside an Activity:
//
//	// window.Window (obtained from onNativeWindowCreated callback)
//	func NewRenderer(windowHandle int64) *Renderer {
//	    win := window.NewWindowFromUintPtr(uintptr(windowHandle))
//	    ...
//	}
//
//	// EGL types (obtained from eglGetDisplay, eglCreateContext, etc.)
//	func (r *Renderer) EGLDisplayHandle() int64 {
//	    return int64(egl.EGLDisplayToUintPtr(r.display))
//	}
//
// This program must run on an Android device.
package main

import (
	"log"
	"runtime"

	"github.com/AndroidGoLab/ndk/config"
	"github.com/AndroidGoLab/ndk/looper"
	"github.com/AndroidGoLab/ndk/sensor"
)

// simulateGomobileBoundary converts a uintptr to int64 and back,
// simulating what happens when a native handle crosses the Go/Java
// boundary via gomobile bind (Go uintptr -> Java long -> Go uintptr).
func simulateGomobileBoundary(original uintptr) uintptr {
	javaLong := int64(original)
	return uintptr(javaLong)
}

func main() {
	log.Println("=== gomobile handle interop ===")

	// --- ALooper: thread-local event loop handle ---

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lp := looper.Prepare(int32(looper.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS))
	if lp == nil {
		log.Fatal("failed to prepare looper")
	}
	defer func() {
		_ = lp.Close()
	}()

	lpRestored := looper.NewLooperFromUintPtr(simulateGomobileBoundary(lp.UintPtr()))
	if lpRestored.UintPtr() != lp.UintPtr() {
		log.Fatalf("FAIL: looper pointer mismatch")
	}
	log.Printf("looper round-trip OK: 0x%x", lp.UintPtr())

	// Use the restored looper: wake and poll.
	lpRestored.Wake()
	result := looper.PollOnce(0, nil, nil, nil)
	log.Printf("looper poll after wake: %d (WAKE=%d)", result, looper.ALOOPER_POLL_WAKE)

	// --- AConfiguration: device configuration handle ---

	cfg := config.NewConfig()
	defer func() {
		_ = cfg.Close()
	}()

	cfgRestored := config.NewConfigFromUintPtr(simulateGomobileBoundary(cfg.UintPtr()))
	if cfgRestored.UintPtr() != cfg.UintPtr() {
		log.Fatalf("FAIL: config pointer mismatch")
	}
	log.Printf("config round-trip OK: 0x%x", cfg.UintPtr())

	// Use the restored config handle.
	log.Printf("config density via restored handle: %d", cfgRestored.Density())
	log.Printf("config sdk version via restored handle: %d", cfgRestored.SdkVersion())

	// --- ASensorManager: singleton handle ---

	mgr := sensor.ASensorManager_getInstanceForPackage("com.example.gomobile")

	mgrRestored := sensor.NewManagerFromUintPtr(simulateGomobileBoundary(mgr.UintPtr()))
	if mgrRestored.UintPtr() != mgr.UintPtr() {
		log.Fatalf("FAIL: sensor manager pointer mismatch")
	}
	log.Printf("sensor manager round-trip OK: 0x%x", mgr.UintPtr())

	// Use the restored sensor manager.
	accel := mgrRestored.DefaultSensor(sensor.Accelerometer)
	if accel.UintPtr() != 0 {
		log.Printf("accelerometer via restored manager: %q", accel.Name())
	} else {
		log.Println("accelerometer not available on this device")
	}

	log.Println("handle-interop done")
}
