// AImageDecoder API overview.
//
// Documents the Android AImageDecoder pipeline provided by the ndk/image
// package. AImageDecoder (introduced in API 30) decodes JPEG, PNG, GIF,
// WebP, BMP, and ICO images into raw pixel buffers without requiring the
// Java Bitmap machinery.
//
// Decoder lifecycle:
//
//  1. Create  -- from a file descriptor or in-memory buffer
//  2. Inspect -- query header info (width, height, MIME type)
//  3. Configure (optional) -- set target size for downscaling
//  4. Stride  -- call MinimumStride to learn the row byte count
//  5. Allocate -- create a pixel buffer of stride * height bytes
//  6. Decode  -- fill the buffer with decoded pixel data
//  7. Close   -- release the decoder
//
// The Decoder cannot be created directly from Go: it requires either a
// POSIX file descriptor (AImageDecoder_createFromFd) or an in-memory
// buffer (AImageDecoder_createFromBuffer). These factory functions are
// not yet exposed in the high-level image package; they can be accessed
// through the raw NDK layer when needed.
//
// Creation pattern (pseudocode):
//
//	var decoder *image.Decoder
//	fd := openImageFile()  // POSIX file descriptor
//	// create decoder from fd (factory not yet in high-level API)
//	defer decoder.Close()
//
// Once the Decoder is obtained, the high-level API takes over.
//
// This program must run on an Android device with API 30+.
package main

import (
	"fmt"

	"github.com/AndroidGoLab/ndk/image"
)

// errorInfo pairs an Error constant with its name and meaning.
type errorInfo struct {
	name   string
	value  image.Error
	detail string
}

// errors lists every image.Error constant with a description of when it
// is returned by the NDK.
var errors = []errorInfo{
	{
		name:   "Incomplete",
		value:  image.ErrIncomplete,
		detail: "The input data was truncated. The decoder may have partially decoded the image.",
	},
	{
		name:   "Error",
		value:  image.ErrError,
		detail: "Generic decoding error. Returned when no more specific code applies.",
	},
	{
		name:   "InvalidConversion",
		value:  image.ErrInvalidConversion,
		detail: "The requested pixel format conversion is not supported.",
	},
	{
		name:   "InvalidScale",
		value:  image.ErrInvalidScale,
		detail: "The target size set via SetTargetSize is invalid (e.g. zero or negative).",
	},
	{
		name:   "BadParameter",
		value:  image.ErrBadParameter,
		detail: "A nil pointer or otherwise invalid argument was passed.",
	},
	{
		name:   "InvalidInput",
		value:  image.ErrInvalidInput,
		detail: "The input data is not a recognized image format.",
	},
	{
		name:   "SeekError",
		value:  image.ErrSeekError,
		detail: "The data source does not support seeking, which is required for decoding.",
	},
	{
		name:   "InternalError",
		value:  image.ErrInternalError,
		detail: "An unexpected internal error occurred in the decoder.",
	},
	{
		name:   "UnsupportedFormat",
		value:  image.ErrUnsupportedFormat,
		detail: "The image format is recognized but not supported for decoding.",
	},
	{
		name:   "Finished",
		value:  image.ErrFinished,
		detail: "The decoder has already finished decoding (e.g. single-frame image decoded twice).",
	},
	{
		name:   "InvalidState",
		value:  image.ErrInvalidState,
		detail: "The decoder is in a state that does not allow the requested operation.",
	},
}

