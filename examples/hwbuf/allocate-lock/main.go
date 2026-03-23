// Hardware buffer allocation and CPU-lock lifecycle.
//
// Demonstrates how to allocate an AHardwareBuffer via the idiomatic hwbuf
// package, lock it for CPU access, write pixel data, unlock it, and release
// the buffer.
//
// AHardwareBuffer lifecycle:
//
//  1. Describe  -- Fill an hwbuf.Desc with width, height, layers, pixel
//     format, and usage flags. Usage flags control which hardware blocks
//     (CPU, GPU, compositor, video encoder) may access the buffer.
//
//  2. Allocate  -- Call hwbuf.Allocate with the descriptor. The driver
//     chooses the optimal memory domain.
//
//  3. Lock      -- Call buf.Lock with the desired CPU usage (read, write,
//     or both) and an optional dirty rectangle. The call returns a raw
//     pointer to the pixel data.
//
//  4. Access    -- Read or write pixels through the returned pointer.
//     Respect the stride returned by buf.Describe.
//
//  5. Unlock    -- Call buf.Unlock to release the CPU mapping and allow
//     the GPU or other consumers to access the buffer. The optional
//     fence parameter signals when the unlock completes asynchronously;
//     pass nil for synchronous.
//
//  6. Share     -- Optionally call buf.Acquire to increment the reference
//     count before passing the buffer to another component (another
//     thread, process, or API such as Vulkan or EGL).
//
//  7. Release   -- Call buf.Close to decrement the reference count. When
//     the count reaches zero the driver frees the memory.
//
// Prerequisites:
//   - Android device with API level 26+ (AHardwareBuffer was added in API 26).
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/AndroidGoLab/ndk/hwbuf"
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

	desc := hwbuf.Desc{
		Width:  1920,
		Height: 1080,
		Layers: 1,
		Format: uint32(hwbuf.R8g8b8a8Unorm),
		Usage:  uint64(hwbuf.CpuWriteOften | hwbuf.GpuSampledImage),
	}

	// ── Step 2: Allocate ────────────────────────────────────────
	//
	// hwbuf.Allocate creates a hardware buffer matching the
	// descriptor. The driver chooses the optimal memory domain.

	buf, err := hwbuf.Allocate(&desc)
	if err != nil {
		log.Fatalf("hwbuf.Allocate failed: %v", err)
	}
	defer func() {
		if err := buf.Close(); err != nil {
			log.Printf("close: %v", err)
		}
	}()

	fmt.Println("Allocation: success")

	// Read back the descriptor to get the stride chosen by the driver.
	actual := buf.Describe()
	fmt.Printf("  Actual stride: %d pixels\n", actual.Stride)
	fmt.Println()

	// ── Step 3: Lock for CPU access ─────────────────────────────
	//
	// Locking maps the buffer into the process address space for
	// CPU access. The usage parameter must be a subset of the
	// usage flags used at allocation time.
	// fence=-1 means no fence; rect=nil means full buffer.

	pixelPtr, err := buf.Lock(uint64(hwbuf.CpuWriteOften), -1, nil)
	if err != nil {
		log.Fatalf("buf.Lock failed: %v", err)
	}

	fmt.Println("Lock: success — buffer mapped for CPU write")

	// ── Step 4: Write pixels ────────────────────────────────────
	//
	// After locking, pixelPtr points to the first pixel. For
	// R8g8b8a8Unorm each pixel is 4 bytes: [R, G, B, A].
	// Use the actual stride (not width) to handle row padding.

	stride := int(actual.Stride) // pixels per row including padding
	pixels := unsafe.Slice((*byte)(pixelPtr), stride*1080*4)
	for y := 0; y < 1080; y++ {
		for x := 0; x < 1920; x++ {
			off := (y*stride + x) * 4
			pixels[off+0] = byte(x % 256) // R
			pixels[off+1] = byte(y % 256) // G
			pixels[off+2] = 0             // B
			pixels[off+3] = 255           // A
		}
	}

	fmt.Printf("  Wrote %dx%d RGBA pixels (stride=%d)\n", 1920, 1080, stride)
	fmt.Println()

	// ── Step 5: Unlock ──────────────────────────────────────────
	//
	// After writing pixels, unlock the buffer so the GPU or other
	// consumers can access it. Passing nil for the fence parameter
	// performs a synchronous unlock.

	if err := buf.Unlock(nil); err != nil {
		log.Fatalf("unlock: %v", err)
	}

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
	// the driver frees the buffer memory. Handled by deferred
	// buf.Close() above.

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
	fmt.Println("allocate-lock lifecycle complete")
}
