// Example: input event type system and processing pattern.
//
// Prints all enum types in the ndk/input package -- event types, key actions,
// input sources, and motion actions -- together with their integer values and
// human-readable string representations. This demonstrates the full set of
// constants that classify input events on Android.
//
// Input events flow through the system in four steps:
//
//  1. The Activity receives an AInputQueue from the framework.
//  2. The Queue is polled for pending events (Queue.HasEvents).
//  3. Each Event is inspected: Event.Type tells you whether it is a Key or
//     Motion event, and the remaining accessors (KeyAction, KeyCode,
//     MotionAction, X, Y, Pressure, etc.) provide the details.
//  4. After processing, the event must be returned to the system with
//     Queue.FinishEvent so Android knows whether the app handled it.
//
// Actually receiving events requires an AInputQueue obtained from an Android
// Activity, so this example focuses on the type system and documents the
// processing pattern in code comments.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/AndroidGoLab/ndk/input"
)

func main() {
	// -- Event types -------------------------------------------------------
	// EventType distinguishes the two broad categories of input events.
	// Key events come from physical or virtual keys; motion events come from
	// touch, mouse, or joystick input.
	fmt.Println("Event types:")
	for _, et := range []input.EventType{input.Key, input.Motion} {
		fmt.Printf("  %-8s = %d\n", et, int32(et))
	}
	fmt.Println()

	// -- Key actions -------------------------------------------------------
	// KeyAction tells you what phase a key event is in. Down fires when the
	// key is first pressed, Up when released, and Multiple for repeated
	// character generation (soft keyboard).
	fmt.Println("Key actions:")
	for _, ka := range []input.KeyAction{input.Down, input.Up, input.Multiple} {
		fmt.Printf("  %-10s = %d\n", ka, int32(ka))
	}
	fmt.Println()

	// -- Input sources -----------------------------------------------------
	// Source identifies the hardware class that produced the event.
	// Keyboard covers both physical keys and soft input methods.
	// Touchscreen is the primary source for finger-driven motion events.
	// Mouse covers USB or Bluetooth mice (and trackpads in mouse mode).
	// Joystick covers game controllers.
	fmt.Println("Input sources:")
	for _, src := range []input.Source{
		input.Keyboard,
		input.Touchscreen,
		input.Mouse,
		input.Joystick,
	} {
		fmt.Printf("  %-14s = %d\n", src, int32(src))
	}
	fmt.Println()

	// -- Motion actions ----------------------------------------------------
	// MotionAction describes the gesture phase. ActionDown / ActionUp mark
	// the first and last finger touching the screen. ActionMove fires as the
	// finger moves. ActionCancel means the gesture was cancelled by the
	// system. ActionPointerDown / ActionPointerUp track secondary fingers in
	// a multi-touch gesture.
	fmt.Println("Motion actions:")
	for _, ma := range []input.MotionAction{
		input.ActionDown,
		input.ActionUp,
		input.ActionMove,
		input.ActionCancel,
		input.ActionPointerDown,
		input.ActionPointerUp,
	} {
		fmt.Printf("  %-20s = %d\n", ma, int32(ma))
	}
	fmt.Println()

	// -- Processing pattern (pseudo-code) ----------------------------------
	// The code below illustrates the typical event loop. It cannot run
	// without a real AInputQueue, so it is guarded by a false constant.
	//
	//   // Assume queue is an *input.Queue obtained from the Activity.
	//   if err := queue.HasEvents(); err == nil {
	//       // Retrieve the next event (via looper poll, not shown).
	//       var event *input.Event
	//       eventType := input.EventType(event.Type())
	//
	//       switch eventType {
	//       case input.Key:
	//           action := input.KeyAction(event.KeyAction())
	//           code := event.KeyCode()
	//           repeat := event.RepeatCount()
	//           fmt.Printf("key: action=%s code=%d repeat=%d\n", action, code, repeat)
	//
	//       case input.Motion:
	//           action := input.MotionAction(event.MotionAction())
	//           source := input.Source(event.Source())
	//           count := event.PointerCount()
	//           for i := int64(0); i < int64(count); i++ {
	//               x, y := event.X(i), event.Y(i)
	//               pressure := event.Pressure(i)
	//               fmt.Printf("motion: action=%s source=%s pointer=%d (%.1f, %.1f) pressure=%.2f\n",
	//                   action, source, i, x, y, pressure)
	//           }
	//       }
	//
	//       // Return the event to the system. Pass 1 if handled, 0 otherwise.
	//       queue.FinishEvent(event, 1)
	//   }
	//
	//   // When the Activity's input queue is being destroyed:
	//   queue.DetachLooper()
	fmt.Println("Event processing pattern printed above in source comments.")
}
