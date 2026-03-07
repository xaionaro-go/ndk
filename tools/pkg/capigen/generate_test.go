package capigen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/ndk/tools/pkg/specmodel"
)

func looperSpec() *specmodel.Spec {
	return &specmodel.Spec{
		Module:        "looper",
		SourcePackage: "github.com/xaionaro-go/ndk/capi/looper",
		Types: map[string]specmodel.TypeDef{
			"ALooper": {
				Kind:   "opaque_ptr",
				CType:  "ALooper",
				GoType: "*C.ALooper",
			},
		},
		Enums: map[string][]specmodel.EnumValue{
			"ALOOPER_POLL": {
				{Name: "ALOOPER_POLL_WAKE", Value: -1},
				{Name: "ALOOPER_POLL_CALLBACK", Value: -2},
				{Name: "ALOOPER_POLL_TIMEOUT", Value: -3},
				{Name: "ALOOPER_POLL_ERROR", Value: -4},
			},
			"ALOOPER_EVENT": {
				{Name: "ALOOPER_EVENT_INPUT", Value: 1},
				{Name: "ALOOPER_EVENT_OUTPUT", Value: 2},
				{Name: "ALOOPER_EVENT_ERROR", Value: 4},
			},
		},
		Functions: map[string]specmodel.FuncDef{
			"ALooper_forThread": {
				CName:   "ALooper_forThread",
				Returns: "*ALooper",
			},
			"ALooper_addFd": {
				CName: "ALooper_addFd",
				Params: []specmodel.Param{
					{Name: "looper", Type: "*ALooper"},
					{Name: "fd", Type: "int32"},
					{Name: "ident", Type: "int32"},
					{Name: "events", Type: "int32"},
					{Name: "callback", Type: "ALooper_callbackFunc"},
					{Name: "data", Type: "unsafe.Pointer"},
				},
				Returns: "int32",
			},
			"ALooper_acquire": {
				CName: "ALooper_acquire",
				Params: []specmodel.Param{
					{Name: "looper", Type: "*ALooper"},
				},
			},
		},
		Callbacks: map[string]specmodel.CallbackDef{
			"ALooper_callbackFunc": {
				Params: []specmodel.Param{
					{Name: "fd", Type: "int32"},
					{Name: "events", Type: "int32"},
					{Name: "data", Type: "unsafe.Pointer"},
				},
				Returns: "int32",
			},
		},
	}
}

func looperManifest() *Manifest {
	m := &Manifest{}
	m.Generator.PackageName = "looper"
	m.Generator.PackageDescription = "Raw CGo bindings for Android looper"
	m.Generator.Includes = []string{"android/looper.h"}
	m.Generator.FlagGroups = []FlagGroup{
		{Name: "LDFLAGS", Flags: []string{"-landroid"}},
	}
	return m
}

func TestGeneratePackage(t *testing.T) {
	spec := looperSpec()
	manifest := looperManifest()

	outDir := t.TempDir()
	err := GeneratePackage(spec, manifest, outDir)
	require.NoError(t, err)

	// Verify all expected files are created.
	expectedFiles := []string{
		"doc.go",
		"types.go",
		"const.go",
		"cgo_helpers.go",
		"cgo_helpers.h",
		"looper.go",
	}
	for _, f := range expectedFiles {
		path := filepath.Join(outDir, f)
		_, err := os.Stat(path)
		assert.NoError(t, err, "missing file: %s", f)
	}
}

func TestGenerateDocGo(t *testing.T) {
	content := generateDocGo("looper", "Raw CGo bindings for Android looper")
	assert.Contains(t, content, "package looper")
	assert.Contains(t, content, "DO NOT EDIT")
	assert.Contains(t, content, "Raw CGo bindings for Android looper")
}

