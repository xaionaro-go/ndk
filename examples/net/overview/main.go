//go:build ignore

// Android Multinetwork API overview.
//
// Documents the Network handle type exported by ndk/net. On Android, the
// system supports multiple simultaneous network interfaces (Wi-Fi, cellular,
// VPN, Ethernet). Each active network is identified by a Network handle -- a
// 64-bit integer assigned by the ConnectivityManager framework service.
//
// How a Network handle is obtained:
//
// The handle originates on the Java/Kotlin side. A typical flow:
//
//  1. Obtain a ConnectivityManager from the Android Context.
//  2. Register a NetworkCallback or call getActiveNetwork() to get a
//     android.net.Network object.
//  3. Call Network.getNetworkHandle() to get the int64 handle.
//  4. Pass the handle to native code via JNI.
//
// In Go code using gomobile or NativeActivity, the JNI bridge delivers
// the int64 value which can be stored directly as a net.Network.
//
// How a Network handle is used in the NDK:
//
// The NDK <android/multinetwork.h> header declares several functions
// that accept a net_handle_t (the C equivalent of net.Network):
//
//	android_setsockoptfornetwork(net_handle_t, int fd)
//	    Binds an existing socket to the specified network so that all
//	    traffic on that socket routes through it.
//
//	android_setprocdns(net_handle_t)
//	    Sets the default DNS servers for the calling process to those
//	    associated with the given network.
//
//	android_getaddrinfofornetwork(net_handle_t, ...)
//	    Performs DNS resolution using the DNS servers of the specified
//	    network, similar to getaddrinfo(3) but network-aware.
//
//	android_res_nquery / android_res_nresult / android_res_nsend / ...
//	    Asynchronous DNS resolution functions bound to a network.
//
// Special values:
//
//	0  (NETWORK_UNSPECIFIED) -- no specific network; the system uses its
//	   default routing. This is the zero value of net.Network.
//
// This program runs on the host (no Android device required) because it
// only inspects the Network type without calling into the NDK.
package main

import (
	"fmt"
	"unsafe"

	"github.com/xaionaro-go/ndk/net"
)

func main() {
	fmt.Println("=== ndk/net API overview ===")
	fmt.Println()

	// -----------------------------------------------------------------
	// Network type
	//
	// net.Network is a type alias for the NDK net_handle_t, which is
	// an int64. It uniquely identifies an Android network interface
	// within the ConnectivityManager's lifetime.
	// -----------------------------------------------------------------

	var n net.Network
	fmt.Printf("Type:       net.Network (alias for int64)\n")
	fmt.Printf("Size:       %d bytes\n", unsafe.Sizeof(n))
	fmt.Printf("Zero value: %d\n", n)
	fmt.Println()

	// -----------------------------------------------------------------
	// Zero value semantics
	//
	// A zero-valued Network (NETWORK_UNSPECIFIED) tells the NDK
	// functions to use the system default network. This matches the
	// NDK constant NETWORK_UNSPECIFIED defined in multinetwork.h.
	// -----------------------------------------------------------------

	fmt.Println("Zero value semantics:")
	fmt.Println("  A net.Network of 0 is NETWORK_UNSPECIFIED.")
	fmt.Println("  When passed to NDK functions it means \"use the default")
	fmt.Println("  network\" -- the system decides which interface to use")
	fmt.Println("  based on its current routing rules.")
	fmt.Println()

	// -----------------------------------------------------------------
	// Obtaining a handle from Java
	// -----------------------------------------------------------------

	fmt.Println("Obtaining a Network handle (Java/Kotlin side):")
	fmt.Println()
	fmt.Println("  ConnectivityManager cm = context.getSystemService(")
	fmt.Println("      ConnectivityManager.class);")
	fmt.Println("  android.net.Network network = cm.getActiveNetwork();")
	fmt.Println("  long handle = network.getNetworkHandle();")
	fmt.Println("  // Pass 'handle' to native via JNI.")
	fmt.Println()

	// -----------------------------------------------------------------
	// Using the handle in Go
	// -----------------------------------------------------------------

	fmt.Println("Using the handle in Go (pseudocode):")
	fmt.Println()
	fmt.Println("  // Receive the handle from JNI as int64.")
	fmt.Println("  var network net.Network = handleFromJNI")
	fmt.Println()
	fmt.Println("  // Bind a socket to this network:")
	fmt.Println("  //   android_setsockoptfornetwork(network, fd)")
	fmt.Println()
	fmt.Println("  // Resolve DNS via this network:")
	fmt.Println("  //   android_getaddrinfofornetwork(network, host, ...)")
	fmt.Println()

	// -----------------------------------------------------------------
	// Example handle values
	//
	// Real handle values are opaque integers assigned by the framework.
	// They are not sequential or predictable. Here are illustrative
	// values to show the type in action.
	// -----------------------------------------------------------------

	fmt.Println("Example handle values (illustrative):")
	fmt.Println()

	examples := []struct {
		name   string
		handle net.Network
	}{
		{"NETWORK_UNSPECIFIED (default)", 0},
		{"Typical Wi-Fi handle", 100},
		{"Typical cellular handle", 200},
	}

	for _, ex := range examples {
		fmt.Printf("  %-35s handle = %d\n", ex.name, ex.handle)
	}
	fmt.Println()

	// -----------------------------------------------------------------
	// Comparison and storage
	//
	// Because Network is an int64, handles can be compared with == and
	// used as map keys. This is useful for maintaining per-network
	// state such as DNS caches or socket pools.
	// -----------------------------------------------------------------

	fmt.Println("Comparison and storage:")

	a := net.Network(100)
	b := net.Network(200)
	fmt.Printf("  net.Network(100) == net.Network(200): %v\n", a == b)
	fmt.Printf("  net.Network(100) == net.Network(100): %v\n", a == net.Network(100))
	fmt.Println()
	fmt.Println("  Network handles can be used as map keys:")
	fmt.Println("    networkState := map[net.Network]string{")
	fmt.Printf("        %d: \"wifi\",\n", a)
	fmt.Printf("        %d: \"cellular\",\n", b)
	fmt.Println("    }")
}
