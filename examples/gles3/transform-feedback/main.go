// Transform feedback example using OpenGL ES 3.0.
//
// Demonstrates the ES3-only transform feedback feature, which captures
// vertex shader outputs into a buffer object on the GPU. This enables
// GPU-side particle systems, physics simulations, and stream-out
// techniques without reading data back to the CPU.
//
// The vertex shader takes input positions and applies a simple
// transformation (scaling + offset). The transformed positions are
// captured into a feedback buffer via TransformFeedbackVaryings. After
// the feedback pass the captured data is read back with MapBufferRange
// to verify correctness.
//
// This program must run on an Android device with OpenGL ES 3.0 support.
package main

import (
	"log"
	"unsafe"

	"github.com/AndroidGoLab/ndk/egl"
	"github.com/AndroidGoLab/ndk/gles2"
	"github.com/AndroidGoLab/ndk/gles3"
)

// Input positions: three 2D points.
var inputPositions = [...]float32{
	1.0, 2.0,
	3.0, 4.0,
	5.0, 6.0,
}

// The vertex shader scales each position by 0.5 and adds (0.1, 0.2).
// The "out" varying outPos is what transform feedback captures.
var vertSrc = `#version 300 es
layout(location=0) in vec2 aPos;
out vec2 outPos;

void main() {
    outPos = aPos * 0.5 + vec2(0.1, 0.2);
    gl_Position = vec4(outPos, 0.0, 1.0);
}
` + "\x00"

// Fragment shader is required for a valid program, but rasterization is
// disabled during transform feedback so it never executes.
var fragSrc = `#version 300 es
precision mediump float;
out vec4 fragColor;
void main() {
    fragColor = vec4(1.0);
}
` + "\x00"

const vertexCount = 3

func main() {
	display, surface, ctx := initEGL()
	defer teardownEGL(display, surface, ctx)

	program := buildTransformFeedbackProgram(vertSrc, fragSrc, "outPos")
	defer gles2.DeleteProgram(program)

	// --- Create VAO + VBO for input positions ---
	var vao gles3.GLuint
	var inputVBO gles2.GLuint
	gles3.GenVertexArrays(1, &vao)
	gles2.GenBuffers(1, &inputVBO)

	gles3.BindVertexArray(vao)
	gles2.BindBuffer(gles2.ArrayBuffer, inputVBO)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(inputPositions)),
		unsafe.Pointer(&inputPositions[0]),
		gles2.StaticDraw,
	)
	gles2.EnableVertexAttribArray(0)
	gles2.VertexAttribPointer(0, 2, gles2.Float, gles2.Boolean(gles2.False), 0, nil)

	// --- Create the transform feedback buffer ---
	// This buffer will receive the vertex shader output (vec2 per vertex).
	feedbackBufSize := gles2.GLsizeiptr(vertexCount * 2 * int(unsafe.Sizeof(float32(0))))
	var feedbackVBO gles2.GLuint
	gles2.GenBuffers(1, &feedbackVBO)
	gles2.BindBuffer(gles2.Enum(gles3.TransformFeedbackBuffer), feedbackVBO)
	gles2.BufferData(
		gles2.Enum(gles3.TransformFeedbackBuffer),
		feedbackBufSize,
		nil, // allocate only, no initial data
		gles2.StaticDraw,
	)

	// --- Create a transform feedback object ---
	var tfo gles3.GLuint
	gles3.GenTransformFeedbacks(1, &tfo)
	gles3.BindTransformFeedback(gles3.TransformFeedback, tfo)

	// Bind the feedback buffer to binding point 0 of the TF object.
	gles3.BindBufferBase(gles3.TransformFeedbackBuffer, 0, gles3.GLuint(feedbackVBO))

	// --- Execute the transform feedback pass ---
	gles2.UseProgram(program)

	// Disable rasterization: we only care about the vertex shader output.
	gles2.Enable(gles3.GL_RASTERIZER_DISCARD)

	gles3.BeginTransformFeedback(gles3.GLenum(gles2.Points))
	gles2.DrawArrays(gles2.Points, 0, gles2.GLsizei(vertexCount))
	gles3.EndTransformFeedback()

	gles2.Disable(gles3.GL_RASTERIZER_DISCARD)
	checkGLError("after transform feedback")

	// --- Read back captured data ---
	// Use a fence to ensure the GPU has finished writing before mapping.
	sync := gles3.FenceSync(gles3.SyncGpuCommandsComplete, 0)
	gles2.Flush()
	gles3.ClientWaitSync(sync, gles3.SyncFlushCommandsBit, gles3.GLuint64(0xFFFFFFFFFFFFFFFF))
	gles3.DeleteSync(sync)

	gles2.BindBuffer(gles2.Enum(gles3.TransformFeedbackBuffer), feedbackVBO)
	ptr := gles3.MapBufferRange(
		gles3.TransformFeedbackBuffer,
		0,
		gles3.GLsizeiptr(feedbackBufSize),
		gles3.GL_MAP_READ_BIT,
	)
	if ptr == nil {
		log.Fatal("MapBufferRange returned nil")
	}

	// Interpret the mapped memory as float32 pairs.
	results := unsafe.Slice((*float32)(ptr), vertexCount*2)
	for i := 0; i < vertexCount; i++ {
		x := results[i*2]
		y := results[i*2+1]
		// Expected: input * 0.5 + (0.1, 0.2)
		log.Printf("  vertex %d: (%.2f, %.2f) -> (%.2f, %.2f)",
			i, inputPositions[i*2], inputPositions[i*2+1], x, y)
	}

	gles3.UnmapBuffer(gles3.TransformFeedbackBuffer)

	log.Println("transform feedback: captured vertex shader output successfully")

	// --- Cleanup ---
	gles3.BindTransformFeedback(gles3.TransformFeedback, 0)
	gles3.DeleteTransformFeedbacks(1, &tfo)
	gles3.BindVertexArray(0)
	gles3.DeleteVertexArrays(1, &vao)
	gles2.DeleteBuffers(1, &inputVBO)
	gles2.DeleteBuffers(1, &feedbackVBO)
}

// buildTransformFeedbackProgram compiles the shaders, configures the
// transform feedback varying, and links the program. The varying must be
// declared before linking.
func buildTransformFeedbackProgram(vertSrc, fragSrc, varyingName string) gles2.GLuint {
	vs := compileShader(gles2.VertexShader, vertSrc)
	fs := compileShader(gles2.FragmentShader, fragSrc)

	prog := gles2.CreateProgram()
	gles2.AttachShader(prog, vs)
	gles2.AttachShader(prog, fs)

	// Specify which vertex shader outputs to capture.
	// The name must be NUL-terminated and passed as a C string array.
	varyingCStr := []byte(varyingName + "\x00")
	varyingPtr := (*gles3.GLchar)(unsafe.Pointer(&varyingCStr[0]))
	gles3.TransformFeedbackVaryings(
		gles3.GLuint(prog),
		1,
		&varyingPtr,
		gles3.InterleavedAttribs,
	)

	gles2.LinkProgram(prog)

	var status gles2.GLint
	gles2.GetProgramiv(prog, gles2.LinkStatus, &status)
	if status == 0 {
		log.Fatal("program link failed (transform feedback)")
	}

	gles2.DeleteShader(vs)
	gles2.DeleteShader(fs)
	return prog
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
