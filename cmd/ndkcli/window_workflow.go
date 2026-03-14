package main

import (
	"fmt"
	"unsafe"

	"github.com/spf13/cobra"
	mediacapi "github.com/xaionaro-go/ndk/capi/media"
	"github.com/xaionaro-go/ndk/window"
)

var windowQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query window properties via an ImageReader-backed ANativeWindow",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		width, _ := cmd.Flags().GetInt32("width")
		height, _ := cmd.Flags().GetInt32("height")
		format, _ := cmd.Flags().GetInt32("format")

		// Create an ImageReader via capi to get an ANativeWindow handle.
		var readerPtr *mediacapi.AImageReader
		status := mediacapi.AImageReader_new(width, height, format, 2, &readerPtr)
		if status < 0 {
			return fmt.Errorf("creating image reader: media error %d", status)
		}
		defer mediacapi.AImageReader_delete(readerPtr)

		var nativeWindow *mediacapi.ANativeWindow
		status = mediacapi.AImageReader_getWindow(readerPtr, &nativeWindow)
		if status < 0 {
			return fmt.Errorf("getting window from image reader: media error %d", status)
		}

		win := window.NewWindowFromPointer(unsafe.Pointer(nativeWindow))

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
	windowQueryCmd.Flags().Int32("format", mediacapi.AIMAGE_FORMAT_RGBA_8888, "image format (1=RGBA_8888)")

	windowCmd.AddCommand(windowQueryCmd)
}