func main() {
	fmt.Println("=== ndk/image API overview ===")
	fmt.Println()

	// ---------------------------------------------------------------
	// Decoder creation
	//
	// The Decoder struct wraps an NDK AImageDecoder pointer. It is
	// NOT created directly from Go. The Android NDK provides two
	// factory functions:
	//
	//   AImageDecoder_createFromFd(fd, out)
	//       Creates a decoder from a POSIX file descriptor. The fd
	//       must be open for reading and seekable. The caller retains
	//       ownership of the fd and must close it after the decoder
	//       is finished.
	//
	//   AImageDecoder_createFromBuffer(buf, length, out)
	//       Creates a decoder from an in-memory buffer. The buffer
	//       must remain valid for the lifetime of the decoder.
	//
	// Both return 0 (ANDROID_IMAGE_DECODER_SUCCESS) on success, or a
	// negative error code that can be converted to image.Error.
	//
	// These factory functions are not yet wrapped in the high-level
	// image package; they are available in the raw NDK layer.
	// ---------------------------------------------------------------

	fmt.Println("Decoder creation:")
	fmt.Println()
	fmt.Println("  From file descriptor:")
	fmt.Println("      AImageDecoder_createFromFd(int32(fd), &decoder)")
	fmt.Println()
	fmt.Println("  From memory buffer:")
	fmt.Println("      AImageDecoder_createFromBuffer(ptr, length, &decoder)")
	fmt.Println()

	// ---------------------------------------------------------------
	// Decoder methods
	// ---------------------------------------------------------------

	fmt.Println("Decoder methods:")
	fmt.Println()
	fmt.Println("  MinimumStride() int32")
	fmt.Println("      Returns the minimum number of bytes per row in the")
	fmt.Println("      decoded output. This depends on the image width and")
	fmt.Println("      the output pixel format (RGBA_8888 by default, so")
	fmt.Println("      stride = width * 4). Call this after any SetTargetSize")
	fmt.Println("      to get the stride for the configured dimensions.")
	fmt.Println()
	fmt.Println("  SetTargetSize(width, height int32) error")
	fmt.Println("      Request the decoder to scale the output to the given")
	fmt.Println("      dimensions. Both values must be positive and no larger")
	fmt.Println("      than the original image size. Returns ErrInvalidScale")
	fmt.Println("      if the dimensions are invalid. This is optional; by")
	fmt.Println("      default the image is decoded at its original size.")
	fmt.Println()
	fmt.Println("  Decode(pixels unsafe.Pointer, stride int64, size int64) error")
	fmt.Println("      Decode the image into the provided pixel buffer.")
	fmt.Println("      - pixels: pointer to the output buffer")
	fmt.Println("      - stride: bytes per row (>= MinimumStride())")
	fmt.Println("      - size:   total buffer size in bytes (>= stride * height)")
	fmt.Println("      Returns nil on success, or an Error on failure.")
	fmt.Println()
	fmt.Println("  Close() error")
	fmt.Println("      Release the underlying AImageDecoder. Safe to call")
	fmt.Println("      multiple times (idempotent). Must be called when")
	fmt.Println("      decoding is complete to free native resources.")
	fmt.Println()

	// ---------------------------------------------------------------
	// HeaderInfo
	//
	// After creating a decoder, the header can be queried through the
	// NDK functions:
	//
	//   AImageDecoder_getHeaderInfo(decoder) -> *AImageDecoderHeaderInfo
	//   AImageDecoderHeaderInfo_getWidth(info) -> int32
	//   AImageDecoderHeaderInfo_getHeight(info) -> int32
	//   AImageDecoderHeaderInfo_getMimeType(info) -> *byte
	//
	// The HeaderInfo pointer is owned by the decoder; it is valid
	// until the decoder is deleted and must NOT be freed separately.
	//
	// These functions are not yet wrapped in the high-level image
	// package; they are available in the raw NDK layer.
	// ---------------------------------------------------------------

	fmt.Println("HeaderInfo (via NDK layer):")
	fmt.Println()
	fmt.Println("  AImageDecoder_getHeaderInfo(decoder) -> *AImageDecoderHeaderInfo")
	fmt.Println("      Returns a pointer to the header info. Owned by the decoder.")
	fmt.Println()
	fmt.Println("  AImageDecoderHeaderInfo_getWidth(info) -> int32")
	fmt.Println("      Original image width in pixels.")
	fmt.Println()
	fmt.Println("  AImageDecoderHeaderInfo_getHeight(info) -> int32")
	fmt.Println("      Original image height in pixels.")
	fmt.Println()
	fmt.Println("  AImageDecoderHeaderInfo_getMimeType(info) -> *byte")
	fmt.Println("      MIME type string (e.g. \"image/png\"). Null-terminated,")
	fmt.Println("      owned by the decoder.")
	fmt.Println()

	// ---------------------------------------------------------------
	// Complete decode pattern (pseudocode)
	// ---------------------------------------------------------------

	fmt.Println("Complete decode pattern (pseudocode):")
	fmt.Println()
	fmt.Println("  // 1. Open the image file.")
	fmt.Println("  fd, _ := syscall.Open(\"photo.jpg\", syscall.O_RDONLY, 0)")
	fmt.Println("  defer syscall.Close(fd)")
	fmt.Println()
	fmt.Println("  // 2. Create the decoder from the fd.")
	fmt.Println("  //    (factory function not yet in high-level image package)")
	fmt.Println()
	fmt.Println("  // 3. Query image dimensions from the header.")
	fmt.Println("  //    (header info not yet in high-level image package)")
	fmt.Println()
	fmt.Println("  // 4. (Optional) Downscale to half size.")
	fmt.Println("  decoder.SetTargetSize(width/2, height/2)")
	fmt.Println()
	fmt.Println("  // 5. Query the stride for the (possibly resized) output.")
	fmt.Println("  stride := decoder.MinimumStride()")
	fmt.Println()
	fmt.Println("  // 6. Allocate the pixel buffer.")
	fmt.Println("  //    For RGBA_8888 at half size: stride * (height/2) bytes.")
	fmt.Println("  pixels := make([]byte, int64(stride)*int64(height/2))")
	fmt.Println()
	fmt.Println("  // 7. Decode into the buffer.")
	fmt.Println("  err := decoder.Decode(")
	fmt.Println("      unsafe.Pointer(&pixels[0]),")
	fmt.Println("      int64(stride),")
	fmt.Println("      int64(len(pixels)),")
	fmt.Println("  )")
	fmt.Println("  // check err == nil")
	fmt.Println()
	fmt.Println("  // 8. Clean up.")
	fmt.Println("  decoder.Close()")
	fmt.Println()

	// ---------------------------------------------------------------
	// Error constants
	// ---------------------------------------------------------------

	fmt.Println("Error constants:")
	fmt.Println()
	for _, e := range errors {
		fmt.Printf("  %-20s  code=%3d   %s\n",
			e.name, int32(e.value), e.detail)
	}
}
