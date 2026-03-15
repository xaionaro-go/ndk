package nnapi

import (
	"runtime"
	"unsafe"

	capi "github.com/AndroidGoLab/ndk/capi/neuralnetworks"
)

// OperandTypeDesc describes the type of an operand for AddOperand.
// It mirrors the C struct ANeuralNetworksOperandType.
type OperandTypeDesc struct {
	Type      OperandType
	Dims      []uint32
	Scale     float32
	ZeroPoint int32
}

// operandTypeC mirrors the C ANeuralNetworksOperandType struct layout.
type operandTypeC struct {
	Type           int32
	DimensionCount uint32
	Dimensions     *uint32
	Scale          float32
	ZeroPoint      int32
}

// AddOperand adds an operand to the model with the given type descriptor.
func (h *Model) AddOperand(desc OperandTypeDesc) error {
	var ct operandTypeC
	ct.Type = int32(desc.Type)
	ct.DimensionCount = uint32(len(desc.Dims))
	if len(desc.Dims) > 0 {
		ct.Dimensions = &desc.Dims[0]
	}
	ct.Scale = desc.Scale
	ct.ZeroPoint = desc.ZeroPoint

	var pinner runtime.Pinner
	pinner.Pin(&ct)
	if ct.Dimensions != nil {
		pinner.Pin(ct.Dimensions)
	}
	defer pinner.Unpin()

	goType := (*capi.ANeuralNetworksOperandType)(unsafe.Pointer(&ct))
	return result(int32(capi.ANeuralNetworksModel_addOperand(h.ptr, goType)))
}
