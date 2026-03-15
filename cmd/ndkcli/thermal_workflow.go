package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/thermal"
)

var thermalMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor thermal status at a regular interval",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		interval, _ := cmd.Flags().GetDuration("interval")
		duration, _ := cmd.Flags().GetDuration("duration")

		mgr := thermal.NewManager()
		defer mgr.Close()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		defer signal.Stop(sigCh)

		var deadline <-chan time.Time
		if duration > 0 {
			timer := time.NewTimer(duration)
			defer timer.Stop()
			deadline = timer.C
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		tick := 0
		for {
			status := mgr.CurrentStatus()
			fmt.Printf("[%d] thermal status: %s (%d)\n", tick, status, int32(status))
			tick++

			select {
			case <-sigCh:
				fmt.Println("interrupted")
				return nil
			case <-deadline:
				fmt.Printf("monitoring complete (%d samples)\n", tick)
				return nil
			case <-ticker.C:
			}
		}
	},
}

func init() {
	thermalMonitorCmd.Flags().Duration("interval", time.Second, "polling interval")
	thermalMonitorCmd.Flags().Duration("duration", 0, "total monitoring duration (0 = until Ctrl+C)")

	thermalCmd.AddCommand(thermalMonitorCmd)
}
