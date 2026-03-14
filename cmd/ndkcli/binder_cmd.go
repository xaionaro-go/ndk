package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var binderCmd = &cobra.Command{
	Use:   "binder",
	Short: "Binder NDK operations",
}

var binderInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show binder information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Binder API surface:")
		fmt.Println("  - Binder:        wraps AIBinder (Close/decStrong)")
		fmt.Println("  - Class:         wraps AIBinder_Class")
		fmt.Println("  - Parcel:        wraps AParcel (Close/delete)")
		fmt.Println("  - Status:        wraps AStatus (Close/delete)")
		fmt.Println("  - ExceptionCode: None, Security, BadParcelable, IllegalArgument, NullPointer, ...")
		fmt.Println("  - Flags:         Oneway")
		fmt.Println("  - TransactionCodeT: transaction code type")
		fmt.Println()
		fmt.Println("Note: Binder handles are obtained from the Android binder driver via")
		fmt.Println("      AIBinder_new or from Java via AIBinder_fromJavaBinder.")
		fmt.Println("      Use NewBinderFromPointer with an existing handle.")
		return nil
	},
}

func init() {
	binderCmd.AddCommand(binderInfoCmd)
	rootCmd.AddCommand(binderCmd)
}
