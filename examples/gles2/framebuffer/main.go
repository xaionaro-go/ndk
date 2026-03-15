// Render-to-texture via an offscreen framebuffer object (FBO).
//
// Demonstrates:
//   - Creating a texture to serve as a color attachment
//   - Creating an FBO and attaching the texture with FramebufferTexture2D
//   - Checking framebuffer completeness with CheckFramebufferStatus
//   - Rendering a triangle into the FBO (offscreen)
//   - Reading back pixels from the FBO with ReadPixels
//   - Binding the FBO texture and sampling it in a second pass
//
// The program performs two render passes:
//  1. Render a green triangle into a 32x32 FBO texture.
//  2. Bind that texture and draw a fullscreen quad onto the default
//     framebuffer (64x64 pbuffer), effectively displaying the FBO contents.
//
// Both passes are verified by reading back pixels.
package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/xaionaro-go/ndk/egl"
	"github.com/xaionaro-go/ndk/gles2"
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

func glStr(s string) *gles2.GLchar {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return (*gles2.GLchar)(unsafe.Pointer(&b[0]))
}

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

// --- Shaders for pass 1: solid-color triangle ---

const pass1VertSrc = `
attribute vec2 aPosition;
void main() {
    gl_Position = vec4(aPosition, 0.0, 1.0);
}
`

const pass1FragSrc = `
precision mediump float;
void main() {
    gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
}
`

// --- Shaders for pass 2: textured fullscreen quad ---

const pass2VertSrc = `
attribute vec2 aPosition;
attribute vec2 aTexCoord;
varying vec2 vTexCoord;
void main() {
    gl_Position = vec4(aPosition, 0.0, 1.0);
    vTexCoord = aTexCoord;
}
`

const pass2FragSrc = `
precision mediump float;
varying vec2 vTexCoord;
uniform sampler2D uTexture;
void main() {
    gl_FragColor = texture2D(uTexture, vTexCoord);
}
`

