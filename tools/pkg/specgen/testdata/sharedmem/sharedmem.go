// Simulates c-for-go output for Android SharedMemory.
// This file is parsed at AST level only; it does not compile.
package sharedmem

// --- SharedMemory functions ---
func ASharedMemory_create(name *byte, size int32) int32   { return 0 }
func ASharedMemory_getSize(fd int32) int32                { return 0 }
func ASharedMemory_setProt(fd int32, prot int32) int32    { return 0 }
