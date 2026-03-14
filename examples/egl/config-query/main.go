// Example: EGL configuration attribute query.
//
// Demonstrates how to choose an EGL config and inspect its attributes.
// After selecting a config that supports pbuffers and OpenGL ES 2, the
// program queries and prints key framebuffer properties: color channel
// sizes, depth/stencil sizes, sample buffers, and more.
package main

import (
	"fmt"
	"os"

	"github.com/xaionaro-go/ndk/egl"
)

// EGL attribute constants not exported by the idiomatic bindings.
const (
	eglSamples          egl.Int = 12337
	eglSampleBuffers    egl.Int = 12338
	eglMaxPbufferWidth  egl.Int = 12330
	eglMaxPbufferHeight egl.Int = 12329
	eglMaxPbufferPixels egl.Int = 12331
	eglNativeRenderable egl.Int = 12346
	eglNativeVisualId   egl.Int = 12349
	eglConformant       egl.Int = 12354
	eglConfigCaveat     egl.Int = 12327
	eglPixmapBit        egl.Int = 2
)

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s (EGL error 0x%04X)\n", msg, egl.GetError())
	os.Exit(1)
}

// getAttrib queries a single config attribute and returns it.
// Terminates the program on failure.
func getAttrib(display egl.EGLDisplay, config egl.EGLConfig, attr egl.Int) egl.Int {
	var val egl.Int
	if egl.GetConfigAttrib(display, config, attr, &val) == egl.False {
		fatal(fmt.Sprintf("GetConfigAttrib(0x%04X) failed", attr))
	}
	return val
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
	defer egl.Terminate(display)
	fmt.Printf("EGL %d.%d initialized\n", major, minor)

	// --- Choose config ---
	// Ask for a config that supports pbuffers and OpenGL ES 2.
	attribs := [...]egl.Int{
		egl.SurfaceType, egl.PbufferBit,
		egl.RenderableType, egl.OpenglEs2Bit,
		egl.None,
	}
	var config egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(display, &attribs[0], &config, 1, &numConfigs) == egl.False {
		fatal("ChooseConfig failed")
	}
	if numConfigs == 0 {
		fatal("no matching configs found")
	}
	fmt.Printf("selected 1 config out of %d match(es)\n\n", numConfigs)

	// --- Query attributes ---
	type attrInfo struct {
		name string
		id   egl.Int
	}
	attributes := []attrInfo{
		{"Config ID", egl.ConfigID},
		{"Buffer Size (total bits)", egl.BufferSize},
		{"Red Size", egl.RedSize},
		{"Green Size", egl.GreenSize},
		{"Blue Size", egl.BlueSize},
		{"Alpha Size", egl.AlphaSize},
		{"Depth Size", egl.DepthSize},
		{"Stencil Size", egl.StencilSize},
		{"Samples", eglSamples},
		{"Sample Buffers", eglSampleBuffers},
		{"Max Pbuffer Width", eglMaxPbufferWidth},
		{"Max Pbuffer Height", eglMaxPbufferHeight},
		{"Max Pbuffer Pixels", eglMaxPbufferPixels},
		{"Native Renderable", eglNativeRenderable},
		{"Native Visual ID", eglNativeVisualId},
		{"Surface Type (bitmask)", egl.SurfaceType},
		{"Renderable Type (bitmask)", egl.RenderableType},
		{"Conformant (bitmask)", eglConformant},
		{"Config Caveat", eglConfigCaveat},
	}

	fmt.Println("Config attributes:")
	for _, a := range attributes {
		val := getAttrib(display, config, a.id)
		fmt.Printf("  %-30s = %d (0x%04X)\n", a.name, val, val)
	}

	// --- Interpret some bitmask fields for readability ---
	fmt.Println("\nSurface type flags:")
	surfType := getAttrib(display, config, egl.SurfaceType)
	printFlag("  Window", surfType, egl.WindowBit)
	printFlag("  Pbuffer", surfType, egl.PbufferBit)
	printFlag("  Pixmap", surfType, eglPixmapBit)

	fmt.Println("Renderable type flags:")
	rendType := getAttrib(display, config, egl.RenderableType)
	printFlag("  OpenGL ES 2", rendType, egl.OpenglEs2Bit)
	printFlag("  OpenGL ES 3", rendType, egl.OpenglEs3Bit)

	fmt.Println("\ndone")
}

func printFlag(label string, bitmask, flag egl.Int) {
	if bitmask&flag != 0 {
		fmt.Printf("%s: yes\n", label)
	} else {
		fmt.Printf("%s: no\n", label)
	}
}
