// Example: EGL initialization and display query.
//
// Demonstrates the basic EGL lifecycle: obtaining the default display,
// initializing it, querying vendor/version/extensions/client-API strings,
// and terminating. This is the minimal starting point for any EGL program.
package main

import (
	"fmt"
	"os"

	"github.com/xaionaro-go/ndk/egl"
)

// Standard EGL query-string name constants (not exported by the bindings).
const (
	eglVendor     egl.Int = 0x3053
	eglVersion    egl.Int = 0x3054
	eglExtensions egl.Int = 0x3055
	eglClientAPIs egl.Int = 0x308D
)

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s (EGL error 0x%04X)\n", msg, egl.GetError())
	os.Exit(1)
}

func main() {
	// 1. Obtain the default display connection.
	var defaultDisplay egl.EGLNativeDisplayType
	display := egl.GetDisplay(defaultDisplay)
	if display == nil {
		fatal("GetDisplay returned EGL_NO_DISPLAY")
	}
	fmt.Println("obtained EGL display")

	// 2. Initialize EGL on that display.
	var major, minor egl.Int
	if egl.Initialize(display, &major, &minor) == egl.False {
		fatal("Initialize failed")
	}
	fmt.Printf("EGL initialized: version %d.%d\n", major, minor)

	// 3. Query and print display strings.
	queries := []struct {
		label string
		name  egl.Int
	}{
		{"Vendor", eglVendor},
		{"Version", eglVersion},
		{"Client APIs", eglClientAPIs},
		{"Extensions", eglExtensions},
	}
	for _, q := range queries {
		s := egl.QueryString(display, q.name)
		fmt.Printf("  %-12s: %s\n", q.label, s)
	}

	// 4. Terminate EGL.
	if egl.Terminate(display) == egl.False {
		fatal("Terminate failed")
	}
	fmt.Println("EGL terminated")
}
