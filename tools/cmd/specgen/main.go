// Command specgen extracts structured YAML specs from c2ffi JSON output
// (or legacy Go AST parsing for test fixtures).
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xaionaro-go/ndk/tools/pkg/c2ffi"
	"github.com/xaionaro-go/ndk/tools/pkg/specgen"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		// c2ffi mode flags.
		c2ffiJSON = flag.String("c2ffi", "", "path to c2ffi JSON output file (c2ffi mode)")
		manifest  = flag.String("manifest", "", "path to capi manifest YAML (required for c2ffi mode)")
		ndkIncDir = flag.String("ndk-include", "", "path to NDK sysroot include directory (for c2ffi invocation and callback params)")
		ndkPath   = flag.String("ndk-path", "", "path to NDK root directory (for finding NDK clang)")
		c2ffiBin  = flag.String("c2ffi-bin", "", "path to c2ffi binary (auto-invoke mode; omit -c2ffi to use)")

		// Legacy mode flags (c-for-go Go AST).
		pkgDir = flag.String("pkg", "", "path to capi package directory (legacy mode)")

		// Shared flags.
		module = flag.String("module", "", "module name (defaults to manifest PackageName or directory name)")
		srcPkg = flag.String("source-package", "", "Go import path of source package")
		out    = flag.String("out", "", "output YAML path")

		// Legacy-only.
		ndkHeaders = flag.String("ndk-headers", "", "path to NDK sysroot include directory for C header struct parsing (legacy mode)")
	)
	flag.Parse()

	if *out == "" {
		fmt.Fprintln(os.Stderr, "usage: specgen -out <file.yaml> [-manifest <yaml> -ndk-include <dir>] [-c2ffi <json>] [-pkg <dir>]")
		os.Exit(1)
	}

	switch {
	case *c2ffiJSON != "":
		// Direct c2ffi JSON file mode.
		runC2FFI(*c2ffiJSON, *manifest, *ndkIncDir, *module, *srcPkg, *out)
	case *manifest != "" && *ndkIncDir != "":
		// Auto-invoke c2ffi mode.
		runC2FFIAuto(*manifest, *ndkIncDir, *ndkPath, *c2ffiBin, *module, *srcPkg, *out)
	case *pkgDir != "":
		runLegacy(*pkgDir, *module, *srcPkg, *out, *ndkHeaders)
	default:
		fmt.Fprintln(os.Stderr, "error: specify -manifest + -ndk-include, -c2ffi, or -pkg")
		os.Exit(1)
	}
}

func runC2FFIAuto(
	manifestPath string,
	ndkIncDir string,
	ndkPath string,
	c2ffiBin string,
	module string,
	srcPkg string,
	outPath string,
) {
	mf := readManifest(manifestPath)

	if c2ffiBin == "" {
		c2ffiBin = "c2ffi"
	}

	// Derive NDK path from sysroot include dir if not provided.
	if ndkPath == "" {
		// .../sysroot/usr/include → .../ndk/VERSION
		ndkPath = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(ndkIncDir)))))
	}

	jsonData, err := c2ffi.Invoke(c2ffi.InvokeOptions{
		C2FFIBin:   c2ffiBin,
		NDKSysroot: ndkIncDir,
		NDKPath:    ndkPath,
		Includes:   mf.Generator.Includes,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "invoke c2ffi: %v\n", err)
		os.Exit(1)
	}

	// Filter by target headers to exclude system types from transitive includes.
	convertAndWrite(jsonData, &mf, ndkIncDir, module, srcPkg, outPath)
}

func runC2FFI(
	jsonPath string,
	manifestPath string,
	ndkIncDir string,
	module string,
	srcPkg string,
	outPath string,
) {
	mf := readManifest(manifestPath)

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read c2ffi JSON: %v\n", err)
		os.Exit(1)
	}

	convertAndWrite(jsonData, &mf, ndkIncDir, module, srcPkg, outPath)
}

func readManifest(path string) c2ffi.Manifest {
	var mf c2ffi.Manifest
	if path == "" {
		return mf
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read manifest: %v\n", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(data, &mf); err != nil {
		fmt.Fprintf(os.Stderr, "parse manifest: %v\n", err)
		os.Exit(1)
	}
	return mf
}

func convertAndWrite(
	jsonData []byte,
	mf *c2ffi.Manifest,
	ndkIncDir string,
	module string,
	srcPkg string,
	outPath string,
) {
	if module == "" {
		module = mf.Generator.PackageName
	}
	if module == "" {
		module = strings.TrimSuffix(filepath.Base(outPath), ".yaml")
	}
	if srcPkg == "" {
		srcPkg = "github.com/xaionaro-go/ndk/capi/" + module
	}

	var headerDirs []string
	if ndkIncDir != "" {
		for _, inc := range mf.Generator.Includes {
			dir := filepath.Join(ndkIncDir, filepath.Dir(inc))
			headerDirs = appendUnique(headerDirs, dir)
		}
	}

	opts := c2ffi.ConvertOptions{
		Module:        module,
		SourcePackage: srcPkg,
		TargetHeaders: mf.Generator.Includes,
		Rules:         mf.Translator.Rules.Global,
		NDKHeaderDirs: headerDirs,
	}

	spec, err := c2ffi.Convert(jsonData, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "convert c2ffi: %v\n", err)
		os.Exit(1)
	}

	if err := specgen.WriteSpec(*spec, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("specgen: wrote %s (%d types, %d enums, %d functions, %d callbacks, %d structs)\n",
		outPath, len(spec.Types), len(spec.Enums), len(spec.Functions), len(spec.Callbacks), len(spec.Structs))
}

func runLegacy(
	pkgDir string,
	module string,
	srcPkg string,
	outPath string,
	ndkHeaders string,
) {
	if module == "" {
		module = filepath.Base(pkgDir)
	}
	if srcPkg == "" {
		srcPkg = "github.com/xaionaro-go/ndk/capi/" + module
	}

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir: %v\n", err)
		os.Exit(1)
	}
	var sources []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".go" {
			sources = append(sources, filepath.Join(pkgDir, e.Name()))
		}
	}

	spec, err := specgen.ParseSources(module, srcPkg, sources)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		os.Exit(1)
	}

	if ndkHeaders != "" {
		headerDirs, err := specgen.ExtractIncludeDirsFromGoFiles(sources, ndkHeaders)
		if err != nil {
			fmt.Fprintf(os.Stderr, "extract include dirs: %v\n", err)
			os.Exit(1)
		}

		for _, dir := range headerDirs {
			structs, err := specgen.ParseStructsFromDir(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse structs from %s: %v\n", dir, err)
				os.Exit(1)
			}
			specgen.MergeStructs(&spec, structs)

			funcs, err := specgen.ParseFunctionsFromDir(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse functions from %s: %v\n", dir, err)
				os.Exit(1)
			}
			specgen.MergeFunctions(&spec, funcs)
		}
	}

	if err := specgen.WriteSpec(spec, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("specgen: wrote %s (%d types, %d enums, %d functions, %d callbacks)\n",
		outPath, len(spec.Types), len(spec.Enums), len(spec.Functions), len(spec.Callbacks))
}

func appendUnique(slice []string, s string) []string {
	for _, existing := range slice {
		if existing == s {
			return slice
		}
	}
	return append(slice, s)
}
