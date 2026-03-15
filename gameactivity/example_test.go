//go:build ignore

package gameactivity_test

import (
	"fmt"
	"unsafe"

	"github.com/AndroidGoLab/ndk/gameactivity"
)

func Example() {
	// In a real application, the GameActivity pointer comes from the
	// framework via GameActivity_onCreate.
	var actPtr unsafe.Pointer
	act := gameactivity.NewActivityFromPointer(actPtr)

	// Register lifecycle and input callbacks.
	gameactivity.SetCallbacks(act, &gameactivity.Callbacks{
		OnResume: func(activity *gameactivity.Activity) {
			fmt.Println("resumed")
		},
		OnTouchEvent: func(activity *gameactivity.Activity, event gameactivity.MotionEvent) bool {
			for _, p := range event.Pointers {
				fmt.Printf("touch pointer %d at (%.0f, %.0f)\n", p.ID, p.X(), p.Y())
			}
			return true
		},
	})

	// Show the software keyboard.
	act.ShowSoftInput(gameactivity.ShowSoftInputImplicit)

	// Update text input state.
	act.SetTextInputState(gameactivity.TextInputState{
		Text:           "Hello",
		SelectionStart: 5,
		SelectionEnd:   5,
	})

	// Make the activity fullscreen.
	act.SetWindowFlags(gameactivity.FlagFullscreen, 0)

	fmt.Println("GameActivity initialized")
}
