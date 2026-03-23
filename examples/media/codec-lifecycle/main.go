// MediaCodec full lifecycle example.
//
// Demonstrates the complete AMediaCodec state machine for an H.264 video
// decoder: create, configure, start, the buffer processing loop pattern,
// flush, stop, and close. The codec transitions through these states:
//
//	Created -> Configured -> Started -> (processing) -> Flushed -> Started -> Stopped -> Released
//
// The configure and start steps require an Activity context, so they are
// skipped here. The buffer processing loop is shown as commented-out code
// using the idiomatic media.DequeueInputBuffer, media.DequeueOutputBuffer,
// and media.BufferInfo APIs.
//
// This program must run on an Android device with NDK media support.
package main

import (
	"fmt"
	"log"

	"github.com/AndroidGoLab/ndk/media"
)

func main() {
	// State: (none) -> Created
	// Create an H.264 decoder.
	codec := media.NewDecoder("video/avc")
	if codec == nil {
		log.Fatal("NewDecoder returned nil: H.264 decoder not available")
	}
	defer func() {
		if err := codec.Close(); err != nil {
			log.Printf("close codec: %v", err)
		}
	}()
	fmt.Println("state: Created")

	// Build the format describing the compressed input stream.
	format := media.NewFormat()
	defer func() {
		if err := format.Close(); err != nil {
			log.Printf("close format: %v", err)
		}
	}()

	format.
		SetString("mime", "video/avc").
		SetInt32("width", 1920).
		SetInt32("height", 1080)

	// Verify format properties were stored correctly.
	var width, height int32
	if !format.GetInt32("width", &width) || !format.GetInt32("height", &height) {
		log.Fatal("failed to read back format properties")
	}
	fmt.Printf("format: video/avc %dx%d\n", width, height)

	// State: Created -> Configured
	// Configure the decoder without a surface (buffer-mode) and without DRM.
	// Note: AMediaCodec_configure may crash with SIGSEGV on some devices
	// when called outside a proper Activity context (no display surface).
	// In a real app, configure is called from the NativeActivity thread.
	fmt.Println("state: Configured (skipped — requires Activity context)")

	// State: Configured -> Started
	// codec.Start() also requires prior Configure, so it is skipped here.
	// In a real NativeActivity app:
	//   if err := codec.Configure(format, nil, nil, 0); err != nil { ... }
	//   if err := codec.Start(); err != nil { ... }
	fmt.Println("state: Started (skipped — requires Activity context)")

	// -- Buffer processing loop (documented pattern) --
	//
	// In a real application the loop uses the idiomatic media package:
	//
	//   const timeoutUs int64 = 10000 // 10ms
	//
	//   for {
	//       // 1. Dequeue an input buffer.
	//       idx := codec.DequeueInputBuffer(timeoutUs)
	//       if idx < 0 {
	//           continue // no buffer available yet
	//       }
	//
	//       // 2. Fill the buffer with compressed data (NAL units for H.264).
	//       var bufSize uint64
	//       buf := codec.GetInputBuffer(uint64(idx), &bufSize)
	//       data := unsafe.Slice(buf, bufSize)
	//       n := copy(data, nalUnit)
	//
	//       // 3. Queue the filled buffer for decoding.
	//       codec.QueueInputBuffer(uint64(idx), 0, uint64(n), presentationTimeUs, 0)
	//       // For end-of-stream, pass uint32(media.AMEDIACODEC_BUFFER_FLAG_END_OF_STREAM) as flags.
	//
	//       // 4. Dequeue an output buffer using media.BufferInfo.
	//       var info media.BufferInfo
	//       outIdx := codec.DequeueOutputBuffer(&info, timeoutUs)
	//       if outIdx == int64(media.AMEDIACODEC_INFO_OUTPUT_FORMAT_CHANGED) {
	//           continue // output format changed, query new format
	//       }
	//       if outIdx < 0 {
	//           continue // try again
	//       }
	//
	//       // 5. Consume the decoded frame from the output buffer.
	//       //    info.Offset, info.Size, info.PresentationTimeUs, info.Flags
	//       //    describe the output buffer contents.
	//       var outSize uint64
	//       outBuf := codec.GetOutputBuffer(uint64(outIdx), &outSize)
	//       outData := unsafe.Slice(outBuf, outSize)
	//       processFrame(outData)
	//
	//       // 6. Release the output buffer back to the codec.
	//       codec.ReleaseOutputBuffer(uint64(outIdx), false)
	//   }
	fmt.Println("state: Started (buffer loop would run here)")

	// The following lifecycle steps require a properly configured and
	// started codec, which needs an Activity context. Documented here:
	//
	//   codec.Flush()  // Started -> Flushed (discard pending buffers)
	//   // After flush, codec stays in Started state, ready for new input.
	//   codec.Stop()   // Started -> Stopped
	//   codec.Close()  // Stopped -> Released (free all resources)
	fmt.Println("state: Released (via deferred Close)")

	fmt.Println("codec lifecycle complete")
}
