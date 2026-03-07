// Simulates c-for-go output for Android Storage Manager (android/storage_manager.h).
// This file is parsed at AST level only; it does not compile.
package storage

import "unsafe"

// Opaque handle types.
type AStorageManager C.AStorageManager

// --- StorageManager functions ---
func AStorageManager_new() *AStorageManager                                                             { return nil }
func AStorageManager_delete(manager *AStorageManager)                                                   {}
func AStorageManager_mountObb(manager *AStorageManager, filename *byte, key *byte, callback unsafe.Pointer, data unsafe.Pointer) {}
func AStorageManager_unmountObb(manager *AStorageManager, filename *byte, force int32, callback unsafe.Pointer, data unsafe.Pointer) {}
func AStorageManager_isObbMounted(manager *AStorageManager, filename *byte) int32                       { return 0 }
func AStorageManager_getMountedObbPath(manager *AStorageManager, filename *byte) *byte                  { return nil }

var _ = unsafe.Pointer(nil)
