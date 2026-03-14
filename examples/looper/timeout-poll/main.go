// Timeout-based ALooper polling example.
//
// Demonstrates how to use PollOnce with a finite timeout so the event loop
// returns periodically even when no events are pending. This pattern is useful
// for performing periodic work (heartbeats, progress checks, housekeeping)
// alongside event-driven I/O.
//
// This program must run on an Android device with NDK looper support.
package main

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/xaionaro-go/ndk/looper"
)

// ALooper_prepare constants.
const prepareAllowNonCallbacks = int32(1) // ALOOPER_PREPARE_ALLOW_NON_CALLBACKS

// ALooper_pollOnce return values.
const (
	pollWake     = int32(-1) // ALOOPER_POLL_WAKE
	pollCallback = int32(-2) // ALOOPER_POLL_CALLBACK
	pollTimeout  = int32(-3) // ALOOPER_POLL_TIMEOUT
	pollError    = int32(-4) // ALOOPER_POLL_ERROR
)

func pollResultString(r int32) string {
	switch r {
	case pollWake:
		return "WAKE"
	case pollCallback:
		return "CALLBACK"
	case pollTimeout:
		return "TIMEOUT"
	case pollError:
		return "ERROR"
	default:
		return "FD_EVENT"
	}
}

func main() {
	// Lock this goroutine to its OS thread because ALooper is thread-local.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lp := looper.Prepare(prepareAllowNonCallbacks)
	if lp == nil {
		log.Fatal("failed to prepare looper")
	}
	defer func() {
		if err := lp.Close(); err != nil {
			log.Printf("close looper: %v", err)
		}
	}()

	const (
		timeout    = 200 * time.Millisecond // poll timeout
		iterations = 5                      // number of poll cycles to run
	)

	log.Printf("polling %d times with %v timeout", iterations, timeout)

	var fd, events int32
	var data unsafe.Pointer

	for i := 1; i <= iterations; i++ {
		result := looper.PollOnce(timeout, &fd, &events, &data)
		log.Printf("poll %d/%d: result=%s (%d)", i, iterations,
			pollResultString(result), result)

		if result == pollTimeout {
			// No events arrived within the timeout window. This is the
			// expected path in this example since nothing is waking the
			// looper or registering file descriptors.
			log.Printf("  -> timeout elapsed, performing periodic work")
			continue
		}

		if result == pollWake {
			log.Printf("  -> looper was woken externally")
			continue
		}

		if result == pollError {
			log.Fatal("  -> poll error, exiting")
		}

		// Non-negative result means a monitored fd has events.
		if result >= 0 {
			log.Printf("  -> fd=%d events=%d", result, events)
		}
	}

	log.Println("timeout poll finished")
}
