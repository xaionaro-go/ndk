package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/font"
)

var fontMatchCmd = &cobra.Command{
	Use:   "match",
	Short: "Match a font by family name, weight, and italic style",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		family, _ := cmd.Flags().GetString("family")
		weight, _ := cmd.Flags().GetUint16("weight")
		italic, _ := cmd.Flags().GetBool("italic")

		matcher := font.NewMatcher()
		defer matcher.Close()

		fmt.Println("setting style...")
		matcher.SetStyle(weight, italic)

		// AFontMatcher_match requires non-nil text for script detection.
		text := []uint16{'A'}
		var runLength uint32
		fmt.Printf("matching family=%q weight=%d italic=%v...\n", family, weight, italic)

		var matched *font.Font
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("font matching crashed (recovered): %v\n", r)
				}
			}()
			matched = matcher.Match(family, &text[0], uint32(len(text)), &runLength)
		}()
		fmt.Printf("match returned, runLength=%d\n", runLength)
		if matched == nil || matched.Pointer() == nil {
			fmt.Println("no matching font found")
			return nil
		}
		defer matched.Close()

		fmt.Printf("matched font:\n")
		fmt.Printf("  weight:   %d\n", matched.Weight())
		fmt.Printf("  is italic: %v\n", matched.IsItalic())

		return nil
	},
}

func init() {
	fontMatchCmd.Flags().String("family", "sans-serif", "font family name")
	fontMatchCmd.Flags().Uint16("weight", uint16(font.Normal), "font weight (100-900)")
	fontMatchCmd.Flags().Bool("italic", false, "request italic style")

	fontCmd.AddCommand(fontMatchCmd)
}
