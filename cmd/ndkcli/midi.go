package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var midiCmd = &cobra.Command{
	Use:   "midi",
	Short: "MIDI NDK operations",
}

var midiInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show MIDI information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("MIDI API surface:")
		fmt.Println("  - Device:     wraps AMidiDevice (NumInputPorts, NumOutputPorts, OpenInputPort, OpenOutputPort)")
		fmt.Println("  - InputPort:  wraps AMidiInputPort (Send, Close)")
		fmt.Println("  - OutputPort: wraps AMidiOutputPort (Close)")
		fmt.Println()
		fmt.Println("Note: MIDI device enumeration requires a Java context (AMidiDevice_fromJava).")
		fmt.Println("      Use NewDeviceFromPointer with a handle obtained via JNI.")
		return nil
	},
}

func init() {
	midiCmd.AddCommand(midiInfoCmd)
	rootCmd.AddCommand(midiCmd)
}
