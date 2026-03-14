//go:build android

package android

import (
	"github.com/xaionaro-go/ndk/jni"
)

// HasPermission checks if a runtime permission is currently granted.
func (app *App) HasPermission(permission string) bool {
	app.mu.Lock()
	act := app.activity
	app.mu.Unlock()

	if act == nil {
		return false
	}
	return jni.HasPermission(act.Pointer(), permission)
}

// RequestPermission shows the system permission dialog for the given permission.
// This is asynchronous; use HasPermission to check the result after the user responds.
func (app *App) RequestPermission(permission string) {
	app.mu.Lock()
	act := app.activity
	app.mu.Unlock()

	if act == nil {
		return
	}
	jni.RequestPermission(act.Pointer(), permission)
}
