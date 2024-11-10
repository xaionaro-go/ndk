// Copyright 2017 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package app

import (
	app "github.com/gooid/gooid/ndk"
)

type Callbacks = app.Callbacks
type Activity = app.Activity
type Window = app.Window
type InputEvent = app.InputEvent
type Context = app.Context

func SetMainCB(fn func(*Context)) {
	app.SetMainCB(fn)
}

func Loop() bool {
	return app.Loop()
}

// getprop
func PropGet(k string) string {
	return app.PropGet(k)
}

// visitor all properties
func PropVisit(cb func(k, v string)) {
	app.PropVisit(cb)
}

// FindMatchLibrary find library path
//
//	see filepath.Glob(pattern string)
func FindMatchLibrary(pattern string) []string {
	return app.FindMatchLibrary(pattern)
}
