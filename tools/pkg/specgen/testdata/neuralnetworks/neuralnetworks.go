// Simulates c-for-go output for Android Neural Networks API (NeuralNetworks.h).
// This file is parsed at AST level only; it does not compile.
package neuralnetworks

import "unsafe"

// Opaque handle types.
type ANeuralNetworksModel C.ANeuralNetworksModel
type ANeuralNetworksCompilation C.ANeuralNetworksCompilation
type ANeuralNetworksExecution C.ANeuralNetworksExecution
type ANeuralNetworksMemory C.ANeuralNetworksMemory
type ANeuralNetworksBurst C.ANeuralNetworksBurst
type ANeuralNetworksDevice C.ANeuralNetworksDevice
type ANeuralNetworksEvent C.ANeuralNetworksEvent

// Integer typedefs.
type Nnapi_result_t int32
type Nnapi_operand_type_t int32
type Nnapi_operation_type_t int32
type Nnapi_fuse_code_t int32
type Nnapi_preference_t int32
type Nnapi_device_type_t int32
type Nnapi_priority_t int32

// Result codes.
const (
	ANEURALNETWORKS_NO_ERROR                       Nnapi_result_t = 0
	ANEURALNETWORKS_OUT_OF_MEMORY                  Nnapi_result_t = 1
	ANEURALNETWORKS_INCOMPLETE                     Nnapi_result_t = 2
	ANEURALNETWORKS_UNEXPECTED_NULL                Nnapi_result_t = 3
	ANEURALNETWORKS_BAD_DATA                       Nnapi_result_t = 4
	ANEURALNETWORKS_OP_FAILED                      Nnapi_result_t = 5
	ANEURALNETWORKS_BAD_STATE                      Nnapi_result_t = 6
	ANEURALNETWORKS_UNMAPPABLE                     Nnapi_result_t = 7
	ANEURALNETWORKS_OUTPUT_INSUFFICIENT_SIZE       Nnapi_result_t = 8
	ANEURALNETWORKS_UNAVAILABLE_DEVICE             Nnapi_result_t = 9
	ANEURALNETWORKS_MISSED_DEADLINE_TRANSIENT      Nnapi_result_t = 10
	ANEURALNETWORKS_MISSED_DEADLINE_PERSISTENT     Nnapi_result_t = 11
	ANEURALNETWORKS_RESOURCE_EXHAUSTED_TRANSIENT   Nnapi_result_t = 12
	ANEURALNETWORKS_RESOURCE_EXHAUSTED_PERSISTENT  Nnapi_result_t = 13
	ANEURALNETWORKS_DEAD_OBJECT                    Nnapi_result_t = 14
)

// Operand types.
const (
	ANEURALNETWORKS_FLOAT32                         Nnapi_operand_type_t = 0
	ANEURALNETWORKS_INT32                           Nnapi_operand_type_t = 1
	ANEURALNETWORKS_UINT32                          Nnapi_operand_type_t = 2
	ANEURALNETWORKS_TENSOR_FLOAT32                  Nnapi_operand_type_t = 3
	ANEURALNETWORKS_TENSOR_INT32                    Nnapi_operand_type_t = 4
	ANEURALNETWORKS_TENSOR_QUANT8_ASYMM             Nnapi_operand_type_t = 5
	ANEURALNETWORKS_BOOL                            Nnapi_operand_type_t = 6
	ANEURALNETWORKS_TENSOR_QUANT16_SYMM             Nnapi_operand_type_t = 7
	ANEURALNETWORKS_TENSOR_FLOAT16                  Nnapi_operand_type_t = 8
	ANEURALNETWORKS_TENSOR_BOOL8                    Nnapi_operand_type_t = 9
	ANEURALNETWORKS_FLOAT16                         Nnapi_operand_type_t = 10
	ANEURALNETWORKS_TENSOR_QUANT8_SYMM_PER_CHANNEL  Nnapi_operand_type_t = 11
	ANEURALNETWORKS_TENSOR_QUANT16_ASYMM            Nnapi_operand_type_t = 12
	ANEURALNETWORKS_TENSOR_QUANT8_SYMM              Nnapi_operand_type_t = 13
	ANEURALNETWORKS_TENSOR_QUANT8_ASYMM_SIGNED      Nnapi_operand_type_t = 14
)

