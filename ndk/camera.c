#include <stdlib.h>

#include <android/log.h>

#include <camera/NdkCameraDevice.h>
#include <camera/NdkCameraManager.h>
#include <media/NdkImageReader.h>

#include "camera_properties.h"

#define TAG "go-ndk"

#define LOG_ERROR(...) __android_log_print(ANDROID_LOG_ERROR, TAG, __VA_ARGS__)
#define LOG_WARN(...)  __android_log_print(ANDROID_LOG_WARN, TAG, __VA_ARGS__)
#define LOG_INFO(...)  __android_log_print(ANDROID_LOG_INFO, TAG, __VA_ARGS__)
#define LOG_DEBUG(...) __android_log_print(ANDROID_LOG_DEBUG, TAG, __VA_ARGS__)

static void camera_device_on_disconnected(void *context, ACameraDevice *device) {
    LOG_INFO("camera %s is disconnected.\n", ACameraDevice_getId(device));
}

static void camera_device_on_error(void *context, ACameraDevice *device, int error) {
    LOG_ERROR("error %d on camera %s.\n", error, ACameraDevice_getId(device));
}

static void camera_session_on_ready(void *context, ACameraCaptureSession *session) {
    LOG_INFO("session is ready. %p\n", session);
}

static void camera_session_on_active(void *context, ACameraCaptureSession *session) {
    LOG_INFO("session is activated. %p\n", session);
}

static void camera_session_on_closed(void *context, ACameraCaptureSession *session) {
    LOG_INFO("session is closed. %p\n", session);
}

static void image_callback(void *context, AImageReader *reader) {
		LOG_ERROR("NOT IMPLEMENTED");
}

struct camera {
    ACameraDevice_stateCallbacks deviceStateCallbacks;
};
typedef struct camera camera_t;

static void *NewCamera() {
    camera_t *camera = (camera_t *)malloc(sizeof(camera_t));
    camera->deviceStateCallbacks.context = camera;
    camera->deviceStateCallbacks.onDisconnected = camera_device_on_disconnected;
    camera->deviceStateCallbacks.onError = camera_device_on_error;
    return camera;
}
