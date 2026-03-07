//go:build ignore

// Example: Vulkan buffer creation with host-visible memory.
//
// Demonstrates the Vulkan buffer-memory workflow: creating a VkBuffer,
// querying its memory requirements, finding a suitable host-visible memory
// type, allocating device memory, binding it to the buffer, mapping the
// memory to write data from the CPU, reading it back for verification, and
// cleaning everything up. This pattern is the basis for staging buffers and
// host-visible uniform/storage buffers.
//
// Build (Android):
//
//	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
//	CC=$NDK/.../aarch64-linux-android35-clang \
//	go build -o compute-buffer ./examples/vulkan/compute-buffer
package main

/*
#include <string.h>
#include <vulkan/vulkan.h>

// Build a VkInstanceCreateInfo with a minimal VkApplicationInfo.
static void make_instance_create_info(void** out) {
	static VkApplicationInfo appInfo = {
		.sType              = VK_STRUCTURE_TYPE_APPLICATION_INFO,
		.pApplicationName   = "compute-buffer",
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

// Query queue families and return the first family index that supports
// compute or graphics. Returns -1 if none found.
static int find_queue_family(void* physDev) {
	uint32_t n = 0;
	vkGetPhysicalDeviceQueueFamilyProperties((VkPhysicalDevice)physDev, &n, NULL);
	VkQueueFamilyProperties families[32];
	if (n > 32) n = 32;
	vkGetPhysicalDeviceQueueFamilyProperties((VkPhysicalDevice)physDev, &n, families);
	for (uint32_t i = 0; i < n; i++) {
		if (families[i].queueFlags & (VK_QUEUE_COMPUTE_BIT | VK_QUEUE_GRAPHICS_BIT)) {
			return (int)i;
		}
	}
	return -1;
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

// Build a VkBufferCreateInfo for a host-visible buffer of the given size.
// Usage includes TRANSFER_SRC | TRANSFER_DST to allow copies.
static void make_buffer_create_info(uint64_t size, void** out) {
	static VkBufferCreateInfo bci;
	bci = (VkBufferCreateInfo){
		.sType       = VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO,
		.size        = size,
		.usage       = VK_BUFFER_USAGE_TRANSFER_SRC_BIT | VK_BUFFER_USAGE_TRANSFER_DST_BIT,
		.sharingMode = VK_SHARING_MODE_EXCLUSIVE,
	};
	*out = &bci;
}

// Query buffer memory requirements. Returns size, alignment, and memory
// type bits through output parameters.
static void get_buffer_memory_requirements(void* device, uint64_t buffer,
                                           uint64_t* size, uint64_t* alignment,
                                           uint32_t* memoryTypeBits) {
	VkMemoryRequirements reqs;
	vkGetBufferMemoryRequirements((VkDevice)device, (VkBuffer)buffer, &reqs);
	*size           = reqs.size;
	*alignment      = reqs.alignment;
	*memoryTypeBits = reqs.memoryTypeBits;
}

// Query physical device memory properties. Returns the number of memory
// types and each type's propertyFlags through the output arrays.
static void get_memory_properties(void* physDev, uint32_t* typeCount,
                                  uint32_t* propertyFlags) {
	VkPhysicalDeviceMemoryProperties memProps;
	vkGetPhysicalDeviceMemoryProperties((VkPhysicalDevice)physDev, &memProps);
	*typeCount = memProps.memoryTypeCount;
	for (uint32_t i = 0; i < memProps.memoryTypeCount && i < 32; i++) {
		propertyFlags[i] = memProps.memoryTypes[i].propertyFlags;
	}
}

// Build a VkMemoryAllocateInfo for the given size and memory type index.
static void make_memory_allocate_info(uint64_t size, uint32_t memTypeIndex, void** out) {
	static VkMemoryAllocateInfo mai;
	mai = (VkMemoryAllocateInfo){
		.sType           = VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO,
		.allocationSize  = size,
		.memoryTypeIndex = memTypeIndex,
	};
	*out = &mai;
}

// Bind buffer memory. Returns VkResult.
static int32_t bind_buffer_memory(void* device, uint64_t buffer,
                                  uint64_t memory, uint64_t offset) {
	return (int32_t)vkBindBufferMemory((VkDevice)device, (VkBuffer)buffer,
	                                   (VkDeviceMemory)memory, (VkDeviceSize)offset);
}
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/xaionaro-go/ndk/vulkan"
)

// Vulkan memory property flags.
const (
	vkMemoryPropertyHostVisible  = 0x02
	vkMemoryPropertyHostCoherent = 0x04
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
	const bufferSize = 256 // bytes

	// --- Instance ---
	var createInfo unsafe.Pointer
	C.make_instance_create_info(&createInfo)

	var instance vulkan.VkInstance
	ret := vulkan.CreateInstance(createInfo, nil, &instance)
	if ret != 0 {
		fatalf("CreateInstance failed (VkResult=%d)", ret)
	}
	defer vulkan.DestroyInstance(instance, nil)
	fmt.Println("Vulkan instance created")

	// --- Physical device ---
	var deviceCount uint32
	ret = vulkan.EnumeratePhysicalDevices(instance, &deviceCount, nil)
	if ret != 0 || deviceCount == 0 {
		fatal("no Vulkan physical devices available")
	}
	var physDev vulkan.VkPhysicalDevice
	count := uint32(1)
	vulkan.EnumeratePhysicalDevices(instance, &count, &physDev)

	// --- Queue family ---
	queueFamily := int(C.find_queue_family(unsafe.Pointer(physDev)))
	if queueFamily < 0 {
		fatal("no compute/graphics queue family found")
	}

	// --- Logical device ---
	var devCreateInfo unsafe.Pointer
	C.make_device_create_info(C.uint32_t(queueFamily), &devCreateInfo)

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

	// --- Create buffer ---
	var bufCI unsafe.Pointer
	C.make_buffer_create_info(C.uint64_t(bufferSize), &bufCI)

	var buffer vulkan.VkBuffer
	ret = vulkan.CreateBuffer(device, bufCI, nil, &buffer)
	if ret != 0 {
		fatalf("CreateBuffer failed (VkResult=%d)", ret)
	}
	defer vulkan.DestroyBuffer(device, buffer, nil)
	fmt.Printf("buffer created (%d bytes requested)\n", bufferSize)

	// --- Query memory requirements ---
	var reqSize, reqAlign C.uint64_t
	var memTypeBits C.uint32_t
	C.get_buffer_memory_requirements(
		unsafe.Pointer(device), C.uint64_t(buffer),
		&reqSize, &reqAlign, &memTypeBits,
	)
	fmt.Printf("memory requirements: size=%d alignment=%d typeBits=0x%X\n",
		reqSize, reqAlign, memTypeBits)

	// --- Find a host-visible, host-coherent memory type ---
	var typeCount C.uint32_t
	var propFlags [32]C.uint32_t
	C.get_memory_properties(unsafe.Pointer(physDev), &typeCount, &propFlags[0])

	requiredFlags := uint32(vkMemoryPropertyHostVisible | vkMemoryPropertyHostCoherent)
	memTypeIndex := -1
	for i := 0; i < int(typeCount); i++ {
		supported := uint32(memTypeBits)&(1<<uint(i)) != 0
		hasFlags := uint32(propFlags[i])&requiredFlags == requiredFlags
		if supported && hasFlags {
			memTypeIndex = i
			break
		}
	}
	if memTypeIndex < 0 {
		fatal("no host-visible coherent memory type available for this buffer")
	}
	fmt.Printf("selected memory type index %d (host-visible + host-coherent)\n", memTypeIndex)

	// --- Allocate memory ---
	var allocInfo unsafe.Pointer
	C.make_memory_allocate_info(reqSize, C.uint32_t(memTypeIndex), &allocInfo)

	var memory vulkan.VkDeviceMemory
	ret = vulkan.AllocateMemory(device, allocInfo, nil, &memory)
	if ret != 0 {
		fatalf("AllocateMemory failed (VkResult=%d)", ret)
	}
	defer vulkan.FreeMemory(device, memory, nil)
	fmt.Printf("allocated %d bytes of device memory\n", reqSize)

	// --- Bind buffer to memory ---
	// vkBindBufferMemory is not in the generated Go bindings, so call via C.
	bindRet := C.bind_buffer_memory(
		unsafe.Pointer(device), C.uint64_t(buffer),
		C.uint64_t(memory), C.uint64_t(0),
	)
	if bindRet != 0 {
		fatalf("vkBindBufferMemory failed (VkResult=%d)", bindRet)
	}
	fmt.Println("buffer bound to memory")

	// --- Map memory, write data, read it back ---
	var mapped unsafe.Pointer
	ret = vulkan.MapMemory(device, memory, 0, vulkan.VkDeviceSize(reqSize), 0, &mapped)
	if ret != 0 {
		fatalf("MapMemory failed (VkResult=%d)", ret)
	}
	fmt.Println("memory mapped")

	// Write a recognizable pattern: 0, 1, 2, ..., 255.
	dst := unsafe.Slice((*byte)(mapped), int(reqSize))
	for i := 0; i < bufferSize && i < len(dst); i++ {
		dst[i] = byte(i)
	}
	fmt.Printf("wrote %d bytes to mapped memory\n", bufferSize)

	// Read back and verify.
	ok := true
	for i := 0; i < bufferSize && i < len(dst); i++ {
		if dst[i] != byte(i) {
			fmt.Printf("  MISMATCH at offset %d: expected %d, got %d\n", i, byte(i), dst[i])
			ok = false
			break
		}
	}
	if ok {
		fmt.Println("readback verification passed")
	}

	vulkan.UnmapMemory(device, memory)
	fmt.Println("memory unmapped")
}
