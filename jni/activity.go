// Package jni provides JNI helper functions for Android NativeActivity
// operations that have no NDK C API equivalent (permission dialogs, Toast
// messages, SurfaceTexture creation from Java).
//
// This package is hand-written (not auto-generated) because JNI calls require
// inline C with JNIEnv method tables.
package jni

/*
#cgo LDFLAGS: -landroid -lnativewindow

#include <stdlib.h>
#include <jni.h>
#include <android/native_activity.h>
#include <android/native_window.h>
#include <android/surface_texture.h>
#include <android/surface_texture_jni.h>

static int hasPermission(ANativeActivity* act, const char* perm) {
	JNIEnv* env = act->env;
	jclass cls = (*env)->GetObjectClass(env, act->clazz);
	jmethodID mid = (*env)->GetMethodID(env, cls, "checkSelfPermission",
		"(Ljava/lang/String;)I");
	if (mid == NULL) {
		(*env)->ExceptionClear(env);
		(*env)->DeleteLocalRef(env, cls);
		return 0;
	}
	jstring jperm = (*env)->NewStringUTF(env, perm);
	jint result = (*env)->CallIntMethod(env, act->clazz, mid, jperm);
	(*env)->DeleteLocalRef(env, jperm);
	(*env)->DeleteLocalRef(env, cls);
	return (result == 0) ? 1 : 0;
}

static void requestPermission(ANativeActivity* act, const char* perm) {
	JNIEnv* env = act->env;
	jclass cls = (*env)->GetObjectClass(env, act->clazz);
	jmethodID mid = (*env)->GetMethodID(env, cls, "requestPermissions",
		"([Ljava/lang/String;I)V");
	if (mid == NULL) {
		(*env)->ExceptionClear(env);
		(*env)->DeleteLocalRef(env, cls);
		return;
	}
	jclass strCls = (*env)->FindClass(env, "java/lang/String");
	jstring jperm = (*env)->NewStringUTF(env, perm);
	jobjectArray arr = (*env)->NewObjectArray(env, 1, strCls, jperm);
	(*env)->CallVoidMethod(env, act->clazz, mid, arr, (jint)1);
	(*env)->DeleteLocalRef(env, arr);
	(*env)->DeleteLocalRef(env, jperm);
	(*env)->DeleteLocalRef(env, strCls);
	(*env)->DeleteLocalRef(env, cls);
}

static void showToast(ANativeActivity* act, const char* msg) {
	JNIEnv* env = act->env;
	jclass actCls = (*env)->GetObjectClass(env, act->clazz);
	jmethodID getCtx = (*env)->GetMethodID(env, actCls,
		"getApplicationContext", "()Landroid/content/Context;");
	jobject ctx = (*env)->CallObjectMethod(env, act->clazz, getCtx);

	jclass toastCls = (*env)->FindClass(env, "android/widget/Toast");
	jmethodID makeText = (*env)->GetStaticMethodID(env, toastCls, "makeText",
		"(Landroid/content/Context;Ljava/lang/CharSequence;I)Landroid/widget/Toast;");
	jmethodID show = (*env)->GetMethodID(env, toastCls, "show", "()V");

	jstring jmsg = (*env)->NewStringUTF(env, msg);
	jobject toast = (*env)->CallStaticObjectMethod(env, toastCls, makeText,
		ctx, jmsg, (jint)1);
	(*env)->CallVoidMethod(env, toast, show);

	(*env)->DeleteLocalRef(env, toast);
	(*env)->DeleteLocalRef(env, jmsg);
	(*env)->DeleteLocalRef(env, toastCls);
	(*env)->DeleteLocalRef(env, ctx);
	(*env)->DeleteLocalRef(env, actCls);
}

static ASurfaceTexture* createSurfaceTexture(
	ANativeActivity* act, int texName, int width, int height
) {
	JNIEnv* env = act->env;
	jclass cls = (*env)->FindClass(env, "android/graphics/SurfaceTexture");
	jmethodID ctor = (*env)->GetMethodID(env, cls, "<init>", "(I)V");
	jobject jst = (*env)->NewObject(env, cls, ctor, (jint)texName);

	jmethodID setSize = (*env)->GetMethodID(env, cls,
		"setDefaultBufferSize", "(II)V");
	(*env)->CallVoidMethod(env, jst, setSize, (jint)width, (jint)height);

	ASurfaceTexture* nst = ASurfaceTexture_fromSurfaceTexture(env, jst);
	(*env)->DeleteLocalRef(env, jst);
	(*env)->DeleteLocalRef(env, cls);
	return nst;
}

static void fillWindowColor(ANativeWindow* win, uint32_t color) {
	ANativeWindow_setBuffersGeometry(win, 0, 0, WINDOW_FORMAT_RGBA_8888);
	ANativeWindow_Buffer buf;
	if (ANativeWindow_lock(win, &buf, NULL) != 0) return;
	uint32_t* px = (uint32_t*)buf.bits;
	for (int y = 0; y < buf.height; y++) {
		for (int x = 0; x < buf.width; x++) {
			px[y * buf.stride + x] = color;
		}
	}
	ANativeWindow_unlockAndPost(win);
}
*/
import "C"

import "unsafe"

// HasPermission checks whether the activity has the given permission granted.
func HasPermission(activityPtr unsafe.Pointer, permission string) bool {
	cperm := C.CString(permission)
	defer C.free(unsafe.Pointer(cperm))
	return C.hasPermission((*C.ANativeActivity)(activityPtr), cperm) != 0
}

// RequestPermission shows the system permission dialog for the given permission.
func RequestPermission(activityPtr unsafe.Pointer, permission string) {
	cperm := C.CString(permission)
	defer C.free(unsafe.Pointer(cperm))
	C.requestPermission((*C.ANativeActivity)(activityPtr), cperm)
}

// ShowToast displays an Android Toast message.
func ShowToast(activityPtr unsafe.Pointer, message string) {
	cmsg := C.CString(message)
	defer C.free(unsafe.Pointer(cmsg))
	C.showToast((*C.ANativeActivity)(activityPtr), cmsg)
}

// CreateSurfaceTexture creates a Java SurfaceTexture wrapping the given GL
// texture name and returns the native ASurfaceTexture handle as unsafe.Pointer.
func CreateSurfaceTexture(
	activityPtr unsafe.Pointer,
	texName, width, height int,
) unsafe.Pointer {
	return unsafe.Pointer(C.createSurfaceTexture(
		(*C.ANativeActivity)(activityPtr),
		C.int(texName), C.int(width), C.int(height),
	))
}

// FillWindowColor fills an ANativeWindow with a solid RGBA color.
func FillWindowColor(windowPtr unsafe.Pointer, color uint32) {
	C.fillWindowColor((*C.ANativeWindow)(windowPtr), C.uint32_t(color))
}
