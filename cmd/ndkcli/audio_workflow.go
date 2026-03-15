package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/audio"
)

var audioRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record raw PCM16 audio to a file",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		output, _ := cmd.Flags().GetString("output")
		duration, _ := cmd.Flags().GetDuration("duration")
		sampleRate, _ := cmd.Flags().GetInt32("sample-rate")
		channels, _ := cmd.Flags().GetInt32("channels")

		builder, err := audio.NewStreamBuilder()
		if err != nil {
			return fmt.Errorf("creating stream builder: %w", err)
		}
		defer builder.Close()

		builder.
			SetDirection(audio.Input).
			SetSampleRate(sampleRate).
			SetChannelCount(channels).
			SetFormat(audio.PcmI16)

		stream, err := builder.Open()
		if err != nil {
			return fmt.Errorf("opening stream: %w", err)
		}
		defer func() {
			if closeErr := stream.Close(); closeErr != nil && _err == nil {
				_err = fmt.Errorf("closing stream: %w", closeErr)
			}
		}()

		if err := stream.Start(); err != nil {
			return fmt.Errorf("starting stream: %w", err)
		}
		defer func() {
			if stopErr := stream.Stop(); stopErr != nil && _err == nil {
				_err = fmt.Errorf("stopping stream: %w", stopErr)
			}
		}()

		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()

		framesPerBurst := stream.FramesPerBurst()
		// PCM16: 2 bytes per sample per channel
		bytesPerFrame := channels * 2
		buf := make([]byte, framesPerBurst*bytesPerFrame)
		timeout := 100 * time.Millisecond

		deadline := time.Now().Add(duration)
		var totalFrames int64
		for time.Now().Before(deadline) {
			framesRead, err := stream.Read(buf, framesPerBurst, timeout)
			if err != nil {
				return fmt.Errorf("reading from stream: %w", err)
			}
			if framesRead > 0 {
				bytesRead := framesRead * bytesPerFrame
				if _, err := f.Write(buf[:bytesRead]); err != nil {
					return fmt.Errorf("writing to file: %w", err)
				}
				totalFrames += int64(framesRead)
			}
		}

		fmt.Printf("recorded %d frames (%d bytes) to %s\n", totalFrames, totalFrames*int64(bytesPerFrame), output)
		return nil
	},
}

var audioPlayCmd = &cobra.Command{
	Use:   "play",
	Short: "Play raw PCM16 audio from a file",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		input, _ := cmd.Flags().GetString("input")
		sampleRate, _ := cmd.Flags().GetInt32("sample-rate")
		channels, _ := cmd.Flags().GetInt32("channels")

		builder, err := audio.NewStreamBuilder()
		if err != nil {
			return fmt.Errorf("creating stream builder: %w", err)
		}
		defer builder.Close()

		builder.
			SetDirection(audio.Output).
			SetSampleRate(sampleRate).
			SetChannelCount(channels).
			SetFormat(audio.PcmI16)

		stream, err := builder.Open()
		if err != nil {
			return fmt.Errorf("opening stream: %w", err)
		}
		defer func() {
			if closeErr := stream.Close(); closeErr != nil && _err == nil {
				_err = fmt.Errorf("closing stream: %w", closeErr)
			}
		}()

		if err := stream.Start(); err != nil {
			return fmt.Errorf("starting stream: %w", err)
		}
		defer func() {
			if stopErr := stream.Stop(); stopErr != nil && _err == nil {
				_err = fmt.Errorf("stopping stream: %w", stopErr)
			}
		}()

		f, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}
		defer f.Close()

		framesPerBurst := stream.FramesPerBurst()
		// PCM16: 2 bytes per sample per channel
		bytesPerFrame := channels * 2
		buf := make([]byte, framesPerBurst*bytesPerFrame)
		timeout := 100 * time.Millisecond

		var totalFrames int64
		for {
			n, err := f.Read(buf)
			switch {
			case err == nil:
			case err == io.EOF:
				fmt.Printf("played %d frames from %s\n", totalFrames, input)
				return nil
			default:
				return fmt.Errorf("reading from file: %w", err)
			}

			framesToWrite := int32(n) / bytesPerFrame
			if framesToWrite == 0 {
				continue
			}

			framesWritten, err := stream.Write(buf[:framesToWrite*bytesPerFrame], framesToWrite, timeout)
			if err != nil {
				return fmt.Errorf("writing to stream: %w", err)
			}
			totalFrames += int64(framesWritten)
		}
	},
}

func init() {
	audioRecordCmd.Flags().String("output", "recording.pcm", "output file path")
	audioRecordCmd.Flags().Duration("duration", 5*time.Second, "recording duration")
	audioRecordCmd.Flags().Int32("sample-rate", 44100, "sample rate in Hz")
	audioRecordCmd.Flags().Int32("channels", 1, "number of audio channels")

	audioPlayCmd.Flags().String("input", "recording.pcm", "input file path")
	audioPlayCmd.Flags().Int32("sample-rate", 44100, "sample rate in Hz")
	audioPlayCmd.Flags().Int32("channels", 1, "number of audio channels")

	audioCmd.AddCommand(audioRecordCmd)
	audioCmd.AddCommand(audioPlayCmd)
}
