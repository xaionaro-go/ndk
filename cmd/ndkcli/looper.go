package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/looper"
)

var looperCmd = &cobra.Command{
	Use:   "looper",
	Short: "Looper NDK operations",
}

var looperInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show looper information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Looper API surface:")
		fmt.Println("  - Looper:   wraps ALooper (Acquire, Wake, RemoveFd, Close)")
		fmt.Println("  - Prepare:  create a looper for the current thread")
		fmt.Println("  - PollOnce: poll for events with timeout")
		fmt.Println("  - Run:      lock OS thread, prepare looper, call callback, then clean up")
		fmt.Println("  - Event constants: EventInput, EventOutput, EventError, EventHangup, EventInvalid")
		fmt.Println()

		looper.Run(func(lp *looper.Looper) {
			fmt.Println("Looper test: prepared looper for current thread successfully.")
		})
		return nil
	},
}

func init() {
	looperCmd.AddCommand(looperInfoCmd)
	rootCmd.AddCommand(looperCmd)
}
