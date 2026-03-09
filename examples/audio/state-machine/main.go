// AAudio stream state machine example.
//
// Demonstrates the lifecycle of an AAudio stream by walking through
// key state transitions and printing the stream state after each:
//
//	Open -> Start -> Stop -> Close
//
// Note that AAudio state transitions are asynchronous: a call to Start()
// requests the transition and the stream may briefly be in the "Starting"
// state before reaching "Started". The state printed here is whatever the
// stream reports immediately after the request returns.
//
// Additional states like Pause, Flush, and their async counterparts
// (Pausing, Flushing) are part of the AAudio state machine. In a real
// application you would use WaitForStateChange to synchronize before
// calling operations that require a specific state (e.g. Flush requires
// the stream to be fully Paused).
//
// This program must run on an Android device with AAudio support.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/audio"
)

// printState logs the current stream state with a label describing the
// operation that was just performed.
func printState(label string, s *audio.Stream) {
	fmt.Printf("  %-20s -> state: %s\n", label, s.State())
}

func main() {
	// Build a simple output stream.
	builder, err := audio.NewStreamBuilder()
	if err != nil {
		log.Fatalf("create stream builder: %v", err)
	}
	defer builder.Close()

	builder.
		SetDirection(audio.Output).
		SetSampleRate(44100).
		SetChannelCount(1).
		SetFormat(audio.PcmI16).
		SetPerformanceMode(audio.LowLatency).
		SetSharingMode(audio.Shared)

	// Open the stream. The state should be Open after this call.
	stream, err := builder.Open()
	if err != nil {
		log.Fatalf("open stream: %v", err)
	}

	fmt.Println("State transitions:")
	printState("Open", stream)

	// Query stream properties set by the system.
	fmt.Printf("  Sample rate:         %d Hz\n", stream.SampleRate())
	fmt.Printf("  Channel count:       %d\n", stream.ChannelCount())
	fmt.Printf("  Frames per burst:    %d\n", stream.FramesPerBurst())
	fmt.Println()

	// Start playback.
	if err := stream.Start(); err != nil {
		log.Fatalf("start: %v", err)
	}
	printState("Start", stream)

	// Stop the stream. Stop is valid from Started or Paused state.
	if err := stream.Stop(); err != nil {
		log.Fatalf("stop: %v", err)
	}
	printState("Stop", stream)

	// Close the stream and release all resources.
	if err := stream.Close(); err != nil {
		log.Fatalf("close: %v", err)
	}
	fmt.Printf("  %-20s -> stream released\n", "Close")

	fmt.Println("all state transitions completed successfully")
}
