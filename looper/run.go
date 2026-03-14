package looper

import "runtime"

const prepareAllowNonCallbacks = int32(1)

// Run locks the current goroutine to its OS thread, prepares a new looper,
// and calls fn with it. The looper is closed when fn returns.
func Run(fn func(*Looper)) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lp := Prepare(prepareAllowNonCallbacks)
	defer lp.Close()

	fn(lp)
}
