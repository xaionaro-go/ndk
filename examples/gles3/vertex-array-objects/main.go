// Vertex Array Object (VAO) example using OpenGL ES 3.0.
//
// Demonstrates the ES3-only VAO feature that encapsulates vertex attribute
// state into named objects. In ES2 every draw call must re-bind vertex
// pointers; with VAOs the driver records the bindings once and replays them
// on BindVertexArray.
//
// This example creates two VAOs -- one for a triangle and one for a quad --
// each backed by its own VBO. Switching between them before drawing shows
// how VAOs eliminate redundant EnableVertexAttribArray/VertexAttribPointer
// calls at render time.
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

// Triangle: three vertices with XY position.
var triangleVertices = [...]float32{
	0.0, 0.5,
	-0.5, -0.5,
	0.5, -0.5,
}

// Quad: four vertices forming a small square (drawn as triangle fan).
var quadVertices = [...]float32{
	-0.3, 0.3,
	-0.3, -0.3,
	0.3, -0.3,
	0.3, 0.3,
}

// vertSrc is a minimal GLSL ES 3.00 vertex shader.
var vertSrc = "#version 300 es\nlayout(location=0) in vec2 aPos;\nvoid main() {\n    gl_Position = vec4(aPos, 0.0, 1.0);\n}\n\x00"

// fragSrc outputs a solid color.
var fragSrc = "#version 300 es\nprecision mediump float;\nout vec4 fragColor;\nuniform vec4 uColor;\nvoid main() {\n    fragColor = uColor;\n}\n\x00"

func main() {
	display, surface, ctx := initEGL()
	defer teardownEGL(display, surface, ctx)

	program := buildProgram(vertSrc, fragSrc)
	defer gles2.DeleteProgram(program)

	colorLoc := gles2.GetUniformLocation(program, strPtr("uColor\x00"))

	// --- Create VAO and VBO for the triangle ---
	var triVAO gles3.GLuint
	var triVBO gles2.GLuint
	gles3.GenVertexArrays(1, &triVAO)
	gles2.GenBuffers(1, &triVBO)

	gles3.BindVertexArray(triVAO)
	gles2.BindBuffer(gles2.ArrayBuffer, triVBO)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(triangleVertices)),
		unsafe.Pointer(&triangleVertices[0]),
		gles2.StaticDraw,
	)
	// Attribute 0: vec2 position, tightly packed.
	gles2.EnableVertexAttribArray(0)
	gles2.VertexAttribPointer(0, 2, gles2.Float, gles2.Boolean(gles2.False), 0, nil)
	gles3.BindVertexArray(0) // unbind so subsequent buffer calls don't affect this VAO

	// --- Create VAO and VBO for the quad ---
	var quadVAO gles3.GLuint
	var quadVBO gles2.GLuint
	gles3.GenVertexArrays(1, &quadVAO)
	gles2.GenBuffers(1, &quadVBO)

	gles3.BindVertexArray(quadVAO)
	gles2.BindBuffer(gles2.ArrayBuffer, quadVBO)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(quadVertices)),
		unsafe.Pointer(&quadVertices[0]),
		gles2.StaticDraw,
	)
	gles2.EnableVertexAttribArray(0)
	gles2.VertexAttribPointer(0, 2, gles2.Float, gles2.Boolean(gles2.False), 0, nil)
	gles3.BindVertexArray(0)

	// --- Render one frame ---
	gles2.Viewport(0, 0, 64, 64)
	gles2.ClearColor(0.1, 0.1, 0.1, 1.0)
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

	gles2.UseProgram(program)

	// Draw the triangle in red. Binding the VAO restores all vertex state.
	red := [4]gles2.GLfloat{1.0, 0.0, 0.0, 1.0}
	gles2.Uniform4fv(colorLoc, 1, &red[0])
	gles3.BindVertexArray(triVAO)
	gles2.DrawArrays(gles2.Triangles, 0, 3)

	// Draw the quad in green. No need to call VertexAttribPointer again.
	green := [4]gles2.GLfloat{0.0, 1.0, 0.0, 1.0}
	gles2.Uniform4fv(colorLoc, 1, &green[0])
	gles3.BindVertexArray(quadVAO)
	gles2.DrawArrays(gles2.TriangleFan, 0, 4)

	gles3.BindVertexArray(0)
	gles2.Flush()
	checkGLError("after draw")

	log.Println("VAO example: drew triangle (red) and quad (green)")

	// --- Cleanup ---
	gles3.DeleteVertexArrays(1, &triVAO)
	gles3.DeleteVertexArrays(1, &quadVAO)
	gles2.DeleteBuffers(1, &triVBO)
	gles2.DeleteBuffers(1, &quadVBO)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// initEGL creates an ES 3.0 context with a 64x64 PBuffer surface.
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

// buildProgram compiles a vertex/fragment pair and links them.
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

// strPtr returns a pointer to the first byte of a NUL-terminated Go string.
func strPtr(s string) *gles2.GLchar {
	b := []byte(s)
	return (*gles2.GLchar)(unsafe.Pointer(&b[0]))
}

func checkGLError(tag string) {
	if e := gles2.GetError(); e != gles2.NoError {
		log.Fatalf("GL error at %s: 0x%x", tag, e)
	}
}