// Operation types.
const (
	ANEURALNETWORKS_ADD                Nnapi_operation_type_t = 0
	ANEURALNETWORKS_AVERAGE_POOL_2D    Nnapi_operation_type_t = 1
	ANEURALNETWORKS_CONCATENATION      Nnapi_operation_type_t = 2
	ANEURALNETWORKS_CONV_2D            Nnapi_operation_type_t = 3
	ANEURALNETWORKS_DEPTHWISE_CONV_2D  Nnapi_operation_type_t = 4
	ANEURALNETWORKS_FULLY_CONNECTED    Nnapi_operation_type_t = 9
	ANEURALNETWORKS_L2_NORMALIZATION   Nnapi_operation_type_t = 11
	ANEURALNETWORKS_LOGISTIC           Nnapi_operation_type_t = 14
	ANEURALNETWORKS_LSTM               Nnapi_operation_type_t = 16
	ANEURALNETWORKS_MAX_POOL_2D        Nnapi_operation_type_t = 17
	ANEURALNETWORKS_MUL                Nnapi_operation_type_t = 18
	ANEURALNETWORKS_RELU               Nnapi_operation_type_t = 19
	ANEURALNETWORKS_RELU1              Nnapi_operation_type_t = 20
	ANEURALNETWORKS_RELU6              Nnapi_operation_type_t = 21
	ANEURALNETWORKS_RESHAPE            Nnapi_operation_type_t = 22
	ANEURALNETWORKS_SOFTMAX            Nnapi_operation_type_t = 25
	ANEURALNETWORKS_TANH               Nnapi_operation_type_t = 28
	ANEURALNETWORKS_BATCH_NORMALIZATION Nnapi_operation_type_t = 29
)

// Fuse codes.
const (
	ANEURALNETWORKS_FUSED_NONE  Nnapi_fuse_code_t = 0
	ANEURALNETWORKS_FUSED_RELU  Nnapi_fuse_code_t = 1
	ANEURALNETWORKS_FUSED_RELU1 Nnapi_fuse_code_t = 2
	ANEURALNETWORKS_FUSED_RELU6 Nnapi_fuse_code_t = 3
)

// Preference constants.
const (
	ANEURALNETWORKS_PREFER_LOW_POWER          Nnapi_preference_t = 0
	ANEURALNETWORKS_PREFER_FAST_SINGLE_ANSWER Nnapi_preference_t = 1
	ANEURALNETWORKS_PREFER_SUSTAINED_SPEED    Nnapi_preference_t = 2
)

// Device types.
const (
	ANEURALNETWORKS_DEVICE_UNKNOWN     Nnapi_device_type_t = 0
	ANEURALNETWORKS_DEVICE_OTHER       Nnapi_device_type_t = 1
	ANEURALNETWORKS_DEVICE_CPU         Nnapi_device_type_t = 2
	ANEURALNETWORKS_DEVICE_GPU         Nnapi_device_type_t = 3
	ANEURALNETWORKS_DEVICE_ACCELERATOR Nnapi_device_type_t = 4
)

// Priority constants.
const (
	ANEURALNETWORKS_PRIORITY_LOW     Nnapi_priority_t = 90
	ANEURALNETWORKS_PRIORITY_MEDIUM  Nnapi_priority_t = 100
	ANEURALNETWORKS_PRIORITY_HIGH    Nnapi_priority_t = 110
	ANEURALNETWORKS_PRIORITY_DEFAULT Nnapi_priority_t = 100
)

// Operand type descriptor (opaque for CGo bridge).
type ANeuralNetworksOperandType C.ANeuralNetworksOperandType

