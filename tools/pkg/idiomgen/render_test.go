package idiomgen_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/xaionaro-go/ndk/tools/pkg/idiomgen"
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
		SourcePackage: "github.com/xaionaro-go/ndk/capi/vulkan",
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
		SourcePackage: "github.com/xaionaro-go/ndk/tools/pkg/capi/aaudio",
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
		"func result(",
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
		SourcePackage: "github.com/xaionaro-go/ndk/capi/aaudio",
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

func TestRenderFunctions(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	tmplPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "templates", "functions.go.tmpl")

	spec := idiomgen.MergedSpec{
		PackageName:   "egl",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/egl",
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
