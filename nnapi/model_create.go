package nnapi

import (
	capi "github.com/AndroidGoLab/ndk/capi/neuralnetworks"
)

// NewModel creates a new empty neural network model.
func NewModel() (*Model, error) {
	var ptr *capi.ANeuralNetworksModel
	if err := result(int32(capi.ANeuralNetworksModel_create(&ptr))); err != nil {
		return nil, err
	}
	return &Model{ptr: ptr}, nil
}

// Close releases the underlying NDK handle.
func (h *Model) Close() error {
	if h.ptr == nil {
		return nil
	}
	capi.ANeuralNetworksModel_free(h.ptr)
	h.ptr = nil
	return nil
}
