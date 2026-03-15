package idiomgen_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/AndroidGoLab/ndk/tools/pkg/idiomgen"
)

func TestRenderTemplate_Enum(t *testing.T) {
	// Test that the inline enum template still works for basic rendering.
	spec := idiomgen.MergedSpec{
		PackageName: "audio",
		ValueEnums: []idiomgen.MergedValueEnum{
			{
				GoName:       "Direction",
				SpecName:     "aaudio_direction_t",
				StripPrefix:  "AAUDIO_DIRECTION_",
				StringMethod: true,
				Values: []idiomgen.MergedEnumValue{
					{GoName: "DirectionOutput", SpecName: "AAUDIO_DIRECTION_OUTPUT", Value: 0},
					{GoName: "DirectionInput", SpecName: "AAUDIO_DIRECTION_INPUT", Value: 1},
				},
			},
		},
	}

	// Use a simple inline template to test RenderTemplate still works.
	tmpl := `package {{ .PackageName }}
{{ range $enum := .ValueEnums }}
type {{ $enum.GoName }} int32
{{ end }}
`
	out, err := idiomgen.RenderTemplate("test.go.tmpl", tmpl, spec)
	if err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}
	if !strings.Contains(out, "type Direction int32") {
		t.Errorf("output missing 'type Direction int32'\n\ngot:\n%s", out)
	}
}

func TestRenderValueEnumFile(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "value_enum_file.go.tmpl")

	data := idiomgen.PerValueEnumData{
		PackageName: "audio",
		Enum: idiomgen.MergedValueEnum{
			GoName:       "Direction",
			SpecName:     "aaudio_direction_t",
			BaseType:     "int32",
			StripPrefix:  "AAUDIO_DIRECTION_",
			StringMethod: true,
			Values: []idiomgen.MergedEnumValue{
				{GoName: "Output", SpecName: "AAUDIO_DIRECTION_OUTPUT", Value: 0},
				{GoName: "Input", SpecName: "AAUDIO_DIRECTION_INPUT", Value: 1},
			},
		},
	}

	out, err := idiomgen.RenderAny(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderAny: %v", err)
	}

	for _, want := range []string{
		"package audio",
		"type Direction int32",
		"Output Direction = 0",
		"Input Direction = 1",
		"func (v Direction) String() string",
		`return "Output"`,
		`return "Input"`,
		`return fmt.Sprintf("Direction(%d)", int(v))`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}
}

func TestRenderValueEnumFile_WithoutStringMethod(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "value_enum_file.go.tmpl")

	data := idiomgen.PerValueEnumData{
		PackageName: "audio",
		Enum: idiomgen.MergedValueEnum{
			GoName:       "Format",
			SpecName:     "aaudio_format_t",
			BaseType:     "int32",
			StringMethod: false,
			Values: []idiomgen.MergedEnumValue{
				{GoName: "FormatPCMI16", Value: 1},
				{GoName: "FormatPCMFloat", Value: 2},
			},
		},
	}

	out, err := idiomgen.RenderAny(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderAny: %v", err)
	}

	if strings.Contains(out, "func (v Format) String()") {
		t.Errorf("String() method should not be generated when StringMethod is false\n\ngot:\n%s", out)
	}
	if strings.Contains(out, `import "fmt"`) {
		t.Errorf("import fmt should not be generated when StringMethod is false\n\ngot:\n%s", out)
	}
	if !strings.Contains(out, "type Format int32") {
		t.Errorf("output missing type declaration\n\ngot:\n%s", out)
	}
}

func TestRenderTypeAliasFile(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "type_alias_file.go.tmpl")

	data := idiomgen.PerTypeAliasData{
		PackageName:   "vulkan",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/vulkan",
		Alias:         idiomgen.MergedTypeAlias{GoName: "VkInstance", CapiType: "VkInstance"},
	}

	out, err := idiomgen.RenderAny(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderAny: %v", err)
	}

	for _, want := range []string{
		"package vulkan",
		"type VkInstance = capi.VkInstance",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}

	if strings.Contains(out, "struct") {
		t.Errorf("type alias file should not contain struct definitions\n\ngot:\n%s", out)
	}
}

