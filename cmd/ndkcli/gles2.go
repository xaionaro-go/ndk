package main

import (
	"fmt"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/egl"
	"github.com/xaionaro-go/ndk/gles2"
)

const (
	glVendor     gles2.Enum = 0x1F00
	glRenderer   gles2.Enum = 0x1F01
	glVersion    gles2.Enum = 0x1F02
	glExtensions gles2.Enum = 0x1F03
)

// glStringToGo converts a *GLubyte (C string) returned by glGetString to a Go string.
func glStringToGo(p unsafe.Pointer) string {
	if p == nil {
		return "<nil>"
	}
	// The pointer is to a null-terminated byte string.
	var buf []byte
	for ptr := (*byte)(p); *ptr != 0; ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
		buf = append(buf, *ptr)
	}
	return string(buf)
}

// initEGLContext creates a minimal EGL pbuffer context for OpenGL ES.
// contextVersion is 2 for ES 2.0 or 3 for ES 3.0.
func initEGLContext(contextVersion egl.Int) (
	dpy egl.EGLDisplay,
	surface egl.EGLSurface,
	ctx egl.EGLContext,
	_ error,
) {
	dpy = egl.GetDisplay(nil)
	if dpy == nil {
		return nil, nil, nil, fmt.Errorf("egl.GetDisplay returned nil")
	}

	if egl.Initialize(dpy, nil, nil) == egl.False {
		return nil, nil, nil, fmt.Errorf("egl.Initialize failed (error 0x%X)", egl.GetError())
	}

	renderableBit := egl.OpenglEs2Bit
	if contextVersion == 3 {
		renderableBit = egl.OpenglEs3Bit
	}

	attribs := []egl.Int{
		egl.RenderableType, egl.Int(renderableBit),
		egl.SurfaceType, egl.PbufferBit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.AlphaSize, 8,
		egl.None,
	}

	var cfg egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(dpy, &attribs[0], &cfg, 1, &numConfigs) == egl.False || numConfigs == 0 {
		egl.Terminate(dpy)
		return nil, nil, nil, fmt.Errorf("egl.ChooseConfig failed (error 0x%X)", egl.GetError())
	}

	pbufAttribs := []egl.Int{
		egl.Width, 1,
		egl.Height, 1,
		egl.None,
	}
	surface = egl.CreatePbufferSurface(dpy, cfg, &pbufAttribs[0])
	if surface == nil {
		egl.Terminate(dpy)
		return nil, nil, nil, fmt.Errorf("egl.CreatePbufferSurface failed (error 0x%X)", egl.GetError())
	}

	ctxAttribs := []egl.Int{
		egl.ContextClientVersion, contextVersion,
		egl.None,
	}
	ctx = egl.CreateContext(dpy, cfg, nil, &ctxAttribs[0])
	if ctx == nil {
		egl.DestroySurface(dpy, surface)
		egl.Terminate(dpy)
		return nil, nil, nil, fmt.Errorf("egl.CreateContext failed (error 0x%X)", egl.GetError())
	}

	if egl.MakeCurrent(dpy, surface, surface, ctx) == egl.False {
		egl.DestroyContext(dpy, ctx)
		egl.DestroySurface(dpy, surface)
		egl.Terminate(dpy)
		return nil, nil, nil, fmt.Errorf("egl.MakeCurrent failed (error 0x%X)", egl.GetError())
	}

	return dpy, surface, ctx, nil
}

// teardownEGLContext releases EGL resources created by initEGLContext.
func teardownEGLContext(dpy egl.EGLDisplay, surface egl.EGLSurface, ctx egl.EGLContext) {
	egl.MakeCurrent(dpy, nil, nil, nil)
	egl.DestroyContext(dpy, ctx)
	egl.DestroySurface(dpy, surface)
	egl.Terminate(dpy)
}

var gles2Cmd = &cobra.Command{
	Use:   "gles2",
	Short: "OpenGL ES 2.0 operations",
}

var gles2InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print OpenGL ES 2.0 information",
	RunE: func(cmd *cobra.Command, args []string) error {
		dpy, surface, ctx, err := initEGLContext(2)
		if err != nil {
			return fmt.Errorf("creating EGL context: %w", err)
		}
		defer teardownEGLContext(dpy, surface, ctx)

		fmt.Printf("GL_VENDOR:     %s\n", glStringToGo(unsafe.Pointer(gles2.GetString(glVendor))))
		fmt.Printf("GL_RENDERER:   %s\n", glStringToGo(unsafe.Pointer(gles2.GetString(glRenderer))))
		fmt.Printf("GL_VERSION:    %s\n", glStringToGo(unsafe.Pointer(gles2.GetString(glVersion))))
		fmt.Printf("GL_EXTENSIONS: %s\n", glStringToGo(unsafe.Pointer(gles2.GetString(glExtensions))))

		return nil
	},
}

func init() {
	gles2Cmd.AddCommand(gles2InfoCmd)
	rootCmd.AddCommand(gles2Cmd)
}
