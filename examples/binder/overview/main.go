// Binder IPC API overview.
//
// Documents the Android Binder IPC architecture and the ndk/binder type
// hierarchy. Binder is Android's primary inter-process communication
// mechanism -- virtually all system service calls (SurfaceFlinger, Activity
// Manager, Package Manager, etc.) flow through Binder transactions.
//
// Architecture:
//
//	Client process                          Server process
//	  |                                       |
//	  |-- AIBinder (proxy) ---[kernel]--> AIBinder (local)
//	  |       |                                |
//	  |   prepareTransaction              onTransact callback
//	  |       |                                |
//	  |   AParcel (write args) --------> AParcel (read args)
//	  |                                        |
//	  |   AParcel (read reply) <-------- AParcel (write reply)
//	  |       |                                |
//	  |   AStatus (check result)          AStatus (return result)
//
// Type hierarchy in ndk/binder:
//
//	Class   - Defines a binder interface (its descriptor string and the
//	          three callbacks: onCreate, onDestroy, onTransact).
//	          Created via AIBinder_Class_define.
//	          Has no Close method -- the class definition is static
//	          and lives for the process lifetime.
//
//	Binder  - A handle to a binder object (local or remote proxy).
//	          Created via AIBinder_new (local) or obtained from the
//	          Android service manager (remote).
//	          Close() decrements the strong reference count.
//
//	Parcel  - A container for serialized transaction data. Supports
//	          writing/reading int32 values and strings.
//	          Close() deletes the parcel.
//
//	Status  - Represents the outcome of a binder transaction.
//	          Close() deletes the status object.
//
//	Error   - An error type wrapping NDK binder status codes.
//	          Implements the error interface. Constants like
//	          ErrDeadObject and ErrPermissionDenied map to the
//	          underlying STATUS_* codes from the NDK.
//
// Transaction flow:
//
//	// 1. Define a binder class with callbacks.
//	//    (AIBinder_Class_define — not yet in high-level API)
//
//	// 2. Create a local binder instance.
//	//    (AIBinder_new — not yet in high-level API)
//
//	// 3. Prepare a transaction — allocates an input parcel.
//	//    (AIBinder_prepareTransaction — not yet in high-level API)
//
//	// 4. Marshal arguments into the input parcel.
//	//    parcel.WriteInt32(42)
//
//	// 5. Execute the transaction — sends input, receives output.
//	//    (AIBinder_transact — not yet in high-level API)
//
//	// 6. Unmarshal the reply from the output parcel.
//	//    parcel.ReadInt32(&result)
//
//	// 7. Clean up.
//	//    parcel.Close()
//	//    binder.Close()
//
// On the server side, the onTransact callback receives the transaction code
// and input parcel, reads the arguments, performs the operation, writes the
// reply into the output parcel, and returns an AStatus.
//
// Prerequisites:
//   - Android device or emulator with API level 29+ (Binder NDK API).
//   - The binder interface must be registered with the Android service
//     manager for cross-process discovery. Local (in-process) binder
//     objects can be created directly with AIBinder_new.
//   - SELinux policy must permit the binder transactions between the
//     client and server security contexts.
package main

import (
	"fmt"
	"log"

	"github.com/AndroidGoLab/ndk/binder"
)

