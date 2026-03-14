package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var surfacetextureCmd = &cobra.Command{
	Use:   "surfacetexture",
	Short: "Surface texture NDK operations",
}

var surfacetextureInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show surface texture information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("SurfaceTexture API surface:")
		fmt.Println("  - SurfaceTexture: wraps ASurfaceTexture (AcquireWindow, AttachToGLContext,")
		fmt.Println("                    DetachFromGLContext, Timestamp, TransformMatrix, UpdateTexImage, Close)")
		fmt.Println("  - NativeWindow:   wraps ANativeWindow obtained from SurfaceTexture")
		fmt.Println()
		fmt.Println("Note: SurfaceTexture requires a handle obtained via JNI (ASurfaceTexture_fromSurfaceTexture).")
		fmt.Println("      Use NewSurfaceTextureFromPointer with an existing handle.")
		return nil
	},
}

func init() {
	surfacetextureCmd.AddCommand(surfacetextureInfoCmd)
	rootCmd.AddCommand(surfacetextureCmd)
}
