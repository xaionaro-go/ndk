// Nested trace sections and counters example.
//
// Demonstrates how to use the NDK trace API for systrace/perfetto
// integration:
//   - Nested synchronous sections (outer "NetworkRequest" containing
//     an inner "ParseResponse")
//   - Overlapping asynchronous sections for concurrent operations
//   - Performance counters for tracking frame counts and byte throughput
//
// This program must run on an Android device. Trace events are only recorded
// when a trace session is active (e.g. via `perfetto` or `atrace`).
package main

import (
	"fmt"
	"time"

	"github.com/xaionaro-go/ndk/trace"
)

// simulateNetworkRequest wraps a simulated network call in nested trace
// sections: an outer section for the whole request and an inner section
// for response parsing.
func simulateNetworkRequest(url string, bytesReceived *int64) {
	trace.BeginSection("NetworkRequest")

	// Simulate network I/O.
	fmt.Printf("  fetching %s...\n", url)
	time.Sleep(20 * time.Millisecond)
	*bytesReceived += 4096

	// The parse step is a nested section inside the request.
	trace.BeginSection("ParseResponse")
	time.Sleep(5 * time.Millisecond)
	*bytesReceived += 512
	trace.EndSection() // ParseResponse

	trace.EndSection() // NetworkRequest
}

func main() {
	// Check whether a trace session is capturing events.
	if trace.IsEnabled() {
		fmt.Println("tracing is enabled")
	} else {
		fmt.Println("tracing is not enabled (events will be silently dropped)")
	}

	// --- Nested synchronous sections ---
	// Sections must be strictly nested (like parentheses). Here we show a
	// top-level "App" section containing two network requests, each of
	// which has its own nested sections.

	trace.BeginSection("App")
	fmt.Println("started top-level section: App")

	var totalBytes int64

	simulateNetworkRequest("https://api.example.com/data", &totalBytes)
	simulateNetworkRequest("https://api.example.com/config", &totalBytes)

	trace.EndSection() // App
	fmt.Printf("nested sections finished, total bytes: %d\n", totalBytes)

	// --- Overlapping asynchronous sections ---
	// Async sections use a cookie to correlate begin/end across time.
	// Unlike synchronous sections, multiple async sections with the same
	// name can overlap.

	const (
		uploadCookie   int32 = 100
		downloadCookie int32 = 200
	)

	trace.BeginAsyncSection("FileTransfer", uploadCookie)
	fmt.Println("started async upload")

	trace.BeginAsyncSection("FileTransfer", downloadCookie)
	fmt.Println("started async download (overlapping with upload)")

	// Download finishes first.
	time.Sleep(10 * time.Millisecond)
	trace.EndAsyncSection("FileTransfer", downloadCookie)
	fmt.Println("async download completed")

	// Upload finishes second.
	time.Sleep(15 * time.Millisecond)
	trace.EndAsyncSection("FileTransfer", uploadCookie)
	fmt.Println("async upload completed")

	// --- Performance counters ---
	// Counters appear as numeric tracks in the trace viewer. They are
	// useful for tracking quantities that change over time.

	for frame := int64(1); frame <= 10; frame++ {
		trace.SetCounter("frameCount", frame)
		trace.SetCounter("bytesProcessed", totalBytes*frame)
		time.Sleep(2 * time.Millisecond)
	}
	fmt.Println("published 10 counter samples")

	fmt.Println("done")
}
