// Example: survey system font coverage across weight values.
//
// Iterates through the standard font weight constants (Thin through Black)
// and asks the system font matcher which font it selects for each requested
// weight against the "sans-serif" family. This reveals which weights the
// device actually has installed and how the system maps missing weights to
// available ones.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/font"
)

func main() {
	matcher := font.NewMatcher()
	if matcher == nil {
		log.Fatal("failed to create font matcher")
	}
	defer matcher.Close()

	text := []uint16{'A', 'B', 'C'}
	length := uint32(len(text))

	type weightInfo struct {
		name   string
		weight font.Weight
	}

	weights := []weightInfo{
		{"Thin", font.Thin},
		{"Light", font.Light},
		{"Normal", font.Normal},
		{"Medium", font.Medium},
		{"Bold", font.Bold},
		{"Black", font.Black},
	}

	fmt.Println("Font weight survey for \"sans-serif\":")
	fmt.Println()
	fmt.Printf("  %-10s  %-10s  %-10s  %s\n",
		"Requested", "Value", "Matched", "Italic")
	fmt.Printf("  %-10s  %-10s  %-10s  %s\n",
		"----------", "-----", "-------", "------")

	for _, w := range weights {
		var runLength uint32
		matcher.SetStyle(uint16(w.weight), false)
		f := matcher.Match("sans-serif", &text[0], length, &runLength)
		if f == nil {
			fmt.Printf("  %-10s  %-10d  %-10s  %s\n",
				w.name, int(w.weight), "nil", "-")
			continue
		}

		italic := "no"
		if f.IsItalic() {
			italic = "yes"
		}
		fmt.Printf("  %-10s  %-10d  %-10d  %s\n",
			w.name, int(w.weight), f.Weight(), italic)
		f.Close()
	}
}
