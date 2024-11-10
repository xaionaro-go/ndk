package main

import (
	"log"
	"strconv"

	"github.com/gooid/gl/egl"
	"github.com/xaionaro-go/ndk"
	"github.com/xaionaro-go/ndk/audio/al"
	"github.com/xaionaro-go/ndk/input"
)

func main() {
	context := ndk.Callbacks{
		WindowDraw:         drawAudio,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
		Create: func(act *ndk.Activity, _ []byte) {
			create()
			// 需手动 Load libopenal.so
			libs := ndk.FindMatchLibrary("libopenal*.so")
			if len(libs) > 0 {
				al.InitPath(libs[0])
			}
		},
	}
	ndk.SetMainCB(func(ctx *ndk.Context) {
		ctx.Debug(true)
		ctx.Run(context)
	})
	for ndk.Loop() {
	}
	log.Println("done.")
}

func event(act *ndk.Activity, e *ndk.InputEvent) {
	if mot := e.Motion(); mot != nil {
		lastX = int(float32(mot.GetX(0)) / WINDOWSCALE)
		lastY = int(float32(mot.GetY(0)) / WINDOWSCALE)
		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			mouseLeft = false
			//log.Println("event:", mot)
		case input.MOTION_EVENT_ACTION_DOWN:
			mouseLeft = true
			//log.Println("event:", mot)
		case input.MOTION_EVENT_ACTION_MOVE:

		default:
			//log.Println("event:", mot)
			return
		}
		drawAudio(act, nil)

		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			lastX, lastY = 0, 0
			drawAudio(act, nil)
		}
	}
}

var curAct *ndk.Activity
var eglctx *egl.EGLContext

func redraw(act *ndk.Activity, win *ndk.Window) {
	curAct = act
	act.Context().Call(func() {
		releaseEGL()

		width, height = win.Width(), win.Height()
		eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
		if eglctx == nil {
			return
		}
		initEGL()
	}, false)
}

func destroyed(act *ndk.Activity, win *ndk.Window) {
	releaseEGL()
}

func drawAudio(act *ndk.Activity, win *ndk.Window) {
	if eglctx != nil {
		draw()
	}
}

const RECORDPATH = "/sdcard/records"

func getDensity() int {
	dstr := ndk.PropGet("hw.lcd.density")
	if dstr == "" {
		dstr = ndk.PropGet("qemu.sf.lcd_density")
	}

	log.Println(" lcd_density:", dstr)
	if dstr != "" {
		density, _ := strconv.Atoi(dstr)
		return density
	}
	return density
}
func SwapBuffers() { eglctx.SwapBuffers() }
func Wake()        { curAct.Context().Wake() }