func main() {
	// ---------------------------------------------------------------
	// Error codes
	//
	// The binder package defines Error constants that map to the
	// underlying NDK STATUS_* codes. These are returned by binder
	// operations and converted to Go errors by the binder package.
	//
	// Use errors.Is to check for specific binder failures:
	//   if errors.Is(err, binder.ErrDeadObject) { reconnect... }
	// ---------------------------------------------------------------
	errors := []struct {
		name string
		err  binder.Error
	}{
		{"ErrUnknownError", binder.ErrUnknownError},
		{"ErrNoMemory", binder.ErrNoMemory},
		{"ErrInvalidOperation", binder.ErrInvalidOperation},
		{"ErrBadValue", binder.ErrBadValue},
		{"ErrDeadObject", binder.ErrDeadObject},
		{"ErrPermissionDenied", binder.ErrPermissionDenied},
		{"ErrNameNotFound", binder.ErrNameNotFound},
		{"ErrWouldBlock", binder.ErrWouldBlock},
		{"ErrFdsNotAllowed", binder.ErrFdsNotAllowed},
	}

	fmt.Println("binder.Error constants:")
	for _, e := range errors {
		fmt.Printf("  %-25s code=%d  msg=%q\n", e.name, int32(e.err), e.err.Error())
	}

	// ---------------------------------------------------------------
	// Class
	//
	// A Class defines a binder interface. It is created once at
	// initialization time via AIBinder_Class_define:
	//
	//   AIBinder_Class_define(
	//       descriptor,   // e.g. "com.example.IMyService"
	//       onCreateFn,
	//       onDestroyFn,
	//       onTransactFn,
	//   )
	//
	// The descriptor string uniquely identifies the interface and must
	// match between client and server. The Class has no Close method
	// because it is a static definition -- it lives for the process
	// lifetime.
	//
	// The binder.Class struct wraps the AIBinder_Class pointer.
	// Class creation is not yet exposed in the high-level API.
	// ---------------------------------------------------------------
	log.Println("Class: defines a binder interface (descriptor + callbacks)")

	// ---------------------------------------------------------------
	// Binder
	//
	// A Binder represents a handle to a binder object. There are two
	// flavors:
	//
	//   Local binder (server side):
	//     Created with AIBinder_new(class, userData). The userData
	//     pointer is passed to onCreate and can be retrieved later
	//     with AIBinder_getUserData.
	//
	//   Remote proxy (client side):
	//     Obtained from the Android service manager or received as
	//     a parameter in a binder transaction. The proxy transparently
	//     forwards transactions across process boundaries via the
	//     kernel binder driver.
	//
	// Reference counting:
	//   AIBinder_incStrong / AIBinder_decStrong manage the binder's
	//   lifetime. The Go wrapper Binder.Close() calls decStrong.
	//   When the strong count reaches zero, the onDestroy callback
	//   fires and the binder is freed.
	//
	// Example usage pattern:
	//
	//   // Server: create a local binder
	//   //   (AIBinder_new — not yet in high-level API)
	//   defer b.Close()  // calls AIBinder_decStrong
	//
	//   // Client: obtain proxy from service manager
	//   //   (AServiceManager_getService — not yet wrapped)
	//   defer b.Close()
	// ---------------------------------------------------------------
	log.Println("Binder: handle to a local or remote binder object")

	// ---------------------------------------------------------------
	// Parcel
	//
	// A Parcel is the data container for binder transactions. It
	// provides typed serialization for marshaling arguments and
	// replies across process boundaries.
	//
	// Supported data types:
	//   - int32:  WriteInt32 / ReadInt32
	//   - string: WriteString / ReadString
	//
	// Parcels are obtained in two ways:
	//   1. AIBinder_prepareTransaction allocates an input parcel for
	//      a client-initiated transaction.
	//   2. The onTransact callback receives input and output parcels
	//      on the server side.
	//
	// Parcel.Close() deletes the parcel. Only delete parcels that
	// you own -- the output parcel from AIBinder_transact is yours,
	// but parcels received in onTransact are owned by the framework.
	// ---------------------------------------------------------------
	log.Println("Parcel: serialization container for transaction data")

	// ---------------------------------------------------------------
	// Status
	//
	// A Status represents the outcome of a binder operation. The NDK
	// provides:
	//
	//   AStatus_newOk()          - create a success status
	//   AStatus_isOk(status)     - check if the status indicates success
	//   AStatus_getStatus(status) - get the underlying status code
	//   AStatus_delete(status)   - free the status object
	//
	// The onTransact callback returns an AStatus to indicate whether
	// the transaction was handled successfully. Callers inspect the
	// status after AIBinder_transact completes.
	//
	// Status.Close() deletes the status object.
	// ---------------------------------------------------------------
	log.Println("Status: transaction outcome (ok or error code)")

	// ---------------------------------------------------------------
	// Cleanup order
	//
	// Resources should be released in reverse creation order:
	//   1. status.Close()     - delete status objects
	//   2. parcel.Close()     - delete owned parcels
	//   3. binder.Close()     - decrement strong reference
	//
	// Class has no Close -- it is a process-lifetime definition.
	// All Close methods are idempotent and nil-safe.
	// ---------------------------------------------------------------

	log.Println("overview complete -- see source comments for full API documentation")
}
