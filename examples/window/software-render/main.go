// Software rendering to an ANativeWindow.
//
// Demonstrates the pattern for rendering pixels directly to an Android
// native window using the CPU (software rendering). This is useful for
// 2D applications, terminal emulators, or any scenario where GPU rendering
// is not needed.
//
// The window package exposes Window with methods for querying dimensions,
// setting buffer geometry, and posting rendered frames. The Lock step
// requires the low-level API because the idiomatic package does not yet
// wrap ANativeWindow_lock.
//
// Software rendering lifecycle:
//
//  1. Obtain     -- Get a *Window from the Activity's native surface.
//     The NDK delivers this via ANativeActivity callbacks
//     (onNativeWindowCreated). Wrap the pointer with
//     window.NewWindowFromPointer.
//
//  2. Query      -- Inspect the window: Width(), Height(), Format().
//     Width and Height return the current surface dimensions.
//     Format returns an int32 matching a window.Format constant.
//
//  3. Configure  -- Call SetBuffersGeometry(width, height, format)
//     to request a specific buffer size and pixel format.
//     Pass 0 for any dimension to keep the current value.
//     The compositor scales the buffer to the window size.
//
//  4. Lock       -- Call ANativeWindow_lock (the low-level layer) to lock the next
//     back-buffer for CPU access. Returns an
//     ANativeWindow_Buffer with pixel memory pointer, width,
//     height, stride, and format.
//     (Not yet exposed in the idiomatic layer.)
//
//  5. Draw       -- Write pixels into the locked buffer. Respect the
//     stride (pixels per row, may be wider than width due
//     to alignment). The pixel layout depends on the format.
//
//  6. Post       -- Call UnlockAndPost() to release the buffer and
//     submit it to the compositor for display.
//
//  7. Repeat     -- Go to step 4 for the next frame.
//
// Pixel formats:
//
//   Rgba8888  -- 4 bytes/pixel: [R, G, B, A]
//   Rgbx8888  -- 4 bytes/pixel: [R, G, B, x] (alpha ignored)
//   Rgb565    -- 2 bytes/pixel: 5 red, 6 green, 5 blue bits
//
// Prerequisites:
//   - Android device with a running Activity that provides a native window.
//   - The window handle is obtained from onNativeWindowCreated.
//
// Because obtaining a window requires an Activity, this example documents
// the complete software rendering pipeline and prints the API calls rather
// than invoking them directly.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/window"
)

