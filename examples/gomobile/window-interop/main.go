// Passing ANativeWindow handles between Java and Go via gomobile bind.
//
// When integrating NDK rendering with a Java/Kotlin Android app,
// the Java side holds the Surface and must pass the native window
// handle to Go. Because gomobile bind does not support unsafe.Pointer
// or uintptr, the handle is transported as int64 (Java long).
//
// # Obtaining ANativeWindow from Java
//
// Java/Kotlin obtains an ANativeWindow* via JNI from an android.view.Surface:
//
//	// Java (JNI helper)
//	public static native long surfaceToNativeWindow(Surface surface);
//
//	// C (JNI implementation)
//	JNIEXPORT jlong JNICALL Java_com_example_NativeLib_surfaceToNativeWindow(
//	    JNIEnv *env, jclass cls, jobject surface) {
//	    return (jlong)ANativeWindow_fromSurface(env, surface);
//	}
//
// # Go library package (bindable)
//
//	package renderer
//
//	import (
//	    "github.com/AndroidGoLab/ndk/egl"
//	    "github.com/AndroidGoLab/ndk/window"
//	)
//
//	// Renderer manages an EGL context attached to a native window.
//	type Renderer struct {
//	    win     *window.Window
//	    display egl.EGLDisplay
//	    surface egl.EGLSurface
//	    ctx     egl.EGLContext
//	}
//
//	// NewRenderer creates a renderer from a native window handle.
//	// The windowHandle is an ANativeWindow* cast to int64 by the Java caller.
//	func NewRenderer(windowHandle int64) (*Renderer, error) {
//	    win := window.NewWindowFromUintPtr(uintptr(windowHandle))
//	    // ... set up EGL display, surface, context using win ...
//	    return &Renderer{win: win}, nil
//	}
//
//	// WindowWidth returns the current window width in pixels.
//	func (r *Renderer) WindowWidth() int32 {
//	    // window.Width() returns error; unwrap for gomobile.
//	    // In production code, handle the error appropriately.
//	    return 0 // placeholder: win.GetWidth() returns error in current API
//	}
//
//	// RenderFrame renders a single frame. Called from Java's render loop.
//	func (r *Renderer) RenderFrame() error {
//	    // ... OpenGL ES rendering calls ...
//	    return nil
//	}
//
//	// Destroy releases all resources. Must be called when the surface
//	// is destroyed to avoid leaking native handles.
//	func (r *Renderer) Destroy() {
//	    // ... tear down EGL ...
//	    // Do NOT close the window — Java owns the Surface lifecycle.
//	}
//
//	// WindowHandle returns the raw ANativeWindow* as int64 so Java
//	// can pass it to another native library.
//	func (r *Renderer) WindowHandle() int64 {
//	    return int64(r.win.UintPtr())
//	}
//
// # Java usage
//
//	import renderer.Renderer;
//
//	public class MySurfaceView extends SurfaceView implements SurfaceHolder.Callback {
//	    private Renderer renderer;
//
//	    @Override
//	    public void surfaceCreated(SurfaceHolder holder) {
//	        long windowPtr = surfaceToNativeWindow(holder.getSurface());
//	        try {
//	            renderer = Renderer.newRenderer(windowPtr);
//	        } catch (Exception e) {
//	            Log.e("Renderer", "init failed", e);
//	        }
//	    }
//
//	    @Override
//	    public void surfaceDestroyed(SurfaceHolder holder) {
//	        if (renderer != null) {
//	            renderer.destroy();
//	            renderer = null;
//	        }
//	    }
//	}
//
// # EGL handle round-tripping
//
// EGL types (EGLDisplay, EGLContext, EGLSurface) are unsafe.Pointer aliases.
// Use the egl package's conversion functions:
//
//	// Go -> Java: extract EGL handle as int64
//	func (r *Renderer) EGLDisplayHandle() int64 {
//	    return int64(egl.EGLDisplayToUintPtr(r.display))
//	}
//
//	// Java -> Go: reconstruct EGL handle from int64
//	func SetEGLDisplay(handle int64) {
//	    display := egl.EGLDisplayFromUintPtr(uintptr(handle))
//	    // use display...
//	}
//
// # Ownership rules
//
//   - The Java side owns Surface lifecycle (create/destroy).
//   - Go must NOT call window.Close() if Java owns the Surface.
//   - Go MUST release its own resources (EGL context, buffers) before
//     surfaceDestroyed returns.
//   - Native handles passed as int64 are raw pointers — no reference
//     counting. The caller must ensure the handle remains valid.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/AndroidGoLab/ndk/egl"
	"github.com/AndroidGoLab/ndk/window"
)

