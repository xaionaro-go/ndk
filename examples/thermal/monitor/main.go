// Example: thermal status monitoring.
//
// Demonstrates how to acquire a thermal manager, query the current thermal
// status, and use the result to make throttling decisions. The main loop
// simulates a game loop that polls thermal status each frame and adjusts
// workload accordingly.
//
// This program must run on an Android device with thermal API support.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xaionaro-go/ndk/thermal"
)

// throttleForStatus returns a workload scale factor (0.0 to 1.0) based on
// the current thermal status. A game engine would use this to reduce draw
// distance, particle counts, or frame rate targets.
func throttleForStatus(status thermal.ThermalStatus) float64 {
	switch status {
	case thermal.StatusNone:
		return 1.0
	case thermal.StatusLight:
		return 0.9
	case thermal.StatusModerate:
		return 0.7
	case thermal.StatusSevere:
		return 0.5
	case thermal.StatusCritical:
		return 0.3
	case thermal.StatusEmergency:
		return 0.1
	case thermal.StatusShutdown:
		return 0.0
	default:
		return 0.5
	}
}

func main() {
	// Acquire a thermal manager handle.
	mgr := thermal.NewManager()
	defer mgr.Close()

	log.Println("thermal manager acquired")

	// Query the current thermal status once upfront.
	raw := mgr.CurrentStatus()
	status := thermal.ThermalStatus(raw)
	log.Printf("current thermal status: %s (%d)", status, raw)

	// Simulate a game loop that polls thermal status periodically.
	// In a real application this loop would be driven by the frame cadence,
	// and the thermal check would happen once every N frames rather than
	// every frame.
	const (
		totalFrames    = 300
		framesPerCheck = 60
		frameDuration  = 16 * time.Millisecond // ~60 fps
	)

	log.Printf("starting game loop: %d frames, checking thermal every %d frames",
		totalFrames, framesPerCheck)

	for frame := 0; frame < totalFrames; frame++ {
		if frame%framesPerCheck == 0 {
			raw = mgr.CurrentStatus()
			status = thermal.ThermalStatus(raw)
			scale := throttleForStatus(status)
			fmt.Printf("  [frame %3d] thermal=%s  workload=%.0f%%\n",
				frame, status, scale*100)
		}

		// Placeholder for actual frame work. A real game loop would render
		// the scene here with the throttled workload parameters.
		time.Sleep(frameDuration)
	}

	log.Println("game loop finished")
}
