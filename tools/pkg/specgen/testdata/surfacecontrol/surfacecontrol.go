// Simulates c-for-go output for Android SurfaceControl.
// This file is parsed at AST level only; it does not compile.
package surfacecontrol

import "unsafe"

// Opaque handle types.
type ASurfaceControl C.ASurfaceControl
type ASurfaceTransaction C.ASurfaceTransaction
type ANativeWindow C.ANativeWindow
type AHardwareBuffer C.AHardwareBuffer

// Integer typedefs.
type ASurfaceTransaction_visibility_t int32
type ASurfaceTransaction_transparency_t int32

// Visibility enum.
const (
	ASURFACE_TRANSACTION_VISIBILITY_SHOW ASurfaceTransaction_visibility_t = 0
	ASURFACE_TRANSACTION_VISIBILITY_HIDE ASurfaceTransaction_visibility_t = 1
)

// Transparency enum.
const (
	ASURFACE_TRANSACTION_TRANSPARENCY_TRANSPARENT ASurfaceTransaction_transparency_t = 0
	ASURFACE_TRANSACTION_TRANSPARENCY_TRANSLUCENT ASurfaceTransaction_transparency_t = 1
	ASURFACE_TRANSACTION_TRANSPARENCY_OPAQUE      ASurfaceTransaction_transparency_t = 2
)

// --- SurfaceControl lifecycle ---
func ASurfaceControl_createFromWindow(parent *ANativeWindow, debugName *byte) *ASurfaceControl { return nil }
func ASurfaceControl_create(parent *ASurfaceControl, debugName *byte) *ASurfaceControl         { return nil }
func ASurfaceControl_acquire(surfaceControl *ASurfaceControl)                                   {}
func ASurfaceControl_release(surfaceControl *ASurfaceControl)                                   {}

// --- SurfaceTransaction lifecycle ---
func ASurfaceTransaction_create() *ASurfaceTransaction        { return nil }
func ASurfaceTransaction_delete(transaction *ASurfaceTransaction) {}
func ASurfaceTransaction_apply(transaction *ASurfaceTransaction)  {}

// --- Transaction operations ---
func ASurfaceTransaction_setVisibility(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, visibility ASurfaceTransaction_visibility_t) {}
func ASurfaceTransaction_setZOrder(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, zOrder int32) {}
func ASurfaceTransaction_setBuffer(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, buffer *AHardwareBuffer, fenceFd int32) {}
// setGeometry removed: uses C++ references (const ARect&) incompatible with CGo.
func ASurfaceTransaction_setBufferTransparency(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, transparency ASurfaceTransaction_transparency_t) {}
func ASurfaceTransaction_setDamageRegion(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, rects unsafe.Pointer, count uint32) {}
func ASurfaceTransaction_setPosition(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, x int32, y int32) {}
func ASurfaceTransaction_setBufferAlpha(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, alpha float32) {}
// setCrop removed: uses C++ references (const ARect&) incompatible with CGo.
func ASurfaceTransaction_setScale(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, xScale float32, yScale float32) {}
func ASurfaceTransaction_setColor(transaction *ASurfaceTransaction, surfaceControl *ASurfaceControl, r float32, g float32, b float32, a float32, dataspace int32) {}

// --- OnComplete callback ---
type ASurfaceTransaction_OnComplete func(context unsafe.Pointer, stats unsafe.Pointer)

func ASurfaceTransaction_setOnComplete(transaction *ASurfaceTransaction, context unsafe.Pointer, callback ASurfaceTransaction_OnComplete) {}

var _ = unsafe.Pointer(nil)
