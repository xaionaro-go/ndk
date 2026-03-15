// Command headerspec generates spec YAML and capi/ Go packages from NDK C headers
// using clang AST.
//
// Usage:
//
//	headerspec -manifest capi/manifests/looper.yaml \
//	    -ndk-sysroot $NDK_SYSROOT \
//	    -module-path github.com/AndroidGoLab/ndk \
//	    -out-spec spec/generated/looper.yaml \
//	    -out-capi capi/looper
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AndroidGoLab/ndk/tools/pkg/headerspec"
)

func main() {
	var (
		manifestPath = flag.String("manifest", "", "path to manifest YAML (e.g. capi/manifests/looper.yaml)")
		ndkSysroot   = flag.String("ndk-sysroot", "", "path to NDK sysroot (e.g. $NDK/toolchains/llvm/prebuilt/linux-x86_64/sysroot)")
		modulePath   = flag.String("module-path", "github.com/AndroidGoLab/ndk", "Go module path")
		outSpec      = flag.String("out-spec", "", "output spec YAML path")
		outCapi      = flag.String("out-capi", "", "output capi/ package directory (e.g. capi/looper)")
		target       = flag.String("target", "aarch64-linux-android26", "clang target triple")
	)
	flag.Parse()

	if *manifestPath == "" || *ndkSysroot == "" || (*outSpec == "" && *outCapi == "") {
		fmt.Fprintln(os.Stderr, "usage: headerspec -manifest <file> -ndk-sysroot <dir> [-out-spec <file>] [-out-capi <dir>] [-module-path <path>] [-target <triple>]")
		os.Exit(1)
	}

	if err := run(*manifestPath, *ndkSysroot, *modulePath, *outSpec, *outCapi, *target); err != nil {
		fmt.Fprintf(os.Stderr, "headerspec: %v\n", err)
		os.Exit(1)
	}
}

func run(
	manifestPath string,
	ndkSysroot string,
	modulePath string,
	outSpec string,
	outCapi string,
	target string,
) error {
	manifest, err := headerspec.ParseManifest(manifestPath)
	if err != nil {
		return err
	}

	clangBin := findClang(ndkSysroot)

	jsonData, err := runClang(clangBin, ndkSysroot, target, manifest.Generator.Includes)
	if err != nil {
		return err
	}

	root, err := headerspec.ParseClangAST(jsonData)
	if err != nil {
		return fmt.Errorf("parsing clang AST: %w", err)
	}

	decls := headerspec.ExtractDeclarations(root, manifest.Generator.Includes)

	rules := manifest.Translator.Rules["global"]
	filtered := headerspec.ApplyRules(decls, rules)

	if outSpec != "" {
		sourcePackage := modulePath + "/capi/" + manifest.Generator.PackageName
		spec := headerspec.GenerateSpec(manifest.Generator.PackageName, sourcePackage, filtered)

		if err := os.MkdirAll(filepath.Dir(outSpec), 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}

		if err := headerspec.WriteSpecYAML(spec, outSpec); err != nil {
			return err
		}

		nTypes := 0
		if spec.Types != nil {
			nTypes = len(spec.Types)
		}
		nEnums := 0
		if spec.Enums != nil {
			nEnums = len(spec.Enums)
		}
		nFuncs := 0
		if spec.Functions != nil {
			nFuncs = len(spec.Functions)
		}
		nCallbacks := 0
		if spec.Callbacks != nil {
			nCallbacks = len(spec.Callbacks)
		}

		fmt.Printf("headerspec: wrote %s (%d types, %d enums, %d functions, %d callbacks)\n",
			outSpec, nTypes, nEnums, nFuncs, nCallbacks)
	}

	if outCapi != "" {
		if err := headerspec.GenerateCapiPackage(manifest, filtered, outCapi); err != nil {
			return fmt.Errorf("generating capi package: %w", err)
		}
		fmt.Printf("headerspec: wrote capi package to %s (%d functions, %d types)\n",
			outCapi, len(filtered.Functions), len(filtered.Typedefs))
	}

	return nil
}

// findClang locates the clang binary relative to the NDK sysroot.
// The sysroot is at .../toolchains/llvm/prebuilt/<host>/sysroot,
// so the binary is at .../toolchains/llvm/prebuilt/<host>/bin/clang.
func findClang(sysroot string) string {
	// Strip trailing /sysroot to get the prebuilt prefix.
	prefix := strings.TrimSuffix(sysroot, "/sysroot")
	prefix = strings.TrimSuffix(prefix, "/")
	return filepath.Join(prefix, "bin", "clang")
}

// runClang invokes clang with -ast-dump=json and returns the JSON output.
func runClang(
	clangBin string,
	sysroot string,
	target string,
	includes []string,
) ([]byte, error) {
	// Create a temporary .c file that includes all requested headers.
	tmpFile, err := os.CreateTemp("", "headerspec-*.c")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	var sb strings.Builder
	for _, inc := range includes {
		fmt.Fprintf(&sb, "#include <%s>\n", inc)
	}
	if _, err := tmpFile.WriteString(sb.String()); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("writing temp file: %w", err)
	}
	_ = tmpFile.Close()

	args := []string{
		"-target", target,
		"--sysroot=" + sysroot,
		"-I" + filepath.Join(sysroot, "usr", "include"),
		"-Xclang", "-ast-dump=json",
		"-fsyntax-only",
		"-x", "c",
		tmpFile.Name(),
	}

	cmd := exec.Command(clangBin, args...)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running clang: %w", err)
	}

	return out, nil
}
