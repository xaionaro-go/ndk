// Performance Hint API overview.
//
// Documents the Android Performance Hint API using the ndk hint package.
// The Performance Hint API allows applications to communicate their
// performance requirements to the system so that it can adjust CPU frequency
// and scheduling in real time. This is the primary mechanism for achieving
// consistent frame timing without over- or under-clocking the CPU.
//
// Core concept:
//
//	The application creates a hint session with a set of thread IDs (the
//	threads doing frame work) and a target work duration. Each frame, the
//	application reports how long the work actually took. The system uses
//	the ratio of actual-to-target duration to adjust CPU frequency:
//
//	  actual < target  ->  system may lower CPU frequency to save power
//	  actual ~ target  ->  CPU frequency is well-matched
//	  actual > target  ->  system may raise CPU frequency to meet deadlines
//
// Typical game loop integration:
//
//  1. Obtain a Manager via APerformanceHintManager_fromContext (JNI).
//  2. Gather thread IDs for your render and game-logic threads.
//  3. Create a Session with those thread IDs and a target duration
//     (e.g., 16_666_666 ns for 60 fps).
//  4. Each frame, measure the wall-clock work time and call
//     ReportActualWorkDuration.
//  5. If the frame rate target changes (e.g., 60 fps -> 90 fps),
//     call UpdateTargetWorkDuration with the new target.
//  6. Close the Session when the render loop exits.
//
// The Manager handle is obtained from the Android system through JNI:
//
//	APerformanceHintManager *mgr = APerformanceHintManager_fromContext(env, context);
//
// There is no Go-side constructor for Manager because the handle must come
// from the Android runtime. In a real application you would pass the pointer
// across the JNI boundary and wrap it in a hint.Manager.
//
// Prerequisites:
//   - Android device with API level 33+ (ADPF performance hint support).
//   - A valid APerformanceHintManager handle obtained via JNI.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AndroidGoLab/ndk/hint"
)

func main() {
	// ---------------------------------------------------------------
	// Manager
	//
	// The Manager is the entry point to the Performance Hint API. On
	// Android it is obtained by calling the JNI function:
	//
	//   APerformanceHintManager_fromContext(env, context)
	//
	// The returned pointer is then wrapped in a hint.Manager on the
	// Go side. There is no Go-side constructor because the handle is
	// owned by the Android runtime.
	//
	// For this example we use a nil Manager to illustrate the API
	// shape. On a real device you would have a valid pointer.
	// ---------------------------------------------------------------
	var mgr *hint.Manager // obtained from APerformanceHintManager_fromContext via JNI

	log.Println("manager handle would be obtained from APerformanceHintManager_fromContext")

	// ---------------------------------------------------------------
	// PreferredUpdateRateNanos
	//
	// The system reports how frequently it wants the application to
	// call ReportActualWorkDuration. Typical values:
	//
	//   ~2_000_000 ns  (2 ms)   - report every frame at 60 fps
	//   ~4_000_000 ns  (4 ms)   - report every other frame
	//
	// Reporting more frequently than this wastes CPU; reporting less
	// frequently makes the system slower to react.
	// ---------------------------------------------------------------
	//
	// preferredRate := mgr.PreferredUpdateRateNanos()
	// log.Printf("preferred update rate: %d ns (%.1f ms)",
	//     preferredRate, float64(preferredRate)/1e6)

	log.Println("preferred update rate would indicate how often to report durations")

	// ---------------------------------------------------------------
	// CreateSession
	//
	// A Session binds a set of threads to a target work duration.
	// The thread IDs identify which OS threads the system should
	// tune. Typically these are the render thread and any worker
	// threads that contribute to frame production.
	//
	// Parameters:
	//   tids     - pointer to the first element of a thread ID array
	//   size     - number of thread IDs
	//   initial  - target work duration in nanoseconds
	//
	// For 60 fps the target is ~16.67 ms = 16_666_666 ns.
	// For 90 fps the target is ~11.11 ms = 11_111_111 ns.
	// For 120 fps the target is ~8.33 ms =  8_333_333 ns.
	// ---------------------------------------------------------------
	const (
		targetFPS        = 60
		targetDurationNs = time.Second / targetFPS // ~16.67 ms
	)

	// In a real application, gather thread IDs from your render and
	// worker threads. On Android you can use gettid() from each
	// thread to obtain its kernel thread ID.
	tids := []int32{
		1234, // example: render thread tid
		1235, // example: physics worker tid
	}

	log.Printf("target: %d fps -> %v per frame", targetFPS, targetDurationNs)
	log.Printf("thread IDs: %v", tids)

	// Create the session. On a real device this returns a valid
	// Session handle. With a nil Manager this would crash, so we
	// guard it and show the call pattern instead.
	if mgr != nil {
		session := mgr.CreateSession(&tids[0], uint64(len(tids)), targetDurationNs)
		defer session.Close()

		runFrameLoop(session, targetDurationNs)
	} else {
		log.Println("skipping session creation (no Manager handle)")
		log.Println("on a real device the call would be:")
		log.Println("  session := mgr.CreateSession(&tids[0], uint64(len(tids)), targetDurationNs)")
	}

	// ---------------------------------------------------------------
	// Frame loop pattern (reference)
	//
	// Each frame:
	//   1. Record the start time.
	//   2. Do rendering / simulation work.
	//   3. Record the end time.
	//   4. Call session.ReportActualWorkDuration(elapsed).
	//
	// The system uses the stream of reports to adjust CPU frequency.
	// If work consistently finishes early, the CPU frequency is
	// lowered to save power. If work consistently overruns, the
	// frequency is raised to help meet the deadline.
	// ---------------------------------------------------------------

	// Print a summary of the API surface for reference.
	fmt.Println()
	fmt.Println("Performance Hint API summary:")
	fmt.Println()
	fmt.Println("  hint.Manager (from APerformanceHintManager_fromContext via JNI)")
	fmt.Println("    .CreateSession(tids, size, targetNs) *Session")
	fmt.Println("    .PreferredUpdateRateNanos() int64")
	fmt.Println()
	fmt.Println("  hint.Session")
	fmt.Println("    .ReportActualWorkDuration(actualNs) error")
	fmt.Println("    .UpdateTargetWorkDuration(targetNs) error")
	fmt.Println("    .Close() error")
}

