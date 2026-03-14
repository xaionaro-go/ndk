package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/camera"
)

// Tag values for commonly-used ACAMERA_* metadata entries.
const (
	tagLensFacing                 = 524293  // ACAMERA_LENS_FACING
	tagSensorOrientation          = 917518  // ACAMERA_SENSOR_ORIENTATION
	tagInfoSupportedHardwareLevel = 1376256 // ACAMERA_INFO_SUPPORTED_HARDWARE_LEVEL
)

func lensFacingString(v int32) string {
	switch v {
	case 0:
		return "Front"
	case 1:
		return "Back"
	case 2:
		return "External"
	default:
		return fmt.Sprintf("Unknown(%d)", v)
	}
}

func hardwareLevelString(v int32) string {
	switch v {
	case 0:
		return "Limited"
	case 1:
		return "Full"
	case 2:
		return "Legacy"
	case 3:
		return "Level3"
	case 4:
		return "External"
	default:
		return fmt.Sprintf("Unknown(%d)", v)
	}
}

var cameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "Query cameras via the NDK Camera2 API",
}

var cameraListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available camera IDs",
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) (_err error) {
		mgr := camera.NewManager()
		defer mgr.Close()

		ids, err := mgr.CameraIDList()
		if err != nil {
			return fmt.Errorf("listing cameras: %w", err)
		}

		for _, id := range ids {
			fmt.Println(id)
		}
		return nil
	},
}

var cameraInfoCmd = &cobra.Command{
	Use:   "info <id>",
	Short: "Print characteristics for a camera",
	Args:  cobra.ExactArgs(1),
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) (_err error) {
		cameraID := args[0]

		mgr := camera.NewManager()
		defer mgr.Close()

		meta, err := mgr.GetCameraCharacteristics(cameraID)
		if err != nil {
			return fmt.Errorf("getting characteristics for camera %q: %w", cameraID, err)
		}
		defer meta.Close()

		fmt.Printf("Camera ID:        %s\n", cameraID)

		if n := meta.I32Count(tagLensFacing); n > 0 {
			fmt.Printf("Lens Facing:      %s\n", lensFacingString(meta.I32At(tagLensFacing, 0)))
		}

		if n := meta.I32Count(tagSensorOrientation); n > 0 {
			fmt.Printf("Sensor Orient.:   %d degrees\n", meta.I32At(tagSensorOrientation, 0))
		}

		if n := meta.I32Count(tagInfoSupportedHardwareLevel); n > 0 {
			fmt.Printf("Hardware Level:   %s\n", hardwareLevelString(meta.I32At(tagInfoSupportedHardwareLevel, 0)))
		}

		return nil
	},
}

func init() {
	cameraCmd.AddCommand(cameraListCmd)
	cameraCmd.AddCommand(cameraInfoCmd)
	rootCmd.AddCommand(cameraCmd)
}
