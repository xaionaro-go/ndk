// ASurfaceTexture API overview.
//
// Documents the Android ASurfaceTexture lifecycle and the methods provided by
// the ndk/surfacetexture package. ASurfaceTexture streams image frames from
// an Android Surface into an OpenGL ES texture (GL_TEXTURE_EXTERNAL_OES),
// enabling GPU-accelerated rendering of camera preview, video playback, or
// any Surface-backed producer.
//
// The NDK ASurfaceTexture handle is obtained from a Java SurfaceTexture
// object via JNI. There is no pure-NDK constructor. The typical pattern is:
//
//	Java side:
//	  SurfaceTexture st = new SurfaceTexture(/* texName */ 0);
//	  st.detachFromGLContext();   // detach so native can re-attach
//
//	JNI bridge:
//	  ASurfaceTexture* ast = ASurfaceTexture_fromSurfaceTexture(env, javaST);
//
//	Go side:
//	  // Wrap the raw *ASurfaceTexture into surfacetexture.SurfaceTexture
//	  st := surfacetexture.NewSurfaceTextureFromPointer(rawPtr)
//
// Because the SurfaceTexture handle requires a JNI environment and a Java
// SurfaceTexture jobject, this example documents the workflow and prints the
// API surface rather than calling the methods directly.
//
// GL texture workflow:
//
//  1. Create GL_TEXTURE_EXTERNAL_OES texture:
//     glGenTextures(1, &texName)
//     glBindTexture(GL_TEXTURE_EXTERNAL_OES, texName)
//
//  2. Attach the SurfaceTexture to the GL context:
//     st.AttachToGLContext(texName)
//
//  3. When a new frame arrives (via onFrameAvailable callback):
//     st.UpdateTexImage()
//
//  4. Retrieve the 4x4 texture transform matrix for the shader:
//     var mtx [16]float32
//     st.TransformMatrix(&mtx[0])
//
//  5. Retrieve the frame timestamp (nanoseconds):
//     ts := st.Timestamp()
//
//  6. Render using the external texture with the transform matrix
//     applied in the fragment shader via samplerExternalOES.
//
//  7. When switching GL contexts or done:
//     st.DetachFromGLContext()
//
//  8. Release the SurfaceTexture:
//     st.Close()
//
// Fragment shader for external textures:
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
// This program must run on an Android device.
package main

import (
	"fmt"

	_ "github.com/AndroidGoLab/ndk/surfacetexture"
)

