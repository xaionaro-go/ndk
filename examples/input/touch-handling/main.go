// Touch and key input event handling example.
//
// Demonstrates the complete input event processing pipeline using the
// ndk/input package. On Android, input events (touches, key presses)
// are delivered through an AInputQueue that the framework attaches to
// the Activity. This example shows:
//
//   - How to wrap a raw AInputQueue pointer into an input.Queue.
//   - The full event polling loop: HasEvents -> GetEvent ->
//     PreDispatchEvent -> handle -> FinishEvent.
//   - Dispatching key events (action, keycode, repeat count).
//   - Dispatching motion/touch events (action, pointer count,
//     per-pointer x/y coordinates and pressure).
//   - How to integrate the input queue with an ALooper for
//     continuous event-driven processing (documented in comments).
//
// The event loop runs against a nil queue pointer in this example, so
// HasEvents returns an error immediately and no events are delivered.
// In a real NativeActivity the queue pointer comes from the framework
// via the onInputQueueCreated callback.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"github.com/xaionaro-go/ndk/input"
)

// looperIdent is the identifier returned by ALooper_pollOnce when the
// input queue has pending events. Any positive integer works; this
// value must match the ident passed to AInputQueue_attachLooper.
const looperIdent = 2

func main() {
	// Lock to the OS thread. Both ALooper and AInputQueue are
	// thread-local on Android, so all calls must happen on the same
	// OS thread that prepared the looper.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// ---------------------------------------------------------------
	// Step 1: Obtain the input queue
	//
	// In a NativeActivity the framework delivers an AInputQueue*
	// through the onInputQueueCreated callback. Wrap it with
	// input.NewQueueFromPointer to get an *input.Queue.
	//
	// For this example we use a nil pointer. The first HasEvents call
	// will return an error, which causes the loop to exit cleanly.
	// ---------------------------------------------------------------

	var rawQueuePtr unsafe.Pointer // would come from onInputQueueCreated
	queue := input.NewQueueFromPointer(rawQueuePtr)

	fmt.Println("=== Input event handling example ===")
	fmt.Println()
	fmt.Println("Queue created from pointer:", queue.Pointer())
	fmt.Println()

	// ---------------------------------------------------------------
	// Step 2: Looper integration (documented pattern)
	//
	// Before polling for events, the queue must be attached to an
	// ALooper so that the looper wakes up when input is available.
	// The idiomatic looper package provides Prepare and PollOnce:
	//
	//   import "github.com/xaionaro-go/ndk/looper"
	//
	//   lp := looper.Prepare(1) // ALOOPER_PREPARE_ALLOW_NON_CALLBACKS
	//   defer lp.Close()
	//
	// Attaching the queue to the looper is done at the capi level
	// (the idiomatic binding does not yet expose AttachLooper):
	//
	//   capi.AInputQueue_attachLooper(
	//       (*capi.AInputQueue)(queue.Pointer()),
	//       (*capi.ALooper)(lp.Pointer()),
	//       looperIdent, // returned by PollOnce when events are ready
	//       nil,         // no C callback; use poll-based reading
	//       nil,         // no user data
	//   )
	//
	// Then the main loop becomes:
	//
	//   for {
	//       ident := looper.PollOnce(-1, nil, nil, nil)
	//       if ident == looperIdent {
	//           processInputEvents(queue)
	//       }
	//   }
	//
	// When the Activity's input queue is being destroyed
	// (onInputQueueDestroyed), detach and stop polling:
	//
	//   queue.DetachLooper()
	// ---------------------------------------------------------------

	// ---------------------------------------------------------------
	// Step 3: Event polling loop
	//
	// processInputEvents drains all pending events from the queue.
	// In a real app this function is called each time the looper
	// returns the input queue's ident.
	// ---------------------------------------------------------------

	if queue.Pointer() == nil {
		fmt.Println("No real AInputQueue available (nil pointer).")
		fmt.Println("In a NativeActivity, the queue comes from onInputQueueCreated.")
	} else {
		processInputEvents(queue)
	}

	fmt.Println()
	fmt.Println("See main.go source for the full processing pattern.")
}

