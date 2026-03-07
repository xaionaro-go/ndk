// Media extractor track inspection example.
//
// Demonstrates how to create an AMediaExtractor, attach it to a file
// descriptor, query the number of tracks, select a track, and read the
// current sample timestamp. This is the typical first step when building
// a media player or transcoder: open the container, discover tracks,
// and choose which ones to process.
//
// Usage: extractor-info <media-file>
//
// This program must run on an Android device with NDK media support.
package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/xaionaro-go/ndk/media"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <media-file>\n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]

	// 1. Open the media file to obtain a file descriptor.
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer syscall.Close(fd)

	// Determine file size for SetDataSourceFd.
	var stat syscall.Stat_t
	if err := syscall.Fstat(fd, &stat); err != nil {
		log.Fatalf("fstat: %v", err)
	}

	// 2. Create an extractor and point it at the file.
	extractor := media.NewExtractor()
	defer func() {
		if err := extractor.Close(); err != nil {
			log.Printf("close extractor: %v", err)
		}
	}()

	if err := extractor.SetDataSourceFd(int32(fd), 0, media.Off64_t(stat.Size)); err != nil {
		log.Fatalf("set data source: %v", err)
	}
	fmt.Printf("opened: %s\n", path)

	// 3. Query the number of tracks in the container.
	trackCount := extractor.TrackCount()
	fmt.Printf("track count: %d\n", trackCount)

	if trackCount == 0 {
		fmt.Println("no tracks found")
		return
	}

	// 4. Select the first track and read its initial sample time.
	if err := extractor.SelectTrack(0); err != nil {
		log.Fatalf("select track 0: %v", err)
	}
	fmt.Printf("selected track 0\n")

	sampleTime := extractor.SampleTime()
	fmt.Printf("first sample time: %d us\n", sampleTime)

	// In a real application the next steps would be:
	//   - Read sample data with ReadSampleData into a buffer
	//   - Advance to the next sample with Advance()
	//   - Feed samples into a Codec for decoding
	fmt.Println("extractor setup complete")
}
