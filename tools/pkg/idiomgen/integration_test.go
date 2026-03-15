package idiomgen_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/AndroidGoLab/ndk/tools/pkg/idiomgen"
	"github.com/AndroidGoLab/ndk/tools/pkg/overlaymodel"
	"github.com/AndroidGoLab/ndk/tools/pkg/specgen"
)

// projectRoot returns the absolute path to the project root.
func projectRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file path")
	}
	// This file is at tools/pkg/idiomgen/integration_test.go, so go up three times.
	return filepath.Join(filepath.Dir(filename), "..", "..", "..")
}

// TestIntegration_FullPipeline runs the complete specgen -> idiomgen pipeline:
//
//  1. Parse the test fixture (specgen stage 2)
//  2. Write the spec to a temp file
//  3. Load the spec back from the temp file
//  4. Run idiomgen.Generate with the loaded spec, empty overlay, and real templates
//  5. Verify the output directory has generated files with expected content
func TestIntegration_FullPipeline(t *testing.T) {
	root := projectRoot(t)

	// Stage 2 (specgen): Parse the test fixture.
	fixturePath := filepath.Join(root, "tools", "pkg", "specgen", "testdata", "simple", "simple.go")
	if _, err := os.Stat(fixturePath); err != nil {
		t.Fatalf("fixture not found: %v", err)
	}

	spec, err := specgen.ParseSources("simple", "github.com/AndroidGoLab/ndk/capi/simple", []string{fixturePath})
	if err != nil {
		t.Fatalf("ParseSources: %v", err)
	}

	// Verify the parsed spec is non-empty.
	if len(spec.Types) == 0 {
		t.Fatal("ParseSources returned empty Types")
	}
	if len(spec.Functions) == 0 {
		t.Fatal("ParseSources returned empty Functions")
	}

	// Write spec to temp file.
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "simple.yaml")
	if err := specgen.WriteSpec(spec, specPath); err != nil {
		t.Fatalf("WriteSpec: %v", err)
	}

	// Verify the spec file was written.
	info, err := os.Stat(specPath)
	if err != nil {
		t.Fatalf("spec file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("spec file is empty")
	}

	// Load spec back (idiomgen input).
	loadedSpec, err := idiomgen.LoadSpec(specPath)
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}

	// Verify round-trip preserved data.
	if loadedSpec.Module != "simple" {
		t.Errorf("round-trip Module = %q, want %q", loadedSpec.Module, "simple")
	}
	if loadedSpec.SourcePackage != "github.com/AndroidGoLab/ndk/capi/simple" {
		t.Errorf("round-trip SourcePackage = %q", loadedSpec.SourcePackage)
	}
	if len(loadedSpec.Types) != len(spec.Types) {
		t.Errorf("round-trip Types count = %d, want %d", len(loadedSpec.Types), len(spec.Types))
	}

	// Stage 3 (idiomgen): Generate with empty overlay and real templates.
	overlay := overlaymodel.Overlay{
		Module: "simple",
		Package: overlaymodel.PackageOverlay{
			GoName:   "simple",
			GoImport: "github.com/AndroidGoLab/ndk/simple",
			Doc:      "Package simple is a test package.",
		},
		Types: map[string]overlaymodel.TypeOverlay{
			"FakeStream": {
				GoName:     "Stream",
				Destructor: "FakeStream_close",
			},
			"FakeBuilder": {
				GoName:      "Builder",
				Constructor: "Fake_createBuilder",
				Destructor:  "FakeBuilder_delete",
				Pattern:     "builder",
			},
			"Fake_result_t": {
				GoError:      true,
				SuccessValue: "FAKE_OK",
				ErrorPrefix:  "fake",
			},
			"Fake_direction_t": {
				GoName:       "Direction",
				StripPrefix:  "FAKE_DIRECTION_",
				StringMethod: true,
			},
		},
		Functions: map[string]overlaymodel.FuncOverlay{
			"FakeBuilder_setDeviceId": {
				Receiver: "Builder",
				GoName:   "SetDeviceID",
				Chain:    true,
			},
			"FakeBuilder_openStream": {
				Receiver:   "Builder",
				GoName:     "Open",
				ReturnsNew: "Stream",
			},
			"FakeStream_close": {
				Receiver: "Stream",
				GoName:   "Close",
			},
			"FakeStream_requestStart": {
				Receiver: "Stream",
				GoName:   "Start",
			},
			"FakeStream_getSampleRate": {
				Receiver: "Stream",
				GoName:   "SampleRate",
				Pure:     true,
			},
		},
	}

	tmplDir := filepath.Join(root, "templates")
	if _, err := os.Stat(tmplDir); err != nil {
		t.Fatalf("templates dir not found: %v", err)
	}

	outDir := filepath.Join(tmpDir, "output")
	if err := idiomgen.Generate(loadedSpec, overlay, tmplDir, outDir); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Verify output directory has generated files.
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("read output dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("output directory is empty")
	}

	generatedFiles := make(map[string]bool)
	for _, e := range entries {
		generatedFiles[e.Name()] = true
	}

	sharedFiles := []string{"package.go", "errors.go", "functions.go"}
	for _, name := range sharedFiles {
		if !generatedFiles[name] {
			t.Errorf("missing generated file %q", name)
		}
	}

	// Per-type files should exist.
	expectedFiles := []string{"builder.go", "stream.go", "direction.go"}
	for _, name := range expectedFiles {
		if !generatedFiles[name] {
			t.Errorf("missing per-type file %q", name)
		}
	}

	// direction.go should contain the Direction enum.
	assertFileContains(t, filepath.Join(outDir, "direction.go"), "package simple", "package declaration")
	assertFileContains(t, filepath.Join(outDir, "direction.go"), "type Direction int32", "Direction enum type")
	assertFileContains(t, filepath.Join(outDir, "direction.go"), "Output", "stripped enum value")

	// errors.go should contain error type for Fake_result_t.
	assertFileContains(t, filepath.Join(outDir, "errors.go"), "type Error int32", "error type")

	// Per-type files should contain opaque wrapper types and methods.
	assertFileContains(t, filepath.Join(outDir, "builder.go"), "type Builder struct", "Builder opaque type")
	assertFileContains(t, filepath.Join(outDir, "builder.go"), "SetDeviceID", "builder setter method")
	assertFileContains(t, filepath.Join(outDir, "stream.go"), "type Stream struct", "Stream opaque type")
	assertFileContains(t, filepath.Join(outDir, "stream.go"), "SampleRate", "getter method")
}

// assertFileContains reads a file and checks that it contains the expected substring.
func assertFileContains(t *testing.T, path, expected, description string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("read %s: %v", filepath.Base(path), err)
		return
	}
	if !strings.Contains(string(data), expected) {
		t.Errorf("%s: %s does not contain %q", filepath.Base(path), description, expected)
	}
}
