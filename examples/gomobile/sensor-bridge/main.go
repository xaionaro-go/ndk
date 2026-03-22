// gomobile bind sensor bridge: round-trip NDK handles through int64.
//
// Demonstrates the pattern for passing NDK sensor handles across the
// Go/Java boundary via gomobile bind. gomobile bind does not support
// unsafe.Pointer or uintptr, so native handles are transported as int64
// (Java long).
//
// This example obtains a real ASensorManager, converts it to int64
// (simulating the gomobile boundary), converts back, and uses the
// restored handle to query sensors — proving the round-trip is lossless.
//
// A gomobile-bindable Go package would expose these operations as
// exported methods on a struct:
//
//	package sensorbridge
//
//	type Bridge struct{ mgr *sensor.Manager }
//
//	func NewBridge(h int64) *Bridge {
//	    return &Bridge{mgr: sensor.NewManagerFromUintPtr(uintptr(h))}
//	}
//
//	func (b *Bridge) AccelerometerName() string {
//	    return b.mgr.DefaultSensor(sensor.Accelerometer).Name()
//	}
//
//	func (b *Bridge) Handle() int64 {
//	    return int64(b.mgr.UintPtr())
//	}
//
// Java side:
//
//	long ptr = nativeGetSensorManager(); // from JNI
//	Bridge b = Sensorbridge.newBridge(ptr);
//	String name = b.accelerometerName();
//
// This program must run on an Android device.
package main

import (
	"log"

	"github.com/AndroidGoLab/ndk/sensor"
)

func main() {
	log.Println("=== gomobile sensor bridge ===")

	// Step 1: Obtain a real sensor manager from the NDK.
	mgr := sensor.ASensorManager_getInstanceForPackage("com.example.gomobile")
	original := mgr.UintPtr()
	log.Printf("sensor manager obtained: 0x%x", original)

	// Step 2: Simulate the gomobile boundary — convert to int64.
	// In a real app, this int64 crosses the Go/Java boundary.
	javaHandle := int64(original)
	log.Printf("converted to int64 (Java long): %d", javaHandle)

	// Step 3: On the "Java side", pass the int64 back to Go.
	// Reconstruct the sensor.Manager from int64.
	restored := sensor.NewManagerFromUintPtr(uintptr(javaHandle))
	log.Printf("restored from int64: 0x%x", restored.UintPtr())

	// Step 4: Verify the round-trip — pointers must match.
	if restored.UintPtr() != original {
		log.Fatalf("FAIL: pointer mismatch: original=0x%x restored=0x%x", original, restored.UintPtr())
	}
	log.Println("round-trip OK: pointers match")

	// Step 5: Use the restored handle to query real sensors.
	sensorTypes := []struct {
		name     string
		typeCode sensor.Type
	}{
		{"Accelerometer", sensor.Accelerometer},
		{"Gyroscope", sensor.Gyroscope},
		{"MagneticField", sensor.MagneticField},
		{"Light", sensor.Light},
		{"Proximity", sensor.Proximity},
	}

	for _, st := range sensorTypes {
		s := restored.DefaultSensor(st.typeCode)
		if s.UintPtr() == 0 {
			log.Printf("  %-15s not available", st.name)
			continue
		}
		log.Printf("  %-15s name=%q vendor=%q", st.name, s.Name(), s.Vendor())
	}

	log.Println("sensor-bridge done")
}
