//go:build ignore

package gameactivity

/*
#include "include/game-activity/GameActivity.h"
*/
import "C"

// KeyEvent represents a key event from GameActivity.
// It maps 1:1 to the Java KeyEvent class.
type KeyEvent struct {
	DeviceID    int32
	Source      int32
	Action      int32
	Flags       int32
	KeyCode     int32
	MetaState   int32
	Modifiers   int32
	RepeatCount int32
	DownTime    int64
	EventTime   int64
}

// keyEventFromC converts a C GameActivityKeyEvent to a Go KeyEvent.
func keyEventFromC(c *C.GameActivityKeyEvent) KeyEvent {
	return KeyEvent{
		DeviceID:    int32(c.deviceId),
		Source:      int32(c.source),
		Action:      int32(c.action),
		Flags:       int32(c.flags),
		KeyCode:     int32(c.keyCode),
		MetaState:   int32(c.metaState),
		Modifiers:   int32(c.modifiers),
		RepeatCount: int32(c.repeatCount),
		DownTime:    int64(c.downTime),
		EventTime:   int64(c.eventTime),
	}
}
