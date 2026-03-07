package idiomgen_test

import (
	"testing"

	"github.com/xaionaro-go/ndk/tools/pkg/idiomgen"
	"github.com/xaionaro-go/ndk/tools/pkg/overlaymodel"
	"github.com/xaionaro-go/ndk/tools/pkg/specmodel"
)

// buildFixture creates AAudio-like spec and overlay data for testing.
func buildFixture() (specmodel.Spec, overlaymodel.Overlay) {
	spec := specmodel.Spec{
		Module:        "aaudio",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/aaudio",
		Types: map[string]specmodel.TypeDef{
			"AAudioStreamBuilder": {
				Kind:   "opaque_ptr",
				CType:  "AAudioStreamBuilder",
				GoType: "*C.AAudioStreamBuilder",
			},
			"Aaudio_result_t": {
				Kind:   "typedef_int32",
				CType:  "aaudio_result_t",
				GoType: "int32",
			},
			"Aaudio_direction_t": {
				Kind:   "typedef_int32",
				CType:  "aaudio_direction_t",
				GoType: "int32",
			},
		},
		Enums: map[string][]specmodel.EnumValue{
			"Aaudio_result_t": {
				{Name: "AAUDIO_OK", Value: 0},
				{Name: "AAUDIO_ERROR_BASE", Value: -900},
				{Name: "AAUDIO_ERROR_DISCONNECTED", Value: -899},
			},
			"Aaudio_direction_t": {
				{Name: "AAUDIO_DIRECTION_OUTPUT", Value: 0},
				{Name: "AAUDIO_DIRECTION_INPUT", Value: 1},
			},
		},
		Functions: map[string]specmodel.FuncDef{
			"AAudioStreamBuilder_setDeviceId": {
				CName: "AAudioStreamBuilder_setDeviceId",
				Params: []specmodel.Param{
					{Name: "builder", Type: "*AAudioStreamBuilder"},
					{Name: "deviceId", Type: "int32"},
				},
				Returns: "void",
			},
		},
		Callbacks: map[string]specmodel.CallbackDef{
			"AAudioStream_dataCallback": {
				Params: []specmodel.Param{
					{Name: "stream", Type: "*AAudioStream"},
					{Name: "userData", Type: "unsafe.Pointer"},
					{Name: "audioData", Type: "*void"},
					{Name: "numFrames", Type: "int32"},
				},
				Returns: "Aaudio_data_callback_result_t",
			},
		},
	}

	overlay := overlaymodel.Overlay{
		Module: "aaudio",
		Package: overlaymodel.PackageOverlay{
			GoName:   "audio",
			GoImport: "github.com/xaionaro-go/ndk/audio",
			Doc:      "Package audio provides Go bindings for Android AAudio.",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"AAudioStreamBuilder": {
				GoName:      "StreamBuilder",
				Constructor: "AAudio_createStreamBuilder",
				Destructor:  "AAudioStreamBuilder_delete",
				Pattern:     "builder",
			},
			"Aaudio_result_t": {
				GoError:      true,
				SuccessValue: "AAUDIO_OK",
				ErrorPrefix:  "audio",
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
		},
		APILevels: map[string]int{
			"AAudioStream_release": 30,
		},
	}

	return spec, overlay
}

func TestMerge_PackageName(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if merged.PackageName != "audio" {
		t.Errorf("PackageName = %q, want %q", merged.PackageName, "audio")
	}
	if merged.PackageImport != "github.com/xaionaro-go/ndk/audio" {
		t.Errorf("PackageImport = %q, want %q", merged.PackageImport, "github.com/xaionaro-go/ndk/audio")
	}
	if merged.PackageDoc != "Package audio provides Go bindings for Android AAudio." {
		t.Errorf("PackageDoc = %q", merged.PackageDoc)
	}
	if merged.SourcePackage != "github.com/xaionaro-go/ndk/capi/aaudio" {
		t.Errorf("SourcePackage = %q, want %q", merged.SourcePackage, "github.com/xaionaro-go/ndk/capi/aaudio")
	}
}

func TestMerge_OpaqueTypes(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.OpaqueTypes) != 1 {
		t.Fatalf("OpaqueTypes count = %d, want 1", len(merged.OpaqueTypes))
	}
	sb, ok := merged.OpaqueTypes["StreamBuilder"]
	if !ok {
		t.Fatal("OpaqueTypes missing StreamBuilder")
	}
	if sb.GoName != "StreamBuilder" {
		t.Errorf("GoName = %q, want %q", sb.GoName, "StreamBuilder")
	}
	if sb.CapiType != "AAudioStreamBuilder" {
		t.Errorf("CapiType = %q, want %q", sb.CapiType, "AAudioStreamBuilder")
	}
	if sb.Constructor != "AAudio_createStreamBuilder" {
		t.Errorf("Constructor = %q", sb.Constructor)
	}
	if sb.Destructor != "AAudioStreamBuilder_delete" {
		t.Errorf("Destructor = %q", sb.Destructor)
	}
	if sb.Pattern != "builder" {
		t.Errorf("Pattern = %q, want %q", sb.Pattern, "builder")
	}
}

