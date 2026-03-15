// PersistableBundle API overview.
//
// Demonstrates the Android PersistableBundle NDK API, which provides a
// type-safe key-value container that can be serialized to a Parcel for
// IPC. PersistableBundles support boolean, int32, int64, float64, and
// string values, as well as nested PersistableBundles.
//
// Available since Android API level 35 (Android 15).
//
// Operations demonstrated:
//
//   - Create a new PersistableBundle
//   - Put and get scalar values (bool, int32, int64, float64, string)
//   - Query the bundle size
//   - Duplicate a bundle and compare with IsEqual
//   - Nest one bundle inside another
//   - Erase a key
//   - Clean up with Close
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/persistablebundle"
)

func main() {
	// ---------------------------------------------------------------
	// Create a new PersistableBundle.
	// ---------------------------------------------------------------
	pb := persistablebundle.NewPersistableBundle()
	defer pb.Close()

	// ---------------------------------------------------------------
	// Put scalar values of each supported type.
	// ---------------------------------------------------------------
	pb.PutBoolean("debug_enabled", true)
	pb.PutInt("retry_count", 3)
	pb.PutLong("session_id", 1234567890123)
	pb.PutDouble("threshold", 0.95)
	pb.PutString("app_name", "MyAndroidApp")

	fmt.Printf("bundle size after inserting 5 keys: %d\n", pb.Size())

	// ---------------------------------------------------------------
	// Get each value back.
	//
	// Get* methods return true if the key was found.
	// ---------------------------------------------------------------
	var boolVal bool
	if pb.GetBoolean("debug_enabled", &boolVal) {
		fmt.Printf("  debug_enabled = %v\n", boolVal)
	} else {
		log.Fatal("expected key 'debug_enabled' to exist")
	}

	var intVal int32
	if pb.GetInt("retry_count", &intVal) {
		fmt.Printf("  retry_count   = %d\n", intVal)
	} else {
		log.Fatal("expected key 'retry_count' to exist")
	}

	var longVal int64
	if pb.GetLong("session_id", &longVal) {
		fmt.Printf("  session_id    = %d\n", longVal)
	} else {
		log.Fatal("expected key 'session_id' to exist")
	}

	var doubleVal float64
	if pb.GetDouble("threshold", &doubleVal) {
		fmt.Printf("  threshold     = %.2f\n", doubleVal)
	} else {
		log.Fatal("expected key 'threshold' to exist")
	}

	// Note: GetString requires a C string allocator callback,
	// so it is only available via the capi package.

	// ---------------------------------------------------------------
	// Duplicate the bundle and verify equality.
	// ---------------------------------------------------------------
	pbCopy := pb.Dup()
	defer pbCopy.Close()

	fmt.Printf("original == copy: %v\n", pb.IsEqual(pbCopy))

	// ---------------------------------------------------------------
	// Nested bundles.
	//
	// PersistableBundles can contain other PersistableBundles,
	// enabling hierarchical configuration structures.
	// ---------------------------------------------------------------
	inner := persistablebundle.NewPersistableBundle()
	inner.PutInt("inner_value", 42)

	pb.PutPersistableBundle("nested", inner)
	// The inner bundle has been copied into pb, so we can close ours.
	inner.Close()

	fmt.Printf("bundle size after adding nested bundle: %d\n", pb.Size())

	retrieved, ok := pb.GetPersistableBundle("nested")
	if !ok {
		log.Fatal("expected key 'nested' to exist")
	}
	defer retrieved.Close()

	var innerVal int32
	if retrieved.GetInt("inner_value", &innerVal) {
		fmt.Printf("  nested.inner_value = %d\n", innerVal)
	} else {
		log.Fatal("expected key 'inner_value' in nested bundle")
	}

	// ---------------------------------------------------------------
	// Erase a key.
	//
	// Erase returns the number of entries removed (0 or 1).
	// ---------------------------------------------------------------
	erased := pb.Erase("retry_count")
	fmt.Printf("erased 'retry_count': count=%d\n", erased)
	fmt.Printf("bundle size after erase: %d\n", pb.Size())

	// Verify the key is gone.
	if pb.GetInt("retry_count", &intVal) {
		log.Fatal("expected key 'retry_count' to be erased")
	}
	fmt.Println("confirmed: 'retry_count' no longer exists")

	fmt.Println("overview complete")
}
