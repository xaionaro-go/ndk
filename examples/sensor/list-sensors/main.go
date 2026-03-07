// Example: list default sensors and their properties.
//
// Obtains the sensor manager singleton and queries default sensors for
// the most common types (accelerometer, gyroscope, light, proximity,
// magnetic field). For each sensor that is present on the device, the
// program prints its name, vendor, type code, resolution, and minimum
// delay between events.
//
// Not every device has every sensor. Because DefaultSensor always
// returns a non-nil *Sensor wrapper (even when the underlying C
// pointer is NULL), calling methods on an absent sensor would crash.
// This example uses a recover guard so that missing sensors are
// reported gracefully instead of terminating the program.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/sensor"
)

// printSensor queries and prints sensor properties. It returns false
// if the sensor's underlying C pointer is NULL (the device lacks this
// sensor type), recovering from the resulting panic.
func printSensor(mgr *sensor.Manager, label string, sensorType sensor.Type) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	s := mgr.DefaultSensor(int32(sensorType))

	// Trigger a method call; if the internal pointer is NULL the NDK
	// dereferences a null pointer and Go's signal handler turns it
	// into a panic that we recover above.
	name := s.Name()
	vendor := s.Vendor()
	if name == "" || vendor == "" {
		return false
	}

	fmt.Printf("  %s:\n", label)
	fmt.Printf("    Name:       %s\n", name)
	fmt.Printf("    Vendor:     %s\n", vendor)
	fmt.Printf("    Type:       %s (%d)\n", sensorType, int32(sensorType))
	fmt.Printf("    Resolution: %g\n", s.Resolution())
	fmt.Printf("    Min delay:  %d us\n", s.MinDelay())
	fmt.Println()
	return true
}

func main() {
	mgr := sensor.GetInstance()

	type sensorInfo struct {
		label      string
		sensorType sensor.Type
	}

	sensors := []sensorInfo{
		{"Accelerometer", sensor.Accelerometer},
		{"Gyroscope", sensor.Gyroscope},
		{"Light", sensor.Light},
		{"Proximity", sensor.Proximity},
		{"Magnetic Field", sensor.MagneticField},
	}

	fmt.Println("Default sensors on this device:")
	fmt.Println()

	found := 0
	for _, info := range sensors {
		if printSensor(mgr, info.label, info.sensorType) {
			found++
		} else {
			fmt.Printf("  %-16s  not available\n", info.label+":")
		}
	}

	if found == 0 {
		fmt.Println("  No default sensors found on this device.")
	}
}
