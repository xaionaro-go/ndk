package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Activity NDK operations",
}

var activityInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show activity information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Activity API surface:")
		fmt.Println("  - Activity:           wraps ANativeActivity (Finish, ShowSoftInput, HideSoftInput,")
		fmt.Println("                        SetWindowFlags, SetWindowFormat)")
		fmt.Println("  - LifecycleCallbacks: Go callback struct for lifecycle events (OnCreate, OnStart,")
		fmt.Println("                        OnResume, OnPause, OnStop, OnDestroy, OnWindowFocusChanged,")
		fmt.Println("                        OnNativeWindowCreated/Resized/RedrawNeeded/Destroyed,")
		fmt.Println("                        OnInputQueueCreated/Destroyed, OnConfigurationChanged, OnLowMemory)")
		fmt.Println("  - SetLifecycleCallbacks: register lifecycle callbacks")
		fmt.Println("  - HideSoftInputFlags:  ImplicitOnly, NotAlways")
		fmt.Println("  - ShowSoftInputFlags:  flag constants for ShowSoftInput")
		fmt.Println()
		fmt.Println("Note: Activity requires a NativeActivity context. The ANativeActivity pointer is")
		fmt.Println("      provided by the Android runtime at application startup.")
		return nil
	},
}

func init() {
	activityCmd.AddCommand(activityInfoCmd)
	rootCmd.AddCommand(activityCmd)
}
