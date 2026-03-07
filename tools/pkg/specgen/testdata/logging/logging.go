// Simulates c-for-go output for Android logging.
// This file is parsed at AST level only; it does not compile.
package logging

// Log priority enum.
type Android_LogPriority int32

const (
	ANDROID_LOG_UNKNOWN Android_LogPriority = 0
	ANDROID_LOG_DEFAULT Android_LogPriority = 1
	ANDROID_LOG_VERBOSE Android_LogPriority = 2
	ANDROID_LOG_DEBUG   Android_LogPriority = 3
	ANDROID_LOG_INFO    Android_LogPriority = 4
	ANDROID_LOG_WARN    Android_LogPriority = 5
	ANDROID_LOG_ERROR   Android_LogPriority = 6
	ANDROID_LOG_FATAL   Android_LogPriority = 7
	ANDROID_LOG_SILENT  Android_LogPriority = 8
)

// --- Log functions ---
func __android_log_write(prio int32, tag *byte, text *byte) int32 { return 0 }
func __android_log_print(prio int32, tag *byte, fmt *byte) int32  { return 0 }
