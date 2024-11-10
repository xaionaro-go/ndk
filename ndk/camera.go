//go:build android
// +build android

package ndk

/*
#include "camera.c"

#cgo android CFLAGS: -D__ANDROID_API__=24
#cgo android LDFLAGS: -lcamera2ndk -lmediandk -llog -landroid
*/
import "C"

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

type Camera uintptr
type CameraCallback func([]byte) bool

func CameraConnect(cameraId int, cb CameraCallback) Camera {
	panic("not implemented")
}

//export cgoCameraCallback
func cgoCameraCallback(buffer *C.char, bufferSize C.size_t, userData unsafe.Pointer) C.int {
	if userData != nil {
		userData = unsafe.Pointer(^(uintptr(userData)))
		buf := (*[1 << 28]byte)(unsafe.Pointer(buffer))[:int(bufferSize)]
		if (*(*CameraCallback)((unsafe.Pointer)(&userData)))(buf) {
			return 1
		}
		return 0
	}
	return 1
}

func (c Camera) Disconnect() {
	panic("not implemented")
}

func (c Camera) getProperty(propIdx int) float64 {
	panic("not implemented")
}

func (c Camera) setProperty(propIdx int, value float64) {
	panic("not implemented")
}

func (c Camera) ApplyProperties() {
	panic("not implemented")
}

func (c Camera) FrameSize() (int, int) {
	return int(c.getProperty(iCAMERA_PROPERTY_FRAMEWIDTH)), int(c.getProperty(iCAMERA_PROPERTY_FRAMEHEIGHT))
}

func (c Camera) SupportedPreviewSizes() string {
	panic("not implemented")
}

func (c Camera) PreviewFormat() string {
	panic("not implemented")
}

func (c Camera) Fps() int {
	return int(c.getProperty(iCAMERA_PROPERTY_FPS))
}

func (c Camera) Exposure() int {
	return int(c.getProperty(iCAMERA_PROPERTY_EXPOSURE))
}

func (c Camera) FlashMode() CameraFlashMode {
	return CameraFlashMode(c.getProperty(iCAMERA_PROPERTY_FLASH_MODE))
}

func (c Camera) FocusMode() CameraFocusMode {
	return CameraFocusMode(c.getProperty(iCAMERA_PROPERTY_FOCUS_MODE))
}

func (c Camera) WhiteBalance() CameraWhiteBalance {
	return CameraWhiteBalance(c.getProperty(iCAMERA_PROPERTY_WHITE_BALANCE))
}

func (c Camera) Antibanding() CameraAntibanding {
	return CameraAntibanding(c.getProperty(iCAMERA_PROPERTY_ANTIBANDING))
}

func (c Camera) FocalLength() float32 {
	return float32(c.getProperty(iCAMERA_PROPERTY_FOCAL_LENGTH))
}

func (c Camera) FocusDistance(index CameraFocusDistance) float32 {
	switch index {
	case CAMERA_FOCUS_DISTANCE_NEAR_INDEX:
		return float32(c.getProperty(iCAMERA_PROPERTY_FOCUS_DISTANCE_NEAR))

	case CAMERA_FOCUS_DISTANCE_OPTIMAL_INDEX:
		return float32(c.getProperty(iCAMERA_PROPERTY_FOCUS_DISTANCE_OPTIMAL))

	case CAMERA_FOCUS_DISTANCE_FAR_INDEX:
		return float32(c.getProperty(iCAMERA_PROPERTY_FOCUS_DISTANCE_FAR))
	}
	return 0
}

func (c Camera) IsExposeLock() bool {
	return 0 != int(c.getProperty(iCAMERA_PROPERTY_EXPOSE_LOCK))
}

func (c Camera) IsWhiteBalanceLock() bool {
	return 0 != int(c.getProperty(iCAMERA_PROPERTY_WHITEBALANCE_LOCK))
}

// SET
func (c Camera) SetFrameSize(w, h int) {
	c.setProperty(iCAMERA_PROPERTY_FRAMEWIDTH, float64(w))
	c.setProperty(iCAMERA_PROPERTY_FRAMEHEIGHT, float64(h))
}

func (c Camera) SetExposure(v int) {
	c.setProperty(iCAMERA_PROPERTY_EXPOSURE, float64(v))
}

func (c Camera) SetFlashMode(v CameraFlashMode) {
	c.setProperty(iCAMERA_PROPERTY_FLASH_MODE, float64(v))
}

func (c Camera) SetFocusMode(v CameraFocusMode) {
	c.setProperty(iCAMERA_PROPERTY_FOCUS_MODE, float64(v))
}

