//go:build ignore

package gameactivity

/*
#include "include/game-activity/GameActivity.h"
#include <android/native_window.h>

// Forward declarations for Go callback trampolines.
extern void goOnCreate(GameActivity* activity, void* savedState, size_t savedStateSize);
extern void goOnDestroy(GameActivity* activity);
extern void goOnStart(GameActivity* activity);
extern void goOnStop(GameActivity* activity);
extern void goOnPause(GameActivity* activity);
extern void goOnResume(GameActivity* activity);
extern void goOnConfigurationChanged(GameActivity* activity);
extern void goOnTrimMemory(GameActivity* activity, int level);
extern void goOnWindowFocusChanged(GameActivity* activity, bool hasFocus);
extern void goOnNativeWindowCreated(GameActivity* activity, ANativeWindow* window);
extern void goOnNativeWindowDestroyed(GameActivity* activity, ANativeWindow* window);
extern void goOnNativeWindowResized(GameActivity* activity, ANativeWindow* window, int32_t w, int32_t h);
extern void goOnNativeWindowRedrawNeeded(GameActivity* activity, ANativeWindow* window);
extern void goOnWindowInsetsChanged(GameActivity* activity);
extern bool goOnTouchEvent(GameActivity* activity, const GameActivityMotionEvent* event);
extern bool goOnKeyDown(GameActivity* activity, const GameActivityKeyEvent* event);
extern bool goOnKeyUp(GameActivity* activity, const GameActivityKeyEvent* event);
extern void goOnTextInputEvent(GameActivity* activity, const GameTextInputState* state);

static void installCallbacks(GameActivityCallbacks* cb) {
	cb->onDestroy = goOnDestroy;
	cb->onStart = goOnStart;
	cb->onStop = goOnStop;
	cb->onPause = goOnPause;
	cb->onResume = goOnResume;
	cb->onConfigurationChanged = goOnConfigurationChanged;
	cb->onTrimMemory = goOnTrimMemory;
	cb->onWindowFocusChanged = goOnWindowFocusChanged;
	cb->onNativeWindowCreated = goOnNativeWindowCreated;
	cb->onNativeWindowDestroyed = goOnNativeWindowDestroyed;
	cb->onNativeWindowResized = goOnNativeWindowResized;
	cb->onNativeWindowRedrawNeeded = goOnNativeWindowRedrawNeeded;
	cb->onWindowInsetsChanged = goOnWindowInsetsChanged;
	cb->onTouchEvent = goOnTouchEvent;
	cb->onKeyDown = goOnKeyDown;
	cb->onKeyUp = goOnKeyUp;
	cb->onTextInputEvent = goOnTextInputEvent;
}
*/
import "C"

import (
	"sync"
	"unsafe"
)

// Callbacks defines the lifecycle and event callbacks for GameActivity.
// All callbacks execute on the main thread.
type Callbacks struct {
	OnDestroy              func(activity *Activity)
	OnStart                func(activity *Activity)
	OnStop                 func(activity *Activity)
	OnPause                func(activity *Activity)
	OnResume               func(activity *Activity)
	OnConfigurationChanged func(activity *Activity)
	OnTrimMemory           func(activity *Activity, level int)
	OnWindowFocusChanged   func(activity *Activity, hasFocus bool)

	OnNativeWindowCreated     func(activity *Activity, window unsafe.Pointer)
	OnNativeWindowDestroyed   func(activity *Activity, window unsafe.Pointer)
	OnNativeWindowResized     func(activity *Activity, window unsafe.Pointer, width, height int32)
	OnNativeWindowRedrawNeeded func(activity *Activity, window unsafe.Pointer)
	OnWindowInsetsChanged     func(activity *Activity)

	OnTouchEvent     func(activity *Activity, event MotionEvent) bool
	OnKeyDown        func(activity *Activity, event KeyEvent) bool
	OnKeyUp          func(activity *Activity, event KeyEvent) bool
	OnTextInputEvent func(activity *Activity, state TextInputState)
}

var (
	callbacksMu       sync.Mutex
	registeredCallbacks *Callbacks
)

// SetCallbacks registers the Go callback handlers and installs
// C trampoline functions into the GameActivity's callback table.
// Only one set of callbacks can be active at a time.
func SetCallbacks(activity *Activity, cb *Callbacks) {
	callbacksMu.Lock()
	registeredCallbacks = cb
	callbacksMu.Unlock()
	C.installCallbacks(activity.ptr.callbacks)
}

func getCallbacks() *Callbacks {
	callbacksMu.Lock()
	defer callbacksMu.Unlock()
	return registeredCallbacks
}

