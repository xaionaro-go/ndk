// E2E test: exercises the full AAudio input stream lifecycle at 48 kHz.
//
// The test opens an AAudio capture stream, reads ~1 second of audio,
// and verifies the lifecycle completes without errors. If a 440 Hz
// tone is being injected (via host PulseAudio or emulator gRPC), the
// test also verifies frequency detection via the Goertzel algorithm.
//
// Pass with -detect-tone to require 440 Hz detection (exit 1 on failure).
// Without the flag, the test passes as long as AAudio works correctly.
//
// Exit 0 = PASS, exit 1 = FAIL.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"
	"unsafe"

	"github.com/xaionaro-go/ndk/audio"
)

// goertzelMagnitude computes the magnitude of a specific frequency in
// a buffer of int16 PCM samples using the Goertzel algorithm.
func goertzelMagnitude(samples []int16, targetFreq, sampleRate float64) float64 {
	n := len(samples)
	k := int(0.5 + float64(n)*targetFreq/sampleRate)
	w := 2.0 * math.Pi * float64(k) / float64(n)
	coeff := 2.0 * math.Cos(w)

	var s0, s1, s2 float64
	for _, sample := range samples {
		s0 = coeff*s1 - s2 + float64(sample)
		s2 = s1
		s1 = s0
	}

	return s1*s1 + s2*s2 - coeff*s1*s2
}

func main() {
	detectTone := flag.Bool("detect-tone", false, "require 440 Hz tone detection (fail if not found)")
	flag.Parse()

	builder, err := audio.NewStreamBuilder()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: create builder: %v\n", err)
		os.Exit(1)
	}
	defer builder.Close()

	builder.
		SetDirection(audio.Input).
		SetSampleRate(48000).
		SetChannelCount(1).
		SetFormat(audio.PcmI16).
		SetPerformanceMode(audio.LowLatency).
		SetSharingMode(audio.Shared)

	stream, err := builder.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: open stream: %v\n", err)
		os.Exit(1)
	}
	defer stream.Close()

	actualRate := float64(stream.SampleRate())
	fmt.Printf("stream opened: rate=%.0f Hz, channels=%d\n", actualRate, stream.ChannelCount())

	if err := stream.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: start stream: %v\n", err)
		os.Exit(1)
	}

	// Read ~1 second of audio.
	totalFrames := int32(actualRate)
	buf := make([]int16, 1024)
	bufBytes := unsafe.Slice((*byte)(unsafe.Pointer(&buf[0])), len(buf)*int(unsafe.Sizeof(buf[0])))
	var captured []int16

	for int32(len(captured)) < totalFrames {
		framesToRead := int32(len(buf))
		if remaining := totalFrames - int32(len(captured)); remaining < framesToRead {
			framesToRead = remaining
		}

		n, err := stream.Read(bufBytes, framesToRead, time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL: read: %v\n", err)
			os.Exit(1)
		}
		if n == 0 {
			fmt.Fprintf(os.Stderr, "FAIL: read returned 0 frames\n")
			os.Exit(1)
		}

		captured = append(captured, buf[:n]...)
	}

	if err := stream.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "WARN: stop: %v\n", err)
	}

	fmt.Printf("captured %d frames\n", len(captured))

	// Compute peak amplitude.
	var peak int16
	for _, s := range captured {
		v := s
		if v < 0 {
			v = -v
		}
		if v > peak {
			peak = v
		}
	}
	fmt.Printf("peak amplitude: %d\n", peak)

	// Goertzel: check 440 Hz vs other frequencies.
	targetFreq := 440.0
	targetMag := goertzelMagnitude(captured, targetFreq, actualRate)

	otherFreqs := []float64{200, 600, 1000, 2000, 5000}
	var maxOtherMag float64
	for _, f := range otherFreqs {
		mag := goertzelMagnitude(captured, f, actualRate)
		fmt.Printf("  %.0f Hz magnitude: %.2e\n", f, mag)
		if mag > maxOtherMag {
			maxOtherMag = mag
		}
	}
	fmt.Printf("  440 Hz magnitude: %.2e\n", targetMag)
	fmt.Printf("  max other magnitude: %.2e\n", maxOtherMag)

	toneDetected := targetMag > 1e10 && (maxOtherMag == 0 || targetMag/maxOtherMag >= 10)

	if toneDetected {
		fmt.Println("PASS: 440 Hz tone detected")
		return
	}

	if *detectTone {
		if targetMag < 1e10 {
			fmt.Fprintf(os.Stderr, "FAIL: 440 Hz energy too low (%.2e < 1e10)\n", targetMag)
		} else {
			fmt.Fprintf(os.Stderr, "FAIL: 440 Hz not dominant (ratio=%.1f, need >=10)\n", targetMag/maxOtherMag)
		}
		os.Exit(1)
	}

	fmt.Println("PASS: AAudio lifecycle OK (tone detection skipped — no injected signal)")
}
