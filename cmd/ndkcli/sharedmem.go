package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sharedmemCmd = &cobra.Command{
	Use:   "sharedmem",
	Short: "Shared memory NDK operations",
}

var sharedmemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show shared memory information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Shared Memory API surface:")
		fmt.Println()
		fmt.Println("Package has no standalone-usable API.")
		fmt.Println("The sharedmem package provides Go bindings for Android shared memory (ASharedMemory)")
		fmt.Println("but currently exports no types or functions in the idiomatic layer.")
		return nil
	},
}

func init() {
	sharedmemCmd.AddCommand(sharedmemInfoCmd)
	rootCmd.AddCommand(sharedmemCmd)
}
