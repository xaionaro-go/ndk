// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go runtime entry point for apps running on android.
// Sets up everything the runtime needs and exposes
// the entry point to JNI.

package ndk

/*
#cgo LDFLAGS: -llog -landroid

#include <android/log.h>
#include <android/configuration.h>
#include <android/native_activity.h>
#include <stdlib.h>
#include <time.h>
#include <dlfcn.h>

extern int cgoLooperCallback(int p0, int p1, void* p2);
extern void* onCreate_gondk(ANativeActivity* p0, void* p1, size_t p2);
extern void onStart_gondk(ANativeActivity* p0);
extern void onResume_gondk(ANativeActivity* p0);
extern void onPause_gondk(ANativeActivity* p0);
extern void onStop_gondk(ANativeActivity* p0);
extern void onDestroy_gondk(ANativeActivity* p0);
extern void onWindowFocusChanged_gondk(ANativeActivity* p0, int p1);
extern void* onSaveInstanceState_gondk(ANativeActivity* p0, size_t* p1);
extern void onNativeWindowCreated_gondk(ANativeActivity* p0, ANativeWindow* p1);
extern void onNativeWindowResized_gondk(ANativeActivity* p0, ANativeWindow* p1);
extern void onNativeWindowRedrawNeeded_gondk(ANativeActivity* p0, ANativeWindow* p1);
extern void onNativeWindowDestroyed_gondk(ANativeActivity* p0, ANativeWindow* p1);
extern void onInputQueueCreated_gondk(ANativeActivity* p0, AInputQueue* p1);
extern void onInputQueueDestroyed_gondk(ANativeActivity* p0, AInputQueue* p1);
extern void onContentRectChanged_gondk(ANativeActivity* p0, ARect* p1);
extern void onConfigurationChanged_gondk(ANativeActivity* p0);
extern void onLowMemory_gondk(ANativeActivity* p0);

static const char* callGetStringMethod(JNIEnv *env, jclass jobj, const char* method) {
	jstring jpath;
	jclass _class = (*env)->GetObjectClass(env, jobj);
	jmethodID m = (*env)->GetMethodID(env, _class, method, "()Ljava/lang/String;");
	if (m == 0) {
		(*env)->ExceptionClear(env);
		return NULL;
	}
	jpath = (jstring)(*env)->CallObjectMethod(env, jobj, m, NULL);
	return (*env)->GetStringUTFChars(env, jpath, NULL);
}

static const char* _GetPackageName(JNIEnv *env, jobject jobj) {
	return callGetStringMethod(env, jobj, "getPackageName");
}

static const char* _GetLocalClassName(JNIEnv *env, jobject jobj) {
	return callGetStringMethod(env, jobj, "getLocalClassName");
}

static void* _GetMainPC() { return dlsym(RTLD_DEFAULT, "main.main"); }

static void* _SetActivityCallbacks(ANativeActivity* activity) {
	activity->callbacks->onStart = onStart_gondk;
	activity->callbacks->onResume = onResume_gondk;
	activity->callbacks->onSaveInstanceState = onSaveInstanceState_gondk;
	activity->callbacks->onPause = onPause_gondk;
	activity->callbacks->onStop = onStop_gondk;
	activity->callbacks->onDestroy = onDestroy_gondk;
	activity->callbacks->onWindowFocusChanged = onWindowFocusChanged_gondk;
	activity->callbacks->onNativeWindowCreated = onNativeWindowCreated_gondk;
	activity->callbacks->onNativeWindowResized = onNativeWindowResized_gondk;
	activity->callbacks->onNativeWindowRedrawNeeded = onNativeWindowRedrawNeeded_gondk;
	activity->callbacks->onNativeWindowDestroyed = onNativeWindowDestroyed_gondk;
	activity->callbacks->onInputQueueCreated = onInputQueueCreated_gondk;
	activity->callbacks->onInputQueueDestroyed = onInputQueueDestroyed_gondk;
	activity->callbacks->onContentRectChanged = (void*)onContentRectChanged_gondk;
	activity->callbacks->onConfigurationChanged = onConfigurationChanged_gondk;
	activity->callbacks->onLowMemory = onLowMemory_gondk;
	return (void*)activity;
}

*/
import "C"

import (
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/xaionaro-go/ndk/callfn"
)

var appContext struct {
	mainPC    unsafe.Pointer
	funcChan  chan func()
	actMaps   map[string]*Activity
	actMainCB func(*Context)
}

//export onStart_gondk
func onStart_gondk(act *Activity) {
	ctx := act.Context()
	Info("onStart:", act)
	ctx.Reset()
	if ctx.Start != nil {
		ctx.RunFunc(func() {
			ctx.Start(ctx.Act)
		}, true)
	}
}

