// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ndk

import (
	"os"
	"runtime"
	"sync"
	"unsafe"
)

// Callbacks is the set of functions called by the Activity.
type Callbacks struct {
	Create               func(*Activity, []byte)
	Start                func(*Activity)
	Stop                 func(*Activity)
	Pause                func(*Activity)
	Resume               func(*Activity)
	Destroy              func(*Activity)
	SaveState            func(*Activity) []byte
	FocusChanged         func(*Activity, bool)
	ContentRectChanged   func(*Activity, *Rect)
	ConfigurationChanged func(*Activity)
	LowMemory            func(*Activity)

	// Window
	WindowCreated      func(*Activity, *Window)
	WindowDestroyed    func(*Activity, *Window)
	WindowDraw         func(*Activity, *Window)
	WindowResized      func(*Activity, *Window)
	WindowRedrawNeeded func(*Activity, *Window)

	// Touch is called by the app when a touch event occurs.
	Event func(*Activity, *InputEvent)
	// Sensor
	Sensor func(*Activity, []SensorEvent)
}

type Context struct {
	Callbacks

	Act    *Activity
	Window *Window // onNativeWindowCreated may not be called after onStart
	Input  *InputQueue

	IsResume,
	IsFocus, // onWindowFocusChanged may not be called after onStart/onStop
	WillDestoryValue bool

	Looper      *Looper
	read, write *os.File
	sensorQueue *SensorEventQueue
	FuncChan    chan func()
	wait        *sync.WaitGroup
	isReady     bool
	isDebug     bool

	ClassName   string
	packageName string
	SavedState  []byte
}

func (ctx *Context) init() {
	ctx.FuncChan = make(chan func(), 3)
	ctx.wait = &sync.WaitGroup{}
}

func (ctx *Context) setCB(cb *Callbacks) {
	if cb != nil {
		ctx.Callbacks = *cb
	}
}

func (ctx *Context) Reset() {
	ctx.IsResume = false
	ctx.WillDestoryValue = false
}

func (ctx *Context) Debug(enable bool) {
	ctx.isDebug = enable
}

func (ctx *Context) WillDestory() bool {
	return ctx.WillDestoryValue
}

const (
	Looper_ID_MAIN   = 1
	Looper_ID_INPUT  = 2
	Looper_ID_SENSOR = 3
)

// Init looper
func (ctx *Context) initLooper() {
	var err error
	ctx.read, ctx.write, err = os.Pipe()
	Assert(err)
	ctx.Looper = LooperPrepare(LOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	ctx.Looper.AddFd(int(ctx.read.Fd()), Looper_ID_MAIN, LOOPER_EVENT_INPUT, nil, unsafe.Pointer(uintptr(Looper_ID_MAIN)))
}

func (ctx *Context) begin(cb Callbacks) {
	ctx.setCB(&cb)
	ctx.initLooper()
	ctx.FuncChan <- func() {}
}

// Run starts the Activity.
//
// It must be called directly from from the main function and will
// block until the app exits.
func (ctx *Context) Run(cb Callbacks) {
	Info("ctx.Run")
	ctx.begin(cb)
	ctx.Loop()
}

func (ctx *Context) Begin(cb Callbacks) {
	Info("ctx.Begin")
	ctx.begin(cb)
}

func (ctx *Context) Release() {
	ctx.Looper.RemoveFd(int(ctx.read.Fd()))
	if ctx.sensorQueue != nil {
		SensorManagerInstance().Destroy(ctx.sensorQueue)
	}
	//ctx.looper.Release()
	ctx.read.Close()
	ctx.write.Close()
}

func (ctx *Context) RunFunc(fun func(), sync bool) {
	if !ctx.isDebug {
		defer func() {
			if r := recover(); r != nil {
				Info("RunFunc.recover:", r)
			}
		}()
	}

	//Info("...wait...")
	//oldM := util.Mode(util.StackTrace)
	//Assert(false)
	//util.Mode(oldM)

	var cmd = []byte{Looper_ID_MAIN}
	ctx.write.Write(cmd)
	if !sync {
		ctx.FuncChan <- fun
	} else {
		ctx.wait.Add(1)
		ctx.FuncChan <- func() { defer ctx.wait.Done(); fun() }
		ctx.wait.Wait()
	}
	//Info("...wait<<<", forceEsace)
}

var entryDoFunc uintptr

func (ctx *Context) doFunc() {
	if !ctx.isDebug {
		defer func() {
			if r := recover(); r != nil {
				Info("execFunc.recover:", r)
			}
		}()
	}

	if entryDoFunc == 0 {
		pc := make([]uintptr, 1)
		n := runtime.Callers(1, pc)
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()
		entryDoFunc = frame.Entry
	}

	(<-ctx.FuncChan)()
}

func (ctx *Context) Call(fun func(), sync bool) {
	if sync {
		// If it is called in doFunc, it will be executed directly, otherwise it will cause deadlock.
		pc := make([]uintptr, 10000)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])
		for {
			frame, more := frames.Next()
			if frame.Entry == entryDoFunc {
				fun()
				return
			}
			if !more {
				break
			}
		}
	}

	ctx.RunFunc(fun, sync)
}

