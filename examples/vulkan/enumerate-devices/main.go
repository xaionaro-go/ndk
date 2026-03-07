//go:build ignore

// Example: Vulkan instance creation and physical device enumeration.
//
// Demonstrates the first steps of any Vulkan program on Android: creating a
// VkInstance, enumerating all physical devices (GPUs), and querying their
// properties (device name and API version). All Vulkan struct creation is
// done in C because struct layout is a C concern; the actual Vulkan API
// calls go through the idiomatic ndk/vulkan Go package.
//
// Build (Android):
//
//	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
//	CC=$NDK/.../aarch64-linux-android35-clang \
//	go build -o enumerate-devices ./examples/vulkan/enumerate-devices
package main

/*
#include <string.h>
#include <vulkan/vulkan.h>

// Build a VkInstanceCreateInfo with a minimal VkApplicationInfo.
static void make_instance_create_info(void** out) {
	static VkApplicationInfo appInfo = {
		.sType              = VK_STRUCTURE_TYPE_APPLICATION_INFO,
		.pApplicationName   = "enumerate-devices",
		.applicationVersion = 1,
		.pEngineName        = "ndk",
		.engineVersion      = 1,
		.apiVersion         = VK_API_VERSION_1_0,
	};
	static VkInstanceCreateInfo ci = {
		.sType            = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
		.pApplicationInfo = &appInfo,
	};
	*out = &ci;
}

// Extract device name and API version from VkPhysicalDeviceProperties.
// The caller passes the VkPhysicalDevice handle; this helper queries
// the properties via the C Vulkan API and copies the results out.
static void get_device_props(void* physDev, char* name, uint32_t* apiVer) {
	VkPhysicalDeviceProperties props;
	vkGetPhysicalDeviceProperties((VkPhysicalDevice)physDev, &props);
	strncpy(name, props.deviceName, 255);
	name[255] = 0;
	*apiVer = props.apiVersion;
}
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/xaionaro-go/ndk/vulkan"
)

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	// 1. Create a Vulkan instance.
	var createInfo unsafe.Pointer
	C.make_instance_create_info(&createInfo)

	var instance vulkan.VkInstance
	ret := vulkan.CreateInstance(createInfo, nil, &instance)
	if ret != 0 {
		fatalf("CreateInstance failed (VkResult=%d)", ret)
	}
	defer vulkan.DestroyInstance(instance, nil)
	fmt.Println("Vulkan instance created")

	// 2. Query the number of physical devices.
	var deviceCount uint32
	ret = vulkan.EnumeratePhysicalDevices(instance, &deviceCount, nil)
	if ret != 0 {
		fatalf("EnumeratePhysicalDevices (count) failed (VkResult=%d)", ret)
	}
	fmt.Printf("found %d physical device(s)\n", deviceCount)

	if deviceCount == 0 {
		fatal("no Vulkan physical devices available")
	}

	// 3. Retrieve all physical device handles.
	devices := make([]vulkan.VkPhysicalDevice, deviceCount)
	ret = vulkan.EnumeratePhysicalDevices(instance, &deviceCount, &devices[0])
	if ret != 0 {
		fatalf("EnumeratePhysicalDevices (list) failed (VkResult=%d)", ret)
	}

	// 4. Print properties for each device.
	for i, dev := range devices {
		var nameBuf [256]byte
		var apiVer C.uint32_t
		C.get_device_props(
			unsafe.Pointer(dev),
			(*C.char)(unsafe.Pointer(&nameBuf[0])),
			&apiVer,
		)

		// Decode Vulkan API version: major(10).minor(10).patch(12).
		major := uint32(apiVer) >> 22
		minor := (uint32(apiVer) >> 12) & 0x3FF
		patch := uint32(apiVer) & 0xFFF

		// Find null terminator for the name.
		nameLen := 0
		for nameLen < len(nameBuf) && nameBuf[nameLen] != 0 {
			nameLen++
		}
		name := string(nameBuf[:nameLen])

		fmt.Printf("  device[%d]: name=%q api=%d.%d.%d\n", i, name, major, minor, patch)
	}

	fmt.Println("instance destroyed")
}
