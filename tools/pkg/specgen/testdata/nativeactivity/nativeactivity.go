// Simulates c-for-go output for Android NativeActivity.
// This file is parsed at AST level only; it does not compile.
package nativeactivity

import "unsafe"

// Opaque handle types.
type ANativeActivity C.ANativeActivity
type ANativeWindow C.ANativeWindow
type AInputQueue C.AInputQueue
type ARect C.ARect

// Integer typedefs.
type Show_soft_input_flags_t uint32
type Hide_soft_input_flags_t uint32

// ShowSoftInput flags.
const (
	ANATIVEACTIVITY_SHOW_SOFT_INPUT_IMPLICIT Show_soft_input_flags_t = 0x0001
	ANATIVEACTIVITY_SHOW_SOFT_INPUT_FORCED   Show_soft_input_flags_t = 0x0002
)

// HideSoftInput flags.
const (
	ANATIVEACTIVITY_HIDE_SOFT_INPUT_IMPLICIT_ONLY Hide_soft_input_flags_t = 0x0001
	ANATIVEACTIVITY_HIDE_SOFT_INPUT_NOT_ALWAYS     Hide_soft_input_flags_t = 0x0002
)

// --- NativeActivity functions ---
func ANativeActivity_finish(activity *ANativeActivity)                                      {}
func ANativeActivity_setWindowFormat(activity *ANativeActivity, format int32)                {}
func ANativeActivity_setWindowFlags(activity *ANativeActivity, addFlags uint32, removeFlags uint32) {}
func ANativeActivity_showSoftInput(activity *ANativeActivity, flags uint32)                 {}
func ANativeActivity_hideSoftInput(activity *ANativeActivity, flags uint32)                 {}

var _ = unsafe.Pointer(nil)