func TestRenderCallbackFile(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "callback_file.go.tmpl")

	data := idiomgen.PerCallbackData{
		PackageName: "audio",
		Callback: idiomgen.MergedCallback{
			SpecName:       "AAudioStream_dataCallback",
			GoCallbackType: "DataCallback",
			GoCallbackSig:  "func(stream *Stream, buf []byte, numFrames int) DataCallbackResult",
		},
	}

	out, err := idiomgen.RenderAny(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderAny: %v", err)
	}

	want := "type DataCallback func(stream *Stream, buf []byte, numFrames int) DataCallbackResult"
	if !strings.Contains(out, want) {
		t.Errorf("output missing %q\n\ngot:\n%s", want, out)
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"StreamBuilder", "stream_builder"},
		{"Stream", "stream"},
		{"CaptureRequest", "capture_request"},
		{"Model", "model"},
		{"Manager", "manager"},
		{"SessionOutputContainer", "session_output_container"},
		{"AUDIO", "audio"},
		{"LOOPER_POLL", "looper_poll"},
		{"IMAGE_FORMATS", "image_formats"},
		{"NDROID_BITMAP_FLAGS", "ndroid_bitmap_flags"},
		{"STATUS", "status"},
		{"", ""},
	}
	fm := idiomgen.FuncMap()
	toSnakeCase := fm["toSnakeCase"].(func(string) string)
	for _, tt := range tests {
		if got := toSnakeCase(tt.in); got != tt.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFixGoAcronyms(t *testing.T) {
	fm := idiomgen.FuncMap()
	safeGoName := fm["safeGoName"].(func(string) string)
	tests := []struct {
		in, want string
	}{
		{"deviceId", "deviceID"},
		{"sessionId", "sessionID"},
		{"SetDeviceId", "SetDeviceID"},
		{"CameraIdList", "CameraIDList"},
		{"VsyncId", "VsyncID"},
		// Must NOT change words where "Id" is followed by lowercase.
		{"Identify", "Identify"},
		{"Idle", "Idle"},
		{"Identity", "Identity"},
		{"DeviceWaitIdle", "DeviceWaitIdle"},
		// Plain names without "Id" are unchanged.
		{"count", "count"},
		{"StreamBuilder", "StreamBuilder"},
	}
	for _, tt := range tests {
		if got := safeGoName(tt.in); got != tt.want {
			t.Errorf("safeGoName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFuncMapUnexport(t *testing.T) {
	fm := idiomgen.FuncMap()
	unexport := fm["unexport"].(func(string) string)

	tests := []struct {
		in, want string
	}{
		{"Direction", "direction"},
		{"ABC", "aBC"},
		{"", ""},
		{"a", "a"},
	}
	for _, tt := range tests {
		if got := unexport(tt.in); got != tt.want {
			t.Errorf("unexport(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFuncMapSkip(t *testing.T) {
	fm := idiomgen.FuncMap()
	skip := fm["skip"].(func([]idiomgen.MergedParam, int) []idiomgen.MergedParam)

	params := []idiomgen.MergedParam{
		{Name: "self", GoType: "*StreamBuilder"},
		{Name: "rate", GoType: "int32"},
		{Name: "channels", GoType: "int32"},
	}

	got := skip(params, 1)
	if len(got) != 2 {
		t.Fatalf("skip(params, 1) returned %d elements, want 2", len(got))
	}
	if got[0].Name != "rate" {
		t.Errorf("skip(params, 1)[0].Name = %q, want %q", got[0].Name, "rate")
	}

	// skip beyond length returns nil
	got = skip(params, 5)
	if got != nil {
		t.Errorf("skip(params, 5) = %v, want nil", got)
	}
}

func TestFuncMapSkipAndFilterOut(t *testing.T) {
	fm := idiomgen.FuncMap()
	skipAndFilterOut := fm["skipAndFilterOut"].(func([]idiomgen.MergedParam) []idiomgen.MergedParam)

	// Receiver is already stripped by the merger, so params here don't include it.
	params := []idiomgen.MergedParam{
		{Name: "rate", GoType: "int32", Direction: "in"},
		{Name: "result", GoType: "*int32", Direction: "out"},
		{Name: "channels", GoType: "int32", Direction: "in"},
	}

	got := skipAndFilterOut(params)
	if len(got) != 2 {
		t.Fatalf("skipAndFilterOut returned %d elements, want 2", len(got))
	}
	if got[0].Name != "rate" || got[1].Name != "channels" {
		t.Errorf("skipAndFilterOut = %v, want [rate, channels]", got)
	}
}

func TestFuncMapLookupCapiType(t *testing.T) {
	fm := idiomgen.FuncMap()
	lookup := fm["lookupCapiType"].(func(map[string]idiomgen.MergedOpaqueType, string) string)

	types := map[string]idiomgen.MergedOpaqueType{
		"StreamBuilder": {CapiType: "AAudioStreamBuilder"},
	}

	if got := lookup(types, "StreamBuilder"); got != "AAudioStreamBuilder" {
		t.Errorf("lookupCapiType(StreamBuilder) = %q, want %q", got, "AAudioStreamBuilder")
	}
	if got := lookup(types, "Unknown"); got != "Unknown" {
		t.Errorf("lookupCapiType(Unknown) = %q, want %q", got, "Unknown")
	}
}

func TestRenderErrors(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "errors.go.tmpl")

	spec := idiomgen.MergedSpec{
		PackageName:   "audio",
		SourcePackage: "github.com/AndroidGoLab/ndk/tools/pkg/capi/aaudio",
		ErrorEnums: []idiomgen.MergedErrorEnum{
			{
				GoName:       "Aaudio_result_t",
				Prefix:       "audio",
				SuccessValue: "AAUDIO_OK",
				Values: []idiomgen.MergedEnumValue{
					{GoName: "AAUDIO_OK", SpecName: "AAUDIO_OK", Value: 0},
					{GoName: "ErrorBase", SpecName: "AAUDIO_ERROR_BASE", Value: -900},
					{GoName: "ErrorDisconnected", SpecName: "AAUDIO_ERROR_DISCONNECTED", Value: -899},
				},
			},
		},
	}

	out, err := idiomgen.RenderTemplateFile(tmplPath, spec)
	if err != nil {
		t.Fatalf("RenderTemplateFile: %v", err)
	}

	for _, want := range []string{
		"type Error int32",
		"func (e Error) Error() string",
		"func result[T ~int32 | ~uint32 | ~int64](",
		"ErrorBase Error = -900",
		"ErrorDisconnected Error = -899",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}

	// The success value must NOT appear in the const block.
	if strings.Contains(out, "AAUDIO_OK Error =") {
		t.Errorf("success value AAUDIO_OK should not appear in const block\n\ngot:\n%s", out)
	}
}

func TestRenderPerTypeFile(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "type_file.go.tmpl")

	data := idiomgen.PerTypeData{
		PackageName:   "audio",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/aaudio",
		Type: idiomgen.MergedOpaqueType{
			GoName:      "StreamBuilder",
			CapiType:    "AAudioStreamBuilder",
			Constructor: "AAudioStreamBuilder_create",
			Destructor:  "AAudioStreamBuilder_delete",
		},
		Methods: []idiomgen.MergedMethod{
			{
				GoName:       "SetDeviceID",
				CName:        "AAudioStreamBuilder_setDeviceId",
				ReceiverType: "StreamBuilder",
				Chain:        true,
				Params: []idiomgen.MergedParam{
					{Name: "deviceID", GoType: "int32", Direction: "in"},
				},
			},
		},
		OpaqueTypes: map[string]idiomgen.MergedOpaqueType{
			"StreamBuilder": {CapiType: "AAudioStreamBuilder"},
		},
	}

	out, err := idiomgen.RenderPerType(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderPerType: %v", err)
	}

	for _, want := range []string{
		"package audio",
		"type StreamBuilder struct",
		"ptr *capi.AAudioStreamBuilder",
		"NewStreamBuilder",
		"func (h *StreamBuilder) Close() error",
		"capi.AAudioStreamBuilder_delete(h.ptr)",
		"func (h *StreamBuilder) SetDeviceID(deviceID int32) *StreamBuilder",
		"return h",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}
}

func TestRenderPerTypeFile_OutputParams(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "type_file.go.tmpl")

	data := idiomgen.PerTypeData{
		PackageName:   "media",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/media",
		Type: idiomgen.MergedOpaqueType{
			GoName:     "ImageReader",
			CapiType:   "AImageReader",
			Destructor: "AImageReader_delete",
		},
		Methods: []idiomgen.MergedMethod{
			{
				GoName:       "AcquireNextImage",
				CName:        "AImageReader_acquireNextImage",
				ReceiverType: "ImageReader",
				Params: []idiomgen.MergedParam{
					{Name: "image", GoType: "**Image", Direction: "out"},
				},
				OutputParams: []idiomgen.MergedOutputParam{
					{CParamName: "image", GoType: "*Image", CapiType: "*capi.AImage", IsHandle: true},
				},
			},
		},
		FreeFunctions: []idiomgen.MergedFreeFunction{
			{
				GoName:     "NewImageReader",
				CName:      "AImageReader_new",
				ReturnsNew: "ImageReader",
				Params: []idiomgen.MergedParam{
					{Name: "width", GoType: "int32"},
					{Name: "height", GoType: "int32"},
					{Name: "reader", GoType: "**ImageReader", Direction: "out"},
				},
				OutputParams: []idiomgen.MergedOutputParam{
					{CParamName: "reader", GoType: "*ImageReader", CapiType: "*capi.AImageReader", IsHandle: true},
				},
			},
		},
		OpaqueTypes: map[string]idiomgen.MergedOpaqueType{
			"ImageReader": {CapiType: "AImageReader"},
			"Image":       {CapiType: "AImage"},
		},
	}

	out, err := idiomgen.RenderPerType(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderPerType: %v", err)
	}

	for _, want := range []string{
		// Method with output_params
		"func (h *ImageReader) AcquireNextImage() (*Image, error)",
		"var imagePtr *capi.AImage",
		"ret := capi.AImageReader_acquireNextImage(h.ptr, &imagePtr)",
		"return nil, err",
		"return &Image{ptr: imagePtr}, nil",
		// Free function with output_params
		"func NewImageReader(width int32, height int32) (*ImageReader, error)",
		"var readerPtr *capi.AImageReader",
		"ret := capi.AImageReader_new(width, height, &readerPtr)",
		"return &ImageReader{ptr: readerPtr}, nil",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}

	// Verify output params are excluded from Go signature.
	if strings.Contains(out, "AcquireNextImage(image") {
		t.Errorf("output param 'image' should NOT appear in method signature\n\ngot:\n%s", out)
	}
}

func TestRenderPerTypeFile_OutputParams_Scalar(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "type_file.go.tmpl")

	data := idiomgen.PerTypeData{
		PackageName:   "media",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/media",
		Type: idiomgen.MergedOpaqueType{
			GoName:     "Image",
			CapiType:   "AImage",
			Destructor: "AImage_delete",
		},
		Methods: []idiomgen.MergedMethod{
			{
				GoName:       "PlaneData",
				CName:        "AImage_getPlaneData",
				ReceiverType: "Image",
				Params: []idiomgen.MergedParam{
					{Name: "planeIdx", GoType: "int32"},
					{Name: "data", GoType: "*uint8", Direction: "out"},
					{Name: "dataLength", GoType: "int32", Direction: "out"},
				},
				OutputParams: []idiomgen.MergedOutputParam{
					{CParamName: "data", GoType: "*uint8", CapiType: "*uint8", IsHandle: false},
					{CParamName: "dataLength", GoType: "int32", CapiType: "int32", IsHandle: false},
				},
			},
		},
		OpaqueTypes: map[string]idiomgen.MergedOpaqueType{
			"Image": {CapiType: "AImage"},
		},
	}

	out, err := idiomgen.RenderPerType(tmplPath, data)
	if err != nil {
		t.Fatalf("RenderPerType: %v", err)
	}

	for _, want := range []string{
		"func (h *Image) PlaneData(planeIdx int32) (*uint8, int32, error)",
		"var dataPtr *uint8",
		"var dataLengthPtr int32",
		"ret := capi.AImage_getPlaneData(h.ptr, planeIdx, &dataPtr, &dataLengthPtr)",
		"return nil, 0, err",
		"return dataPtr, dataLengthPtr, nil",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}
}

func TestRenderFunctions(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "functions.go.tmpl")

	spec := idiomgen.MergedSpec{
		PackageName:   "egl",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/egl",
		FreeFunctions: []idiomgen.MergedFreeFunction{
			{
				GoName:  "GetDisplay",
				CName:   "eglGetDisplay",
				Returns: "EGLDisplay",
				Params: []idiomgen.MergedParam{
					{Name: "displayID", GoType: "unsafe.Pointer"},
				},
			},
			{
				GoName:  "Terminate",
				CName:   "eglTerminate",
				Returns: "EGLBoolean",
				Params: []idiomgen.MergedParam{
					{Name: "display", GoType: "EGLDisplay"},
				},
			},
		},
	}

	out, err := idiomgen.RenderTemplateFile(tmplPath, spec)
	if err != nil {
		t.Fatalf("RenderTemplateFile: %v", err)
	}

	for _, want := range []string{
		"func GetDisplay(displayID unsafe.Pointer) EGLDisplay",
		"func Terminate(display EGLDisplay) EGLBoolean",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\ngot:\n%s", want, out)
		}
	}

	// functions.go should NOT contain methods (those are in per-type files now).
	if strings.Contains(out, "func (h *") {
		t.Errorf("functions.go should not contain methods\n\ngot:\n%s", out)
	}
}
