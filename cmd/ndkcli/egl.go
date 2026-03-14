package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/egl"
)

const (
	eglVendor     egl.Int = 0x3053
	eglVersion    egl.Int = 0x3054
	eglExtensions egl.Int = 0x3055
	eglClientAPIs egl.Int = 0x308D
)

var eglCmd = &cobra.Command{
	Use:   "egl",
	Short: "EGL display and context operations",
}

var eglInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print EGL display information",
	RunE: func(cmd *cobra.Command, args []string) error {
		dpy := egl.GetDisplay(nil)
		if dpy == nil {
			return fmt.Errorf("egl.GetDisplay returned nil")
		}

		if egl.Initialize(dpy, nil, nil) == egl.False {
			return fmt.Errorf("egl.Initialize failed (error 0x%X)", egl.GetError())
		}
		defer egl.Terminate(dpy)

		fmt.Printf("Vendor:      %s\n", egl.QueryString(dpy, eglVendor))
		fmt.Printf("Version:     %s\n", egl.QueryString(dpy, eglVersion))
		fmt.Printf("Extensions:  %s\n", egl.QueryString(dpy, eglExtensions))
		fmt.Printf("Client APIs: %s\n", egl.QueryString(dpy, eglClientAPIs))

		return nil
	},
}

func init() {
	eglCmd.AddCommand(eglInfoCmd)
	rootCmd.AddCommand(eglCmd)
}
