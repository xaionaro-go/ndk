package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/xaionaro-go/ndk"
	"github.com/xaionaro-go/ndk/storage"
)

type nativeInfo struct {
	win *ndk.Window
}

func (e *nativeInfo) NativeDisplay() unsafe.Pointer {
	return nil
}
func (e *nativeInfo) NativeWindow() unsafe.Pointer {
	return unsafe.Pointer(e.win)
}
func (e *nativeInfo) WindowSize() (w, h int) {
	return e.win.Width(), e.win.Height()
}

func (e *nativeInfo) SetBuffersGeometry(format int) int {
	return e.win.SetBuffersGeometry(0, 0, format)
}

// create
func preCreate(ctx interface{}) {
	//act := ctx.(*ndk.Activity)
}

func postCreate(ctx interface{}) {
	act := ctx.(*ndk.Activity)
	libPath := ""
	ls := ndk.FindMatchLibrary("libopencv_*.so")
	if len(ls) > 0 {
		libPath = ls[0]
	}

	cascadePath := loadCascade(act)
	faceDetect.LoadOpenCV(libPath, cascadePath)
}

// Load CascadeClassifier, but since there is only an interface for loading files,
// it can only be read from the asset and loaded by generating a temporary file.
func loadCascade(act *ndk.Activity) string {
	assetMgr := act.AssetManager()
	altname := "haarcascade_frontalface_alt2.xml"
	tmpPath := filepath.Join("/data/data", act.Context().Package(), "files")
	os.MkdirAll(tmpPath, 0666)

	fname := filepath.Join(tmpPath, altname)
	fw, err := os.Create(fname)
	if err != nil {
		tmpPath = `/sdcard`
		fw, err = os.Create(filepath.Join(tmpPath, altname))
		if err != nil {
			panic(err)
		}
	}

	r := assetMgr.Open(altname, storage.ASSET_MODE_STREAMING)
	if r == nil {
		panic(fmt.Errorf("open ", altname, " fail"))
	}
	io.Copy(fw, r)
	r.Close()
	fw.Close()

	return fname
}
