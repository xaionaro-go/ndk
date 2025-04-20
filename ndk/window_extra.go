//go:build ndk_window_extra
// +build ndk_window_extra

package ndk

import (
	"syscall"
)

/*
#include <android/window.h>
#include <android/native_window.h>
#include <android/native_window_jni.h>

#cgo LDFLAGS: -lnativewindow
*/
import "C"

// See https://developer.android.com/ndk/reference/group/a-native-window#anativewindow_setbufferstransform
func (w *Window) SetBuffersTransform(transform WindowTransform) error {
	rc := int32(C.ANativeWindow_setBuffersTransform(w.cptr(), C.int32_t(transform)))
	if rc != 0 {
		return syscall.Errno(-rc)
	}
	return nil
}
