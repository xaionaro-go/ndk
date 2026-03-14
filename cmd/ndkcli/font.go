package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/font"
)

var fontCmd = &cobra.Command{
	Use:   "font",
	Short: "Font NDK operations",
}

var (
	fontMatchFamily string
	fontMatchWeight uint16
	fontMatchItalic bool
)

var fontMatchCmd = &cobra.Command{
	Use:   "match",
	Short: "Match a font by family name, weight, and italic",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := font.NewMatcher()
		defer m.Close()

		m.SetStyle(fontMatchWeight, fontMatchItalic)

		// Match requires text to shape; pass a single space as minimal input.
		text := []uint16{' '}
		var runLength uint32
		f := m.Match(fontMatchFamily, &text[0], uint32(len(text)), &runLength)
		if f == nil {
			fmt.Println("No font matched.")
			return nil
		}
		defer f.Close()

		fmt.Printf("Matched font for family %q:\n", fontMatchFamily)
		fmt.Printf("  Weight: %d\n", f.Weight())
		fmt.Printf("  Italic: %v\n", f.IsItalic())
		return nil
	},
}

func init() {
	fontMatchCmd.Flags().StringVar(&fontMatchFamily, "family", "sans-serif", "font family name")
	fontMatchCmd.Flags().Uint16Var(&fontMatchWeight, "weight", uint16(font.Normal), "font weight (100-900)")
	fontMatchCmd.Flags().BoolVar(&fontMatchItalic, "italic", false, "request italic variant")

	fontCmd.AddCommand(fontMatchCmd)
	rootCmd.AddCommand(fontCmd)
}
