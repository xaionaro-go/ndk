// Simulates c-for-go output for Android ALooper.
// This file is parsed at AST level only; it does not compile.
package looper

import "unsafe"

// Opaque handle types.
type ALooper C.ALooper

// Integer typedefs.
type Looper_event_t int32

// Event enum.
const (
	ALOOPER_EVENT_INPUT   Looper_event_t = 1
	ALOOPER_EVENT_OUTPUT  Looper_event_t = 2
	ALOOPER_EVENT_ERROR   Looper_event_t = 4
	ALOOPER_EVENT_HANGUP  Looper_event_t = 8
	ALOOPER_EVENT_INVALID Looper_event_t = 16
)

// Poll result constants.
const (
	ALOOPER_POLL_WAKE     int32 = -1
	ALOOPER_POLL_CALLBACK int32 = -2
	ALOOPER_POLL_TIMEOUT  int32 = -3
	ALOOPER_POLL_ERROR    int32 = -4
)

// Callback type.
type ALooper_callbackFunc func(fd int32, events int32, data unsafe.Pointer) int32

// --- Looper functions ---
func ALooper_prepare(opts int32) *ALooper                    { return nil }
func ALooper_acquire(looper *ALooper)                        {}
func ALooper_release(looper *ALooper)                        {}
func ALooper_pollOnce(timeoutMillis int32, outFd *int32, outEvents *int32, outData *unsafe.Pointer) int32 {
	return 0
}
func ALooper_wake(looper *ALooper) {}
func ALooper_addFd(looper *ALooper, fd int32, ident int32, events int32, callback ALooper_callbackFunc, data unsafe.Pointer) int32 {
	return 0
}
func ALooper_removeFd(looper *ALooper, fd int32) int32 { return 0 }

var _ = unsafe.Pointer(nil)
