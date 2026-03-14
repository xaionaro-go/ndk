package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var surfacecontrolCmd = &cobra.Command{
	Use:   "surfacecontrol",
	Short: "Surface control NDK operations",
}

var surfacecontrolInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show surface control information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("SurfaceControl API surface:")
		fmt.Println("  - SurfaceControl: wraps ASurfaceControl (Acquire, CreateChild, Close)")
		fmt.Println("  - Transaction:    wraps ASurfaceTransaction (Apply, SetBufferAlpha, SetCrop,")
		fmt.Println("                    SetDamageRegion, SetPosition, SetScale, SetZOrder, Close)")
		fmt.Println("  - ARect:          rectangle type for crop/damage regions")
		fmt.Println("  - Transparency:   Transparent, Translucent, Opaque")
		fmt.Println("  - Visibility:     Hide, Show")
		fmt.Println()
		fmt.Println("Note: SurfaceControl handles are created from an existing ANativeWindow surface")
		fmt.Println("      via ASurfaceControl_createFromWindow. Use NewSurfaceControlFromPointer")
		fmt.Println("      with a handle obtained from the windowing system.")
		return nil
	},
}

func init() {
	surfacecontrolCmd.AddCommand(surfacecontrolInfoCmd)
	rootCmd.AddCommand(surfacecontrolCmd)
}
