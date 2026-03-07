/*
 * Copyright (C) 2015 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 */

#ifndef _NDK_CAMERA_DEVICE_H
#define _NDK_CAMERA_DEVICE_H

#include <sys/cdefs.h>

__BEGIN_DECLS

/**
 * ACameraDevice is opaque type that provides access to a camera device.
 */
typedef struct ACameraDevice ACameraDevice;

/**
 * Camera device state callbacks to be used in {@link ACameraDevice_StateCallbacks}.
 *
 * @param context The optional application context provided by the user when
 *                the callbacks were registered.
 * @param device The {@link ACameraDevice} that is being disconnected.
 */
typedef void (*ACameraDevice_StateCallback)(void* context, ACameraDevice* device);

/**
 * Camera device error state callbacks to be used in {@link ACameraDevice_StateCallbacks}.
 *
 * @param context The optional application context provided by the user when
 *                the callbacks were registered.
 * @param device The {@link ACameraDevice} that is in error state.
 * @param error The error code associated with this error state.
 *              One of {@link ERROR_CAMERA_IN_USE}, {@link ERROR_MAX_CAMERAS_IN_USE},
 *              {@link ERROR_CAMERA_DISABLED}, {@link ERROR_CAMERA_DEVICE},
 *              {@link ERROR_CAMERA_SERVICE}.
 */
typedef void (*ACameraDevice_ErrorStateCallback)(void* context, ACameraDevice* device, int error);

/**
 * A set of callbacks for monitoring events on camera device.
 *
 * <p>When a callback is invoked, the calling code should return as quickly as possible
 * and avoid blocking or running expensive operations.</p>
 */
typedef struct ACameraDevice_StateCallbacks {
    /// optional application context.
    void*                             context;

    /**
     * Called when a camera device is no longer available for use.
     *
     * <p>Any attempt to call API methods on this ACameraDevice will return
     * {@link ACAMERA_ERROR_CAMERA_DISCONNECTED}. The disconnection could be
     * due to a change in security policy or permissions.</p>
     */
    ACameraDevice_StateCallback       onDisconnected;

    /**
     * Called when a camera device has encountered a serious error.
     *
     * <p>This callback indicates that the device can no longer be used.</p>
     */
    ACameraDevice_ErrorStateCallback  onError;
} ACameraDevice_StateCallbacks;

/**
 * ACameraIdList represents a list of camera device Ids.
 */
typedef struct ACameraIdList {
    int numCameras;          ///< Number of camera device Ids
    const char** cameraIds;  ///< list of camera device Ids
} ACameraIdList;

__END_DECLS

#endif /* _NDK_CAMERA_DEVICE_H */
