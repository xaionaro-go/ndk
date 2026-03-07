// Simulates c-for-go output for OpenGL ES 2.0.
// This file is parsed at AST level only; it does not compile.
package gles2

import "unsafe"

// Integer typedefs.
type GLenum uint32
type GLint int32
type GLuint uint32
type GLsizei int32
type GLfloat float32
type GLboolean uint8
type GLbitfield uint32
type GLclampf float32
type GLsizeiptr int64
type GLintptr int64

// Error code enum.
const (
	GL_NO_ERROR                      GLenum = 0
	GL_INVALID_ENUM                  GLenum = 0x0500
	GL_INVALID_VALUE                 GLenum = 0x0501
	GL_INVALID_OPERATION             GLenum = 0x0502
	GL_OUT_OF_MEMORY                 GLenum = 0x0505
)

// Primitive type enum.
const (
	GL_POINTS         GLenum = 0x0000
	GL_LINES          GLenum = 0x0001
	GL_LINE_STRIP     GLenum = 0x0003
	GL_TRIANGLES      GLenum = 0x0004
	GL_TRIANGLE_STRIP GLenum = 0x0005
	GL_TRIANGLE_FAN   GLenum = 0x0006
)

// Buffer target enum.
const (
	GL_ARRAY_BUFFER         GLenum = 0x8892
	GL_ELEMENT_ARRAY_BUFFER GLenum = 0x8893
)

// Buffer usage enum.
const (
	GL_STREAM_DRAW  GLenum = 0x88E0
	GL_STATIC_DRAW  GLenum = 0x88E4
	GL_DYNAMIC_DRAW GLenum = 0x88E8
)

// Shader type enum.
const (
	GL_FRAGMENT_SHADER GLenum = 0x8B30
	GL_VERTEX_SHADER   GLenum = 0x8B31
)

// Shader/program status enum.
const (
	GL_COMPILE_STATUS GLenum = 0x8B81
	GL_LINK_STATUS    GLenum = 0x8B82
)

// Texture target and parameter enum.
const (
	GL_TEXTURE_2D         GLenum = 0x0DE1
	GL_TEXTURE0           GLenum = 0x84C0
	GL_TEXTURE_MAG_FILTER GLenum = 0x2800
	GL_TEXTURE_MIN_FILTER GLenum = 0x2801
	GL_TEXTURE_WRAP_S     GLenum = 0x2802
	GL_TEXTURE_WRAP_T     GLenum = 0x2803
	GL_NEAREST            GLenum = 0x2600
	GL_LINEAR             GLenum = 0x2601
	GL_CLAMP_TO_EDGE      GLenum = 0x812F
	GL_REPEAT             GLenum = 0x2901
)

// Framebuffer and renderbuffer enum.
const (
	GL_FRAMEBUFFER          GLenum = 0x8D40
	GL_RENDERBUFFER         GLenum = 0x8D41
	GL_FRAMEBUFFER_COMPLETE GLenum = 0x8CD5
)

// Clear buffer bits.
const (
	GL_DEPTH_BUFFER_BIT   GLbitfield = 0x00000100
	GL_STENCIL_BUFFER_BIT GLbitfield = 0x00000400
	GL_COLOR_BUFFER_BIT   GLbitfield = 0x00004000
)

// Capability enum.
const (
	GL_DEPTH_TEST   GLenum = 0x0B71
	GL_STENCIL_TEST GLenum = 0x0B90
	GL_SCISSOR_TEST GLenum = 0x0C11
	GL_BLEND        GLenum = 0x0BE2
	GL_CULL_FACE    GLenum = 0x0B44
)

// Pixel format enum.
const (
	GL_RGB  GLenum = 0x1907
	GL_RGBA GLenum = 0x1908
)

// Data type enum.
const (
	GL_UNSIGNED_BYTE  GLenum = 0x1401
	GL_UNSIGNED_SHORT GLenum = 0x1403
	GL_FLOAT          GLenum = 0x1406
)

// Boolean enum.
const (
	GL_FALSE GLboolean = 0
	GL_TRUE  GLboolean = 1
)

// Blend factor enum.
const (
	GL_ZERO                GLenum = 0
	GL_ONE                 GLenum = 1
	GL_SRC_ALPHA           GLenum = 0x0302
	GL_ONE_MINUS_SRC_ALPHA GLenum = 0x0303
)

// --- State functions ---
func GlClearColor(red GLclampf, green GLclampf, blue GLclampf, alpha GLclampf) {}
func GlClear(mask GLbitfield)                                                  {}
func GlEnable(cap GLenum)                                                      {}
func GlDisable(cap GLenum)                                                     {}
func GlViewport(x GLint, y GLint, width GLsizei, height GLsizei)              {}
func GlScissor(x GLint, y GLint, width GLsizei, height GLsizei)               {}

// --- Shader functions ---
func GlCreateShader(shaderType GLenum) GLuint                                  { return 0 }
func GlDeleteShader(shader GLuint)                                             {}
func GlShaderSource(shader GLuint, count GLsizei, str **byte, length *GLint)   {}
func GlCompileShader(shader GLuint)                                            {}
func GlGetShaderiv(shader GLuint, pname GLenum, params *GLint)                 {}

