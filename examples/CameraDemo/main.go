package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gooid/gl/egl"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
	"github.com/xaionaro-go/ndk"
	"github.com/xaionaro-go/ndk/camera"
	"github.com/xaionaro-go/ndk/input"
)

func main() {
	context := ndk.Callbacks{
		WindowDraw: draw,
		//WindowCreated:      winCreate,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
		Create:             create,
	}
	ndk.SetMainCB(func(ctx *ndk.Context) {
		ctx.Debug(true)
		ctx.Run(context)
	})
	for ndk.Loop() {
	}
	log.Println("done.")
}

var lastTouch time.Time
var mouseLeft = false
var mouseRight = false
var lastX, lastY int

func event(act *ndk.Activity, e *ndk.InputEvent) {
	if mot := e.Motion(); mot != nil {
		lastTouch = time.Now()

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
		draw(act, nil)

		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			lastX, lastY = 0, 0
			draw(act, nil)
		}
	}
}

const WINDOWSCALE = 1
const AUTOHIDETIME = time.Second * 5

const title = "Camera"

var (
	width, height int
	density       = 160
	eglctx        *egl.EGLContext

	im      *util.Render
	cam     *cameraObj
	flashOn bool
)

func initEGL(act *ndk.Activity, win *ndk.Window) {
	width, height = win.Width(), win.Height()
	log.Println("WINSIZE:", width, height)
	width = int(float32(width) / WINDOWSCALE)
	height = int(float32(height) / WINDOWSCALE)

	eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
	if eglctx == nil {
		return
	}

	log.Println("RENDERER:", gl.GoStr(gl.GetString(gl.RENDERER)))
	log.Println("VENDOR:", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Println("VERSION:", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Println("EXTENSIONS:", gl.GoStr(gl.GetString(gl.EXTENSIONS)))

	log.Printf("%s %s", gl.GoStr(gl.GetString(gl.RENDERER)), gl.GoStr(gl.GetString(gl.VERSION)))

	// After setting the font, call
	im.CreateDeviceObjects()

	io := imgui.GetIO()

	dsize := imgui.NewVec2(float32(width), float32(height))
	defer dsize.Delete()
	io.SetDisplaySize(dsize)

	// Setup time step
	io.SetDeltaTime(float32(time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)))

	if cam == nil {
		cam = cameraInit(camera.FACING_BACK, func(w, h int, img []byte) bool {
			act.Context().Wake()
			return false
		})
	}
	cam.ResetRender()
}

func releaseEGL() {
	if im != nil {
		im.DestroyDeviceObjects()
	}
	if eglctx != nil {
		eglctx.Terminate()
		eglctx = nil
	}
}

func create(act *ndk.Activity, _ []byte) {
	// gl init
	gl.Init()

	// Setup Dear ImGui binding
	imgui.CreateContext()

	// Setup style
	imgui.StyleColorsDark()
	//imgui.StyleColorsClassic()

	io := imgui.GetIO()
	fonts := io.GetFonts()

	log.Println("GOOS", runtime.GOOS)
	if runtime.GOOS == "android" {
		fontName := "/system/fonts/DroidSansFallback.ttf"
		if _, err := os.Stat(fontName); err != nil {
			fontName = "/system/fonts/NotoSansCJK-Regular.ttc"
			if _, err := os.Stat(fontName); err != nil {
				fontName = "/system/fonts/DroidSans.ttf"
				if _, err := os.Stat(fontName); err != nil {
					fontName = ""
				}
			}
		}
		if fontName != "" {
			_ = fonts
			// Loading all Chinese Glyphs, but the memory overhead is too high
			//fonts.AddFontFromFileTTF(fontName, 24.0, imgui.SwigcptrFontConfig(0), fonts.GetGlyphRangesChineseSimplifiedCommon())
			// Only load the Chinese Glyphs that need to be displayed
			fonts.AddFontFromFileTTF(fontName, 24.0, imgui.SwigcptrFontConfig(0), util.GetFontGlyphRanges(title))
		}
	}

	if runtime.GOOS == "android" {
		dstr := ndk.PropGet("hw.lcd.density")
		if dstr == "" {
			dstr = ndk.PropGet("qemu.sf.lcd_density")
		}

		log.Println(" lcd_density:", dstr)
		if dstr != "" {
			density, _ = strconv.Atoi(dstr)
		}
	}

	// By adjusting the size of the elements in Style, you can control the display size, but you also need to adjust the font size.
	if density > 160 {
		iScale := float32(density) / 160 / float32(WINDOWSCALE)
		io.SetFontGlobalScale(iScale)
		style := imgui.GetStyle()
		style.ScaleAllSizes(iScale)
	}

	// Controlling display size by scaling the DisplayFramebuffer
	scale := imgui.NewVec2((float32)(WINDOWSCALE), (float32)(WINDOWSCALE))
	defer scale.Delete()
	io.SetDisplayFramebufferScale(scale)

	io.SetConfigFlags(io.GetConfigFlags() | int(imgui.ConfigFlags_IsTouchScreen))

	// Render only needs to be initialized once
	im = util.NewRender("#version 100")

	lastTouch = time.Now()
}

func redraw(act *ndk.Activity, win *ndk.Window) {
	act.Context().Call(func() {
		releaseEGL()
		initEGL(act, win)
	}, false)
	act.Context().Call(func() {
		draw(act, nil)
	}, false)
}

func destroyed(act *ndk.Activity, win *ndk.Window) {
	releaseEGL()
}

func draw(act *ndk.Activity, _ *ndk.Window) {
	if eglctx != nil {
		io := imgui.GetIO()

		pos := imgui.NewVec2(float32(lastX), float32(lastY))
		defer pos.Delete()
		io.SetMousePos(pos)
		io.SetMouseDown([]bool{mouseLeft, false, mouseRight, false, false})

		// Setup time step
		io.SetDeltaTime(float32(time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)))

		// Margin
		MARGIN := float32(width / 20)

		imgui.NewFrame()
		imgui.Begin("Camera", (*bool)(nil), int(imgui.WindowFlags_NoSavedSettings|
			imgui.WindowFlags_NoTitleBar|
			imgui.WindowFlags_NoResize|
			imgui.WindowFlags_NoMove))

		ds := imgui.GetWindowDrawList()
		ds.AddCallback(func(in interface{}) bool {
			cam.DrawImage()
			return true
		}, nil)

		imgui.End()

		if time.Now().Sub(lastTouch) < AUTOHIDETIME {
			curpos := imgui.NewVec2(MARGIN, MARGIN)
			defer curpos.Delete()
			imgui.SetNextWindowPos(curpos)
			isize := imgui.NewVec2(float32(width)-2*MARGIN, float32(height)-2*MARGIN)
			defer isize.Delete()
			imgui.SetNextWindowSize(isize)

			imgui.Begin(title, (*bool)(nil), int(imgui.WindowFlags_NoSavedSettings|
				imgui.WindowFlags_NoTitleBar|
				imgui.WindowFlags_NoResize|
				imgui.WindowFlags_NoCollapse|
				imgui.WindowFlags_NoMove))

			if imgui.Checkbox("flash", &flashOn) {
				if cam != nil {
					log.Println(" flash:", flashOn)
					if flashOn {
						cam.SetFlashMode(camera.FLASH_MODE_TORCH)
					} else {
						cam.SetFlashMode(camera.FLASH_MODE_OFF)
						//cam.SetFlashMode(camera.FLASH_MODE_ON)
					}
					cam.ApplyProperties()
				}
			}
			cam.Draw()

			imgui.End()
		}

		imgui.Render()

		// Rendering
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		im.Render(imgui.GetDrawData())

		eglctx.SwapBuffers()
	}
}
