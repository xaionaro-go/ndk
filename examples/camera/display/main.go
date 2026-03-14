// Camera preview display using NativeActivity with OpenGL ES rendering.
//
// Opens the first available camera and renders its preview to the screen using
// GPU-based compositing for correct aspect ratio at full screen. The camera
// writes to a SurfaceTexture (zero-copy GPU path), which is sampled as an OES
// external texture and rendered as a fullscreen quad with UV coordinates that
// center-crop to the display aspect ratio.
//
// The app requests CAMERA permission at runtime via JNI. If permission is
// denied, a red screen and Toast message are shown.
//
// Build:
//
//	make apk-displaycamera
//
// Install & run:
//
//	adb install -r examples/camera/display/displaycamera.apk
//	adb shell am start -n com.example.displaycamera/android.app.NativeActivity
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/xaionaro-go/ndk/activity"
	"github.com/xaionaro-go/ndk/camera"
	"github.com/xaionaro-go/ndk/egl"
	"github.com/xaionaro-go/ndk/gles2"
	"github.com/xaionaro-go/ndk/input"
	"github.com/xaionaro-go/ndk/jni"
	ndklog "github.com/xaionaro-go/ndk/log"
	"github.com/xaionaro-go/ndk/surfacetexture"
	"github.com/xaionaro-go/ndk/window"
)

const logTag = "ndk-camera"

const (
	vertexShaderSrc = `
attribute vec4 aPosition;
attribute vec2 aTexCoord;
varying vec2 vTexCoord;
uniform mat4 uSTMatrix;
void main() {
    gl_Position = aPosition;
    vTexCoord = (uSTMatrix * vec4(aTexCoord, 0.0, 1.0)).xy;
}
`
	fragmentShaderSrc = `
#extension GL_OES_EGL_image_external : require
precision mediump float;
varying vec2 vTexCoord;
uniform samplerExternalOES sTexture;
void main() {
    gl_FragColor = texture2D(sTexture, vTexCoord);
}
`
)

var (
	currentActivity *activity.Activity
	currentWindow   unsafe.Pointer
	cameraStarted   bool

	permRequested  bool
	permLostFocus  bool
	permDialogDone bool

	camMgr           *camera.Manager
	camDev           *camera.Device
	captureReq       *camera.CaptureRequest
	outTarget        *camera.OutputTarget
	sessOut          *camera.SessionOutput
	sessOutContainer *camera.SessionOutputContainer
	captureSess      *camera.CaptureSession

	eglDisp    egl.EGLDisplay
	eglSurf    egl.EGLSurface
	eglCtx     egl.EGLContext
	eglInited  bool
	eglHasSurf bool
	eglHasCtx  bool
	glProg     gles2.GLuint
	oesTex     gles2.GLuint

	surfTex   *surfacetexture.SurfaceTexture
	camWindow *surfacetexture.NativeWindow

	renderMu   sync.Mutex
	renderStop chan struct{}

	inputQueue *input.Queue
	inputStop  chan struct{}
)

func init() {
	activity.SetLifecycleCallbacks(activity.LifecycleCallbacks{
		OnNativeWindowCreated: func(act *activity.Activity, win unsafe.Pointer) {
			currentActivity = act
			currentWindow = win
			tryStartCamera()
		},
		OnResume: func(act *activity.Activity) {
			currentActivity = act
			if currentWindow != nil && !cameraStarted && (!permRequested || permDialogDone) {
				tryStartCamera()
			}
		},
		OnWindowFocusChanged: func(_ *activity.Activity, hasFocus int32) {
			if hasFocus == 0 && permRequested && !permDialogDone {
				permLostFocus = true
			}
			if hasFocus != 0 && permRequested && permLostFocus && !permDialogDone {
				permDialogDone = true
				if currentWindow != nil && !cameraStarted {
					tryStartCamera()
				}
			}
		},
		OnNativeWindowDestroyed: func(_ *activity.Activity, _ unsafe.Pointer) {
			stopPreview()
			currentWindow = nil
		},
		OnInputQueueCreated: func(_ *activity.Activity, queuePtr unsafe.Pointer) {
			inputQueue = input.NewQueueFromPointer(queuePtr)
			inputStop = make(chan struct{})
			go drainInputEvents()
		},
		OnInputQueueDestroyed: func(_ *activity.Activity, _ unsafe.Pointer) {
			if inputStop != nil {
				close(inputStop)
				inputStop = nil
			}
			inputQueue = nil
		},
		OnDestroy: func(_ *activity.Activity) {
			stopPreview()
			currentActivity = nil
			currentWindow = nil
		},
	})
}

