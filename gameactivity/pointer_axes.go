//go:build ignore

package gameactivity

/*
#include "include/game-activity/GameActivity.h"
*/
import "C"

const (
	// PointerAxisCount is the number of axis values per pointer.
	PointerAxisCount = C.GAME_ACTIVITY_POINTER_INFO_AXIS_COUNT
)

// PointerAxes describes a single pointer within a MotionEvent.
type PointerAxes struct {
	// ID is the pointer identifier.
	ID int32

	// RawX is the raw X coordinate.
	RawX float32

	// RawY is the raw Y coordinate.
	RawY float32

	// AxisValues contains per-axis values; indexed by AMOTION_EVENT_AXIS_* constants.
	AxisValues [PointerAxisCount]float32
}

// X returns the X coordinate of this pointer.
func (p *PointerAxes) X() float32 {
	return p.AxisValues[0] // AMOTION_EVENT_AXIS_X
}

// Y returns the Y coordinate of this pointer.
func (p *PointerAxes) Y() float32 {
	return p.AxisValues[1] // AMOTION_EVENT_AXIS_Y
}

// pointerAxesFromC converts a C GameActivityPointerAxes to a Go PointerAxes.
func pointerAxesFromC(c *C.GameActivityPointerAxes) PointerAxes {
	var axes [PointerAxisCount]float32
	for i := range PointerAxisCount {
		axes[i] = float32(c.axisValues[i])
	}
	return PointerAxes{
		ID:         int32(c.id),
		RawX:       float32(c.rawX),
		RawY:       float32(c.rawY),
		AxisValues: axes,
	}
}
