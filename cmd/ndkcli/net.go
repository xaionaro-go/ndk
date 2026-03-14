package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Network NDK operations",
}

var netInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show network information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Network (multinetwork) API surface:")
		fmt.Println("  - net_handle_t: opaque network handle type")
		fmt.Println()
		fmt.Println("Package has no standalone-usable API.")
		fmt.Println("The net package exposes the net_handle_t type alias for use with")
		fmt.Println("ConnectivityManager network handles obtained via Java/JNI.")
		return nil
	},
}

func init() {
	netCmd.AddCommand(netInfoCmd)
	rootCmd.AddCommand(netCmd)
}
