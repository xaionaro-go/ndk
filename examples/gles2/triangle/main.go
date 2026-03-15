// Hello Triangle: the classic first OpenGL ES 2.0 program.
//
// Demonstrates:
//   - Compiling a vertex shader and a fragment shader
//   - Linking them into a program
//   - Uploading triangle vertex data into a VBO
//   - Setting up vertex attribute pointers
//   - Drawing with DrawArrays
//   - Reading back pixels to verify the triangle was rasterized
//
// The triangle covers the center of a 64x64 pbuffer. After rendering, the
// program reads back the center pixel (which should be inside the triangle)
// and a corner pixel (which should be the clear color) to confirm correctness.
package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/AndroidGoLab/ndk/egl"
	"github.com/AndroidGoLab/ndk/gles2"
)

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "fatal: %s (EGL error 0x%04X)\n", msg, egl.GetError())
	os.Exit(1)
}

func fatalGL(msg string) {
	fmt.Fprintf(os.Stderr, "fatal: %s (GL error 0x%04X)\n", msg, int32(gles2.GetError()))
	os.Exit(1)
}

func checkGL(context string) {
	if err := gles2.GetError(); err != gles2.NoError {
		fmt.Fprintf(os.Stderr, "fatal: GL error 0x%04X at %s\n", int32(err), context)
		os.Exit(1)
	}
}

// glStr converts a Go string to a null-terminated GLchar slice and returns
// a pointer to the first element, suitable for passing to ShaderSource.
func glStr(s string) *gles2.GLchar {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return (*gles2.GLchar)(unsafe.Pointer(&b[0]))
}

// compileShader creates, sources, and compiles a shader of the given type.
// It returns the shader name or exits on failure.
func compileShader(shaderType gles2.Enum, source string) gles2.GLuint {
	shader := gles2.CreateShader(shaderType)
	if shader == 0 {
		fatalGL("CreateShader returned 0")
	}

	src := glStr(source)
	length := gles2.GLint(len(source))
	gles2.ShaderSource(shader, 1, &src, &length)
	gles2.CompileShader(shader)

	var status gles2.GLint
	gles2.GetShaderiv(shader, gles2.CompileStatus, &status)
	if status == 0 {
		gles2.DeleteShader(shader)
		fatalGL("shader compilation failed")
	}

	return shader
}

// linkProgram creates a program, attaches both shaders, links, and returns
// the program name or exits on failure.
func linkProgram(vertShader, fragShader gles2.GLuint) gles2.GLuint {
	program := gles2.CreateProgram()
	if program == 0 {
		fatalGL("CreateProgram returned 0")
	}

	gles2.AttachShader(program, vertShader)
	gles2.AttachShader(program, fragShader)
	gles2.LinkProgram(program)

	var status gles2.GLint
	gles2.GetProgramiv(program, gles2.LinkStatus, &status)
	if status == 0 {
		gles2.DeleteProgram(program)
		fatalGL("program linking failed")
	}

	return program
}

// Vertex shader: passes through 2D positions in clip space.
const vertSrc = `
attribute vec2 aPosition;
void main() {
    gl_Position = vec4(aPosition, 0.0, 1.0);
}
`

// Fragment shader: outputs solid red.
const fragSrc = `
precision mediump float;
void main() {
    gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
}
`

