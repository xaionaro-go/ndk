// Simulates c-for-go output for Android Sensor NDK.
// This file is parsed at AST level only; it does not compile.
package sensor

import "unsafe"

// Opaque handle types.
type ASensorManager C.ASensorManager
type ASensor C.ASensor
type ASensorEventQueue C.ASensorEventQueue

// ASensorList is a pointer to a list of sensor pointers.
type ASensorList **ASensor

// Integer typedefs.
type Sensor_type_t int32
type Sensor_status_t int32

// Sensor type enum.
const (
	ASENSOR_TYPE_ACCELEROMETER        Sensor_type_t = 1
	ASENSOR_TYPE_MAGNETIC_FIELD       Sensor_type_t = 2
	ASENSOR_TYPE_GYROSCOPE            Sensor_type_t = 4
	ASENSOR_TYPE_LIGHT                Sensor_type_t = 5
	ASENSOR_TYPE_PRESSURE             Sensor_type_t = 6
	ASENSOR_TYPE_PROXIMITY            Sensor_type_t = 8
	ASENSOR_TYPE_GRAVITY              Sensor_type_t = 9
	ASENSOR_TYPE_LINEAR_ACCELERATION  Sensor_type_t = 10
	ASENSOR_TYPE_ROTATION_VECTOR      Sensor_type_t = 11
	ASENSOR_TYPE_RELATIVE_HUMIDITY    Sensor_type_t = 12
	ASENSOR_TYPE_AMBIENT_TEMPERATURE  Sensor_type_t = 13
	ASENSOR_TYPE_GAME_ROTATION_VECTOR Sensor_type_t = 15
	ASENSOR_TYPE_STEP_DETECTOR        Sensor_type_t = 18
	ASENSOR_TYPE_STEP_COUNTER         Sensor_type_t = 19
	ASENSOR_TYPE_HEART_RATE           Sensor_type_t = 21
)

// Sensor status enum.
const (
	ASENSOR_STATUS_NO_CONTACT      Sensor_status_t = -1
	ASENSOR_STATUS_UNRELIABLE      Sensor_status_t = 0
	ASENSOR_STATUS_ACCURACY_LOW    Sensor_status_t = 1
	ASENSOR_STATUS_ACCURACY_MEDIUM Sensor_status_t = 2
	ASENSOR_STATUS_ACCURACY_HIGH   Sensor_status_t = 3
)

// --- SensorManager functions ---
func ASensorManager_getInstance() *ASensorManager                    { return nil }
func ASensorManager_getDefaultSensor(manager *ASensorManager, sensorType int32) *ASensor {
	return nil
}
func ASensorManager_getSensorList(manager *ASensorManager, list **ASensor) int32 { return 0 }
func ASensorManager_createEventQueue(manager *ASensorManager, looper unsafe.Pointer, ident int32, callback unsafe.Pointer, data unsafe.Pointer) *ASensorEventQueue {
	return nil
}
func ASensorManager_destroyEventQueue(manager *ASensorManager, queue *ASensorEventQueue) int32 {
	return 0
}

// --- SensorEventQueue functions ---
func ASensorEventQueue_enableSensor(queue *ASensorEventQueue, sensor *ASensor) int32  { return 0 }
func ASensorEventQueue_disableSensor(queue *ASensorEventQueue, sensor *ASensor) int32 { return 0 }
func ASensorEventQueue_setEventRate(queue *ASensorEventQueue, sensor *ASensor, usec int32) int32 {
	return 0
}
func ASensorEventQueue_hasEvents(queue *ASensorEventQueue) int32 { return 0 }

// --- Sensor info functions ---
func ASensor_getName(sensor *ASensor) *byte       { return nil }
func ASensor_getVendor(sensor *ASensor) *byte     { return nil }
func ASensor_getType(sensor *ASensor) int32       { return 0 }
func ASensor_getResolution(sensor *ASensor) float32 { return 0 }
func ASensor_getMinDelay(sensor *ASensor) int32   { return 0 }

var _ = unsafe.Pointer(nil)
