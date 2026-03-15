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
		matcher.SetLocales("en-US")
		fmt.Println("style set OK")

		// AFontMatcher_match requires the Android font system (libminikin).
		// On headless CLI binaries (no app/Activity), the internal
		// FontCollection is null, causing SIGSEGV in getFamilyForChar.
		// This command only works within an Android app context.
		fmt.Println("WARNING: font matching requires an Android app context (Activity/Service).")
		fmt.Println("On headless CLI binaries, AFontMatcher_match will crash with SIGSEGV")
		fmt.Println("because the system FontCollection is not initialized.")
		text := []uint16{'A'}
		var runLength uint32
		fmt.Printf("matching family=%q weight=%d italic=%v...\n", family, weight, italic)

		matched := matcher.Match(family, &text[0], uint32(len(text)), &runLength)
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

var fontListCmd = &cobra.Command{
	Use:   "list",
	Short: "List system fonts using ASystemFontIterator",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		iter := font.ASystemFontIterator_open()
		if iter == nil || iter.Pointer() == nil {
			return fmt.Errorf("ASystemFontIterator_open returned nil")
		}
		defer iter.Close()

		count := 0
		for {
			f := iter.Next()
			if f == nil || f.Pointer() == nil {
				break
			}
			fmt.Printf("  [%d] weight=%d italic=%v path=%s\n",
				count, f.Weight(), f.IsItalic(), f.GetFontFilePath())
			f.Close()
			count++
			if count >= 20 {
				fmt.Println("  ... (truncated)")
				break
			}
		}
		fmt.Printf("listed %d fonts\n", count)
		return nil
	},
}

func init() {
	fontMatchCmd.Flags().String("family", "sans-serif", "font family name")
	fontMatchCmd.Flags().Uint16("weight", uint16(font.Normal), "font weight (100-900)")
	fontMatchCmd.Flags().Bool("italic", false, "request italic style")

	fontCmd.AddCommand(fontMatchCmd)
	fontCmd.AddCommand(fontListCmd)
}
