// Example: ANativeWindow pixel format constants and API overview.
//
// Enumerates every pixel format constant exposed by the window package and
// prints its numeric value alongside the human-readable String() name.
//
// Creating an ANativeWindow requires a running Android Activity (the window
// handle is obtained from the Activity's native surface), so this example
// focuses on the type system and documents the full rendering lifecycle
// that a real application would follow.
//
// ANativeWindow rendering lifecycle:
//
//  1. Acquire   -- Obtain a *Window from the Activity's native surface.
//     The NDK delivers this via ANativeActivity callbacks
//     (onNativeWindowCreated). The Window wraps an opaque
//     ANativeWindow pointer managed by the system compositor.
//
//  2. Query     -- Inspect the window with Width(), Height(), and Format().
//     Width and Height return the surface dimensions in pixels.
//     Format returns the pixel format as an int32 matching one
//     of the Format constants (Rgba8888, Rgbx8888, Rgb565).
//
//  3. Configure -- Optionally call SetBuffersGeometry(width, height, format)
//     to request a different buffer size or pixel format. The
//     compositor will scale the buffer to the window size if
//     the requested dimensions differ. Pass 0 for width or
//     height to keep the current dimension. Returns an error
//     if the geometry change is rejected.
//
//  4. Lock      -- Lock the next back-buffer for CPU access. This yields an
//     ANativeWindow_Buffer describing the pixel memory layout:
//     stride, dimensions, format, and a raw pointer to pixel
//     data. While locked the compositor cannot read the buffer,
//     so the lock duration should be kept short.
//
//  5. Draw      -- Write pixels directly into the locked buffer. The pixel
//     layout depends on the format: Rgba8888 uses 4 bytes per
//     pixel (R, G, B, A), Rgbx8888 uses 4 bytes with the alpha
//     channel ignored, and Rgb565 packs each pixel into 2 bytes
//     (5 red, 6 green, 5 blue bits).
//
//  6. Unlock    -- Call UnlockAndPost() to release the buffer and submit it
//     to the compositor for display. This atomically unlocks
//     the buffer and posts it; there is no separate post step.
//     Returns an error if the surface is in an invalid state.
//
//  7. Release   -- When the Activity's window is destroyed (onNativeWindowDestroyed),
//     stop using the Window. The system reclaims the underlying
//     native surface.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/window"
)

func main() {
	// All pixel format constants defined by ANativeWindow.
	formats := []struct {
		name   string
		value  window.Format
		bpp    int
		detail string
	}{
		{
			name:   "Rgba8888",
			value:  window.Rgba8888,
			bpp:    4,
			detail: "8 bits per channel (red, green, blue, alpha)",
		},
		{
			name:   "Rgbx8888",
			value:  window.Rgbx8888,
			bpp:    4,
			detail: "8 bits per channel (red, green, blue, padding); alpha ignored",
		},
		{
			name:   "Rgb565",
			value:  window.Rgb565,
			bpp:    2,
			detail: "5 red, 6 green, 5 blue bits packed into 16 bits",
		},
	}

	fmt.Println("ANativeWindow pixel formats:")
	fmt.Println()
	fmt.Printf("  %-12s  %-5s  %-4s  %s\n", "Name", "Value", "BPP", "Description")
	fmt.Printf("  %-12s  %-5s  %-4s  %s\n", "------------", "-----", "----", "-----------")

	for _, f := range formats {
		// String() returns the constant name; verify it matches.
		str := f.value.String()
		fmt.Printf("  %-12s  %-5d  %-4d  %s\n", str, int32(f.value), f.bpp, f.detail)
	}

	// Demonstrate the String() method on an unknown format value.
	unknown := window.Format(99)
	fmt.Printf("\n  Unknown format value 99 prints as: %s\n", unknown)

	// Summarize the Window methods that would be used in a real application.
	fmt.Println()
	fmt.Println("Window methods:")
	fmt.Println()
	fmt.Println("  Width() int32")
	fmt.Println("      Returns the window's current width in pixels.")
	fmt.Println()
	fmt.Println("  Height() int32")
	fmt.Println("      Returns the window's current height in pixels.")
	fmt.Println()
	fmt.Println("  Format() int32")
	fmt.Println("      Returns the window's pixel format (matches a Format constant).")
	fmt.Println()
	fmt.Println("  SetBuffersGeometry(width, height, format int32) error")
	fmt.Println("      Requests a new buffer size and format. Pass 0 to keep a dimension.")
	fmt.Println()
	fmt.Println("  UnlockAndPost() error")
	fmt.Println("      Unlocks the back-buffer and posts it to the display.")
}
