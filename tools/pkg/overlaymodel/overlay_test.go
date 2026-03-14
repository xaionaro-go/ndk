package overlaymodel_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/xaionaro-go/ndk/tools/pkg/overlaymodel"
)

func TestOverlayRoundTrip(t *testing.T) {
	ov := overlaymodel.Overlay{
		Module: "aaudio",
		Package: overlaymodel.PackageOverlay{
			GoName:   "audio",
			GoImport: "github.com/xaionaro-go/ndk/audio",
			Doc:      "Package audio provides Go bindings for Android AAudio.",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"Aaudio_result_t": {
				GoError:      true,
				SuccessValue: "AAUDIO_OK",
				ErrorPrefix:  "audio",
			},
			"AAudioStreamBuilder": {
				GoName:      "StreamBuilder",
				Constructor: "AAudio_createStreamBuilder",
				Destructor:  "AAudioStreamBuilder_delete",
				Pattern:     "builder",
			},
			"AAudioStream": {
				GoName:     "Stream",
				Destructor: "AAudioStream_close",
				Interfaces: []string{"io.Closer"},
			},
			"Aaudio_direction_t": {
				GoName:       "Direction",
				StripPrefix:  "AAUDIO_DIRECTION_",
				StringMethod: true,
			},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"AAudioStreamBuilder_setDeviceId": {
				Receiver: "StreamBuilder",
				GoName:   "SetDeviceID",
				Chain:    true,
			},
			"AAudioStreamBuilder_openStream": {
				Receiver:   "StreamBuilder",
				GoName:     "Open",
				ReturnsNew: "Stream",
			},
			"AAudioStreamBuilder_setDataCallback": {
				Receiver:       "StreamBuilder",
				Chain:          true,
				CallbackParam:  "callback",
				UserdataParam:  "userData",
				GoCallbackType: "DataCallback",
				GoCallbackSig:  "func(stream *Stream, buf []byte, numFrames int) DataCallbackResult",
			},
			"AAudioStream_getState": {
				Receiver: "Stream",
				GoName:   "State",
				Pure:     true,
			},
		},
		APILevels: map[string]int{
			"AAudioStream_release":        30,
			"AAudioStream_getChannelMask": 32,
		},
	}

	data, err := yaml.Marshal(&ov)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got overlaymodel.Overlay
	if err := yaml.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Module != "aaudio" {
		t.Errorf("module = %q", got.Module)
	}
	if got.Package.GoName != "audio" {
		t.Errorf("go_name = %q", got.Package.GoName)
	}
	td := got.Types["AAudioStreamBuilder"]
	if td.Pattern != "builder" || td.Constructor != "AAudio_createStreamBuilder" {
		t.Errorf("builder type = %+v", td)
	}
	fd := got.Functions["AAudioStreamBuilder_setDeviceId"]
	if !fd.Chain || fd.Receiver != "StreamBuilder" {
		t.Errorf("builder setter = %+v", fd)
	}
	fd2 := got.Functions["AAudioStreamBuilder_setDataCallback"]
	if fd2.GoCallbackType != "DataCallback" {
		t.Errorf("callback func = %+v", fd2)
	}
	if got.APILevels["AAudioStream_release"] != 30 {
		t.Errorf("api_level = %d", got.APILevels["AAudioStream_release"])
	}
}

// repoRoot returns the project root (two levels up from this test file).
func repoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file path")
	}
	// overlay_test.go is in tools/pkg/overlaymodel/, so go up three times.
	return filepath.Join(filepath.Dir(filename), "..", "..", "..")
}

