//go:build android

package android

import (
	"sync"
	"unsafe"

	"github.com/xaionaro-go/ndk/activity"
	"github.com/xaionaro-go/ndk/window"
)

// App provides a high-level interface to an Android NativeActivity.
type App struct {
	mu        sync.Mutex
	activity  *activity.Activity
	window    *window.Window
	windowPtr unsafe.Pointer
}

// Activity returns the underlying Activity.
func (app *App) Activity() *activity.Activity {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.activity
}

// Window returns the current native window, or nil if no window is available.
func (app *App) Window() *window.Window {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.window
}

// WindowSize returns the current window dimensions.
// Returns (0, 0) if no window is available.
func (app *App) WindowSize() (width, height int32) {
	app.mu.Lock()
	w := app.window
	app.mu.Unlock()

	if w == nil {
		return 0, 0
	}
	return w.Width(), w.Height()
}

func (app *App) setActivity(act *activity.Activity) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.activity = act
}

func (app *App) setWindow(ptr unsafe.Pointer) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if ptr == nil {
		app.window = nil
		app.windowPtr = nil
	} else {
		app.windowPtr = ptr
		app.window = window.NewWindowFromPointer(ptr)
	}
}
