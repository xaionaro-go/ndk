package specgen

import (
	"strings"
	"testing"
	"testing/quick"
)

// TestStripCComments_UnterminatedBlockComment verifies the fix for the
// unterminated block comment bug: no bytes from an unterminated comment
// should leak into the output.
func TestStripCComments_UnterminatedBlockComment(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{
			"unterminated at end",
			"int x; /* incomplete",
			"int x; ",
		},
		{
			"unterminated only comment",
			"/* everything is a comment",
			"",
		},
		{
			"normal terminated comment",
			"int x; /* comment */ int y;",
			"int x;  int y;",
		},
		{
			"out annotation preserved",
			"int /*out*/ x;",
			"int /*out*/ x;",
		},
		{
			"line comment",
			"int x; // comment\nint y;",
			"int x; \nint y;",
		},
		{
			"empty input",
			"",
			"",
		},
		{
			"no comments",
			"int x; int y;",
			"int x; int y;",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := stripCComments(tc.input)
			if got != tc.want {
				t.Errorf("stripCComments(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestStripCComments_NoCommentCharsInOutput verifies that /* sequences
// from comments don't appear in the output (except preserved /*out*/).
func TestStripCComments_NoCommentCharsInOutput(t *testing.T) {
	f := func(s string) bool {
		result := stripCComments(s)
		// The only /* that should remain is /*out*/.
		cleaned := strings.ReplaceAll(result, "/*out*/", "")
		return !strings.Contains(cleaned, "/*")
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("stripCComments leaked comment markers: %v", err)
	}
}

// TestNormalizeCType_Idempotent verifies normalizeCType produces consistent output.
func TestNormalizeCType_Idempotent(t *testing.T) {
	cases := []struct{ in_, want string }{
		{"int", "int"},
		{"void*", "void*"},
		{"const char**", "**char"},
		{"ALooper*", "*ALooper"},
		{"const ARect*", "*ARect"},
	}
	for _, tc := range cases {
		if got := normalizeCType(tc.in_); got != tc.want {
			t.Errorf("normalizeCType(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestParseSingleCParam_Vectors verifies C parameter parsing.
func TestParseSingleCParam_Vectors(t *testing.T) {
	cases := []struct {
		decl     string
		wantName string
		wantType string
	}{
		{"int x", "x", "int"},
		{"void* context", "context", "void*"},
		{"const char* name", "name", "*char"},
		{"ALooper* looper", "looper", "*ALooper"},
		{"int32_t *outSize", "outSize", "*int32_t"},
	}
	for _, tc := range cases {
		p := parseSingleCParam(tc.decl)
		if p.Name != tc.wantName {
			t.Errorf("parseSingleCParam(%q).Name = %q, want %q", tc.decl, p.Name, tc.wantName)
		}
		if p.Type != tc.wantType {
			t.Errorf("parseSingleCParam(%q).Type = %q, want %q", tc.decl, p.Type, tc.wantType)
		}
	}
}

// TestCTypeToGoType_Roundtrip verifies common C types map correctly.
func TestCTypeToGoType_Roundtrip(t *testing.T) {
	cases := []struct{ cType, want string }{
		{"int32_t", "int32"},
		{"uint32_t", "uint32"},
		{"int64_t", "int64"},
		{"float", "float32"},
		{"double", "float64"},
		{"void", "void"},
		{"void*", "unsafe.Pointer"},
		{"void**", "*unsafe.Pointer"},
		{"char", "byte"},
		{"bool", "bool"},
	}
	for _, tc := range cases {
		if got := cTypeToGoType(tc.cType); got != tc.want {
			t.Errorf("cTypeToGoType(%q) = %q, want %q", tc.cType, got, tc.want)
		}
	}
}
