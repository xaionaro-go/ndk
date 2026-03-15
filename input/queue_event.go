package input

import (
	capi "github.com/AndroidGoLab/ndk/capi/input"
)

// GetEvent returns the next available input event from the queue.
// Returns nil if no events are available or an error occurred.
func (h *Queue) GetEvent() *Event {
	var ev *capi.AInputEvent
	if capi.AInputQueue_getEvent(h.ptr, &ev) < 0 {
		return nil
	}
	return &Event{ptr: ev}
}

// PreDispatchEvent performs pre-dispatching of the given event.
// Returns true if the event was consumed by pre-dispatching (e.g. IME)
// and should NOT be passed to FinishEvent.
func (h *Queue) PreDispatchEvent(event *Event) bool {
	return capi.AInputQueue_preDispatchEvent(h.ptr, event.ptr) != 0
}
