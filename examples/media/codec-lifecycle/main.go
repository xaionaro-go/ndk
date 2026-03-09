// MediaCodec full lifecycle example.
//
// Demonstrates the complete AMediaCodec state machine for an H.264 video
// decoder: create, configure, start, the buffer processing loop pattern,
// flush, stop, and close. The codec transitions through these states:
//
//	Created -> Configured -> Started -> (processing) -> Flushed -> Started -> Stopped -> Released
//
// The input/output buffer loop is shown as documented comments because
// DequeueInputBuffer and DequeueOutputBuffer are not yet wrapped in the
// idiomatic layer. All other lifecycle methods are exercised for real.
//
// This program must run on an Android device with NDK media support.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/media"
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
	noCrypto := &media.Crypto{}
	if err := codec.Configure(format, nil, noCrypto, 0); err != nil {
		log.Fatalf("configure: %v", err)
	}
	fmt.Println("state: Configured")

	// State: Configured -> Started
	if err := codec.Start(); err != nil {
		log.Fatalf("start: %v", err)
	}
	fmt.Println("state: Started")

	// -- Buffer processing loop (documented pattern) --
	//
	// In a real application the loop looks like this:
	//
	//   for {
	//       // 1. Dequeue an input buffer (not yet in idiomatic layer).
	//       //    idx := codec.DequeueInputBuffer(timeoutUs)
	//       //    if idx < 0 { continue }  // no buffer available yet
	//
	//       // 2. Fill the buffer with compressed data (NAL units for H.264).
	//       //    buf := codec.GetInputBuffer(idx)
	//       //    n := copy(buf, nalUnit)
	//
	//       // 3. Queue the filled buffer for decoding.
	//       //    codec.QueueInputBuffer(idx, 0, uint64(n), presentationTimeUs, 0)
	//       //    For end-of-stream, pass media.BufferFlagEndOfStream as flags.
	//
	//       // 4. Dequeue an output buffer (not yet in idiomatic layer).
	//       //    outIdx := codec.DequeueOutputBuffer(&info, timeoutUs)
	//       //    if outIdx < 0 { continue }  // format change or try again
	//
	//       // 5. Consume the decoded frame from the output buffer.
	//       //    outBuf := codec.GetOutputBuffer(outIdx)
	//       //    processFrame(outBuf)
	//
	//       // 6. Release the output buffer back to the codec.
	//       //    codec.ReleaseOutputBuffer(outIdx, false)
	//   }
	fmt.Println("state: Started (buffer loop would run here)")

	// State: Started -> Flushed
	// Flush discards all pending buffers, allowing the codec to be
	// re-used from a new position (e.g. after a seek).
	if err := codec.Flush(); err != nil {
		log.Fatalf("flush: %v", err)
	}
	fmt.Println("state: Flushed")

	// State: Flushed -> Started (codec is ready for a new buffer loop)
	// After flush the codec remains in the Started state and can accept
	// new input immediately. A second Start() call is not needed.
	fmt.Println("state: Started (ready for new input after flush)")

	// State: Started -> Stopped
	if err := codec.Stop(); err != nil {
		log.Fatalf("stop: %v", err)
	}
	fmt.Println("state: Stopped")

	// State: Stopped -> Released
	// Close() releases the codec and frees all resources.
	// The deferred Close() above will handle this, but we call it
	// explicitly here to show the final state transition.
	if err := codec.Close(); err != nil {
		log.Fatalf("close: %v", err)
	}
	fmt.Println("state: Released")

	fmt.Println("codec lifecycle complete")
}
