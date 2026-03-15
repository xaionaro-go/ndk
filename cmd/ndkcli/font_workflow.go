package main

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/font"
)

// matchFontFromIterator iterates all system fonts and finds the best
// match for the given family, weight, and italic style.
//
// This replaces AFontMatcher_match which crashes (SIGSEGV) from headless
// CLI binaries. Root cause: AFontMatcher_match calls
// minikin::SystemFonts::findFontCollection(), which accesses a singleton
// populated by Java's Typeface static initializer via JNI
// (nativeAddFontCollections → minikin::SystemFonts::registerFallback).
// Without Java/ART running, the singleton's mDefaultFallback is null,
// causing a null pointer dereference in FontCollection::getFamilyForChar.
//
// ASystemFontIterator works from CLI because it reads /system/etc/fonts.xml
// directly from the filesystem, bypassing minikin's singleton entirely.
func matchFontFromIterator(
	family string,
	weight uint16,
	italic bool,
) (*font.Font, error) {
	iter := font.ASystemFontIterator_open()
	if iter == nil || iter.Pointer() == nil {
		return nil, fmt.Errorf("ASystemFontIterator_open returned nil")
	}
	defer iter.Close()

	var bestFont *font.Font
	bestScore := math.MaxInt32

	for {
		f := iter.Next()
		if f == nil || f.Pointer() == nil {
			break
		}

		// Check family match via file path (font files are named after families).
		path := f.GetFontFilePath()
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		baseLower := strings.ToLower(base)
		familyLower := strings.ToLower(family)

		// Skip fonts that don't match the family name.
		// Match by checking if the file name contains the family name.
		if !strings.Contains(baseLower, familyLower) &&
			!strings.Contains(familyLower, baseLower) {
			f.Close()
			continue
		}

		// Score: lower is better.
		// Weight distance + italic mismatch penalty.
		weightDist := int(weight) - int(f.Weight())
		if weightDist < 0 {
			weightDist = -weightDist
		}
		score := weightDist
		if f.IsItalic() != italic {
			score += 1000
		}

		if score < bestScore {
			if bestFont != nil {
				bestFont.Close()
			}
			bestFont = f
			bestScore = score
		} else {
			f.Close()
		}
	}

	return bestFont, nil
}

var fontMatchCmd = &cobra.Command{
	Use:   "match",
	Short: "Match a font by family name, weight, and italic style",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		family, _ := cmd.Flags().GetString("family")
		weight, _ := cmd.Flags().GetUint16("weight")
		italic, _ := cmd.Flags().GetBool("italic")

		fmt.Printf("matching family=%q weight=%d italic=%v...\n", family, weight, italic)

		matched, err := matchFontFromIterator(family, weight, italic)
		if err != nil {
			return fmt.Errorf("font matching: %w", err)
		}
		if matched == nil {
			fmt.Println("no matching font found")
			return nil
		}
		defer matched.Close()

		fmt.Printf("matched font:\n")
		fmt.Printf("  path:     %s\n", matched.GetFontFilePath())
		fmt.Printf("  weight:   %d\n", matched.Weight())
		fmt.Printf("  italic:   %v\n", matched.IsItalic())
		fmt.Printf("  locale:   %s\n", matched.GetLocale())

		return nil
	},
}

var fontListCmd = &cobra.Command{
	Use:   "list",
	Short: "List system fonts using ASystemFontIterator",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		limit, _ := cmd.Flags().GetInt("limit")

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
			fmt.Printf("  [%d] weight=%d italic=%-5v %s\n",
				count, f.Weight(), f.IsItalic(), f.GetFontFilePath())
			f.Close()
			count++
			if limit > 0 && count >= limit {
				fmt.Println("  ... (truncated)")
				break
			}
		}
		fmt.Printf("listed %d fonts\n", count)
		return nil
	},
}

func init() {
	fontMatchCmd.Flags().String("family", "Roboto", "font family name (matched against file name)")
	fontMatchCmd.Flags().Uint16("weight", 400, "font weight (100-900)")
	fontMatchCmd.Flags().Bool("italic", false, "request italic style")

	fontListCmd.Flags().Int("limit", 0, "max fonts to list (0=all)")

	fontCmd.AddCommand(fontMatchCmd)
	fontCmd.AddCommand(fontListCmd)
}
