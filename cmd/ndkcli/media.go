package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/media"
)

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "MediaCodec operations",
}

var mediaCodecsCmd = &cobra.Command{
	Use:   "codecs",
	Short: "List available media codecs",
	RunE: func(cmd *cobra.Command, args []string) error {
		mimeTypes := []string{
			"video/avc",
			"video/hevc",
			"video/x-vnd.on2.vp8",
			"video/x-vnd.on2.vp9",
			"video/av01",
			"audio/mp4a-latm",
			"audio/mpeg",
			"audio/opus",
			"audio/vorbis",
		}

		for _, mime := range mimeTypes {
			enc := media.NewEncoder(mime)
			encAvail := enc.Pointer() != nil
			if encAvail {
				_ = enc.Close()
			}

			dec := media.NewDecoder(mime)
			decAvail := dec.Pointer() != nil
			if decAvail {
				_ = dec.Close()
			}

			fmt.Printf("%-25s encoder=%-5v decoder=%v\n", mime, encAvail, decAvail)
		}

		return nil
	},
}

func init() {
	mediaCmd.AddCommand(mediaCodecsCmd)
	rootCmd.AddCommand(mediaCmd)
}