func main() {
	// --- EGL setup (offscreen 64x64 pbuffer) ---
	var defaultDisplay egl.EGLNativeDisplayType
	display := egl.GetDisplay(defaultDisplay)
	if display == nil {
		fatal("GetDisplay returned EGL_NO_DISPLAY")
	}
	var major, minor egl.Int
	if egl.Initialize(display, &major, &minor) == 0 {
		fatal("Initialize failed")
	}
	defer egl.Terminate(display)

	if egl.BindAPI(egl.EGLenum(egl.OpenglEsApi)) == 0 {
		fatal("BindAPI failed")
	}

	configAttribs := [...]egl.Int{
		egl.SurfaceType, egl.PbufferBit,
		egl.RenderableType, egl.OpenglEs2Bit,
		egl.RedSize, 8,
		egl.GreenSize, 8,
		egl.BlueSize, 8,
		egl.AlphaSize, 8,
		egl.None,
	}
	var config egl.EGLConfig
	var numConfigs egl.Int
	if egl.ChooseConfig(display, &configAttribs[0], &config, 1, &numConfigs) == 0 || numConfigs == 0 {
		fatal("ChooseConfig failed")
	}

	const width, height = 64, 64
	pbufAttribs := [...]egl.Int{
		egl.Width, width,
		egl.Height, height,
		egl.None,
	}
	surface := egl.CreatePbufferSurface(display, config, &pbufAttribs[0])
	if surface == nil {
		fatal("CreatePbufferSurface failed")
	}
	defer egl.DestroySurface(display, surface)

	ctxAttribs := [...]egl.Int{
		egl.ContextClientVersion, 2,
		egl.None,
	}
	ctx := egl.CreateContext(display, config, nil, &ctxAttribs[0])
	if ctx == nil {
		fatal("CreateContext failed")
	}
	defer func() {
		egl.MakeCurrent(display, nil, nil, nil)
		egl.DestroyContext(display, ctx)
	}()

	if egl.MakeCurrent(display, surface, surface, ctx) == 0 {
		fatal("MakeCurrent failed")
	}
	fmt.Println("GL context ready")

	// --- Compile shaders and link program ---
	vs := compileShader(gles2.VertexShader, vertSrc)
	defer gles2.DeleteShader(vs)
	fmt.Println("vertex shader compiled")

	fs := compileShader(gles2.FragmentShader, fragSrc)
	defer gles2.DeleteShader(fs)
	fmt.Println("fragment shader compiled")

	program := linkProgram(vs, fs)
	defer gles2.DeleteProgram(program)
	fmt.Println("program linked")

	// --- Look up the attribute location ---
	posLoc := gles2.GetAttribLocation(program, glStr("aPosition"))
	if posLoc < 0 {
		fatalGL("GetAttribLocation(aPosition) returned -1")
	}
	fmt.Printf("aPosition location = %d\n", posLoc)

	// --- Create a VBO with triangle vertices ---
	//
	//        (0, 0.8)
	//       /       \
	//  (-0.8,-0.8)--(0.8,-0.8)
	//
	// These clip-space coordinates produce a triangle that covers the center
	// of the viewport but leaves the corners uncovered.
	vertices := [...]float32{
		0.0, 0.8, //  top
		-0.8, -0.8, // bottom-left
		0.8, -0.8, // bottom-right
	}

	var vbo gles2.GLuint
	gles2.GenBuffers(1, &vbo)
	gles2.BindBuffer(gles2.ArrayBuffer, vbo)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(vertices)),
		unsafe.Pointer(&vertices[0]),
		gles2.StaticDraw,
	)
	defer func() {
		gles2.DeleteBuffers(1, &vbo)
	}()
	checkGL("VBO setup")
	fmt.Println("VBO created")

	// --- Render ---
	gles2.Viewport(0, 0, width, height)
	gles2.ClearColor(0.0, 0.0, 0.0, 1.0) // black background
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

	gles2.UseProgram(program)
	gles2.EnableVertexAttribArray(gles2.GLuint(posLoc))
	// 2 floats per vertex, tightly packed (stride=0).
	gles2.VertexAttribPointer(
		gles2.GLuint(posLoc),
		2,                          // components per vertex
		gles2.Float,                // type
		gles2.Boolean(gles2.False), // not normalized
		0,                          // stride (tightly packed)
		nil,                        // offset into the bound VBO
	)

	gles2.DrawArrays(gles2.Triangles, 0, 3)
	checkGL("DrawArrays")

	gles2.DisableVertexAttribArray(gles2.GLuint(posLoc))
	gles2.Finish()

	// --- Read back and verify ---
	// Center pixel should be red (inside the triangle).
	var center [4]byte
	gles2.ReadPixels(width/2, height/2, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&center[0]))
	checkGL("ReadPixels center")
	fmt.Printf("center pixel: R=%d G=%d B=%d A=%d\n", center[0], center[1], center[2], center[3])

	if center[0] == 255 && center[1] == 0 && center[2] == 0 {
		fmt.Println("verified: center pixel is red (inside triangle)")
	} else {
		fmt.Println("warning: center pixel is not the expected red")
	}

	// Corner (0,0) should be black (outside the triangle).
	var corner [4]byte
	gles2.ReadPixels(0, 0, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&corner[0]))
	fmt.Printf("corner pixel: R=%d G=%d B=%d A=%d\n", corner[0], corner[1], corner[2], corner[3])

	if corner[0] == 0 && corner[1] == 0 && corner[2] == 0 {
		fmt.Println("verified: corner pixel is black (outside triangle)")
	} else {
		fmt.Println("warning: corner pixel is not the expected black")
	}

	fmt.Println("done")
}
