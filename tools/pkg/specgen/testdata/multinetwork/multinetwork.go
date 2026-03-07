// Simulates c-for-go output for Android Multinetwork.
// This file is parsed at AST level only; it does not compile.
package multinetwork

import "unsafe"

// Integer typedefs.
type Net_handle_t int64

// --- Network functions ---
func android_getaddrinfofornetwork(network int64, node *byte, service *byte, hints unsafe.Pointer, res *unsafe.Pointer) int32 {
	return 0
}
func android_setprocnetwork(network int64) int32           { return 0 }
func android_setsocknetwork(network int64, fd int32) int32 { return 0 }
func android_getprocnetwork(network *int64) int32          { return 0 }
func android_tag_socket(sockfd int32, tag uint32) int32    { return 0 }
func android_untag_socket(sockfd int32) int32              { return 0 }

var _ = unsafe.Pointer(nil)
