package edgecases

type FakeStream struct{}

// Untyped const block (should be skipped).
const (
	untypedConst = 42
)

// String const (non-integer, should be skipped).
const (
	stringConst = "hello"
)

// Function with no params and multiple returns (multi-return not captured).
func Edge_multiReturn() (int32, bool) { return 0, false }

// Unexported function (should be skipped).
func unexported() {}

// Function with unnamed params.
func Edge_unnamedParam(int32) {}

// c-for-go ref-helper function (should be filtered).
func NewFakeStreamRef(ref *FakeStream) *FakeStream { return ref }
