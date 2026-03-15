package nnapi

import (
	capi "github.com/AndroidGoLab/ndk/capi/neuralnetworks"
)

// Close releases the underlying NDK handle.
func (h *Compilation) Close() error {
	if h.ptr == nil {
		return nil
	}
	capi.ANeuralNetworksCompilation_free(h.ptr)
	h.ptr = nil
	return nil
}
