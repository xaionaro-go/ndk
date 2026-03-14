//go:build ignore

package gameactivity

/*
#include "include/game-activity/GameActivity.h"

static GameActivityPointerAxes* getPointer(
	GameActivityMotionEvent* event, int index
) {
	return &event->pointers[index];
}
*/
import "C"

const (
	// MaxPointersInMotionEvent is the maximum number of pointers in a motion event.
	MaxPointersInMotionEvent = C.GAMEACTIVITY_MAX_NUM_POINTERS_IN_MOTION_EVENT
)

// MotionEvent represents a motion (touch/mouse) event from GameActivity.
// It maps 1:1 to the Java MotionEvent class.
type MotionEvent struct {
	DeviceID       int32
	Source         int32
	Action         int32
	ActionButton   int32
	Flags          int32
	MetaState      int32
	EdgeFlags      int32
	ButtonState    int32
	Classification int32
	PrecisionX     float32
	PrecisionY     float32
	DownTime       int64
	EventTime      int64
	Pointers       []PointerAxes
}

// motionEventFromC converts a C GameActivityMotionEvent to a Go MotionEvent.
func motionEventFromC(c *C.GameActivityMotionEvent) MotionEvent {
	count := int(c.pointerCount)
	pointers := make([]PointerAxes, count)
	for i := range count {
		p := C.getPointer(c, C.int(i))
		pointers[i] = pointerAxesFromC(p)
	}
	return MotionEvent{
		DeviceID:       int32(c.deviceId),
		Source:         int32(c.source),
		Action:         int32(c.action),
		ActionButton:   int32(c.actionButton),
		Flags:          int32(c.flags),
		MetaState:      int32(c.metaState),
		EdgeFlags:      int32(c.edgeFlags),
		ButtonState:    int32(c.buttonState),
		Classification: int32(c.classification),
		PrecisionX:     float32(c.precisionX),
		PrecisionY:     float32(c.precisionY),
		DownTime:       int64(c.downTime),
		EventTime:      int64(c.eventTime),
		Pointers:       pointers,
	}
}
