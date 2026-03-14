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

		matcher.SetStyle(weight, italic)

		matched := matcher.Match(family, nil, 0, nil)
		if matched.Pointer() == nil {
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
