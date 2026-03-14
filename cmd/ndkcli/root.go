package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ndkcli",
	Short: "CLI tool for querying and interacting with the Android NDK surface",
	Long: `ndkcli provides access to all Android NDK modules from the command line.

Subcommands are organized by NDK module (camera, sensor, audio, etc.).
Run on an Android device or emulator.`,
}