func TestLoadAAudioOverlay(t *testing.T) {
	path := filepath.Join(repoRoot(t), "spec", "overlays", "aaudio.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read overlay: %v", err)
	}

	var ov overlaymodel.Overlay
	if err := yaml.Unmarshal(data, &ov); err != nil {
		t.Fatalf("unmarshal overlay: %v", err)
	}

	// Module
	if ov.Module != "aaudio" {
		t.Errorf("module = %q, want %q", ov.Module, "aaudio")
	}

	// Package
	if ov.Package.GoName != "audio" {
		t.Errorf("package.go_name = %q, want %q", ov.Package.GoName, "audio")
	}
	if ov.Package.GoImport != "github.com/xaionaro-go/ndk/audio" {
		t.Errorf("package.go_import = %q", ov.Package.GoImport)
	}

	// Types: 9 entries
	expectedTypes := []string{
		"aaudio_result_t", "AAudioStreamBuilder", "AAudioStream",
		"aaudio_direction_t", "aaudio_format_t", "aaudio_sharing_mode_t",
		"aaudio_performance_mode_t", "aaudio_stream_state_t",
		"aaudio_data_callback_result_t",
	}
	if len(ov.Types) != len(expectedTypes) {
		t.Errorf("types count = %d, want %d", len(ov.Types), len(expectedTypes))
	}
	for _, name := range expectedTypes {
		if _, ok := ov.Types[name]; !ok {
			t.Errorf("missing type %q", name)
		}
	}

	// Spot-check type annotations
	if !ov.Types["aaudio_result_t"].GoError {
		t.Error("aaudio_result_t.go_error should be true")
	}
	if ov.Types["AAudioStreamBuilder"].Pattern != "builder" {
		t.Error("AAudioStreamBuilder.pattern should be 'builder'")
	}
	if ov.Types["AAudioStream"].GoName != "Stream" {
		t.Error("AAudioStream.go_name should be 'Stream'")
	}
	if len(ov.Types["AAudioStream"].Interfaces) != 1 || ov.Types["AAudioStream"].Interfaces[0] != "io.Closer" {
		t.Errorf("AAudioStream.interfaces = %v", ov.Types["AAudioStream"].Interfaces)
	}
	if !ov.Types["aaudio_stream_state_t"].StringMethod {
		t.Error("aaudio_stream_state_t.string_method should be true")
	}

	// Functions: 21 entries (8 builder setters + openStream + setDataCallback +
	// 5 stream getters + 4 stream control + read + write)
	if len(ov.Functions) != 21 {
		t.Errorf("functions count = %d, want 21", len(ov.Functions))
	}

	// Spot-check function annotations
	setDir := ov.Functions["AAudioStreamBuilder_setDirection"]
	if !setDir.Chain || setDir.Receiver != "StreamBuilder" {
		t.Errorf("setDirection = %+v", setDir)
	}
	open := ov.Functions["AAudioStreamBuilder_openStream"]
	if open.ReturnsNew != "Stream" {
		t.Errorf("openStream.returns_new = %q", open.ReturnsNew)
	}
	cb := ov.Functions["AAudioStreamBuilder_setDataCallback"]
	if !cb.Skip {
		t.Errorf("setDataCallback should be skipped, got %+v", cb)
	}
	state := ov.Functions["AAudioStream_getState"]
	if !state.Pure || state.GoName != "State" {
		t.Errorf("getState = %+v", state)
	}
	write := ov.Functions["AAudioStream_write"]
	if write.BufParam != "buffer" || !write.ReturnsFrames || write.BufGoType != "[]byte" {
		t.Errorf("write = %+v", write)
	}

	// API levels: 6 entries
	if len(ov.APILevels) != 6 {
		t.Errorf("api_levels count = %d, want 6", len(ov.APILevels))
	}
	if ov.APILevels["AAudioStream_release"] != 30 {
		t.Errorf("api_level[release] = %d, want 30", ov.APILevels["AAudioStream_release"])
	}
	if ov.APILevels["AAudioStream_getHardwareFormat"] != 34 {
		t.Errorf("api_level[getHardwareFormat] = %d, want 34", ov.APILevels["AAudioStream_getHardwareFormat"])
	}
}

func TestOverlay_NewFields(t *testing.T) {
	yamlData := `
module: camera
package:
  go_name: camera
  go_import: github.com/xaionaro-go/ndk/camera
callback_structs:
  ACameraDevice_StateCallbacks:
    go_name: DeviceStateCallbacks
    context_field: context
    fields:
      onDisconnected:
        go_name: OnDisconnected
        go_signature: "func()"
      onError:
        go_name: OnError
        go_signature: "func(int)"
struct_accessors:
  ACameraIdList:
    count_field: numCameras
    items_field: cameraIds
    item_type: string
    delete_func: ACameraManager_deleteCameraIdList
functions:
  someFunc:
    receiver: Manager
    go_name: SomeFunc
    fixed_params:
      callbacks: "nil"
      numRequests: "1"
lifecycle:
  entry_point: ANativeActivity_onCreate
  activity_type: ANativeActivity
  callbacks_accessor: "activity->callbacks"
  callback_struct: ANativeActivityCallbacks
`
	var ov overlaymodel.Overlay
	err := yaml.Unmarshal([]byte(yamlData), &ov)
	require.NoError(t, err)

	// CallbackStructs
	require.Len(t, ov.CallbackStructs, 1)
	cs := ov.CallbackStructs["ACameraDevice_StateCallbacks"]
	assert.Equal(t, "DeviceStateCallbacks", cs.GoName)
	assert.Equal(t, "context", cs.ContextField)
	require.Len(t, cs.Fields, 2)
	assert.Equal(t, "OnDisconnected", cs.Fields["onDisconnected"].GoName)
	assert.Equal(t, "func()", cs.Fields["onDisconnected"].GoSignature)

	// StructAccessors
	require.Len(t, ov.StructAccessors, 1)
	sa := ov.StructAccessors["ACameraIdList"]
	assert.Equal(t, "numCameras", sa.CountField)
	assert.Equal(t, "cameraIds", sa.ItemsField)
	assert.Equal(t, "string", sa.ItemType)
	assert.Equal(t, "ACameraManager_deleteCameraIdList", sa.DeleteFunc)

	// FixedParams
	require.Len(t, ov.Functions["someFunc"].FixedParams, 2)
	assert.Equal(t, "nil", ov.Functions["someFunc"].FixedParams["callbacks"])
	assert.Equal(t, "1", ov.Functions["someFunc"].FixedParams["numRequests"])

	// Lifecycle
	require.NotNil(t, ov.Lifecycle)
	assert.Equal(t, "ANativeActivity_onCreate", ov.Lifecycle.EntryPoint)
	assert.Equal(t, "ANativeActivity", ov.Lifecycle.ActivityType)
	assert.Equal(t, "activity->callbacks", ov.Lifecycle.CallbacksAccessor)
	assert.Equal(t, "ANativeActivityCallbacks", ov.Lifecycle.CallbackStruct)
}
