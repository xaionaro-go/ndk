// Simulates c-for-go output for Android ImageDecoder (android/imagedecoder.h).
// This file is parsed at AST level only; it does not compile.
package imagedecoder

import "unsafe"

// Opaque handle types.
type AImageDecoder C.AImageDecoder
type AImageDecoderHeaderInfo C.AImageDecoderHeaderInfo

// Integer typedefs.
type ImageDecoder_result_t int32

// Result codes.
const (
	ANDROID_IMAGE_DECODER_SUCCESS              ImageDecoder_result_t = 0
	ANDROID_IMAGE_DECODER_INCOMPLETE           ImageDecoder_result_t = -1
	ANDROID_IMAGE_DECODER_ERROR                ImageDecoder_result_t = -2
	ANDROID_IMAGE_DECODER_INVALID_CONVERSION   ImageDecoder_result_t = -3
	ANDROID_IMAGE_DECODER_INVALID_SCALE        ImageDecoder_result_t = -4
	ANDROID_IMAGE_DECODER_BAD_PARAMETER        ImageDecoder_result_t = -5
	ANDROID_IMAGE_DECODER_INVALID_INPUT        ImageDecoder_result_t = -6
	ANDROID_IMAGE_DECODER_SEEK_ERROR           ImageDecoder_result_t = -7
	ANDROID_IMAGE_DECODER_INTERNAL_ERROR       ImageDecoder_result_t = -8
	ANDROID_IMAGE_DECODER_UNSUPPORTED_FORMAT   ImageDecoder_result_t = -9
	ANDROID_IMAGE_DECODER_FINISHED             ImageDecoder_result_t = -10
	ANDROID_IMAGE_DECODER_INVALID_STATE        ImageDecoder_result_t = -11
)

// --- Decoder functions ---
func AImageDecoder_createFromFd(fd int32, outDecoder **AImageDecoder) int32                             { return 0 }
func AImageDecoder_createFromBuffer(buffer unsafe.Pointer, length int64, outDecoder **AImageDecoder) int32 { return 0 }
func AImageDecoder_delete(decoder *AImageDecoder)                                                       {}
func AImageDecoder_setTargetSize(decoder *AImageDecoder, width int32, height int32) int32               { return 0 }
func AImageDecoder_getHeaderInfo(decoder *AImageDecoder) *AImageDecoderHeaderInfo                       { return nil }
func AImageDecoder_getMinimumStride(decoder *AImageDecoder) int32                                       { return 0 }
func AImageDecoder_decodeImage(decoder *AImageDecoder, pixels unsafe.Pointer, stride int64, size int64) int32 { return 0 }

// --- HeaderInfo functions ---
func AImageDecoderHeaderInfo_getWidth(info *AImageDecoderHeaderInfo) int32    { return 0 }
func AImageDecoderHeaderInfo_getHeight(info *AImageDecoderHeaderInfo) int32   { return 0 }
func AImageDecoderHeaderInfo_getMimeType(info *AImageDecoderHeaderInfo) *byte { return nil }

var _ = unsafe.Pointer(nil)
