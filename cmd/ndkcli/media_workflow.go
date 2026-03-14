package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/media"
)

var mediaCodecsCmd = &cobra.Command{
	Use:   "codecs",
	Short: "Probe available media codecs (encoder/decoder) for common MIME types",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
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

		fmt.Printf("%-25s %-10s %-10s\n", "MIME", "Encoder", "Decoder")
		fmt.Printf("%-25s %-10s %-10s\n", "----", "-------", "-------")

		for _, mime := range mimeTypes {
			enc := media.NewEncoder(mime)
			hasEncoder := enc.Pointer() != nil
			if hasEncoder {
				enc.Close()
			}

			dec := media.NewDecoder(mime)
			hasDecoder := dec.Pointer() != nil
			if hasDecoder {
				dec.Close()
			}

			fmt.Printf("%-25s %-10v %-10v\n", mime, hasEncoder, hasDecoder)
		}

		return nil
	},
}

var mediaProbeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Probe a media file for track information",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		filePath, _ := cmd.Flags().GetString("file")

		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			return fmt.Errorf("stat file: %w", err)
		}

		extractor := media.NewExtractor()
		defer extractor.Close()

		if err := extractor.SetDataSourceFd(int32(f.Fd()), 0, media.Off64_t(info.Size())); err != nil {
			return fmt.Errorf("setting data source: %w", err)
		}

		trackCount := extractor.TrackCount()
		fmt.Printf("file:   %s\n", filePath)
		fmt.Printf("tracks: %d\n", trackCount)

		for i := uint64(0); i < trackCount; i++ {
			if err := extractor.SelectTrack(i); err != nil {
				fmt.Printf("  track %d: error selecting: %v\n", i, err)
				continue
			}
			sampleTime := extractor.SampleTime()
			fmt.Printf("  track %d: sample_time=%d us\n", i, sampleTime)
		}

		return nil
	},
}

func init() {
	mediaProbeCmd.Flags().String("file", "", "path to media file")
	_ = mediaProbeCmd.MarkFlagRequired("file")

	mediaCmd.AddCommand(mediaCodecsCmd)
	mediaCmd.AddCommand(mediaProbeCmd)
}
