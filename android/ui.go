//go:build android

package android

import (
	"github.com/xaionaro-go/ndk/jni"
)

// ShowToast displays an Android Toast message.
func (app *App) ShowToast(message string) {
	app.mu.Lock()
	act := app.activity
	app.mu.Unlock()

	if act == nil {
		return
	}
	jni.ShowToast(act.Pointer(), message)
}

// FillWindowColor fills the current window with a solid RGBA color.
// Does nothing if no window is available.
func (app *App) FillWindowColor(color uint32) {
	app.mu.Lock()
	wPtr := app.windowPtr
	app.mu.Unlock()

	if wPtr == nil {
		return
	}
	jni.FillWindowColor(wPtr, color)
}
