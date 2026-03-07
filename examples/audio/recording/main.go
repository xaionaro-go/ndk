// Audio recording (capture) stream setup example.
//
// Demonstrates how to configure an AAudio input stream for recording audio.
// The stream is set to Input direction with 16-bit PCM format at 16 kHz,
// which is a common configuration for voice capture on Android.
//
// After opening and starting the stream, its negotiated properties are
// printed. A real application would read captured frames from the stream
// using a data callback or a read loop.
//
// This program must run on an Android device with AAudio support and
// the RECORD_AUDIO permission granted.
package main

import (
	"log"

	"github.com/xaionaro-go/ndk/audio"
)

func main() {
	builder, err := audio.NewStreamBuilder()
	if err != nil {
		log.Fatalf("create stream builder: %v", err)
	}
	defer builder.Close()

	// Configure for mono voice capture at 16 kHz.
	builder.
		SetDirection(audio.Input).
		SetSampleRate(16000).
		SetChannelCount(1).
		SetFormat(audio.PcmI16).
		SetPerformanceMode(audio.LowLatency).
		SetSharingMode(audio.Shared)

	stream, err := builder.Open()
	if err != nil {
		log.Fatalf("open stream: %v", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Printf("close stream: %v", err)
		}
	}()

	log.Printf("capture stream opened")
	log.Printf("  sample rate:      %d Hz", stream.SampleRate())
	log.Printf("  channel count:    %d", stream.ChannelCount())
	log.Printf("  frames per burst: %d", stream.FramesPerBurst())
	log.Printf("  state:            %s", stream.State())

	if err := stream.Start(); err != nil {
		log.Fatalf("start stream: %v", err)
	}
	log.Printf("  state after start: %s", stream.State())

	// A real application would now read captured audio frames from the
	// stream, for example by using a data callback registered on the
	// builder before Open(), or by calling a read method in a loop.

	if err := stream.Stop(); err != nil {
		log.Fatalf("stop stream: %v", err)
	}
	log.Printf("  state after stop: %s", stream.State())
	log.Printf("  xrun count:       %d", stream.XRunCount())

	log.Println("recording example finished")
}
