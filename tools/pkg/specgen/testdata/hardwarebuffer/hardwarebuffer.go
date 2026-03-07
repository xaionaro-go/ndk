// Simulates c-for-go output for Android HardwareBuffer.
// This file is parsed at AST level only; it does not compile.
package hardwarebuffer

import "unsafe"

// Opaque handle type.
type AHardwareBuffer C.AHardwareBuffer

// Struct type for buffer description.
type AHardwareBuffer_Desc C.AHardwareBuffer_Desc

// Integer typedefs.
type HardwareBuffer_format_t int32
type HardwareBuffer_usage_t uint64

// Format enum.
const (
	AHARDWAREBUFFER_FORMAT_R8G8B8A8_UNORM    HardwareBuffer_format_t = 1
	AHARDWAREBUFFER_FORMAT_R8G8B8X8_UNORM    HardwareBuffer_format_t = 2
	AHARDWAREBUFFER_FORMAT_R8G8B8_UNORM      HardwareBuffer_format_t = 3
	AHARDWAREBUFFER_FORMAT_R5G6B5_UNORM      HardwareBuffer_format_t = 4
	AHARDWAREBUFFER_FORMAT_R16G16B16A16_FLOAT HardwareBuffer_format_t = 0x16
	AHARDWAREBUFFER_FORMAT_BLOB              HardwareBuffer_format_t = 0x21
)

// Usage enum.
const (
	AHARDWAREBUFFER_USAGE_GPU_SAMPLED_IMAGE HardwareBuffer_usage_t = 0x100
	AHARDWAREBUFFER_USAGE_GPU_COLOR_OUTPUT  HardwareBuffer_usage_t = 0x200
	AHARDWAREBUFFER_USAGE_CPU_READ_OFTEN    HardwareBuffer_usage_t = 3
	AHARDWAREBUFFER_USAGE_CPU_WRITE_OFTEN   HardwareBuffer_usage_t = 0x30
)

// --- HardwareBuffer functions ---
func AHardwareBuffer_allocate(desc *AHardwareBuffer_Desc, outBuffer **AHardwareBuffer) int32 { return 0 }
func AHardwareBuffer_release(buffer *AHardwareBuffer)                                        {}
func AHardwareBuffer_acquire(buffer *AHardwareBuffer)                                        {}
func AHardwareBuffer_describe(buffer *AHardwareBuffer, outDesc *AHardwareBuffer_Desc)        {}
func AHardwareBuffer_lock(buffer *AHardwareBuffer, usage uint64, fence int32, rect unsafe.Pointer, outVirtualAddress *unsafe.Pointer) int32 { return 0 }
func AHardwareBuffer_unlock(buffer *AHardwareBuffer, fence *int32) int32                     { return 0 }
func AHardwareBuffer_isSupported(desc *AHardwareBuffer_Desc) int32                           { return 0 }

var _ = unsafe.Pointer(nil)
