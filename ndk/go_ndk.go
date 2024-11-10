// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go runtime entry point for apps running on android.
// Sets up everything the runtime needs and exposes
// the entry point to JNI.

//go:build android && !fyne
// +build android,!fyne

package ndk

/*
#cgo LDFLAGS: -llog -landroid

#include <android/log.h>
#include <android/configuration.h>
#include <android/native_activity.h>
#include <stdlib.h>
#include <time.h>

static jint _JNI_OnLoad_gondk(JavaVM* vm, void* reserved) {
	JNIEnv* env;
	if ((*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6) != JNI_OK) {
		return -1;
	}
	return JNI_VERSION_1_6;
}
*/
import "C"

import (
	"unsafe"
)

//export JNI_OnLoad
func JNI_OnLoad(vm *C.JavaVM, reserved unsafe.Pointer) C.jint {
	return C._JNI_OnLoad_gondk(vm, reserved)
}

//export ANativeActivity_onCreate
func ANativeActivity_onCreate(act *Activity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	CreateActivity(act, savedState, savedStateSize)
}

//export ANativeActivity_onCreateB
func ANativeActivity_onCreateB(act *Activity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	Info("ANativeActivity_onCreateB...")

	CreateActivity(act, savedState, savedStateSize)
}