//export goOnDestroy
func goOnDestroy(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnDestroy != nil {
		cb.OnDestroy(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnStart
func goOnStart(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnStart != nil {
		cb.OnStart(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnStop
func goOnStop(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnStop != nil {
		cb.OnStop(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnPause
func goOnPause(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnPause != nil {
		cb.OnPause(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnResume
func goOnResume(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnResume != nil {
		cb.OnResume(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnConfigurationChanged
func goOnConfigurationChanged(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnConfigurationChanged != nil {
		cb.OnConfigurationChanged(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnTrimMemory
func goOnTrimMemory(activity *C.GameActivity, level C.int) {
	if cb := getCallbacks(); cb != nil && cb.OnTrimMemory != nil {
		cb.OnTrimMemory(NewActivityFromPointer(unsafe.Pointer(activity)), int(level))
	}
}

//export goOnWindowFocusChanged
func goOnWindowFocusChanged(activity *C.GameActivity, hasFocus C.bool) {
	if cb := getCallbacks(); cb != nil && cb.OnWindowFocusChanged != nil {
		cb.OnWindowFocusChanged(NewActivityFromPointer(unsafe.Pointer(activity)), bool(hasFocus))
	}
}

//export goOnNativeWindowCreated
func goOnNativeWindowCreated(activity *C.GameActivity, window *C.ANativeWindow) {
	if cb := getCallbacks(); cb != nil && cb.OnNativeWindowCreated != nil {
		cb.OnNativeWindowCreated(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			unsafe.Pointer(window),
		)
	}
}

//export goOnNativeWindowDestroyed
func goOnNativeWindowDestroyed(activity *C.GameActivity, window *C.ANativeWindow) {
	if cb := getCallbacks(); cb != nil && cb.OnNativeWindowDestroyed != nil {
		cb.OnNativeWindowDestroyed(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			unsafe.Pointer(window),
		)
	}
}

//export goOnNativeWindowResized
func goOnNativeWindowResized(
	activity *C.GameActivity,
	window *C.ANativeWindow,
	w, h C.int32_t,
) {
	if cb := getCallbacks(); cb != nil && cb.OnNativeWindowResized != nil {
		cb.OnNativeWindowResized(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			unsafe.Pointer(window),
			int32(w), int32(h),
		)
	}
}

//export goOnNativeWindowRedrawNeeded
func goOnNativeWindowRedrawNeeded(activity *C.GameActivity, window *C.ANativeWindow) {
	if cb := getCallbacks(); cb != nil && cb.OnNativeWindowRedrawNeeded != nil {
		cb.OnNativeWindowRedrawNeeded(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			unsafe.Pointer(window),
		)
	}
}

//export goOnWindowInsetsChanged
func goOnWindowInsetsChanged(activity *C.GameActivity) {
	if cb := getCallbacks(); cb != nil && cb.OnWindowInsetsChanged != nil {
		cb.OnWindowInsetsChanged(NewActivityFromPointer(unsafe.Pointer(activity)))
	}
}

//export goOnTouchEvent
func goOnTouchEvent(activity *C.GameActivity, event *C.GameActivityMotionEvent) C.bool {
	if cb := getCallbacks(); cb != nil && cb.OnTouchEvent != nil {
		return C.bool(cb.OnTouchEvent(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			motionEventFromC(event),
		))
	}
	return C.bool(false)
}

//export goOnKeyDown
func goOnKeyDown(activity *C.GameActivity, event *C.GameActivityKeyEvent) C.bool {
	if cb := getCallbacks(); cb != nil && cb.OnKeyDown != nil {
		return C.bool(cb.OnKeyDown(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			keyEventFromC(event),
		))
	}
	return C.bool(false)
}

//export goOnKeyUp
func goOnKeyUp(activity *C.GameActivity, event *C.GameActivityKeyEvent) C.bool {
	if cb := getCallbacks(); cb != nil && cb.OnKeyUp != nil {
		return C.bool(cb.OnKeyUp(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			keyEventFromC(event),
		))
	}
	return C.bool(false)
}

//export goOnTextInputEvent
func goOnTextInputEvent(activity *C.GameActivity, state *C.GameTextInputState) {
	if cb := getCallbacks(); cb != nil && cb.OnTextInputEvent != nil {
		goState := TextInputState{
			Text:           C.GoStringN(state.text_UTF8, C.int(state.text_length)),
			SelectionStart: int32(state.selection.start),
			SelectionEnd:   int32(state.selection.end),
			ComposingStart: int32(state.composingRegion.start),
			ComposingEnd:   int32(state.composingRegion.end),
		}
		cb.OnTextInputEvent(
			NewActivityFromPointer(unsafe.Pointer(activity)),
			goState,
		)
	}
}

//export goOnCreate
func goOnCreate(activity *C.GameActivity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	// This is the entry point called by the framework.
	// The application should call SetCallbacks from here or from
	// whatever initialization function is triggered by GameActivity_onCreate.
}
