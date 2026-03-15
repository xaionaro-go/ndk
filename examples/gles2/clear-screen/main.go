// Simplest GLES2 example: clear the framebuffer to a known color and verify.
//
// Demonstrates:
//   - EGL offscreen context setup (pbuffer surface)
//   - Setting the clear color with ClearColor
//   - Clearing the color buffer with Clear
//   - Reading back pixels with ReadPixels to confirm the result
//
// This is the "hello world" of OpenGL ES -- no shaders, no geometry, just
// proving that the GPU pipeline is alive by writing and reading a pixel.
package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/AndroidGoLab/ndk/egl"
	"github.com/AndroidGoLab/ndk/gles2"
)

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "fatal: %s (EGL error 0x%04X)\n", msg, egl.GetError())
	os.Exit(1)
}

func checkGL(context string) {
	if err := gles2.GetError(); err != gles2.NoError {
		fmt.Fprintf(os.Stderr, "fatal: GL error 0x%04X at %s\n", int32(err), context)
		os.Exit(1)
	}
}

func main() {
	// --- EGL initialization ---
	var defaultDisplay egl.EGLNativeDisplayType
	display := egl.GetDisplay(defaultDisplay)
	if display == nil {
		fatal("GetDisplay returned EGL_NO_DISPLAY")
	}

	var major, minor egl.Int
	if egl.Initialize(display, &major, &minor) == 0 {
		fatal("Initialize failed")
	}
	defer egl.Terminate(display)
	fmt.Printf("EGL %d.%d\n", major, minor)

	if egl.BindAPI(egl.EGLenum(egl.OpenglEsApi)) == 0 {
		fatal("BindAPI failed")
	}

	// Request an RGBA8 pbuffer config with ES2 support.
	configAttribs := [...]egl.Int{
		egl.SurfaceType, egl.PbufferBit,
		egl.RenderableType, egl.OpenglEs2Bit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.AlphaSize, 8,
		egl.None,
	}
	var config egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(display, &configAttribs[0], &config, 1, &numConfigs) == 0 || numConfigs == 0 {
		fatal("ChooseConfig failed")
	}

	// Create a small offscreen surface.
	const width, height = 4, 4
	pbufAttribs := [...]egl.Int{
		egl.Width, width,
		egl.Height, height,
		egl.None,
	}
	surface := egl.CreatePbufferSurface(display, config, &pbufAttribs[0])
	if surface == nil {
		fatal("CreatePbufferSurface failed")
	}
	defer egl.DestroySurface(display, surface)

	ctxAttribs := [...]egl.Int{
		egl.ContextClientVersion, 2,
		egl.None,
	}
	ctx := egl.CreateContext(display, config, nil, &ctxAttribs[0])
	if ctx == nil {
		fatal("CreateContext failed")
	}
	defer func() {
		egl.MakeCurrent(display, nil, nil, nil)
		egl.DestroyContext(display, ctx)
	}()

	if egl.MakeCurrent(display, surface, surface, ctx) == 0 {
		fatal("MakeCurrent failed")
	}
	fmt.Println("GL context is current")

	// --- Clear to cornflower blue (R=100, G=149, B=237, A=255) ---
	gles2.Viewport(0, 0, width, height)
	gles2.ClearColor(100.0/255.0, 149.0/255.0, 237.0/255.0, 1.0)
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))
	checkGL("Clear")

	// --- Read back pixel (0,0) and verify ---
	var pixel [4]byte
	gles2.ReadPixels(0, 0, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&pixel[0]))
	checkGL("ReadPixels")

	fmt.Printf("pixel at (0,0): R=%d G=%d B=%d A=%d\n", pixel[0], pixel[1], pixel[2], pixel[3])

	// Allow small rounding differences from float-to-unorm conversion.
	if pixel[0] != 100 || pixel[1] != 149 || pixel[2] != 237 || pixel[3] != 255 {
		fmt.Fprintf(os.Stderr, "warning: pixel does not match expected cornflower blue (100,149,237,255)\n")
	} else {
		fmt.Println("verified: clear color matches")
	}

	fmt.Println("done")
}