func tryStartCamera() {
	if cameraStarted {
		return
	}

	actPtr := currentActivity.Pointer()

	if !jni.HasPermission(actPtr, "android.permission.CAMERA") {
		switch {
		case !permRequested:
			logInfo("requesting camera permission")
			jni.RequestPermission(actPtr, "android.permission.CAMERA")
			permRequested = true
		case permDialogDone:
			showError("camera permission denied")
		}
		return
	}

	if err := startPreview(currentWindow); err != nil {
		showError(fmt.Sprintf("camera preview failed: %v", err))
	}
}

func startPreview(win unsafe.Pointer) (_err error) {
	defer func() {
		if _err != nil {
			stopPreview()
		}
	}()

	nw := window.NewWindowFromPointer(win)
	winW := nw.Width()
	winH := nw.Height()
	logInfo(fmt.Sprintf("window: %dx%d", winW, winH))

	camMgr = camera.NewManager()

	ids, err := camMgr.CameraIDList()
	if err != nil {
		return fmt.Errorf("listing cameras: %w", err)
	}
	if len(ids) == 0 {
		return fmt.Errorf("no cameras available")
	}
	logInfo(fmt.Sprintf("opening camera %s", ids[0]))

	sensorOrient, bestW, bestH := getCameraInfo(camMgr, ids[0], winW, winH)
	if bestW == 0 || bestH == 0 {
		bestW = 1920
		bestH = 1080
	}
	logInfo(fmt.Sprintf("camera output: %dx%d orient=%d", bestW, bestH, sensorOrient))

	if err := initEGL(win); err != nil {
		return fmt.Errorf("EGL init: %w", err)
	}

	// Create OES texture for camera frames.
	gles2.GenTextures(1, &oesTex)
	gles2.BindTexture(gles2.TextureExternalOes, oesTex)
	gles2.TexParameteri(gles2.TextureExternalOes, gles2.TextureMinFilter, gles2.GLint(gles2.Linear))
	gles2.TexParameteri(gles2.TextureExternalOes, gles2.TextureMagFilter, gles2.GLint(gles2.Linear))
	gles2.TexParameteri(gles2.TextureExternalOes, gles2.TextureWrapS, gles2.GLint(gles2.ClampToEdge))
	gles2.TexParameteri(gles2.TextureExternalOes, gles2.TextureWrapT, gles2.GLint(gles2.ClampToEdge))

	// Create SurfaceTexture via JNI, wrapping our OES texture.
	stPtr := jni.CreateSurfaceTexture(
		currentActivity.Pointer(),
		int(oesTex), int(bestW), int(bestH),
	)
	if stPtr == nil {
		return fmt.Errorf("failed to create SurfaceTexture")
	}
	surfTex = surfacetexture.NewSurfaceTextureFromPointer(stPtr)

	camWindow = surfTex.AcquireWindow()
	if camWindow.Pointer() == nil {
		return fmt.Errorf("failed to acquire SurfaceTexture window")
	}

	prog, err := createShaderProgram()
	if err != nil {
		return fmt.Errorf("shader program: %w", err)
	}
	glProg = prog

	// Release EGL context from this thread so the render goroutine can use it.
	var noSurf egl.EGLSurface
	var noCtx egl.EGLContext
	egl.MakeCurrent(eglDisp, noSurf, noSurf, noCtx)

	// Open camera and create capture session targeting the SurfaceTexture.
	camDev, err = camMgr.OpenCamera(ids[0], camera.DeviceStateCallbacks{
		OnDisconnected: func() { logError("camera disconnected") },
		OnError:        func(code int) { logError(fmt.Sprintf("camera error: %d", code)) },
	})
	if err != nil {
		return fmt.Errorf("opening camera %s: %w", ids[0], err)
	}

	captureReq, err = camDev.CreateCaptureRequest(camera.Preview)
	if err != nil {
		return fmt.Errorf("creating capture request: %w", err)
	}

	camWinPtr := (*camera.ANativeWindow)(camWindow.Pointer())

	outTarget, err = camera.NewOutputTarget(camWinPtr)
	if err != nil {
		return fmt.Errorf("creating output target: %w", err)
	}
	captureReq.AddTarget(outTarget)

	sessOut, err = camera.NewSessionOutput(camWinPtr)
	if err != nil {
		return fmt.Errorf("creating session output: %w", err)
	}

	sessOutContainer, err = camera.NewSessionOutputContainer()
	if err != nil {
		return fmt.Errorf("creating output container: %w", err)
	}
	if err := sessOutContainer.Add(sessOut); err != nil {
		return fmt.Errorf("adding session output: %w", err)
	}

	captureSess, err = camDev.CreateCaptureSession(
		sessOutContainer,
		camera.SessionStateCallbacks{
			OnReady:  func() { logInfo("capture session ready") },
			OnActive: func() { logInfo("capture session active") },
		},
	)
	if err != nil {
		return fmt.Errorf("creating capture session: %w", err)
	}

	if err := captureSess.SetRepeatingRequest(captureReq); err != nil {
		return fmt.Errorf("setting repeating request: %w", err)
	}

	// Compute vertex scaling for center-crop at correct aspect ratio.
	// The SurfaceTexture transform matrix handles orientation; we scale
	// the quad beyond [-1,1] so the viewport clips the excess.
	var effectiveCamAR float32
	if sensorOrient == 90 || sensorOrient == 270 {
		effectiveCamAR = float32(bestH) / float32(bestW)
	} else {
		effectiveCamAR = float32(bestW) / float32(bestH)
	}
	displayAR := float32(winW) / float32(winH)

	var scaleX, scaleY float32
	ratio := effectiveCamAR / displayAR
	if ratio > 1 {
		scaleX = ratio
		scaleY = 1.0
	} else {
		scaleX = 1.0
		scaleY = 1.0 / ratio
	}

	logInfo(fmt.Sprintf("scale: x=%.3f y=%.3f camAR=%.3f displayAR=%.3f",
		scaleX, scaleY, effectiveCamAR, displayAR))

	renderStop = make(chan struct{})
	go renderLoop(winW, winH, scaleX, scaleY)

	logInfo("camera preview started")
	cameraStarted = true
	return nil
}

