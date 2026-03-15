package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/looper"
	"github.com/AndroidGoLab/ndk/sensor"
)

// sensorEventData extracts x, y, z from an ASensorEvent.
// ASensorEvent layout: version(4) + sensor(4) + type(4) + reserved(4) + timestamp(8) + data[16]float32
// data offset = 24 bytes from the start.
func sensorEventData(event *sensor.SensorEvent) (x, y, z float32, timestamp int64) {
	base := (*byte)(event.Pointer())
	// timestamp at offset 16 (after version+sensor+type+reserved = 4*4 = 16)
	timestamp = *(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(base)) + 16))
	// data[0..2] at offset 24
	dataPtr := uintptr(unsafe.Pointer(base)) + 24
	x = *(*float32)(unsafe.Pointer(dataPtr))
	y = *(*float32)(unsafe.Pointer(dataPtr + 4))
	z = *(*float32)(unsafe.Pointer(dataPtr + 8))
	return
}

var sensorReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read real-time sensor data",
	Long: `Reads real-time sensor events using a looper and event queue.
Prints x, y, z values for each event at the configured sample rate.`,
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		sensorTypeInt, _ := cmd.Flags().GetInt32("type")
		duration, _ := cmd.Flags().GetDuration("duration")
		rate, _ := cmd.Flags().GetInt32("rate")

		if duration == 0 {
			duration = 5 * time.Second
		}

		looper.Run(func(lp *looper.Looper) {
			_err = readSensorEvents(lp, sensor.Type(sensorTypeInt), duration, rate)
		})
		return
	},
}

func readSensorEvents(
	lp *looper.Looper,
	sensorType sensor.Type,
	duration time.Duration,
	rateUs int32,
) error {
	mgr := sensor.GetInstance()

	s := mgr.DefaultSensor(sensorType)
	if s.Pointer() == nil {
		return fmt.Errorf("no default sensor found for type %d", sensorType)
	}

	fmt.Printf("sensor: %s (%s)\n", s.Name(), sensor.Type(s.Type()))
	fmt.Printf("vendor: %s\n", s.Vendor())
	fmt.Printf("resolution: %g, min delay: %d us\n", s.Resolution(), s.MinDelay())

	// Create event queue on the current looper thread.
	const looperIdent = 1
	queue := mgr.CreateEventQueue(
		(*sensor.ALooper)(lp.Pointer()),
		looperIdent,
		sensor.ALooper_callbackFunc(nil), // no callback — we poll manually
		nil,
	)
	if queue == nil || queue.Pointer() == nil {
		return fmt.Errorf("failed to create sensor event queue")
	}
	defer mgr.DestroyEventQueue(queue)

	// Enable the sensor on the queue and set the sample rate.
	if err := queue.EnableSensor(s); err != nil {
		return fmt.Errorf("enabling sensor: %w", err)
	}
	if err := queue.SetEventRate(s, rateUs); err != nil {
		return fmt.Errorf("setting event rate: %w", err)
	}

	fmt.Printf("\nstreaming for %v (rate=%d us)...\n", duration, rateUs)

	deadline := time.Now().Add(duration)
	count := 0
	// Allocate a raw ASensorEvent buffer (104 bytes per event on arm64).
	// SensorEvent.ptr must point to valid C memory for GetEvents to write into.
	const eventSize = 104 // sizeof(ASensorEvent)
	eventBuf := make([]byte, eventSize)
	event := sensor.NewSensorEventFromPointer(unsafe.Pointer(&eventBuf[0]))

	for time.Now().Before(deadline) {
		// Poll the looper for events (100ms timeout).
		looper.PollOnce(100*time.Millisecond, nil, nil, nil)

		// Drain all available events.
		for {
			n := queue.GetEvents(event, 1)
			if n <= 0 {
				break
			}
			x, y, z, ts := sensorEventData(event)
			fmt.Printf("  [%d] t=%d  x=%.6f  y=%.6f  z=%.6f\n", count, ts, x, y, z)
			count++
		}
	}

	fmt.Printf("received %d events\n", count)
	return nil
}

func init() {
	sensorReadCmd.Flags().Int32("type", int32(sensor.Accelerometer), "sensor type ID (1=accelerometer, 2=magnetic, 3=orientation, 4=gyroscope, 5=light, 8=proximity)")
	sensorReadCmd.Flags().Duration("duration", 5*time.Second, "how long to stream")
	sensorReadCmd.Flags().Int32("rate", 100000, "sample period in microseconds (default 100ms)")

	sensorCmd.AddCommand(sensorReadCmd)
}
