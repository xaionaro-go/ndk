//go:build ignore

// Example: full Vulkan device creation pipeline.
//
// Demonstrates the complete device setup sequence: instance creation,
// physical device selection, queue family discovery, logical device
// creation with one queue, and retrieving the queue handle. This is the
// foundation for any Vulkan workload (graphics or compute). All Vulkan
// struct creation is done in C; the actual API calls go through ndk/vulkan.
//
// Build (Android):
//
//	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
//	CC=$NDK/.../aarch64-linux-android35-clang \
//	go build -o create-device ./examples/vulkan/create-device
package main

/*
#include <string.h>
#include <vulkan/vulkan.h>

// Build a VkInstanceCreateInfo with a minimal VkApplicationInfo.
static void make_instance_create_info(void** out) {
	static VkApplicationInfo appInfo = {
		.sType              = VK_STRUCTURE_TYPE_APPLICATION_INFO,
		.pApplicationName   = "create-device",
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

// Extract device name from VkPhysicalDeviceProperties.
static void get_device_name(void* physDev, char* name) {
	VkPhysicalDeviceProperties props;
	vkGetPhysicalDeviceProperties((VkPhysicalDevice)physDev, &props);
	strncpy(name, props.deviceName, 255);
	name[255] = 0;
}

// Query queue family properties. Returns the count via *count. The caller
// must provide a buffer of VkQueueFamilyProperties structs (max entries).
// Each family's queueFlags and queueCount are copied to the parallel output
// arrays for consumption by Go.
static void get_queue_families(void* physDev, uint32_t* count,
                               uint32_t maxFamilies,
                               uint32_t* flags, uint32_t* queueCounts) {
	VkQueueFamilyProperties families[32];
	uint32_t n = 0;
	vkGetPhysicalDeviceQueueFamilyProperties((VkPhysicalDevice)physDev, &n, NULL);
	if (n > 32) n = 32;
	vkGetPhysicalDeviceQueueFamilyProperties((VkPhysicalDevice)physDev, &n, families);
	if (n > maxFamilies) n = maxFamilies;
	*count = n;
	for (uint32_t i = 0; i < n; i++) {
		flags[i]       = families[i].queueFlags;
		queueCounts[i] = families[i].queueCount;
	}
}

// Build a VkDeviceCreateInfo requesting one queue from the specified family.
static void make_device_create_info(uint32_t queueFamily, void** out) {
	static float priority = 1.0f;
	static VkDeviceQueueCreateInfo qci;
	qci = (VkDeviceQueueCreateInfo){
		.sType            = VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO,
		.queueFamilyIndex = queueFamily,
		.queueCount       = 1,
		.pQueuePriorities = &priority,
	};
	static VkDeviceCreateInfo dci;
	dci = (VkDeviceCreateInfo){
		.sType                = VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO,
		.queueCreateInfoCount = 1,
		.pQueueCreateInfos    = &qci,
	};
	*out = &dci;
}
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/xaionaro-go/ndk/vulkan"
)

// Vulkan queue capability bits.
const (
	vkQueueGraphics = 0x01
	vkQueueCompute  = 0x02
	vkQueueTransfer = 0x04
)

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

// flagsString formats VkQueueFlags as a human-readable string.
func flagsString(flags uint32) string {
	s := ""
	if flags&vkQueueGraphics != 0 {
		s += "GRAPHICS "
	}
	if flags&vkQueueCompute != 0 {
		s += "COMPUTE "
	}
	if flags&vkQueueTransfer != 0 {
		s += "TRANSFER "
	}
	if s == "" {
		return fmt.Sprintf("0x%X", flags)
	}
	return s[:len(s)-1] // trim trailing space
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

	// 2. Pick the first physical device.
	var deviceCount uint32
	ret = vulkan.EnumeratePhysicalDevices(instance, &deviceCount, nil)
	if ret != 0 || deviceCount == 0 {
		fatal("no Vulkan physical devices available")
	}

	var physDev vulkan.VkPhysicalDevice
	count := uint32(1)
	vulkan.EnumeratePhysicalDevices(instance, &count, &physDev)

	var nameBuf [256]byte
	C.get_device_name(unsafe.Pointer(physDev), (*C.char)(unsafe.Pointer(&nameBuf[0])))
	nameLen := 0
	for nameLen < len(nameBuf) && nameBuf[nameLen] != 0 {
		nameLen++
	}
	fmt.Printf("selected device: %s\n", string(nameBuf[:nameLen]))

	// 3. Enumerate queue families and find one that supports graphics.
	var familyCount C.uint32_t
	var flags [32]C.uint32_t
	var queueCounts [32]C.uint32_t
	C.get_queue_families(
		unsafe.Pointer(physDev), &familyCount, 32,
		&flags[0], &queueCounts[0],
	)
	fmt.Printf("found %d queue family/families\n", familyCount)

	graphicsFamily := -1
	for i := 0; i < int(familyCount); i++ {
		f := uint32(flags[i])
		qc := uint32(queueCounts[i])
		fmt.Printf("  family[%d]: flags=%s queues=%d\n", i, flagsString(f), qc)
		if graphicsFamily < 0 && f&vkQueueGraphics != 0 {
			graphicsFamily = i
		}
	}
	if graphicsFamily < 0 {
		fatal("no queue family supports graphics")
	}
	fmt.Printf("using queue family %d (supports graphics)\n", graphicsFamily)

	// 4. Create a logical device with one queue from the chosen family.
	var devCreateInfo unsafe.Pointer
	C.make_device_create_info(C.uint32_t(graphicsFamily), &devCreateInfo)

	var device vulkan.VkDevice
	ret = vulkan.CreateDevice(physDev, devCreateInfo, nil, &device)
	if ret != 0 {
		fatalf("CreateDevice failed (VkResult=%d)", ret)
	}
	defer func() {
		vulkan.DeviceWaitIdle(device)
		vulkan.DestroyDevice(device, nil)
		fmt.Println("device destroyed")
	}()
	fmt.Println("logical device created")

	// 5. Retrieve the queue handle.
	var queue vulkan.VkQueue
	vulkan.GetDeviceQueue(device, uint32(graphicsFamily), 0, &queue)
	if queue == nil {
		fatal("GetDeviceQueue returned nil")
	}
	fmt.Println("obtained graphics queue")

	// 6. Wait for idle to verify the device is functional.
	ret = vulkan.QueueWaitIdle(queue)
	if ret != 0 {
		fatalf("QueueWaitIdle failed (VkResult=%d)", ret)
	}
	fmt.Println("queue idle -- device is operational")
}
