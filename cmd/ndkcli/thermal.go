package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/thermal"
)

var thermalCmd = &cobra.Command{
	Use:   "thermal",
	Short: "Query thermal status via the NDK Thermal API",
}

var thermalStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the current thermal status",
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) (_err error) {
		mgr := thermal.NewManager()
		defer mgr.Close()

		status := mgr.CurrentStatus()
		fmt.Printf("Thermal Status: %s\n", status)
		return nil
	},
}

func init() {
	thermalCmd.AddCommand(thermalStatusCmd)
	rootCmd.AddCommand(thermalCmd)
}
