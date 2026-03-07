// Media format creation and property access example.
//
// Demonstrates how to create an AMediaFormat, populate it with MIME type,
// resolution, and bitrate using the chaining API, and then query values
// back with GetInt32. The idiomatic API accepts Go strings for key names.
//
// This program must run on an Android device with NDK media support.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xaionaro-go/ndk/media"
)

func main() {
	// NDK format keys as Go strings.
	keyWidth := "width"
	keyHeight := "height"
	keyBitrate := "bitrate"
	keyFrameRate := "frame-rate"

	// Create a format and configure it with chaining.
	format := media.NewFormat()
	defer func() {
		if err := format.Close(); err != nil {
			log.Printf("close format: %v", err)
		}
	}()

	format.
		SetString("mime", "video/avc").
		SetInt32(keyWidth, 1920).
		SetInt32(keyHeight, 1080).
		SetInt32(keyBitrate, 5_000_000). // 5 Mbps
		SetInt32(keyFrameRate, 30)

	// Read properties back to verify they were stored.
	queries := []struct {
		label string
		key   string
	}{
		{"Width", keyWidth},
		{"Height", keyHeight},
		{"Bitrate", keyBitrate},
		{"Frame rate", keyFrameRate},
	}

	fmt.Fprintf(os.Stdout, "Format properties:\n")
	for _, q := range queries {
		var val int32
		if ok := format.GetInt32(q.key, &val); !ok {
			log.Fatalf("GetInt32(%s): key not found", q.label)
		}
		fmt.Fprintf(os.Stdout, "  %-12s: %d\n", q.label, val)
	}

	fmt.Println("format created and verified successfully")
}
