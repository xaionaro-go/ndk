package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/trace"
)

var traceSectionCmd = &cobra.Command{
	Use:   "section",
	Short: "Begin a named trace section (ended on process exit)",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		name, _ := cmd.Flags().GetString("name")

		trace.BeginSection(name)
		fmt.Println("trace section '" + name + "' started")
		return nil
	},
}

func init() {
	traceSectionCmd.Flags().String("name", "ndkcli", "trace section name")

	traceCmd.AddCommand(traceSectionCmd)
}
