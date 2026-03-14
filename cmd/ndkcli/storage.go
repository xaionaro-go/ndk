package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Storage NDK operations",
}

var storageInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show storage information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Storage API surface:")
		fmt.Println("  - Manager: wraps AStorageManager (NewManager, MountObb, UnmountObb, IsObbMounted, MountedObbPath)")
		fmt.Println()
		fmt.Println("Note: OBB operations require a running NativeActivity to receive callbacks.")
		fmt.Println("      NewManager can be called standalone, but mount/unmount need an active Android environment.")
		return nil
	},
}

func init() {
	storageCmd.AddCommand(storageInfoCmd)
	rootCmd.AddCommand(storageCmd)
}
