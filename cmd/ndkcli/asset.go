package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "Asset NDK operations",
}

var assetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show asset information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Asset API surface:")
		fmt.Println("  - Manager: wraps AAssetManager (Open, OpenDir)")
		fmt.Println("  - Asset:   wraps AAsset (Read, Seek, Buffer, Length, RemainingLength, Close)")
		fmt.Println("  - Dir:     wraps AAssetDir (NextFileName, Rewind, Close)")
		fmt.Println()
		fmt.Println("Note: Asset access requires a NativeActivity context.")
		fmt.Println("      Obtain the AAssetManager pointer from ANativeActivity and use NewManagerFromPointer.")
		return nil
	},
}

func init() {
	assetCmd.AddCommand(assetInfoCmd)
	rootCmd.AddCommand(assetCmd)
}
