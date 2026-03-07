// Simulates c-for-go output for Android thermal API.
// This file is parsed at AST level only; it does not compile.
package thermal

import "unsafe"

// Opaque handle types.
type AThermalManager C.AThermalManager

// Thermal status enum.
type Thermal_status_t int32

const (
	ATHERMAL_STATUS_NONE      Thermal_status_t = 0
	ATHERMAL_STATUS_LIGHT     Thermal_status_t = 1
	ATHERMAL_STATUS_MODERATE  Thermal_status_t = 2
	ATHERMAL_STATUS_SEVERE    Thermal_status_t = 3
	ATHERMAL_STATUS_CRITICAL  Thermal_status_t = 4
	ATHERMAL_STATUS_EMERGENCY Thermal_status_t = 5
	ATHERMAL_STATUS_SHUTDOWN  Thermal_status_t = 6
)

// Callback type.
type AThermal_StatusCallback func(data unsafe.Pointer, status int32)

// --- Thermal functions ---
func AThermal_acquireManager() *AThermalManager                                                                  { return nil }
func AThermal_releaseManager(manager *AThermalManager)                                                           {}
func AThermal_getCurrentThermalStatus(manager *AThermalManager) int32                                            { return 0 }
func AThermal_registerThermalStatusListener(manager *AThermalManager, callback AThermal_StatusCallback, data unsafe.Pointer) int32   { return 0 }
func AThermal_unregisterThermalStatusListener(manager *AThermalManager, callback AThermal_StatusCallback, data unsafe.Pointer) int32 { return 0 }

var _ = unsafe.Pointer(nil)
