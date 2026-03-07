// Simple NNAPI model building example.
//
// Demonstrates the basic pattern for constructing a neural network model
// using the Android Neural Networks API: create a model, add operands
// (two float32 input tensors, one activation scalar, and one output tensor),
// wire them together with an ADD operation, identify model inputs/outputs,
// and finalize the model.
//
// The ADD operation computes: output = input0 + input1, with a fuse code
// that controls optional activation (FUSED_NONE here means no activation).
//
// This program must run on an Android device with API level 27+.
package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/xaionaro-go/ndk/nnapi"
)

func main() {
	// --- Step 1: Create a new empty model ---
	model, err := nnapi.NewModel()
	if err != nil {
		log.Fatalf("create model: %v", err)
	}
	defer func() {
		model.Close()
		fmt.Println("model closed")
	}()
	fmt.Println("model created")

	// --- Step 2: Add operands ---
	//
	// The ADD operation takes three inputs and produces one output:
	//   input0:     TENSOR_FLOAT32 [1, 3]
	//   input1:     TENSOR_FLOAT32 [1, 3]
	//   activation: INT32 scalar (fuse code: 0 = FUSED_NONE)
	//   output:     TENSOR_FLOAT32 [1, 3]
	//
	// Operands are identified by index in the order they are added.
	//   operand 0 -> input0
	//   operand 1 -> input1
	//   operand 2 -> activation (fuse code)
	//   operand 3 -> output

	tensorDims := []uint32{1, 3}

	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.TensorFloat32, Dims: tensorDims}); err != nil {
		log.Fatalf("add operand 0 (input0): %v", err)
	}
	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.TensorFloat32, Dims: tensorDims}); err != nil {
		log.Fatalf("add operand 1 (input1): %v", err)
	}
	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.Int32}); err != nil {
		log.Fatalf("add operand 2 (activation): %v", err)
	}
	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.TensorFloat32, Dims: tensorDims}); err != nil {
		log.Fatalf("add operand 3 (output): %v", err)
	}
	fmt.Println("operands added: 2 input tensors [1,3], 1 activation scalar, 1 output tensor [1,3]")

	// --- Step 3: Set the activation operand value ---
	//
	// The fuse code is a constant known at model-build time, so it is
	// set with SetOperandValue rather than provided at execution time.
	fuseCode := int32(nnapi.FusedNone)
	if err := model.SetOperandValue(2, unsafe.Pointer(&fuseCode), uint64(unsafe.Sizeof(fuseCode))); err != nil {
		log.Fatalf("SetOperandValue (activation): %v", err)
	}
	fmt.Println("activation fuse code set to FUSED_NONE")

	// --- Step 4: Add the ADD operation ---
	//
	// Inputs:  [operand 0, operand 1, operand 2]
	// Outputs: [operand 3]
	inputs := []uint32{0, 1, 2}
	outputs := []uint32{3}
	if err := model.AddOperation(nnapi.ANeuralNetworksOperationType(nnapi.Add), uint32(len(inputs)), &inputs[0], uint32(len(outputs)), &outputs[0]); err != nil {
		log.Fatalf("AddOperation: %v", err)
	}
	fmt.Println("ADD operation added: output = input0 + input1")

	// --- Step 5: Identify model inputs and outputs ---
	//
	// Model inputs are the operands that will be provided at execution
	// time. Model outputs are the operands that will be read back.
	// The activation operand (index 2) is a constant, not a model input.
	modelInputs := []uint32{0, 1}
	modelOutputs := []uint32{3}
	if err := model.IdentifyInputsAndOutputs(uint32(len(modelInputs)), &modelInputs[0], uint32(len(modelOutputs)), &modelOutputs[0]); err != nil {
		log.Fatalf("IdentifyInputsAndOutputs: %v", err)
	}
	fmt.Println("model inputs/outputs identified: inputs=[0,1] outputs=[3]")

	// --- Step 6: Finalize the model ---
	//
	// After Finish(), the model is immutable. No more operands or
	// operations can be added. The model can now be compiled.
	if err := model.Finish(); err != nil {
		log.Fatalf("Finish: %v", err)
	}
	fmt.Println("model finished (immutable)")

	// The deferred Close will release the model.
	fmt.Println("simple model built successfully")
}
