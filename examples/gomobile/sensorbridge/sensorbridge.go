// Package sensorbridge exposes NDK sensor access to Java/Kotlin
// via gomobile bind.
//
// Build with:
//
//	gomobile bind -target=android -androidapi=26 ./examples/gomobile/sensorbridge
//
// This produces sensorbridge.aar + sensorbridge-sources.jar.
//
// Java usage:
//
//	import sensorbridge.Sensorbridge;
//	import sensorbridge.Bridge;
//
//	long ptr = nativeGetSensorManager(); // from JNI
//	Bridge bridge = Sensorbridge.newBridge(ptr);
//	String name = bridge.accelerometerName();
//	bridge.close();
package sensorbridge

import (
	"github.com/AndroidGoLab/ndk/sensor"
)

// Bridge wraps an ASensorManager for cross-language use.
// Unexported fields are invisible to Java/Kotlin.
type Bridge struct {
	mgr *sensor.Manager
}

// NewBridge creates a Bridge from a raw ASensorManager* passed as int64
// (Java long). gomobile bind does not support unsafe.Pointer or uintptr,
// so native handles must be transported as int64.
func NewBridge(managerHandle int64) *Bridge {
	return &Bridge{
		mgr: sensor.NewManagerFromUintPtr(uintptr(managerHandle)),
	}
}

// Handle returns the underlying ASensorManager* as int64 for passing
// back to Java/JNI code.
func (b *Bridge) Handle() int64 {
	return int64(b.mgr.UintPtr())
}

// AccelerometerName returns the name of the default accelerometer,
// or empty string if unavailable.
func (b *Bridge) AccelerometerName() string {
	return b.sensorName(sensor.Accelerometer)
}

// GyroscopeName returns the name of the default gyroscope,
// or empty string if unavailable.
func (b *Bridge) GyroscopeName() string {
	return b.sensorName(sensor.Gyroscope)
}

// LightSensorName returns the name of the default light sensor,
// or empty string if unavailable.
func (b *Bridge) LightSensorName() string {
	return b.sensorName(sensor.Light)
}

// SensorName returns the name of a sensor by its Android type code,
// or empty string if the sensor is not available.
func (b *Bridge) SensorName(sensorType int32) string {
	return b.sensorName(sensor.Type(sensorType))
}

func (b *Bridge) sensorName(t sensor.Type) string {
	s := b.mgr.DefaultSensor(t)
	if s.UintPtr() == 0 {
		return ""
	}
	return s.Name()
}

// Close is a no-op — the ASensorManager is a singleton owned by the
// system. Provided so the Bridge satisfies Java's Closeable pattern.
func (b *Bridge) Close() {
	b.mgr = nil
}