//export onResume_gondk
func onResume_gondk(act *Activity) {
	ctx := act.Context()
	Info("onResume:", act)
	if ctx.Resume != nil {
		ctx.RunFunc(func() {
			ctx.Resume(ctx.Act)
		}, true)
	}
	ctx.IsResume = true
}

//export onPause_gondk
func onPause_gondk(act *Activity) {
	ctx := act.Context()
	Info("onPause:", act)
	ctx.IsResume = false
	ctx.RunFunc(func() {
		if ctx.Pause != nil {
			ctx.Pause(ctx.Act)
		}
	}, true)
}

//export onStop_gondk
func onStop_gondk(act *Activity) {
	ctx := act.Context()
	Info("onStop:", act)
	ctx.IsResume = false
	ctx.RunFunc(func() {
		if ctx.Stop != nil {
			ctx.Stop(ctx.Act)
		}
	}, true)
}

//export onDestroy_gondk
func onDestroy_gondk(act *Activity) {
	ctx := act.Context()
	Info("onDestroy:", act)
	ctx.RunFunc(func() {
		ctx.WillDestoryValue = true
		if ctx.Input != nil {
			ctx.Input.DetachLooper()
			ctx.Input = nil
		}
		if ctx.Destroy != nil {
			ctx.Destroy(ctx.Act)
		}
		ctx.Release()
	}, true)

	ctx.FuncChan <- func() {
		Info("onDestroy:", act, "complete")
	}
	ClrActivity(act, ctx.ClassName)
}

//export onWindowFocusChanged_gondk
func onWindowFocusChanged_gondk(act *Activity, hasFocus C.int) {
	ctx := act.Context()
	Info("onWindowFocusChanged:", act, hasFocus)
	focus := hasFocus != 0
	ctx.IsFocus = focus
	ctx.RunFunc(func() {
		if ctx.FocusChanged != nil {
			ctx.FocusChanged(ctx.Act, focus)
		}
	}, true)
}

//export onSaveInstanceState_gondk
func onSaveInstanceState_gondk(act *Activity, outSize *C.size_t) unsafe.Pointer {
	ctx := act.Context()
	Info("onSaveInstanceState:", act)
	if ctx.SaveState != nil {
		ctx.RunFunc(func() {
			ctx.SavedState = ctx.SaveState(ctx.Act)
		}, true)

		if len(ctx.SavedState) > 0 {
			size := len(ctx.SavedState)
			Info("\t\tsize =", size)
			*outSize = C.size_t(size)
			ptr := C.malloc(C.size_t(size))
			copy((*[1 << 30]byte)(unsafe.Pointer(ptr))[:size], ctx.SavedState)
			return ptr
		}
	}
	return nil
}

//export onNativeWindowCreated_gondk
func onNativeWindowCreated_gondk(act *Activity, window *Window) {
	ctx := act.Context()
	Info("onNativeWindowCreated:", act, window)
	ctx.Window = window
	if ctx.WindowCreated != nil {
		ctx.RunFunc(func() {
			ctx.WindowCreated(act, window)
		}, true)
	}
}

//export onNativeWindowResized_gondk
func onNativeWindowResized_gondk(act *Activity, window *Window) {
	ctx := act.Context()
	Info("onNativeWindowResized:", act, window)
	if ctx.WindowResized != nil {
		ctx.RunFunc(func() {
			ctx.WindowResized(act, window)
		}, true)
	}
}

//export onNativeWindowRedrawNeeded_gondk
func onNativeWindowRedrawNeeded_gondk(act *Activity, window *Window) {
	ctx := act.Context()
	Info("onNativeWindowRedrawNeeded:", act, window)
	Assert(ctx.Window == window)
	if ctx.Window != window {
		ctx.Window = window
	}

	if ctx.WindowRedrawNeeded != nil {
		ctx.RunFunc(func() {
			ctx.WindowRedrawNeeded(act, window)
		}, true)
	}
}

//export onNativeWindowDestroyed_gondk
func onNativeWindowDestroyed_gondk(act *Activity, window *Window) {
	ctx := act.Context()
	Info("onNativeWindowDestroyed:", act, window)
	ctx.RunFunc(func() {
		//Info("onNativeWindowDestroyed.func")
		if ctx.WindowDestroyed != nil {
			ctx.WindowDestroyed(act, window)
		}
		ctx.Window = nil
	}, true)
}

//export onInputQueueCreated_gondk
func onInputQueueCreated_gondk(act *Activity, queue *InputQueue) {
	ctx := act.Context()
	Info("onInputQueueCreated:", act, queue)
	ctx.RunFunc(func() {
		ctx.Input = (*InputQueue)(queue)
		ctx.Input.AttachLooper(ctx.Looper, Looper_ID_INPUT, nil, unsafe.Pointer(uintptr(Looper_ID_INPUT)))
	}, true)
}