func TestMerge_ErrorEnums(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.ErrorEnums) != 1 {
		t.Fatalf("ErrorEnums count = %d, want 1", len(merged.ErrorEnums))
	}
	ee := merged.ErrorEnums[0]
	if ee.GoName != "Aaudio_result_t" {
		t.Errorf("GoName = %q, want %q", ee.GoName, "Aaudio_result_t")
	}
	if ee.Prefix != "audio" {
		t.Errorf("Prefix = %q, want %q", ee.Prefix, "audio")
	}
	if ee.SuccessValue != "AAUDIO_OK" {
		t.Errorf("SuccessValue = %q", ee.SuccessValue)
	}
	if len(ee.Values) != 3 {
		t.Fatalf("error enum values count = %d, want 3", len(ee.Values))
	}
	// Values should preserve the spec names as both GoName and SpecName.
	if ee.Values[0].SpecName != "AAUDIO_OK" {
		t.Errorf("Values[0].SpecName = %q", ee.Values[0].SpecName)
	}
}

func TestMerge_ValueEnums(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.ValueEnums) != 1 {
		t.Fatalf("ValueEnums count = %d, want 1", len(merged.ValueEnums))
	}
	ve := merged.ValueEnums[0]
	if ve.GoName != "Direction" {
		t.Errorf("GoName = %q, want %q", ve.GoName, "Direction")
	}
	if ve.SpecName != "Aaudio_direction_t" {
		t.Errorf("SpecName = %q, want %q", ve.SpecName, "Aaudio_direction_t")
	}
	if !ve.StringMethod {
		t.Error("StringMethod = false, want true")
	}
	if len(ve.Values) != 2 {
		t.Fatalf("value enum values count = %d, want 2", len(ve.Values))
	}
	// Prefix stripping: "AAUDIO_DIRECTION_OUTPUT" → "Output"
	if ve.Values[0].GoName != "Output" {
		t.Errorf("Values[0].GoName = %q, want %q", ve.Values[0].GoName, "Output")
	}
	if ve.Values[0].SpecName != "AAUDIO_DIRECTION_OUTPUT" {
		t.Errorf("Values[0].SpecName = %q", ve.Values[0].SpecName)
	}
	if ve.Values[1].GoName != "Input" {
		t.Errorf("Values[1].GoName = %q, want %q", ve.Values[1].GoName, "Input")
	}
}

func TestMerge_Methods(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	m := merged.Methods[0]
	if m.GoName != "SetDeviceID" {
		t.Errorf("GoName = %q, want %q", m.GoName, "SetDeviceID")
	}
	if m.CName != "AAudioStreamBuilder_setDeviceId" {
		t.Errorf("CName = %q", m.CName)
	}
	if m.ReceiverType != "StreamBuilder" {
		t.Errorf("ReceiverType = %q, want %q", m.ReceiverType, "StreamBuilder")
	}
	if !m.Chain {
		t.Error("Chain = false, want true")
	}
	// The receiver param should be excluded from Params.
	// The function has builder (*AAudioStreamBuilder) and deviceId (int32).
	// builder is the receiver, so only deviceId should remain.
	if len(m.Params) != 1 {
		t.Fatalf("Params count = %d, want 1", len(m.Params))
	}
	if m.Params[0].Name != "deviceId" {
		t.Errorf("Params[0].Name = %q, want %q", m.Params[0].Name, "deviceId")
	}
}

