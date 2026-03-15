package specgen

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSource(t *testing.T) {
	spec, err := ParseSources("fake", "github.com/AndroidGoLab/ndk/capi/fake", []string{
		"testdata/simple/simple.go",
	})
	if err != nil {
		t.Fatalf("ParseSources: %v", err)
	}

	// Module and source package.
	if spec.Module != "fake" {
		t.Errorf("module = %q, want %q", spec.Module, "fake")
	}
	if spec.SourcePackage != "github.com/AndroidGoLab/ndk/capi/fake" {
		t.Errorf("source_package = %q, want %q", spec.SourcePackage, "github.com/AndroidGoLab/ndk/capi/fake")
	}

	// Types: at least 4 (2 opaque + 2 typedef).
	if len(spec.Types) < 4 {
		t.Fatalf("types count = %d, want >= 4; got: %v", len(spec.Types), spec.Types)
	}

	// FakeStream is opaque_ptr.
	td, ok := spec.Types["FakeStream"]
	require.True(t, ok, "missing type FakeStream")
	assert.Equal(t, "opaque_ptr", td.Kind, "FakeStream.Kind")

	// FakeBuilder is opaque_ptr.
	td, ok = spec.Types["FakeBuilder"]
	require.True(t, ok, "missing type FakeBuilder")
	assert.Equal(t, "opaque_ptr", td.Kind, "FakeBuilder.Kind")

	// Fake_result_t is typedef_int32.
	td, ok = spec.Types["Fake_result_t"]
	require.True(t, ok, "missing type Fake_result_t")
	assert.Equal(t, "typedef_int32", td.Kind, "Fake_result_t.Kind")

	// Fake_direction_t is typedef_int32.
	td, ok = spec.Types["Fake_direction_t"]
	require.True(t, ok, "missing type Fake_direction_t")
	assert.Equal(t, "typedef_int32", td.Kind, "Fake_direction_t.Kind")

	// Enums: 2 groups with 2 values each.
	if len(spec.Enums) != 2 {
		t.Fatalf("enum groups = %d, want 2; got: %v", len(spec.Enums), spec.Enums)
	}

	resultEnums := spec.Enums["Fake_result_t"]
	if len(resultEnums) != 2 {
		t.Fatalf("Fake_result_t enum values = %d, want 2", len(resultEnums))
	}
	// Check FAKE_OK = 0.
	if resultEnums[0].Name != "FAKE_OK" || resultEnums[0].Value != 0 {
		t.Errorf("resultEnums[0] = %+v, want FAKE_OK=0", resultEnums[0])
	}
	// Check FAKE_ERROR_BASE = -900 (negative value via UnaryExpr).
	if resultEnums[1].Name != "FAKE_ERROR_BASE" || resultEnums[1].Value != -900 {
		t.Errorf("resultEnums[1] = %+v, want FAKE_ERROR_BASE=-900", resultEnums[1])
	}

	dirEnums := spec.Enums["Fake_direction_t"]
	if len(dirEnums) != 2 {
		t.Fatalf("Fake_direction_t enum values = %d, want 2", len(dirEnums))
	}
	if dirEnums[0].Name != "FAKE_DIRECTION_OUTPUT" || dirEnums[0].Value != 0 {
		t.Errorf("dirEnums[0] = %+v, want FAKE_DIRECTION_OUTPUT=0", dirEnums[0])
	}
	if dirEnums[1].Name != "FAKE_DIRECTION_INPUT" || dirEnums[1].Value != 1 {
		t.Errorf("dirEnums[1] = %+v, want FAKE_DIRECTION_INPUT=1", dirEnums[1])
	}

	// Functions: at least 7.
	if len(spec.Functions) < 7 {
		t.Fatalf("functions count = %d, want >= 7; got: %v", len(spec.Functions), spec.Functions)
	}

	// Fake_createBuilder: returns Fake_result_t, 1 out param.
	fn, ok := spec.Functions["Fake_createBuilder"]
	if !ok {
		t.Fatal("missing function Fake_createBuilder")
	}
	if fn.Returns != "Fake_result_t" {
		t.Errorf("Fake_createBuilder.Returns = %q, want Fake_result_t", fn.Returns)
	}
	if len(fn.Params) != 1 {
		t.Fatalf("Fake_createBuilder params = %d, want 1", len(fn.Params))
	}
	if fn.Params[0].Direction != "out" {
		t.Errorf("Fake_createBuilder param direction = %q, want out", fn.Params[0].Direction)
	}
	if fn.Params[0].Name != "builder" {
		t.Errorf("Fake_createBuilder param name = %q, want builder", fn.Params[0].Name)
	}

	// FakeBuilder_setDeviceId: no return, 2 params.
	fn2, ok := spec.Functions["FakeBuilder_setDeviceId"]
	if !ok {
		t.Fatal("missing function FakeBuilder_setDeviceId")
	}
	if fn2.Returns != "" {
		t.Errorf("FakeBuilder_setDeviceId.Returns = %q, want empty", fn2.Returns)
	}
	if len(fn2.Params) != 2 {
		t.Fatalf("FakeBuilder_setDeviceId params = %d, want 2", len(fn2.Params))
	}

	// FakeBuilder_openStream: returns Fake_result_t, 2 params, second is out.
	fn3, ok := spec.Functions["FakeBuilder_openStream"]
	if !ok {
		t.Fatal("missing function FakeBuilder_openStream")
	}
	if fn3.Returns != "Fake_result_t" {
		t.Errorf("FakeBuilder_openStream.Returns = %q, want Fake_result_t", fn3.Returns)
	}
	if len(fn3.Params) != 2 {
		t.Fatalf("FakeBuilder_openStream params = %d, want 2", len(fn3.Params))
	}
	if fn3.Params[1].Direction != "out" {
		t.Errorf("FakeBuilder_openStream param[1].Direction = %q, want out", fn3.Params[1].Direction)
	}

	// FakeStream_getSampleRate: returns int32.
	fn4, ok := spec.Functions["FakeStream_getSampleRate"]
	if !ok {
		t.Fatal("missing function FakeStream_getSampleRate")
	}
	if fn4.Returns != "int32" {
		t.Errorf("FakeStream_getSampleRate.Returns = %q, want int32", fn4.Returns)
	}

	// Callbacks: FakeStream_dataCallback.
	if len(spec.Callbacks) < 1 {
		t.Fatalf("callbacks count = %d, want >= 1", len(spec.Callbacks))
	}
	cb, ok := spec.Callbacks["FakeStream_dataCallback"]
	if !ok {
		t.Fatal("missing callback FakeStream_dataCallback")
	}
	if cb.Returns != "Fake_result_t" {
		t.Errorf("FakeStream_dataCallback.Returns = %q, want Fake_result_t", cb.Returns)
	}
	if len(cb.Params) != 4 {
		t.Fatalf("FakeStream_dataCallback params = %d, want 4", len(cb.Params))
	}
}

