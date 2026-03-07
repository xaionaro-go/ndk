// Example: sharedmem package overview
//
// The ndk/sharedmem package provides Go bindings for Android's
// ASharedMemory API (available since API level 26).
//
// ASharedMemory enables creating and managing anonymous shared memory
// regions that can be shared between processes via file descriptors.
// This is useful for zero-copy IPC, large data transfers between
// services, and sharing buffers with the Android framework.
//
// Typical workflow:
//  1. Create a shared memory region: ASharedMemory_create(name, size)
//  2. Get the file descriptor for sharing
//  3. Map into process address space with mmap()
//  4. Share the fd with another process (via Binder/socket)
//  5. The other process mmaps the same fd
//  6. Both processes read/write the shared region
//  7. Set protection with ASharedMemory_setProt() to make read-only
//
// The underlying ASharedMemory functions operate on file descriptors
// (int), not opaque handles.
//
// This program must run on an Android device with API level 26+.
package main

import (
	"fmt"

	_ "github.com/xaionaro-go/ndk/sharedmem"
)

func main() {
	fmt.Println("github.com/xaionaro-go/ndk/sharedmem — Android ASharedMemory API")
	fmt.Println()
	fmt.Println("ASharedMemory provides anonymous shared memory for IPC.")
	fmt.Println()
	fmt.Println("Key NDK functions:")
	fmt.Println("  ASharedMemory_create(name, size)    — create a region, returns fd")
	fmt.Println("  ASharedMemory_getSize(fd)            — query region size")
	fmt.Println("  ASharedMemory_setProt(fd, prot)      — restrict access (e.g. read-only)")
	fmt.Println()
	fmt.Println("The fd can be passed to another process via Binder or Unix socket,")
	fmt.Println("then both processes mmap() it for zero-copy shared access.")
}
