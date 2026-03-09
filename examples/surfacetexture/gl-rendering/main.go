// SurfaceTexture GL rendering pipeline example.
//
// Demonstrates how to use ASurfaceTexture to stream image frames from a
// producer (camera, video decoder) into an OpenGL ES external texture and
// render them on screen. This is the standard pattern for camera preview
// or video playback in Android NDK applications.
//
// ASurfaceTexture bridges the Android Surface (a queue of graphic buffers)
// and OpenGL ES by presenting each incoming frame as a GL_TEXTURE_EXTERNAL_OES
// texture. The texture coordinates require a 4x4 transform matrix that
// compensates for buffer rotation, scaling, and coordinate-system differences.
//
// The NDK ASurfaceTexture handle is obtained from a Java SurfaceTexture via
// JNI (ASurfaceTexture_fromSurfaceTexture). There is no pure-NDK constructor.
// Once the handle is wrapped with surfacetexture.NewSurfaceTextureFromPointer,
// the idiomatic API provides the full lifecycle.
//
// GL rendering workflow:
//
//  1. Create a GL_TEXTURE_EXTERNAL_OES texture on the render thread.
//
//  2. Attach the SurfaceTexture to the current GL context with
//     AttachToGLContext(texName). This binds the SurfaceTexture to the
//     texture object so incoming frames update it.
//
//  3. Acquire a NativeWindow from the SurfaceTexture with AcquireWindow().
//     Hand this window to the frame producer (camera, codec). The producer
//     enqueues frames into the Surface backed by this window.
//
//  4. On the render thread, when a new frame is available:
//     a. Call UpdateTexImage() to latch the latest frame into the texture.
//     b. Call TransformMatrix() to get the 4x4 matrix for texture coords.
//     c. Call Timestamp() to get the frame's presentation timestamp (ns).
//     d. Bind the texture, upload the matrix as a uniform, draw.
//
//  5. When switching GL contexts: DetachFromGLContext(), then re-attach.
//
//  6. When done: DetachFromGLContext(), then Close().
//
// Fragment shader for GL_TEXTURE_EXTERNAL_OES:
//
//	#extension GL_OES_EGL_image_external : require
//	precision mediump float;
//	varying vec2 vTexCoord;
//	uniform samplerExternalOES uTexture;
//	uniform mat4 uTexMatrix;
//	void main() {
//	    vec2 tc = (uTexMatrix * vec4(vTexCoord, 0.0, 1.0)).xy;
//	    gl_FragColor = texture2D(uTexture, tc);
//	}
//
// Prerequisites:
//   - Android device with API level 28+ (ASurfaceTexture was added in API 28).
//   - A Java SurfaceTexture obtained from the application's Java layer.
//   - A current EGL context with GL_OES_EGL_image_external support.
//
// Because the SurfaceTexture handle requires JNI to obtain, this example
// documents the complete rendering pipeline and prints the API calls that
// a real application would make, rather than invoking them directly.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	_ "github.com/xaionaro-go/ndk/surfacetexture"
)

