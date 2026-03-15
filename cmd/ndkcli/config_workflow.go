package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/config"
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display device configuration (density, orientation, screen, SDK version)",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		cfg := config.NewConfig()
		defer cfg.Close()

		fmt.Printf("density:          %d\n", cfg.Density())
		fmt.Printf("orientation:      %d (%s)\n", cfg.Orientation(), config.Orientation(cfg.Orientation()))
		fmt.Printf("screen width dp:  %d\n", cfg.ScreenWidthDp())
		fmt.Printf("screen height dp: %d\n", cfg.ScreenHeightDp())
		fmt.Printf("screen size:      %d (%s)\n", cfg.ScreenSize(), config.ScreenSize(cfg.ScreenSize()))
		fmt.Printf("SDK version:      %d\n", cfg.SdkVersion())

		// Country and Language take a string buffer; pass a zero-filled
		// string of sufficient length so the C function can write into it.
		countryBuf := strings.Repeat("\x00", 4)
		cfg.Country(countryBuf)
		fmt.Printf("country:          %s\n", strings.TrimRight(countryBuf, "\x00"))

		languageBuf := strings.Repeat("\x00", 4)
		cfg.Language(languageBuf)
		fmt.Printf("language:         %s\n", strings.TrimRight(languageBuf, "\x00"))

		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}
