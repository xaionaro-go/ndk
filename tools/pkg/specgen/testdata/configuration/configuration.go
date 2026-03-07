// Simulates c-for-go output for Android AConfiguration.
// This file is parsed at AST level only; it does not compile.
package configuration

import "unsafe"

// Opaque handle types.
type AConfiguration C.AConfiguration
type AAssetManager C.AAssetManager

// Integer typedefs.
type Configuration_orientation int32
type Configuration_screenSize int32

// Orientation enum.
const (
	ACONFIGURATION_ORIENTATION_ANY    Configuration_orientation = 0
	ACONFIGURATION_ORIENTATION_PORT   Configuration_orientation = 1
	ACONFIGURATION_ORIENTATION_LAND   Configuration_orientation = 2
	ACONFIGURATION_ORIENTATION_SQUARE Configuration_orientation = 3
)

// Screen size enum.
const (
	ACONFIGURATION_SCREENSIZE_ANY    Configuration_screenSize = 0
	ACONFIGURATION_SCREENSIZE_SMALL  Configuration_screenSize = 1
	ACONFIGURATION_SCREENSIZE_NORMAL Configuration_screenSize = 2
	ACONFIGURATION_SCREENSIZE_LARGE  Configuration_screenSize = 3
	ACONFIGURATION_SCREENSIZE_XLARGE Configuration_screenSize = 4
)

// --- Configuration functions ---
func AConfiguration_new() *AConfiguration                                       { return nil }
func AConfiguration_delete(config *AConfiguration)                              {}
func AConfiguration_fromAssetManager(config *AConfiguration, mgr *AAssetManager) {}
func AConfiguration_getLanguage(config *AConfiguration, outLanguage *byte)      {}
func AConfiguration_getCountry(config *AConfiguration, outCountry *byte)        {}
func AConfiguration_getOrientation(config *AConfiguration) int32                { return 0 }
func AConfiguration_getDensity(config *AConfiguration) int32                    { return 0 }
func AConfiguration_getScreenSize(config *AConfiguration) int32                 { return 0 }
func AConfiguration_getScreenWidthDp(config *AConfiguration) int32              { return 0 }
func AConfiguration_getScreenHeightDp(config *AConfiguration) int32             { return 0 }
func AConfiguration_getSdkVersion(config *AConfiguration) int32                 { return 0 }

var _ = unsafe.Pointer(nil)
