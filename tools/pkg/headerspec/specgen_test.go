package headerspec

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runClangForTest(t *testing.T, includes []string) []byte {
	t.Helper()

	ndkBase := os.ExpandEnv("${HOME}/Android/Sdk/ndk/28.0.13004108/toolchains/llvm/prebuilt/linux-x86_64")
	sysroot := filepath.Join(ndkBase, "sysroot")
	clang := filepath.Join(ndkBase, "bin/clang")

	if _, err := os.Stat(clang); err != nil {
		t.Skipf("clang not found: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "test-*.c")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	for _, inc := range includes {
		_, err := tmpFile.WriteString("#include <" + inc + ">\n")
		require.NoError(t, err)
	}
	_ = tmpFile.Close()

	args := []string{
		"-target", "aarch64-linux-android26",
		"--sysroot=" + sysroot,
		"-I" + filepath.Join(sysroot, "usr", "include"),
		"-Xclang", "-ast-dump=json",
		"-fsyntax-only",
		"-x", "c",
		tmpFile.Name(),
	}

	cmd := exec.Command(clang, args...)
	out, err := cmd.Output()
	require.NoError(t, err)

	return out
}

func TestGenerateSpecLooper(t *testing.T) {
	jsonData := runClangForTest(t, []string{"android/looper.h"})

	root, err := ParseClangAST(jsonData)
	require.NoError(t, err)

	decls := ExtractDeclarations(root, []string{"android/looper.h"})
	rules := []Rule{
		{Action: "ignore", From: "^ALooper_pollAll$"},
		{Action: "accept", From: "^ALooper"},
		{Action: "accept", From: "^ALOOPER_"},
	}
	filtered := ApplyRules(decls, rules)
	spec := GenerateSpec("looper", "github.com/AndroidGoLab/ndk/capi/looper", filtered)

	// Types.
	require.Contains(t, spec.Types, "ALooper")
	assert.Equal(t, "opaque_ptr", spec.Types["ALooper"].Kind)
	assert.Equal(t, "*C.ALooper", spec.Types["ALooper"].GoType)

	// Functions: ALooper_pollAll should be ignored, everything else accepted.
	require.Contains(t, spec.Functions, "ALooper_pollOnce")
	require.NotContains(t, spec.Functions, "ALooper_pollAll")
	require.Contains(t, spec.Functions, "ALooper_addFd")
	require.Contains(t, spec.Functions, "ALooper_acquire")
	require.Contains(t, spec.Functions, "ALooper_release")
	require.Contains(t, spec.Functions, "ALooper_wake")
	require.Contains(t, spec.Functions, "ALooper_forThread")
	require.Contains(t, spec.Functions, "ALooper_prepare")
	require.Contains(t, spec.Functions, "ALooper_removeFd")

	// ALooper_pollOnce parameter types.
	pollOnce := spec.Functions["ALooper_pollOnce"]
	require.Len(t, pollOnce.Params, 4)
	assert.Equal(t, "int32", pollOnce.Params[0].Type)
	assert.Equal(t, "[]int32", pollOnce.Params[1].Type)
	assert.Equal(t, "[]int32", pollOnce.Params[2].Type)
	assert.Equal(t, "[]unsafe.Pointer", pollOnce.Params[3].Type)

	// ALooper_addFd parameters.
	addFd := spec.Functions["ALooper_addFd"]
	require.Len(t, addFd.Params, 6)
	assert.Equal(t, "*ALooper", addFd.Params[0].Type)
	assert.Equal(t, "int32", addFd.Params[1].Type)
	assert.Equal(t, "ALooper_callbackFunc", addFd.Params[4].Type)
	assert.Equal(t, "unsafe.Pointer", addFd.Params[5].Type)

	// Callbacks.
	require.Contains(t, spec.Callbacks, "ALooper_callbackFunc")
	cb := spec.Callbacks["ALooper_callbackFunc"]
	assert.Equal(t, "int32", cb.Returns)
	require.Len(t, cb.Params, 3)
	assert.Equal(t, "int32", cb.Params[0].Type)
	assert.Equal(t, "int32", cb.Params[1].Type)
	assert.Equal(t, "unsafe.Pointer", cb.Params[2].Type)

	// Enums: all three groups should be present.
	require.NotNil(t, spec.Enums)
	assert.GreaterOrEqual(t, len(spec.Enums), 3)

	// Check that ALOOPER_ enum constants are present in some group.
	foundConstants := map[string]bool{}
	for _, values := range spec.Enums {
		for _, v := range values {
			foundConstants[v.Name] = true
		}
	}
	assert.True(t, foundConstants["ALOOPER_EVENT_INPUT"])
	assert.True(t, foundConstants["ALOOPER_POLL_WAKE"])
	assert.True(t, foundConstants["ALOOPER_PREPARE_ALLOW_NON_CALLBACKS"])
	assert.Equal(t, int64(-1), findEnumValue(spec.Enums, "ALOOPER_POLL_WAKE"))
	assert.Equal(t, int64(1), findEnumValue(spec.Enums, "ALOOPER_EVENT_INPUT"))
}

func findEnumValue(enums map[string][]specmodel.EnumValue, name string) int64 {
	for _, values := range enums {
		for _, v := range values {
			if v.Name == name {
				return v.Value
			}
		}
	}
	return 0
}

func TestApplyRulesIgnoreOverridesAccept(t *testing.T) {
	decls := &Declarations{
		Functions: []FuncDecl{
			{Name: "ALooper_pollAll"},
			{Name: "ALooper_pollOnce"},
			{Name: "ALooper_acquire"},
			{Name: "Unrelated_func"},
		},
	}
	rules := []Rule{
		{Action: "ignore", From: "^ALooper_pollAll$"},
		{Action: "accept", From: "^ALooper"},
	}

	filtered := ApplyRules(decls, rules)

	funcNames := map[string]bool{}
	for _, fn := range filtered.Functions {
		funcNames[fn.Name] = true
	}

	assert.True(t, funcNames["ALooper_pollOnce"])
	assert.True(t, funcNames["ALooper_acquire"])
	assert.False(t, funcNames["ALooper_pollAll"], "ALooper_pollAll should be ignored")
	assert.False(t, funcNames["Unrelated_func"], "unmatched functions should be rejected")
}

func TestCTypeToGoType(t *testing.T) {
	tests := []struct {
		cType  string
		goType string
	}{
		{"int", "int32"},
		{"unsigned int", "uint32"},
		{"int32_t", "int32"},
		{"uint32_t", "uint32"},
		{"int64_t", "int64"},
		{"uint64_t", "uint64"},
		{"int8_t", "int8"},
		{"uint8_t", "uint8"},
		{"int16_t", "int16"},
		{"uint16_t", "uint16"},
		{"float", "float32"},
		{"double", "float64"},
		{"size_t", "uint64"},
		{"ssize_t", "int64"},
		{"_Bool", "bool"},
		{"bool", "bool"},
		{"void", ""},
		{"void *", "unsafe.Pointer"},
		{"const void *", "unsafe.Pointer"},
		{"const char *", "string"},
		{"char *", "string"},
		{"ALooper *", "*ALooper"},
		{"struct ALooper", "ALooper"},
	}

	for _, tt := range tests {
		t.Run(tt.cType, func(t *testing.T) {
			got := cTypeToGoType(tt.cType)
			assert.Equal(t, tt.goType, got)
		})
	}
}

func TestDeriveEnumGroupName(t *testing.T) {
	tests := []struct {
		name      string
		constants []EnumConstant
		want      string
	}{
		{
			name: "multiple with shared prefix",
			constants: []EnumConstant{
				{Name: "ALOOPER_POLL_WAKE"},
				{Name: "ALOOPER_POLL_CALLBACK"},
				{Name: "ALOOPER_POLL_TIMEOUT"},
			},
			want: "ALOOPER_POLL",
		},
		{
			name: "single constant",
			constants: []EnumConstant{
				{Name: "ALOOPER_PREPARE_ALLOW_NON_CALLBACKS"},
			},
			want: "ALOOPER_PREPARE_ALLOW",
		},
		{
			name:      "empty",
			constants: nil,
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveEnumGroupName(tt.constants)
			assert.Equal(t, tt.want, got)
		})
	}
}
