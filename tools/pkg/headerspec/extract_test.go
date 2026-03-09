package headerspec

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestExtractDeclarationsLooper(t *testing.T) {
	ndkBase := os.ExpandEnv("${HOME}/Android/Sdk/ndk/28.0.13004108/toolchains/llvm/prebuilt/linux-x86_64")
	sysroot := filepath.Join(ndkBase, "sysroot")
	clang := filepath.Join(ndkBase, "bin/clang")

	if _, err := os.Stat(clang); err != nil {
		t.Skipf("clang not found at %s: %v", clang, err)
	}

	tmpFile, err := os.CreateTemp("", "test-*.c")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	if _, err := tmpFile.WriteString("#include <android/looper.h>\n"); err != nil {
		t.Fatalf("writing temp file: %v", err)
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
	if err != nil {
		t.Fatalf("clang: %v", err)
	}

	root, err := ParseClangAST(out)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	decls := ExtractDeclarations(root, []string{"android/looper.h"})

	data, _ := json.MarshalIndent(decls, "", "  ")
	t.Logf("Extracted declarations:\n%s", string(data))

	// Verify functions.
	funcNames := map[string]bool{}
	for _, fn := range decls.Functions {
		funcNames[fn.Name] = true
	}
	expectedFuncs := []string{
		"ALooper_forThread", "ALooper_prepare", "ALooper_acquire",
		"ALooper_release", "ALooper_pollOnce", "ALooper_pollAll",
		"ALooper_wake", "ALooper_addFd", "ALooper_removeFd",
	}
	for _, name := range expectedFuncs {
		if !funcNames[name] {
			t.Errorf("missing function: %s", name)
		}
	}

	// Verify typedefs.
	typedefNames := map[string]bool{}
	for _, td := range decls.Typedefs {
		typedefNames[td.Name] = true
	}
	if !typedefNames["ALooper"] {
		t.Error("missing typedef: ALooper")
	}
	if !typedefNames["ALooper_callbackFunc"] {
		t.Error("missing typedef: ALooper_callbackFunc")
	}

	// Verify enums.
	t.Logf("Found %d enums", len(decls.Enums))
	totalConstants := 0
	for _, e := range decls.Enums {
		totalConstants += len(e.Constants)
		t.Logf("Enum: name=%q typedefName=%q constants=%d", e.Name, e.TypedefName, len(e.Constants))
		for _, c := range e.Constants {
			t.Logf("  %s = %d", c.Name, c.Value)
		}
	}
	if totalConstants < 10 {
		t.Errorf("expected at least 10 enum constants, got %d", totalConstants)
	}

	// Verify the ALooper_callbackFunc is detected as a function pointer.
	for _, td := range decls.Typedefs {
		if td.Name == "ALooper_callbackFunc" {
			if !td.IsFuncPtr {
				t.Error("ALooper_callbackFunc should be detected as function pointer")
			}
			if td.FuncReturn != "int" {
				t.Errorf("ALooper_callbackFunc return type: got %q, want %q", td.FuncReturn, "int")
			}
			if len(td.FuncParams) != 3 {
				t.Errorf("ALooper_callbackFunc params: got %d, want 3", len(td.FuncParams))
			}
		}
	}
}
