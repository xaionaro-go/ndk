// Simulates c-for-go output for OpenGL ES 3.0 (additions beyond ES 2.0).
// This file is parsed at AST level only; it does not compile.
package gles3

import "unsafe"

// Integer typedefs (GL types).
type GLenum uint32
type GLint int32
type GLuint uint32
type GLsizei int32
type GLfloat float32
type GLsizeiptr int64
type GLintptr int64
type GLbitfield uint32
type GLboolean uint8
type GLsync unsafe.Pointer

// Buffer target and framebuffer enum.
const (
	GL_READ_BUFFER       GLenum = 0x0C02
	GL_DRAW_BUFFER       GLenum = 0x0C01
	GL_READ_FRAMEBUFFER  GLenum = 0x8CA8
	GL_DRAW_FRAMEBUFFER  GLenum = 0x8CA9
	GL_COLOR_ATTACHMENT0 GLenum = 0x8CE0
	GL_DEPTH_ATTACHMENT  GLenum = 0x8D00
	GL_STENCIL_ATTACHMENT GLenum = 0x8D20
	GL_FRAMEBUFFER_DEFAULT GLenum = 0x8218
)

// Buffer mapping bits enum.
const (
	GL_MAP_READ_BIT              GLbitfield = 0x0001
	GL_MAP_WRITE_BIT             GLbitfield = 0x0002
	GL_MAP_INVALIDATE_RANGE_BIT  GLbitfield = 0x0004
	GL_MAP_INVALIDATE_BUFFER_BIT GLbitfield = 0x0008
	GL_MAP_FLUSH_EXPLICIT_BIT    GLbitfield = 0x0010
	GL_MAP_UNSYNCHRONIZED_BIT    GLbitfield = 0x0020
)

// Buffer binding target enum.
const (
	GL_UNIFORM_BUFFER            GLenum = 0x8A11
	GL_TRANSFORM_FEEDBACK_BUFFER GLenum = 0x8C8E
	GL_PIXEL_PACK_BUFFER         GLenum = 0x88EB
	GL_PIXEL_UNPACK_BUFFER       GLenum = 0x88EC
	GL_COPY_READ_BUFFER          GLenum = 0x8F36
	GL_COPY_WRITE_BUFFER         GLenum = 0x8F37
)

// Texture target and format enum.
const (
	GL_TEXTURE_3D       GLenum = 0x806F
	GL_TEXTURE_2D_ARRAY GLenum = 0x8C1A
	GL_TEXTURE_WRAP_R   GLenum = 0x8072
	GL_R8               GLenum = 0x8229
	GL_RG8              GLenum = 0x822B
	GL_RGB8             GLenum = 0x8051
	GL_RGBA8            GLenum = 0x8058
	GL_R16F             GLenum = 0x822D
	GL_RG16F            GLenum = 0x822F
	GL_RGBA16F          GLenum = 0x881A
	GL_R32F             GLenum = 0x822E
	GL_RG32F            GLenum = 0x8230
	GL_RGBA32F          GLenum = 0x8814
)

// Vertex array binding enum.
const (
	GL_VERTEX_ARRAY_BINDING GLenum = 0x85B5
)

// Sync enum.
const (
	GL_SYNC_GPU_COMMANDS_COMPLETE GLenum     = 0x9117
	GL_ALREADY_SIGNALED           GLenum     = 0x911A
	GL_TIMEOUT_EXPIRED            GLenum     = 0x911B
	GL_CONDITION_SATISFIED        GLenum     = 0x911C
	GL_WAIT_FAILED                GLenum     = 0x911D
	GL_SYNC_FLUSH_COMMANDS_BIT    GLbitfield = 0x00000001
	GL_TIMEOUT_IGNORED            int64      = -1
)

// Transform feedback and query enum.
const (
	GL_TRANSFORM_FEEDBACK  GLenum = 0x8E22
	GL_ANY_SAMPLES_PASSED  GLenum = 0x8C2F
	GL_SAMPLER_BINDING     GLenum = 0x8919
)

// --- Vertex array object functions ---
func GlGenVertexArrays(n GLsizei, arrays *GLuint)    {}
func GlDeleteVertexArrays(n GLsizei, arrays *GLuint) {}
func GlBindVertexArray(array GLuint)                  {}

