package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/audio"
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Query audio capabilities via the NDK AAudio API",
}

var audioInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Open a probe stream and print audio system properties",
	RunE: func(
		cmd *cobra.Command,
		args []string,
	) (_err error) {
		builder, err := audio.NewStreamBuilder()
		if err != nil {
			return fmt.Errorf("creating stream builder: %w", err)
		}
		defer builder.Close()

		builder.
			SetSampleRate(44100).
			SetChannelCount(2).
			SetFormat(audio.PcmFloat).
			SetDirection(audio.Output).
			SetPerformanceMode(audio.LowLatency)

		stream, err := builder.Open()
		if err != nil {
			return fmt.Errorf("opening audio stream: %w", err)
		}
		defer stream.Close()

		fmt.Printf("Sample Rate:      %d Hz\n", stream.SampleRate())
		fmt.Printf("Channel Count:    %d\n", stream.ChannelCount())
		fmt.Printf("Frames Per Burst: %d\n", stream.FramesPerBurst())
		fmt.Printf("State:            %s\n", stream.State())
		fmt.Printf("XRun Count:       %d\n", stream.XRunCount())

		return nil
	},
}

func init() {
	audioCmd.AddCommand(audioInfoCmd)
	rootCmd.AddCommand(audioCmd)
}
