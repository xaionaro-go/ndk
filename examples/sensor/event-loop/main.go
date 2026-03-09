// Sensor event loop pattern combining sensor and looper packages.
//
// Demonstrates sensor discovery and type enumeration using the idiomatic
// sensor package, and documents the complete ALooper-based event queue
// pattern in comments. The event queue creation requires capi-level access
// (ASensorManager_createEventQueue), so the actual event loop is shown as
// a documented pattern rather than executable code.
//
// The executable portion of this example:
//   1. Gets the sensor manager singleton.
//   2. Queries the default accelerometer and prints its properties.
//   3. Enumerates common sensor types, checking availability on the device.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/sensor"
)

// probeSensor checks whether a sensor type is present on the device.
// DefaultSensor always returns a non-nil wrapper; if the underlying
// C pointer is NULL, calling Name() will panic. We recover from the
// panic to treat it as "not available."
func probeSensor(
	mgr *sensor.Manager,
	sensorType sensor.Type,
) (name, vendor string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	s := mgr.DefaultSensor(int32(sensorType))
	name = s.Name()
	vendor = s.Vendor()
	if name == "" || vendor == "" {
		return "", "", false
	}
	return name, vendor, true
}

func main() {
	// --- 1. Get sensor manager singleton ---
	mgr := sensor.GetInstance()
	log.Println("sensor manager obtained")

	// --- 2. Query default accelerometer ---
	accel := mgr.DefaultSensor(int32(sensor.Accelerometer))
	name := accel.Name()
	if name == "" {
		log.Fatal("no default accelerometer on this device")
	}

	fmt.Println("Default accelerometer:")
	fmt.Printf("  Name:       %s\n", name)
	fmt.Printf("  Vendor:     %s\n", accel.Vendor())
	fmt.Printf("  Type:       %s (%d)\n", sensor.Accelerometer, int32(sensor.Accelerometer))
	fmt.Printf("  Resolution: %g\n", accel.Resolution())
	fmt.Printf("  Min delay:  %d us\n", accel.MinDelay())

	minDelay := accel.MinDelay()
	if minDelay > 0 {
		maxHz := 1_000_000.0 / float64(minDelay)
		fmt.Printf("  Max rate:   %.1f Hz\n", maxHz)
	} else {
		fmt.Println("  Max rate:   on-change sensor (no periodic sampling)")
	}
	fmt.Println()

	// --- 3. Enumerate common sensor types ---
	//
	// Try each common sensor type and report whether it is available.
	type sensorInfo struct {
		label      string
		sensorType sensor.Type
	}

	allTypes := []sensorInfo{
		{"Accelerometer", sensor.Accelerometer},
		{"Magnetic Field", sensor.MagneticField},
		{"Gyroscope", sensor.Gyroscope},
		{"Light", sensor.Light},
		{"Pressure", sensor.Pressure},
		{"Proximity", sensor.Proximity},
		{"Gravity", sensor.Gravity},
		{"Linear Acceleration", sensor.LinearAcceleration},
		{"Rotation Vector", sensor.RotationVector},
		{"Relative Humidity", sensor.RelativeHumidity},
		{"Ambient Temperature", sensor.AmbientTemperature},
		{"Game Rotation Vector", sensor.GameRotationVector},
		{"Step Detector", sensor.StepDetector},
		{"Step Counter", sensor.StepCounter},
		{"Geomagnetic Rotation", sensor.GeomagneticRotationVector},
		{"Heart Rate", sensor.HeartRate},
		{"Hinge Angle", sensor.HingeAngle},
		{"Head Tracker", sensor.HeadTracker},
		{"Heading", sensor.Heading},
	}

	fmt.Println("Sensor type availability:")
	available := 0
	for _, info := range allTypes {
		name, vendor, ok := probeSensor(mgr, info.sensorType)
		if ok {
			fmt.Printf("  [+] %-25s  %s (%s)\n", info.label, name, vendor)
			available++
		} else {
			fmt.Printf("  [-] %-25s  not available\n", info.label)
		}
	}
	fmt.Printf("\n%d of %d sensor types available\n", available, len(allTypes))

	// --- 4. Event queue setup pattern (documented) ---
	//
	// Reading continuous sensor data on Android requires an ALooper-based
	// event queue. The complete pattern is documented below. The event
	// queue creation function (ASensorManager_createEventQueue) is not
	// yet exposed in the idiomatic sensor package, so this section serves
	// as a reference for the expected API shape.
	//
	// Step 1 -- Lock the goroutine to an OS thread (ALooper is thread-local):
	//
	//   runtime.LockOSThread()
	//   defer runtime.UnlockOSThread()
	//
	// Step 2 -- Prepare a looper for the current thread:
	//
	//   import "github.com/xaionaro-go/ndk/looper"
	//
	//   const prepareAllowNonCallbacks = int32(1)
	//   lp := looper.Prepare(prepareAllowNonCallbacks)
	//   defer lp.Close()
	//
	// Step 3 -- Create an event queue bound to the looper:
	//
	//   // ASensorManager_createEventQueue takes the looper, an ident
	//   // (integer to distinguish event sources), and optional callback
	//   // and user-data pointers (nil for poll-based reading).
	//   const sensorIdent = 1
	//   queue := mgr.CreateEventQueue(lp.Pointer(), sensorIdent, nil, nil)
	//   defer mgr.DestroyEventQueue(queue)
	//
	// Step 4 -- Enable the sensor on the queue and set sampling rate:
	//
	//   queue.EnableSensor(accel)
	//   queue.SetEventRate(accel, 100000) // 100 ms between events
	//
	// Step 5 -- Poll the looper and read events:
	//
	//   for {
	//       ident := looper.PollOnce(-1, nil, nil, nil)
	//       if ident == sensorIdent {
	//           // Read ASensorEvent structs from the queue.
	//           // Each event contains: type, timestamp, and a data
	//           // union with up to 16 float32 values.
	//           //
	//           // For accelerometer:
	//           //   data[0] = acceleration X (m/s^2)
	//           //   data[1] = acceleration Y (m/s^2)
	//           //   data[2] = acceleration Z (m/s^2)
	//           for queue.HasEvents() > 0 {
	//               // Read and process events.
	//           }
	//       }
	//   }
	//
	// Step 6 -- Disable sensor and destroy queue (handled by defers):
	//
	//   queue.DisableSensor(accel)

	fmt.Println("\nEvent queue setup pattern documented in source comments.")
	fmt.Println("See main.go for the full ALooper-based sensor read loop pattern.")
}
