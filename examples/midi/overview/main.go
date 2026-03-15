// MIDI API overview and reference example.
//
// Demonstrates the full device-to-port lifecycle of the ndk/midi package,
// which wraps Android's native MIDI API (AMidi). The example shows how to:
//
//   - Obtain an AMidiDevice handle from the Java MidiManager via JNI
//   - Inspect the number of input and output ports on the device
//   - Open an input port and send MIDI messages (Note On, Note Off,
//     Control Change, Program Change)
//   - Construct MIDI messages as byte arrays with the correct status bytes
//   - Close ports and the device in the correct order
//
// # Obtaining the AMidiDevice handle
//
// The Android NDK does not provide a way to discover or open MIDI devices
// from C/native code. Device discovery happens on the Java side through
// android.media.midi.MidiManager. A typical pattern is:
//
//  1. In Java/Kotlin, use MidiManager.openDevice() to get a MidiDevice.
//  2. Call AMidiDevice_fromJava(env, javaMidiDevice) via JNI to obtain the
//     native AMidiDevice pointer.
//  3. Pass that pointer to Go (e.g. via a cgo bridge or as an unsafe.Pointer
//     stored in a global).
//
// Because this example cannot run standalone (it requires a live Android
// MIDI device obtained through JNI), it prints reference information about
// MIDI message formats and demonstrates the API calls that would be made
// once a Device handle is available.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"unsafe"

	"github.com/AndroidGoLab/ndk/midi"
)

// MIDI status byte constants (channel voice messages, channel 0).
const (
	noteOff       = 0x80 // Note Off, channel 0
	noteOn        = 0x90 // Note On, channel 0
	controlChange = 0xB0 // Control Change, channel 0
	programChange = 0xC0 // Program Change, channel 0
)

// MIDI note and controller constants used in the examples below.
const (
	middleC     = 60  // MIDI note number for Middle C (C4)
	velocity    = 100 // Note velocity (0-127)
	velocityOff = 0   // Velocity 0 acts as Note Off in many contexts

	ccModWheel    = 1   // Control Change number: Modulation Wheel
	ccVolume      = 7   // Control Change number: Channel Volume
	ccSustain     = 64  // Control Change number: Sustain Pedal
	ccAllNotesOff = 123 // Control Change number: All Notes Off
)

// midiMsg is a convenience alias for building short MIDI messages.
type midiMsg = [3]byte

// noteOnMsg constructs a 3-byte Note On message.
// Channel is 0-15, note is 0-127, vel is 0-127.
func noteOnMsg(channel, note, vel byte) midiMsg {
	return midiMsg{noteOn | (channel & 0x0F), note & 0x7F, vel & 0x7F}
}

// noteOffMsg constructs a 3-byte Note Off message.
func noteOffMsg(channel, note, vel byte) midiMsg {
	return midiMsg{noteOff | (channel & 0x0F), note & 0x7F, vel & 0x7F}
}

// ccMsg constructs a 3-byte Control Change message.
func ccMsg(channel, controller, value byte) midiMsg {
	return midiMsg{controlChange | (channel & 0x0F), controller & 0x7F, value & 0x7F}
}

// printMsgTable prints a reference table of MIDI message formats.
func printMsgTable() {
	fmt.Println("MIDI channel voice message format (channel 0):")
	fmt.Println()
	fmt.Printf("  %-18s  Status  Data1       Data2\n", "Message")
	fmt.Printf("  %-18s  ------  ----------  ----------\n", repeatByte(18, '-'))
	fmt.Printf("  %-18s  0x%02X    note 0-127  vel  0-127\n", "Note Off", noteOff)
	fmt.Printf("  %-18s  0x%02X    note 0-127  vel  0-127\n", "Note On", noteOn)
	fmt.Printf("  %-18s  0x%02X    cc#  0-127  val  0-127\n", "Control Change", controlChange)
	fmt.Printf("  %-18s  0x%02X    pgm  0-127  (none)\n", "Program Change", programChange)
	fmt.Println()
	fmt.Println("  Status byte high nibble selects the message type.")
	fmt.Println("  Status byte low nibble selects the MIDI channel (0-15).")
	fmt.Println()
}

// repeatByte returns a string of n copies of byte c. Used for table formatting.
func repeatByte(n int, c byte) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = c
	}
	return string(b)
}

// sendMsg sends a MIDI message through an InputPort. InputPort.Send takes a
// pointer to the first byte and the message length.
func sendMsg(port *midi.InputPort, msg []byte) error {
	return port.Send(&msg[0], uint64(len(msg)))
}

// showDeviceInfo prints port counts for a MIDI device.
func showDeviceInfo(dev *midi.Device) {
	fmt.Printf("  Input ports:  %d  (send MIDI data TO the device)\n", dev.NumInputPorts())
	fmt.Printf("  Output ports: %d  (receive MIDI data FROM the device)\n", dev.NumOutputPorts())
	fmt.Println()
}

