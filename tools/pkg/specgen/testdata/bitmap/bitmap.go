// Simulates c-for-go output for Android Bitmap (android/bitmap.h, jnigraphics).
// This file is parsed at AST level only; it does not compile.
package bitmap

// Integer typedefs.
type AndroidBitmap_result_t int32
type AndroidBitmap_format_t int32

// Result codes.
const (
	ANDROID_BITMAP_RESULT_SUCCESS           AndroidBitmap_result_t = 0
	ANDROID_BITMAP_RESULT_BAD_PARAMETER     AndroidBitmap_result_t = -1
	ANDROID_BITMAP_RESULT_JNI_EXCEPTION     AndroidBitmap_result_t = -2
	ANDROID_BITMAP_RESULT_ALLOCATION_FAILED AndroidBitmap_result_t = -3
)

// Format constants.
const (
	ANDROID_BITMAP_FORMAT_NONE      AndroidBitmap_format_t = 0
	ANDROID_BITMAP_FORMAT_RGBA_8888 AndroidBitmap_format_t = 1
	ANDROID_BITMAP_FORMAT_RGB_565   AndroidBitmap_format_t = 4
	ANDROID_BITMAP_FORMAT_RGBA_4444 AndroidBitmap_format_t = 7
	ANDROID_BITMAP_FORMAT_A_8       AndroidBitmap_format_t = 8
	ANDROID_BITMAP_FORMAT_RGBA_F16  AndroidBitmap_format_t = 9
)
