package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var hwbufCmd = &cobra.Command{
	Use:   "hwbuf",
	Short: "Hardware buffer NDK operations",
}

var hwbufInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show hardware buffer information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hardware Buffer API surface:")
		fmt.Println("  - Buffer: wraps AHardwareBuffer (Acquire, Unlock, Close)")
		fmt.Println("  - Format: pixel format constants (R8g8b8a8Unorm, Blob, D32Float, YcbcrP010, ...)")
		fmt.Println("  - Usage:  buffer usage flags (CpuReadOften, GpuSampledImage, GpuFramebuffer, VideoEncode, ...)")
		fmt.Println()
		fmt.Println("Note: Buffers are created via AHardwareBuffer_allocate (not yet exposed in idiomatic API)")
		fmt.Println("      or received from other APIs. Use NewBufferFromPointer with an existing handle.")
		return nil
	},
}

func init() {
	hwbufCmd.AddCommand(hwbufInfoCmd)
	rootCmd.AddCommand(hwbufCmd)
}
