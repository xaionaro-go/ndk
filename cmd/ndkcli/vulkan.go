package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var vulkanCmd = &cobra.Command{
	Use:   "vulkan",
	Short: "Vulkan NDK operations",
}

var vulkanInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Vulkan information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Vulkan API surface:")
		fmt.Println()
		fmt.Println("Package has no standalone-usable API.")
		fmt.Println("The vulkan package provides Go bindings for Vulkan on Android but currently")
		fmt.Println("exports no types or functions in the idiomatic layer.")
		fmt.Println("Use a dedicated Vulkan Go library (e.g., vulkan-go) for Vulkan rendering,")
		fmt.Println("and this package for Android-specific Vulkan integration once bindings are generated.")
		return nil
	},
}

func init() {
	vulkanCmd.AddCommand(vulkanInfoCmd)
	rootCmd.AddCommand(vulkanCmd)
}
