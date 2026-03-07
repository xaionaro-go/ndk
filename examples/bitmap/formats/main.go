// Example: Android bitmap format constants and error codes.
//
// Enumerates every bitmap pixel format defined by the NDK and prints its
// name and underlying integer value. Each format corresponds to an
// ANDROID_BITMAP_FORMAT_* constant used when creating or inspecting
// bitmaps through the jnigraphics API.
//
// The program also demonstrates the Error type, which implements the
// standard error interface and wraps NDK result codes returned by
// bitmap operations.
//
// This program does not require a live Android device; it only reads
// compile-time constants.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/bitmap"
)

// formatInfo pairs a Format constant with a human-readable description.
type formatInfo struct {
	name   string
	value  bitmap.Format
	detail string
}

// formats lists every bitmap.Format constant together with a short
// explanation of its pixel layout and typical use case.
var formats = []formatInfo{
	{
		name:   "None",
		value:  bitmap.None,
		detail: "No format / unknown. Returned when the bitmap format cannot be determined.",
	},
	{
		// 32-bit, 8 bits per channel (red, green, blue, alpha).
		// The most common format; suitable for general-purpose rendering
		// and compositing where full color depth and transparency are needed.
		name:   "RGBA_8888",
		value:  bitmap.Rgba8888,
		detail: "32-bit standard color (8 bits per channel). Default for most bitmaps.",
	},
	{
		// 16-bit, 5 bits red / 6 bits green / 5 bits blue, no alpha.
		// Uses half the memory of RGBA_8888 at the cost of color precision.
		// Common in memory-constrained scenarios such as live wallpapers or
		// thumbnail caches.
		name:   "RGB_565",
		value:  bitmap.Rgb565,
		detail: "16-bit color (5-6-5), no alpha. Half the memory of RGBA_8888.",
	},
	{
		// 16-bit, 4 bits per channel. Deprecated in API 13 in favor of
		// RGBA_8888 because the limited color depth introduces visible
		// banding. Included for completeness.
		name:   "RGBA_4444",
		value:  bitmap.Rgba4444,
		detail: "16-bit color (4 bits per channel). Deprecated; prefer RGBA_8888.",
	},
	{
		// 8-bit alpha-only. Each pixel stores a single alpha value; no
		// color information. Used for masks and text rendering.
		name:   "A_8",
		value:  bitmap.A8,
		detail: "8-bit alpha only. Used for masks and glyph rendering.",
	},
	{
		// 64-bit, 16-bit IEEE 754 half-float per channel.
		// Required for HDR content and wide color gamut (Display P3,
		// scRGB). Available from API 26.
		name:   "RGBA_F16",
		value:  bitmap.RgbaF16,
		detail: "64-bit HDR (16-bit float per channel). For wide-gamut / HDR content.",
	},
}

// errorInfo pairs an Error constant with its meaning.
type errorInfo struct {
	name  string
	value bitmap.Error
}

// errors lists every bitmap.Error constant.
var errors = []errorInfo{
	{"BadParameter", bitmap.ErrBadParameter},
	{"JniException", bitmap.ErrJniException},
	{"AllocationFailed", bitmap.ErrAllocationFailed},
}

func main() {
	fmt.Println("Bitmap pixel formats")
	fmt.Println()
	for _, f := range formats {
		fmt.Printf("  %-12s  value=%d   %s\n", f.name, int32(f.value), f.detail)
	}

	fmt.Println()
	fmt.Println("Bitmap error codes")
	fmt.Println()
	for _, e := range errors {
		// Error implements the error interface; call it to show the
		// formatted message that would appear in practice.
		fmt.Printf("  %-20s  code=%2d   error.Error() = %q\n",
			e.name, int32(e.value), e.value.Error())
	}
}
