// Example: obtain the Android Choreographer singleton.
//
// AChoreographer is the central timing coordinator for frame production on
// Android. The display subsystem raises a vsync interrupt at the panel's
// refresh rate (typically 60, 90, or 120 Hz). AChoreographer exposes that
// signal to user-space so that game loops and rendering engines can start
// each frame's work at the optimal moment, avoiding jank caused by
// submitting frames out of phase with the display.
//
// A typical usage pattern in a game engine:
//
//  1. Obtain the singleton with GetInstance (one per Looper thread).
//  2. Register a frame callback via postFrameCallback.
//  3. Inside the callback, compute the delta from the previous vsync
//     timestamp, advance simulation, record draw commands, and re-register
//     the callback for the next frame.
//
// Because the choreographer is a process-wide singleton tied to the calling
// thread's Looper, it does not require explicit cleanup -- there is no
// Close method.
//
// This program must run on an Android device with NDK choreographer support.
package main

import (
	"log"

	"github.com/xaionaro-go/ndk/choreographer"
)

func main() {
	// Obtain the per-thread Choreographer singleton. On Android this
	// returns the AChoreographer associated with the current thread's
	// Looper. The returned handle is valid for the lifetime of the
	// thread and must not be freed by the caller.
	ch := choreographer.GetInstance()
	if ch == nil {
		log.Fatal("failed to get choreographer instance")
	}

	log.Println("choreographer instance acquired")
	log.Println("the application is now ready to post vsync-aligned frame callbacks")
}
