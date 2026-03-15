package idiomgen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AndroidGoLab/ndk/tools/pkg/idiomgen"
	"github.com/AndroidGoLab/ndk/tools/pkg/overlaymodel"
	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
)

func TestGenerate(t *testing.T) {
	spec, overlay := buildFixture()

	tmplDir := filepath.Join("..", "..", "..", "templates")
	// Verify templates directory exists before running.
	if _, err := os.Stat(tmplDir); os.IsNotExist(err) {
		t.Fatalf("templates dir not found at %s", tmplDir)
	}

	outDir := t.TempDir()

	if err := idiomgen.Generate(spec, overlay, tmplDir, outDir); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Verify shared output files exist.
	sharedFiles := []string{
		"package.go",
		"errors.go",
		"functions.go",
	}
	for _, name := range sharedFiles {
		path := filepath.Join(outDir, name)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected file %s: %v", name, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("file %s is empty", name)
		}
	}

	// Verify per-type file exists for StreamBuilder.
	sbPath := filepath.Join(outDir, "stream_builder.go")
	if _, err := os.Stat(sbPath); err != nil {
		t.Errorf("expected per-type file stream_builder.go: %v", err)
	}

	// Verify per-type file contains struct + method.
	sbContent, err := os.ReadFile(sbPath)
	if err != nil {
		t.Fatalf("read stream_builder.go: %v", err)
	}
	sbStr := string(sbContent)
	if !strings.Contains(sbStr, "type StreamBuilder struct") {
		t.Error("stream_builder.go missing 'type StreamBuilder struct'")
	}
	if !strings.Contains(sbStr, "SetDeviceID") {
		t.Error("stream_builder.go missing 'SetDeviceID' method")
	}

	// Verify per-enum file contains expected enum type.
	dirPath := filepath.Join(outDir, "direction.go")
	dirContent, err := os.ReadFile(dirPath)
	if err != nil {
		t.Fatalf("read direction.go: %v", err)
	}
	dirStr := string(dirContent)
	if !strings.Contains(dirStr, "type Direction int32") {
		t.Error("direction.go missing 'type Direction int32'")
	}
	if !strings.Contains(dirStr, "package audio") {
		t.Error("direction.go missing 'package audio'")
	}

	// Verify errors.go contains the error type.
	errorsContent, err := os.ReadFile(filepath.Join(outDir, "errors.go"))
	if err != nil {
		t.Fatalf("read errors.go: %v", err)
	}
	errorsStr := string(errorsContent)
	if !strings.Contains(errorsStr, "type Error int32") {
		t.Error("errors.go missing 'type Error int32'")
	}
	if !strings.Contains(errorsStr, "func result[") {
		t.Error("errors.go missing result() helper")
	}
}

func TestGenerate_APILevelSplitting(t *testing.T) {
	spec, overlay := buildFixture()

	// Add a higher-API method to test splitting.
	spec.Functions["AAudioStreamBuilder_setNewFeature"] = specmodel.FuncDef{
		CName: "AAudioStreamBuilder_setNewFeature",
		Params: []specmodel.Param{
			{Name: "builder", Type: "*AAudioStreamBuilder"},
			{Name: "value", Type: "int32"},
		},
	}
	overlay.Functions["AAudioStreamBuilder_setNewFeature"] = overlaymodel.FuncOverlay{
		Receiver: "StreamBuilder",
		GoName:   "SetNewFeature",
		Chain:    true,
	}
	overlay.APILevels["AAudioStreamBuilder_setNewFeature"] = 36

	tmplDir := filepath.Join("..", "..", "..", "templates")
	if _, err := os.Stat(tmplDir); os.IsNotExist(err) {
		t.Fatalf("templates dir not found at %s", tmplDir)
	}

	outDir := t.TempDir()
	if err := idiomgen.Generate(spec, overlay, tmplDir, outDir); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Base per-type file should have base methods but NOT the API-36 method.
	sbPath := filepath.Join(outDir, "stream_builder.go")
	sbContent, err := os.ReadFile(sbPath)
	if err != nil {
		t.Fatalf("read stream_builder.go: %v", err)
	}
	sbStr := string(sbContent)
	if !strings.Contains(sbStr, "SetDeviceID") {
		t.Error("stream_builder.go missing base method SetDeviceID")
	}
	if strings.Contains(sbStr, "SetNewFeature") {
		t.Error("stream_builder.go must NOT contain API-36 method SetNewFeature")
	}

	// API-36 file should exist with build tag and the higher-API method.
	api36Path := filepath.Join(outDir, "stream_builder_api36.go")
	api36Content, err := os.ReadFile(api36Path)
	if err != nil {
		t.Fatalf("expected stream_builder_api36.go: %v", err)
	}
	api36Str := string(api36Content)
	if !strings.Contains(api36Str, "//go:build android_ndk36") {
		t.Error("stream_builder_api36.go missing build tag")
	}
	if !strings.Contains(api36Str, "SetNewFeature") {
		t.Error("stream_builder_api36.go missing API-36 method SetNewFeature")
	}
}

func TestGenerate_MissingTemplate_Skipped(t *testing.T) {
	spec, overlay := buildFixture()

	// Use an empty template directory; no templates means no output files,
	// but Generate should not error.
	tmplDir := t.TempDir()
	outDir := t.TempDir()

	if err := idiomgen.Generate(spec, overlay, tmplDir, outDir); err != nil {
		t.Fatalf("Generate with empty templates: %v", err)
	}

	// Output dir should exist but contain no generated files.
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("read outDir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty outDir, got %d entries", len(entries))
	}
}

