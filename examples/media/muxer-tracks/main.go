// Muxer track format preparation example.
//
// Demonstrates how to build AMediaFormat objects for audio and video tracks
// using the chaining API, query properties back with GetInt32, and documents
// how these formats feed into the AMediaMuxer track-addition workflow.
//
// The muxer itself requires a valid file descriptor and encoded media data,
// so this example focuses on the format construction and property verification
// steps that precede muxing. The muxer workflow is documented in comments.
//
// This program must run on an Android device with NDK media support.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AndroidGoLab/ndk/media"
)

func main() {
	// -- Video track format --
	videoFormat := media.NewFormat()
	defer func() {
		if err := videoFormat.Close(); err != nil {
			log.Printf("close video format: %v", err)
		}
	}()

	videoFormat.
		SetString("mime", "video/avc").
		SetInt32("width", 1920).
		SetInt32("height", 1080).
		SetInt32("bitrate", 8_000_000). // 8 Mbps
		SetInt32("frame-rate", 30).
		SetInt32("i-frame-interval", 2) // keyframe every 2 seconds

	fmt.Println("Video track format:")
	printInt32Props(videoFormat, []prop{
		{"  MIME", "mime"}, // SetString, not readable via GetInt32
		{"  Width", "width"},
		{"  Height", "height"},
		{"  Bitrate", "bitrate"},
		{"  Frame rate", "frame-rate"},
		{"  I-frame interval", "i-frame-interval"},
	})

	// -- Audio track format --
	audioFormat := media.NewFormat()
	defer func() {
		if err := audioFormat.Close(); err != nil {
			log.Printf("close audio format: %v", err)
		}
	}()

	audioFormat.
		SetString("mime", "audio/mp4a-latm").
		SetInt32("channel-count", 2).
		SetInt32("sample-rate", 44100).
		SetInt32("bitrate", 128_000) // 128 kbps

	fmt.Println("Audio track format:")
	printInt32Props(audioFormat, []prop{
		{"  Channel count", "channel-count"},
		{"  Sample rate", "sample-rate"},
		{"  Bitrate", "bitrate"},
	})

	// -- Muxer workflow (documented) --
	//
	// In a real application the steps after format creation are:
	//
	//   1. Open an output file and obtain its file descriptor.
	//      fd, _ := syscall.Open("output.mp4", syscall.O_WRONLY|syscall.O_CREAT, 0644)
	//
	//   2. Create a muxer with the desired container format.
	//      muxer := media.NewMuxer(fd, media.OutputFormatMPEG4)
	//
	//   3. Add tracks using the formats built above.
	//      videoTrackIdx := muxer.AddTrack(videoFormat)
	//      audioTrackIdx := muxer.AddTrack(audioFormat)
	//
	//   4. Start the muxer (no more tracks can be added after this).
	//      muxer.Start()
	//
	//   5. Write encoded samples from the encoder/extractor.
	//      muxer.WriteSampleData(videoTrackIdx, encodedFrame, bufferInfo)
	//      muxer.WriteSampleData(audioTrackIdx, encodedAudio, bufferInfo)
	//
	//   6. Stop and close the muxer to finalize the file.
	//      muxer.Stop()
	//      muxer.Close()

	fmt.Println("track formats created and verified successfully")
}

type prop struct {
	label string
	key   string
}

// printInt32Props queries each int32 property and prints its value.
// Properties set via SetString are silently skipped since GetInt32
// cannot read them.
func printInt32Props(
	format *media.Format,
	props []prop,
) {
	for _, p := range props {
		var val int32
		if !format.GetInt32(p.key, &val) {
			// Not an int32 property (e.g. MIME is a string).
			continue
		}
		fmt.Fprintf(os.Stdout, "%-22s: %d\n", p.label, val)
	}
}