// getCameraInfo reads sensor orientation and picks the best output size.
func getCameraInfo(
	mgr *camera.Manager,
	camID string,
	winW, winH int32,
) (orient, bestW, bestH int32) {
	chars, err := mgr.GetCameraCharacteristics(camID)
	if err != nil {
		return 0, 0, 0
	}
	defer chars.Close()

	// Sensor orientation.
	if chars.I32Count(uint32(camera.SensorOrientation)) > 0 {
		orient = chars.I32At(uint32(camera.SensorOrientation), 0)
	}

	// Target aspect ratio (landscape).
	var targetAR float64
	if orient == 90 || orient == 270 {
		targetAR = float64(winH) / float64(winW)
	} else {
		targetAR = float64(winW) / float64(winH)
	}
	if targetAR < 1.0 {
		targetAR = 1.0 / targetAR
	}

	// Stream configurations: tuples of (format, width, height, isInput).
	scTag := uint32(camera.ScalerAvailableStreamConfigurations)
	count := chars.I32Count(scTag)
	bestDiff := 1e9
	var bestPixels int32
	for i := int32(0); i+3 < count; i += 4 {
		format := chars.I32At(scTag, i)
		w := chars.I32At(scTag, i+1)
		h := chars.I32At(scTag, i+2)
		isInput := chars.I32At(scTag, i+3)
		if isInput != 0 {
			continue
		}
		// 0x22 = YUV_420_888, 0x23 = RAW_SENSOR
		if format != 0x22 && format != 0x23 {
			continue
		}
		ar := float64(w) / float64(h)
		if w < h {
			ar = float64(h) / float64(w)
		}
		maxDim := w
		if h > w {
			maxDim = h
		}
		if maxDim < 720 {
			continue
		}
		diff := (ar - targetAR) * (ar - targetAR)
		pixels := w * h
		if diff < bestDiff-0.001 || (diff < bestDiff+0.001 && pixels > bestPixels) {
			bestDiff = diff
			bestPixels = pixels
			bestW = w
			bestH = h
		}
	}

	logInfo(fmt.Sprintf("best output: %dx%d (targetAR=%.3f)", bestW, bestH, targetAR))
	return orient, bestW, bestH
}

