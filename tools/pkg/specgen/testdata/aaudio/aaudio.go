// Simulates c-for-go output for Android AAudio.
// This file is parsed at AST level only; it does not compile.
package aaudio

import "unsafe"

// Opaque handle types.
type AAudioStreamBuilder C.AAudioStreamBuilder
type AAudioStream C.AAudioStream

// Integer typedefs.
type Aaudio_result_t int32
type Aaudio_direction_t int32
type Aaudio_format_t int32
type Aaudio_sharing_mode_t int32
type Aaudio_performance_mode_t int32
type Aaudio_stream_state_t int32
type Aaudio_data_callback_result_t int32

// Result codes.
const (
	AAUDIO_OK                        Aaudio_result_t = 0
	AAUDIO_ERROR_BASE                Aaudio_result_t = -900
	AAUDIO_ERROR_DISCONNECTED        Aaudio_result_t = -899
	AAUDIO_ERROR_ILLEGAL_ARGUMENT    Aaudio_result_t = -898
	AAUDIO_ERROR_INTERNAL            Aaudio_result_t = -896
	AAUDIO_ERROR_INVALID_STATE       Aaudio_result_t = -895
	AAUDIO_ERROR_INVALID_HANDLE      Aaudio_result_t = -892
	AAUDIO_ERROR_UNIMPLEMENTED       Aaudio_result_t = -890
	AAUDIO_ERROR_UNAVAILABLE         Aaudio_result_t = -889
	AAUDIO_ERROR_NO_FREE_HANDLES     Aaudio_result_t = -888
	AAUDIO_ERROR_NO_MEMORY           Aaudio_result_t = -887
	AAUDIO_ERROR_NULL                Aaudio_result_t = -886
	AAUDIO_ERROR_TIMEOUT             Aaudio_result_t = -885
	AAUDIO_ERROR_WOULD_BLOCK         Aaudio_result_t = -884
	AAUDIO_ERROR_INVALID_FORMAT      Aaudio_result_t = -883
	AAUDIO_ERROR_OUT_OF_RANGE        Aaudio_result_t = -882
	AAUDIO_ERROR_NO_SERVICE          Aaudio_result_t = -881
	AAUDIO_ERROR_INVALID_RATE        Aaudio_result_t = -880
)

// Direction enum.
const (
	AAUDIO_DIRECTION_OUTPUT Aaudio_direction_t = 0
	AAUDIO_DIRECTION_INPUT  Aaudio_direction_t = 1
)

// Format enum.
const (
	AAUDIO_FORMAT_INVALID   Aaudio_format_t = -1
	AAUDIO_FORMAT_UNSPECIFIED Aaudio_format_t = 0
	AAUDIO_FORMAT_PCM_I16   Aaudio_format_t = 1
	AAUDIO_FORMAT_PCM_FLOAT Aaudio_format_t = 2
)

// Sharing mode enum.
const (
	AAUDIO_SHARING_MODE_EXCLUSIVE Aaudio_sharing_mode_t = 0
	AAUDIO_SHARING_MODE_SHARED    Aaudio_sharing_mode_t = 1
)

// Performance mode enum.
const (
	AAUDIO_PERFORMANCE_MODE_NONE          Aaudio_performance_mode_t = 10
	AAUDIO_PERFORMANCE_MODE_POWER_SAVING  Aaudio_performance_mode_t = 11
	AAUDIO_PERFORMANCE_MODE_LOW_LATENCY   Aaudio_performance_mode_t = 12
)

