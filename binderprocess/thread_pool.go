// Package binderprocess provides access to ABinderProcess_* functions
// that are not part of the public NDK headers but available at runtime
// via dlsym in libbinder_ndk.so.
package binderprocess

/*
#include <dlfcn.h>
#include <stdint.h>

static int thread_pool_start(unsigned thread_pool_size)
{
	typedef int (*set_thread_pool_max_fn)(uint32_t);
	typedef void (*start_thread_pool_fn)(void);

	void *h = dlopen("libbinder_ndk.so", RTLD_NOW | RTLD_LOCAL);
	if (h == NULL)
		return 1;

	start_thread_pool_fn start_thread_pool =
		(start_thread_pool_fn)dlsym(h, "ABinderProcess_startThreadPool");
	if (start_thread_pool == NULL)
		return 2;

	if (thread_pool_size != 0) {
		set_thread_pool_max_fn set_thread_pool_max =
			(set_thread_pool_max_fn)dlsym(h, "ABinderProcess_setThreadPoolMaxThreadCount");
		if (set_thread_pool_max == NULL)
			return 3;
		set_thread_pool_max(thread_pool_size);
	}

	start_thread_pool();
	return 0;
}
*/
import "C"

import "fmt"

// ErrDLOpenLibbinderNDK is returned when libbinder_ndk.so cannot be loaded.
type ErrDLOpenLibbinderNDK struct{}

func (ErrDLOpenLibbinderNDK) Error() string {
	return "unable to dlopen libbinder_ndk.so"
}

// ErrSymbolStartThreadPool is returned when ABinderProcess_startThreadPool is not found.
type ErrSymbolStartThreadPool struct{}

func (ErrSymbolStartThreadPool) Error() string {
	return "unable to find ABinderProcess_startThreadPool symbol"
}

// ErrSymbolSetThreadPoolMax is returned when ABinderProcess_setThreadPoolMaxThreadCount is not found.
type ErrSymbolSetThreadPoolMax struct{}

func (ErrSymbolSetThreadPoolMax) Error() string {
	return "unable to find ABinderProcess_setThreadPoolMaxThreadCount symbol"
}

// StartThreadPool starts the binder thread pool with the given number of threads.
// If threads is 0, the default thread count is used (only startThreadPool is called).
func StartThreadPool(threads int) error {
	rc := C.thread_pool_start(C.uint(threads))
	switch rc {
	case 0:
		return nil
	case 1:
		return ErrDLOpenLibbinderNDK{}
	case 2:
		return ErrSymbolStartThreadPool{}
	case 3:
		return ErrSymbolSetThreadPoolMax{}
	default:
		return fmt.Errorf("thread_pool_start returned unexpected code %d", rc)
	}
}
