// Stream information query example.
//
// Demonstrates how to open an AAudio stream and inspect its runtime
// properties. The actual values returned by SampleRate, ChannelCount, etc.
// may differ from the values requested on the builder because AAudio is
// free to choose the closest supported configuration.
//
// This program must run on an Android device with AAudio support.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/audio"
)

func main() {
	builder, err := audio.NewStreamBuilder()
	if err != nil {
		log.Fatalf("create stream builder: %v", err)
	}
	defer builder.Close()

	// Request specific parameters; AAudio may negotiate different values.
	builder.
		SetDirection(audio.Output).
		SetSampleRate(48000).
		SetChannelCount(2).
		SetFormat(audio.PcmFloat).
		SetPerformanceMode(audio.LowLatency)

	stream, err := builder.Open()
	if err != nil {
		log.Fatalf("open stream: %v", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Printf("close stream: %v", err)
		}
	}()

	// Query the negotiated stream properties.
	fmt.Println("Stream properties:")
	fmt.Printf("  Sample rate:      %d Hz\n", stream.SampleRate())
	fmt.Printf("  Channel count:    %d\n", stream.ChannelCount())
	fmt.Printf("  State:            %s\n", stream.State())
	fmt.Printf("  Frames per burst: %d\n", stream.FramesPerBurst())
	fmt.Printf("  XRun count:       %d\n", stream.XRunCount())

	// Start the stream and observe the state change.
	if err := stream.Start(); err != nil {
		log.Fatalf("start stream: %v", err)
	}
	fmt.Printf("\nAfter Start():\n")
	fmt.Printf("  State:      %s\n", stream.State())
	fmt.Printf("  XRun count: %d\n", stream.XRunCount())

	if err := stream.Stop(); err != nil {
		log.Fatalf("stop stream: %v", err)
	}
	fmt.Printf("\nAfter Stop():\n")
	fmt.Printf("  State: %s\n", stream.State())
}