//export onInputQueueDestroyed_gondk
func onInputQueueDestroyed_gondk(act *Activity, queue *InputQueue) {
	ctx := act.Context()
	Info("onInputQueueDestroyed:", act, queue)
	ctx.RunFunc(func() {
		ctx.Input.DetachLooper()
		ctx.Input = nil
	}, true)
}

//export onContentRectChanged_gondk
func onContentRectChanged_gondk(act *Activity, rect *Rect) {
	ctx := act.Context()
	Info("onContentRectChanged:", act, rect)
	if ctx.ContentRectChanged != nil {
		ctx.RunFunc(func() {
			ctx.ContentRectChanged(act, rect)
		}, true)
	}
}

//export onConfigurationChanged_gondk
func onConfigurationChanged_gondk(act *Activity) {
	ctx := act.Context()
	Info("onConfigurationChanged:", act)
	if ctx.ConfigurationChanged != nil {
		ctx.RunFunc(func() {
			ctx.ConfigurationChanged(act)
		}, true)
	}
}

//export onLowMemory_gondk
func onLowMemory_gondk(act *Activity) {
	ctx := act.Context()
	Info("onLowMemory:", act)
	if ctx.LowMemory != nil {
		ctx.RunFunc(func() {
			ctx.LowMemory(act)
		}, true)
	}
}

func SetMainCB(fn func(*Context)) {
	if fn == nil {
		Fatal("SetMainCB(nil) is incorrect")
	}
	if appContext.actMainCB != nil {
		Info("MainCB is ready")
		return
	}
	appContext.funcChan <- func() { appContext.actMainCB = fn }
	appContext.funcChan <- func() {}
}

// Loop
// applation will close, return false
func Loop() bool {
	if fn, ok := <-appContext.funcChan; ok {
		fn()
		return true
	}
	return false
}

func addActivity(act *Activity, n string) {
	appContext.actMaps[n] = act
}

func ClrActivity(act *Activity, n string) {
	delete(appContext.actMaps, n)
	if len(appContext.actMaps) == 0 {
		close(appContext.funcChan)
	}
}

func callMain() {
	if appContext.mainPC != nil {
		Info("main is running")
		return
	}

	var waitMainCB func()
	appContext.actMaps = make(map[string]*Activity)
	appContext.funcChan, waitMainCB = make(chan func(), 1), func() { (<-appContext.funcChan)() }
	appContext.mainPC = C._GetMainPC()
	if appContext.mainPC == nil {
		Fatal("missing main.main")
	}

	for _, name := range []string{"TMPDIR", "PATH", "LD_LIBRARY_PATH", "BOOTCLASSPATH"} {
		n := C.CString(name)
		os.Setenv(name, C.GoString(C.getenv(n)))
		C.free(unsafe.Pointer(n))
	}

	// Set timezone.
	//
	// Note that Android zoneinfo is stored in /system/usr/share/zoneinfo,
	// but it is in some kind of packed TZiff file that we do not support
	// yet. As a stopgap, we build a fixed zone using the tm_zone name.
	var curtime C.time_t
	var curtm C.struct_tm
	C.time(&curtime)
	C.localtime_r(&curtime, &curtm)
	tzOffset := int(curtm.tm_gmtoff)
	tz := C.GoString(curtm.tm_zone)
	time.Local = time.FixedZone(tz, tzOffset)

	go func() {
		runtime.LockOSThread()
		callfn.CallFn(uintptr(appContext.mainPC))
		os.Exit(0)
	}()

	// This will hang until SetMainCB is completed
	waitMainCB()
}

func CreateActivity(act *Activity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	lname := C.GoString(C._GetLocalClassName(act.env, act.clazz))
	pname := C.GoString(C._GetPackageName(act.env, act.clazz))
	Info("ANativeActivity_onCreate:", pname+"/"+lname)

	callMain()

	C._SetActivityCallbacks(act.cptr())

	buf := []byte{}
	if savedStateSize > 0 {
		buf = (*[1 << 30]byte)(unsafe.Pointer(savedState))[:savedStateSize]
	}

	ctx := &Context{}
	act.instance = unsafe.Pointer(ctx)
	ctx.init()

	ctx.ClassName = lname
	ctx.packageName = pname
	addActivity(act, ctx.ClassName)
	appContext.funcChan <- func() {
		go func() {
			runtime.LockOSThread()
			appContext.actMainCB(ctx)
		}()
	}

	// Wait for the message queue initialization (ctx.begin) to complete
	ctx.doFunc()
	OnCreate(act, buf)
}

func OnCreate(act *Activity, buf []byte) {
	ctx := act.Context()
	ctx.RunFunc(func() {}, true)
	Info("onCreate:", act, len(buf))

	ctx.Act = (*Activity)(act)
	if len(buf) > 0 {
		ctx.SavedState = make([]byte, len(buf))
		copy(ctx.SavedState, buf)
	}

	if ctx.Create != nil {
		ctx.RunFunc(func() {
			ctx.Create(act, ctx.SavedState)
		}, true)
	}
}
