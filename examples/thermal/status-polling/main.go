// Thermal status polling example.
//
// Demonstrates periodic polling of the device thermal status and using the
// result to make adaptive workload decisions. When the device reaches Severe
// or higher, the simulated workload is progressively reduced.
//
// This program must run on an Android device with thermal API support.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xaionaro-go/ndk/thermal"
)

// workloadScale returns a multiplier (0.0 to 1.0) indicating how much work
// should be performed at the given thermal status. Higher thermal pressure
// yields a lower multiplier.
func workloadScale(status thermal.ThermalStatus) float64 {
	switch status {
	case thermal.StatusNone:
		return 1.0
	case thermal.StatusLight:
		return 0.9
	case thermal.StatusModerate:
		return 0.75
	case thermal.StatusSevere:
		return 0.5
	case thermal.StatusCritical:
		return 0.25
	case thermal.StatusEmergency:
		return 0.1
	case thermal.StatusShutdown:
		return 0.0
	default:
		return 0.5
	}
}

func main() {
	mgr := thermal.NewManager()
	defer mgr.Close()

	log.Println("thermal manager acquired")

	// Read initial status.
	status := mgr.CurrentStatus()
	log.Printf("initial thermal status: %s", status)

	// Poll thermal status at a fixed interval. In a real application this
	// would be driven by a game loop or render tick; here we use a simple
	// timer to keep the example self-contained.
	const (
		pollInterval = 500 * time.Millisecond
		pollCount    = 10
	)

	log.Printf("polling thermal status %d times at %v intervals", pollCount, pollInterval)

	for i := 0; i < pollCount; i++ {
		status = mgr.CurrentStatus()
		scale := workloadScale(status)

		fmt.Printf("  poll %2d: thermal=%s  workload=%.0f%%",
			i+1, status, scale*100)

		// Demonstrate adaptive throttling: log a warning when the device
		// is under significant thermal pressure.
		switch {
		case status >= thermal.StatusSevere:
			fmt.Printf("  ** throttling (severe+)\n")
		case status >= thermal.StatusLight:
			fmt.Printf("  (light pressure)\n")
		default:
			fmt.Printf("\n")
		}

		time.Sleep(pollInterval)
	}

	log.Println("polling finished")
}
