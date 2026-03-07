// Example: AAssetManager API overview and directory iteration.
//
// Demonstrates the full ndk/asset API surface: opening assets by path,
// reading content, querying length and allocation state, seeking within an
// asset, and iterating directories.
//
// AAssetManager lifecycle on Android:
//
//  1. Obtain the AAssetManager pointer from the Activity's ANativeActivity
//     struct (activity.AssetManager in ndk/activity). The manager is owned
//     by the framework and must not be freed by application code.
//  2. Open an individual asset with Manager.Open(filename, mode). The
//     filename is relative to the APK's "assets/" directory (e.g.
//     "textures/wood.png"). The mode controls how the NDK buffers the data.
//  3. Read the asset with Asset.Read, or access it directly via
//     Asset.Buffer (only valid when opened with Buffer mode). Query the
//     size with Asset.Length / Asset.Length64 and the unread portion with
//     Asset.RemainingLength / Asset.RemainingLength64. Reposition with
//     Asset.Seek using standard SEEK_SET/SEEK_CUR/SEEK_END whence values.
//  4. Close the asset with Asset.Close when finished.
//
// Directory iteration:
//
//  1. Open a directory with Manager.OpenDir(dirName). Pass an empty string
//     to list the root of the assets/ tree, or a subdirectory path.
//  2. Call Dir.NextFileName repeatedly. Each call returns the next filename
//     as a string, or "" when the listing is exhausted.
//  3. Call Dir.Rewind to reset the iterator back to the first entry.
//  4. Close with Dir.Close when finished.
//
// This program prints all Mode constants and documents the usage patterns
// in executable form. It cannot run without a real AAssetManager obtained
// from an Android Activity.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/xaionaro-go/ndk/asset"
)

// modeName returns a human-readable description for an asset open mode.
func modeName(m asset.Mode) string {
	switch m {
	case asset.Unknown:
		return "Unknown"
	case asset.Random:
		return "Random"
	case asset.Streaming:
		return "Streaming"
	case asset.Buffer:
		return "Buffer"
	default:
		return fmt.Sprintf("Mode(%d)", int32(m))
	}
}

// modeDescription returns a short explanation of when to use each mode.
func modeDescription(m asset.Mode) string {
	switch m {
	case asset.Unknown:
		return "no specific access pattern; NDK chooses defaults"
	case asset.Random:
		return "random access via Seek; NDK may memory-map the asset"
	case asset.Streaming:
		return "sequential forward reads; NDK may use a smaller buffer"
	case asset.Buffer:
		return "entire asset loaded into memory; Buffer() returns a direct pointer"
	default:
		return "unknown mode"
	}
}

func main() {
	// -- Mode constants ----------------------------------------------------
	// Mode controls how the NDK accesses the underlying APK data for an
	// asset. Choose the mode that matches your read pattern for best
	// performance.
	fmt.Println("Asset open modes:")
	fmt.Println()
	fmt.Printf("  %-12s  %5s  %s\n", "Name", "Value", "Description")
	fmt.Printf("  %-12s  %5s  %s\n", "------------", "-----", "-----------")

	for _, m := range []asset.Mode{
		asset.Unknown,
		asset.Random,
		asset.Streaming,
		asset.Buffer,
	} {
		fmt.Printf("  %-12s  %5d  %s\n", modeName(m), int32(m), modeDescription(m))
	}
	fmt.Println()

	// -- Asset open/read/close pattern ------------------------------------
	// The code below shows the typical sequence for reading an asset file.
	// It requires a valid *asset.Manager, which comes from the Activity:
	//
	//   mgr := activity.AssetManager  // from ANativeActivity
	//
	// Since we cannot obtain a real manager in a standalone example, the
	// remainder is guarded by a nil check and serves as a code reference.

	var mgr *asset.Manager // would come from the Activity in a real app

	if mgr != nil {
		// Open a text file from the APK's assets/ directory.
		a := mgr.Open("data/config.json", int32(asset.Streaming))
		defer a.Close()

		// Query total size. Length returns off_t (32-bit on older ABIs),
		// Length64 always returns a 64-bit value.
		totalLen := a.Length64()
		fmt.Printf("  Asset length: %d bytes\n", totalLen)

		// Check how much data remains unread from the current position.
		remaining := a.RemainingLength64()
		fmt.Printf("  Remaining:    %d bytes\n", remaining)

		// IsAllocated reports whether the asset data lives in a dedicated
		// memory allocation (1) or is memory-mapped from the APK (0).
		fmt.Printf("  Allocated:    %d\n", a.IsAllocated())

		// Read the entire asset into a byte slice.
		buf := make([]byte, totalLen)
		err := a.Read(unsafe.Pointer(&buf[0]), uint64(totalLen))
		if err != nil {
			fmt.Printf("  Read error: %v\n", err)
		} else {
			fmt.Printf("  Read %d bytes successfully\n", totalLen)
		}

		// Seek back to the beginning using standard io.SeekStart (0).
		// Seek returns the new position or -1 on error.
		newPos := a.Seek(0, int32(io.SeekStart))
		fmt.Printf("  Seek to start: position = %d\n", newPos)

		// With Buffer mode, the entire asset is accessible through a
		// direct pointer without calling Read:
		//
		//   a := mgr.Open("data/config.json", int32(asset.Buffer))
		//   ptr := a.Buffer()  // unsafe.Pointer to the raw data
		//   data := unsafe.Slice((*byte)(ptr), a.Length64())
	}

	// -- Directory iteration pattern --------------------------------------
	// OpenDir lists files inside an assets/ subdirectory.

	if mgr != nil {
		// Open the root assets directory (empty string = root).
		dir := mgr.OpenDir("")
		defer dir.Close()

		fmt.Println("  Files in assets root:")
		for {
			name := dir.NextFileName()
			if name == "" {
				break
			}
			fmt.Printf("    %s\n", name)
		}

		// Rewind resets the iterator so you can walk the listing again.
		dir.Rewind()
		fmt.Println("  After Rewind, first file again:")
		if name := dir.NextFileName(); name != "" {
			fmt.Printf("    %s\n", name)
		}
	}

	// -- Subdirectory iteration -------------------------------------------
	// Pass a path relative to assets/ to list a subdirectory.

	if mgr != nil {
		dir := mgr.OpenDir("textures")
		defer dir.Close()

		fmt.Println("  Files in assets/textures/:")
		for {
			name := dir.NextFileName()
			if name == "" {
				break
			}
			fmt.Printf("    %s\n", name)
		}
	}

	fmt.Println()
	fmt.Println("Asset lifecycle and directory iteration patterns printed above in source.")
}
