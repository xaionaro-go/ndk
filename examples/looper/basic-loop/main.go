// Basic ALooper event loop example.
//
// Demonstrates how to prepare an ALooper for the current thread, wake it from
// a separate goroutine, and poll for the wake event. This is the fundamental
// pattern for any ALooper-based event loop on Android.
//
// This program must run on an Android device with NDK looper support.
package main

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/AndroidGoLab/ndk/looper"
)

func main() {
	// Lock this goroutine to its OS thread. ALooper is thread-local, so the
	// prepare and poll calls must happen on the same OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Prepare a looper for the current thread. The non-callbacks flag allows
	// PollOnce to return file descriptor events directly instead of requiring
	// a C callback.
	lp := looper.Prepare(int32(looper.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS))
	if lp == nil {
		log.Fatal("failed to prepare looper")
	}
	defer func() {
		if err := lp.Close(); err != nil {
			log.Printf("close looper: %v", err)
		}
	}()

	log.Println("looper prepared, acquiring extra reference for wake goroutine")

	// Acquire an extra reference so the wake goroutine can safely call Wake
	// while the main goroutine still owns the looper.
	lp.Acquire()

	// Launch a goroutine that wakes the looper after a short delay.
	go func() {
		time.Sleep(100 * time.Millisecond)
		log.Println("goroutine: waking looper")
		lp.Wake()
	}()

	log.Println("entering poll loop (blocking until wake)")

	// Poll with an indefinite timeout (-1). PollOnce blocks until an event
	// arrives or the looper is woken.
	var fd, events int32
	var data unsafe.Pointer

	result := looper.LOOPER_POLL(looper.PollOnce(-1, &fd, &events, &data))

	switch result {
	case looper.ALOOPER_POLL_WAKE:
		log.Println("poll returned: WAKE -- looper was woken successfully")
	case looper.ALOOPER_POLL_CALLBACK:
		log.Println("poll returned: CALLBACK")
	case looper.ALOOPER_POLL_TIMEOUT:
		log.Println("poll returned: TIMEOUT (unexpected with infinite wait)")
	case looper.ALOOPER_POLL_ERROR:
		log.Fatal("poll returned: ERROR")
	default:
		// A non-negative value is a file descriptor that has events.
		log.Printf("poll returned fd=%d events=%d", result, events)
	}

	log.Println("basic loop finished")
}
