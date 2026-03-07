// Simulates c-for-go output for Android trace.
// This file is parsed at AST level only; it does not compile.
package trace

// --- Trace functions ---
func ATrace_isEnabled() bool                                      { return false }
func ATrace_beginSection(sectionName *byte)                       {}
func ATrace_endSection()                                          {}
func ATrace_beginAsyncSection(sectionName *byte, cookie int32)    {}
func ATrace_endAsyncSection(sectionName *byte, cookie int32)      {}
func ATrace_setCounter(counterName *byte, counterValue int64)     {}
