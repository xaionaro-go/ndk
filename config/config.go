// Copyright 2018-2024 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/xaionaro-go/ndk/ndk"
)

const (
	ORIENTATION_ANY    = ndk.CONFIGURATION_ORIENTATION_ANY
	ORIENTATION_PORT   = ndk.CONFIGURATION_ORIENTATION_PORT
	ORIENTATION_LAND   = ndk.CONFIGURATION_ORIENTATION_LAND
	ORIENTATION_SQUARE = ndk.CONFIGURATION_ORIENTATION_SQUARE

	TOUCHSCREEN_ANY     = ndk.CONFIGURATION_TOUCHSCREEN_ANY
	TOUCHSCREEN_NOTOUCH = ndk.CONFIGURATION_TOUCHSCREEN_NOTOUCH
	TOUCHSCREEN_STYLUS  = ndk.CONFIGURATION_TOUCHSCREEN_STYLUS
	TOUCHSCREEN_FINGER  = ndk.CONFIGURATION_TOUCHSCREEN_FINGER

	DENSITY_DEFAULT = ndk.CONFIGURATION_DENSITY_DEFAULT
	DENSITY_LOW     = ndk.CONFIGURATION_DENSITY_LOW
	DENSITY_MEDIUM  = ndk.CONFIGURATION_DENSITY_MEDIUM
	DENSITY_TV      = ndk.CONFIGURATION_DENSITY_TV
	DENSITY_HIGH    = ndk.CONFIGURATION_DENSITY_HIGH
	DENSITY_XHIGH   = ndk.CONFIGURATION_DENSITY_XHIGH
	DENSITY_XXHIGH  = ndk.CONFIGURATION_DENSITY_XXHIGH
	DENSITY_XXXHIGH = ndk.CONFIGURATION_DENSITY_XXXHIGH
	DENSITY_NONE    = ndk.CONFIGURATION_DENSITY_NONE

	KEYBOARD_ANY    = ndk.CONFIGURATION_KEYBOARD_ANY
	KEYBOARD_NOKEYS = ndk.CONFIGURATION_KEYBOARD_NOKEYS
	KEYBOARD_QWERTY = ndk.CONFIGURATION_KEYBOARD_QWERTY
	KEYBOARD_12KEY  = ndk.CONFIGURATION_KEYBOARD_12KEY

	NAVIGATION_ANY       = ndk.CONFIGURATION_NAVIGATION_ANY
	NAVIGATION_NONAV     = ndk.CONFIGURATION_NAVIGATION_NONAV
	NAVIGATION_DPAD      = ndk.CONFIGURATION_NAVIGATION_DPAD
	NAVIGATION_TRACKBALL = ndk.CONFIGURATION_NAVIGATION_TRACKBALL
	NAVIGATION_WHEEL     = ndk.CONFIGURATION_NAVIGATION_WHEEL

	KEYSHIDDEN_ANY  = ndk.CONFIGURATION_KEYSHIDDEN_ANY
	KEYSHIDDEN_NO   = ndk.CONFIGURATION_KEYSHIDDEN_NO
	KEYSHIDDEN_YES  = ndk.CONFIGURATION_KEYSHIDDEN_YES
	KEYSHIDDEN_SOFT = ndk.CONFIGURATION_KEYSHIDDEN_SOFT

	NAVHIDDEN_ANY = ndk.CONFIGURATION_NAVHIDDEN_ANY
	NAVHIDDEN_NO  = ndk.CONFIGURATION_NAVHIDDEN_NO
	NAVHIDDEN_YES = ndk.CONFIGURATION_NAVHIDDEN_YES

	SCREENSIZE_ANY    = ndk.CONFIGURATION_SCREENSIZE_ANY
	SCREENSIZE_SMALL  = ndk.CONFIGURATION_SCREENSIZE_SMALL
	SCREENSIZE_NORMAL = ndk.CONFIGURATION_SCREENSIZE_NORMAL
	SCREENSIZE_LARGE  = ndk.CONFIGURATION_SCREENSIZE_LARGE
	SCREENSIZE_XLARGE = ndk.CONFIGURATION_SCREENSIZE_XLARGE

	SCREENLONG_ANY = ndk.CONFIGURATION_SCREENLONG_ANY
	SCREENLONG_NO  = ndk.CONFIGURATION_SCREENLONG_NO
	SCREENLONG_YES = ndk.CONFIGURATION_SCREENLONG_YES

	UI_MODE_TYPE_ANY        = ndk.CONFIGURATION_UI_MODE_TYPE_ANY
	UI_MODE_TYPE_NORMAL     = ndk.CONFIGURATION_UI_MODE_TYPE_NORMAL
	UI_MODE_TYPE_DESK       = ndk.CONFIGURATION_UI_MODE_TYPE_DESK
	UI_MODE_TYPE_CAR        = ndk.CONFIGURATION_UI_MODE_TYPE_CAR
	UI_MODE_TYPE_TELEVISION = ndk.CONFIGURATION_UI_MODE_TYPE_TELEVISION
	UI_MODE_TYPE_APPLIANCE  = ndk.CONFIGURATION_UI_MODE_TYPE_APPLIANCE

	UI_MODE_NIGHT_ANY = ndk.CONFIGURATION_UI_MODE_NIGHT_ANY
	UI_MODE_NIGHT_NO  = ndk.CONFIGURATION_UI_MODE_NIGHT_NO
	UI_MODE_NIGHT_YES = ndk.CONFIGURATION_UI_MODE_NIGHT_YES

	SCREEN_WIDTH_DP_ANY = ndk.CONFIGURATION_SCREEN_WIDTH_DP_ANY

	SCREEN_HEIGHT_DP_ANY = ndk.CONFIGURATION_SCREEN_HEIGHT_DP_ANY

	SMALLEST_SCREEN_WIDTH_DP_ANY = ndk.CONFIGURATION_SMALLEST_SCREEN_WIDTH_DP_ANY

	LAYOUTDIR_ANY = ndk.CONFIGURATION_LAYOUTDIR_ANY
	LAYOUTDIR_LTR = ndk.CONFIGURATION_LAYOUTDIR_LTR
	LAYOUTDIR_RTL = ndk.CONFIGURATION_LAYOUTDIR_RTL

	MCC                  = ndk.CONFIGURATION_MCC
	MNC                  = ndk.CONFIGURATION_MNC
	LOCALE               = ndk.CONFIGURATION_LOCALE
	TOUCHSCREEN          = ndk.CONFIGURATION_TOUCHSCREEN
	KEYBOARD             = ndk.CONFIGURATION_KEYBOARD
	KEYBOARD_HIDDEN      = ndk.CONFIGURATION_KEYBOARD_HIDDEN
	NAVIGATION           = ndk.CONFIGURATION_NAVIGATION
	ORIENTATION          = ndk.CONFIGURATION_ORIENTATION
	DENSITY              = ndk.CONFIGURATION_DENSITY
	SCREEN_SIZE          = ndk.CONFIGURATION_SCREEN_SIZE
	VERSION              = ndk.CONFIGURATION_VERSION
	SCREEN_LAYOUT        = ndk.CONFIGURATION_SCREEN_LAYOUT
	UI_MODE              = ndk.CONFIGURATION_UI_MODE
	SMALLEST_SCREEN_SIZE = ndk.CONFIGURATION_SMALLEST_SCREEN_SIZE
	LAYOUTDIR            = ndk.CONFIGURATION_LAYOUTDIR
)

type Configuration = ndk.Configuration