func TestGenerateTypesGo(t *testing.T) {
	spec := looperSpec()
	preamble := buildCGoPreamble(looperManifest())
	callbackSet := map[string]bool{"ALooper_callbackFunc": true}

	content := generateTypesGo("looper", preamble, spec, callbackSet, nil, nil)

	assert.Contains(t, content, "type ALooper C.ALooper")
	assert.Contains(t, content, "type ALooper_callbackFunc func(")
	assert.Contains(t, content, "fd int32")
	assert.Contains(t, content, "data unsafe.Pointer")
	assert.Contains(t, content, ") int32")
}

func TestGenerateConstGo(t *testing.T) {
	spec := looperSpec()
	preamble := buildCGoPreamble(looperManifest())
	enumTypedefSet := buildEnumTypedefSet(spec)

	content := generateConstGo("looper", preamble, spec, enumTypedefSet)

	assert.Contains(t, content, "ALOOPER_POLL_WAKE = -1")
	assert.Contains(t, content, "ALOOPER_POLL_CALLBACK = -2")
	assert.Contains(t, content, "ALOOPER_EVENT_INPUT = 1")
}

func TestGenerateConstGoTypedEnum(t *testing.T) {
	spec := &specmodel.Spec{
		Types: map[string]specmodel.TypeDef{
			"camera_status_t": {
				Kind:   "typedef_int32",
				CType:  "camera_status_t",
				GoType: "int32",
			},
		},
		Enums: map[string][]specmodel.EnumValue{
			"camera_status_t": {
				{Name: "ACAMERA_OK", Value: 0},
				{Name: "ACAMERA_ERROR_UNKNOWN", Value: -10000},
			},
		},
	}
	enumTypedefSet := buildEnumTypedefSet(spec)
	content := generateConstGo("camera", "", spec, enumTypedefSet)

	assert.Contains(t, content, "ACAMERA_OK Camera_status_t = 0")
	assert.Contains(t, content, "ACAMERA_ERROR_UNKNOWN Camera_status_t = -10000")
}

func TestGenerateCgoHelpersH(t *testing.T) {
	spec := looperSpec()
	manifest := looperManifest()

	content := generateCgoHelpersH("looper", manifest, spec, nil)

	assert.Contains(t, content, "#include \"android/looper.h\"")
	assert.Contains(t, content, "#pragma once")
	assert.Contains(t, content, "#define __CGOGEN 1")
	// Callback proxy declaration.
	assert.Contains(t, content, "ALooper_callbackFunc_")
	assert.Contains(t, content, "int fd")
	assert.Contains(t, content, "void* data")
}

func TestGenerateFunctionsGo(t *testing.T) {
	spec := looperSpec()
	callbackSet := map[string]bool{"ALooper_callbackFunc": true}
	preamble := buildCGoPreamble(looperManifest())

	content := generateFunctionsGo("looper", preamble, spec, callbackSet, nil, nil)

	// ALooper_forThread returns *ALooper.
	assert.Contains(t, content, "func ALooper_forThread() *ALooper")
	assert.Contains(t, content, "C.ALooper_forThread()")
	assert.Contains(t, content, "unsafe.Pointer(&__ret)")

	// ALooper_addFd has parameters.
	assert.Contains(t, content, "func ALooper_addFd(")
	assert.Contains(t, content, "looper *ALooper")
	assert.Contains(t, content, "fd int32")
	assert.Contains(t, content, "callback ALooper_callbackFunc")
	assert.Contains(t, content, "data unsafe.Pointer")
	assert.Contains(t, content, "C.ALooper_addFd(")
	assert.Contains(t, content, "(*C.ALooper)(unsafe.Pointer(looper))")
	assert.Contains(t, content, ".PassValue()")
	assert.Contains(t, content, "runtime.KeepAlive(")

	// ALooper_acquire has no return.
	assert.Contains(t, content, "func ALooper_acquire(")
}