func TestMerge_Callbacks(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.Callbacks) != 1 {
		t.Fatalf("Callbacks count = %d, want 1", len(merged.Callbacks))
	}
	cb := merged.Callbacks[0]
	if cb.SpecName != "AAudioStream_dataCallback" {
		t.Errorf("SpecName = %q", cb.SpecName)
	}
	if len(cb.Params) != 4 {
		t.Fatalf("Params count = %d, want 4", len(cb.Params))
	}
	if cb.Returns != "Aaudio_data_callback_result_t" {
		t.Errorf("Returns = %q", cb.Returns)
	}
}

func TestMerge_APILevels(t *testing.T) {
	spec, overlay := buildFixture()
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.APILevels) != 1 {
		t.Fatalf("APILevels count = %d, want 1", len(merged.APILevels))
	}
	if merged.APILevels["AAudioStream_release"] != 30 {
		t.Errorf("APILevel = %d, want 30", merged.APILevels["AAudioStream_release"])
	}
}

func TestMerge_FunctionWithoutReceiver_Skipped(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Functions: map[string]specmodel.FuncDef{
			"freeFunction": {
				CName:   "freeFunction",
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"freeFunction": {
				GoName: "FreeFunction",
				// No Receiver → should be skipped
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Methods) != 0 {
		t.Errorf("Methods count = %d, want 0 (free functions skipped)", len(merged.Methods))
	}
}

func TestMerge_TypeWithoutOverlayGoName_UsesSpecName(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"SomeOpaqueHandle": {
				Kind:   "opaque_ptr",
				CType:  "SomeOpaqueHandle",
				GoType: "*C.SomeOpaqueHandle",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"SomeOpaqueHandle": {
				// No GoName → should use spec name
				Pattern: "ref_counted",
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.OpaqueTypes) != 1 {
		t.Fatalf("OpaqueTypes count = %d, want 1", len(merged.OpaqueTypes))
	}
	ot, ok := merged.OpaqueTypes["SomeOpaqueHandle"]
	if !ok {
		t.Fatal("OpaqueTypes missing SomeOpaqueHandle")
	}
	if ot.GoName != "SomeOpaqueHandle" {
		t.Errorf("GoName = %q, want %q", ot.GoName, "SomeOpaqueHandle")
	}
	if ot.Pattern != "ref_counted" {
		t.Errorf("Pattern = %q, want %q", ot.Pattern, "ref_counted")
	}
}

func TestMerge_MethodAPILevel(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"Handle": {Kind: "opaque_ptr", CType: "Handle", GoType: "*C.Handle"},
		},
		Functions: map[string]specmodel.FuncDef{
			"Handle_doSomething": {
				CName: "Handle_doSomething",
				Params: []specmodel.Param{
					{Name: "h", Type: "*Handle"},
				},
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"Handle": {GoName: "Handle"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"Handle_doSomething": {
				Receiver: "Handle",
				GoName:   "DoSomething",
			},
		},
		APILevels: map[string]int{
			"Handle_doSomething": 31,
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	if merged.Methods[0].APILevel != 31 {
		t.Errorf("APILevel = %d, want 31", merged.Methods[0].APILevel)
	}
}

func TestMerge_ValueEnum_NoStripPrefix(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"Color_t": {Kind: "typedef_int32", CType: "color_t", GoType: "int32"},
		},
		Enums: map[string][]specmodel.EnumValue{
			"Color_t": {
				{Name: "RED", Value: 0},
				{Name: "GREEN", Value: 1},
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"Color_t": {GoName: "Color"},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.ValueEnums) != 1 {
		t.Fatalf("ValueEnums count = %d, want 1", len(merged.ValueEnums))
	}
	// Without strip_prefix, GoName should be the raw spec name.
	if merged.ValueEnums[0].Values[0].GoName != "RED" {
		t.Errorf("Values[0].GoName = %q, want %q", merged.ValueEnums[0].Values[0].GoName, "RED")
	}
}

func TestMerge_CallbackWithoutOverlayAnnotation(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Callbacks: map[string]specmodel.CallbackDef{
			"myCallback": {
				Params:  []specmodel.Param{{Name: "x", Type: "int32"}},
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Callbacks) != 1 {
		t.Fatalf("Callbacks count = %d, want 1", len(merged.Callbacks))
	}
	cb := merged.Callbacks[0]
	if cb.SpecName != "myCallback" {
		t.Errorf("SpecName = %q", cb.SpecName)
	}
	if cb.GoCallbackType != "" {
		t.Errorf("GoCallbackType = %q, want empty", cb.GoCallbackType)
	}
}

func TestMerge_CallbackAnnotationMatchesByType(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Functions: map[string]specmodel.FuncDef{
			"setDataCallback": {
				CName: "setDataCallback",
				Params: []specmodel.Param{
					{Name: "h", Type: "*Handle"},
					{Name: "cb", Type: "myDataCallback"},
					{Name: "userData", Type: "unsafe.Pointer"},
				},
				Returns: "void",
			},
			"setErrorCallback": {
				CName: "setErrorCallback",
				Params: []specmodel.Param{
					{Name: "h", Type: "*Handle"},
					{Name: "cb", Type: "myErrorCallback"},
					{Name: "userData", Type: "unsafe.Pointer"},
				},
				Returns: "void",
			},
		},
		Callbacks: map[string]specmodel.CallbackDef{
			"myDataCallback": {
				Params:  []specmodel.Param{{Name: "data", Type: "*byte"}},
				Returns: "int32",
			},
			"myErrorCallback": {
				Params:  []specmodel.Param{{Name: "err", Type: "int32"}},
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"Handle": {GoName: "Handle"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"setDataCallback": {
				Receiver:       "Handle",
				GoName:         "SetDataCallback",
				Chain:          true,
				CallbackParam:  "cb",
				GoCallbackType: "DataCallback",
				GoCallbackSig:  "func(data []byte) int",
			},
			"setErrorCallback": {
				Receiver:       "Handle",
				GoName:         "SetErrorCallback",
				Chain:          true,
				CallbackParam:  "cb",
				GoCallbackType: "ErrorCallback",
				GoCallbackSig:  "func(err error)",
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.Callbacks) != 2 {
		t.Fatalf("Callbacks count = %d, want 2", len(merged.Callbacks))
	}

	// Callbacks are sorted by name, so myDataCallback comes first.
	dataCb := merged.Callbacks[0]
	if dataCb.SpecName != "myDataCallback" {
		t.Errorf("Callbacks[0].SpecName = %q, want %q", dataCb.SpecName, "myDataCallback")
	}
	if dataCb.GoCallbackType != "DataCallback" {
		t.Errorf("Callbacks[0].GoCallbackType = %q, want %q", dataCb.GoCallbackType, "DataCallback")
	}
	if dataCb.GoCallbackSig != "func(data []byte) int" {
		t.Errorf("Callbacks[0].GoCallbackSig = %q", dataCb.GoCallbackSig)
	}

	errCb := merged.Callbacks[1]
	if errCb.SpecName != "myErrorCallback" {
		t.Errorf("Callbacks[1].SpecName = %q, want %q", errCb.SpecName, "myErrorCallback")
	}
	if errCb.GoCallbackType != "ErrorCallback" {
		t.Errorf("Callbacks[1].GoCallbackType = %q, want %q", errCb.GoCallbackType, "ErrorCallback")
	}
}

func TestMerge_FunctionNotInOverlay_Skipped(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Functions: map[string]specmodel.FuncDef{
			"unknownFunc": {
				CName:   "unknownFunc",
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "test",
			GoImport: "github.com/xaionaro-go/ndk/test",
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Methods) != 0 {
		t.Errorf("Methods count = %d, want 0", len(merged.Methods))
	}
}

func TestMerge_ErrorEnumWithStripPrefix(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"Result_t": {Kind: "typedef_int32"},
		},
		Enums: map[string][]specmodel.EnumValue{
			"Result_t": {
				{Name: "OK", Value: 0},
				{Name: "ERROR_DISCONNECTED", Value: -1},
				{Name: "ERROR_TIMEOUT", Value: -2},
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "test",
		Package: overlaymodel.PackageOverlay{GoName: "test", GoImport: "github.com/xaionaro-go/ndk/test"},
		Types: map[string]overlaymodel.TypeOverlay{
			"Result_t": {
				GoError:      true,
				SuccessValue: "OK",
				ErrorPrefix:  "test",
				StripPrefix:  "ERROR_",
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.ErrorEnums) != 1 {
		t.Fatalf("ErrorEnums count = %d, want 1", len(merged.ErrorEnums))
	}
	ee := merged.ErrorEnums[0]
	// Success value keeps its raw name.
	if ee.Values[0].GoName != "OK" {
		t.Errorf("Values[0].GoName = %q, want %q", ee.Values[0].GoName, "OK")
	}
	// Error values get Err prefix + stripped + title-cased.
	if ee.Values[1].GoName != "ErrDisconnected" {
		t.Errorf("Values[1].GoName = %q, want %q", ee.Values[1].GoName, "ErrDisconnected")
	}
	if ee.Values[2].GoName != "ErrTimeout" {
		t.Errorf("Values[2].GoName = %q, want %q", ee.Values[2].GoName, "ErrTimeout")
	}
}

func TestMerge_AutoDeriveGoName(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"MyHandle": {Kind: "opaque_ptr", CType: "MyHandle", GoType: "*C.MyHandle"},
		},
		Functions: map[string]specmodel.FuncDef{
			"MyHandle_setChannelCount": {
				CName: "MyHandle_setChannelCount",
				Params: []specmodel.Param{
					{Name: "h", Type: "*MyHandle"},
					{Name: "count", Type: "int32"},
				},
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "test",
		Package: overlaymodel.PackageOverlay{GoName: "test", GoImport: "github.com/xaionaro-go/ndk/test"},
		Types: map[string]overlaymodel.TypeOverlay{
			"MyHandle": {GoName: "Handle"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"MyHandle_setChannelCount": {
				Receiver: "Handle",
				// No GoName — should auto-derive "SetChannelCount"
				Chain: true,
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	if merged.Methods[0].GoName != "SetChannelCount" {
		t.Errorf("GoName = %q, want %q", merged.Methods[0].GoName, "SetChannelCount")
	}
}

func TestMerge_ParamTypeMapping(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"Handle":      {Kind: "opaque_ptr", CType: "Handle", GoType: "*C.Handle"},
			"Direction_t": {Kind: "typedef_int32", CType: "direction_t", GoType: "int32"},
		},
		Enums: map[string][]specmodel.EnumValue{
			"Direction_t": {{Name: "DIR_OUT", Value: 0}},
		},
		Functions: map[string]specmodel.FuncDef{
			"Handle_setDirection": {
				CName: "Handle_setDirection",
				Params: []specmodel.Param{
					{Name: "h", Type: "*Handle"},
					{Name: "dir", Type: "Direction_t"},
				},
				Returns: "void",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "test",
		Package: overlaymodel.PackageOverlay{GoName: "test", GoImport: "github.com/xaionaro-go/ndk/test"},
		Types: map[string]overlaymodel.TypeOverlay{
			"Handle":      {GoName: "Handle"},
			"Direction_t": {GoName: "Direction", StripPrefix: "DIR_"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"Handle_setDirection": {
				Receiver: "Handle",
				GoName:   "SetDirection",
				Chain:    true,
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	// Param type should be resolved from "Direction_t" to "Direction".
	if len(merged.Methods[0].Params) != 1 {
		t.Fatalf("Params count = %d, want 1", len(merged.Methods[0].Params))
	}
	if merged.Methods[0].Params[0].GoType != "Direction" {
		t.Errorf("Param GoType = %q, want %q", merged.Methods[0].Params[0].GoType, "Direction")
	}
}

func TestMerge_EnumWithoutOverlay_Skipped(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"Orphan_t": {Kind: "typedef_int32", CType: "orphan_t", GoType: "int32"},
		},
		Enums: map[string][]specmodel.EnumValue{
			"Orphan_t": {{Name: "ORPHAN_A", Value: 0}},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "test",
		Package: overlaymodel.PackageOverlay{GoName: "test", GoImport: "github.com/xaionaro-go/ndk/test"},
		// No type overlay for Orphan_t
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.ValueEnums) != 0 {
		t.Errorf("ValueEnums count = %d, want 0 (orphaned enum should be skipped)", len(merged.ValueEnums))
	}
}

func TestMerge_CallbackStructs(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "camera",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/camera",
		Types: map[string]specmodel.TypeDef{
			"ACameraDevice_StateCallbacks": {Kind: "opaque_ptr", CType: "ACameraDevice_StateCallbacks"},
		},
		Structs: map[string]specmodel.StructDef{
			"ACameraDevice_StateCallbacks": {
				Fields: []specmodel.StructField{
					{Name: "context", Type: "void*"},
					{Name: "onDisconnected", Type: "func_ptr", Params: []specmodel.Param{
						{Name: "context", Type: "void*"},
						{Name: "device", Type: "*ACameraDevice"},
					}},
					{Name: "onError", Type: "func_ptr", Params: []specmodel.Param{
						{Name: "context", Type: "void*"},
						{Name: "device", Type: "*ACameraDevice"},
						{Name: "error", Type: "int"},
					}},
				},
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "camera",
		Package: overlaymodel.PackageOverlay{GoName: "camera", GoImport: "github.com/xaionaro-go/ndk/camera"},
		Types: map[string]overlaymodel.TypeOverlay{
			"ACameraDevice_StateCallbacks": {GoName: "DeviceStateCallbacks"},
		},
		CallbackStructs: map[string]overlaymodel.CallbackStructOverlay{
			"ACameraDevice_StateCallbacks": {
				GoName:       "DeviceStateCallbacks",
				ContextField: "context",
				Fields: map[string]overlaymodel.CallbackFieldOverlay{
					"onDisconnected": {GoName: "OnDisconnected", GoSignature: "func()"},
					"onError":        {GoName: "OnError", GoSignature: "func(int)"},
				},
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.CallbackStructs) != 1 {
		t.Fatalf("CallbackStructs count = %d, want 1", len(merged.CallbackStructs))
	}
	cs := merged.CallbackStructs[0]
	if cs.GoName != "DeviceStateCallbacks" {
		t.Errorf("GoName = %q, want %q", cs.GoName, "DeviceStateCallbacks")
	}
	if cs.ContextField != "context" {
		t.Errorf("ContextField = %q", cs.ContextField)
	}
	if len(cs.Fields) != 2 {
		t.Fatalf("Fields count = %d, want 2", len(cs.Fields))
	}
	if cs.Fields[0].CName != "onDisconnected" {
		t.Errorf("Fields[0].CName = %q", cs.Fields[0].CName)
	}
	if cs.Fields[0].GoName != "OnDisconnected" {
		t.Errorf("Fields[0].GoName = %q", cs.Fields[0].GoName)
	}
	if cs.Fields[0].GoSignature != "func()" {
		t.Errorf("Fields[0].GoSignature = %q", cs.Fields[0].GoSignature)
	}
	if len(cs.Fields[0].Params) != 2 {
		t.Fatalf("Fields[0].Params count = %d, want 2", len(cs.Fields[0].Params))
	}
}

func TestMerge_StructAccessors(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "camera",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/camera",
	}
	overlay := overlaymodel.Overlay{
		Module:  "camera",
		Package: overlaymodel.PackageOverlay{GoName: "camera", GoImport: "github.com/xaionaro-go/ndk/camera"},
		StructAccessors: map[string]overlaymodel.StructAccessorOverlay{
			"ACameraIdList": {
				CountField: "numCameras",
				ItemsField: "cameraIds",
				ItemType:   "string",
				DeleteFunc: "ACameraManager_deleteCameraIdList",
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.StructAccessors) != 1 {
		t.Fatalf("StructAccessors count = %d, want 1", len(merged.StructAccessors))
	}
	sa := merged.StructAccessors[0]
	if sa.SpecName != "ACameraIdList" {
		t.Errorf("SpecName = %q", sa.SpecName)
	}
	if sa.CountField != "numCameras" {
		t.Errorf("CountField = %q", sa.CountField)
	}
	if sa.ItemType != "string" {
		t.Errorf("ItemType = %q", sa.ItemType)
	}
	if sa.DeleteFunc != "ACameraManager_deleteCameraIdList" {
		t.Errorf("DeleteFunc = %q", sa.DeleteFunc)
	}
}

func TestMerge_Lifecycle(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "nativeactivity",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/nativeactivity",
		Types: map[string]specmodel.TypeDef{
			"ANativeActivity": {Kind: "opaque_ptr", CType: "ANativeActivity"},
		},
		Structs: map[string]specmodel.StructDef{
			"ANativeActivityCallbacks": {
				Fields: []specmodel.StructField{
					{Name: "onStart", Type: "func_ptr", Params: []specmodel.Param{
						{Name: "activity", Type: "*ANativeActivity"},
					}},
					{Name: "onResume", Type: "func_ptr", Params: []specmodel.Param{
						{Name: "activity", Type: "*ANativeActivity"},
					}},
				},
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "nativeactivity",
		Package: overlaymodel.PackageOverlay{GoName: "activity", GoImport: "github.com/xaionaro-go/ndk/activity"},
		Types: map[string]overlaymodel.TypeOverlay{
			"ANativeActivity": {GoName: "Activity"},
		},
		Lifecycle: &overlaymodel.LifecycleOverlay{
			EntryPoint:        "ANativeActivity_onCreate",
			ActivityType:      "ANativeActivity",
			CallbacksAccessor: "activity->callbacks",
			CallbackStruct:    "ANativeActivityCallbacks",
		},
	}
	merged := idiomgen.Merge(spec, overlay)

	if merged.Lifecycle == nil {
		t.Fatal("Lifecycle is nil")
	}
	if merged.Lifecycle.EntryPoint != "ANativeActivity_onCreate" {
		t.Errorf("EntryPoint = %q", merged.Lifecycle.EntryPoint)
	}
	if merged.Lifecycle.GoActivityType != "Activity" {
		t.Errorf("GoActivityType = %q", merged.Lifecycle.GoActivityType)
	}
	if len(merged.Lifecycle.Fields) != 2 {
		t.Fatalf("Lifecycle fields = %d, want 2", len(merged.Lifecycle.Fields))
	}
	if merged.Lifecycle.Fields[0].CName != "onStart" {
		t.Errorf("Fields[0].CName = %q", merged.Lifecycle.Fields[0].CName)
	}
	if merged.Lifecycle.Fields[0].GoName != "OnStart" {
		t.Errorf("Fields[0].GoName = %q", merged.Lifecycle.Fields[0].GoName)
	}
}

func TestMerge_FixedParams(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "camera",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/camera",
		Types: map[string]specmodel.TypeDef{
			"CaptureSession": {Kind: "opaque_ptr", CType: "CaptureSession"},
		},
		Functions: map[string]specmodel.FuncDef{
			"CaptureSession_setRepeatingRequest": {
				CName: "CaptureSession_setRepeatingRequest",
				Params: []specmodel.Param{
					{Name: "session", Type: "*CaptureSession"},
					{Name: "callbacks", Type: "*CaptureCallbacks"},
					{Name: "numRequests", Type: "int"},
					{Name: "requests", Type: "**CaptureRequest"},
					{Name: "sequenceId", Type: "*int"},
				},
				Returns: "Camera_status_t",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "camera",
		Package: overlaymodel.PackageOverlay{GoName: "camera", GoImport: "github.com/xaionaro-go/ndk/camera"},
		Types: map[string]overlaymodel.TypeOverlay{
			"CaptureSession": {GoName: "CaptureSession"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"CaptureSession_setRepeatingRequest": {
				Receiver: "CaptureSession",
				GoName:   "SetRepeatingRequest",
				FixedParams: map[string]string{
					"callbacks":   "nil",
					"numRequests": "1",
					"sequenceId":  "nil",
				},
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	m := merged.Methods[0]
	if len(m.FixedParams) != 3 {
		t.Fatalf("FixedParams count = %d, want 3", len(m.FixedParams))
	}
	if m.FixedParams["callbacks"] != "nil" {
		t.Errorf("FixedParams[callbacks] = %q", m.FixedParams["callbacks"])
	}
	if m.FixedParams["numRequests"] != "1" {
		t.Errorf("FixedParams[numRequests] = %q", m.FixedParams["numRequests"])
	}
}

func TestMerge_MethodWithCallbackParam(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "camera",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/camera",
		Types: map[string]specmodel.TypeDef{
			"ACameraManager": {Kind: "opaque_ptr", CType: "ACameraManager"},
		},
		Functions: map[string]specmodel.FuncDef{
			"ACameraManager_openCamera": {
				CName: "ACameraManager_openCamera",
				Params: []specmodel.Param{
					{Name: "manager", Type: "*ACameraManager"},
					{Name: "cameraId", Type: "string"},
					{Name: "callbacks", Type: "*ACameraDevice_StateCallbacks"},
					{Name: "device", Type: "**ACameraDevice", Direction: "out"},
				},
				Returns: "Camera_status_t",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "camera",
		Package: overlaymodel.PackageOverlay{GoName: "camera", GoImport: "github.com/xaionaro-go/ndk/camera"},
		Types: map[string]overlaymodel.TypeOverlay{
			"ACameraManager": {GoName: "Manager"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"ACameraManager_openCamera": {
				Receiver:      "Manager",
				GoName:        "OpenCamera",
				ReturnsNew:    "Device",
				CallbackParam: "callbacks",
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)

	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	m := merged.Methods[0]
	if m.CallbackParam != "callbacks" {
		t.Errorf("CallbackParam = %q, want %q", m.CallbackParam, "callbacks")
	}
	if m.CallbackStruct != "ACameraDevice_StateCallbacks" {
		t.Errorf("CallbackStruct = %q, want %q", m.CallbackStruct, "ACameraDevice_StateCallbacks")
	}
}

func TestMerge_ReturnsFrames(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/test",
		Types: map[string]specmodel.TypeDef{
			"Stream": {Kind: "opaque_ptr", CType: "Stream", GoType: "*C.Stream"},
		},
		Functions: map[string]specmodel.FuncDef{
			"Stream_write": {
				CName: "Stream_write",
				Params: []specmodel.Param{
					{Name: "s", Type: "*Stream"},
					{Name: "data", Type: "unsafe.Pointer"},
					{Name: "frames", Type: "int32"},
				},
				Returns: "int32",
			},
		},
	}
	overlay := overlaymodel.Overlay{
		Module:  "test",
		Package: overlaymodel.PackageOverlay{GoName: "test", GoImport: "github.com/xaionaro-go/ndk/test"},
		Types: map[string]overlaymodel.TypeOverlay{
			"Stream": {GoName: "Stream"},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"Stream_write": {
				Receiver:      "Stream",
				GoName:        "Write",
				ReturnsFrames: true,
			},
		},
	}
	merged := idiomgen.Merge(spec, overlay)
	if len(merged.Methods) != 1 {
		t.Fatalf("Methods count = %d, want 1", len(merged.Methods))
	}
	if !merged.Methods[0].ReturnsFrames {
		t.Error("ReturnsFrames = false, want true")
	}
}
