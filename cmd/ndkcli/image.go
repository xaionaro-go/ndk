package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image decoder NDK operations",
}

var imageInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show image decoder information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Image Decoder API surface:")
		fmt.Println("  - Decoder:    wraps AImageDecoder (Decode, MinimumStride, SetTargetSize, Close)")
		fmt.Println("  - HeaderInfo: wraps AImageDecoderHeaderInfo (image metadata)")
		fmt.Println()
		fmt.Println("Note: Decoders are created from an AAsset or file descriptor via")
		fmt.Println("      AImageDecoder_createFromAAsset / AImageDecoder_createFromFd.")
		fmt.Println("      Use NewDecoderFromPointer with a handle obtained from these functions.")
		return nil
	},
}

func init() {
	imageCmd.AddCommand(imageInfoCmd)
	rootCmd.AddCommand(imageCmd)
}
