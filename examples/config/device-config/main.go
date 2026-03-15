// Example: query device configuration properties.
//
// Creates an AConfiguration, reads every available property (language,
// country, density, orientation, screen size, screen dimensions, and SDK
// version), and prints them. Language and Country are two-character ISO
// codes written through an out-parameter.
//
// On a real device the values reflect the current system locale and
// display settings. When run with a freshly created (empty) config, most
// fields return zero or blank, which is the expected default behavior.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/AndroidGoLab/ndk/config"
)

// orientationName returns a human-readable label for an orientation code.
func orientationName(v int32) string {
	switch config.Orientation(v) {
	case config.OrientationAny:
		return "Any"
	case config.OrientationPort:
		return "Portrait"
	case config.OrientationLand:
		return "Landscape"
	case config.OrientationSquare:
		return "Square"
	default:
		return fmt.Sprintf("Unknown(%d)", v)
	}
}

// screenSizeName returns a human-readable label for a screen-size code.
func screenSizeName(v int32) string {
	switch config.ScreenSize(v) {
	case config.ScreensizeAny:
		return "Any"
	case config.ScreensizeSmall:
		return "Small"
	case config.ScreensizeNormal:
		return "Normal"
	case config.ScreensizeLarge:
		return "Large"
	case config.ScreensizeXlarge:
		return "Xlarge"
	default:
		return fmt.Sprintf("Unknown(%d)", v)
	}
}

func main() {
	cfg := config.NewConfig()
	defer cfg.Close()

	// Language and Country write a 2-character code into a string buffer.
	// Allocate 2-byte buffers as strings for each.
	langBuf := string(make([]byte, 2))
	countryBuf := string(make([]byte, 2))
	cfg.Language(langBuf)
	cfg.Country(countryBuf)

	fmt.Println("Device configuration:")
	fmt.Println()
	fmt.Printf("  Language:        %q\n", langBuf)
	fmt.Printf("  Country:         %q\n", countryBuf)
	fmt.Printf("  Density:         %d dpi\n", cfg.Density())
	fmt.Printf("  Orientation:     %s (%d)\n", orientationName(cfg.Orientation()), cfg.Orientation())
	fmt.Printf("  Screen size:     %s (%d)\n", screenSizeName(cfg.ScreenSize()), cfg.ScreenSize())
	fmt.Printf("  Screen width:    %d dp\n", cfg.ScreenWidthDp())
	fmt.Printf("  Screen height:   %d dp\n", cfg.ScreenHeightDp())
	fmt.Printf("  SDK version:     %d\n", cfg.SdkVersion())
}
