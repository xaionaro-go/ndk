package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/looper"
)

var looperTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Prepare a looper, wake it, and poll once to verify looper functionality",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		// Looper must be used on a locked OS thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		lp := looper.Prepare(1) // 1 = ALOOPER_PREPARE_ALLOW_NON_CALLBACKS
		if lp.Pointer() == nil {
			return fmt.Errorf("looper.Prepare returned nil")
		}
		defer lp.Close()

		fmt.Println("looper prepared successfully")

		lp.Wake()
		fmt.Println("looper.Wake called")

		result := looper.PollOnce(0, nil, nil, nil)
		fmt.Printf("looper.PollOnce(0) result: %d\n", result)

		return nil
	},
}

func init() {
	looperCmd.AddCommand(looperTestCmd)
}
