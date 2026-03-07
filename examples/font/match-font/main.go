// Example: match system fonts by family name and style.
//
// Creates a font matcher, requests a normal-weight non-italic font from the
// "sans-serif" family for a sample text string, and prints the matched font's
// weight and italic status. Then requests a bold italic match to show how the
// SetStyle chaining pattern works with different style parameters.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/font"
)

// sampleText returns a uint16 slice containing the Unicode code points for
// "Hello" followed by the text length. The AFontMatcher_match function
// takes UTF-16 text to determine which font best covers the characters.
func sampleText() ([]uint16, uint32) {
	text := []uint16{'H', 'e', 'l', 'l', 'o'}
	return text, uint32(len(text))
}

func printFont(label string, f *font.Font) {
	if f == nil {
		log.Printf("%s: no font matched (nil)", label)
		return
	}
	defer f.Close()

	italic := "no"
	if f.IsItalic() {
		italic = "yes"
	}
	fmt.Printf("  %s:\n", label)
	fmt.Printf("    Weight: %d\n", f.Weight())
	fmt.Printf("    Italic: %s\n", italic)
}

func main() {
	matcher := font.NewMatcher()
	if matcher == nil {
		log.Fatal("failed to create font matcher")
	}
	defer matcher.Close()

	text, length := sampleText()

	// Match a normal-weight, non-italic font.
	fmt.Println("Font matching results for \"sans-serif\":")
	fmt.Println()

	var runLength uint32
	matcher.SetStyle(uint16(font.Normal), false)
	normalFont := matcher.Match("sans-serif", &text[0], length, &runLength)
	printFont("Normal (400, upright)", normalFont)
	fmt.Printf("    Run length: %d\n", runLength)
	fmt.Println()

	// Match a bold italic font using the chaining pattern.
	// SetStyle returns the matcher, so the call reads fluently.
	boldItalicFont := matcher.
		SetStyle(uint16(font.Bold), true).
		Match("sans-serif", &text[0], length, &runLength)
	printFont("Bold italic (700, italic)", boldItalicFont)
	fmt.Printf("    Run length: %d\n", runLength)
}
