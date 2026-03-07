package simple

import "unsafe"

// Opaque handle types (c-for-go pattern: type alias to C struct pointer).
type FakeStream C.FakeStream
type FakeBuilder C.FakeBuilder

// Integer typedefs.
type Fake_result_t int32
type Fake_direction_t int32

// Enum values.
const (
	FAKE_OK         Fake_result_t = 0
	FAKE_ERROR_BASE Fake_result_t = -900
)

const (
	FAKE_DIRECTION_OUTPUT Fake_direction_t = 0
	FAKE_DIRECTION_INPUT  Fake_direction_t = 1
)

// Callback type.
type FakeStream_dataCallback func(stream *FakeStream, userData unsafe.Pointer, audioData unsafe.Pointer, numFrames int32) Fake_result_t

// Functions.
func Fake_createBuilder(builder **FakeBuilder) Fake_result_t                    { return 0 }
func FakeBuilder_setDeviceId(builder *FakeBuilder, deviceId int32)              {}
func FakeBuilder_openStream(builder *FakeBuilder, stream **FakeStream) Fake_result_t { return 0 }
func FakeBuilder_delete(builder *FakeBuilder)                                   {}
func FakeStream_close(stream *FakeStream) Fake_result_t                         { return 0 }
func FakeStream_requestStart(stream *FakeStream) Fake_result_t                  { return 0 }
func FakeStream_getSampleRate(stream *FakeStream) int32                         { return 0 }

var _ = unsafe.Pointer(nil)
