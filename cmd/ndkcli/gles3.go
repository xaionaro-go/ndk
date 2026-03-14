package main

import (
	"fmt"
	"unsafe"

	"github.com/spf13/cobra"
	capigles3 "github.com/xaionaro-go/ndk/capi/gles3"
)

// gles3 idiomatic package does not expose GetString (only GetStringi),
// so we use the capi directly for querying GL strings.

const (
	gl3Vendor     capigles3.GLenum = 0x1F00
	gl3Renderer   capigles3.GLenum = 0x1F01
	gl3Version    capigles3.GLenum = 0x1F02
	gl3Extensions capigles3.GLenum = 0x1F03
)

var gles3Cmd = &cobra.Command{
	Use:   "gles3",
	Short: "OpenGL ES 3.0 operations",
}

var gles3InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print OpenGL ES 3.0 information",
	RunE: func(cmd *cobra.Command, args []string) error {
		dpy, surface, ctx, err := initEGLContext(3)
		if err != nil {
			return fmt.Errorf("creating EGL context: %w", err)
		}
		defer teardownEGLContext(dpy, surface, ctx)

		fmt.Printf("GL_VENDOR:     %s\n", glStringToGo(unsafe.Pointer(capigles3.GlGetString(gl3Vendor))))
		fmt.Printf("GL_RENDERER:   %s\n", glStringToGo(unsafe.Pointer(capigles3.GlGetString(gl3Renderer))))
		fmt.Printf("GL_VERSION:    %s\n", glStringToGo(unsafe.Pointer(capigles3.GlGetString(gl3Version))))
		fmt.Printf("GL_EXTENSIONS: %s\n", glStringToGo(unsafe.Pointer(capigles3.GlGetString(gl3Extensions))))

		return nil
	},
}

func init() {
	gles3Cmd.AddCommand(gles3InfoCmd)
	rootCmd.AddCommand(gles3Cmd)
}
