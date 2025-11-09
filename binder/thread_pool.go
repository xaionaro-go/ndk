package binder

/*
#include <dlfcn.h>
#include <stdint.h>

static void *dlopen_libbinder_ndk(void)
{
    void *h = dlopen("libbinder_ndk.so", RTLD_NOW | RTLD_LOCAL);
    if (h != NULL)
        return h;

    return NULL;
}

int thread_pool_start(unsigned thead_pool_size)
{
    typedef int (*set_thread_pool_max_fn)(uint32_t);
    typedef void (*start_thread_pool_fn)(void);

    set_thread_pool_max_fn set_thread_pool_max = NULL;
    start_thread_pool_fn start_thread_pool = NULL;

    void *h = dlopen_libbinder_ndk();
    if (h == NULL) {
        return 1;
    }

    start_thread_pool =
        (start_thread_pool_fn) dlsym(h, "ABinderProcess_startThreadPool");

    if (start_thread_pool == NULL) {
        return 2;
    }

    if (thead_pool_size != 0) {
        set_thread_pool_max = (set_thread_pool_max_fn) dlsym(h, "ABinderProcess_setThreadPoolMaxThreadCount");
        if (set_thread_pool_max == NULL)
            return 3;
        set_thread_pool_max(thead_pool_size);
    }

    start_thread_pool();
    return 0;
}
*/
import "C"
import "fmt"

type ErrUnableToDLLoadLibbinderNDK struct{}
type ErrUnableToFindStartThreadPoolSymbol struct{}
type ErrUnableToFindSetThreadPoolMaxSymbol struct{}

func (e ErrUnableToDLLoadLibbinderNDK) Error() string {
	return "unable to dlopen libbinder_ndk.so"
}
func (e ErrUnableToFindStartThreadPoolSymbol) Error() string {
	return "unable to find ABinderProcess_startThreadPool symbol"
}
func (e ErrUnableToFindSetThreadPoolMaxSymbol) Error() string {
	return "unable to find ABinderProcess_setThreadPoolMaxThreadCount symbol"
}

func ThreadPoolStart(threads int) error {
	rc := C.thread_pool_start(C.uint(threads))
	switch rc {
	case 0:
		return nil
	case 1:
		return ErrUnableToDLLoadLibbinderNDK{}
	case 2:
		return ErrUnableToFindStartThreadPoolSymbol{}
	case 3:
		return ErrUnableToFindSetThreadPoolMaxSymbol{}
	default:
		return fmt.Errorf("unknown error code %d from thread_pool_start", rc)
	}
}
