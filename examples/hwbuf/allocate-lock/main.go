// Hardware buffer allocation and CPU-lock lifecycle.
//
// Demonstrates how to allocate an AHardwareBuffer via the capi layer,
// wrap it in the idiomatic hwbuf.Buffer, lock it for CPU access, write
// pixel data, unlock it, and release the buffer.
//
// The idiomatic hwbuf package exposes the Buffer type with Acquire,
// Unlock, and Close methods. Allocation and locking require the capi
// layer because the idiomatic package does not yet wrap
// AHardwareBuffer_allocate or AHardwareBuffer_lock. This example
// documents the full pattern and shows how to bridge the gap.
//
// AHardwareBuffer lifecycle:
//
//  1. Describe  -- Fill an AHardwareBuffer_Desc with width, height,
//     layers, pixel format, and usage flags. Usage flags
//     control which hardware blocks (CPU, GPU, compositor,
//     video encoder) may access the buffer.
//
//  2. Allocate  -- Call AHardwareBuffer_allocate with the descriptor.
//     The driver chooses the optimal memory domain.
//     (Not yet exposed in the idiomatic layer; use capi.)
//
//  3. Lock      -- Call AHardwareBuffer_lock with the desired CPU
//     usage (read, write, or both) and an optional dirty
//     rectangle. The call returns a raw pointer to the
//     pixel data. (Not yet exposed in the idiomatic layer.)
//
//  4. Access    -- Read or write pixels through the returned pointer.
//     Respect the stride returned in the descriptor or via
//     AHardwareBuffer_lockAndGetInfo.
//
//  5. Unlock    -- Call Buffer.Unlock to release the CPU mapping and
//     allow the GPU or other consumers to access the buffer.
//     The optional fence parameter signals when the unlock
//     completes asynchronously; pass nil for synchronous.
//
//  6. Share     -- Optionally call Buffer.Acquire to increment the
//     reference count before passing the buffer to another
//     component (another thread, process, or API such as
//     Vulkan or EGL).
//
//  7. Release   -- Call Buffer.Close to decrement the reference count.
//     When the count reaches zero the driver frees the memory.
//
// This example exercises only the type system and documents the
// allocation pattern. Because AHardwareBuffer_allocate is not
// available in the idiomatic layer, the allocation step is shown as
// pseudocode in comments. The Acquire/Unlock/Close lifecycle is
// demonstrated with a hypothetical buffer pointer.
//
// Prerequisites:
//   - Android device with API level 26+ (AHardwareBuffer was added in API 26).
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/hwbuf"
)

