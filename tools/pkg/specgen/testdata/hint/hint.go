// Simulates c-for-go output for Android performance hint API.
// This file is parsed at AST level only; it does not compile.
package hint

import "unsafe"

// Opaque handle types.
type APerformanceHintManager C.APerformanceHintManager
type APerformanceHintSession C.APerformanceHintSession

// --- Manager functions ---
func APerformanceHint_getManager() *APerformanceHintManager                                                                                                 { return nil }
func APerformanceHint_createSession(manager *APerformanceHintManager, tids *int32, size int32, initialTargetWorkDurationNanos int64) *APerformanceHintSession { return nil }
func APerformanceHint_getPreferredUpdateRateNanos(manager *APerformanceHintManager) int64                                                                   { return 0 }

// --- Session functions ---
func APerformanceHint_updateTargetWorkDuration(session *APerformanceHintSession, targetDurationNanos int64) int32 { return 0 }
func APerformanceHint_reportActualWorkDuration(session *APerformanceHintSession, actualDurationNanos int64) int32 { return 0 }
func APerformanceHint_closeSession(session *APerformanceHintSession)                                              {}

var _ = unsafe.Pointer(nil)