// --- Buffer mapping functions ---
func GlMapBufferRange(target GLenum, offset GLintptr, length GLsizeiptr, access GLbitfield) unsafe.Pointer { return nil }
func GlUnmapBuffer(target GLenum) GLboolean                                                               { return 0 }
func GlFlushMappedBufferRange(target GLenum, offset GLintptr, length GLsizeiptr)                          {}

// --- Framebuffer functions ---
func GlBlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 GLint, mask GLbitfield, filter GLenum) {}
func GlRenderbufferStorageMultisample(target GLenum, samples GLsizei, internalformat GLenum, width, height GLsizei)   {}
func GlFramebufferTextureLayer(target, attachment GLenum, texture GLuint, level GLint, layer GLint)                    {}
func GlInvalidateFramebuffer(target GLenum, numAttachments GLsizei, attachments *GLenum)                              {}

// --- Buffer binding functions ---
func GlBindBufferBase(target GLenum, index GLuint, buffer GLuint)                                {}
func GlBindBufferRange(target GLenum, index GLuint, buffer GLuint, offset GLintptr, size GLsizeiptr) {}
func GlCopyBufferSubData(readTarget, writeTarget GLenum, readOffset, writeOffset GLintptr, size GLsizeiptr) {}

// --- Transform feedback functions ---
func GlBeginTransformFeedback(primitiveMode GLenum)                                    {}
func GlEndTransformFeedback()                                                          {}
func GlTransformFeedbackVaryings(program GLuint, count GLsizei, varyings **byte, bufferMode GLenum) {}
func GlGenTransformFeedbacks(n GLsizei, ids *GLuint)                                   {}
func GlDeleteTransformFeedbacks(n GLsizei, ids *GLuint)                                {}
func GlBindTransformFeedback(target GLenum, id GLuint)                                 {}

// --- Sync functions ---
func GlFenceSync(condition GLenum, flags GLbitfield) GLsync           { return nil }
func GlDeleteSync(sync GLsync)                                       {}
func GlClientWaitSync(sync GLsync, flags GLbitfield, timeout uint64) GLenum { return 0 }
func GlWaitSync(sync GLsync, flags GLbitfield, timeout uint64)       {}

// --- Sampler functions ---
func GlGenSamplers(count GLsizei, samplers *GLuint)                    {}
func GlDeleteSamplers(count GLsizei, samplers *GLuint)                 {}
func GlBindSampler(unit GLuint, sampler GLuint)                        {}
func GlSamplerParameteri(sampler GLuint, pname GLenum, param GLint)    {}

// --- Query functions ---
func GlBeginQuery(target GLenum, id GLuint)                             {}
func GlEndQuery(target GLenum)                                          {}
func GlGenQueries(n GLsizei, ids *GLuint)                               {}
func GlDeleteQueries(n GLsizei, ids *GLuint)                            {}
func GlGetQueryObjectuiv(id GLuint, pname GLenum, params *GLuint)       {}

// --- Texture functions ---
func GlTexStorage2D(target GLenum, levels GLsizei, internalformat GLenum, width, height GLsizei)               {}
func GlTexStorage3D(target GLenum, levels GLsizei, internalformat GLenum, width, height, depth GLsizei)        {}
func GlTexSubImage3D(target GLenum, level GLint, xoffset, yoffset, zoffset GLint, width, height, depth GLsizei, format, typ GLenum, pixels unsafe.Pointer) {}

// --- Draw functions ---
func GlDrawArraysInstanced(mode GLenum, first GLint, count GLsizei, instancecount GLsizei)                    {}
func GlDrawElementsInstanced(mode GLenum, count GLsizei, typ GLenum, indices unsafe.Pointer, instancecount GLsizei) {}
func GlDrawRangeElements(mode GLenum, start, end GLuint, count GLsizei, typ GLenum, indices unsafe.Pointer)   {}

// --- Clear functions ---
func GlClearBufferfv(buffer GLenum, drawbuffer GLint, value *GLfloat)           {}
func GlClearBufferiv(buffer GLenum, drawbuffer GLint, value *GLint)             {}
func GlClearBufferfi(buffer GLenum, drawbuffer GLint, depth GLfloat, stencil GLint) {}

// --- Read/Draw buffer functions ---
func GlReadBuffer(src GLenum)                       {}
func GlDrawBuffers(n GLsizei, bufs *GLenum)         {}

// --- String query functions ---
func GlGetStringi(name GLenum, index GLuint) *byte { return nil }

var _ = unsafe.Pointer(nil)
