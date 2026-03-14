package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync primitives NDK operations",
}

var syncInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show sync information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Sync API surface:")
		fmt.Println()
		fmt.Println("Package has no standalone-usable API.")
		fmt.Println("The sync package provides Go bindings for Android sync primitives (sync fence FDs)")
		fmt.Println("but currently exports no types or functions in the idiomatic layer.")
		return nil
	},
}

func init() {
	syncCmd.AddCommand(syncInfoCmd)
	rootCmd.AddCommand(syncCmd)
}