// demonstrateLifecycle shows the full open-send-close sequence on a device.
// In production code the Device would be obtained via JNI; here we accept
// it as a parameter to document the API calls.
func demonstrateLifecycle(dev *midi.Device) error {
	showDeviceInfo(dev)

	if dev.NumInputPorts() == 0 {
		fmt.Println("  Device has no input ports; skipping send demo.")
		return nil
	}

	// Open the first input port (port 0). An "input port" from the NDK's
	// perspective is a port that accepts data sent FROM the app TO the
	// device (e.g. a synthesizer receiving note events).
	port, err := dev.OpenInputPort(0)
	if err != nil {
		return fmt.Errorf("open input port 0: %w", err)
	}
	defer port.Close()

	fmt.Println("  Opened input port 0.")

	// Build and send example MIDI messages.
	messages := []struct {
		label string
		data  midiMsg
	}{
		{"Note On  (Middle C, vel 100)", noteOnMsg(0, middleC, velocity)},
		{"CC Mod Wheel = 64", ccMsg(0, ccModWheel, 64)},
		{"CC Volume = 100", ccMsg(0, ccVolume, 100)},
		{"Note Off (Middle C)", noteOffMsg(0, middleC, 0)},
		{"CC All Notes Off", ccMsg(0, ccAllNotesOff, 0)},
	}

	for _, m := range messages {
		msg := m.data[:]
		fmt.Printf("  Send %-32s  [% 02X]\n", m.label, msg)
		if err := sendMsg(port, msg); err != nil {
			return fmt.Errorf("send %s: %w", m.label, err)
		}
	}

	// Program Change is a 2-byte message (status + program number).
	pc := []byte{programChange, 0} // select program 0
	fmt.Printf("  Send %-32s  [% 02X]\n", "Program Change 0", pc)
	if err := sendMsg(port, pc); err != nil {
		return fmt.Errorf("send program change: %w", err)
	}

	fmt.Println()
	fmt.Println("  All messages sent. Closing port.")
	return nil
}

func main() {
	fmt.Println("github.com/AndroidGoLab/ndk/midi API Overview")
	fmt.Println("=======================")
	fmt.Println()

	// -- Reference: MIDI message format table --
	printMsgTable()

	// -- Reference: constructing MIDI messages as byte arrays --
	fmt.Println("Example MIDI messages as Go byte arrays:")
	fmt.Println()

	on := noteOnMsg(0, middleC, velocity)
	off := noteOffMsg(0, middleC, 0)
	cc := ccMsg(0, ccSustain, 127)
	pc := [2]byte{programChange, 42}

	fmt.Printf("  Note On  (C4, vel %d):    [% 02X]\n", velocity, on[:])
	fmt.Printf("  Note Off (C4):            [% 02X]\n", off[:])
	fmt.Printf("  CC Sustain Pedal On:      [% 02X]\n", cc[:])
	fmt.Printf("  Program Change 42:        [% 02X]\n", pc[:])
	fmt.Println()

	// -- Reference: the JNI bridge pattern --
	fmt.Println("Obtaining a Device handle (JNI bridge pattern):")
	fmt.Println()
	fmt.Println("  // Java side:")
	fmt.Println("  //   MidiManager mgr = (MidiManager) getSystemService(MIDI_SERVICE);")
	fmt.Println("  //   mgr.openDevice(deviceInfo, device -> {")
	fmt.Println("  //       nativeOnDeviceOpened(device);  // JNI call into C/Go")
	fmt.Println("  //   }, handler);")
	fmt.Println()
	fmt.Println("  // C/Go JNI bridge (cgo or hand-written):")
	fmt.Println("  //   AMidiDevice *nativeDevice;")
	fmt.Println("  //   media_status_t status = AMidiDevice_fromJava(env, javaMidiDevice, &nativeDevice);")
	fmt.Println("  //")
	fmt.Println("  // The resulting AMidiDevice pointer is wrapped by midi.Device.")
	fmt.Println("  // In Go, you would store it via unsafe.Pointer and construct")
	fmt.Println("  // the Device value on the Go side.")
	fmt.Println()

	// -- Reference: device lifecycle --
	fmt.Println("Device lifecycle (requires a live Android MIDI device):")
	fmt.Println()
	fmt.Println("  dev := obtainDeviceViaJNI()  // see pattern above")
	fmt.Println("  defer dev.Close()")
	fmt.Println()
	fmt.Println("  port, err := dev.OpenInputPort(0)")
	fmt.Println("  if err != nil { log.Fatal(err) }")
	fmt.Println("  defer port.Close()")
	fmt.Println()
	fmt.Println("  msg := []byte{0x90, 60, 100}  // Note On, C4, velocity 100")
	fmt.Println("  err = port.Send(&msg[0], uint64(len(msg)))")
	fmt.Println()

	// -- Demonstrate the full API with a nil device guard --
	// In a real app, pass a Device obtained via JNI. Here we show what
	// the code would look like; calling demonstrateLifecycle with an
	// actual device is left to the integrator.
	fmt.Println("InputPort.Send signature:")
	fmt.Printf("  func (h *InputPort) Send(buffer *uint8, numBytes uint64) error\n")
	fmt.Println()
	fmt.Println("  buffer:   pointer to the first byte of the MIDI message")
	fmt.Println("  numBytes: length of the message in bytes")
	fmt.Println()

	// Show how unsafe.Pointer can convert a slice for the Send call.
	fmt.Println("Converting a Go slice for Send():")
	fmt.Println()
	buf := []byte{noteOn, middleC, velocity}
	ptr := &buf[0]
	fmt.Printf("  buf  = %v\n", buf)
	fmt.Printf("  &buf[0] is %T (value %p)\n", ptr, unsafe.Pointer(ptr))
	fmt.Printf("  port.Send(&buf[0], %d)\n", len(buf))
	fmt.Println()

	fmt.Println("Error codes:")
	fmt.Printf("  midi.ErrUnknown    = %d\n", midi.ErrUnknown)
	fmt.Printf("  midi.ErrWouldBlock = %d\n", midi.ErrWouldBlock)
	fmt.Println()
	fmt.Println("Done.")
}
