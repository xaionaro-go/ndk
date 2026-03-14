package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bitmapCmd = &cobra.Command{
	Use:   "bitmap",
	Short: "Bitmap NDK operations",
}

var bitmapInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show bitmap information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Bitmap API surface:")
		fmt.Println("  - GetInfo(env, jbitmap, info):         retrieve AndroidBitmapInfo for a Java Bitmap")
		fmt.Println("  - LockPixels(env, jbitmap, addrPtr):   lock pixel buffer for direct access")
		fmt.Println("  - UnlockPixels(env, jbitmap):          unlock pixel buffer")
		fmt.Println("  - AndroidBitmapInfo:                    width, height, stride, format, flags")
		fmt.Println("  - Format constants:                     Rgba8888, Rgb565, Rgba4444, A8, RgbaF16, Rgba1010102")
		fmt.Println()
		fmt.Println("Note: Bitmap operations require a JNI environment (*JNIEnv) and a Java Bitmap jobject.")
		fmt.Println("      Use within a NativeActivity or JNI callback.")
		return nil
	},
}

func init() {
	bitmapCmd.AddCommand(bitmapInfoCmd)
	rootCmd.AddCommand(bitmapCmd)
}
