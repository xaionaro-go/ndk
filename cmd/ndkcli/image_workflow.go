package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/image"
)

var imageDecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode an image file and print dimensions, stride, and format info",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		filePath, _ := cmd.Flags().GetString("file")
		targetWidth, _ := cmd.Flags().GetInt32("width")
		targetHeight, _ := cmd.Flags().GetInt32("height")

		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		decoder, err := image.NewDecoderFromFd(int32(f.Fd()))
		if err != nil {
			return fmt.Errorf("creating decoder: %w", err)
		}
		defer decoder.Close()

		header := decoder.HeaderInfo()
		width := header.Width()
		height := header.Height()
		mimeType := header.MimeType()

		fmt.Printf("file:          %s\n", filePath)
		fmt.Printf("mime type:     %s\n", mimeType)
		fmt.Printf("dimensions:    %d x %d\n", width, height)

		if targetWidth > 0 && targetHeight > 0 {
			if err := decoder.SetTargetSize(targetWidth, targetHeight); err != nil {
				return fmt.Errorf("setting target size: %w", err)
			}
			fmt.Printf("target size:   %d x %d\n", targetWidth, targetHeight)
		}

		stride := decoder.MinimumStride()
		fmt.Printf("min stride:    %d bytes\n", stride)

		effectiveHeight := height
		if targetHeight > 0 {
			effectiveHeight = targetHeight
		}
		bufSize := stride * uint64(effectiveHeight)
		buf := make([]byte, bufSize)

		if err := decoder.Decode(unsafe.Pointer(&buf[0]), stride, bufSize); err != nil {
			return fmt.Errorf("decoding image: %w", err)
		}
		fmt.Printf("decoded:       %d bytes\n", bufSize)

		return nil
	},
}

func init() {
	imageDecodeCmd.Flags().String("file", "", "path to image file")
	_ = imageDecodeCmd.MarkFlagRequired("file")
	imageDecodeCmd.Flags().Int32("width", 0, "target decode width (0 = original)")
	imageDecodeCmd.Flags().Int32("height", 0, "target decode height (0 = original)")

	imageCmd.AddCommand(imageDecodeCmd)
}
