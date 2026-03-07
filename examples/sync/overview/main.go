// Example: sync package overview
//
// The ndk/sync package provides Go bindings for Android's native
// sync framework (sync fences / sync file descriptors).
//
// Android sync fences are kernel-level synchronization primitives
// used to coordinate GPU and display operations. A sync fence
// represents a point in time on a GPU command stream — it signals
// when all commands submitted before it have completed.
//
// Typical usage in the graphics pipeline:
//  1. GPU driver creates a fence when submitting work
//  2. Fence fd is passed to SurfaceFlinger / HWC
//  3. Consumer waits on the fence before reading the buffer
//  4. This enables pipelined rendering without CPU stalls
//
// The sync fence API operates on file descriptors:
//   - sync_merge(name, fd1, fd2)  — merge two fences into one
//   - sync_wait(fd, timeout_ms)   — block until fence signals
//   - sync_file_info(fd)          — query fence metadata
//
// Sync fences are typically obtained from EGL (eglDupNativeFenceFD),
// Vulkan (vkGetFenceFdKHR), or AHardwareBuffer operations.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	_ "github.com/xaionaro-go/ndk/sync"
)

func main() {
	fmt.Println("github.com/xaionaro-go/ndk/sync — Android Native Sync Fence API")
	fmt.Println()
	fmt.Println("Sync fences coordinate GPU work completion with buffer consumers.")
	fmt.Println()
	fmt.Println("Key NDK functions:")
	fmt.Println("  sync_wait(fd, timeout)     — wait for fence to signal")
	fmt.Println("  sync_merge(name, fd1, fd2) — combine two fences")
	fmt.Println("  sync_file_info(fd)         — query fence details")
	fmt.Println()
	fmt.Println("Fences are obtained from:")
	fmt.Println("  - EGL:    eglDupNativeFenceFD()")
	fmt.Println("  - Vulkan: vkGetFenceFdKHR()")
	fmt.Println("  - HWB:    AHardwareBuffer_unlock() fence out-param")
}
