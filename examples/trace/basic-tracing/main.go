// Basic systrace example: synchronous sections, asynchronous sections, and counters.
//
// Demonstrates:
//   - Checking whether tracing is currently enabled
//   - Wrapping synchronous work in BeginSection / EndSection
//   - Tracking overlapping asynchronous operations with BeginAsyncSection / EndAsyncSection
//   - Publishing numeric counters visible in the trace timeline with SetCounter
//
// This program must run on an Android device. Trace events are only recorded
// when a trace session is active (e.g. via `perfetto` or `atrace`).
package main

import (
	"fmt"
	"time"

	"github.com/xaionaro-go/ndk/trace"
)

func main() {
	// Check whether a trace session is currently capturing events.
	// All trace functions are safe to call regardless, but this lets you
	// skip expensive argument preparation when nobody is listening.
	if trace.IsEnabled() {
		fmt.Println("tracing is enabled")
	} else {
		fmt.Println("tracing is not enabled (events will be silently dropped)")
	}

	// --- Synchronous section ---
	// BeginSection / EndSection mark a region on the calling thread's
	// timeline. They must be strictly nested (like parentheses).

	trace.BeginSection("computeFrame")

	// Simulate some work.
	total := 0
	for i := 0; i < 1_000_000; i++ {
		total += i
	}

	trace.EndSection()
	fmt.Printf("synchronous section finished (total=%d)\n", total)

	// --- Asynchronous sections ---
	// BeginAsyncSection / EndAsyncSection track work that may overlap
	// across goroutines or threads. The cookie identifies which begin
	// matches which end, so multiple instances of the same section name
	// can be in flight simultaneously.

	const (
		cookieA int32 = 1
		cookieB int32 = 2
	)

	// Start two overlapping async operations.
	trace.BeginAsyncSection("networkRequest", cookieA)
	trace.BeginAsyncSection("networkRequest", cookieB)
	fmt.Println("started two async sections")

	// Simulate the first request finishing before the second.
	time.Sleep(10 * time.Millisecond)
	trace.EndAsyncSection("networkRequest", cookieA)
	fmt.Println("async section A ended")

	time.Sleep(20 * time.Millisecond)
	trace.EndAsyncSection("networkRequest", cookieB)
	fmt.Println("async section B ended")

	// --- Counter ---
	// SetCounter publishes a named integer value that appears as a
	// counter track in the trace viewer. Useful for tracking queue
	// depths, cache sizes, frame counts, etc.

	for frame := int64(1); frame <= 5; frame++ {
		trace.SetCounter("framesRendered", frame)
		time.Sleep(5 * time.Millisecond)
	}
	fmt.Println("counter updated 5 times")

	fmt.Println("done")
}