func TestCallbackHashSuffix(t *testing.T) {
	hash := callbackHashSuffix("looper", "ALooper_callbackFunc")
	// Hash is md5-based, not CRC32 like c-for-go. Both sides (caller and
	// proxy) are generated by capigen so consistency is all that matters.
	assert.Len(t, hash, 8)

	// Same callback name in different packages must produce different hashes.
	hash2 := callbackHashSuffix("otherpackage", "ALooper_callbackFunc")
	assert.NotEqual(t, hash, hash2, "different packages must produce different hashes")
}

func TestBuildCGoPreamble(t *testing.T) {
	manifest := looperManifest()
	preamble := buildCGoPreamble(manifest)

	assert.Contains(t, preamble, "#cgo LDFLAGS: -landroid")
	assert.Contains(t, preamble, "#include \"android/looper.h\"")
	assert.Contains(t, preamble, "#include <stdlib.h>")
	assert.Contains(t, preamble, "#include \"cgo_helpers.h\"")
}

func TestCallbackProxy(t *testing.T) {
	spec := looperSpec()
	cb := spec.Callbacks["ALooper_callbackFunc"]

	proxy := generateCallbackProxy("looper", "ALooper_callbackFunc", cb, spec, nil)

	hash := callbackHashSuffix("looper", "ALooper_callbackFunc")
	upperHash := strings.ToUpper(hash)

	assert.Contains(t, proxy, "PassRef()")
	assert.Contains(t, proxy, "PassValue()")
	assert.Contains(t, proxy, "//export ALooper_callbackFunc"+upperHash)
	assert.Contains(t, proxy, "func ALooper_callbackFunc"+upperHash+"(")
	assert.Contains(t, proxy, "aLooper_callbackFunc"+upperHash+"Func")

	// Verify it has parameter conversion.
	lines := strings.Split(proxy, "\n")
	hasExport := false
	for _, l := range lines {
		if strings.Contains(l, "//export") {
			hasExport = true
		}
	}
	assert.True(t, hasExport)
}

func TestParamConversionOpaque(t *testing.T) {
	spec := looperSpec()
	callbackSet := map[string]bool{"ALooper_callbackFunc": true}

	p := specmodel.Param{Name: "looper", Type: "*ALooper"}
	conv := paramConversion(p, spec, callbackSet, nil)

	assert.Contains(t, conv.code, "(*C.ALooper)(unsafe.Pointer(looper))")
	assert.Equal(t, "clooper", conv.cVarName)
}

func TestParamConversionScalar(t *testing.T) {
	spec := looperSpec()
	callbackSet := map[string]bool{}

	p := specmodel.Param{Name: "fd", Type: "int32"}
	conv := paramConversion(p, spec, callbackSet, nil)

	assert.Contains(t, conv.code, "(C.int)(fd)")
	assert.Equal(t, "cfd", conv.cVarName)
}

func TestParamConversionCallback(t *testing.T) {
	spec := looperSpec()
	callbackSet := map[string]bool{"ALooper_callbackFunc": true}

	p := specmodel.Param{Name: "callback", Type: "ALooper_callbackFunc"}
	conv := paramConversion(p, spec, callbackSet, nil)

	assert.Contains(t, conv.code, ".PassValue()")
	assert.Equal(t, "ccallback", conv.cVarName)
}

func TestParamConversionUnsafePointer(t *testing.T) {
	spec := looperSpec()
	callbackSet := map[string]bool{}

	p := specmodel.Param{Name: "data", Type: "unsafe.Pointer"}
	conv := paramConversion(p, spec, callbackSet, nil)

	assert.Contains(t, conv.code, "data, cgoAllocsUnknown")
	assert.Equal(t, "cdata", conv.cVarName)
}

func TestReturnConversionPointer(t *testing.T) {
	spec := looperSpec()
	result := returnConversion("*ALooper", spec)
	assert.Contains(t, result, "unsafe.Pointer(&__ret)")
	assert.Contains(t, result, "return __v")
}

func TestReturnConversionScalar(t *testing.T) {
	spec := looperSpec()
	result := returnConversion("int32", spec)
	assert.Contains(t, result, "(int32)(__ret)")
}
