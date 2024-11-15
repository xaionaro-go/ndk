// Copyright 2018-2024 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package camera

import (
	"github.com/xaionaro-go/ndk/ndk"
)

const (
	FACING_FRONT = 1
	FACING_BACK  = 0

	FLASH_MODE_AUTO    = ndk.CAMERA_FLASH_MODE_AUTO
	FLASH_MODE_OFF     = ndk.CAMERA_FLASH_MODE_OFF
	FLASH_MODE_ON      = ndk.CAMERA_FLASH_MODE_ON
	FLASH_MODE_RED_EYE = ndk.CAMERA_FLASH_MODE_RED_EYE
	FLASH_MODE_TORCH   = ndk.CAMERA_FLASH_MODE_TORCH
	FLASH_MODES_NUM    = ndk.CAMERA_FLASH_MODES_NUM

	FOCUS_MODE_AUTO               = ndk.CAMERA_FOCUS_MODE_AUTO
	FOCUS_MODE_CONTINUOUS_VIDEO   = ndk.CAMERA_FOCUS_MODE_CONTINUOUS_VIDEO
	FOCUS_MODE_EDOF               = ndk.CAMERA_FOCUS_MODE_EDOF
	FOCUS_MODE_FIXED              = ndk.CAMERA_FOCUS_MODE_FIXED
	FOCUS_MODE_INFINITY           = ndk.CAMERA_FOCUS_MODE_INFINITY
	FOCUS_MODE_MACRO              = ndk.CAMERA_FOCUS_MODE_MACRO
	FOCUS_MODE_CONTINUOUS_PICTURE = ndk.CAMERA_FOCUS_MODE_CONTINUOUS_PICTURE
	FOCUS_MODES_NUM               = ndk.CAMERA_FOCUS_MODES_NUM

	WHITE_BALANCE_AUTO             = ndk.CAMERA_WHITE_BALANCE_AUTO
	WHITE_BALANCE_CLOUDY_DAYLIGHT  = ndk.CAMERA_WHITE_BALANCE_CLOUDY_DAYLIGHT
	WHITE_BALANCE_DAYLIGHT         = ndk.CAMERA_WHITE_BALANCE_DAYLIGHT
	WHITE_BALANCE_FLUORESCENT      = ndk.CAMERA_WHITE_BALANCE_FLUORESCENT
	WHITE_BALANCE_INCANDESCENT     = ndk.CAMERA_WHITE_BALANCE_INCANDESCENT
	WHITE_BALANCE_SHADE            = ndk.CAMERA_WHITE_BALANCE_SHADE
	WHITE_BALANCE_TWILIGHT         = ndk.CAMERA_WHITE_BALANCE_TWILIGHT
	WHITE_BALANCE_WARM_FLUORESCENT = ndk.CAMERA_WHITE_BALANCE_WARM_FLUORESCENT
	WHITE_BALANCE_MODES_NUM        = ndk.CAMERA_WHITE_BALANCE_MODES_NUM

	ANTIBANDING_50HZ      = ndk.CAMERA_ANTIBANDING_50HZ
	ANTIBANDING_60HZ      = ndk.CAMERA_ANTIBANDING_60HZ
	ANTIBANDING_AUTO      = ndk.CAMERA_ANTIBANDING_AUTO
	ANTIBANDING_OFF       = ndk.CAMERA_ANTIBANDING_OFF
	ANTIBANDING_MODES_NUM = ndk.CAMERA_ANTIBANDING_MODES_NUM

	FOCUS_DISTANCE_NEAR_INDEX    = ndk.CAMERA_FOCUS_DISTANCE_NEAR_INDEX
	FOCUS_DISTANCE_OPTIMAL_INDEX = ndk.CAMERA_FOCUS_DISTANCE_OPTIMAL_INDEX
	FOCUS_DISTANCE_FAR_INDEX     = ndk.CAMERA_FOCUS_DISTANCE_FAR_INDEX
)

type Camera = ndk.Camera
type CameraCallback = ndk.CameraCallback

func Connect(cameraId int, cb CameraCallback) Camera {
	return ndk.CameraConnect(cameraId, cb)
}
