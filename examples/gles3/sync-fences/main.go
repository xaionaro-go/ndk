// Sync fence example using OpenGL ES 3.0.
//
// Demonstrates the ES3-only FenceSync/ClientWaitSync mechanism for CPU-GPU
// synchronization. In ES2 the only way to guarantee GPU completion is the
// heavyweight glFinish. ES3 fence syncs let the CPU poll or block until a
// specific point in the GPU command stream has been reached.
//
// The example issues a clear + draw, inserts a fence, then waits for the
// GPU to finish processing those commands. This pattern is essential for
// safe buffer mapping (MapBufferRange) where the CPU must not read data
// that the GPU is still writing.
//
// This program must run on an Android device with OpenGL ES 3.0 support.
package main

import (
	"log"
	"unsafe"

	"github.com/xaionaro-go/ndk/egl"
	"github.com/xaionaro-go/ndk/gles2"
	"github.com/xaionaro-go/ndk/gles3"
)


var triVertices = [...]float32{
	0.0, 0.5,
	-0.5, -0.5,
	0.5, -0.5,
}

var vertSrc = "#version 300 es\nlayout(location=0) in vec2 aPos;\nvoid main() {\n    gl_Position = vec4(aPos, 0.0, 1.0);\n}\n\x00"

var fragSrc = "#version 300 es\nprecision mediump float;\nout vec4 fragColor;\nvoid main() {\n    fragColor = vec4(0.0, 0.8, 0.4, 1.0);\n}\n\x00"

func main() {
	display, surface, ctx := initEGL()
	defer teardownEGL(display, surface, ctx)

	program := buildProgram(vertSrc, fragSrc)
	defer gles2.DeleteProgram(program)

	// Set up geometry in a VAO.
	var vao gles3.GLuint
	var vbo gles2.GLuint
	gles3.GenVertexArrays(1, &vao)
	gles2.GenBuffers(1, &vbo)

	gles3.BindVertexArray(vao)
	gles2.BindBuffer(gles2.ArrayBuffer, vbo)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(triVertices)),
		unsafe.Pointer(&triVertices[0]),
		gles2.StaticDraw,
	)
	gles2.EnableVertexAttribArray(0)
	gles2.VertexAttribPointer(0, 2, gles2.Float, gles2.Boolean(gles2.False), 0, nil)

	// --- Issue GPU work ---
	gles2.Viewport(0, 0, 64, 64)
	gles2.ClearColor(0.0, 0.0, 0.0, 1.0)
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))
	gles2.UseProgram(program)
	gles2.DrawArrays(gles2.Triangles, 0, 3)

	checkGLError("after draw")

	// --- Insert a fence sync right after the draw ---
	// GL_SYNC_GPU_COMMANDS_COMPLETE means the fence is signalled when all
	// previously submitted GL commands have finished executing on the GPU.
	sync := gles3.FenceSync(gles3.SyncGpuCommandsComplete, 0)
	if sync == nil {
		log.Fatal("FenceSync returned nil")
	}
	log.Println("fence inserted after draw commands")

	// Flush to ensure the fence command reaches the GPU. Without this the
	// driver may buffer commands indefinitely and ClientWaitSync would hang.
	gles2.Flush()

	// --- Wait for the GPU to reach the fence ---
	// SYNC_FLUSH_COMMANDS_BIT tells the driver to flush if the sync is not
	// already signalled, which avoids a potential deadlock when the command
	// queue has not been flushed yet.
	result := gles3.ClientWaitSync(sync, gles3.SyncFlushCommandsBit, gles3.GLuint64(gles3.GL_TIMEOUT_IGNORED))
	switch result {
	case gles3.AlreadySignaled:
		log.Println("ClientWaitSync: already signalled (GPU was fast)")
	case gles3.ConditionSatisfied:
		log.Println("ClientWaitSync: condition satisfied (GPU caught up)")
	case gles3.TimeoutExpired:
		log.Println("ClientWaitSync: timeout expired (should not happen with TIMEOUT_IGNORED)")
	case gles3.WaitFailed:
		log.Fatal("ClientWaitSync: wait failed")
	default:
		log.Fatalf("ClientWaitSync: unexpected result 0x%x", result)
	}

	// At this point all GPU work before the fence has completed. It is now
	// safe to read back pixel data or map buffers that the GPU was writing.

	gles3.DeleteSync(sync)
	log.Println("fence deleted, GPU synchronization complete")

	// --- Cleanup ---
	gles3.BindVertexArray(0)
	gles3.DeleteVertexArrays(1, &vao)
	gles2.DeleteBuffers(1, &vbo)
}

