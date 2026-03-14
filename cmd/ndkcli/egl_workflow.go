package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/egl"
)

var eglInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display EGL information (vendor, version, extensions, client APIs)",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		var defaultDisplay egl.EGLNativeDisplayType
		dpy := egl.GetDisplay(defaultDisplay)
		if dpy == nil {
			return fmt.Errorf("egl.GetDisplay returned no display")
		}

		var major, minor egl.Int
		if egl.Initialize(dpy, &major, &minor) == egl.False {
			return fmt.Errorf("egl.Initialize failed (error 0x%x)", egl.GetError())
		}
		defer egl.Terminate(dpy)

		fmt.Printf("EGL version: %d.%d\n", major, minor)
		fmt.Printf("vendor:      %s\n", egl.QueryString(dpy, 0x3053))
		fmt.Printf("version:     %s\n", egl.QueryString(dpy, 0x3054))
		fmt.Printf("extensions:  %s\n", egl.QueryString(dpy, 0x3055))
		fmt.Printf("client APIs: %s\n", egl.QueryString(dpy, 0x308D))

		return nil
	},
}

var eglConfigsCmd = &cobra.Command{
	Use:   "configs",
	Short: "List available EGL configurations",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		var defaultDisplay egl.EGLNativeDisplayType
		dpy := egl.GetDisplay(defaultDisplay)
		if dpy == nil {
			return fmt.Errorf("egl.GetDisplay returned no display")
		}

		var major, minor egl.Int
		if egl.Initialize(dpy, &major, &minor) == egl.False {
			return fmt.Errorf("egl.Initialize failed (error 0x%x)", egl.GetError())
		}
		defer egl.Terminate(dpy)

		var config egl.EGLConfig
		var numConfig egl.Int
		if egl.ChooseConfig(dpy, nil, &config, 1, &numConfig) == egl.False {
			return fmt.Errorf("egl.ChooseConfig failed (error 0x%x)", egl.GetError())
		}

		if numConfig == 0 {
			fmt.Println("no EGL configs found")
			return nil
		}

		var redSize, greenSize, blueSize, depthSize, surfaceType egl.Int
		egl.GetConfigAttrib(dpy, config, egl.RedSize, &redSize)
		egl.GetConfigAttrib(dpy, config, egl.GreenSize, &greenSize)
		egl.GetConfigAttrib(dpy, config, egl.BlueSize, &blueSize)
		egl.GetConfigAttrib(dpy, config, egl.DepthSize, &depthSize)
		egl.GetConfigAttrib(dpy, config, egl.SurfaceType, &surfaceType)

		fmt.Printf("config 0:\n")
		fmt.Printf("  red size:     %d\n", redSize)
		fmt.Printf("  green size:   %d\n", greenSize)
		fmt.Printf("  blue size:    %d\n", blueSize)
		fmt.Printf("  depth size:   %d\n", depthSize)
		fmt.Printf("  surface type: 0x%x\n", surfaceType)

		return nil
	},
}

func init() {
	eglCmd.AddCommand(eglInfoCmd)
	eglCmd.AddCommand(eglConfigsCmd)
}
