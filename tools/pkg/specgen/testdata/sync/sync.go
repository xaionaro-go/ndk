// Simulates c-for-go output for Android Sync.
// This file is parsed at AST level only; it does not compile.
package sync

import "unsafe"

// --- Sync functions ---
func sync_merge(name *byte, fd1 int32, fd2 int32) int32 { return 0 }
func sync_file_info_free(info unsafe.Pointer)            {}

var _ = unsafe.Pointer(nil)
