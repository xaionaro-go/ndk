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

	"github.com/AndroidGoLab/ndk/egl"
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
		{"Samples", egl.EGL_SAMPLES},
		{"Sample Buffers", egl.EGL_SAMPLE_BUFFERS},
		{"Max Pbuffer Width", egl.EGL_MAX_PBUFFER_WIDTH},
		{"Max Pbuffer Height", egl.EGL_MAX_PBUFFER_HEIGHT},
		{"Max Pbuffer Pixels", egl.EGL_MAX_PBUFFER_PIXELS},
		{"Native Renderable", egl.EGL_NATIVE_RENDERABLE},
		{"Native Visual ID", egl.EGL_NATIVE_VISUAL_ID},
		{"Surface Type (bitmask)", egl.SurfaceType},
		{"Renderable Type (bitmask)", egl.RenderableType},
		{"Conformant (bitmask)", egl.EGL_CONFORMANT},
		{"Config Caveat", egl.EGL_CONFIG_CAVEAT},
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
	printFlag("  Pixmap", surfType, egl.EGL_PIXMAP_BIT)

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