func TestClassifyIdent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"int32", "typedef_int32"},
		{"uint32", "typedef_uint32"},
		{"int64", "typedef_int64"},
		{"uint64", "typedef_uint64"},
		{"int", "typedef_int"},
		{"uint", "typedef_uint"},
		{"int8", "typedef_int8"},
		{"uint8", "typedef_uint8"},
		{"int16", "typedef_int16"},
		{"uint16", "typedef_uint16"},
		{"float32", "typedef_float32"},
		{"float64", "typedef_float64"},
		{"string", ""},
		{"bool", ""},
		{"MyCustomType", ""},
	}
	for _, tt := range tests {
		got := classifyIdent(tt.input)
		if got != tt.want {
			t.Errorf("classifyIdent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want string
	}{
		{
			name: "ident",
			expr: &ast.Ident{Name: "int32"},
			want: "int32",
		},
		{
			name: "star",
			expr: &ast.StarExpr{X: &ast.Ident{Name: "Foo"}},
			want: "*Foo",
		},
		{
			name: "double star",
			expr: &ast.StarExpr{X: &ast.StarExpr{X: &ast.Ident{Name: "Foo"}}},
			want: "**Foo",
		},
		{
			name: "selector",
			expr: &ast.SelectorExpr{X: &ast.Ident{Name: "unsafe"}, Sel: &ast.Ident{Name: "Pointer"}},
			want: "unsafe.Pointer",
		},
		{
			name: "slice",
			expr: &ast.ArrayType{Elt: &ast.Ident{Name: "byte"}},
			want: "[]byte",
		},
		{
			name: "map",
			expr: &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}},
			want: "map[string]int",
		},
		{
			name: "interface",
			expr: &ast.InterfaceType{Methods: &ast.FieldList{}},
			want: "interface{}",
		},
		{
			name: "ellipsis",
			expr: &ast.Ellipsis{Elt: &ast.Ident{Name: "int"}},
			want: "...int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typeString(tt.expr)
			if got != tt.want {
				t.Errorf("typeString = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseSourceEdgeCases(t *testing.T) {
	spec, err := ParseSources("edge", "github.com/AndroidGoLab/ndk/capi/edge", []string{
		"testdata/edgecases/edgecases.go",
	})
	if err != nil {
		t.Fatalf("ParseSources: %v", err)
	}

	// No enums from untyped or string const blocks.
	if len(spec.Enums) != 0 {
		t.Errorf("enums = %v, want empty", spec.Enums)
	}

	// Unexported function should not appear.
	if _, ok := spec.Functions["unexported"]; ok {
		t.Error("unexported function should not be extracted")
	}

	// Multi-return function: Returns should be empty (we don't capture multi-return).
	if fn, ok := spec.Functions["Edge_multiReturn"]; ok {
		if fn.Returns != "" {
			t.Errorf("Edge_multiReturn.Returns = %q, want empty", fn.Returns)
		}
	} else {
		t.Error("missing function Edge_multiReturn")
	}

	// c-for-go ref-helper (New* prefix) should be filtered.
	require.NotContains(t, spec.Functions, "NewFakeStreamRef")

	// Unnamed param function.
	if fn, ok := spec.Functions["Edge_unnamedParam"]; ok {
		if len(fn.Params) != 1 {
			t.Fatalf("Edge_unnamedParam params = %d, want 1", len(fn.Params))
		}
		if fn.Params[0].Name != "" {
			t.Errorf("Edge_unnamedParam param name = %q, want empty", fn.Params[0].Name)
		}
		if fn.Params[0].Type != "int32" {
			t.Errorf("Edge_unnamedParam param type = %q, want int32", fn.Params[0].Type)
		}
	} else {
		t.Error("missing function Edge_unnamedParam")
	}
}

func TestParseSourcesError(t *testing.T) {
	_, err := ParseSources("bad", "pkg", []string{"nonexistent.go"})
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestEvalConstInt(t *testing.T) {
	// Test unsupported expression type.
	_, ok := evalConstInt(&ast.Ident{Name: "iota"})
	if ok {
		t.Error("expected false for unsupported expression")
	}

	// Test non-SUB unary.
	_, ok = evalConstInt(&ast.UnaryExpr{
		Op: 17, // some non-SUB op
		X:  &ast.BasicLit{Kind: 5, Value: "1"},
	})
	if ok {
		t.Error("expected false for non-SUB unary")
	}
}

func TestMergeStructs_IncludesAllStructs(t *testing.T) {
	spec := specmodel.Spec{
		Types: map[string]specmodel.TypeDef{
			"ACameraDevice_StateCallbacks": {
				Kind:   "opaque_ptr",
				CType:  "ACameraDevice_StateCallbacks",
				GoType: "*C.ACameraDevice_StateCallbacks",
			},
		},
	}

	structs := map[string]specmodel.StructDef{
		"ACameraDevice_StateCallbacks": {
			Fields: []specmodel.StructField{
				{Name: "context", Type: "void*"},
				{Name: "onDisconnected", Type: "func_ptr"},
			},
		},
		"UnrelatedStruct": {
			Fields: []specmodel.StructField{
				{Name: "x", Type: "int"},
			},
		},
	}

	MergeStructs(&spec, structs)

	require.Len(t, spec.Structs, 2)
	require.Contains(t, spec.Structs, "ACameraDevice_StateCallbacks")
	require.Contains(t, spec.Structs, "UnrelatedStruct")
	assert.Len(t, spec.Structs["ACameraDevice_StateCallbacks"].Fields, 2)
}

func TestMergeStructs_NilStructsMap(t *testing.T) {
	spec := specmodel.Spec{
		Types: map[string]specmodel.TypeDef{
			"Foo": {Kind: "opaque_ptr"},
		},
	}

	structs := map[string]specmodel.StructDef{
		"Foo": {Fields: []specmodel.StructField{{Name: "x", Type: "int"}}},
	}

	MergeStructs(&spec, structs)

	require.NotNil(t, spec.Structs)
	require.Contains(t, spec.Structs, "Foo")
}

func TestMergeStructs_EmptyStructs(t *testing.T) {
	spec := specmodel.Spec{
		Types: map[string]specmodel.TypeDef{
			"Foo": {Kind: "opaque_ptr"},
		},
	}

	MergeStructs(&spec, map[string]specmodel.StructDef{})

	require.NotNil(t, spec.Structs)
	assert.Empty(t, spec.Structs)
}

func TestExtractIncludeDirsFromGoFiles(t *testing.T) {
	// Create a temporary sysroot with a header directory.
	sysroot := t.TempDir()
	cameraDir := filepath.Join(sysroot, "camera")
	require.NoError(t, os.MkdirAll(cameraDir, 0o755))

	// Create a Go file with CGo includes.
	goDir := t.TempDir()
	goFile := filepath.Join(goDir, "camera.go")
	content := `package camera

// #include <camera/NdkCameraDevice.h>
// #include <camera/NdkCameraManager.h>
import "C"
`
	require.NoError(t, os.WriteFile(goFile, []byte(content), 0o644))

	dirs, err := ExtractIncludeDirsFromGoFiles([]string{goFile}, sysroot)
	require.NoError(t, err)

	require.Len(t, dirs, 1)
	assert.Equal(t, cameraDir, dirs[0])
}

func TestExtractIncludeDirsFromGoFiles_NoCGoIncludes(t *testing.T) {
	goDir := t.TempDir()
	goFile := filepath.Join(goDir, "noop.go")
	content := `package noop

func Noop() {}
`
	require.NoError(t, os.WriteFile(goFile, []byte(content), 0o644))

	dirs, err := ExtractIncludeDirsFromGoFiles([]string{goFile}, "/nonexistent")
	require.NoError(t, err)
	assert.Empty(t, dirs)
}
