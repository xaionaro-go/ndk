// Example: offscreen pbuffer surface with an OpenGL ES 2 context.
//
// Demonstrates the full EGL setup for offscreen rendering: display init,
// config selection, pbuffer surface creation, context creation, making
// the context current, verifying it, and orderly cleanup. This pattern
// is the foundation for headless GPU compute or offscreen render-to-texture.
package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/xaionaro-go/ndk/egl"
)

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s (EGL error 0x%04X)\n", msg, egl.GetError())
	os.Exit(1)
}

func main() {
	// --- Display ---
	var defaultDisplay egl.EGLNativeDisplayType
	display := egl.GetDisplay(defaultDisplay)
	if display == nil {
		fatal("GetDisplay returned EGL_NO_DISPLAY")
	}

	var major, minor egl.Int
	if egl.Initialize(display, &major, &minor) == egl.False {
		fatal("Initialize failed")
	}
	defer func() {
		egl.Terminate(display)
		fmt.Println("EGL terminated")
	}()
	fmt.Printf("EGL %d.%d initialized\n", major, minor)

	// --- Bind the OpenGL ES API ---
	if egl.BindAPI(egl.EGLenum(egl.OpenglEsApi)) == egl.False {
		fatal("BindAPI(EGL_OPENGL_ES_API) failed")
	}

	// --- Config selection ---
	// Request an RGBA8 config that supports pbuffer surfaces and ES2 rendering.
	attribs := [...]egl.Int{
		egl.SurfaceType, egl.PbufferBit,
		egl.RenderableType, egl.OpenglEs2Bit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.AlphaSize, 8,
		egl.DepthSize, 0,
		egl.StencilSize, 0,
		egl.None,
	}
	var config egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(display, &attribs[0], &config, 1, &numConfigs) == egl.False {
		fatal("ChooseConfig failed")
	}
	if numConfigs == 0 {
		fatal("ChooseConfig returned 0 matching configs")
	}
	fmt.Printf("found %d matching config(s)\n", numConfigs)

	// --- Pbuffer surface ---
	pbufAttribs := [...]egl.Int{
		egl.Width, 64,
		egl.Height, 64,
		egl.None,
	}
	surface := egl.CreatePbufferSurface(display, config, &pbufAttribs[0])
	if surface == nil {
		fatal("CreatePbufferSurface failed")
	}
	defer func() {
		egl.DestroySurface(display, surface)
		fmt.Println("surface destroyed")
	}()
	fmt.Println("created 64x64 pbuffer surface")

	// --- Context ---
	ctxAttribs := [...]egl.Int{
		egl.ContextClientVersion, 2,
		egl.None,
	}
	ctx := egl.CreateContext(display, config, nil, &ctxAttribs[0])
	if ctx == nil {
		fatal("CreateContext failed")
	}
	defer func() {
		// Unbind before destroying.
		egl.MakeCurrent(display, nil, nil, nil)
		egl.DestroyContext(display, ctx)
		fmt.Println("context destroyed")
	}()
	fmt.Println("created OpenGL ES 2 context")

	// --- Make current ---
	if egl.MakeCurrent(display, surface, surface, ctx) == egl.False {
		fatal("MakeCurrent failed")
	}
	fmt.Println("context is now current")

	// --- Verify ---
	cur := egl.GetCurrentContext()
	if cur == nil {
		fatal("GetCurrentContext returned EGL_NO_CONTEXT")
	}
	if unsafe.Pointer(cur) != unsafe.Pointer(ctx) {
		fatal("GetCurrentContext does not match the context we created")
	}
	fmt.Println("verified: GetCurrentContext matches our context")

	curDisplay := egl.GetCurrentDisplay()
	if curDisplay == nil {
		fatal("GetCurrentDisplay returned EGL_NO_DISPLAY")
	}
	fmt.Println("verified: GetCurrentDisplay is valid")

	// Query the actual pbuffer dimensions.
	var w, h egl.Int
	egl.QuerySurface(display, surface, egl.Width, &w)
	egl.QuerySurface(display, surface, egl.Height, &h)
	fmt.Printf("pbuffer surface dimensions: %dx%d\n", w, h)

	// At this point OpenGL ES commands can be issued.
	fmt.Println("ready for OpenGL ES rendering (not demonstrated here)")
}