// ---------------------------------------------------------------------------
// EGL and shader helpers
// ---------------------------------------------------------------------------

func initEGL() (egl.EGLDisplay, egl.EGLSurface, egl.EGLContext) {
	var defaultDisplay egl.EGLNativeDisplayType
	display := egl.GetDisplay(defaultDisplay)
	if egl.Initialize(display, nil, nil) == 0 {
		log.Fatalf("eglInitialize failed: 0x%x", egl.GetError())
	}

	configAttribs := [...]egl.Int{
		egl.RenderableType, egl.OpenglEs3Bit,
		egl.SurfaceType, egl.PbufferBit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.AlphaSize, 8,
		egl.None,
	}
	var config egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(display, &configAttribs[0], &config, 1, &numConfigs) == 0 || numConfigs == 0 {
		log.Fatalf("eglChooseConfig failed: 0x%x", egl.GetError())
	}

	surfAttribs := [...]egl.Int{egl.Width, 64, egl.Height, 64, egl.None}
	surface := egl.CreatePbufferSurface(display, config, &surfAttribs[0])

	ctxAttribs := [...]egl.Int{egl.ContextMajorVersion, 3, egl.None}
	ctx := egl.CreateContext(display, config, nil, &ctxAttribs[0])

	if egl.MakeCurrent(display, surface, surface, ctx) == 0 {
		log.Fatalf("eglMakeCurrent failed: 0x%x", egl.GetError())
	}
	return display, surface, ctx
}

func teardownEGL(display egl.EGLDisplay, surface egl.EGLSurface, ctx egl.EGLContext) {
	egl.MakeCurrent(display, nil, nil, nil)
	egl.DestroySurface(display, surface)
	egl.DestroyContext(display, ctx)
	egl.Terminate(display)
}

func buildProgram(vertSrc, fragSrc string) gles2.GLuint {
	vs := compileShader(gles2.VertexShader, vertSrc)
	fs := compileShader(gles2.FragmentShader, fragSrc)

	prog := gles2.CreateProgram()
	gles2.AttachShader(prog, vs)
	gles2.AttachShader(prog, fs)
	gles2.LinkProgram(prog)

	var status gles2.GLint
	gles2.GetProgramiv(prog, gles2.LinkStatus, &status)
	if status == 0 {
		log.Fatal("program link failed")
	}

	gles2.DeleteShader(vs)
	gles2.DeleteShader(fs)
	return prog
}

func compileShader(shaderType gles2.Enum, src string) gles2.GLuint {
	s := gles2.CreateShader(shaderType)
	srcPtr := strPtr(src)
	gles2.ShaderSource(s, 1, &srcPtr, nil)
	gles2.CompileShader(s)

	var status gles2.GLint
	gles2.GetShaderiv(s, gles2.CompileStatus, &status)
	if status == 0 {
		log.Fatalf("shader compile failed (type 0x%x)", shaderType)
	}
	return s
}

func strPtr(s string) *gles2.GLchar {
	b := []byte(s)
	return (*gles2.GLchar)(unsafe.Pointer(&b[0]))
}

func checkGLError(tag string) {
	if e := gles2.GetError(); e != gles2.NoError {
		log.Fatalf("GL error at %s: 0x%x", tag, e)
	}
}