func renderLoop(
	winW, winH int32,
	scaleX, scaleY float32,
) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	renderMu.Lock()
	defer renderMu.Unlock()

	egl.MakeCurrent(eglDisp, eglSurf, eglSurf, eglCtx)
	gles2.Viewport(gles2.GLint(0), gles2.GLint(0), gles2.GLsizei(winW), gles2.GLsizei(winH))

	aPos := gles2.GetAttribLocation(glProg, glStr("aPosition"))
	aTex := gles2.GetAttribLocation(glProg, glStr("aTexCoord"))
	uST := gles2.GetUniformLocation(glProg, glStr("uSTMatrix"))
	sSamp := gles2.GetUniformLocation(glProg, glStr("sTexture"))

	// Vertex positions scaled beyond [-1,1] for center-crop; viewport clips excess.
	verts := [8]float32{
		-scaleX, -scaleY,
		scaleX, -scaleY,
		-scaleX, scaleY,
		scaleX, scaleY,
	}

	// Standard OpenGL UVs: (0,0) at bottom-left, (1,1) at top-right.
	// The SurfaceTexture transform matrix handles orientation.
	uvs := [8]float32{
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	}

	var stMat [16]float32

	for {
		select {
		case <-renderStop:
			return
		default:
		}

		if err := surfTex.UpdateTexImage(); err != nil {
			continue
		}
		surfTex.TransformMatrix(&stMat)

		gles2.ClearColor(0, 0, 0, 1)
		gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

		gles2.UseProgram(glProg)
		gles2.ActiveTexture(gles2.Texture0)
		gles2.BindTexture(gles2.TextureExternalOes, oesTex)
		gles2.Uniform1i(sSamp, 0)
		gles2.UniformMatrix4fv(uST, 1, gles2.Boolean(gles2.False), (*gles2.GLfloat)(&stMat[0]))

		gles2.EnableVertexAttribArray(gles2.GLuint(aPos))
		gles2.VertexAttribPointer(gles2.GLuint(aPos), 2, gles2.Float, gles2.Boolean(gles2.False), 0, unsafe.Pointer(&verts[0]))

		gles2.EnableVertexAttribArray(gles2.GLuint(aTex))
		gles2.VertexAttribPointer(gles2.GLuint(aTex), 2, gles2.Float, gles2.Boolean(gles2.False), 0, unsafe.Pointer(&uvs[0]))

		gles2.DrawArrays(gles2.TriangleStrip, 0, 4)

		gles2.DisableVertexAttribArray(gles2.GLuint(aPos))
		gles2.DisableVertexAttribArray(gles2.GLuint(aTex))

		egl.SwapBuffers(eglDisp, eglSurf)
	}
}

func initEGL(win unsafe.Pointer) error {
	var noDisplay egl.EGLNativeDisplayType
	eglDisp = egl.GetDisplay(noDisplay)
	if egl.Initialize(eglDisp, nil, nil) == 0 {
		return fmt.Errorf("eglInitialize failed: 0x%x", egl.GetError())
	}
	eglInited = true

	attribs := []egl.Int{
		egl.RenderableType, egl.OpenglEs2Bit,
		egl.SurfaceType, egl.WindowBit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.AlphaSize, 8,
		egl.None,
	}

	var config egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(eglDisp, &attribs[0], &config, 1, &numConfigs) == 0 || numConfigs == 0 {
		return fmt.Errorf("eglChooseConfig failed: 0x%x", egl.GetError())
	}

	ctxAttribs := []egl.Int{
		egl.ContextClientVersion, 2,
		egl.None,
	}
	eglCtx = egl.CreateContext(eglDisp, config, nil, &ctxAttribs[0])
	if errCode := egl.GetError(); errCode != 0x3000 { // EGL_SUCCESS
		return fmt.Errorf("eglCreateContext failed: 0x%x", errCode)
	}
	eglHasCtx = true

	eglSurf = egl.CreateWindowSurface(eglDisp, config, egl.EGLNativeWindowType(win), nil)
	if errCode := egl.GetError(); errCode != 0x3000 { // EGL_SUCCESS
		return fmt.Errorf("eglCreateWindowSurface failed: 0x%x", errCode)
	}
	eglHasSurf = true

	if egl.MakeCurrent(eglDisp, eglSurf, eglSurf, eglCtx) == 0 {
		return fmt.Errorf("eglMakeCurrent failed: 0x%x", egl.GetError())
	}

	egl.SwapInterval(eglDisp, 1)
	return nil
}

func cleanupEGL() {
	if eglInited {
		var noSurface egl.EGLSurface
		var noContext egl.EGLContext
		egl.MakeCurrent(eglDisp, noSurface, noSurface, noContext)
		if eglHasCtx {
			egl.DestroyContext(eglDisp, eglCtx)
			eglHasCtx = false
		}
		if eglHasSurf {
			egl.DestroySurface(eglDisp, eglSurf)
			eglHasSurf = false
		}
		egl.Terminate(eglDisp)
		eglInited = false
	}
}

