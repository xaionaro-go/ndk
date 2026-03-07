// Simulates c-for-go output for Android Camera2 NDK.
// This file is parsed at AST level only; it does not compile.
package camera

import "unsafe"

// Opaque handle types.
type ACameraManager C.ACameraManager
type ACameraDevice C.ACameraDevice
type ACameraCaptureSession C.ACameraCaptureSession
type ACaptureRequest C.ACaptureRequest
type ACameraOutputTarget C.ACameraOutputTarget
type ACaptureSessionOutputContainer C.ACaptureSessionOutputContainer
type ACaptureSessionOutput C.ACaptureSessionOutput
type ACameraMetadata C.ACameraMetadata

// Integer typedefs.
type Camera_status_t int32

// Status codes.
const (
	ACAMERA_OK                              Camera_status_t = 0
	ACAMERA_ERROR_UNKNOWN                   Camera_status_t = -10000
	ACAMERA_ERROR_INVALID_PARAMETER         Camera_status_t = -10001
	ACAMERA_ERROR_CAMERA_DISCONNECTED       Camera_status_t = -10002
	ACAMERA_ERROR_NOT_ENOUGH_MEMORY         Camera_status_t = -10003
	ACAMERA_ERROR_METADATA_NOT_FOUND        Camera_status_t = -10004
	ACAMERA_ERROR_CAMERA_DEVICE             Camera_status_t = -10005
	ACAMERA_ERROR_CAMERA_SERVICE            Camera_status_t = -10006
	ACAMERA_ERROR_SESSION_CLOSED            Camera_status_t = -10007
	ACAMERA_ERROR_INVALID_OPERATION         Camera_status_t = -10008
	ACAMERA_ERROR_STREAM_CONFIGURE_FAIL     Camera_status_t = -10009
	ACAMERA_ERROR_CAMERA_IN_USE             Camera_status_t = -10010
	ACAMERA_ERROR_MAX_CAMERA_IN_USE         Camera_status_t = -10011
	ACAMERA_ERROR_CAMERA_DISABLED           Camera_status_t = -10012
	ACAMERA_ERROR_PERMISSION_DENIED         Camera_status_t = -10013
)

// Template ID enum.
type ACameraDevice_request_template int32

const (
	TEMPLATE_PREVIEW           ACameraDevice_request_template = 1
	TEMPLATE_STILL_CAPTURE     ACameraDevice_request_template = 2
	TEMPLATE_RECORD            ACameraDevice_request_template = 3
	TEMPLATE_VIDEO_SNAPSHOT    ACameraDevice_request_template = 4
	TEMPLATE_ZERO_SHUTTER_LAG  ACameraDevice_request_template = 5
	TEMPLATE_MANUAL            ACameraDevice_request_template = 6
)

// Camera ID list returned by ACameraManager_getCameraIdList.
type ACameraIdList C.ACameraIdList

// Callback structs (opaque to the generator).
type ACameraDevice_StateCallbacks C.ACameraDevice_StateCallbacks
type ACameraCaptureSession_stateCallbacks C.ACameraCaptureSession_stateCallbacks
type ACameraCaptureSession_captureCallbacks C.ACameraCaptureSession_captureCallbacks

// --- CameraManager functions ---
func ACameraManager_create() *ACameraManager                                                         { return nil }
func ACameraManager_delete(manager *ACameraManager)                                                  {}
func ACameraManager_getCameraIdList(manager *ACameraManager, cameraIdList **ACameraIdList) Camera_status_t {
	return 0
}
func ACameraManager_deleteCameraIdList(list *ACameraIdList) {}
func ACameraManager_openCamera(manager *ACameraManager, cameraId *byte, callbacks *ACameraDevice_StateCallbacks, device **ACameraDevice) Camera_status_t {
	return 0
}

// --- CameraDevice functions ---
func ACameraDevice_close(device *ACameraDevice) Camera_status_t { return 0 }
func ACameraDevice_createCaptureRequest(device *ACameraDevice, templateId ACameraDevice_request_template, request **ACaptureRequest) Camera_status_t {
	return 0
}
func ACameraDevice_createCaptureSession(device *ACameraDevice, outputs *ACaptureSessionOutputContainer, callbacks *ACameraCaptureSession_stateCallbacks, session **ACameraCaptureSession) Camera_status_t {
	return 0
}

// --- CaptureRequest functions ---
func ACaptureRequest_free(request *ACaptureRequest)                                                        {}
func ACaptureRequest_addTarget(request *ACaptureRequest, target *ACameraOutputTarget) Camera_status_t      { return 0 }
func ACaptureRequest_removeTarget(request *ACaptureRequest, target *ACameraOutputTarget) Camera_status_t   { return 0 }

// --- OutputTarget functions ---
func ACameraOutputTarget_create(window unsafe.Pointer, target **ACameraOutputTarget) Camera_status_t { return 0 }
func ACameraOutputTarget_free(target *ACameraOutputTarget)                                           {}

// --- CaptureSessionOutputContainer functions ---
func ACaptureSessionOutputContainer_create(container **ACaptureSessionOutputContainer) Camera_status_t { return 0 }
func ACaptureSessionOutputContainer_free(container *ACaptureSessionOutputContainer)                    {}
func ACaptureSessionOutputContainer_add(container *ACaptureSessionOutputContainer, output *ACaptureSessionOutput) Camera_status_t {
	return 0
}
func ACaptureSessionOutputContainer_remove(container *ACaptureSessionOutputContainer, output *ACaptureSessionOutput) Camera_status_t {
	return 0
}

// --- CaptureSessionOutput functions ---
func ACaptureSessionOutput_create(window unsafe.Pointer, output **ACaptureSessionOutput) Camera_status_t { return 0 }
func ACaptureSessionOutput_free(output *ACaptureSessionOutput)                                           {}

// --- CameraCaptureSession functions ---
func ACameraCaptureSession_close(session *ACameraCaptureSession) {}
func ACameraCaptureSession_setRepeatingRequest(session *ACameraCaptureSession, callbacks *ACameraCaptureSession_captureCallbacks, numRequests int32, requests **ACaptureRequest, sequenceId *int32) Camera_status_t {
	return 0
}
func ACameraCaptureSession_stopRepeating(session *ACameraCaptureSession) Camera_status_t { return 0 }
func ACameraCaptureSession_capture(session *ACameraCaptureSession, callbacks *ACameraCaptureSession_captureCallbacks, numRequests int32, requests **ACaptureRequest, sequenceId *int32) Camera_status_t {
	return 0
}

// --- Metadata functions ---
func ACameraMetadata_free(metadata *ACameraMetadata) {}

var _ = unsafe.Pointer(nil)
