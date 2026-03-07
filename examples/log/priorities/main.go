// Example: Android log priority levels and their string representations.
//
// Enumerates every Priority constant exported by ndk/log and prints
// its name, integer value, and the equivalent Android logcat letter.
//
// The NDK function __android_log_write(int prio, const char *tag,
// const char *text) accepts these priority values as its first argument.
// In Go code using ndk you would pass androidlog.Debug, androidlog.Info,
// etc. to any API that requires a log priority.
//
// Priority levels from lowest to highest:
//
//	Unknown (0) – should not normally appear in output
//	Default (1) – the default priority; logcat maps it to the minimum
//	Verbose (2) – extremely detailed, usually disabled in production
//	Debug   (3) – development diagnostics
//	Info    (4) – normal operational messages
//	Warn    (5) – unexpected but recoverable situations
//	Error   (6) – failures that need attention
//	Fatal   (7) – unrecoverable errors; typically followed by abort
//	Silent  (8) – suppresses all output; used only as a filter, never for writing
//
// This program runs on the host (no Android device required) because
// it only inspects the constant values without calling into the NDK.
package main

import (
	"fmt"

	androidlog "github.com/xaionaro-go/ndk/log"
)

// logcatLetter returns the single-character abbreviation that logcat
// uses for each priority level (e.g. "V" for Verbose, "D" for Debug).
func logcatLetter(p androidlog.Priority) string {
	switch p {
	case androidlog.Verbose:
		return "V"
	case androidlog.Debug:
		return "D"
	case androidlog.Info:
		return "I"
	case androidlog.Warn:
		return "W"
	case androidlog.Error:
		return "E"
	case androidlog.Fatal:
		return "F"
	case androidlog.Silent:
		return "S"
	default:
		return "?"
	}
}

func main() {
	// Every priority constant in ascending order.
	priorities := []androidlog.Priority{
		androidlog.Unknown,
		androidlog.Default,
		androidlog.Verbose,
		androidlog.Debug,
		androidlog.Info,
		androidlog.Warn,
		androidlog.Error,
		androidlog.Fatal,
		androidlog.Silent,
	}

	fmt.Println("Android log priority levels:")
	fmt.Println()
	fmt.Printf("  %-10s  %5s  %s\n", "Name", "Value", "Logcat")
	fmt.Printf("  %-10s  %5s  %s\n", "----------", "-----", "------")

	for _, p := range priorities {
		fmt.Printf("  %-10s  %5d  %s\n", p.String(), int(p), logcatLetter(p))
	}

	// Demonstrate that the String() method handles out-of-range values.
	fmt.Println()
	fmt.Printf("  Out-of-range value: %s\n", androidlog.Priority(42))

	// Show typical usage: choosing a minimum log level for filtering.
	//
	// On Android, logcat filters by priority. For example, setting the
	// minimum to Warn suppresses Verbose, Debug, and Info messages:
	//
	//   adb logcat *:W
	//
	// In application code the same comparison decides whether to emit
	// a message:
	//
	//   if msgPriority >= minPriority { /* write to log */ }
	fmt.Println()
	minLevel := androidlog.Warn
	fmt.Printf("  Filter example: minimum level = %s (%d)\n", minLevel, int(minLevel))

	for _, p := range priorities {
		if p >= minLevel && p != androidlog.Silent {
			fmt.Printf("    %s would be printed\n", p)
		}
	}
}
