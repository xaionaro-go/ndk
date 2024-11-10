// Copyright 2017-2024 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ndk

import (
	"github.com/xaionaro-go/ndk/ndk"
)

type Callbacks = ndk.Callbacks
type Activity = ndk.Activity
type Window = ndk.Window
type InputEvent = ndk.InputEvent
type Context = ndk.Context

func SetMainCB(fn func(*Context)) {
	ndk.SetMainCB(fn)
}

func Loop() bool {
	return ndk.Loop()
}

// getprop
func PropGet(k string) string {
	return ndk.PropGet(k)
}

// visitor all properties
func PropVisit(cb func(k, v string)) {
	ndk.PropVisit(cb)
}

// FindMatchLibrary find library path
//
//	see filepath.Glob(pattern string)
func FindMatchLibrary(pattern string) []string {
	return ndk.FindMatchLibrary(pattern)
}
