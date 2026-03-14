package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var hintCmd = &cobra.Command{
	Use:   "hint",
	Short: "Performance hint NDK operations",
}

var hintInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show performance hint information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Performance Hint API surface:")
		fmt.Println("  - Manager: wraps APerformanceHintManager (CreateSession, PreferredUpdateRateNanos)")
		fmt.Println("  - Session: wraps APerformanceHintSession (ReportActualWorkDuration, UpdateTargetWorkDuration, Close)")
		fmt.Println()
		fmt.Println("Note: The hint Manager requires a pointer obtained from Java via")
		fmt.Println("      APerformanceHint_getManager(). Use NewManagerFromPointer with that handle.")
		return nil
	},
}

func init() {
	hintCmd.AddCommand(hintInfoCmd)
	rootCmd.AddCommand(hintCmd)
}
