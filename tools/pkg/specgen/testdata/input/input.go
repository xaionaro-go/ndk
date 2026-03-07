// Simulates c-for-go output for Android Input.
// This file is parsed at AST level only; it does not compile.
package input

import "unsafe"

// Opaque handle types.
type AInputEvent C.AInputEvent
type AInputQueue C.AInputQueue

// Integer typedefs.
type Input_event_type_t int32
type Key_action_t int32
type Motion_action_t int32
type Key_source_t int32

// Event type enum.
const (
	AINPUT_EVENT_TYPE_KEY    Input_event_type_t = 1
	AINPUT_EVENT_TYPE_MOTION Input_event_type_t = 2
)

// Key action enum.
const (
	AKEY_EVENT_ACTION_DOWN     Key_action_t = 0
	AKEY_EVENT_ACTION_UP       Key_action_t = 1
	AKEY_EVENT_ACTION_MULTIPLE Key_action_t = 2
)

// Motion action enum.
const (
	AMOTION_EVENT_ACTION_DOWN         Motion_action_t = 0
	AMOTION_EVENT_ACTION_UP           Motion_action_t = 1
	AMOTION_EVENT_ACTION_MOVE         Motion_action_t = 2
	AMOTION_EVENT_ACTION_CANCEL       Motion_action_t = 3
	AMOTION_EVENT_ACTION_POINTER_DOWN Motion_action_t = 5
	AMOTION_EVENT_ACTION_POINTER_UP   Motion_action_t = 6
)

// Input source enum.
const (
	AINPUT_SOURCE_KEYBOARD    Key_source_t = 257
	AINPUT_SOURCE_TOUCHSCREEN Key_source_t = 4098
	AINPUT_SOURCE_MOUSE       Key_source_t = 8194
	AINPUT_SOURCE_JOYSTICK    Key_source_t = 16777232
)

// --- Input event functions ---
func AInputEvent_getType(event *AInputEvent) int32   { return 0 }
func AInputEvent_getSource(event *AInputEvent) int32  { return 0 }
func AKeyEvent_getAction(event *AInputEvent) int32    { return 0 }
func AKeyEvent_getKeyCode(event *AInputEvent) int32   { return 0 }
func AKeyEvent_getRepeatCount(event *AInputEvent) int32 { return 0 }
func AMotionEvent_getAction(event *AInputEvent) int32 { return 0 }
func AMotionEvent_getPointerCount(event *AInputEvent) int32 { return 0 }
func AMotionEvent_getX(event *AInputEvent, pointerIndex int64) float32 { return 0 }
func AMotionEvent_getY(event *AInputEvent, pointerIndex int64) float32 { return 0 }
func AMotionEvent_getPressure(event *AInputEvent, pointerIndex int64) float32 { return 0 }

// --- Input queue functions ---
func AInputQueue_detachLooper(queue *AInputQueue)                                  {}
func AInputQueue_hasEvents(queue *AInputQueue) int32                               { return 0 }
func AInputQueue_getEvent(queue *AInputQueue, outEvent **AInputEvent) int32        { return 0 }
func AInputQueue_finishEvent(queue *AInputQueue, event *AInputEvent, handled int32) {}

var _ = unsafe.Pointer(nil)
