// Simulates c-for-go output for Android NativeWindow.
// This file is parsed at AST level only; it does not compile.
package nativewindow

import "unsafe"

// Opaque handle types.
type ANativeWindow C.ANativeWindow
type ANativeWindow_Buffer C.ANativeWindow_Buffer
type ARect C.ARect

// Integer typedefs.
type Window_format_t int32

// Window format enum.
const (
	WINDOW_FORMAT_RGBA_8888 Window_format_t = 1
	WINDOW_FORMAT_RGBX_8888 Window_format_t = 2
	WINDOW_FORMAT_RGB_565   Window_format_t = 4
)

// --- Window functions ---
func ANativeWindow_acquire(window *ANativeWindow)                           {}
func ANativeWindow_release(window *ANativeWindow)                           {}
func ANativeWindow_getWidth(window *ANativeWindow) int32                    { return 0 }
func ANativeWindow_getHeight(window *ANativeWindow) int32                   { return 0 }
func ANativeWindow_getFormat(window *ANativeWindow) int32                   { return 0 }
func ANativeWindow_setBuffersGeometry(window *ANativeWindow, width int32, height int32, format int32) int32 { return 0 }
func ANativeWindow_lock(window *ANativeWindow, outBuffer *ANativeWindow_Buffer, inOutDirtyBounds unsafe.Pointer) int32 { return 0 }
func ANativeWindow_unlockAndPost(window *ANativeWindow) int32               { return 0 }

var _ = unsafe.Pointer(nil)
