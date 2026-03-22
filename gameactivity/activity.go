//go:build ignore

package gameactivity

/*
#cgo LDFLAGS: -landroid -llog

#include "include/game-activity/GameActivity.h"
#include <stdlib.h>
*/
import "C"

import "unsafe"

// Activity wraps a GameActivity pointer.
type Activity struct {
	ptr *C.GameActivity
}

// NewActivityFromPointer wraps a raw GameActivity pointer.
func NewActivityFromPointer(ptr unsafe.Pointer) *Activity {
	return &Activity{ptr: (*C.GameActivity)(ptr)}
}

// Pointer returns the underlying pointer as unsafe.Pointer.
func (a *Activity) Pointer() unsafe.Pointer {
	return unsafe.Pointer(a.ptr)
}

// UintPtr returns the underlying pointer as a uintptr.
func (a *Activity) UintPtr() uintptr {
	return uintptr(unsafe.Pointer(a.ptr))
}

// NewActivityFromUintPtr wraps a uintptr as an Activity.
func NewActivityFromUintPtr(ptr uintptr) *Activity {
	return &Activity{ptr: (*C.GameActivity)(unsafe.Pointer(ptr))}
}

// VM returns the Java VM handle.
func (a *Activity) VM() unsafe.Pointer {
	return unsafe.Pointer(a.ptr.vm)
}

// VMUintPtr returns the Java VM handle as uintptr.
func (a *Activity) VMUintPtr() uintptr {
	return uintptr(unsafe.Pointer(a.ptr.vm))
}

// Env returns the JNI environment for the main thread.
// This must only be used from the main thread.
func (a *Activity) Env() unsafe.Pointer {
	return unsafe.Pointer(a.ptr.env)
}

// EnvUintPtr returns the JNI environment as uintptr.
// This must only be used from the main thread.
func (a *Activity) EnvUintPtr() uintptr {
	return uintptr(unsafe.Pointer(a.ptr.env))
}

// JavaGameActivity returns the Java GameActivity jobject handle.
func (a *Activity) JavaGameActivity() uintptr {
	return uintptr(a.ptr.javaGameActivity)
}

// SDKVersion returns the platform SDK version code.
func (a *Activity) SDKVersion() int32 {
	return int32(a.ptr.sdkVersion)
}

// InternalDataPath returns the path to the application's internal data directory.
func (a *Activity) InternalDataPath() string {
	return C.GoString(a.ptr.internalDataPath)
}

// ExternalDataPath returns the path to the application's external data directory.
func (a *Activity) ExternalDataPath() string {
	return C.GoString(a.ptr.externalDataPath)
}

// OBBPath returns the path to the application's OBB directory.
func (a *Activity) OBBPath() string {
	return C.GoString(a.ptr.obbPath)
}

// AssetManager returns the AAssetManager pointer.
func (a *Activity) AssetManager() unsafe.Pointer {
	return unsafe.Pointer(a.ptr.assetManager)
}

// Instance returns the application's custom instance data.
func (a *Activity) Instance() unsafe.Pointer {
	return a.ptr.instance
}

// SetInstance stores custom instance data on the activity.
func (a *Activity) SetInstance(data unsafe.Pointer) {
	a.ptr.instance = data
}

// Finish requests that the GameActivity be finished.
func (a *Activity) Finish() {
	C.GameActivity_finish(a.ptr)
}

// SetWindowFlags sets the activity's window flags.
func (a *Activity) SetWindowFlags(addFlags, removeFlags uint32) {
	C.GameActivity_setWindowFlags(a.ptr, C.uint32_t(addFlags), C.uint32_t(removeFlags))
}

// ShowSoftInput shows the software keyboard.
func (a *Activity) ShowSoftInput(flags ShowSoftInputFlags) {
	C.GameActivity_showSoftInput(a.ptr, C.uint32_t(flags))
}

// HideSoftInput hides the software keyboard.
func (a *Activity) HideSoftInput(flags HideSoftInputFlags) {
	C.GameActivity_hideSoftInput(a.ptr, C.uint32_t(flags))
}

// SetImeEditorInfo configures the IME editor info.
func (a *Activity) SetImeEditorInfo(
	inputType int,
	actionID int,
	imeOptions int,
) {
	C.GameActivity_setImeEditorInfo(
		a.ptr,
		C.int(inputType),
		C.int(actionID),
		C.int(imeOptions),
	)
}

// SetTextInputState sets the current text input state.
func (a *Activity) SetTextInputState(state TextInputState) {
	var cState C.GameTextInputState
	cText := C.CString(state.Text)
	defer C.free(unsafe.Pointer(cText))
	cState.text_UTF8 = cText
	cState.text_length = C.int32_t(len(state.Text))
	cState.selection.start = C.int32_t(state.SelectionStart)
	cState.selection.end = C.int32_t(state.SelectionEnd)
	cState.composingRegion.start = C.int32_t(state.ComposingStart)
	cState.composingRegion.end = C.int32_t(state.ComposingEnd)
	C.GameActivity_setTextInputState(a.ptr, &cState)
}

// GetWindowInsets retrieves the window insets for the given type.
// Returns left, top, right, bottom.
func (a *Activity) GetWindowInsets(insetsType InsetsType) (left, top, right, bottom int32) {
	var rect C.ARect
	C.GameActivity_getWindowInsets(a.ptr, C.GameCommonInsetsType(insetsType), &rect)
	return int32(rect.left), int32(rect.top), int32(rect.right), int32(rect.bottom)
}
