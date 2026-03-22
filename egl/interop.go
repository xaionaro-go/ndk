package egl

import "unsafe"

// EGLDisplayToUintPtr converts an EGLDisplay to uintptr for framework interop.
func EGLDisplayToUintPtr(d EGLDisplay) uintptr { return uintptr(unsafe.Pointer(d)) }

// EGLDisplayFromUintPtr converts a uintptr to EGLDisplay.
func EGLDisplayFromUintPtr(ptr uintptr) EGLDisplay { return EGLDisplay(unsafe.Pointer(ptr)) }

// EGLContextToUintPtr converts an EGLContext to uintptr for framework interop.
func EGLContextToUintPtr(c EGLContext) uintptr { return uintptr(unsafe.Pointer(c)) }

// EGLContextFromUintPtr converts a uintptr to EGLContext.
func EGLContextFromUintPtr(ptr uintptr) EGLContext { return EGLContext(unsafe.Pointer(ptr)) }

// EGLSurfaceToUintPtr converts an EGLSurface to uintptr for framework interop.
func EGLSurfaceToUintPtr(s EGLSurface) uintptr { return uintptr(unsafe.Pointer(s)) }

// EGLSurfaceFromUintPtr converts a uintptr to EGLSurface.
func EGLSurfaceFromUintPtr(ptr uintptr) EGLSurface { return EGLSurface(unsafe.Pointer(ptr)) }

// EGLConfigToUintPtr converts an EGLConfig to uintptr for framework interop.
func EGLConfigToUintPtr(c EGLConfig) uintptr { return uintptr(unsafe.Pointer(c)) }

// EGLConfigFromUintPtr converts a uintptr to EGLConfig.
func EGLConfigFromUintPtr(ptr uintptr) EGLConfig { return EGLConfig(unsafe.Pointer(ptr)) }

// EGLImageToUintPtr converts an EGLImage to uintptr for framework interop.
func EGLImageToUintPtr(i EGLImage) uintptr { return uintptr(unsafe.Pointer(i)) }

// EGLImageFromUintPtr converts a uintptr to EGLImage.
func EGLImageFromUintPtr(ptr uintptr) EGLImage { return EGLImage(unsafe.Pointer(ptr)) }

// EGLSyncToUintPtr converts an EGLSync to uintptr for framework interop.
func EGLSyncToUintPtr(s EGLSync) uintptr { return uintptr(unsafe.Pointer(s)) }

// EGLSyncFromUintPtr converts a uintptr to EGLSync.
func EGLSyncFromUintPtr(ptr uintptr) EGLSync { return EGLSync(unsafe.Pointer(ptr)) }

// EGLClientBufferToUintPtr converts an EGLClientBuffer to uintptr for framework interop.
func EGLClientBufferToUintPtr(b EGLClientBuffer) uintptr { return uintptr(unsafe.Pointer(b)) }

// EGLClientBufferFromUintPtr converts a uintptr to EGLClientBuffer.
func EGLClientBufferFromUintPtr(ptr uintptr) EGLClientBuffer {
	return EGLClientBuffer(unsafe.Pointer(ptr))
}