// runFrameLoop demonstrates the per-frame reporting pattern. It simulates
// a game loop that measures work duration and reports it to the performance
// hint session so the system can tune CPU frequency.
func runFrameLoop(session *hint.Session, target time.Duration) {
	const totalFrames = 180 // 3 seconds at 60 fps

	log.Printf("starting frame loop: %d frames, target %v", totalFrames, target)

	for frame := 0; frame < totalFrames; frame++ {
		start := time.Now()

		// Simulate frame work. In a real engine this would be the
		// render and simulation pass.
		doFrameWork(frame)

		elapsed := time.Since(start)

		// Report how long the frame actually took. The system
		// compares this against the target to decide whether to
		// adjust CPU frequency up or down.
		if err := session.ReportActualWorkDuration(elapsed); err != nil {
			log.Printf("frame %d: report failed: %v", frame, err)
		}

		if frame%60 == 0 {
			fmt.Printf("  [frame %3d] actual=%5.2f ms  target=%5.2f ms  ratio=%.2f\n",
				frame,
				float64(elapsed.Nanoseconds())/1e6,
				float64(target.Nanoseconds())/1e6,
				float64(elapsed)/float64(target))
		}
	}

	// If the application switches frame rate targets (e.g., the user
	// changes quality settings from 60 fps to 90 fps), update the
	// session target rather than creating a new session.
	const newTarget = time.Second / 90 // ~11.11 ms for 90 fps
	if err := session.UpdateTargetWorkDuration(newTarget); err != nil {
		log.Printf("update target failed: %v", err)
	} else {
		log.Printf("target updated to %v (90 fps)", newTarget)
	}

	log.Println("frame loop finished")
}

// doFrameWork simulates variable-duration frame work. Real engines would
// perform scene traversal, draw-call submission, and GPU synchronization.
func doFrameWork(frame int) {
	// Vary the simulated work to show how actual durations fluctuate.
	base := 14 * time.Millisecond
	jitter := time.Duration(frame%5) * time.Millisecond
	time.Sleep(base + jitter)
}
