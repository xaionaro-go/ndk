// Audio recording (capture) example.
//
// Opens an AAudio input stream at 48 kHz mono 16-bit PCM, reads one second
// of audio frames, and prints basic statistics (total frames, peak amplitude).
//
// This program must run on an Android device with AAudio support and
// the RECORD_AUDIO permission granted.
package main

import (
	"log"
	"math"
	"time"
	"unsafe"

	"github.com/xaionaro-go/ndk/audio"
)

func main() {
	builder, err := audio.NewStreamBuilder()
	if err != nil {
		log.Fatalf("create stream builder: %v", err)
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
		log.Fatalf("open stream: %v", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Printf("close stream: %v", err)
		}
	}()

	rate := stream.SampleRate()
	log.Printf("capture stream opened (rate=%d Hz, ch=%d)", rate, stream.ChannelCount())

	if err := stream.Start(); err != nil {
		log.Fatalf("start stream: %v", err)
	}

	// Read approximately 1 second of audio.
	totalFrames := rate
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
			log.Fatalf("read: %v", err)
		}
		captured = append(captured, buf[:n]...)
	}

	if err := stream.Stop(); err != nil {
		log.Fatalf("stop stream: %v", err)
	}

	// Compute peak amplitude.
	var peak int16
	for _, s := range captured {
		if s < 0 {
			s = -s
		}
		if s > peak {
			peak = s
		}
	}

	log.Printf("captured %d frames", len(captured))
	log.Printf("peak amplitude: %d (%.1f dBFS)", peak, 20*math.Log10(float64(peak)/32767.0))
	log.Println("recording example finished")
}
