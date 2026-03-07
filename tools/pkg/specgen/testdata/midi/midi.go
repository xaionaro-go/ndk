// Simulates c-for-go output for Android MIDI (amidi/AMidi.h).
// This file is parsed at AST level only; it does not compile.
package midi

import "unsafe"

// Opaque handle types.
type AMidiDevice C.AMidiDevice
type AMidiInputPort C.AMidiInputPort
type AMidiOutputPort C.AMidiOutputPort

// Integer typedefs.
type Midi_result_t int32

// Result codes.
const (
	AMIDI_OK               Midi_result_t = 0
	AMIDI_ERROR_UNKNOWN    Midi_result_t = -1
	AMIDI_ERROR_WOULD_BLOCK Midi_result_t = -2
)

// --- Device functions ---
func AMidiDevice_release(device *AMidiDevice)                                                              {}
func AMidiDevice_getNumInputPorts(device *AMidiDevice) int32                                               { return 0 }
func AMidiDevice_getNumOutputPorts(device *AMidiDevice) int32                                              { return 0 }

// --- InputPort functions ---
func AMidiInputPort_open(device *AMidiDevice, portNumber int32, outPort **AMidiInputPort) int32 { return 0 }
func AMidiInputPort_close(port *AMidiInputPort)                                                 {}
func AMidiInputPort_send(port *AMidiInputPort, buffer *byte, numBytes int32) int32              { return 0 }

// --- OutputPort functions ---
func AMidiOutputPort_open(device *AMidiDevice, portNumber int32, outPort **AMidiOutputPort) int32 { return 0 }
func AMidiOutputPort_close(port *AMidiOutputPort)                                                 {}

var _ = unsafe.Pointer(nil)
