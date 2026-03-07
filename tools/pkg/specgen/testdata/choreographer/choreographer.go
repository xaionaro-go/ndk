// Simulates c-for-go output for Android AChoreographer.
// This file is parsed at AST level only; it does not compile.
package choreographer

import "unsafe"

// Opaque handle types.
type AChoreographer C.AChoreographer

// Callback types.
type AChoreographer_frameCallback func(frameTimeNanos int64, data unsafe.Pointer)
type AChoreographer_frameCallback64 func(frameTimeNanos int64, data unsafe.Pointer)

// --- Choreographer functions ---
func AChoreographer_getInstance() *AChoreographer { return nil }
func AChoreographer_postFrameCallback(choreographer *AChoreographer, callback AChoreographer_frameCallback, data unsafe.Pointer) {
}
func AChoreographer_postFrameCallbackDelayed(choreographer *AChoreographer, callback AChoreographer_frameCallback, data unsafe.Pointer, delayMillis int64) {
}
func AChoreographer_postFrameCallback64(choreographer *AChoreographer, callback AChoreographer_frameCallback64, data unsafe.Pointer) {
}
func AChoreographer_postFrameCallbackDelayed64(choreographer *AChoreographer, callback AChoreographer_frameCallback64, data unsafe.Pointer, delayMillis uint32) {
}

var _ = unsafe.Pointer(nil)
