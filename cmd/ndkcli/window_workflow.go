package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/media"
	"github.com/AndroidGoLab/ndk/window"
)

// imageFormatRGBA8888 is AIMAGE_FORMAT_RGBA_8888 from the NDK header.
const imageFormatRGBA8888 = 1

var windowQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query window properties via an ImageReader-backed ANativeWindow",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		width, _ := cmd.Flags().GetInt32("width")
		height, _ := cmd.Flags().GetInt32("height")
		format, _ := cmd.Flags().GetInt32("format")

		// Create an ImageReader to get an ANativeWindow handle.
		reader, err := media.NewImageReader(width, height, format, 2)
		if err != nil {
			return fmt.Errorf("creating image reader: %w", err)
		}
		defer reader.Close()

		mediaWin, err := reader.Window()
		if err != nil {
			return fmt.Errorf("getting window from image reader: %w", err)
		}

		// Convert media.Window to window.Window via unsafe.Pointer for the
		// window package's query methods.
		win := window.NewWindowFromPointer(mediaWin.Pointer())

		fmt.Printf("window properties:\n")
		fmt.Printf("  width:  %d\n", win.Width())
		fmt.Printf("  height: %d\n", win.Height())
		fmt.Printf("  format: %d (%s)\n", win.Format(), window.Format(win.Format()))

		return nil
	},
}

func init() {
	windowQueryCmd.Flags().Int32("width", 640, "image reader width")
	windowQueryCmd.Flags().Int32("height", 480, "image reader height")
	windowQueryCmd.Flags().Int32("format", imageFormatRGBA8888, "image format (1=RGBA_8888)")

	windowCmd.AddCommand(windowQueryCmd)
}
