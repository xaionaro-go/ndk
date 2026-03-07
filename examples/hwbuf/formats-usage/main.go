// Hardware buffer formats and usage flags overview.
//
// Enumerates every AHardwareBuffer pixel format with its String()
// representation, then shows how to combine Usage flags with bitwise
// OR to describe the intended access pattern for a buffer.
//
// AHardwareBuffer is Android's cross-process, cross-API buffer primitive.
// It enables zero-copy sharing of pixel and blob data between the CPU,
// GPU (OpenGL ES / Vulkan), hardware composer, camera, and codec
// pipelines.  A single buffer can be written by the GPU and read by the
// CPU (or vice-versa) without an intermediate copy, provided the correct
// Usage flags are set at allocation time.
//
// This program exercises only the type system (Format, Usage,
// AHardwareBuffer_Desc); it does not call Allocate, so it compiles and
// runs on any host for demonstration purposes.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/hwbuf"
)

func main() {
	// ── Pixel formats ────────────────────────────────────────────
	//
	// Each Format constant maps to the NDK AHARDWAREBUFFER_FORMAT_*
	// value.  The String() method returns a human-readable name.

	formats := []struct {
		f    hwbuf.Format
		bits string // short description of the layout
	}{
		{hwbuf.R8g8b8a8Unorm, "32-bit RGBA, 8 bits per channel"},
		{hwbuf.R8g8b8x8Unorm, "32-bit RGBx, alpha channel ignored"},
		{hwbuf.R8g8b8Unorm, "24-bit RGB, no alpha"},
		{hwbuf.R5g6b5Unorm, "16-bit RGB (5-6-5)"},
		{hwbuf.R16g16b16a16Float, "64-bit RGBA, 16-bit float per channel (HDR)"},
		{hwbuf.Blob, "opaque byte buffer (non-image data)"},
	}

	fmt.Println("AHardwareBuffer pixel formats")
	fmt.Println("─────────────────────────────────────────────────────")
	fmt.Printf("  %-25s  %5s  %s\n", "Name", "Value", "Description")
	for _, entry := range formats {
		fmt.Printf("  %-25s  %5d  %s\n", entry.f, int32(entry.f), entry.bits)
	}

	// ── Usage flags ──────────────────────────────────────────────
	//
	// Usage constants are a bitmask.  Combine them with | to declare
	// how the buffer will be accessed.  The driver uses these hints to
	// place the allocation in the most efficient memory domain.

	fmt.Println()
	fmt.Println("Individual usage flags")
	fmt.Println("─────────────────────────────────────────────────────")

	usages := []struct {
		u    hwbuf.Usage
		name string
	}{
		{hwbuf.CpuReadOften, "CpuReadOften"},
		{hwbuf.CpuWriteOften, "CpuWriteOften"},
		{hwbuf.GpuSampledImage, "GpuSampledImage"},
		{hwbuf.GpuColorOutput, "GpuColorOutput"},
	}

	for _, entry := range usages {
		fmt.Printf("  %-20s = 0x%04X (%d)\n", entry.name, int32(entry.u), int32(entry.u))
	}

	// ── Composing usage masks ────────────────────────────────────
	//
	// Real-world buffers typically need more than one flag.  A GPU
	// render target that the CPU will read back for screenshots, for
	// example, combines GpuColorOutput | CpuReadOften.

	fmt.Println()
	fmt.Println("Combined usage examples")
	fmt.Println("─────────────────────────────────────────────────────")

	// GPU renders into the buffer, CPU reads the pixels back.
	gpuRenderCpuReadback := hwbuf.GpuColorOutput | hwbuf.CpuReadOften
	fmt.Printf("  GPU render + CPU readback : 0x%04X\n", int32(gpuRenderCpuReadback))

	// CPU writes pixel data, GPU samples it as a texture.
	cpuUploadGpuTexture := hwbuf.CpuWriteOften | hwbuf.GpuSampledImage
	fmt.Printf("  CPU upload + GPU texture  : 0x%04X\n", int32(cpuUploadGpuTexture))

	// Full bidirectional: CPU reads and writes, GPU renders and samples.
	bidirectional := hwbuf.CpuReadOften | hwbuf.CpuWriteOften |
		hwbuf.GpuSampledImage | hwbuf.GpuColorOutput
	fmt.Printf("  Full bidirectional        : 0x%04X\n", int32(bidirectional))

	// ── Descriptor sketch ────────────────────────────────────────
	//
	// On a real device the descriptor would be passed to Allocate:
	//
	//   desc := hwbuf.AHardwareBuffer_Desc{
	//       Width:  1920,
	//       Height: 1080,
	//       Layers: 1,
	//       Format: hwbuf.R8g8b8a8Unorm,
	//       Usage:  gpuRenderCpuReadback,
	//   }
	//   buf, err := hwbuf.Allocate(&desc)
	//
	// After allocation the buffer can be:
	//   - Locked for CPU access   (Lock / Unlock)
	//   - Bound as a GL texture   (via EGLClientBuffer)
	//   - Imported into Vulkan    (VkImportAndroidHardwareBufferInfoANDROID)
	//   - Sent across processes   (Binder / Unix sockets)
	//
	// All without copying the underlying pixel data.

	fmt.Println()
	fmt.Println("done")
}
