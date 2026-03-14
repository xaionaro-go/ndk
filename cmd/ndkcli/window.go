package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var windowCmd = &cobra.Command{
	Use:   "window",
	Short: "Window NDK operations",
}

var windowInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show window information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Window API surface:")
		fmt.Println("  - Window: wraps ANativeWindow (Width, Height, Format, SetBuffersGeometry, UnlockAndPost)")
		fmt.Println("  - Format constants: Rgba8888, Rgbx8888, Rgb565")
		fmt.Println()
		fmt.Println("Note: Window requires a pointer from NativeActivity (OnNativeWindowCreated callback).")
		fmt.Println("      Use NewWindowFromPointer with the ANativeWindow handle.")
		return nil
	},
}

func init() {
	windowCmd.AddCommand(windowInfoCmd)
	rootCmd.AddCommand(windowCmd)
}
