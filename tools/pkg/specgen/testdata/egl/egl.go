// Simulates c-for-go output for EGL (EGL/egl.h and EGL/eglext.h).
// This file is parsed at AST level only; it does not compile.
package egl

import "unsafe"

// Opaque handle types — pointer-sized handles, not C struct pointers.
type EGLDisplay unsafe.Pointer
type EGLSurface unsafe.Pointer
type EGLContext unsafe.Pointer
type EGLConfig unsafe.Pointer
type EGLNativeWindowType unsafe.Pointer
type EGLNativeDisplayType unsafe.Pointer

// Integer typedefs.
type EGLint int32
type EGLBoolean uint32

// Error codes.
const (
	EGL_SUCCESS             EGLint = 0x3000
	EGL_NOT_INITIALIZED     EGLint = 0x3001
	EGL_BAD_ACCESS          EGLint = 0x3002
	EGL_BAD_ALLOC           EGLint = 0x3003
	EGL_BAD_ATTRIBUTE       EGLint = 0x3004
	EGL_BAD_CONFIG          EGLint = 0x3005
	EGL_BAD_CONTEXT         EGLint = 0x3006
	EGL_BAD_CURRENT_SURFACE EGLint = 0x3007
	EGL_BAD_DISPLAY         EGLint = 0x3008
	EGL_BAD_MATCH           EGLint = 0x3009
	EGL_BAD_NATIVE_PIXMAP   EGLint = 0x300A
	EGL_BAD_NATIVE_WINDOW   EGLint = 0x300B
	EGL_BAD_PARAMETER       EGLint = 0x300C
	EGL_BAD_SURFACE         EGLint = 0x300D
	EGL_CONTEXT_LOST        EGLint = 0x300E
)

// Config attributes.
const (
	EGL_BUFFER_SIZE         EGLint = 0x3020
	EGL_ALPHA_SIZE          EGLint = 0x3021
	EGL_BLUE_SIZE           EGLint = 0x3022
	EGL_GREEN_SIZE          EGLint = 0x3023
	EGL_RED_SIZE            EGLint = 0x3024
	EGL_DEPTH_SIZE          EGLint = 0x3025
	EGL_STENCIL_SIZE        EGLint = 0x3026
	EGL_CONFIG_CAVEAT       EGLint = 0x3027
	EGL_CONFIG_ID           EGLint = 0x3028
	EGL_MAX_PBUFFER_HEIGHT  EGLint = 0x302A
	EGL_MAX_PBUFFER_PIXELS  EGLint = 0x302B
	EGL_MAX_PBUFFER_WIDTH   EGLint = 0x302C
	EGL_NATIVE_RENDERABLE   EGLint = 0x302D
	EGL_NATIVE_VISUAL_ID    EGLint = 0x302E
	EGL_NATIVE_VISUAL_TYPE  EGLint = 0x302F
	EGL_SAMPLES             EGLint = 0x3031
	EGL_SAMPLE_BUFFERS      EGLint = 0x3032
	EGL_SURFACE_TYPE        EGLint = 0x3033
	EGL_RENDERABLE_TYPE     EGLint = 0x3040
	EGL_CONFORMANT          EGLint = 0x3042
	EGL_NONE                EGLint = 0x3038
)

// Renderable type bits.
const (
	EGL_OPENGL_ES2_BIT EGLint = 0x0004
	EGL_OPENGL_ES3_BIT EGLint = 0x0040
)

// Context attributes.
const (
	EGL_CONTEXT_CLIENT_VERSION EGLint = 0x3098
	EGL_CONTEXT_MAJOR_VERSION  EGLint = 0x3098
	EGL_CONTEXT_MINOR_VERSION  EGLint = 0x30FB
)

// Special values.
const (
	EGL_DEFAULT_DISPLAY EGLint = 0
	EGL_NO_DISPLAY      EGLint = 0
	EGL_NO_SURFACE      EGLint = 0
	EGL_NO_CONTEXT      EGLint = 0
	EGL_TRUE            EGLint = 1
	EGL_FALSE           EGLint = 0
)

// Surface type bits.
const (
	EGL_WINDOW_BIT EGLint = 0x0004
	EGL_PBUFFER_BIT EGLint = 0x0001
	EGL_PIXMAP_BIT  EGLint = 0x0002
)

// Surface query attributes.
const (
	EGL_WIDTH  EGLint = 0x3057
	EGL_HEIGHT EGLint = 0x3056
)

// --- Display lifecycle ---
func EglGetDisplay(displayId EGLNativeDisplayType) EGLDisplay    { return unsafe.Pointer(nil) }
func EglInitialize(display EGLDisplay, major *EGLint, minor *EGLint) EGLBoolean { return 0 }
func EglTerminate(display EGLDisplay) EGLBoolean                { return 0 }
func EglGetError() EGLint                                       { return 0 }

// --- Config functions ---
func EglChooseConfig(display EGLDisplay, attribList *EGLint, configs *EGLConfig, configSize EGLint, numConfig *EGLint) EGLBoolean {
	return 0
}
func EglGetConfigAttrib(display EGLDisplay, config EGLConfig, attribute EGLint, value *EGLint) EGLBoolean {
	return 0
}

// --- Surface functions ---
func EglCreateWindowSurface(display EGLDisplay, config EGLConfig, win EGLNativeWindowType, attribList *EGLint) EGLSurface {
	return unsafe.Pointer(nil)
}
func EglCreatePbufferSurface(display EGLDisplay, config EGLConfig, attribList *EGLint) EGLSurface {
	return unsafe.Pointer(nil)
}
func EglDestroySurface(display EGLDisplay, surface EGLSurface) EGLBoolean { return 0 }
func EglQuerySurface(display EGLDisplay, surface EGLSurface, attribute EGLint, value *EGLint) EGLBoolean {
	return 0
}

// --- Context functions ---
func EglCreateContext(display EGLDisplay, config EGLConfig, shareContext EGLContext, attribList *EGLint) EGLContext {
	return unsafe.Pointer(nil)
}
func EglDestroyContext(display EGLDisplay, ctx EGLContext) EGLBoolean { return 0 }
func EglMakeCurrent(display EGLDisplay, draw EGLSurface, read EGLSurface, ctx EGLContext) EGLBoolean {
	return 0
}

// --- Swap functions ---
func EglSwapBuffers(display EGLDisplay, surface EGLSurface) EGLBoolean { return 0 }
func EglSwapInterval(display EGLDisplay, interval EGLint) EGLBoolean   { return 0 }

// --- Query current state ---
func EglGetCurrentDisplay() EGLDisplay               { return unsafe.Pointer(nil) }
func EglGetCurrentSurface(readdraw EGLint) EGLSurface { return unsafe.Pointer(nil) }
func EglGetCurrentContext() EGLContext                { return unsafe.Pointer(nil) }

// --- String query ---
func EglQueryString(display EGLDisplay, name EGLint) *byte { return nil }

// --- Thread and sync ---
func EglReleaseThread() EGLBoolean          { return 0 }
func EglWaitGL() EGLBoolean                 { return 0 }
func EglWaitNative(engine EGLint) EGLBoolean { return 0 }

// --- API binding ---
func EglBindAPI(api uint32) EGLBoolean { return 0 }
func EglQueryAPI() uint32              { return 0 }

var _ = unsafe.Pointer(nil)
