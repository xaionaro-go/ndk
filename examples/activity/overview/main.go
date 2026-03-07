// NativeActivity lifecycle overview.
//
// Documents the Android NativeActivity lifecycle and the methods provided by
// the ndk/activity package. NativeActivity is the entry point for purely
// native Android applications: the system creates an ANativeActivity handle
// and delivers it to the native code via the ANativeActivity_onCreate
// callback. All Activity methods require this system-provided handle.
//
// Because the Activity handle is only available inside a running Android
// NativeActivity process, this example documents the pattern and prints the
// available constants rather than calling the methods directly.
//
// Lifecycle summary (system-driven):
//
//	ANativeActivity_onCreate         <- system creates the activity
//	  onStart
//	  onResume
//	  onNativeWindowCreated          <- ANativeWindow becomes available
//	  onInputQueueCreated            <- AInputQueue becomes available
//	    ... application runs ...
//	  onInputQueueDestroyed          <- AInputQueue revoked
//	  onNativeWindowDestroyed        <- ANativeWindow revoked
//	  onPause
//	  onStop
//	  onDestroy                      <- activity is torn down
//
// The Activity handle is valid from onCreate through onDestroy. The
// ANativeWindow and AInputQueue handles are valid only between their
// respective created/destroyed callback pairs.
//
// This program must run on an Android device.
package main

import (
	"fmt"

	"github.com/xaionaro-go/ndk/activity"
)

func main() {
	// ---------------------------------------------------------------
	// Activity handle
	//
	// The Activity struct wraps the NDK ANativeActivity pointer. It is
	// NOT created from Go code. The Android framework passes the
	// ANativeActivity pointer to the native entry point:
	//
	//   void ANativeActivity_onCreate(
	//       ANativeActivity* nativeActivity,
	//       void*            savedState,
	//       size_t           savedStateSize);
	//
	// The ndk runtime wraps this pointer into an Activity value and
	// delivers it through the lifecycle callbacks.
	// ---------------------------------------------------------------

	fmt.Println("=== ndk/activity API overview ===")
	fmt.Println()

	// ---------------------------------------------------------------
	// Activity methods
	// ---------------------------------------------------------------

	fmt.Println("Activity methods:")
	fmt.Println()
	fmt.Println("  Finish()")
	fmt.Println("      Request that the activity be finished and closed.")
	fmt.Println("      Equivalent to calling Activity.finish() from Java.")
	fmt.Println()
	fmt.Println("  ShowSoftInput(flags uint32)")
	fmt.Println("      Show the soft (on-screen) keyboard.")
	fmt.Println("      Pass ShowSoftInputFlags values to control behavior.")
	fmt.Println()
	fmt.Println("  HideSoftInput(flags uint32)")
	fmt.Println("      Hide the soft keyboard.")
	fmt.Println("      Pass HideSoftInputFlags values to control behavior.")
	fmt.Println()
	fmt.Println("  SetWindowFlags(addFlags, removeFlags uint32)")
	fmt.Println("      Add and/or remove window flags (e.g. fullscreen,")
	fmt.Println("      keep-screen-on). These correspond to the")
	fmt.Println("      WindowManager.LayoutParams FLAG_* constants.")
	fmt.Println()
	fmt.Println("  SetWindowFormat(format int32)")
	fmt.Println("      Set the pixel format of the activity's window.")
	fmt.Println("      Common formats are defined in the Android NDK")
	fmt.Println("      window.h header (AHARDWAREBUFFER_FORMAT_*).")
	fmt.Println()

	// ---------------------------------------------------------------
	// ShowSoftInputFlags
	// ---------------------------------------------------------------

	fmt.Println("ShowSoftInputFlags constants:")
	fmt.Printf("  Implicit = %d  -- show only if the user has not explicitly dismissed it\n",
		activity.Implicit)
	fmt.Printf("  Forced   = %d  -- force the keyboard visible even if user dismissed it\n",
		activity.Forced)
	fmt.Println()

	// ---------------------------------------------------------------
	// HideSoftInputFlags
	// ---------------------------------------------------------------

	fmt.Println("HideSoftInputFlags constants:")
	fmt.Printf("  ImplicitOnly = %d  -- hide only if it was shown implicitly (not by the user)\n",
		activity.ImplicitOnly)
	fmt.Printf("  NotAlways    = %d  -- hide unless the user explicitly requested it\n",
		activity.NotAlways)
	fmt.Println()

	// ---------------------------------------------------------------
	// Associated types
	// ---------------------------------------------------------------

	fmt.Println("Associated types:")
	fmt.Println()
	fmt.Println("  AInputQueue")
	fmt.Println("      Wraps the NDK AInputQueue handle. Delivered via the")
	fmt.Println("      onInputQueueCreated / onInputQueueDestroyed callbacks.")
	fmt.Println("      Used to read touch, key, and motion events.")
	fmt.Println()
	fmt.Println("  ANativeWindow")
	fmt.Println("      Wraps the NDK ANativeWindow handle. Delivered via the")
	fmt.Println("      onNativeWindowCreated / onNativeWindowDestroyed callbacks.")
	fmt.Println("      Used as the rendering surface for EGL, Vulkan, or direct")
	fmt.Println("      pixel buffer access.")
	fmt.Println()
	fmt.Println("  ARect")
	fmt.Println("      Wraps the NDK ARect struct (left, top, right, bottom).")
	fmt.Println("      Used in content-rect-changed callbacks to report the")
	fmt.Println("      visible area of the activity window.")
	fmt.Println()

	// ---------------------------------------------------------------
	// Typical usage pattern
	// ---------------------------------------------------------------

	fmt.Println("Typical usage pattern (pseudocode):")
	fmt.Println()
	fmt.Println("  func onCreated(act *activity.Activity) {")
	fmt.Println("      // Store the activity handle for later use.")
	fmt.Println("      // Set fullscreen flags:")
	fmt.Println("      //   act.SetWindowFlags(FLAG_FULLSCREEN, 0)")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  func onTextInput(act *activity.Activity) {")
	fmt.Println("      // Show the keyboard for text entry:")
	fmt.Println("      //   act.ShowSoftInput(uint32(activity.Implicit))")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  func onTextDone(act *activity.Activity) {")
	fmt.Println("      // Dismiss the keyboard:")
	fmt.Println("      //   act.HideSoftInput(uint32(activity.NotAlways))")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  func onQuit(act *activity.Activity) {")
	fmt.Println("      // Finish the activity:")
	fmt.Println("      //   act.Finish()")
	fmt.Println("  }")
}
