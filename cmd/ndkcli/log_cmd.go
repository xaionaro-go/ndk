package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/log"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Android logging operations",
}

var logWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a message to the Android log",
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}
		priority, err := cmd.Flags().GetInt32("priority")
		if err != nil {
			return err
		}

		ret := log.Write(priority, tag, message)
		fmt.Printf("log.Write returned: %d\n", ret)
		return nil
	},
}

func init() {
	logWriteCmd.Flags().String("tag", "", "Log tag")
	logWriteCmd.Flags().String("message", "", "Log message")
	logWriteCmd.Flags().Int32("priority", int32(log.Info), "Log priority (2=Verbose, 3=Debug, 4=Info, 5=Warn, 6=Error)")

	logCmd.AddCommand(logWriteCmd)
	rootCmd.AddCommand(logCmd)
}
