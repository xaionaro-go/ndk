// Instanced drawing example using OpenGL ES 3.0.
//
// Demonstrates DrawArraysInstanced -- the ES3 call that renders multiple
// copies of geometry in a single draw call. Each instance receives a unique
// gl_InstanceID in the vertex shader, which is used here to position copies
// of a small triangle in a grid pattern.
//
// The vertex shader uses gl_InstanceID (built-in in GLSL ES 3.00) to compute
// a per-instance offset. No VertexAttribDivisor is needed for this approach
// because the offset is derived from the built-in instance index rather than
// from a per-instance attribute buffer.
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

// A small triangle centered at the origin, scaled down so many fit on screen.
var triVertices = [...]float32{
	0.0, 0.08,
	-0.06, -0.04,
	0.06, -0.04,
}

// vertSrc positions each instance in a 4x4 grid using gl_InstanceID.
var vertSrc = `#version 300 es
layout(location=0) in vec2 aPos;

void main() {
    // Lay out instances in a 4-column grid.
    int col = gl_InstanceID % 4;
    int row = gl_InstanceID / 4;

    // Map grid cell to clip-space: columns span [-0.8, 0.8], rows likewise.
    float cx = -0.6 + float(col) * 0.4;
    float cy =  0.6 - float(row) * 0.4;

    gl_Position = vec4(aPos.x + cx, aPos.y + cy, 0.0, 1.0);
}
` + "\x00"

// fragSrc colours each instance based on gl_InstanceID via flat interpolation.
var fragSrc = `#version 300 es
precision mediump float;
out vec4 fragColor;

void main() {
    fragColor = vec4(0.2, 0.7, 1.0, 1.0);
}
` + "\x00"

const instanceCount = 16 // 4x4 grid

func main() {
	display, surface, ctx := initEGL()
	defer teardownEGL(display, surface, ctx)

	program := buildProgram(vertSrc, fragSrc)
	defer gles2.DeleteProgram(program)

	// Create a VAO + VBO for the triangle template geometry.
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

	// --- Render ---
	gles2.Viewport(0, 0, 128, 128)
	gles2.ClearColor(0.05, 0.05, 0.15, 1.0)
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

	gles2.UseProgram(program)

	// DrawArraysInstanced draws 3 vertices (one triangle) x instanceCount
	// times. The vertex shader offsets each instance to a grid cell.
	gles3.DrawArraysInstanced(
		gles3.GLenum(gles2.Triangles), // mode
		0,                             // first
		3,                             // count (vertices per instance)
		instanceCount,                 // instancecount
	)

	gles3.BindVertexArray(0)
	gles2.Flush()
	checkGLError("after instanced draw")

	log.Printf("instanced drawing: rendered %d triangle instances", instanceCount)

	// --- Cleanup ---
	gles3.DeleteVertexArrays(1, &vao)
	gles2.DeleteBuffers(1, &vbo)
}

// ---------------------------------------------------------------------------
// EGL and shader helpers (same pattern as the VAO example)
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

	surfAttribs := [...]egl.Int{egl.Width, 128, egl.Height, 128, egl.None}
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
