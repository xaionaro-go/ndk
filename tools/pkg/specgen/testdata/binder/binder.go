// Simulates c-for-go output for Android Binder.
// This file is parsed at AST level only; it does not compile.
package binder

import "unsafe"

// Opaque handle types.
type AIBinder C.AIBinder
type AIBinder_Class C.AIBinder_Class
type AParcel C.AParcel
type AStatus C.AStatus

// Integer typedefs.
type Binder_status_t int32

// Status codes.
const (
	STATUS_OK                Binder_status_t = 0
	STATUS_UNKNOWN_ERROR     Binder_status_t = -2147483647
	STATUS_NO_MEMORY         Binder_status_t = -12
	STATUS_INVALID_OPERATION Binder_status_t = -38
	STATUS_BAD_VALUE         Binder_status_t = -22
	STATUS_DEAD_OBJECT       Binder_status_t = -32
	STATUS_PERMISSION_DENIED Binder_status_t = -1
	STATUS_NAME_NOT_FOUND    Binder_status_t = -2
	STATUS_WOULD_BLOCK       Binder_status_t = -11
	STATUS_FDS_NOT_ALLOWED   Binder_status_t = -2147483641
)

// Callback type.
type AIBinder_Class_onTransact func(binder *AIBinder, code uint32, in *AParcel, out *AParcel) int32

// --- Class functions ---
func AIBinder_Class_define(interfaceDescriptor *byte, onCreate unsafe.Pointer, onDestroy unsafe.Pointer, onTransact AIBinder_Class_onTransact) *AIBinder_Class {
	return nil
}

// --- Binder functions ---
func AIBinder_new(clazz *AIBinder_Class, args unsafe.Pointer) *AIBinder      { return nil }
func AIBinder_incStrong(binder *AIBinder)                                     {}
func AIBinder_decStrong(binder *AIBinder)                                     {}
func AIBinder_prepareTransaction(binder *AIBinder, in **AParcel) int32        { return 0 }
func AIBinder_transact(binder *AIBinder, code uint32, in **AParcel, out **AParcel, flags uint32) int32 {
	return 0
}
func AIBinder_getUserData(binder *AIBinder) unsafe.Pointer { return nil }

// --- Parcel functions ---
func AParcel_writeInt32(parcel *AParcel, value int32) int32                              { return 0 }
func AParcel_writeString(parcel *AParcel, string *byte, length int32) int32              { return 0 }
func AParcel_readInt32(parcel *AParcel, value *int32) int32                              { return 0 }
func AParcel_readString(parcel *AParcel, stringData unsafe.Pointer, allocator unsafe.Pointer) int32 {
	return 0
}
func AParcel_delete(parcel *AParcel) {}

// --- Status functions ---
func AStatus_newOk() *AStatus              { return nil }
func AStatus_delete(status *AStatus)        {}
func AStatus_isOk(status *AStatus) bool     { return false }
func AStatus_getStatus(status *AStatus) int32 { return 0 }

var _ = unsafe.Pointer(nil)
