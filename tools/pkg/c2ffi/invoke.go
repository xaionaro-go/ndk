package c2ffi

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// InvokeOptions configures how to run c2ffi on NDK headers.
type InvokeOptions struct {
	// C2FFIBin is the path to the c2ffi binary.
	C2FFIBin string

	// NDKSysroot is the NDK sysroot include directory
	// (e.g., .../sysroot/usr/include).
	NDKSysroot string

	// NDKPath is the root NDK directory (e.g., .../ndk/28.0.13004108).
	// Used to find NDK clang for preprocessing.
	NDKPath string

	// Includes are the header files to process
	// (e.g., ["android/looper.h", "camera/NdkCameraDevice.h"]).
	Includes []string

	// Target is the clang target triple (default: aarch64-linux-android35).
	Target string
}

// Invoke preprocesses each NDK header with NDK clang (to strip availability
// attributes that c2ffi can't handle), then runs c2ffi on each preprocessed
// header and merges the JSON arrays.
func Invoke(opts InvokeOptions) ([]byte, error) {
	if opts.Target == "" {
		opts.Target = "aarch64-linux-android35"
	}
	if opts.C2FFIBin == "" {
		var err error
		opts.C2FFIBin, err = exec.LookPath("c2ffi")
		if err != nil {
			return nil, fmt.Errorf("c2ffi not found in PATH: %w", err)
		}
	}

	ndkClang := findNDKClang(opts.NDKPath)

	var allDecls []Declaration
	for _, inc := range opts.Includes {
		headerPath := filepath.Join(opts.NDKSysroot, inc)

		// Preprocess with NDK clang to strip __attribute__((availability(...)))
		// annotations. c2ffi silently drops declarations with these attributes.
		ppPath, err := preprocessHeader(ndkClang, headerPath, opts)
		if err != nil {
			return nil, fmt.Errorf("preprocess %s: %w", inc, err)
		}
		defer func() { _ = os.Remove(ppPath) }()

		args := []string{ppPath}

		cmd := exec.Command(opts.C2FFIBin, args...)
		out, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return nil, fmt.Errorf("c2ffi failed on %s: %s\n%s", inc, err, string(exitErr.Stderr))
			}
			return nil, fmt.Errorf("c2ffi failed on %s: %w", inc, err)
		}

		var decls []Declaration
		if err := json.Unmarshal(out, &decls); err != nil {
			return nil, fmt.Errorf("parsing c2ffi JSON for %s: %w", inc, err)
		}
		allDecls = append(allDecls, decls...)
	}

	return json.Marshal(allDecls)
}

// preprocessHeader uses NDK clang to preprocess a header, stripping
// __attribute__ annotations that c2ffi can't handle.
func preprocessHeader(
	ndkClang string,
	headerPath string,
	opts InvokeOptions,
) (string, error) {
	sysroot := filepath.Dir(filepath.Dir(opts.NDKSysroot)) // .../sysroot

	tmpFile, err := os.CreateTemp("", "c2ffi-pp-*.h")
	if err != nil {
		return "", err
	}
	_ = tmpFile.Close()

	args := []string{
		"-E",
		"-target", opts.Target,
		"--sysroot", sysroot,
		"-D__attribute__(x)=",
		"-D_Nonnull=",
		"-D_Nullable=",
		"-D_Null_unspecified=",
		headerPath,
	}

	cmd := exec.Command(ndkClang, args...)
	out, err := cmd.Output()
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s", ndkClang, string(exitErr.Stderr))
		}
		return "", err
	}

	if err := os.WriteFile(tmpFile.Name(), out, 0o644); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// findNDKClang locates the NDK clang binary.
func findNDKClang(ndkPath string) string {
	if ndkPath != "" {
		clang := filepath.Join(ndkPath, "toolchains/llvm/prebuilt/linux-x86_64/bin/clang")
		if _, err := os.Stat(clang); err == nil {
			return clang
		}
	}

	// Fall back to PATH.
	if p, err := exec.LookPath("clang"); err == nil {
		return p
	}

	return "clang"
}
