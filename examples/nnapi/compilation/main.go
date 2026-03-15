// NNAPI model compilation example.
//
// Builds a simple ADD model (same as the simple-model example), then
// creates a compilation from it with the FastSingleAnswer preference,
// and finalizes the compilation. This demonstrates the full path from
// model construction to a ready-to-execute compiled graph.
//
// The compilation step translates the abstract model into device-specific
// instructions. The preference hint tells the runtime whether to optimize
// for low power, fast single inference, or sustained throughput.
//
// Cleanup order matters: compilation must be closed before the model,
// because the compilation references the model internally.
//
// This program must run on an Android device with API level 27+.
package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/xaionaro-go/ndk/nnapi"
)

// buildAddModel constructs a simple ADD model: output = input0 + input1.
// Returns the finished model ready for compilation.
func buildAddModel() (*nnapi.Model, error) {
	model, err := nnapi.NewModel()
	if err != nil {
		return nil, fmt.Errorf("NewModel: %w", err)
	}

	// Operand layout:
	//   0: input0     TENSOR_FLOAT32 [1, 4]
	//   1: input1     TENSOR_FLOAT32 [1, 4]
	//   2: activation INT32 scalar (fuse code)
	//   3: output     TENSOR_FLOAT32 [1, 4]
	tensorDims := []uint32{1, 4}

	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.TensorFloat32, Dims: tensorDims}); err != nil {
		model.Close()
		return nil, fmt.Errorf("operand 0: %w", err)
	}
	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.TensorFloat32, Dims: tensorDims}); err != nil {
		model.Close()
		return nil, fmt.Errorf("operand 1: %w", err)
	}
	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.Int32}); err != nil {
		model.Close()
		return nil, fmt.Errorf("operand 2: %w", err)
	}
	if err := model.AddOperand(nnapi.OperandTypeDesc{Type: nnapi.TensorFloat32, Dims: tensorDims}); err != nil {
		model.Close()
		return nil, fmt.Errorf("operand 3: %w", err)
	}

	// Set the activation fuse code to FUSED_NONE (no activation).
	fuseCode := int32(nnapi.FusedNone)
	if err := model.SetOperandValue(2, unsafe.Pointer(&fuseCode), uint64(unsafe.Sizeof(fuseCode))); err != nil {
		model.Close()
		return nil, fmt.Errorf("SetOperandValue: %w", err)
	}

	// Wire the ADD operation.
	inputs := []uint32{0, 1, 2}
	outputs := []uint32{3}
	if err := model.AddOperation(nnapi.ANeuralNetworksOperationType(nnapi.Add), uint32(len(inputs)), &inputs[0], uint32(len(outputs)), &outputs[0]); err != nil {
		model.Close()
		return nil, fmt.Errorf("AddOperation: %w", err)
	}

	// Declare model-level inputs (operands 0, 1) and outputs (operand 3).
	modelInputs := []uint32{0, 1}
	modelOutputs := []uint32{3}
	if err := model.IdentifyInputsAndOutputs(uint32(len(modelInputs)), &modelInputs[0], uint32(len(modelOutputs)), &modelOutputs[0]); err != nil {
		model.Close()
		return nil, fmt.Errorf("IdentifyInputsAndOutputs: %w", err)
	}

	if err := model.Finish(); err != nil {
		model.Close()
		return nil, fmt.Errorf("Finish: %w", err)
	}

	return model, nil
}

func main() {
	// --- Step 1: Build the model ---
	model, err := buildAddModel()
	if err != nil {
		log.Fatalf("build model: %v", err)
	}
	defer func() {
		model.Close()
		fmt.Println("model closed")
	}()
	fmt.Println("model built and finished")

	// --- Step 2: Create a compilation from the model ---
	compilation, err := model.NewCompilation()
	if err != nil {
		log.Fatalf("NewCompilation: %v", err)
	}
	// Compilation must be closed before the model.
	defer func() {
		compilation.Close()
		fmt.Println("compilation closed")
	}()
	fmt.Println("compilation created")

	// --- Step 3: Set compilation preferences ---
	//
	// FastSingleAnswer tells the runtime to optimize for the lowest
	// latency of a single inference. Other options:
	//   nnapi.LowPower       - minimize battery usage
	//   nnapi.SustainedSpeed - optimize for repeated inferences
	if err := compilation.SetPreference(nnapi.FastSingleAnswer); err != nil {
		log.Fatalf("SetPreference: %v", err)
	}
	fmt.Println("preference set to FastSingleAnswer")

	// --- Step 4: Finalize the compilation ---
	//
	// After Finish(), the compilation is immutable and ready to create
	// Execution instances for running inference.
	if err := compilation.Finish(); err != nil {
		log.Fatalf("compilation Finish: %v", err)
	}
	fmt.Println("compilation finished")

	// At this point, an Execution could be created from the compilation
	// to run actual inference:
	//
	//   exec, err := compilation.NewExecution()
	//   // exec.SetInput(...)  -- not yet in high-level API
	//   // exec.SetOutput(...) -- not yet in high-level API
	//   // exec.Compute()
	//   // exec.Close()

	// Deferred closes run in LIFO order: compilation first, then model.
	fmt.Println("compilation pipeline complete")
}
