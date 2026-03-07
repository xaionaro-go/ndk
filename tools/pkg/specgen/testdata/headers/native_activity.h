/*
 * Copyright (C) 2010 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 */

#ifndef ANDROID_NATIVE_ACTIVITY_H
#define ANDROID_NATIVE_ACTIVITY_H

#include <stdint.h>
#include <sys/types.h>

#ifdef __cplusplus
extern "C" {
#endif

struct ANativeActivity;
struct ANativeWindow;

/**
 * These are the callbacks the framework makes into a native application.
 * All of these callbacks happen on the main thread of the application.
 * By default, all callbacks are NULL; set to a function pointer to have
 * them called.
 */
typedef struct ANativeActivityCallbacks {
    /**
     * NativeActivity has started.  See Java documentation for Activity.onStart()
     * for more information.
     */
    void (*onStart)(ANativeActivity* activity);

    /**
     * The framework is asking NativeActivity to save its current instance state.
     * See Java documentation for Activity.onSaveInstanceState() for more
     * information. The returned pointer needs to be created with malloc();
     * the framework will call free() on it for you. You also must fill in
     * outSize with the number of bytes in the allocation. Note that the
     * saved state will be persisted, so it can not contain any active
     * entities (pointers to memory, file descriptors, etc).
     */
    void* (*onSaveInstanceState)(ANativeActivity* activity, size_t* outSize);

    /**
     * NativeActivity window focus has changed.
     */
    void (*onWindowFocusChanged)(ANativeActivity* activity, int hasFocus);

    /**
     * The drawing window for this native activity has been created.  You
     * can use the given native window object to start drawing.
     */
    void (*onNativeWindowCreated)(ANativeActivity* activity, ANativeWindow* window);

    /**
     * The rectangle in the window in which content should be placed has changed.
     */
    void (*onContentRectChanged)(ANativeActivity* activity, const ARect* rect);
} ANativeActivityCallbacks;

#ifdef __cplusplus
}
#endif

#endif // ANDROID_NATIVE_ACTIVITY_H
