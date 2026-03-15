// Textured quad: upload pixel data to a texture and render it on a quad.
//
// Demonstrates:
//   - Creating and configuring a GL texture (GenTextures, BindTexture,
//     TexParameteri, TexImage2D)
//   - Using texture coordinates in shaders to sample from the texture
//   - Drawing a quad as two triangles via an element buffer (DrawElements)
//   - Reading back pixels to verify texture sampling
//
// A 2x2 checkerboard texture (red/green/blue/white) is created from raw
// pixel data and mapped onto a fullscreen quad. The program reads back the
// four quadrant centers to verify correct texturing.
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

// Vertex shader: fullscreen quad with texture coordinates.
// Each vertex carries a 2D position and a 2D texcoord.
const vertSrc = `
attribute vec2 aPosition;
attribute vec2 aTexCoord;
varying vec2 vTexCoord;
void main() {
    gl_Position = vec4(aPosition, 0.0, 1.0);
    vTexCoord = aTexCoord;
}
`

// Fragment shader: samples a 2D texture.
const fragSrc = `
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

	// --- Create a 2x2 RGBA checkerboard texture ---
	//
	//   (0,1) green  | (1,1) white
	//   -------------|-------------
	//   (0,0) red    | (1,0) blue
	//
	// In memory, texels are stored bottom-to-top, left-to-right.
	texPixels := [...]byte{
		// Row 0 (bottom): red, blue
		255, 0, 0, 255, 0, 0, 255, 255,
		// Row 1 (top): green, white
		0, 255, 0, 255, 255, 255, 255, 255,
	}

	var tex gles2.GLuint
	gles2.GenTextures(1, &tex)
	gles2.ActiveTexture(gles2.Texture0)
	gles2.BindTexture(gles2.Texture2d, tex)
	// Use nearest filtering so that each quadrant samples exactly one texel.
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureMinFilter, gles2.GLint(gles2.Nearest))
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureMagFilter, gles2.GLint(gles2.Nearest))
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureWrapS, gles2.GLint(gles2.ClampToEdge))
	gles2.TexParameteri(gles2.Texture2d, gles2.TextureWrapT, gles2.GLint(gles2.ClampToEdge))
	gles2.TexImage2D(
		gles2.Texture2d,
		0,                       // mip level 0
		gles2.GLint(gles2.Rgba), // internal format
		2, 2,                    // width, height
		0,                  // border (must be 0 in ES2)
		gles2.Rgba,         // format
		gles2.UnsignedByte, // type
		unsafe.Pointer(&texPixels[0]),
	)
	checkGL("TexImage2D")
	defer gles2.DeleteTextures(1, &tex)
	fmt.Println("2x2 texture created")

	// --- Compile shaders and link program ---
	vs := compileShader(gles2.VertexShader, vertSrc)
	defer gles2.DeleteShader(vs)

	fs := compileShader(gles2.FragmentShader, fragSrc)
	defer gles2.DeleteShader(fs)

	program := linkProgram(vs, fs)
	defer gles2.DeleteProgram(program)
	fmt.Println("shader program linked")

	posLoc := gles2.GetAttribLocation(program, glStr("aPosition"))
	texLoc := gles2.GetAttribLocation(program, glStr("aTexCoord"))
	samplerLoc := gles2.GetUniformLocation(program, glStr("uTexture"))
	if posLoc < 0 || texLoc < 0 || samplerLoc < 0 {
		fatalGL("failed to look up shader locations")
	}

	// --- Set up quad geometry ---
	// A fullscreen quad from two triangles. Each vertex has:
	//   position (x, y) and texcoord (s, t).
	//
	//  (-1,1)---(1,1)     (0,1)---(1,1)
	//    |  \    |          |  \    |
	//    |   \   |   =>     |   \   |     texcoords
	//    |    \  |          |    \  |
	//  (-1,-1)--(1,-1)    (0,0)---(1,0)
	//
	type vertex struct {
		x, y float32 // position
		s, t float32 // texcoord
	}
	quadVerts := [...]vertex{
		{-1, -1, 0, 0}, // bottom-left
		{1, -1, 1, 0},  // bottom-right
		{1, 1, 1, 1},   // top-right
		{-1, 1, 0, 1},  // top-left
	}

	// Element indices for two triangles.
	indices := [...]uint16{
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	var vbo gles2.GLuint
	gles2.GenBuffers(1, &vbo)
	gles2.BindBuffer(gles2.ArrayBuffer, vbo)
	gles2.BufferData(
		gles2.ArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(quadVerts)),
		unsafe.Pointer(&quadVerts[0]),
		gles2.StaticDraw,
	)
	defer gles2.DeleteBuffers(1, &vbo)

	var ebo gles2.GLuint
	gles2.GenBuffers(1, &ebo)
	gles2.BindBuffer(gles2.ElementArrayBuffer, ebo)
	gles2.BufferData(
		gles2.ElementArrayBuffer,
		gles2.GLsizeiptr(unsafe.Sizeof(indices)),
		unsafe.Pointer(&indices[0]),
		gles2.StaticDraw,
	)
	defer gles2.DeleteBuffers(1, &ebo)
	checkGL("buffer setup")

	// --- Render the textured quad ---
	gles2.Viewport(0, 0, width, height)
	gles2.ClearColor(0, 0, 0, 1)
	gles2.Clear(gles2.Bitfield(gles2.ColorBufferBit))

	gles2.UseProgram(program)
	gles2.Uniform1i(samplerLoc, 0) // texture unit 0

	stride := gles2.GLsizei(unsafe.Sizeof(vertex{}))

	gles2.EnableVertexAttribArray(gles2.GLuint(posLoc))
	gles2.VertexAttribPointer(
		gles2.GLuint(posLoc),
		2, // vec2
		gles2.Float,
		gles2.Boolean(gles2.False),
		stride,
		nil, // offset 0
	)

	// texcoord offset = 2 floats = 8 bytes into each vertex
	texOffset := unsafe.Pointer(uintptr(2 * unsafe.Sizeof(float32(0))))
	gles2.EnableVertexAttribArray(gles2.GLuint(texLoc))
	gles2.VertexAttribPointer(
		gles2.GLuint(texLoc),
		2, // vec2
		gles2.Float,
		gles2.Boolean(gles2.False),
		stride,
		texOffset,
	)

	gles2.DrawElements(gles2.Triangles, 6, gles2.GL_UNSIGNED_SHORT, nil)
	checkGL("DrawElements")

	gles2.DisableVertexAttribArray(gles2.GLuint(posLoc))
	gles2.DisableVertexAttribArray(gles2.GLuint(texLoc))
	gles2.Finish()

	// --- Verify by reading back quadrant centers ---
	// The 2x2 texture maps to 4 quadrants of the 64x64 viewport.
	type check struct {
		name    string
		x, y    gles2.GLint
		r, g, b byte
	}
	checks := []check{
		{"bottom-left (red)", width / 4, height / 4, 255, 0, 0},
		{"bottom-right (blue)", 3 * width / 4, height / 4, 0, 0, 255},
		{"top-left (green)", width / 4, 3 * height / 4, 0, 255, 0},
		{"top-right (white)", 3 * width / 4, 3 * height / 4, 255, 255, 255},
	}

	allOK := true
	for _, c := range checks {
		var pixel [4]byte
		gles2.ReadPixels(c.x, c.y, 1, 1, gles2.Rgba, gles2.UnsignedByte, unsafe.Pointer(&pixel[0]))
		match := pixel[0] == c.r && pixel[1] == c.g && pixel[2] == c.b
		result := "OK"
		if !match {
			result = "MISMATCH"
			allOK = false
		}
		fmt.Printf("  %-25s got (%3d,%3d,%3d) expected (%3d,%3d,%3d) [%s]\n",
			c.name, pixel[0], pixel[1], pixel[2], c.r, c.g, c.b, result)
	}

	if allOK {
		fmt.Println("verified: all quadrants match the checkerboard texture")
	} else {
		fmt.Println("warning: some quadrants did not match")
	}

	fmt.Println("done")
}
