// Function-level tracing with defer.
//
// Demonstrates a reusable pattern for instrumenting entire functions: call
// BeginSection at the top and defer EndSection so the section is always
// closed, even if the function returns early or panics.
//
// This program must run on an Android device. Trace events are only recorded
// when a trace session is active (e.g. via `perfetto` or `atrace`).
package main

import (
	"fmt"
	"time"

	"github.com/AndroidGoLab/ndk/trace"
)

// traceSection begins a named trace section and returns a function that ends
// it. The intended usage is:
//
//	defer traceSection("myWork")()
func traceSection(name string) func() {
	trace.BeginSection(name)
	return func() {
		trace.EndSection()
	}
}

// loadConfig simulates reading configuration from disk.
func loadConfig() map[string]string {
	defer traceSection("loadConfig")()

	time.Sleep(5 * time.Millisecond) // simulate I/O
	return map[string]string{
		"server": "192.168.1.1",
		"port":   "8080",
	}
}

// connectToServer simulates establishing a network connection.
func connectToServer(addr string) error {
	defer traceSection("connectToServer")()

	fmt.Printf("connecting to %s...\n", addr)
	time.Sleep(15 * time.Millisecond) // simulate network latency
	return nil
}

// processItems simulates a batch processing loop. Each iteration gets its
// own nested trace section so individual items are visible in the timeline.
func processItems(n int) int {
	defer traceSection("processItems")()

	total := 0
	for i := 0; i < n; i++ {
		trace.BeginSection("processItem")

		// Simulate per-item work.
		time.Sleep(2 * time.Millisecond)
		total += i + 1

		trace.EndSection()
	}
	return total
}

func main() {
	defer traceSection("main")()

	if trace.IsEnabled() {
		fmt.Println("tracing is enabled")
	} else {
		fmt.Println("tracing is not enabled (events will be silently dropped)")
	}

	cfg := loadConfig()
	fmt.Printf("config loaded: %v\n", cfg)

	addr := cfg["server"] + ":" + cfg["port"]
	if err := connectToServer(addr); err != nil {
		fmt.Printf("connection failed: %v\n", err)
		return // deferred EndSection still fires
	}
	fmt.Println("connected")

	total := processItems(5)
	fmt.Printf("processed 5 items, total=%d\n", total)

	fmt.Println("done")
}