func main() {
	fmt.Println("=== gomobile bind: window interop ===")
	fmt.Println()

	// Demonstrate the handle conversion chain for window handles.

	fmt.Println("Handle flow: Java Surface -> JNI ANativeWindow* -> Go int64 -> window.Window")
	fmt.Println()

	// Step 1: Java obtains ANativeWindow from Surface via JNI.
	fmt.Println("Step 1: Java obtains ANativeWindow* via JNI")
	fmt.Println("  Java: long ptr = surfaceToNativeWindow(holder.getSurface());")
	fmt.Println("  C:    return (jlong)ANativeWindow_fromSurface(env, surface);")
	fmt.Println()

	// Step 2: Java passes handle to Go as int64.
	fmt.Println("Step 2: Java passes handle to Go via gomobile bind")
	fmt.Println("  Java: Renderer renderer = Renderer.newRenderer(ptr);")
	fmt.Println("  Go:   func NewRenderer(windowHandle int64) (*Renderer, error)")
	fmt.Println()

	// Step 3: Go wraps int64 as window.Window.
	fmt.Println("Step 3: Go converts int64 -> window.Window")
	fmt.Println("  Go:   win := window.NewWindowFromUintPtr(uintptr(windowHandle))")
	fmt.Println()

	// Step 4: Go returns handle to Java as int64.
	fmt.Println("Step 4: Go returns handle to Java as int64")
	fmt.Println("  Go:   return int64(win.UintPtr())")
	fmt.Println("  Java: long handle = renderer.windowHandle();")
	fmt.Println()

	// Show the window.Window and egl interop API.
	fmt.Println("window.Window interop API:")
	fmt.Printf("  %-35s  %s\n", "window.NewWindowFromUintPtr(uintptr)", "wrap uintptr")
	fmt.Printf("  %-35s  %s\n", "window.NewWindowFromPointer(ptr)", "wrap unsafe.Pointer")
	fmt.Printf("  %-35s  %s\n", "win.UintPtr()", "extract as uintptr")
	fmt.Printf("  %-35s  %s\n", "win.Pointer()", "extract as unsafe.Pointer")
	fmt.Println()

	// Demonstrate that window.Window and egl.ANativeWindow are separate
	// types wrapping the same C type, and handles can be shared.
	fmt.Println("Cross-package handle sharing:")
	fmt.Println("  // window.Window and egl.ANativeWindow both wrap ANativeWindow*.")
	fmt.Println("  // Transfer via uintptr:")
	fmt.Println("  eglWin := egl.NewANativeWindowFromUintPtr(win.UintPtr())")
	fmt.Println()

	// Show EGL type conversion API.
	fmt.Println("EGL type interop API:")
	eglTypes := []string{
		"EGLDisplay", "EGLContext", "EGLSurface",
		"EGLConfig", "EGLImage", "EGLSync", "EGLClientBuffer",
	}
	for _, t := range eglTypes {
		fmt.Printf("  egl.%sToUintPtr(h) / egl.%sFromUintPtr(p)\n", t, t)
	}
	fmt.Println()

	// Show gomobile bind pattern for EGL.
	fmt.Println("gomobile bind pattern for EGL handles:")
	fmt.Println("  // Go: export as int64 for Java")
	fmt.Println("  func (r *Renderer) EGLDisplayHandle() int64 {")
	fmt.Println("      return int64(egl.EGLDisplayToUintPtr(r.display))")
	fmt.Println("  }")
	fmt.Println()

	// Suppress unused import errors by referencing the packages.
	_ = window.NewWindowFromUintPtr
	_ = egl.EGLDisplayToUintPtr

	fmt.Println("window-interop example complete")
}
