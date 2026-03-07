// Simulates c-for-go output for Android SurfaceTexture.
// This file is parsed at AST level only; it does not compile.
package surfacetexture

import "unsafe"

// Opaque handle types.
type ASurfaceTexture C.ASurfaceTexture
type ANativeWindow C.ANativeWindow

// --- SurfaceTexture lifecycle ---
func ASurfaceTexture_release(surfaceTexture *ASurfaceTexture) {}

// --- SurfaceTexture operations ---
func ASurfaceTexture_acquireANativeWindow(surfaceTexture *ASurfaceTexture) *ANativeWindow { return nil }
func ASurfaceTexture_attachToGLContext(surfaceTexture *ASurfaceTexture, texName uint32) int32 {
	return 0
}
func ASurfaceTexture_detachFromGLContext(surfaceTexture *ASurfaceTexture) int32 { return 0 }
func ASurfaceTexture_updateTexImage(surfaceTexture *ASurfaceTexture) int32      { return 0 }
func ASurfaceTexture_getTransformMatrix(surfaceTexture *ASurfaceTexture, mtx *float32) {
}
func ASurfaceTexture_getTimestamp(surfaceTexture *ASurfaceTexture) int64 { return 0 }
var _ = unsafe.Pointer(nil)