// Stream state enum.
const (
	AAUDIO_STREAM_STATE_UNINITIALIZED Aaudio_stream_state_t = 0
	AAUDIO_STREAM_STATE_UNKNOWN       Aaudio_stream_state_t = 1
	AAUDIO_STREAM_STATE_OPEN          Aaudio_stream_state_t = 2
	AAUDIO_STREAM_STATE_STARTING      Aaudio_stream_state_t = 3
	AAUDIO_STREAM_STATE_STARTED       Aaudio_stream_state_t = 4
	AAUDIO_STREAM_STATE_PAUSING       Aaudio_stream_state_t = 5
	AAUDIO_STREAM_STATE_PAUSED        Aaudio_stream_state_t = 6
	AAUDIO_STREAM_STATE_FLUSHING      Aaudio_stream_state_t = 7
	AAUDIO_STREAM_STATE_FLUSHED       Aaudio_stream_state_t = 8
	AAUDIO_STREAM_STATE_STOPPING      Aaudio_stream_state_t = 9
	AAUDIO_STREAM_STATE_STOPPED       Aaudio_stream_state_t = 10
	AAUDIO_STREAM_STATE_CLOSING       Aaudio_stream_state_t = 11
	AAUDIO_STREAM_STATE_CLOSED        Aaudio_stream_state_t = 12
	AAUDIO_STREAM_STATE_DISCONNECTED  Aaudio_stream_state_t = 13
)

// Data callback result enum.
const (
	AAUDIO_CALLBACK_RESULT_CONTINUE Aaudio_data_callback_result_t = 0
	AAUDIO_CALLBACK_RESULT_STOP     Aaudio_data_callback_result_t = 1
)

// Callback type.
type AAudioStream_dataCallback func(stream *AAudioStream, userData unsafe.Pointer, audioData unsafe.Pointer, numFrames int32) Aaudio_data_callback_result_t

// --- StreamBuilder functions ---
func AAudio_createStreamBuilder(builder **AAudioStreamBuilder) Aaudio_result_t { return 0 }
func AAudioStreamBuilder_setDeviceId(builder *AAudioStreamBuilder, deviceId int32) {}
func AAudioStreamBuilder_setDirection(builder *AAudioStreamBuilder, direction Aaudio_direction_t) {}
func AAudioStreamBuilder_setSharingMode(builder *AAudioStreamBuilder, sharingMode Aaudio_sharing_mode_t) {}
func AAudioStreamBuilder_setPerformanceMode(builder *AAudioStreamBuilder, mode Aaudio_performance_mode_t) {}
func AAudioStreamBuilder_setSampleRate(builder *AAudioStreamBuilder, sampleRate int32) {}
func AAudioStreamBuilder_setChannelCount(builder *AAudioStreamBuilder, channelCount int32) {}
func AAudioStreamBuilder_setFormat(builder *AAudioStreamBuilder, format Aaudio_format_t) {}
func AAudioStreamBuilder_setBufferCapacityInFrames(builder *AAudioStreamBuilder, numFrames int32) {}
func AAudioStreamBuilder_setDataCallback(builder *AAudioStreamBuilder, callback AAudioStream_dataCallback, userData unsafe.Pointer) {}
func AAudioStreamBuilder_openStream(builder *AAudioStreamBuilder, stream **AAudioStream) Aaudio_result_t { return 0 }
func AAudioStreamBuilder_delete(builder *AAudioStreamBuilder) {}

// --- Stream functions ---
func AAudioStream_close(stream *AAudioStream) Aaudio_result_t   { return 0 }
func AAudioStream_getState(stream *AAudioStream) Aaudio_stream_state_t { return 0 }
func AAudioStream_getFramesPerBurst(stream *AAudioStream) int32 { return 0 }
func AAudioStream_getSampleRate(stream *AAudioStream) int32     { return 0 }
func AAudioStream_getChannelCount(stream *AAudioStream) int32   { return 0 }
func AAudioStream_getXRunCount(stream *AAudioStream) int32      { return 0 }
func AAudioStream_requestStart(stream *AAudioStream) Aaudio_result_t { return 0 }
func AAudioStream_requestPause(stream *AAudioStream) Aaudio_result_t { return 0 }
func AAudioStream_requestStop(stream *AAudioStream) Aaudio_result_t  { return 0 }
func AAudioStream_requestFlush(stream *AAudioStream) Aaudio_result_t { return 0 }
func AAudioStream_write(stream *AAudioStream, audioData unsafe.Pointer, numFrames int32, timeoutNanos int64) Aaudio_result_t { return 0 }

var _ = unsafe.Pointer(nil)
