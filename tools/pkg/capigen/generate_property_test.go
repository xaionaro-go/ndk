package capigen

import (
	"strings"
	"testing"
	"testing/quick"
	"unicode"
)

// TestExportName_Idempotent verifies that ExportName(ExportName(x)) == ExportName(x)
// for all inputs. This is the idempotency property proven in Lean 4.
func TestExportName_Idempotent(t *testing.T) {
	f := func(name string) bool {
		return ExportName(ExportName(name)) == ExportName(name)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("ExportName idempotency violated: %v", err)
	}
}

// TestExportName_NonEmpty verifies that ExportName never returns "" for non-empty input.
func TestExportName_NonEmpty(t *testing.T) {
	f := func(name string) bool {
		if name == "" {
			return true
		}
		return ExportName(name) != ""
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("ExportName returned empty for non-empty input: %v", err)
	}
}

// TestExportName_UppercaseStart verifies that the result starts with an uppercase
// letter (when the input has any non-underscore characters).
func TestExportName_UppercaseStart(t *testing.T) {
	f := func(name string) bool {
		result := ExportName(name)
		if result == "" {
			return true
		}
		// If all underscores, result is unchanged.
		stripped := strings.TrimLeft(name, "_")
		if stripped == "" {
			return result == name
		}
		first := rune(result[0])
		// First char should be uppercase if it was lowercase.
		if stripped[0] >= 'a' && stripped[0] <= 'z' {
			return unicode.IsUpper(first)
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("ExportName uppercase-start violated: %v", err)
	}
}

// TestExportName_AllUnderscores verifies that ExportName returns the input
// unchanged when it consists entirely of underscores.
func TestExportName_AllUnderscores(t *testing.T) {
	for n := 0; n <= 10; n++ {
		name := strings.Repeat("_", n)
		if got := ExportName(name); got != name {
			t.Errorf("ExportName(%q) = %q, want %q", name, got, name)
		}
	}
}

// TestStripPointers_Idempotent verifies that stripPointers(stripPointers(x)) == stripPointers(x).
func TestStripPointers_Idempotent(t *testing.T) {
	f := func(typ string) bool {
		return stripPointers(stripPointers(typ)) == stripPointers(typ)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("stripPointers idempotency violated: %v", err)
	}
}

// TestStripPointers_NoLeadingStar verifies that the result has no leading '*'.
func TestStripPointers_NoLeadingStar(t *testing.T) {
	f := func(typ string) bool {
		result := stripPointers(typ)
		return result == "" || result[0] != '*'
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("stripPointers has leading star: %v", err)
	}
}

// TestIsScalarGoType_Exhaustive verifies all expected scalar types are recognized.
func TestIsScalarGoType_Exhaustive(t *testing.T) {
	scalars := []string{
		"int8", "uint8", "int16", "uint16", "int32", "uint32",
		"int64", "uint64", "float32", "float64", "bool", "int", "uint",
	}
	for _, s := range scalars {
		if !isScalarGoType(s) {
			t.Errorf("isScalarGoType(%q) = false, want true", s)
		}
	}

	nonScalars := []string{
		"string", "unsafe.Pointer", "ALooper", "camera_status_t",
		"*int32", "[]byte", "",
	}
	for _, s := range nonScalars {
		if isScalarGoType(s) {
			t.Errorf("isScalarGoType(%q) = true, want false", s)
		}
	}
}

// TestSanitizeGoIdent_CoversAllKeywords verifies that all Go keywords are sanitized.
func TestSanitizeGoIdent_CoversAllKeywords(t *testing.T) {
	keywords := []string{
		"break", "case", "chan", "const", "continue",
		"default", "defer", "else", "fallthrough", "for",
		"func", "go", "goto", "if", "import",
		"interface", "map", "package", "range", "return",
		"select", "struct", "switch", "type", "var",
	}
	for _, kw := range keywords {
		result := sanitizeGoIdent(kw)
		if result == kw {
			t.Errorf("sanitizeGoIdent(%q) returned keyword unchanged", kw)
		}
		if result != "_"+kw {
			t.Errorf("sanitizeGoIdent(%q) = %q, want %q", kw, result, "_"+kw)
		}
	}
}

// TestSanitizeGoIdent_NonKeyword verifies non-keywords are unchanged.
func TestSanitizeGoIdent_NonKeyword(t *testing.T) {
	f := func(name string) bool {
		// Skip actual keywords.
		switch name {
		case "type", "func", "var", "const", "import", "package",
			"return", "range", "map", "chan", "interface", "struct",
			"error", "string", "bool", "int", "uint", "byte", "rune",
			"select", "case", "default", "switch", "break", "continue",
			"for", "go", "goto", "if", "else", "defer", "fallthrough":
			return true
		}
		return sanitizeGoIdent(name) == name
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("sanitizeGoIdent modified non-keyword: %v", err)
	}
}
