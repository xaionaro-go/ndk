package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Input NDK operations",
}

var inputInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show input information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Input API surface:")
		fmt.Println("  - Queue:      wraps AInputQueue (HasEvents, FinishEvent, DetachLooper)")
		fmt.Println("  - Event:      wraps AInputEvent (Type, Source, KeyAction, KeyCode, RepeatCount,")
		fmt.Println("                MotionAction, PointerCount, Pressure, X, Y)")
		fmt.Println("  - EventType:  Key, Motion, Focus, Capture, Drag, TouchMode")
		fmt.Println("  - KeyAction:  key event action constants")
		fmt.Println("  - MotionAction: motion event action constants")
		fmt.Println("  - Source:     input source constants")
		fmt.Println()
		fmt.Println("Note: Input queue requires a NativeActivity context (OnInputQueueCreated callback).")
		fmt.Println("      Use NewQueueFromPointer with the AInputQueue handle.")
		return nil
	},
}

func init() {
	inputCmd.AddCommand(inputInfoCmd)
	rootCmd.AddCommand(inputCmd)
}
