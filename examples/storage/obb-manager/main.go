// OBB storage manager lifecycle example.
//
// Demonstrates how to create an Android storage manager, check whether an
// OBB (Opaque Binary Blob) file is mounted, and query its mounted path.
//
// OBB files are large binary assets (up to 2 GB each) that Android games and
// apps use to store resources that exceed the APK size limit. The system
// mounts an OBB as a virtual filesystem at a path under
// /storage/emulated/0/Android/obb/<package>/. Once mounted, the application
// reads assets from the mounted path like ordinary files.
//
// Typical OBB workflow:
//  1. Create a storage manager.
//  2. Check if the OBB is already mounted (IsObbMounted).
//  3. If not mounted, mount it (MountObb with an async callback).
//  4. Obtain the mounted path (MountedObbPath) to read assets.
//  5. When done, unmount the OBB (UnmountObb with an async callback).
//  6. Close the manager to release native resources.
//
// This example covers steps 1, 2, 4, and 6. Steps 3 and 5 (asynchronous
// mount/unmount) require a running ALooper on the calling thread, which is
// beyond the scope of this minimal example.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/storage"
)

func main() {
	// Create a storage manager. This allocates the underlying
	// AStorageManager and is the entry point for all OBB operations.
	mgr := storage.NewManager()
	if mgr == nil {
		log.Fatal("storage.NewManager returned nil")
	}
	defer func() {
		if err := mgr.Close(); err != nil {
			log.Fatalf("manager close: %v", err)
		}
		log.Println("storage manager closed")
	}()
	log.Println("storage manager created")

	// The OBB filename is an absolute path on the device filesystem.
	obbPath := "/storage/emulated/0/Android/obb/com.example.game/main.1.com.example.game.obb"

	// Check whether this OBB file is currently mounted.
	// IsObbMounted returns 1 if mounted, 0 if not.
	mounted := mgr.IsObbMounted(obbPath)
	if mounted != 0 {
		fmt.Println("OBB is mounted")
	} else {
		fmt.Println("OBB is not mounted")
	}

	// Query the mounted path. When the OBB is mounted, the NDK returns a
	// path such as "/storage/emulated/0/Android/obb/com.example.game/".
	// When it is not mounted the return value is typically an empty string.
	mountedPath := mgr.MountedObbPath(obbPath)
	if mountedPath != "" {
		fmt.Printf("Mounted path: %s\n", mountedPath)
	} else {
		fmt.Println("Mounted path is empty (OBB not mounted)")
	}

	// In a full application the next steps would be:
	//
	//   // Mount the OBB (asynchronous, requires an ALooper on this thread).
	//   // The key is "" for unencrypted OBBs.
	//   mgr.MountObb(obbPath, "", callback, data)
	//
	//   // After the callback fires confirming success, read assets from
	//   // the mounted path returned by MountedObbPath.
	//
	//   // Unmount when finished. Pass force=1 to force-unmount even if
	//   // files are still open.
	//   mgr.UnmountObb(obbPath, 0, callback, data)

	// Demonstrate idempotent close: the deferred Close above will be a
	// harmless no-op after this explicit close.
	if err := mgr.Close(); err != nil {
		log.Fatalf("close: %v", err)
	}
	log.Println("storage manager closed (first close)")
}