func TestLoadSpec(t *testing.T) {
	yaml := `module: aaudio
source_package: github.com/AndroidGoLab/ndk/capi/aaudio
types:
  AAudioStream:
    kind: opaque_ptr
    c_type: AAudioStream
    go_type: "*C.AAudioStream"
enums:
  Aaudio_direction_t:
    - name: AAUDIO_DIRECTION_OUTPUT
      value: 0
    - name: AAUDIO_DIRECTION_INPUT
      value: 1
functions:
  AAudioStream_close:
    c_name: AAudioStream_close
    params:
      - name: stream
        type: "*AAudioStream"
    returns: Aaudio_result_t
`
	path := filepath.Join(t.TempDir(), "spec.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	spec, err := idiomgen.LoadSpec(path)
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	if spec.Module != "aaudio" {
		t.Errorf("Module = %q, want %q", spec.Module, "aaudio")
	}
	if spec.SourcePackage != "github.com/AndroidGoLab/ndk/capi/aaudio" {
		t.Errorf("SourcePackage = %q", spec.SourcePackage)
	}
	if len(spec.Types) != 1 {
		t.Errorf("Types count = %d, want 1", len(spec.Types))
	}
	if spec.Types["AAudioStream"].Kind != "opaque_ptr" {
		t.Errorf("Types[AAudioStream].Kind = %q", spec.Types["AAudioStream"].Kind)
	}
	if len(spec.Enums) != 1 {
		t.Errorf("Enums count = %d, want 1", len(spec.Enums))
	}
	if len(spec.Functions) != 1 {
		t.Errorf("Functions count = %d, want 1", len(spec.Functions))
	}
}

func TestLoadSpec_FileNotFound(t *testing.T) {
	_, err := idiomgen.LoadSpec("/nonexistent/spec.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadOverlay(t *testing.T) {
	yaml := `module: aaudio
package:
  go_name: audio
  go_import: github.com/AndroidGoLab/ndk/audio
  doc: "Package audio provides Go bindings for Android AAudio."
types:
  AAudioStream:
    go_name: Stream
    pattern: ref_counted
functions:
  AAudioStream_close:
    receiver: Stream
    go_name: Close
api_levels:
  AAudioStream_close: 26
`
	path := filepath.Join(t.TempDir(), "overlay.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	ov, err := idiomgen.LoadOverlay(path)
	if err != nil {
		t.Fatalf("LoadOverlay: %v", err)
	}
	if ov.Module != "aaudio" {
		t.Errorf("Module = %q, want %q", ov.Module, "aaudio")
	}
	if ov.Package.GoName != "audio" {
		t.Errorf("Package.GoName = %q", ov.Package.GoName)
	}
	if ov.Package.GoImport != "github.com/AndroidGoLab/ndk/audio" {
		t.Errorf("Package.GoImport = %q", ov.Package.GoImport)
	}
	if len(ov.Types) != 1 {
		t.Errorf("Types count = %d, want 1", len(ov.Types))
	}
	if ov.Types["AAudioStream"].GoName != "Stream" {
		t.Errorf("Types[AAudioStream].GoName = %q", ov.Types["AAudioStream"].GoName)
	}
	if len(ov.Functions) != 1 {
		t.Errorf("Functions count = %d, want 1", len(ov.Functions))
	}
	if ov.APILevels["AAudioStream_close"] != 26 {
		t.Errorf("APILevels[AAudioStream_close] = %d, want 26", ov.APILevels["AAudioStream_close"])
	}
}

func TestLoadOverlay_FileNotFound(t *testing.T) {
	_, err := idiomgen.LoadOverlay("/nonexistent/overlay.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadOverlayOrEmpty_EmptyPath(t *testing.T) {
	ov, err := idiomgen.LoadOverlayOrEmpty("")
	if err != nil {
		t.Fatalf("LoadOverlayOrEmpty: %v", err)
	}
	if ov.Module != "" {
		t.Errorf("Module = %q, want empty", ov.Module)
	}
}

func TestLoadOverlayOrEmpty_WithPath(t *testing.T) {
	yaml := `module: test
package:
  go_name: testpkg
  go_import: github.com/AndroidGoLab/ndk/testpkg
`
	path := filepath.Join(t.TempDir(), "overlay.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	ov, err := idiomgen.LoadOverlayOrEmpty(path)
	if err != nil {
		t.Fatalf("LoadOverlayOrEmpty: %v", err)
	}
	if ov.Module != "test" {
		t.Errorf("Module = %q, want %q", ov.Module, "test")
	}
}

func TestGenerate_CreatesOutputDir(t *testing.T) {
	spec := specmodel.Spec{
		Module:        "test",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/test",
	}
	overlay := overlaymodel.Overlay{
		Module: "test",
		Package: overlaymodel.PackageOverlay{
			GoName:   "testpkg",
			GoImport: "github.com/AndroidGoLab/ndk/testpkg",
		},
	}

	// Use a nested output dir that does not yet exist.
	outDir := filepath.Join(t.TempDir(), "sub", "dir", "out")
	tmplDir := t.TempDir() // empty templates dir

	if err := idiomgen.Generate(spec, overlay, tmplDir, outDir); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	info, err := os.Stat(outDir)
	if err != nil {
		t.Fatalf("outDir does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("outDir is not a directory")
	}
}