func main() {
	// --- EGL setup ---
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

	const defaultWidth, defaultHeight = 64, 64
	pbufAttribs := [...]egl.Int{
		egl.Width, defaultWidth,
		egl.Height, defaultHeight,
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

	// =====================================================================
	// Create FBO with a color texture attachment
	// =====================================================================

	const fboWidth, fboHeight = 32, 32

	// Create the texture that will receive the FBO output.
	var fboTex gles2.GLuint
	gles2.GenTextures(1, &fboTex)
	gles2.BindTexture(gles2.Texture2d, fboTex)
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureMinFilter, gles2.GLint(gles2.Nearest))
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureMagFilter, gles2.GLint(gles2.Nearest))
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureWrapS, gles2.GLint(gles2.ClampToEdge))
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureWrapT, gles2.GLint(gles2.ClampToEdge))
	// Allocate storage with no initial data.
	gles2.TexImage2D(
		gles2.Texture2d, 0, gles2.GLint(gles2.Rgba),
		fboWidth, fboHeight, 0,
		gles2.Rgba, gles2.UnsignedByte, nil,
	)
	checkGL("FBO texture creation")
	defer gles2.DeleteTextures(1, &fboTex)

	// Create the framebuffer and attach the texture.
	var fbo gles2.GLuint
	gles2.GenFramebuffers(1, &fbo)
	gles2.BindFramebuffer(gles2.Framebuffer, fbo)
	gles2.FramebufferTexture2D(
		gles2.Framebuffer, gles2.ColorAttachment0,
		gles2.Texture2d, fboTex, 0,
	)

	status := gles2.CheckFramebufferStatus(gles2.Framebuffer)
	if status != gles2.FramebufferComplete {
		fmt.Fprintf(os.Stderr, "fatal: framebuffer incomplete (status 0x%04X)\n", int32(status))
		os.Exit(1)
	}
	defer gles2.DeleteFramebuffers(1, &fbo)
	fmt.Println("FBO created and complete")

	// =====================================================================
	// Pass 1: Render a green triangle into the FBO
	// =====================================================================

	vs1 := compileShader(gles2.VertexShader, pass1VertSrc)
	defer gles2.DeleteShader(vs1)
	fs1 := compileShader(gles2.FragmentShader, pass1FragSrc)
	defer gles2.DeleteShader(fs1)
	prog1 := linkProgram(vs1, fs1)
	defer gles2.DeleteProgram(prog1)

	posLoc1 := gles2.GetAttribLocation(prog1, glStr("aPosition"))
	if posLoc1 < 0 {
		fatalGL("GetAttribLocation(aPosition) for pass1")
	}

	// Triangle vertices (clip space).
	triVerts := [...]float32{
		0.0, 0.8,
		-0.8, -0.8,
		0.8, -0.8,
	}
	var triVBO gles2.GLuint
	gles2.GenBuffers(1, &triVBO)
	gles2.BindBuffer(gles2.ArrayBuffer, triVBO)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(triVerts)),
		unsafe.Pointer(&triVerts[0]),
		gles2.StaticDraw,
	)
	defer gles2.DeleteBuffers(1, &triVBO)

	// Render into the FBO (already bound).
	gles2.Viewport(0, 0, fboWidth, fboHeight)
	gles2.ClearColor(0.2, 0.2, 0.2, 1.0) // dark gray background
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

	gles2.UseProgram(prog1)
	gles2.EnableVertexAttribArray(gles2.GLuint(posLoc1))
	gles2.VertexAttribPointer(gles2.GLuint(posLoc1), 2, gles2.Float, gles2.Boolean(gles2.False), 0, nil)
	gles2.DrawArrays(gles2.Triangles, 0, 3)
	gles2.DisableVertexAttribArray(gles2.GLuint(posLoc1))
	gles2.Finish()
	checkGL("pass 1 render")

	// Verify: read the center pixel from the FBO (should be green).
	var fboCenter [4]byte
	gles2.ReadPixels(fboWidth/2, fboHeight/2, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&fboCenter[0]))
	fmt.Printf("FBO center pixel: R=%d G=%d B=%d A=%d\n", fboCenter[0], fboCenter[1], fboCenter[2], fboCenter[3])
	if fboCenter[1] == 255 && fboCenter[0] == 0 && fboCenter[2] == 0 {
		fmt.Println("verified: FBO center is green (triangle rendered correctly)")
	}

	// Verify: read a corner pixel from the FBO (should be dark gray).
	var fboCorner [4]byte
	gles2.ReadPixels(0, 0, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&fboCorner[0]))
	fmt.Printf("FBO corner pixel: R=%d G=%d B=%d A=%d\n", fboCorner[0], fboCorner[1], fboCorner[2], fboCorner[3])

	// =====================================================================
	// Pass 2: Draw the FBO texture onto the default framebuffer
	// =====================================================================

	vs2 := compileShader(gles2.VertexShader, pass2VertSrc)
	defer gles2.DeleteShader(vs2)
	fs2 := compileShader(gles2.FragmentShader, pass2FragSrc)
	defer gles2.DeleteShader(fs2)
	prog2 := linkProgram(vs2, fs2)
	defer gles2.DeleteProgram(prog2)

	posLoc2 := gles2.GetAttribLocation(prog2, glStr("aPosition"))
	texLoc2 := gles2.GetAttribLocation(prog2, glStr("aTexCoord"))
	samplerLoc2 := gles2.GetUniformLocation(prog2, glStr("uTexture"))
	if posLoc2 < 0 || texLoc2 < 0 || samplerLoc2 < 0 {
		fatalGL("failed to look up pass2 shader locations")
	}

	// Fullscreen quad: position + texcoord interleaved.
	type vertex struct {
		x, y float32
		s, t float32
	}
	quadVerts := [...]vertex{
		{-1, -1, 0, 0},
		{1, -1, 1, 0},
		{1, 1, 1, 1},
		{-1, 1, 0, 1},
	}
	quadIndices := [...]uint16{0, 1, 2, 0, 2, 3}

	var quadVBO gles2.GLuint
	gles2.GenBuffers(1, &quadVBO)
	gles2.BindBuffer(gles2.ArrayBuffer, quadVBO)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(quadVerts)),
		unsafe.Pointer(&quadVerts[0]),
		gles2.StaticDraw,
	)
	defer gles2.DeleteBuffers(1, &quadVBO)

	var quadEBO gles2.GLuint
	gles2.GenBuffers(1, &quadEBO)
	gles2.BindBuffer(gles2.ElementArrayBuffer, quadEBO)
	gles2.BufferData(
		gles2.ElementArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(quadIndices)),
		unsafe.Pointer(&quadIndices[0]),
		gles2.StaticDraw,
	)
	defer gles2.DeleteBuffers(1, &quadEBO)

	// Switch to the default framebuffer.
	gles2.BindFramebuffer(gles2.Framebuffer, 0)
	gles2.Viewport(0, 0, defaultWidth, defaultHeight)
	gles2.ClearColor(0, 0, 0, 1)
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

	// Bind the FBO texture.
	gles2.ActiveTexture(gles2.Texture0)
	gles2.BindTexture(gles2.Texture2d, fboTex)

	gles2.UseProgram(prog2)
	gles2.Uniform1i(samplerLoc2, 0)

	stride := gles2.GLsizei(unsafe.Sizeof(vertex{}))

	gles2.EnableVertexAttribArray(gles2.GLuint(posLoc2))
	gles2.VertexAttribPointer(gles2.GLuint(posLoc2), 2, gles2.Float, gles2.Boolean(gles2.False), stride, nil)

	texOffset := unsafe.Pointer(uintptr(2 * unsafe.Sizeof(float32(0))))
	gles2.EnableVertexAttribArray(gles2.GLuint(texLoc2))
	gles2.VertexAttribPointer(gles2.GLuint(texLoc2), 2, gles2.Float, gles2.Boolean(gles2.False), stride, texOffset)

	gles2.DrawElements(gles2.Triangles, 6, gles2.GL_UNSIGNED_SHORT, nil)
	checkGL("pass 2 render")

	gles2.DisableVertexAttribArray(gles2.GLuint(posLoc2))
	gles2.DisableVertexAttribArray(gles2.GLuint(texLoc2))
	gles2.Finish()

	// Verify: the default framebuffer should now contain the FBO image.
	var defaultCenter [4]byte
	gles2.ReadPixels(defaultWidth/2, defaultHeight/2, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&defaultCenter[0]))
	fmt.Printf("default FB center pixel: R=%d G=%d B=%d A=%d\n",
		defaultCenter[0], defaultCenter[1], defaultCenter[2], defaultCenter[3])

	if defaultCenter[1] == 255 && defaultCenter[0] == 0 && defaultCenter[2] == 0 {
		fmt.Println("verified: FBO contents transferred to default framebuffer")
	}

	fmt.Println("done")
}
