package main

import (
	"fmt"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/AndroidGoLab/ndk/egl"
	"github.com/AndroidGoLab/ndk/gles2"
	"github.com/AndroidGoLab/ndk/gles3"
)

// glUbyteToString converts a *GLubyte (C string) to a Go string by walking
// bytes until a NUL terminator is found.
func glUbyteToString(p *gles2.GLubyte) string {
	if p == nil {
		return ""
	}

	var buf []byte
	ptr := unsafe.Pointer(p)
	for {
		b := *(*byte)(ptr)
		if b == 0 {
			break
		}
		buf = append(buf, b)
		ptr = unsafe.Add(ptr, 1)
	}
	return string(buf)
}

// initEGLForES creates an EGL display + pbuffer context for the given ES
// major version. The caller must call the returned cleanup function.
func initEGLForES(esVersion egl.Int) (
	dpy egl.EGLDisplay,
	cleanup func(),
	_err error,
) {
	var defaultDisplay egl.EGLNativeDisplayType
	dpy = egl.GetDisplay(defaultDisplay)
	if dpy == nil {
		return nil, nil, fmt.Errorf("egl.GetDisplay returned no display")
	}

	var major, minor egl.Int
	if egl.Initialize(dpy, &major, &minor) == egl.False {
		return nil, nil, fmt.Errorf("egl.Initialize failed (error 0x%x)", egl.GetError())
	}

	egl.BindAPI(egl.EGLenum(egl.OpenglEsApi))

	renderableBit := egl.OpenglEs2Bit
	if esVersion >= 3 {
		renderableBit = egl.OpenglEs3Bit
	}
	configAttribs := []egl.Int{
		egl.RenderableType, renderableBit,
		egl.SurfaceType, egl.PbufferBit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.None,
	}

	var config egl.EGLConfig
	var numConfig egl.Int
	if egl.ChooseConfig(dpy, &configAttribs[0], &config, 1, &numConfig) == egl.False || numConfig == 0 {
		egl.Terminate(dpy)
		return nil, nil, fmt.Errorf("egl.ChooseConfig failed (error 0x%x)", egl.GetError())
	}

	contextAttribs := []egl.Int{
		egl.ContextClientVersion, esVersion,
		egl.None,
	}
	ctx := egl.CreateContext(dpy, config, nil, &contextAttribs[0])
	if ctx == nil {
		egl.Terminate(dpy)
		return nil, nil, fmt.Errorf("egl.CreateContext failed (error 0x%x)", egl.GetError())
	}

	pbufAttribs := []egl.Int{
		egl.Width, 1,
		egl.Height, 1,
		egl.None,
	}
	surface := egl.CreatePbufferSurface(dpy, config, &pbufAttribs[0])
	if surface == nil {
		egl.DestroyContext(dpy, ctx)
		egl.Terminate(dpy)
		return nil, nil, fmt.Errorf("egl.CreatePbufferSurface failed (error 0x%x)", egl.GetError())
	}

	if egl.MakeCurrent(dpy, surface, surface, ctx) == egl.False {
		egl.DestroySurface(dpy, surface)
		egl.DestroyContext(dpy, ctx)
		egl.Terminate(dpy)
		return nil, nil, fmt.Errorf("egl.MakeCurrent failed (error 0x%x)", egl.GetError())
	}

	cleanup = func() {
		egl.MakeCurrent(dpy, nil, nil, nil)
		egl.DestroySurface(dpy, surface)
		egl.DestroyContext(dpy, ctx)
		egl.Terminate(dpy)
	}
	return dpy, cleanup, nil
}

var gles2InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display OpenGL ES 2 information (vendor, renderer, version, extensions)",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		_, cleanup, err := initEGLForES(2)
		if err != nil {
			return fmt.Errorf("initializing EGL for ES2: %w", err)
		}
		defer cleanup()

		fmt.Printf("GL_VENDOR:     %s\n", glUbyteToString(gles2.GetString(0x1F00)))
		fmt.Printf("GL_RENDERER:   %s\n", glUbyteToString(gles2.GetString(0x1F01)))
		fmt.Printf("GL_VERSION:    %s\n", glUbyteToString(gles2.GetString(0x1F02)))
		fmt.Printf("GL_EXTENSIONS: %s\n", glUbyteToString(gles2.GetString(0x1F03)))

		return nil
	},
}

var gles3InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display OpenGL ES 3 information (vendor, renderer, version, extensions)",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		_, cleanup, err := initEGLForES(3)
		if err != nil {
			return fmt.Errorf("initializing EGL for ES3: %w", err)
		}
		defer cleanup()

		getString := func(name uint32) string {
			p := gles3.GetString(gles3.GLenum(name))
			if p == nil {
				return ""
			}
			var buf []byte
			ptr := unsafe.Pointer(p)
			for {
				b := *(*byte)(ptr)
				if b == 0 {
					break
				}
				buf = append(buf, b)
				ptr = unsafe.Add(ptr, 1)
			}
			return string(buf)
		}

		fmt.Printf("GL_VENDOR:     %s\n", getString(0x1F00))
		fmt.Printf("GL_RENDERER:   %s\n", getString(0x1F01))
		fmt.Printf("GL_VERSION:    %s\n", getString(0x1F02))
		fmt.Printf("GL_EXTENSIONS: %s\n", getString(0x1F03))

		return nil
	},
}

func init() {
	gles2Cmd.AddCommand(gles2InfoCmd)
	gles3Cmd.AddCommand(gles3InfoCmd)
}
