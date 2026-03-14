// Audio injector: streams a sine wave to the Android emulator's virtual
// microphone via the gRPC injectAudio API.
//
// Usage: audio-inject -port 8554 -freq 440 -duration 3s
package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	emupb "github.com/xaionaro-go/ndk/tests/e2e/audio-inject/proto/emupb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.Int("port", 8554, "emulator gRPC port")
	freq := flag.Float64("freq", 440, "sine wave frequency in Hz")
	dur := flag.Duration("duration", 3*time.Second, "injection duration")
	sampleRate := flag.Uint64("rate", 48000, "sample rate in Hz")
	flag.Parse()

	addr := fmt.Sprintf("localhost:%d", *port)
	log.Printf("connecting to emulator gRPC at %s", addr)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc dial: %v", err)
	}
	defer conn.Close()

	client := emupb.NewEmulatorControllerClient(conn)
	stream, err := client.InjectAudio(context.Background())
	if err != nil {
		log.Fatalf("InjectAudio: %v", err)
	}

	// Send audio in 10ms chunks.
	samplesPerChunk := int(*sampleRate / 100)
	totalSamples := int(float64(*sampleRate) * dur.Seconds())
	sent := 0
	phase := 0.0
	phaseInc := 2.0 * math.Pi * *freq / float64(*sampleRate)

	format := &emupb.AudioFormat{
		SamplingRate: *sampleRate,
		Channels:     emupb.AudioFormat_Mono,
		Format:       emupb.AudioFormat_AUD_FMT_S16,
	}

	for sent < totalSamples {
		n := samplesPerChunk
		if sent+n > totalSamples {
			n = totalSamples - sent
		}

		buf := make([]byte, n*2) // 2 bytes per S16LE sample
		for i := 0; i < n; i++ {
			sample := int16(math.Sin(phase) * 0.8 * 32767)
			binary.LittleEndian.PutUint16(buf[i*2:], uint16(sample))
			phase += phaseInc
		}

		pkt := &emupb.AudioPacket{
			Audio: buf,
		}
		if sent == 0 {
			pkt.Format = format
		}

		if err := stream.Send(pkt); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("send: %v", err)
		}
		sent += n

		// Pace at roughly real-time to avoid overflowing the 300ms buffer.
		time.Sleep(10 * time.Millisecond)
	}

	if _, err := stream.CloseAndRecv(); err != nil && err != io.EOF {
		log.Printf("close stream: %v", err)
	}
	log.Printf("injected %d samples (%.2fs) at %d Hz", sent, float64(sent)/float64(*sampleRate), *sampleRate)
}