func (c Camera) SetWhiteBalance(v CameraWhiteBalance) {
	c.setProperty(iCAMERA_PROPERTY_WHITE_BALANCE, float64(v))
}

func (c Camera) SetAntibanding(v CameraAntibanding) {
	c.setProperty(iCAMERA_PROPERTY_ANTIBANDING, float64(v))
}

func (c Camera) SetExposeLock(b bool) {
	v := float64(0)
	if b {
		v = 1
	}
	c.setProperty(iCAMERA_PROPERTY_EXPOSE_LOCK, v)
}

func (c Camera) SetWhiteBalanceLock(b bool) {
	v := float64(0)
	if b {
		v = 1
	}
	c.setProperty(iCAMERA_PROPERTY_WHITEBALANCE_LOCK, v)
}

type CameraFlashMode int
type CameraFocusMode int
type CameraWhiteBalance int
type CameraAntibanding int
type CameraFocusDistance int

const (
	iCAMERA_PROPERTY_FRAMEWIDTH                     = C.ANDROID_CAMERA_PROPERTY_FRAMEWIDTH
	iCAMERA_PROPERTY_FRAMEHEIGHT                    = C.ANDROID_CAMERA_PROPERTY_FRAMEHEIGHT
	iCAMERA_PROPERTY_SUPPORTED_PREVIEW_SIZES_STRING = C.ANDROID_CAMERA_PROPERTY_SUPPORTED_PREVIEW_SIZES_STRING
	iCAMERA_PROPERTY_PREVIEW_FORMAT_STRING          = C.ANDROID_CAMERA_PROPERTY_PREVIEW_FORMAT_STRING
	iCAMERA_PROPERTY_FPS                            = C.ANDROID_CAMERA_PROPERTY_FPS
	iCAMERA_PROPERTY_EXPOSURE                       = C.ANDROID_CAMERA_PROPERTY_EXPOSURE
	iCAMERA_PROPERTY_FLASH_MODE                     = C.ANDROID_CAMERA_PROPERTY_FLASH_MODE
	iCAMERA_PROPERTY_FOCUS_MODE                     = C.ANDROID_CAMERA_PROPERTY_FOCUS_MODE
	iCAMERA_PROPERTY_WHITE_BALANCE                  = C.ANDROID_CAMERA_PROPERTY_WHITE_BALANCE
	iCAMERA_PROPERTY_ANTIBANDING                    = C.ANDROID_CAMERA_PROPERTY_ANTIBANDING
	iCAMERA_PROPERTY_FOCAL_LENGTH                   = C.ANDROID_CAMERA_PROPERTY_FOCAL_LENGTH
	iCAMERA_PROPERTY_FOCUS_DISTANCE_NEAR            = C.ANDROID_CAMERA_PROPERTY_FOCUS_DISTANCE_NEAR
	iCAMERA_PROPERTY_FOCUS_DISTANCE_OPTIMAL         = C.ANDROID_CAMERA_PROPERTY_FOCUS_DISTANCE_OPTIMAL
	iCAMERA_PROPERTY_FOCUS_DISTANCE_FAR             = C.ANDROID_CAMERA_PROPERTY_FOCUS_DISTANCE_FAR
	iCAMERA_PROPERTY_EXPOSE_LOCK                    = C.ANDROID_CAMERA_PROPERTY_EXPOSE_LOCK
	iCAMERA_PROPERTY_WHITEBALANCE_LOCK              = C.ANDROID_CAMERA_PROPERTY_WHITEBALANCE_LOCK

	CAMERA_FLASH_MODE_AUTO    = CameraFlashMode(C.ANDROID_CAMERA_FLASH_MODE_AUTO)
	CAMERA_FLASH_MODE_OFF     = CameraFlashMode(C.ANDROID_CAMERA_FLASH_MODE_OFF)
	CAMERA_FLASH_MODE_ON      = CameraFlashMode(C.ANDROID_CAMERA_FLASH_MODE_ON)
	CAMERA_FLASH_MODE_RED_EYE = CameraFlashMode(C.ANDROID_CAMERA_FLASH_MODE_RED_EYE)
	CAMERA_FLASH_MODE_TORCH   = CameraFlashMode(C.ANDROID_CAMERA_FLASH_MODE_TORCH)
	CAMERA_FLASH_MODES_NUM    = CameraFlashMode(C.ANDROID_CAMERA_FLASH_MODES_NUM)

	CAMERA_FOCUS_MODE_AUTO               = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_AUTO)
	CAMERA_FOCUS_MODE_CONTINUOUS_VIDEO   = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_CONTINUOUS_VIDEO)
	CAMERA_FOCUS_MODE_EDOF               = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_EDOF)
	CAMERA_FOCUS_MODE_FIXED              = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_FIXED)
	CAMERA_FOCUS_MODE_INFINITY           = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_INFINITY)
	CAMERA_FOCUS_MODE_MACRO              = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_MACRO)
	CAMERA_FOCUS_MODE_CONTINUOUS_PICTURE = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODE_CONTINUOUS_PICTURE)
	CAMERA_FOCUS_MODES_NUM               = CameraFocusMode(C.ANDROID_CAMERA_FOCUS_MODES_NUM)

	CAMERA_WHITE_BALANCE_AUTO             = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_AUTO)
	CAMERA_WHITE_BALANCE_CLOUDY_DAYLIGHT  = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_CLOUDY_DAYLIGHT)
	CAMERA_WHITE_BALANCE_DAYLIGHT         = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_DAYLIGHT)
	CAMERA_WHITE_BALANCE_FLUORESCENT      = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_FLUORESCENT)
	CAMERA_WHITE_BALANCE_INCANDESCENT     = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_INCANDESCENT)
	CAMERA_WHITE_BALANCE_SHADE            = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_SHADE)
	CAMERA_WHITE_BALANCE_TWILIGHT         = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_TWILIGHT)
	CAMERA_WHITE_BALANCE_WARM_FLUORESCENT = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_WARM_FLUORESCENT)
	CAMERA_WHITE_BALANCE_MODES_NUM        = CameraWhiteBalance(C.ANDROID_CAMERA_WHITE_BALANCE_MODES_NUM)

	CAMERA_ANTIBANDING_50HZ      = CameraAntibanding(C.ANDROID_CAMERA_ANTIBANDING_50HZ)
	CAMERA_ANTIBANDING_60HZ      = CameraAntibanding(C.ANDROID_CAMERA_ANTIBANDING_60HZ)
	CAMERA_ANTIBANDING_AUTO      = CameraAntibanding(C.ANDROID_CAMERA_ANTIBANDING_AUTO)
	CAMERA_ANTIBANDING_OFF       = CameraAntibanding(C.ANDROID_CAMERA_ANTIBANDING_OFF)
	CAMERA_ANTIBANDING_MODES_NUM = CameraAntibanding(C.ANDROID_CAMERA_ANTIBANDING_MODES_NUM)

	CAMERA_FOCUS_DISTANCE_NEAR_INDEX    = CameraFocusDistance(C.ANDROID_CAMERA_FOCUS_DISTANCE_NEAR_INDEX)
	CAMERA_FOCUS_DISTANCE_OPTIMAL_INDEX = CameraFocusDistance(C.ANDROID_CAMERA_FOCUS_DISTANCE_OPTIMAL_INDEX)
	CAMERA_FOCUS_DISTANCE_FAR_INDEX     = CameraFocusDistance(C.ANDROID_CAMERA_FOCUS_DISTANCE_FAR_INDEX)
)