func main() {
	fmt.Println("=== ndk/surfacetexture API overview ===")
	fmt.Println()

	// ---------------------------------------------------------------
	// Obtaining a SurfaceTexture
	//
	// The SurfaceTexture struct wraps the NDK ASurfaceTexture pointer.
	// It is NOT created from Go code. The handle comes from the Java
	// layer through JNI:
	//
	//   // C/JNI:
	//   ASurfaceTexture* ast = ASurfaceTexture_fromSurfaceTexture(
	//       env,         // JNIEnv*
	//       javaST);     // jobject (android.graphics.SurfaceTexture)
	//
	// ASurfaceTexture_fromSurfaceTexture is intentionally excluded from
	// the idiomatic Go layer because it requires JNI types (JNIEnv*,
	// jobject) that are specific to the application's JNI bridge. The
	// application wraps the returned pointer using
	// surfacetexture.NewSurfaceTextureFromPointer in its own JNI glue.
	// ---------------------------------------------------------------

	fmt.Println("SurfaceTexture handle:")
	fmt.Println("  Obtained via JNI: ASurfaceTexture_fromSurfaceTexture(env, javaST)")
	fmt.Println("  The Java SurfaceTexture must be created on the Java side and")
	fmt.Println("  passed to native code through a JNI call.")
	fmt.Println()

	// ---------------------------------------------------------------
	// SurfaceTexture methods
	// ---------------------------------------------------------------

	fmt.Println("SurfaceTexture methods:")
	fmt.Println()
	fmt.Println("  AttachToGLContext(texName uint32) error")
	fmt.Println("      Attach the SurfaceTexture to the current EGL context,")
	fmt.Println("      using texName as the GL_TEXTURE_EXTERNAL_OES texture.")
	fmt.Println("      The texture must already be created via glGenTextures.")
	fmt.Println("      The SurfaceTexture must be detached before calling this.")
	fmt.Println()
	fmt.Println("  DetachFromGLContext() error")
	fmt.Println("      Detach from the current EGL context. The texture name")
	fmt.Println("      is released and the SurfaceTexture can be re-attached")
	fmt.Println("      to a different context or texture name.")
	fmt.Println()
	fmt.Println("  UpdateTexImage() error")
	fmt.Println("      Update the texture image to the most recent frame from")
	fmt.Println("      the image stream. Call this when a new frame is available")
	fmt.Println("      (signaled by the onFrameAvailable callback on the Java")
	fmt.Println("      side). Must be called from the thread that owns the GL")
	fmt.Println("      context the SurfaceTexture is attached to.")
	fmt.Println()
	fmt.Println("  TransformMatrix(mtx *[16]float32)")
	fmt.Println("      Retrieve the 4x4 texture coordinate transform matrix.")
	fmt.Println("      The matrix compensates for any rotation, scaling, or")
	fmt.Println("      coordinate-system differences between the producer and")
	fmt.Println("      GL. Apply it in the shader to the texture coordinates")
	fmt.Println("      before sampling.")
	fmt.Println()
	fmt.Println("  Timestamp() int64")
	fmt.Println("      Return the timestamp (in nanoseconds) of the current")
	fmt.Println("      texture image, as set by the image producer. For camera")
	fmt.Println("      frames this is the capture time; for video it is the")
	fmt.Println("      presentation timestamp.")
	fmt.Println()
	fmt.Println("  AcquireWindow() *NativeWindow")
	fmt.Println("      Acquire an ANativeWindow from this SurfaceTexture.")
	fmt.Println("      Producers (camera, codec) write frames into this window.")
	fmt.Println()
	fmt.Println("  Close() error")
	fmt.Println("      Release the underlying NDK handle. After Close, no")
	fmt.Println("      other methods may be called. Idempotent and nil-safe.")
	fmt.Println()

	// ---------------------------------------------------------------
	// GL texture workflow
	// ---------------------------------------------------------------

	fmt.Println("GL texture workflow (pseudocode):")
	fmt.Println()
	fmt.Println("  // 1. Create the external texture.")
	fmt.Println("  var texName uint32")
	fmt.Println("  glGenTextures(1, &texName)")
	fmt.Println("  glBindTexture(GL_TEXTURE_EXTERNAL_OES, texName)")
	fmt.Println("  glTexParameteri(GL_TEXTURE_EXTERNAL_OES, GL_TEXTURE_MIN_FILTER, GL_LINEAR)")
	fmt.Println("  glTexParameteri(GL_TEXTURE_EXTERNAL_OES, GL_TEXTURE_MAG_FILTER, GL_LINEAR)")
	fmt.Println()
	fmt.Println("  // 2. Attach the SurfaceTexture to the current GL context.")
	fmt.Println("  if err := st.AttachToGLContext(texName); err != nil {")
	fmt.Println("      log.Fatalf(\"attach: %v\", err)")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  // 3. Acquire an ANativeWindow for producers to write into.")
	fmt.Println("  window := st.AcquireWindow()")
	fmt.Println("  // Pass this window to the camera or video decoder.")
	fmt.Println()
	fmt.Println("  // 4. Render loop: when a new frame arrives:")
	fmt.Println("  if err := st.UpdateTexImage(); err != nil {")
	fmt.Println("      log.Printf(\"update: %v\", err)")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  var mtx [16]float32")
	fmt.Println("  st.TransformMatrix(&mtx)")
	fmt.Println()
	fmt.Println("  ts := st.Timestamp()")
	fmt.Println("  _ = ts // use for A/V sync or frame pacing")
	fmt.Println()
	fmt.Println("  // Upload mtx as a mat4 uniform; bind texName as")
	fmt.Println("  // GL_TEXTURE_EXTERNAL_OES; draw the quad.")
	fmt.Println()
	fmt.Println("  // 5. Cleanup.")
	fmt.Println("  if err := st.DetachFromGLContext(); err != nil {")
	fmt.Println("      log.Printf(\"detach: %v\", err)")
	fmt.Println("  }")
	fmt.Println("  st.Close()")
	fmt.Println()

	// ---------------------------------------------------------------
	// Camera preview pipeline
	// ---------------------------------------------------------------

	fmt.Println("Camera preview pipeline (typical integration):")
	fmt.Println()
	fmt.Println("  Java:")
	fmt.Println("    SurfaceTexture st = new SurfaceTexture(0);")
	fmt.Println("    st.detachFromGLContext();")
	fmt.Println("    nativeInit(st);  // pass to native via JNI")
	fmt.Println()
	fmt.Println("  Native JNI bridge (C):")
	fmt.Println("    ASurfaceTexture* ast =")
	fmt.Println("        ASurfaceTexture_fromSurfaceTexture(env, javaST);")
	fmt.Println()
	fmt.Println("  Go side:")
	fmt.Println("    st := surfacetexture.NewSurfaceTextureFromPointer(rawPtr)")
	fmt.Println("    window := st.AcquireWindow()")
	fmt.Println("    // Pass window to ACameraDevice as output target.")
	fmt.Println()
	fmt.Println("  Go render thread:")
	fmt.Println("    st.AttachToGLContext(texName)")
	fmt.Println("    for frame := range frames {")
	fmt.Println("        st.UpdateTexImage()")
	fmt.Println("        st.TransformMatrix(&mtx)")
	fmt.Println("        // render with external texture + transform matrix")
	fmt.Println("    }")
	fmt.Println("    st.DetachFromGLContext()")
	fmt.Println("    st.Close()")
	fmt.Println()

	// ---------------------------------------------------------------
	// Error handling
	// ---------------------------------------------------------------

	fmt.Println("Error handling:")
	fmt.Println()
	fmt.Println("  Methods that return error use surfacetexture.Error, which")
	fmt.Println("  wraps the NDK int32 result code. A negative value indicates")
	fmt.Println("  failure. Common causes:")
	fmt.Println("    - AttachToGLContext: no current EGL context, or the")
	fmt.Println("      SurfaceTexture is already attached.")
	fmt.Println("    - DetachFromGLContext: not currently attached.")
	fmt.Println("    - UpdateTexImage: called from wrong thread, or the")
	fmt.Println("      SurfaceTexture is not attached to a GL context.")
	fmt.Println()

	fmt.Println("overview complete -- see source comments for full details")
}
