package c2ffi

import (
	"testing"

	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertLooper(t *testing.T) {
	input := `[
		{"tag":"struct","ns":0,"name":"ALooper","id":0,"location":"android/looper.h:41:8","bit-size":0,"bit-alignment":0,"fields":[]},
		{"tag":"typedef","ns":0,"name":"ALooper","location":"android/looper.h:55:24","type":{"tag":"struct","ns":0,"name":"ALooper","id":0}},
		{"tag":"function","name":"ALooper_forThread","ns":0,"location":"android/looper.h:61:10","variadic":false,"inline":false,"storage-class":"none","parameters":[],"return-type":{"tag":":pointer","type":{"tag":"ALooper"}}},
		{"tag":"enum","ns":0,"name":"","id":3,"location":"android/looper.h:85:1","fields":[
			{"tag":"field","name":"ALOOPER_POLL_WAKE","value":4294967295},
			{"tag":"field","name":"ALOOPER_POLL_CALLBACK","value":4294967294}
		]},
		{"tag":"typedef","ns":0,"name":"ALooper_callbackFunc","location":"android/looper.h:179:15","type":{"tag":":function-pointer"}},
		{"tag":"function","name":"ALooper_addFd","ns":0,"location":"android/looper.h:275:5","variadic":false,"inline":false,"storage-class":"none","parameters":[
			{"tag":"parameter","name":"looper","type":{"tag":":pointer","type":{"tag":"ALooper"}}},
			{"tag":"parameter","name":"fd","type":{"tag":":int","bit-size":32,"bit-alignment":32}},
			{"tag":"parameter","name":"callback","type":{"tag":"ALooper_callbackFunc"}},
			{"tag":"parameter","name":"data","type":{"tag":":pointer","type":{"tag":":void"}}}
		],"return-type":{"tag":":int","bit-size":32,"bit-alignment":32}}
	]`

	opts := ConvertOptions{
		Module:        "looper",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/looper",
		TargetHeaders: []string{"android/looper.h"},
		Rules: []Rule{
			{Action: "accept", From: "^ALooper"},
			{Action: "accept", From: "^ALOOPER_"},
		},
	}

	spec, err := Convert([]byte(input), opts)
	require.NoError(t, err)

	// Type: ALooper → opaque_ptr.
	assert.Contains(t, spec.Types, "ALooper")
	assert.Equal(t, "opaque_ptr", spec.Types["ALooper"].Kind)

	// Functions.
	assert.Contains(t, spec.Functions, "ALooper_forThread")
	assert.Equal(t, "*ALooper", spec.Functions["ALooper_forThread"].Returns)

	assert.Contains(t, spec.Functions, "ALooper_addFd")
	fd := spec.Functions["ALooper_addFd"]
	assert.Equal(t, "int32", fd.Returns)
	assert.Len(t, fd.Params, 4)
	assert.Equal(t, "*ALooper", fd.Params[0].Type)
	assert.Equal(t, "int32", fd.Params[1].Type)
	assert.Equal(t, "ALooper_callbackFunc", fd.Params[2].Type)
	assert.Equal(t, "unsafe.Pointer", fd.Params[3].Type)

	// Enums: negative values via unsigned-to-signed conversion.
	assert.Contains(t, spec.Enums, "ALOOPER_POLL")
	pollVals := spec.Enums["ALOOPER_POLL"]
	assert.Len(t, pollVals, 2)
	assert.Equal(t, int64(-1), pollVals[0].Value)
	assert.Equal(t, int64(-2), pollVals[1].Value)

	// Callback (empty params — no header dirs for supplement).
	assert.Contains(t, spec.Callbacks, "ALooper_callbackFunc")
}

func TestConvertStructInlineUnionFields(t *testing.T) {
	input := `[
		{"tag":"struct","ns":0,"name":"ACameraMetadata_const_entry","id":1,"location":"camera/NdkCameraMetadata.h:143:16","bit-size":192,"bit-alignment":64,"fields":[
			{"tag":"field","name":"tag","type":{"tag":"uint32_t"}},
			{"tag":"field","name":"count","type":{"tag":"uint32_t"}},
			{"tag":"field","name":"data","type":{"tag":"union","fields":[
				{"tag":"field","name":"u8","type":{"tag":":pointer","type":{"tag":"uint8_t"}}},
				{"tag":"field","name":"i32","type":{"tag":":pointer","type":{"tag":"int32_t"}}},
				{"tag":"field","name":"f","type":{"tag":":pointer","type":{"tag":":float","bit-size":32,"bit-alignment":32}}}
			]}}
		]}
	]`

	opts := ConvertOptions{
		Module:        "camera",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/camera",
		TargetHeaders: []string{"camera/NdkCameraMetadata.h"},
	}

	spec, err := Convert([]byte(input), opts)
	require.NoError(t, err)

	entry, ok := spec.Structs["ACameraMetadata_const_entry"]
	require.True(t, ok)

	require.Len(t, entry.Fields, 3)
	data := entry.Fields[2]
	assert.Equal(t, "data", data.Name)
	assert.Equal(t, "union", data.Type)

	require.Len(t, data.Fields, 3)
	assert.Equal(t, specmodel.StructField{Name: "u8", Type: "*uint8"}, data.Fields[0])
	assert.Equal(t, specmodel.StructField{Name: "i32", Type: "*int32"}, data.Fields[1])
	assert.Equal(t, specmodel.StructField{Name: "f", Type: "*float32"}, data.Fields[2])
}

func TestToSignedInt64(t *testing.T) {
	assert.Equal(t, int64(-1), toSignedInt64(4294967295))
	assert.Equal(t, int64(-4), toSignedInt64(4294967292))
	assert.Equal(t, int64(1), toSignedInt64(1))
	assert.Equal(t, int64(0), toSignedInt64(0))
	assert.Equal(t, int64(2147483647), toSignedInt64(2147483647))
}

func TestTypeRefToGoType(t *testing.T) {
	tests := []struct {
		name string
		ref  TypeRef
		want string
	}{
		{"void", TypeRef{Tag: ":void"}, ""},
		{"int32", TypeRef{Tag: ":int", BitSize: 32}, "int32"},
		{"uint32", TypeRef{Tag: ":unsigned-int", BitSize: 32}, "uint32"},
		{"int64", TypeRef{Tag: ":long", BitSize: 64}, "int64"},
		{"float32", TypeRef{Tag: ":float", BitSize: 32}, "float32"},
		{"float64", TypeRef{Tag: ":double", BitSize: 64}, "float64"},
		{"bool", TypeRef{Tag: ":_Bool"}, "bool"},
		{"void*", TypeRef{Tag: ":pointer", Type: &TypeRef{Tag: ":void"}}, "unsafe.Pointer"},
		{"char*", TypeRef{Tag: ":pointer", Type: &TypeRef{Tag: ":char"}}, "string"},
		{"int*", TypeRef{Tag: ":pointer", Type: &TypeRef{Tag: ":int", BitSize: 32}}, "*int32"},
		{"typedef ref", TypeRef{Tag: "ALooper"}, "ALooper"},
		{"int32_t", TypeRef{Tag: "int32_t"}, "int32"},
		{"size_t", TypeRef{Tag: "size_t"}, "uint64"},
		{"func ptr", TypeRef{Tag: ":function-pointer"}, "unsafe.Pointer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typeRefToGoType(&tt.ref)
			assert.Equal(t, tt.want, got)
		})
	}
}