func GetPackageName() string {
	smapsName := fmt.Sprint("/proc/", os.Getpid(), "/cmdline")
	b, err := os.ReadFile(smapsName)
	if err != nil {
		return ""
	}
	i := bytes.Index(b, []byte{0})
	if i >= 0 {
		b = b[:i]
	}
	return string(b)
}

func GetLibraryPath() string {
	smapsName := fmt.Sprint("/proc/", os.Getpid(), "/smaps")
	b, err := os.ReadFile(smapsName)
	if err != nil {
		return ""
	}

	for {
		base := []byte("/data/")
		i := bytes.Index(b, base)
		if i < 0 {
			break
		}
		j := bytes.Index(b[i:], []byte("\n"))
		if j < 0 {
			break
		}
		path := string(b[i : i+j])
		if strings.HasSuffix(path, ".so") {
			return filepath.Dir(path)
		}

		b = b[i+j:]
	}
	return ""
}

func FindMatchLibrary(pattern string) []string {
	paths := os.Getenv("LD_LIBRARY_PATH")
	if paths == "" {
		paths = "/system/lib"
	}
	dirs := strings.Split(paths, ":")
	dirs = append([]string{GetLibraryPath()}, dirs...)
	for _, dir := range dirs {
		fns, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			Info("FindMatchLibrary:", err)
		}
		if len(fns) > 0 {
			return fns
		}
	}
	return nil
}