func main() {
	fmt.Println("=== SurfaceTexture GL rendering pipeline ===")
	fmt.Println()

	// ── Step 1: Obtain the SurfaceTexture handle ────────────────
	//
	// The handle comes from Java through JNI. The Java side creates
	// a SurfaceTexture detached from any GL context:
	//
	//   Java:
	//     SurfaceTexture st = new SurfaceTexture(/* texName */ 0);
	//     st.detachFromGLContext();
	//     nativeSetSurfaceTexture(st);  // JNI call
	//
	//   C/JNI bridge:
	//     ASurfaceTexture* ast = ASurfaceTexture_fromSurfaceTexture(env, javaST);
	//
	//   Go:
	//     st := surfacetexture.NewSurfaceTextureFromPointer(rawPtr)

	fmt.Println("Step 1: Obtain SurfaceTexture from Java via JNI")
	fmt.Println("  Java: SurfaceTexture st = new SurfaceTexture(0);")
	fmt.Println("  JNI:  ASurfaceTexture_fromSurfaceTexture(env, javaST)")
	fmt.Println("  Go:   surfacetexture.NewSurfaceTextureFromPointer(rawPtr)")
	fmt.Println()

	// ── Step 2: Create GL texture and attach ────────────────────
	//
	//   // On the render thread with a current EGL context:
	//   var texName uint32
	//   gl.GenTextures(1, &texName)
	//   gl.BindTexture(gl.TEXTURE_EXTERNAL_OES, texName)
	//   gl.TexParameteri(gl.TEXTURE_EXTERNAL_OES, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	//   gl.TexParameteri(gl.TEXTURE_EXTERNAL_OES, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//   gl.TexParameteri(gl.TEXTURE_EXTERNAL_OES, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//   gl.TexParameteri(gl.TEXTURE_EXTERNAL_OES, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	//
	//   if err := st.AttachToGLContext(texName); err != nil {
	//       log.Fatalf("attach: %v", err)
	//   }

	fmt.Println("Step 2: Create GL_TEXTURE_EXTERNAL_OES and attach")
	fmt.Println("  gl.GenTextures(1, &texName)")
	fmt.Println("  gl.BindTexture(GL_TEXTURE_EXTERNAL_OES, texName)")
	fmt.Println("  st.AttachToGLContext(texName)")
	fmt.Println()

	// ── Step 3: Acquire window for producers ────────────────────
	//
	// The NativeWindow provides a Surface that producers write into.
	// For camera preview, pass the window's pointer to the camera
	// output target. For video, pass it to the MediaCodec output.
	//
	//   window := st.AcquireWindow()
	//   // window.Pointer() returns unsafe.Pointer to the ANativeWindow.
	//   // Pass this to camera or codec APIs as the output surface.

	fmt.Println("Step 3: Acquire NativeWindow for frame producers")
	fmt.Println("  window := st.AcquireWindow()")
	fmt.Println("  // Pass window.Pointer() to camera or video decoder output.")
	fmt.Println()

	// ── Step 4: Render loop ─────────────────────────────────────
	//
	// The render loop runs on the GL thread. When the Java-side
	// onFrameAvailable callback fires, signal the render thread to
	// call UpdateTexImage.
	//
	//   for running {
	//       // Wait for onFrameAvailable signal...
	//
	//       // Latch the latest frame into the GL texture.
	//       if err := st.UpdateTexImage(); err != nil {
	//           log.Printf("update: %v", err)
	//           continue
	//       }
	//
	//       // Get the 4x4 texture coordinate transform matrix.
	//       // This matrix compensates for buffer rotation (e.g. camera
	//       // sensor orientation) and coordinate system differences.
	//       var mtx [16]float32
	//       st.TransformMatrix(&mtx)
	//
	//       // Get the frame timestamp for A/V sync or frame pacing.
	//       timestamp := st.Timestamp()  // nanoseconds
	//
	//       // Render: bind the external texture, upload mtx as a mat4
	//       // uniform, draw a full-screen quad.
	//       gl.ActiveTexture(gl.TEXTURE0)
	//       gl.BindTexture(gl.TEXTURE_EXTERNAL_OES, texName)
	//       gl.UniformMatrix4fv(uTexMatrixLoc, 1, false, &mtx[0])
	//       // ... draw quad ...
	//       egl.SwapBuffers(display, surface)
	//   }

	fmt.Println("Step 4: Render loop")
	fmt.Println("  for each frame:")
	fmt.Println("    st.UpdateTexImage()         -- latch frame into texture")
	fmt.Println("    st.TransformMatrix(&mtx)    -- get 4x4 transform for shader")
	fmt.Println("    st.Timestamp()              -- frame timestamp in nanoseconds")
	fmt.Println("    // bind texture, upload matrix, draw quad, swap buffers")
	fmt.Println()

	// ── Step 5: Detach and close ────────────────────────────────
	//
	// When rendering is complete or the GL context is being torn down:
	//
	//   if err := st.DetachFromGLContext(); err != nil {
	//       log.Printf("detach: %v", err)
	//   }
	//   if err := st.Close(); err != nil {
	//       log.Printf("close: %v", err)
	//   }
	//
	// After Close, no other methods may be called on the SurfaceTexture.
	// The NativeWindow returned by AcquireWindow becomes invalid once the
	// SurfaceTexture is released.

	fmt.Println("Step 5: Cleanup")
	fmt.Println("  st.DetachFromGLContext()  -- unbind from GL context")
	fmt.Println("  st.Close()               -- release the NDK handle")
	fmt.Println()

	// ── Transform matrix details ────────────────────────────────
	//
	// The 4x4 matrix returned by TransformMatrix is column-major
	// (OpenGL convention). It must be applied to texture coordinates
	// in the vertex or fragment shader:
	//
	//   vec2 tc = (uTexMatrix * vec4(texCoord, 0.0, 1.0)).xy;
	//   gl_FragColor = texture2D(uTexture, tc);
	//
	// Common transforms:
	//   - Identity:         no rotation, standard orientation
	//   - Y-flip:           Surface coordinate system differs from GL
	//   - 90/180/270 rotate: camera sensor orientation compensation
	//
	// Without the transform matrix, the rendered image may appear
	// upside-down, rotated, or stretched.

	fmt.Println("Transform matrix:")
	fmt.Println("  Column-major 4x4 (OpenGL convention)")
	fmt.Println("  Apply in shader: (uTexMatrix * vec4(tc, 0, 1)).xy")
	fmt.Println("  Compensates for: Y-flip, rotation, scaling")
	fmt.Println()

	// ── Error handling ──────────────────────────────────────────
	//
	// Methods that return error use surfacetexture.Error, wrapping
	// the NDK int32 result code. A negative value indicates failure.

	fmt.Println("Error handling:")
	fmt.Println("  AttachToGLContext: fails if no EGL context or already attached")
	fmt.Println("  DetachFromGLContext: fails if not attached")
	fmt.Println("  UpdateTexImage: fails if wrong thread or not attached")
	fmt.Println("  Close: nil-safe and idempotent")
	fmt.Println()

	fmt.Println("gl-rendering pipeline overview complete")
}
