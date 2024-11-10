// Copyright 2018-2024 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sensor

import (
	"github.com/xaionaro-go/ndk/ndk"
)

type TYPE = ndk.SENSOR_TYPE

type AdditionalInfoEvent = ndk.AdditionalInfoEvent
type DynamicSensorEvent = ndk.DynamicSensorEvent
type HeartRateEvent = ndk.HeartRateEvent
type MetaDataEvent = ndk.MetaDataEvent
type UncalibratedEvent = ndk.UncalibratedEvent
type Vector = ndk.SensorVector
type MagneticVector = ndk.MagneticVector

type Sensor = ndk.Sensor
type Event = ndk.SensorEvent
type Manager = ndk.SensorManager

const (
	FIFO_COUNT_INVALID = ndk.SENSOR_FIFO_COUNT_INVALID
	DELAY_INVALID      = ndk.SENSOR_DELAY_INVALID

	/**
	 * Invalid sensor type. Returned by {@link ASensor_getType} as error value.
	 */
	TYPE_INVALID = ndk.SENSOR_TYPE_INVALID
	/**
	 * {@link ASENSOR_TYPE_ACCELEROMETER}
	 * reporting-mode: continuous
	 *
	 *  All values are in SI units (m/s^2) and measure the acceleration of the
	 *  device minus the force of gravity.
	 */
	TYPE_ACCELEROMETER = ndk.SENSOR_TYPE_ACCELEROMETER
	/**
	 * {@link ASENSOR_TYPE_MAGNETIC_FIELD}
	 * reporting-mode: continuous
	 *
	 *  All values are in micro-Tesla (uT) and measure the geomagnetic
	 *  field in the X, Y and Z axis.
	 */
	TYPE_MAGNETIC_FIELD = ndk.SENSOR_TYPE_MAGNETIC_FIELD
	/**
	 * {@link ASENSOR_TYPE_GYROSCOPE}
	 * reporting-mode: continuous
	 *
	 *  All values are in radians/second and measure the rate of rotation
	 *  around the X, Y and Z axis.
	 */
	TYPE_GYROSCOPE = ndk.SENSOR_TYPE_GYROSCOPE
	/**
	 * {@link ASENSOR_TYPE_LIGHT}
	 * reporting-mode: on-change
	 *
	 * The light sensor value is returned in SI lux units.
	 */
	TYPE_LIGHT = ndk.SENSOR_TYPE_LIGHT
	/**
	 * {@link ASENSOR_TYPE_PROXIMITY}
	 * reporting-mode: on-change
	 *
	 * The proximity sensor which turns the screen off and back on during calls is the
	 * wake-up proximity sensor. Implement wake-up proximity sensor before implementing
	 * a non wake-up proximity sensor. For the wake-up proximity sensor set the flag
	 * SENSOR_FLAG_WAKE_UP.
	 * The value corresponds to the distance to the nearest object in centimeters.
	 */
	TYPE_PROXIMITY = ndk.SENSOR_TYPE_PROXIMITY
	/**
	 * {@link ASENSOR_TYPE_LINEAR_ACCELERATION}
	 * reporting-mode: continuous
	 *
	 *  All values are in SI units (m/s^2) and measure the acceleration of the
	 *  device not including the force of gravity.
	 */
	TYPE_LINEAR_ACCELERATION = ndk.SENSOR_TYPE_LINEAR_ACCELERATION

	// java
	TYPE_ORIENTATION                 = ndk.SENSOR_TYPE_ORIENTATION
	TYPE_PRESSURE                    = ndk.SENSOR_TYPE_PRESSURE
	TYPE_TEMPERATURE                 = ndk.SENSOR_TYPE_TEMPERATURE
	TYPE_GRAVITY                     = ndk.SENSOR_TYPE_GRAVITY
	TYPE_ROTATION_VECTOR             = ndk.SENSOR_TYPE_ROTATION_VECTOR
	TYPE_RELATIVE_HUMIDITY           = ndk.SENSOR_TYPE_RELATIVE_HUMIDITY
	TYPE_AMBIENT_TEMPERATURE         = ndk.SENSOR_TYPE_AMBIENT_TEMPERATURE
	TYPE_MAGNETIC_FIELD_UNCALIBRATED = ndk.SENSOR_TYPE_MAGNETIC_FIELD_UNCALIBRATED
	TYPE_GAME_ROTATION_VECTOR        = ndk.SENSOR_TYPE_GAME_ROTATION_VECTOR
	TYPE_GYROSCOPE_UNCALIBRATED      = ndk.SENSOR_TYPE_GYROSCOPE_UNCALIBRATED
	TYPE_SIGNIFICANT_MOTION          = ndk.SENSOR_TYPE_SIGNIFICANT_MOTION
	TYPE_STEP_DETECTOR               = ndk.SENSOR_TYPE_STEP_DETECTOR
	TYPE_STEP_COUNTER                = ndk.SENSOR_TYPE_STEP_COUNTER
	TYPE_GEOMAGNETIC_ROTATION_VECTOR = ndk.SENSOR_TYPE_GEOMAGNETIC_ROTATION_VECTOR
	TYPE_HEART_RATE                  = ndk.SENSOR_TYPE_HEART_RATE
	TYPE_TILT_DETECTOR               = ndk.SENSOR_TYPE_TILT_DETECTOR
	TYPE_WAKE_GESTURE                = ndk.SENSOR_TYPE_WAKE_GESTURE
	TYPE_GLANCE_GESTURE              = ndk.SENSOR_TYPE_GLANCE_GESTURE
	TYPE_PICK_UP_GESTURE             = ndk.SENSOR_TYPE_PICK_UP_GESTURE
	TYPE_WRIST_TILT_GESTURE          = ndk.SENSOR_TYPE_WRIST_TILT_GESTURE
	TYPE_DEVICE_ORIENTATION          = ndk.SENSOR_TYPE_DEVICE_ORIENTATION
	TYPE_POSE_6DOF                   = ndk.SENSOR_TYPE_POSE_6DOF
	TYPE_STATIONARY_DETECT           = ndk.SENSOR_TYPE_STATIONARY_DETECT
	TYPE_MOTION_DETECT               = ndk.SENSOR_TYPE_MOTION_DETECT
	TYPE_HEART_BEAT                  = ndk.SENSOR_TYPE_HEART_BEAT
	TYPE_DYNAMIC_SENSOR_META         = ndk.SENSOR_TYPE_DYNAMIC_SENSOR_META
	TYPE_LOW_LATENCY_OFFBODY_DETECT  = ndk.SENSOR_TYPE_LOW_LATENCY_OFFBODY_DETECT
	TYPE_ACCELEROMETER_UNCALIBRATED  = ndk.SENSOR_TYPE_ACCELEROMETER_UNCALIBRATED

	/** no contact */
	STATUS_NO_CONTACT = ndk.SENSOR_STATUS_NO_CONTACT
	/** unreliable */
	STATUS_UNRELIABLE = ndk.SENSOR_STATUS_UNRELIABLE
	/** low accuracy */
	STATUS_ACCURACY_LOW = ndk.SENSOR_STATUS_ACCURACY_LOW
	/** medium accuracy */
	STATUS_ACCURACY_MEDIUM = ndk.SENSOR_STATUS_ACCURACY_MEDIUM
	/** high accuracy */
	STATUS_ACCURACY_HIGH = ndk.SENSOR_STATUS_ACCURACY_HIGH

	/** invalid reporting mode */
	REPORTING_MODE_INVALID = ndk.REPORTING_MODE_INVALID
	/** continuous reporting */
	REPORTING_MODE_CONTINUOUS = ndk.REPORTING_MODE_CONTINUOUS
	/** reporting on change */
	REPORTING_MODE_ON_CHANGE = ndk.REPORTING_MODE_ON_CHANGE
	/** on shot reporting */
	REPORTING_MODE_ONE_SHOT = ndk.REPORTING_MODE_ONE_SHOT
	/** special trigger reporting */
	REPORTING_MODE_SPECIAL_TRIGGER = ndk.REPORTING_MODE_SPECIAL_TRIGGER

	/** stopped */
	DIRECT_RATE_STOP = ndk.SENSOR_DIRECT_RATE_STOP
	/** nominal 50Hz */
	DIRECT_RATE_NORMAL = ndk.SENSOR_DIRECT_RATE_NORMAL
	/** nominal 200Hz */
	DIRECT_RATE_FAST = ndk.SENSOR_DIRECT_RATE_FAST
	/** nominal 800Hz */
	DIRECT_RATE_VERY_FAST = ndk.SENSOR_DIRECT_RATE_VERY_FAST

	/** shared memory created by ASharedMemory_create */
	DIRECT_CHANNEL_TYPE_SHARED_MEMORY = ndk.SENSOR_DIRECT_CHANNEL_TYPE_SHARED_MEMORY
	/** AHardwareBuffer */
	DIRECT_CHANNEL_TYPE_HARDWARE_BUFFER = ndk.SENSOR_DIRECT_CHANNEL_TYPE_HARDWARE_BUFFER

	/** Earth's gravity in m/s^2 */
	STANDARD_GRAVITY = ndk.SENSOR_STANDARD_GRAVITY
	/** Maximum magnetic field on Earth's surface in uT */
	MAGNETIC_FIELD_EARTH_MAX = ndk.SENSOR_MAGNETIC_FIELD_EARTH_MAX
	/** Minimum magnetic field on Earth's surface in uT*/
	MAGNETIC_FIELD_EARTH_MIN = ndk.SENSOR_MAGNETIC_FIELD_EARTH_MIN
)

func ManagerInstance() *Manager {
	return ndk.SensorManagerInstance()
}
