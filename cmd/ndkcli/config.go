package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration NDK operations",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show device configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.NewConfig()
		defer c.Close()

		fmt.Println("Device configuration (default AConfiguration):")
		fmt.Printf("  Orientation:    %d\n", c.Orientation())
		fmt.Printf("  ScreenWidthDp:  %d\n", c.ScreenWidthDp())
		fmt.Printf("  ScreenHeightDp: %d\n", c.ScreenHeightDp())
		fmt.Printf("  Density:        %d\n", c.Density())
		fmt.Printf("  ScreenSize:     %d\n", c.ScreenSize())
		fmt.Printf("  SdkVersion:     %d\n", c.SdkVersion())
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
