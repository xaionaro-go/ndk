package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/storage"
)

var storageObbCmd = &cobra.Command{
	Use:   "obb",
	Short: "Check OBB mount status for a given file path",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		filePath, _ := cmd.Flags().GetString("file")

		mgr := storage.NewManager()
		defer mgr.Close()

		mounted := mgr.IsObbMounted(filePath)
		mountedPath := mgr.MountedObbPath(filePath)

		fmt.Printf("file:         %s\n", filePath)
		fmt.Printf("is mounted:   %v (raw: %d)\n", mounted != 0, mounted)
		fmt.Printf("mounted path: %s\n", mountedPath)

		return nil
	},
}

func init() {
	storageObbCmd.Flags().String("file", "", "path to OBB file")
	_ = storageObbCmd.MarkFlagRequired("file")

	storageCmd.AddCommand(storageObbCmd)
}
