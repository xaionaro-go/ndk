// Basic AAudio playback example.
//
// Demonstrates how to create an AAudio stream configured for audio output,
// write silence to it, and tear it down. The StreamBuilder uses a chaining
// pattern so all configuration can be expressed in a single fluent call.
//
// This program must run on an Android device with AAudio support.
package main

import (
	"log"
	"time"
	"unsafe"

	"github.com/xaionaro-go/ndk/audio"
)

func main() {
	const (
		sampleRate   = 44100
		channelCount = 2                    // stereo
		durationSec  = 2                    // seconds of silence to write
		writeTimeout = time.Second
	)

	// Create a builder and configure it for playback in one chain.
	builder, err := audio.NewStreamBuilder()
	if err != nil {
		log.Fatalf("create stream builder: %v", err)
	}
	defer builder.Close()

	builder.
		SetDirection(audio.Output).
		SetSampleRate(sampleRate).
		SetChannelCount(channelCount).
		SetFormat(audio.PcmFloat).
		SetPerformanceMode(audio.LowLatency).
		SetSharingMode(audio.Shared)

	// Open the stream from the configured builder.
	stream, err := builder.Open()
	if err != nil {
		log.Fatalf("open stream: %v", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Printf("close stream: %v", err)
		}
	}()

	log.Printf("stream opened: %d Hz, %d ch, state=%s",
		stream.SampleRate(), stream.ChannelCount(), stream.State())

	if err := stream.Start(); err != nil {
		log.Fatalf("start stream: %v", err)
	}

	// Prepare a buffer of silence (all zeros). Each frame has channelCount
	// float32 samples, so the buffer holds one burst worth of frames.
	framesPerBurst := stream.FramesPerBurst()
	buf := make([]float32, int(framesPerBurst)*channelCount)
	bufBytes := unsafe.Slice((*byte)(unsafe.Pointer(&buf[0])), len(buf)*int(unsafe.Sizeof(buf[0])))

	totalFrames := int32(sampleRate * durationSec)
	written := int32(0)

	for written < totalFrames {
		framesToWrite := framesPerBurst
		if remaining := totalFrames - written; remaining < framesToWrite {
			framesToWrite = remaining
		}

		n, err := stream.Write(bufBytes, framesToWrite, writeTimeout)
		if err != nil {
			log.Fatalf("write: %v", err)
		}
		written += n
	}

	log.Printf("wrote %d frames of silence", written)
	log.Printf("xrun count: %d", stream.XRunCount())

	if err := stream.Stop(); err != nil {
		log.Fatalf("stop stream: %v", err)
	}

	log.Println("playback finished")
}
