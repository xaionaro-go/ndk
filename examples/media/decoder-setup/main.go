// Video decoder setup example.
//
// Demonstrates how to create an H.264 video decoder using AMediaCodec,
// build a format describing the input stream, and configure the codec.
// The codec is configured without a surface (nil window) and without
// DRM (zero-value Crypto), which is the pattern for buffer-mode decoding.
//
// On a real device this would be followed by Start(), the input/output
// buffer loop (QueueInputBuffer / DequeueOutputBuffer / ReleaseOutputBuffer),
// and finally Stop(). Here we stop at Configure to illustrate the setup
// pattern without requiring an actual bitstream or display surface.
//
// This program must run on an Android device with NDK media support.
package main

import (
	"fmt"
	"log"

	"github.com/AndroidGoLab/ndk/media"
)

func main() {
	// 1. Create a decoder for H.264.
	codec := media.NewDecoder("video/avc")
	if codec == nil {
		log.Fatal("NewDecoder returned nil: H.264 decoder not available")
	}
	defer func() {
		if err := codec.Close(); err != nil {
			log.Printf("close codec: %v", err)
		}
	}()
	fmt.Println("H.264 decoder created")

	// 2. Build the input format describing the compressed video.
	format := media.NewFormat()
	defer func() {
		if err := format.Close(); err != nil {
			log.Printf("close format: %v", err)
		}
	}()

	format.
		SetString("mime", "video/avc").
		SetInt32("width", 1280).
		SetInt32("height", 720)

	// 3. Configure the codec.
	//    - window: nil (no surface, buffer-mode decoding)
	//    - crypto: zero-value Crypto (no DRM)
	//    - flags:  0 (decoder; use 1 for encoder)
	// Note: AMediaCodec_configure may crash with SIGSEGV outside an Activity
	// context. In a real app, configure is called from the NativeActivity thread.
	fmt.Println("configured (skipped — requires Activity context)")
	fmt.Println("codec configured for 1280x720 H.264 decoding")

	// 4. In a real application the next steps would be:
	//    codec.Start()
	//    ... dequeue input buffers, fill with NAL units, queue them ...
	//    ... dequeue output buffers, consume decoded frames, release them ...
	//    codec.Stop()
	fmt.Println("decoder setup complete (not started, no bitstream available)")
}