func main() {
	fmt.Println("=== ANativeWindow software rendering ===")
	fmt.Println()

	// ── Step 1: Obtain a Window ─────────────────────────────────
	//
	// The Activity delivers a raw ANativeWindow* through the
	// onNativeWindowCreated callback:
	//
	//   func onNativeWindowCreated(rawPtr unsafe.Pointer) {
	//       win := window.NewWindowFromPointer(rawPtr)
	//       // start rendering
	//   }

	fmt.Println("Step 1: Obtain window from Activity")
	fmt.Println("  win := window.NewWindowFromPointer(rawPtr)")
	fmt.Println("  // rawPtr comes from onNativeWindowCreated callback")
	fmt.Println()

	// ── Step 2: Query window properties ─────────────────────────
	//
	// Before rendering, query the window dimensions and format:
	//
	//   w := win.Width()    // surface width in pixels
	//   h := win.Height()   // surface height in pixels
	//   f := win.Format()   // pixel format (int32)
	//
	// The format value matches one of the window.Format constants.

	fmt.Println("Step 2: Query window properties")
	fmt.Println("  win.Width()   -- surface width in pixels")
	fmt.Println("  win.Height()  -- surface height in pixels")
	fmt.Println("  win.Format()  -- pixel format (int32)")
	fmt.Println()

	// Print the available pixel format constants.
	formats := []struct {
		name   string
		value  window.Format
		bpp    int
		detail string
	}{
		{"Rgba8888", window.Rgba8888, 4, "8 bits per channel: R, G, B, A"},
		{"Rgbx8888", window.Rgbx8888, 4, "8 bits per channel: R, G, B, padding"},
		{"Rgb565", window.Rgb565, 2, "16-bit packed: 5R, 6G, 5B"},
	}

	fmt.Println("  Available pixel formats:")
	for _, f := range formats {
		fmt.Printf("    %-10s = %d  (%d BPP)  %s\n", f.name, int32(f.value), f.bpp, f.detail)
	}
	fmt.Println()

	// ── Step 3: Configure buffer geometry ───────────────────────
	//
	// Request a specific buffer size and format. This is optional
	// but useful for controlling the rendering resolution or
	// switching pixel formats.
	//
	//   // Set 720p resolution with RGBA8888 format.
	//   if err := win.SetBuffersGeometry(1280, 720, window.Rgba8888); err != nil {
	//       log.Fatalf("set geometry: %v", err)
	//   }
	//
	//   // Or keep current dimensions and only change the format:
	//   win.SetBuffersGeometry(0, 0, window.Rgb565)

	fmt.Println("Step 3: Configure buffer geometry")
	fmt.Println("  win.SetBuffersGeometry(1280, 720, window.Rgba8888)")
	fmt.Println("  // Pass 0 for width/height to keep current dimensions.")
	fmt.Println()

	// ── Step 4: Lock the back-buffer (idiomatic) ────────────────
	//
	// The idiomatic window package does not yet expose Lock. In a
	// real application, use the low-level API:
	//
	//   // use idiomatic package instead
	//
	//   var buf ANativeWindow_Buffer
	//   rc := ANativeWindow_lock(
	//       (*ANativeWindow)(win.Pointer()),
	//       &buf,
	//       nil,  // dirtyBounds: nil = full surface
	//   )
	//   if rc != 0 {
	//       log.Fatalf("lock: %d", rc)
	//   }
	//
	// The ANativeWindow_Buffer struct contains:
	//   - Width, Height: buffer dimensions in pixels
	//   - Stride: pixels per row (may be > Width due to alignment)
	//   - Format: pixel format
	//   - Bits: pointer to the pixel data

	fmt.Println("Step 4: Lock the back-buffer (requires the low-level layer)")
	fmt.Println("  ANativeWindow_lock(nativeWin, &buf, nil)")
	fmt.Println("  // buf.Bits  = pointer to pixel data")
	fmt.Println("  // buf.Stride = pixels per row (includes padding)")
	fmt.Println()

	// ── Step 5: Write pixels ────────────────────────────────────
	//
	// Write directly into the buffer memory. Always use the stride
	// (not width) to compute row offsets, as the buffer may have
	// alignment padding.
	//
	// Example: fill with a gradient (RGBA8888):
	//
	//   pixels := unsafe.Slice((*[4]byte)(buf.Bits), stride*height)
	//   for y := 0; y < height; y++ {
	//       for x := 0; x < width; x++ {
	//           off := y*stride + x
	//           pixels[off] = [4]byte{
	//               byte(x * 255 / width),   // R: left-to-right gradient
	//               byte(y * 255 / height),  // G: top-to-bottom gradient
	//               128,                     // B: constant
	//               255,                     // A: opaque
	//           }
	//       }
	//   }
	//
	// For Rgb565, pack each pixel into 2 bytes:
	//
	//   r5 := uint16(r >> 3)
	//   g6 := uint16(g >> 2)
	//   b5 := uint16(b >> 3)
	//   pixel565 := (r5 << 11) | (g6 << 5) | b5

	fmt.Println("Step 5: Write pixels into the locked buffer")
	fmt.Println("  // Use stride (not width) for row offset calculation.")
	fmt.Println("  // RGBA8888: 4 bytes per pixel [R, G, B, A]")
	fmt.Println("  // RGB565:   2 bytes per pixel, bit-packed")
	fmt.Println()

	// ── Step 6: Unlock and post ─────────────────────────────────
	//
	// UnlockAndPost atomically releases the CPU lock and submits the
	// buffer to the system compositor (SurfaceFlinger) for display.
	//
	//   if err := win.UnlockAndPost(); err != nil {
	//       log.Fatalf("unlock and post: %v", err)
	//   }

	fmt.Println("Step 6: Unlock and post")
	fmt.Println("  win.UnlockAndPost()")
	fmt.Println("  // Releases the buffer and submits it for display.")
	fmt.Println()

	// ── Step 7: Frame loop ──────────────────────────────────────
	//
	// Repeat steps 4-6 for each frame:
	//
	//   for running {
	//       // Lock
	//       ANativeWindow_lock(nativeWin, &buf, nil)
	//
	//       // Draw into buf.Bits
	//       renderFrame(&buf)
	//
	//       // Post
	//       win.UnlockAndPost()
	//   }
	//
	// For smooth animation, synchronize with the display refresh
	// rate using AChoreographer or eglSwapInterval. Without pacing,
	// the compositor may throttle or drop frames.

	fmt.Println("Step 7: Frame loop")
	fmt.Println("  for each frame: Lock -> Draw -> UnlockAndPost")
	fmt.Println("  Use AChoreographer for vsync-aligned frame pacing.")
	fmt.Println()

	// ── Window lifecycle notes ──────────────────────────────────
	//
	// The window handle is valid between onNativeWindowCreated and
	// onNativeWindowDestroyed callbacks. Applications must:
	//   - Stop rendering before onNativeWindowDestroyed returns.
	//   - Never call Lock or UnlockAndPost on a destroyed window.
	//   - Handle window size changes via onNativeWindowResized
	//     by re-querying Width/Height and adjusting rendering.

	fmt.Println("Lifecycle notes:")
	fmt.Println("  Window is valid between onNativeWindowCreated/Destroyed.")
	fmt.Println("  Stop rendering before onNativeWindowDestroyed returns.")
	fmt.Println("  Handle size changes via onNativeWindowResized.")
	fmt.Println()

	fmt.Println("software-render overview complete")
}
