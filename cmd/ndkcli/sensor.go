package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/sensor"
)

// commonSensorTypes lists sensor types to probe when listing all sensors.
// The sensor manager has no "list all" API, so we check each known type.
var commonSensorTypes = []sensor.Type{
	sensor.Accelerometer,
	sensor.MagneticField,
	sensor.Gyroscope,
	sensor.Light,
	sensor.Pressure,
	sensor.Proximity,
	sensor.Gravity,
	sensor.LinearAcceleration,
	sensor.RotationVector,
	sensor.RelativeHumidity,
	sensor.AmbientTemperature,
	sensor.MagneticFieldUncalibrated,
	sensor.GameRotationVector,
	sensor.GyroscopeUncalibrated,
	sensor.SignificantMotion,
	sensor.StepDetector,
	sensor.StepCounter,
	sensor.GeomagneticRotationVector,
	sensor.HeartRate,
	sensor.AccelerometerUncalibrated,
	sensor.HingeAngle,
}

func printSensor(s *sensor.Sensor) {
	fmt.Printf("  Name:       %s\n", s.Name())
	fmt.Printf("  Vendor:     %s\n", s.Vendor())
	fmt.Printf("  Type:       %d\n", s.Type())
	fmt.Printf("  Resolution: %g\n", s.Resolution())
	fmt.Printf("  MinDelay:   %d us\n", s.MinDelay())
}

var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "Query sensors via the NDK Sensor API",
}

var sensorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available sensors by probing known types",
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) (_err error) {
		mgr := sensor.GetInstance()

		found := 0
		for _, t := range commonSensorTypes {
			s := mgr.DefaultSensor(int32(t))
			if s.Pointer() == nil {
				continue
			}
			found++
			fmt.Printf("[%s]\n", t)
			printSensor(s)
			fmt.Println()
		}

		if found == 0 {
			fmt.Println("No sensors found.")
		}
		return nil
	},
}

var sensorInfoCmd = &cobra.Command{
	Use:   "info <type>",
	Short: "Print details for a sensor type (integer value)",
	Args:  cobra.ExactArgs(1),
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) (_err error) {
		typeVal, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("parsing sensor type: %w", err)
		}

		mgr := sensor.GetInstance()
		s := mgr.DefaultSensor(int32(typeVal))
		if s.Pointer() == nil {
			return fmt.Errorf("no default sensor for type %d", typeVal)
		}

		fmt.Printf("[Sensor type %d]\n", typeVal)
		printSensor(s)
		return nil
	},
}

func init() {
	sensorCmd.AddCommand(sensorListCmd)
	sensorCmd.AddCommand(sensorInfoCmd)
	rootCmd.AddCommand(sensorCmd)
}