// --- Model functions ---
func ANeuralNetworksModel_create(model **ANeuralNetworksModel) Nnapi_result_t { return 0 }
func ANeuralNetworksModel_free(model *ANeuralNetworksModel)                   {}
func ANeuralNetworksModel_finish(model *ANeuralNetworksModel) Nnapi_result_t  { return 0 }
func ANeuralNetworksModel_addOperand(model *ANeuralNetworksModel, operandType *ANeuralNetworksOperandType) Nnapi_result_t {
	return 0
}
func ANeuralNetworksModel_setOperandValue(model *ANeuralNetworksModel, index int32, buffer unsafe.Pointer, length uint32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksModel_addOperation(model *ANeuralNetworksModel, operationType Nnapi_operation_type_t, inputCount uint32, inputs *uint32, outputCount uint32, outputs *uint32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksModel_identifyInputsAndOutputs(model *ANeuralNetworksModel, inputCount uint32, inputs *uint32, outputCount uint32, outputs *uint32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksModel_relaxComputationFloat32toFloat16(model *ANeuralNetworksModel, allow int32) Nnapi_result_t {
	return 0
}

// --- Compilation functions ---
func ANeuralNetworksCompilation_create(model *ANeuralNetworksModel, compilation **ANeuralNetworksCompilation) Nnapi_result_t {
	return 0
}
func ANeuralNetworksCompilation_free(compilation *ANeuralNetworksCompilation) {}
func ANeuralNetworksCompilation_setPreference(compilation *ANeuralNetworksCompilation, preference Nnapi_preference_t) Nnapi_result_t {
	return 0
}
func ANeuralNetworksCompilation_finish(compilation *ANeuralNetworksCompilation) Nnapi_result_t {
	return 0
}
func ANeuralNetworksCompilation_setPriority(compilation *ANeuralNetworksCompilation, priority Nnapi_priority_t) Nnapi_result_t {
	return 0
}
func ANeuralNetworksCompilation_setTimeout(compilation *ANeuralNetworksCompilation, duration uint64) Nnapi_result_t {
	return 0
}
func ANeuralNetworksCompilation_setCaching(compilation *ANeuralNetworksCompilation, cacheDir *byte, token *byte) Nnapi_result_t {
	return 0
}

// --- Execution functions ---
func ANeuralNetworksExecution_create(compilation *ANeuralNetworksCompilation, execution **ANeuralNetworksExecution) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_free(execution *ANeuralNetworksExecution) {}
func ANeuralNetworksExecution_setInput(execution *ANeuralNetworksExecution, index int32, operandType *ANeuralNetworksOperandType, buffer unsafe.Pointer, length uint32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_setOutput(execution *ANeuralNetworksExecution, index int32, operandType *ANeuralNetworksOperandType, buffer unsafe.Pointer, length uint32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_compute(execution *ANeuralNetworksExecution) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_startCompute(execution *ANeuralNetworksExecution, event **ANeuralNetworksEvent) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_setTimeout(execution *ANeuralNetworksExecution, duration uint64) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_setMeasureTiming(execution *ANeuralNetworksExecution, measure int32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_getDuration(execution *ANeuralNetworksExecution, durationCode int32, duration *uint64) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_setInputFromMemory(execution *ANeuralNetworksExecution, index int32, operandType *ANeuralNetworksOperandType, memory *ANeuralNetworksMemory, offset uint32, length uint32) Nnapi_result_t {
	return 0
}
func ANeuralNetworksExecution_setOutputFromMemory(execution *ANeuralNetworksExecution, index int32, operandType *ANeuralNetworksOperandType, memory *ANeuralNetworksMemory, offset uint32, length uint32) Nnapi_result_t {
	return 0
}

// --- Memory functions ---
func ANeuralNetworksMemory_createFromFd(size uint32, prot int32, fd int32, offset uint32, memory **ANeuralNetworksMemory) Nnapi_result_t {
	return 0
}
func ANeuralNetworksMemory_free(memory *ANeuralNetworksMemory) {}

// --- Burst functions ---
func ANeuralNetworksBurst_create(compilation *ANeuralNetworksCompilation, burst **ANeuralNetworksBurst) Nnapi_result_t {
	return 0
}
func ANeuralNetworksBurst_free(burst *ANeuralNetworksBurst) {}

// --- Event functions ---
func ANeuralNetworksEvent_wait(event *ANeuralNetworksEvent) Nnapi_result_t { return 0 }
func ANeuralNetworksEvent_free(event *ANeuralNetworksEvent)                {}

// --- Device enumeration functions ---
func ANeuralNetworks_getDeviceCount(numDevices *uint32) Nnapi_result_t { return 0 }
func ANeuralNetworks_getDevice(devIndex uint32, device **ANeuralNetworksDevice) Nnapi_result_t {
	return 0
}
func ANeuralNetworksDevice_getName(device *ANeuralNetworksDevice, name **byte) Nnapi_result_t {
	return 0
}
func ANeuralNetworksDevice_getType(device *ANeuralNetworksDevice, devType *Nnapi_device_type_t) Nnapi_result_t {
	return 0
}
func ANeuralNetworksDevice_getVersion(device *ANeuralNetworksDevice, version **byte) Nnapi_result_t {
	return 0
}
func ANeuralNetworksDevice_getFeatureLevel(device *ANeuralNetworksDevice, featureLevel *int64) Nnapi_result_t {
	return 0
}

var _ = unsafe.Pointer(nil)
