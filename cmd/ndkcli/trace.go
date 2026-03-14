package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/trace"
)

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Systrace tracing operations",
}

var traceEnabledCmd = &cobra.Command{
	Use:   "enabled",
	Short: "Check whether tracing is enabled",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Tracing enabled: %v\n", trace.IsEnabled())
		return nil
	},
}

func init() {
	traceCmd.AddCommand(traceEnabledCmd)
	rootCmd.AddCommand(traceCmd)
}
