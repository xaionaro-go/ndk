// Example: accelerometer sensor detail and event-queue setup.
//
// Obtains the default accelerometer, prints its detailed properties,
// and documents how to set up an event queue for continuous reading.
//
// The Android sensor API requires an ALooper for the event queue. The
// typical pattern is:
//
//  1. Prepare a looper on the current thread (looper.Prepare).
//  2. Create an event queue via the sensor manager, passing the looper.
//  3. Enable the sensor on the queue and set a sampling rate.
//  4. Poll the looper; when data is ready, read events from the queue.
//  5. Disable the sensor and destroy the queue on cleanup.
//
// This example focuses on sensor querying and shows the event-queue
// setup pattern in commented code. The event-reading loop is a sketch
// because it requires a running Android looper thread.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"

	"github.com/AndroidGoLab/ndk/sensor"
)

func main() {
	// --- Obtain the sensor manager (singleton) ---
	mgr := sensor.GetInstance()

	// --- Query the default accelerometer ---
	//
	// DefaultSensor always returns a non-nil *Sensor wrapper. When the
	// device lacks the requested sensor type, the underlying C pointer
	// is NULL and any method call will crash. We call Name() first and
	// treat an empty result (or panic) as "sensor not available."
	accel := mgr.DefaultSensor(sensor.Accelerometer)

	name := accel.Name()
	if name == "" {
		log.Fatal("no default accelerometer on this device")
	}
	vendor := accel.Vendor()
	if vendor == "" {
		log.Fatal("accelerometer returned empty vendor")
	}

	fmt.Println("Accelerometer details:")
	fmt.Printf("  Name:       %s\n", name)
	fmt.Printf("  Vendor:     %s\n", vendor)
	fmt.Printf("  Type:       %s (%d)\n", sensor.Accelerometer, int32(sensor.Accelerometer))
	fmt.Printf("  Resolution: %g\n", accel.Resolution())
	fmt.Printf("  Min delay:  %d us\n", accel.MinDelay())

	// Compute the maximum sampling frequency from the minimum delay.
	minDelay := accel.MinDelay()
	if minDelay > 0 {
		maxHz := 1_000_000.0 / float64(minDelay)
		fmt.Printf("  Max rate:   %.1f Hz\n", maxHz)
	} else {
		fmt.Println("  Max rate:   on-change sensor (no periodic sampling)")
	}

	// --- Event queue setup pattern ---
	//
	// Reading sensor data on Android requires an ALooper-based event
	// queue. The high-level sensor package exposes the EventQueue struct,
	// but creation currently requires the NDK Manager. Below is the
	// documented pattern using ndk packages.
	//
	// Step 1 -- Prepare a looper for the current thread:
	//
	//   import "github.com/AndroidGoLab/ndk/looper"
	//
	//   lp := looper.Prepare(0)
	//   defer lp.Close()
	//
	// Step 2 -- Create an event queue bound to the looper:
	//
	//   // The ident parameter distinguishes event sources when
	//   // multiple file descriptors share a single looper.
	//   // Pass nil for callback and data to use poll-based reading.
	//   const sensorIdent = 1
	//   queue := mgr.CreateEventQueue(looperPtr, sensorIdent, nil, nil)
	//   defer mgr.DestroyEventQueue(queue)
	//
	// Step 3 -- Enable the sensor and set the sampling rate:
	//
	//   queue.EnableSensor(accel)
	//   // Request events every 100 ms (100000 us).
	//   queue.SetEventRate(accel, 100000)
	//
	// Step 4 -- Poll the looper and read events:
	//
	//   for {
	//       // Block until an event is available (-1 = infinite timeout).
	//       ident := looper.PollOnce(-1, nil, nil, nil)
	//       if ident == sensorIdent {
	//           for queue.HasEvents() > 0 {
	//               // Read ASensorEvent structs from the queue.
	//               // (ASensorEvent reading is not yet wrapped at the
	//               // Go level; use unsafe + the C struct layout.)
	//           }
	//       }
	//   }
	//
	// Step 5 -- Disable and clean up (handled by defers above):
	//
	//   queue.DisableSensor(accel)

	fmt.Println()
	fmt.Println("Event queue setup pattern documented in source comments.")
	fmt.Println("See main.go for the full ALooper-based read loop sketch.")
}