func createShaderProgram() (gles2.GLuint, error) {
	vs, err := compileShader(gles2.VertexShader, vertexShaderSrc)
	if err != nil {
		return 0, fmt.Errorf("vertex shader: %w", err)
	}

	fs, err := compileShader(gles2.FragmentShader, fragmentShaderSrc)
	if err != nil {
		gles2.DeleteShader(vs)
		return 0, fmt.Errorf("fragment shader: %w", err)
	}

	prog := gles2.CreateProgram()
	gles2.AttachShader(prog, vs)
	gles2.AttachShader(prog, fs)
	gles2.LinkProgram(prog)

	var status gles2.GLint
	gles2.GetProgramiv(prog, gles2.LinkStatus, &status)
	if status == 0 {
		gles2.DeleteProgram(prog)
		gles2.DeleteShader(vs)
		gles2.DeleteShader(fs)
		return 0, fmt.Errorf("program link failed")
	}

	gles2.DeleteShader(vs)
	gles2.DeleteShader(fs)
	return prog, nil
}

func compileShader(shaderType gles2.Enum, source string) (gles2.GLuint, error) {
	shader := gles2.CreateShader(shaderType)

	src := glStr(source)
	// Pin the Go-allocated string data so cgocheck allows passing
	// &src (Go pointer to Go pointer) to the C function.
	var pinner runtime.Pinner
	pinner.Pin(src)
	defer pinner.Unpin()

	length := gles2.GLint(len(source))
	gles2.ShaderSource(shader, 1, &src, &length)
	gles2.CompileShader(shader)

	var status gles2.GLint
	gles2.GetShaderiv(shader, gles2.CompileStatus, &status)
	if status == 0 {
		gles2.DeleteShader(shader)
		return 0, fmt.Errorf("compile failed for shader type %d", shaderType)
	}
	return shader, nil
}

func stopPreview() {
	if renderStop != nil {
		close(renderStop)
		renderMu.Lock()
		renderMu.Unlock()
		renderStop = nil
	}

	cameraStarted = false

	if captureSess != nil {
		captureSess.StopRepeating()
		captureSess.Close()
		captureSess = nil
	}
	if sessOutContainer != nil {
		sessOutContainer.Close()
		sessOutContainer = nil
	}
	if sessOut != nil {
		sessOut.Close()
		sessOut = nil
	}
	if outTarget != nil {
		outTarget.Close()
		outTarget = nil
	}
	if captureReq != nil {
		captureReq.Close()
		captureReq = nil
	}
	if camDev != nil {
		camDev.Close()
		camDev = nil
	}
	if camWindow != nil {
		camWindow = nil
	}
	if surfTex != nil {
		surfTex.Close()
		surfTex = nil
	}
	if glProg != 0 {
		gles2.DeleteProgram(glProg)
		glProg = 0
	}
	if oesTex != 0 {
		gles2.DeleteTextures(1, &oesTex)
		oesTex = 0
	}

	cleanupEGL()

	if camMgr != nil {
		camMgr.Close()
		camMgr = nil
	}
}

func showError(msg string) {
	logError(msg)

	if currentWindow != nil {
		jni.FillWindowColor(currentWindow, 0xFF0000CC)
	}

	if currentActivity != nil {
		jni.ShowToast(currentActivity.Pointer(), msg)
	}
}

func logInfo(msg string) {
	ndklog.Write(int32(ndklog.Info), logTag, msg)
}

func logError(msg string) {
	ndklog.Write(int32(ndklog.Error), logTag, msg)
}

// drainInputEvents polls the input queue and discards all events so
// Android's InputDispatcher does not trigger an ANR dialog.
func drainInputEvents() {
	q := inputQueue
	for {
		select {
		case <-inputStop:
			return
		default:
		}

		ev := q.GetEvent()
		if ev == nil {
			time.Sleep(16 * time.Millisecond)
			continue
		}

		if !q.PreDispatchEvent(ev) {
			q.FinishEvent(ev, 0)
		}
	}
}

// glStr converts a Go string to a null-terminated *gles2.GLchar.
func glStr(s string) *gles2.GLchar {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return (*gles2.GLchar)(unsafe.Pointer(&b[0]))
}

func main() {}
