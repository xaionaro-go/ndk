package examples_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestExamplesCompile cross-compiles all examples for Android arm64
// to catch compile errors. Requires the Android NDK to be installed.
func TestExamplesCompile(t *testing.T) {
	ndkHome := os.Getenv("ANDROID_NDK_HOME")
	if ndkHome == "" {
		androidHome := os.Getenv("ANDROID_HOME")
		if androidHome == "" {
			androidHome = filepath.Join(os.Getenv("HOME"), "Android", "Sdk")
		}
		matches, _ := filepath.Glob(filepath.Join(androidHome, "ndk", "*"))
		if len(matches) == 0 {
			t.Skip("Android NDK not found; set ANDROID_NDK_HOME or ANDROID_HOME")
		}
		ndkHome = matches[len(matches)-1]
	}

	cc := filepath.Join(
		ndkHome,
		"toolchains", "llvm", "prebuilt", runtime.GOOS+"-x86_64",
		"bin", "aarch64-linux-android35-clang",
	)
	if _, err := os.Stat(cc); err != nil {
		t.Skipf("NDK clang not found at %s", cc)
	}

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file path")
	}
	projectRoot := filepath.Join(filepath.Dir(thisFile), "..")

	cmd := exec.Command("go", "build", "./examples/...")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"GOOS=android",
		"GOARCH=arm64",
		"CC="+cc,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("examples failed to compile:\n%s", out)
	}
}