func main() {
	fmt.Println("=== AHardwareBuffer allocate-lock lifecycle ===")
	fmt.Println()

	// ── Step 1: Describe the buffer ─────────────────────────────
	//
	// The descriptor specifies the buffer dimensions, format, and
	// usage flags. Usage flags are a bitmask that tells the driver
	// which hardware blocks will access the buffer.

	format := hwbuf.R8g8b8a8Unorm
	usage := hwbuf.CpuWriteOften | hwbuf.GpuSampledImage

	fmt.Println("Buffer descriptor:")
	fmt.Printf("  Width:   %d\n", 1920)
	fmt.Printf("  Height:  %d\n", 1080)
	fmt.Printf("  Layers:  %d\n", 1)
	fmt.Printf("  Format:  %s (value %d)\n", format, int32(format))
	fmt.Printf("  Usage:   CpuWriteOften | GpuSampledImage = 0x%04X\n", uint64(usage))
	fmt.Println()

	// ── Step 2: Allocate (capi only) ────────────────────────────
	//
	// The idiomatic layer does not yet expose Allocate. In a real
	// application you would call the capi function:
	//
	//   import capi "github.com/xaionaro-go/ndk/capi/hardwarebuffer"
	//
	//   desc := capi.AHardwareBuffer_Desc{
	//       Width:  1920,
	//       Height: 1080,
	//       Layers: 1,
	//       Format: uint32(hwbuf.R8g8b8a8Unorm),
	//       Usage:  uint64(hwbuf.CpuWriteOften | hwbuf.GpuSampledImage),
	//   }
	//   var rawBuf *capi.AHardwareBuffer
	//   rc := capi.AHardwareBuffer_allocate(&desc, &rawBuf)
	//   if rc != 0 {
	//       log.Fatalf("allocate failed: %d", rc)
	//   }
	//   buf := hwbuf.NewBufferFromPointer(unsafe.Pointer(rawBuf))

	fmt.Println("Allocation: requires capi.AHardwareBuffer_allocate")
	fmt.Println("  (see comments in source for the full pattern)")
	fmt.Println()

	// ── Step 3: Lock for CPU access (capi only) ─────────────────
	//
	// Locking maps the buffer into the process address space for
	// CPU access. The usage parameter must be a subset of the
	// usage flags used at allocation time.
	//
	//   var pixelPtr unsafe.Pointer
	//   rc = capi.AHardwareBuffer_lock(
	//       rawBuf,
	//       uint64(hwbuf.CpuWriteOften), // lock for writing
	//       -1,                          // fence: -1 = no fence
	//       nil,                         // rect: nil = full buffer
	//       &pixelPtr,
	//   )
	//   if rc != 0 {
	//       log.Fatalf("lock failed: %d", rc)
	//   }
	//
	// After locking, pixelPtr points to the first pixel. For
	// R8g8b8a8Unorm each pixel is 4 bytes: [R, G, B, A].
	//
	// ── Step 4: Write pixels ────────────────────────────────────
	//
	//   pixels := unsafe.Slice((*byte)(pixelPtr), 1920*1080*4)
	//   for y := 0; y < 1080; y++ {
	//       for x := 0; x < 1920; x++ {
	//           off := (y*1920 + x) * 4
	//           pixels[off+0] = byte(x % 256) // R
	//           pixels[off+1] = byte(y % 256) // G
	//           pixels[off+2] = 0             // B
	//           pixels[off+3] = 255           // A
	//       }
	//   }

	fmt.Println("Lock + write: requires capi.AHardwareBuffer_lock")
	fmt.Println("  (see comments in source for the full pattern)")
	fmt.Println()

	// ── Step 5: Unlock ──────────────────────────────────────────
	//
	// After writing pixels, unlock the buffer so the GPU or other
	// consumers can access it. Passing nil for the fence parameter
	// performs a synchronous unlock.
	//
	//   if err := buf.Unlock(nil); err != nil {
	//       log.Fatalf("unlock: %v", err)
	//   }

	fmt.Println("Unlock: buf.Unlock(nil) -- synchronous release of CPU mapping")
	fmt.Println()

	// ── Step 6: Acquire (reference counting) ────────────────────
	//
	// Acquire increments the buffer's reference count. This is
	// needed when passing the buffer to another component that
	// will independently call Close (decrement).
	//
	//   buf.Acquire()  // ref count: 1 -> 2
	//   // pass buf to another goroutine / API
	//   // that goroutine calls buf.Close() when done

	fmt.Println("Acquire: buf.Acquire() -- increment reference count")
	fmt.Println("  Use this before sharing the buffer across components.")
	fmt.Println()

	// ── Step 7: Release ─────────────────────────────────────────
	//
	// Close decrements the reference count. When it reaches zero
	// the driver frees the buffer memory.
	//
	//   if err := buf.Close(); err != nil {
	//       log.Fatalf("close: %v", err)
	//   }

	fmt.Println("Close: buf.Close() -- decrement reference count and free if zero")
	fmt.Println()

	// ── Usage flag combinations ─────────────────────────────────
	//
	// Common real-world usage flag combinations:

	combos := []struct {
		name  string
		usage hwbuf.Usage
	}{
		{"CPU write + GPU texture", hwbuf.CpuWriteOften | hwbuf.GpuSampledImage},
		{"GPU render + CPU readback", hwbuf.GpuColorOutput | hwbuf.CpuReadOften},
		{"GPU render + compositor", hwbuf.GpuColorOutput | hwbuf.ComposerOverlay},
		{"CPU write + video encode", hwbuf.CpuWriteOften | hwbuf.VideoEncode},
		{"GPU render target", hwbuf.GpuFramebuffer | hwbuf.GpuSampledImage},
	}

	fmt.Println("Common usage flag combinations:")
	for _, c := range combos {
		fmt.Printf("  %-30s = 0x%08X\n", c.name, uint64(c.usage))
	}
	fmt.Println()

	// ── Supported pixel formats for CPU access ──────────────────

	cpuFormats := []struct {
		f   hwbuf.Format
		bpp int
	}{
		{hwbuf.R8g8b8a8Unorm, 4},
		{hwbuf.R8g8b8x8Unorm, 4},
		{hwbuf.R8g8b8Unorm, 3},
		{hwbuf.R5g6b5Unorm, 2},
		{hwbuf.R16g16b16a16Float, 8},
		{hwbuf.R8Unorm, 1},
	}

	fmt.Println("Pixel formats suitable for CPU access:")
	fmt.Printf("  %-25s  %s\n", "Format", "BPP")
	for _, entry := range cpuFormats {
		fmt.Printf("  %-25s  %d\n", entry.f, entry.bpp)
	}

	fmt.Println()
	fmt.Println("allocate-lock lifecycle overview complete")
}
