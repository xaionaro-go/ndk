package main

import (
	"log"

	"github.com/xaionaro-go/ndk"
	"github.com/xaionaro-go/ndk/input"
)

func main() {
	context := ndk.Callbacks{
		Event: event,
	}
	ndk.SetMainCB(func(ctx *ndk.Context) {
		ctx.Run(context)
	})
	for ndk.Loop() {
	}
	log.Println("done.")
}

var bWallpaper = false

func event(act *ndk.Activity, e *ndk.InputEvent) {
	if mot := e.Motion(); mot != nil {
		if mot.GetAction() == input.MOTION_EVENT_ACTION_UP {
			log.Println("event:", mot)

			bWallpaper = !bWallpaper
			if bWallpaper {
				act.SetWindowFlags(ndk.FLAG_SHOW_WALLPAPER, 0)
			} else {
				act.SetWindowFlags(0, ndk.FLAG_SHOW_WALLPAPER)
			}
		}
	}
}
