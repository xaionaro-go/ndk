package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/sensor"
)

var sensorReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read sensor information and poll default sensor",
	Long: `Queries the sensor manager for the default sensor of the given type,
prints its static properties, and optionally polls for the specified duration.

Note: Real-time sensor data streaming requires looper and event queue
integration, which is beyond what a simple CLI poll can provide.
This command demonstrates the sensor query API.`,
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		sensorType, _ := cmd.Flags().GetInt32("type")
		duration, _ := cmd.Flags().GetDuration("duration")

		mgr := sensor.GetInstance()

		s := mgr.DefaultSensor(sensorType)
		if s.Pointer() == nil {
			return fmt.Errorf("no default sensor found for type %d", sensorType)
		}

		fmt.Printf("sensor info:\n")
		fmt.Printf("  name:       %s\n", s.Name())
		fmt.Printf("  vendor:     %s\n", s.Vendor())
		fmt.Printf("  type:       %d (%s)\n", s.Type(), sensor.Type(s.Type()))
		fmt.Printf("  resolution: %g\n", s.Resolution())
		fmt.Printf("  min delay:  %d us\n", s.MinDelay())

		fmt.Println("\nnote: real-time sensor data streaming requires looper + event queue integration")

		if duration > 0 {
			fmt.Printf("\npolling sensor info for %v...\n", duration)
			deadline := time.Now().Add(duration)
			tick := 0
			for time.Now().Before(deadline) {
				polled := mgr.DefaultSensor(sensorType)
				if polled.Pointer() == nil {
					fmt.Printf("  [%d] sensor no longer available\n", tick)
				} else {
					fmt.Printf("  [%d] %s (type=%d, min_delay=%d us)\n",
						tick, polled.Name(), polled.Type(), polled.MinDelay())
				}
				tick++
				time.Sleep(time.Second)
			}
			fmt.Printf("polling complete (%d samples)\n", tick)
		}

		return nil
	},
}

func init() {
	sensorReadCmd.Flags().Int32("type", int32(sensor.Accelerometer), "sensor type ID (1=accelerometer, 2=magnetic, 4=gyroscope, 5=light, ...)")
	sensorReadCmd.Flags().Duration("duration", 0, "how long to poll (0 = just print info)")

	sensorCmd.AddCommand(sensorReadCmd)
}