// processInputEvents drains all pending input events from the queue,
// dispatches them based on type (key or motion), and returns each
// event to the system.
func processInputEvents(queue *input.Queue) {
	for {
		// Check whether the queue has pending events.
		if err := queue.HasEvents(); err != nil {
			log.Printf("HasEvents: %v (no more events or queue unavailable)", err)
			return
		}

		// Retrieve the next event. GetEvent returns nil when the
		// internal AInputQueue_getEvent call fails (e.g. no events).
		event := queue.GetEvent()
		if event == nil {
			log.Println("GetEvent returned nil, stopping")
			return
		}

		// Pre-dispatch gives the system (e.g. the IME) a chance to
		// consume the event before the app sees it. If it returns
		// true the event was consumed and must NOT be passed to
		// FinishEvent.
		if queue.PreDispatchEvent(event) {
			log.Println("event consumed by pre-dispatch (e.g. IME)")
			continue
		}

		// Dispatch based on event type.
		handled := handleEvent(event)

		// Return the event to the system. Pass 1 if the app handled
		// it, 0 to let the system perform the default action.
		var handledFlag int32
		if handled {
			handledFlag = 1
		}
		queue.FinishEvent(event, handledFlag)
	}
}

// handleEvent inspects the event type and delegates to the
// appropriate handler. Returns true if the event was handled.
func handleEvent(event *input.Event) bool {
	eventType := input.EventType(event.Type())
	source := input.Source(event.Source())

	switch eventType {
	case input.Key:
		return handleKeyEvent(event, source)

	case input.Motion:
		return handleMotionEvent(event, source)

	default:
		fmt.Printf("unknown event type: %s (source=%s)\n", eventType, source)
		return false
	}
}

// handleKeyEvent processes a key press/release event.
func handleKeyEvent(event *input.Event, source input.Source) bool {
	action := input.KeyAction(event.KeyAction())
	keyCode := event.KeyCode()
	repeatCount := event.RepeatCount()

	fmt.Printf("KEY event: action=%-8s keycode=%-4d repeat=%-2d source=%s\n",
		action, keyCode, repeatCount, source)

	// Example: treat Down and Up as handled, ignore Multiple.
	switch action {
	case input.Down:
		fmt.Printf("  -> key %d pressed\n", keyCode)
		return true

	case input.Up:
		fmt.Printf("  -> key %d released\n", keyCode)
		return true

	case input.Multiple:
		fmt.Printf("  -> key %d repeated %d times\n", keyCode, repeatCount)
		return true

	default:
		return false
	}
}

// handleMotionEvent processes touch, mouse, or stylus motion events.
func handleMotionEvent(event *input.Event, source input.Source) bool {
	action := input.MotionAction(event.MotionAction())

	// The raw action value encodes both the action and (for
	// pointer-down/up) the pointer index in the upper bits.
	// Mask out the action to get the gesture phase.
	maskedAction := input.MotionAction(int32(action) & int32(input.ActionMask))

	pointerCount := event.PointerCount()

	fmt.Printf("MOTION event: action=%-18s pointers=%-2d source=%s\n",
		maskedAction, pointerCount, source)

	// Print per-pointer coordinates and pressure.
	for i := uint64(0); i < pointerCount; i++ {
		x := event.X(i)
		y := event.Y(i)
		pressure := event.Pressure(i)
		fmt.Printf("  pointer[%d]: x=%8.2f  y=%8.2f  pressure=%.3f\n",
			i, x, y, pressure)
	}

	switch maskedAction {
	case input.ActionDown:
		fmt.Println("  -> first finger touched the screen")

	case input.ActionUp:
		fmt.Println("  -> last finger lifted from the screen")

	case input.ActionMove:
		fmt.Println("  -> finger(s) moved")

	case input.ActionCancel:
		fmt.Println("  -> gesture cancelled by the system")

	case input.ActionPointerDown:
		fmt.Println("  -> additional finger touched the screen")

	case input.ActionPointerUp:
		fmt.Println("  -> non-primary finger lifted")
	}

	return true
}
