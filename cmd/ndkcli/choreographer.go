package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/choreographer"
)

var choreographerCmd = &cobra.Command{
	Use:   "choreographer",
	Short: "Choreographer NDK operations",
}

var choreographerInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show choreographer information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Choreographer API surface:")
		fmt.Println("  - Choreographer: wraps AChoreographer")
		fmt.Println("  - GetInstance:   obtain the choreographer for the current thread's looper")
		fmt.Println("  - AVsyncId:     vsync identifier type")
		fmt.Println()

		ch := choreographer.GetInstance()
		if ch.Pointer() == nil {
			fmt.Println("GetInstance: returned nil (no looper attached to current thread)")
		} else {
			fmt.Println("GetInstance: choreographer obtained successfully")
		}
		return nil
	},
}

func init() {
	choreographerCmd.AddCommand(choreographerInfoCmd)
	rootCmd.AddCommand(choreographerCmd)
}