func (ctx *Context) loopEvents(timeoutMillis int) bool {
	Assert(!ctx.WillDestoryValue, "!ctx.willDestory")

	for {
		ident, _, event, data := LooperPollOnce(timeoutMillis)
		//Info("LooperPollAll is", int(ident))
		switch ident {
		case Looper_ID_MAIN:
			//Info("LooperPollAll is MAIN")
			var cmd = []byte{0}
			n, err := ctx.read.Read(cmd)
			Assert(n == 1 && err == nil)
			Assert(ident == int(cmd[0]))
			Assert(ident == int(data))
			Assert(event == int(data))
			ctx.doFunc()
			if ctx.WillDestoryValue {
				return false
			}
			return true

		case Looper_ID_INPUT:
			//Info("LooperPollAll is INPUT")
			for ctx.Input != nil {
				event := ctx.Input.GetEvent()
				if event == nil {
					break
				}
				if !ctx.Input.PreDispatchEvent(event) {
					if !ctx.WillDestoryValue {
						ctx.processEvent(event)
					}
					ctx.Input.FinishEvent(event, 0)
				}
			}
			if ctx.WillDestoryValue {
				return false
			}
			return true

		case Looper_ID_SENSOR:
			//Info("LooperPollAll is SENSOR")
			if ctx.sensorQueue != nil {
				var events []SensorEvent
				errno := ctx.sensorQueue.HasEvents()
				if errno > 0 {
					events, errno = ctx.sensorQueue.GetEvents(32)
				}

				if errno < 0 {
					Info("SensorEventQueue GetEvents error :", errno)
				} else {
					if ctx.Sensor != nil && len(events) > 0 {
						//for _, event := range events {
						//}
						ctx.Sensor(ctx.Act, events)
					}
				}
			} else {
				Info("looper_ID_SENSOR ...")
			}
			return true

		case LOOPER_POLL_WAKE:
			//info("LooperPollAll is ALOOPER_POLL_WAKE")
			return true
		case LOOPER_POLL_TIMEOUT:
			//Info("LooperPollAll is LOOPER_POLL_TIMEOUT")
			return true
		case LOOPER_POLL_ERROR:
			Info("LooperPollAll is LOOPER_POLL_ERROR")
			return true
		default:
			Info("LooperPollAll is ", ident)
			return true
		}
	}
}

func (ctx *Context) processEvent(e *InputEvent) {
	//Info("processEvent:", e)
	if ctx.Event != nil {
		ctx.Event(ctx.Act, e)
	}
}

func (ctx *Context) pollEvent(timeoutMillis int) bool {
	if ctx.loopEvents(timeoutMillis) {
		if ctx.IsFocus && ctx.IsResume {
			if ctx.Window != nil &&
				ctx.WindowDraw != nil {
				ctx.WindowDraw(ctx.Act, ctx.Window)
			}
		}
		return true
	} else {
		//info("pollEvent: return false")
		ctx.doFunc()
		return false
	}
}

func (ctx *Context) Loop() {
	for ctx.pollEvent(-1) {
		// nothing
	}
}

func (ctx *Context) PollEvent() bool {
	return ctx.pollEvent(0)
}

func (ctx *Context) WaitEvent() bool {
	return ctx.pollEvent(-1)
}

func (ctx *Context) Wake() {
	ctx.Looper.Wake()
}

func (ctx *Context) Name() string {
	return ctx.ClassName
}

func (ctx *Context) Package() string {
	return ctx.packageName
}
