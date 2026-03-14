//go:build android

package android

import (
	"unsafe"

	"github.com/xaionaro-go/ndk/activity"
	"github.com/xaionaro-go/ndk/input"
)

// Handlers defines callbacks for Android lifecycle events and input.
// All callbacks receive the App instance for accessing platform APIs.
type Handlers struct {
	OnCreate             func(app *App)
	OnResume             func(app *App)
	OnPause              func(app *App)
	OnDestroy            func(app *App)
	OnWindowCreated      func(app *App)
	OnWindowDestroyed    func(app *App)
	OnWindowFocusChanged func(app *App, hasFocus bool)
	OnInputEvent         func(app *App, event *input.Event) bool
}

// Run sets up NativeActivity lifecycle callbacks and manages the App instance.
// Call this from an init() function in your main package.
func Run(handlers Handlers) {
	app := &App{}

	activity.SetLifecycleCallbacks(activity.LifecycleCallbacks{
		OnCreate: func(act *activity.Activity) {
			app.setActivity(act)
			if handlers.OnCreate != nil {
				handlers.OnCreate(app)
			}
		},
		OnNativeWindowCreated: func(act *activity.Activity, win unsafe.Pointer) {
			app.setActivity(act)
			app.setWindow(win)
			if handlers.OnWindowCreated != nil {
				handlers.OnWindowCreated(app)
			}
		},
		OnResume: func(act *activity.Activity) {
			app.setActivity(act)
			if handlers.OnResume != nil {
				handlers.OnResume(app)
			}
		},
		OnPause: func(act *activity.Activity) {
			if handlers.OnPause != nil {
				handlers.OnPause(app)
			}
		},
		OnWindowFocusChanged: func(_ *activity.Activity, hasFocus int32) {
			if handlers.OnWindowFocusChanged != nil {
				handlers.OnWindowFocusChanged(app, hasFocus != 0)
			}
		},
		OnNativeWindowDestroyed: func(_ *activity.Activity, _ unsafe.Pointer) {
			if handlers.OnWindowDestroyed != nil {
				handlers.OnWindowDestroyed(app)
			}
			app.setWindow(nil)
		},
		OnInputQueueCreated: func(_ *activity.Activity, queuePtr unsafe.Pointer) {
			if handlers.OnInputEvent != nil {
				q := input.NewQueueFromPointer(queuePtr)
				go drainInput(app, q, handlers.OnInputEvent)
			}
		},
		OnDestroy: func(_ *activity.Activity) {
			if handlers.OnDestroy != nil {
				handlers.OnDestroy(app)
			}
			app.setActivity(nil)
			app.setWindow(nil)
		},
	})
}

func drainInput(
	app *App,
	q *input.Queue,
	handler func(app *App, event *input.Event) bool,
) {
	for {
		ev := q.GetEvent()
		if ev == nil {
			return
		}

		if q.PreDispatchEvent(ev) {
			continue
		}

		handled := handler(app, ev)
		var result int32
		if handled {
			result = 1
		}
		q.FinishEvent(ev, result)
	}
}
