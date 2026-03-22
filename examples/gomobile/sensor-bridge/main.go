// Exposing NDK sensors to Java/Kotlin via gomobile bind.
//
// gomobile bind generates Java/Kotlin bindings for exported Go functions.
// It supports signed integers, floats, strings, booleans, []byte, errors,
// interfaces, and structs — but NOT unsafe.Pointer or uintptr.
//
// To pass NDK handles across the Go/Java boundary, use int64 (Java long)
// as the transport type and convert via uintptr inside Go:
//
//	func NewSensorBridge(managerHandle int64) *SensorBridge {
//	    mgr := sensor.NewManagerFromUintPtr(uintptr(managerHandle))
//	    return &SensorBridge{mgr: mgr}
//	}
//
// The Java side obtains native handles as long values (e.g. from JNI)
// and passes them to Go:
//
//	// Java
//	long sensorManagerPtr = nativeGetSensorManager();
//	SensorBridge bridge = Sensorbridge.newSensorBridge(sensorManagerPtr);
//	String name = bridge.defaultAccelerometerName();
//
// # Go library package (bindable via gomobile bind)
//
// The pattern for a gomobile-bindable Go package:
//
//	// Package sensorbridge exposes NDK sensor access to Java/Kotlin.
//	// Build with: gomobile bind -target=android ./sensorbridge
//	package sensorbbridge
//
//	import "github.com/AndroidGoLab/ndk/sensor"
//
//	// SensorBridge wraps a sensor.Manager for cross-language use.
//	// Exported struct fields and methods with supported types are
//	// automatically exposed to Java/Kotlin by gomobile bind.
//	type SensorBridge struct {
//	    mgr *sensor.Manager // unexported: invisible to Java
//	}
//
//	// NewSensorBridge creates a bridge from a raw ASensorManager handle.
//	// The handle is passed as int64 because gomobile bind does not
//	// support unsafe.Pointer or uintptr.
//	func NewSensorBridge(managerHandle int64) *SensorBridge {
//	    mgr := sensor.NewManagerFromUintPtr(uintptr(managerHandle))
//	    return &SensorBridge{mgr: mgr}
//	}
//
//	// DefaultAccelerometerName returns the name of the default accelerometer.
//	func (b *SensorBridge) DefaultAccelerometerName() string {
//	    s := b.mgr.DefaultSensor(sensor.Accelerometer)
//	    return s.Name()
//	}
//
//	// Handle returns the underlying ASensorManager handle for passing
//	// back to Java/JNI code.
//	func (b *SensorBridge) Handle() int64 {
//	    return int64(b.mgr.UintPtr())
//	}
//
// # Java usage
//
//	import sensorbbridge.SensorBridge;
//
//	// Obtain the native ASensorManager* from JNI or the NDK.
//	long ptr = nativeGetSensorManager();
//
//	// Create the Go bridge, passing the handle as a long.
//	SensorBridge bridge = Sensorbbridge.newSensorBridge(ptr);
//
//	// Call Go methods — gomobile generates Java proxy methods.
//	String name = bridge.defaultAccelerometerName();
//	Log.d("Sensor", "Accelerometer: " + name);
//
//	// Retrieve the handle back (e.g. to pass to another native call).
//	long handle = bridge.handle();
//
// # Key rules
//
//   - Use int64 (Java long) to transport native pointers across the boundary.
//   - Convert int64 <-> uintptr inside Go; never expose uintptr to gomobile.
//   - Keep NDK handle fields unexported — gomobile cannot serialize them.
//   - Exported methods must use only gomobile-supported types.
//   - Close/release NDK resources from Go, not from Java.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/AndroidGoLab/ndk/sensor"
)

func main() {
	fmt.Println("=== gomobile bind: sensor bridge ===")
	fmt.Println()

	// Demonstrate the int64 <-> uintptr <-> NDK handle conversion chain.
	//
	// In a real app, the int64 comes from Java (JNI native handle).
	// Here we show the Go-side conversion pattern.

	fmt.Println("Pattern: Java long -> Go int64 -> uintptr -> NDK handle")
	fmt.Println()

	// Step 1: Receive a handle from Java as int64.
	fmt.Println("Step 1: Java passes native handle as long (int64)")
	fmt.Println("  Java:  long ptr = nativeGetSensorManager();")
	fmt.Println("  Java:  SensorBridge bridge = Sensorbbridge.newSensorBridge(ptr);")
	fmt.Println()

	// Step 2: Convert int64 -> uintptr -> NDK type in Go.
	fmt.Println("Step 2: Go converts int64 -> uintptr -> sensor.Manager")
	fmt.Println("  Go:    mgr := sensor.NewManagerFromUintPtr(uintptr(handle))")
	fmt.Println()

	// Step 3: Use the NDK type normally.
	fmt.Println("Step 3: Use NDK types normally in Go")
	fmt.Println("  Go:    s := mgr.DefaultSensor(sensor.Accelerometer)")
	fmt.Println("  Go:    name := s.Name()")
	fmt.Println()

	// Step 4: Return results to Java using supported types.
	fmt.Println("Step 4: Return results as gomobile-supported types (string, int64, error)")
	fmt.Println("  Go:    return s.Name()  // string -> Java String")
	fmt.Println("  Go:    return int64(mgr.UintPtr())  // handle -> Java long")
	fmt.Println()

	// Show the available sensor type constants that a bridge might expose.
	types := []struct {
		name  string
		value sensor.Type
	}{
		{"Accelerometer", sensor.Accelerometer},
		{"Gyroscope", sensor.Gyroscope},
		{"MagneticField", sensor.MagneticField},
		{"Light", sensor.Light},
		{"Proximity", sensor.Proximity},
	}

	fmt.Println("Sensor type constants (use int32 for gomobile):")
	for _, t := range types {
		fmt.Printf("  %-25s = %d\n", t.name, int32(t.value))
	}
	fmt.Println()

	// gomobile bind type mapping summary.
	fmt.Println("gomobile bind type mapping:")
	fmt.Println("  Go int64       <-> Java long       (native handle transport)")
	fmt.Println("  Go string      <-> Java String")
	fmt.Println("  Go int32       <-> Java int         (sensor type constants)")
	fmt.Println("  Go float32     <-> Java float       (sensor values)")
	fmt.Println("  Go float64     <-> Java double")
	fmt.Println("  Go bool        <-> Java boolean")
	fmt.Println("  Go []byte      <-> Java byte[]")
	fmt.Println("  Go error       <-> Java Exception")
	fmt.Println("  Go struct      <-> Java class       (exported methods only)")
	fmt.Println("  Go interface   <-> Java interface")
	fmt.Println()

	fmt.Println("sensor-bridge example complete")
}
