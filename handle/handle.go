// Package handle defines interfaces satisfied by all NDK handle types,
// enabling generic code and interoperability with external Android Go packages.
package handle

import "unsafe"

// NativeHandle is satisfied by any NDK opaque handle wrapper.
// All generated handle types (Window, Looper, Sensor, Stream, etc.)
// implement this interface via their Pointer() and UintPtr() methods.
type NativeHandle interface {
	// Pointer returns the underlying C pointer as unsafe.Pointer.
	Pointer() unsafe.Pointer

	// UintPtr returns the underlying C pointer as uintptr,
	// suitable for gomobile bind and framework interop
	// (golang.org/x/mobile, gioui.org, github.com/xlab/android-go).
	UintPtr() uintptr
}
