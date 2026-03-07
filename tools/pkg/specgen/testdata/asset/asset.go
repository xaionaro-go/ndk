// Simulates c-for-go output for Android Asset Manager.
// This file is parsed at AST level only; it does not compile.
package asset

import "unsafe"

// Opaque handle types.
type AAssetManager C.AAssetManager
type AAssetDir C.AAssetDir
type AAsset C.AAsset

// Integer typedefs.
type Asset_mode_t int32

// Mode enum.
const (
	AASSET_MODE_UNKNOWN   Asset_mode_t = 0
	AASSET_MODE_RANDOM    Asset_mode_t = 1
	AASSET_MODE_STREAMING Asset_mode_t = 2
	AASSET_MODE_BUFFER    Asset_mode_t = 3
)

// --- AssetManager functions ---
func AAssetManager_openDir(mgr *AAssetManager, dirName *byte) *AAssetDir { return nil }
func AAssetManager_open(mgr *AAssetManager, filename *byte, mode int32) *AAsset { return nil }

// --- AssetDir functions ---
func AAssetDir_getNextFileName(assetDir *AAssetDir) *byte { return nil }
func AAssetDir_rewind(assetDir *AAssetDir)                {}
func AAssetDir_close(assetDir *AAssetDir)                 {}

// --- Asset functions ---
func AAsset_read(asset *AAsset, buf unsafe.Pointer, count int64) int32 { return 0 }
func AAsset_seek(asset *AAsset, offset int64, whence int32) int64      { return 0 }
func AAsset_close(asset *AAsset)                                       {}
func AAsset_getBuffer(asset *AAsset) unsafe.Pointer                    { return nil }
func AAsset_getLength(asset *AAsset) int64                             { return 0 }
func AAsset_getLength64(asset *AAsset) int64                           { return 0 }
func AAsset_getRemainingLength(asset *AAsset) int64                    { return 0 }
func AAsset_getRemainingLength64(asset *AAsset) int64                  { return 0 }
func AAsset_isAllocated(asset *AAsset) int32                           { return 0 }

var _ = unsafe.Pointer(nil)