// --- Program functions ---
func GlCreateProgram() GLuint                                                  { return 0 }
func GlDeleteProgram(program GLuint)                                           {}
func GlAttachShader(program GLuint, shader GLuint)                             {}
func GlLinkProgram(program GLuint)                                             {}
func GlUseProgram(program GLuint)                                              {}
func GlGetProgramiv(program GLuint, pname GLenum, params *GLint)               {}

// --- Uniform functions ---
func GlGetUniformLocation(program GLuint, name *byte) GLint                    { return 0 }
func GlUniform1f(location GLint, v0 GLfloat)                                   {}
func GlUniform1i(location GLint, v0 GLint)                                     {}
func GlUniform4fv(location GLint, count GLsizei, value *GLfloat)               {}
func GlUniformMatrix4fv(location GLint, count GLsizei, transpose GLboolean, value *GLfloat) {}

// --- Vertex attribute functions ---
func GlGetAttribLocation(program GLuint, name *byte) GLint                     { return 0 }
func GlEnableVertexAttribArray(index GLuint)                                   {}
func GlDisableVertexAttribArray(index GLuint)                                  {}
func GlVertexAttribPointer(index GLuint, size GLint, xtype GLenum, normalized GLboolean, stride GLsizei, pointer unsafe.Pointer) {}

// --- Buffer functions ---
func GlGenBuffers(n GLsizei, buffers *GLuint)                                  {}
func GlDeleteBuffers(n GLsizei, buffers *GLuint)                               {}
func GlBindBuffer(target GLenum, buffer GLuint)                                {}
func GlBufferData(target GLenum, size GLsizeiptr, data unsafe.Pointer, usage GLenum) {}
func GlBufferSubData(target GLenum, offset GLintptr, size GLsizeiptr, data unsafe.Pointer) {}

// --- Texture functions ---
func GlGenTextures(n GLsizei, textures *GLuint)                                {}
func GlDeleteTextures(n GLsizei, textures *GLuint)                             {}
func GlBindTexture(target GLenum, texture GLuint)                              {}
func GlTexImage2D(target GLenum, level GLint, internalformat GLint, width GLsizei, height GLsizei, border GLint, format GLenum, xtype GLenum, pixels unsafe.Pointer) {}
func GlTexParameteri(target GLenum, pname GLenum, param GLint)                 {}
func GlActiveTexture(texture GLenum)                                           {}

// --- Draw functions ---
func GlDrawArrays(mode GLenum, first GLint, count GLsizei)                     {}
func GlDrawElements(mode GLenum, count GLsizei, xtype GLenum, indices unsafe.Pointer) {}

// --- Framebuffer functions ---
func GlGenFramebuffers(n GLsizei, framebuffers *GLuint)                        {}
func GlDeleteFramebuffers(n GLsizei, framebuffers *GLuint)                     {}
func GlBindFramebuffer(target GLenum, framebuffer GLuint)                      {}
func GlFramebufferTexture2D(target GLenum, attachment GLenum, textarget GLenum, texture GLuint, level GLint) {}
func GlCheckFramebufferStatus(target GLenum) GLenum                            { return 0 }

// --- Renderbuffer functions ---
func GlGenRenderbuffers(n GLsizei, renderbuffers *GLuint)                      {}
func GlDeleteRenderbuffers(n GLsizei, renderbuffers *GLuint)                   {}
func GlBindRenderbuffer(target GLenum, renderbuffer GLuint)                    {}
func GlRenderbufferStorage(target GLenum, internalformat GLenum, width GLsizei, height GLsizei) {}
func GlFramebufferRenderbuffer(target GLenum, attachment GLenum, renderbuffertarget GLenum, renderbuffer GLuint) {}

// --- Query functions ---
func GlGetError() GLenum                                                       { return 0 }
func GlGetString(name GLenum) *byte                                            { return nil }
func GlGetIntegerv(pname GLenum, data *GLint)                                  {}
func GlFlush()                                                                 {}
func GlFinish()                                                                {}

// --- Blending and depth functions ---
func GlBlendFunc(sfactor GLenum, dfactor GLenum)                               {}
func GlBlendFuncSeparate(srcRGB GLenum, dstRGB GLenum, srcAlpha GLenum, dstAlpha GLenum) {}
func GlDepthFunc(fn GLenum)                                                    {}
func GlDepthMask(flag GLboolean)                                               {}
func GlStencilFunc(fn GLenum, ref GLint, mask GLuint)                          {}
func GlStencilOp(fail GLenum, zfail GLenum, zpass GLenum)                      {}

// --- Pixel functions ---
func GlPixelStorei(pname GLenum, param GLint)                                  {}
func GlReadPixels(x GLint, y GLint, width GLsizei, height GLsizei, format GLenum, xtype GLenum, pixels unsafe.Pointer) {}
func GlLineWidth(width GLfloat)                                                {}

var _ = unsafe.Pointer(nil)
