package idiomgen

import (
	"strings"
	"testing"
	"testing/quick"
)

// TestResolveType_EmptyMap verifies resolveType returns input when map is empty.
func TestResolveType_EmptyMap(t *testing.T) {
	f := func(specType string) bool {
		return resolveType(specType, nil) == specType
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("resolveType with nil map changed input: %v", err)
	}
}

// TestResolveType_DirectPrecedence verifies direct lookup takes precedence.
func TestResolveType_DirectPrecedence(t *testing.T) {
	typeMap := map[string]string{
		"*foo": "Direct",
		"foo":  "Prefix",
	}
	if got := resolveType("*foo", typeMap); got != "Direct" {
		t.Errorf("resolveType(*foo) = %q, want Direct", got)
	}
}

// TestResolveType_PrefixPreservation verifies prefix structure is preserved.
func TestResolveType_PrefixPreservation(t *testing.T) {
	typeMap := map[string]string{"foo": "Bar"}

	cases := []struct{ in_, want string }{
		{"foo", "Bar"},
		{"*foo", "*Bar"},
		{"**foo", "**Bar"},
		{"[]foo", "[]Bar"},
		{"unknown", "unknown"},
	}
	for _, tc := range cases {
		if got := resolveType(tc.in_, typeMap); got != tc.want {
			t.Errorf("resolveType(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestResolveType_Identity verifies unmapped types are returned unchanged.
func TestResolveType_Identity(t *testing.T) {
	typeMap := map[string]string{"foo": "Bar"}
	unmapped := []string{"baz", "*baz", "**baz", "[]baz", "int32", "unsafe.Pointer"}
	for _, typ := range unmapped {
		if got := resolveType(typ, typeMap); got != typ {
			t.Errorf("resolveType(%q) = %q, want %q", typ, got, typ)
		}
	}
}

// TestIsConstructorDestructor_Disjoint verifies the critical invariant that
// no function name can be both a constructor and destructor for the same type.
func TestIsConstructorDestructor_Disjoint(t *testing.T) {
	types := []string{
		"ASensor", "ALooper", "ACameraManager", "AImageReader",
		"AAudioStreamBuilder", "ANativeWindow", "AMediaCodec",
	}
	suffixes := append(
		[]string{"_create", "_new"},
		"_delete", "_free", "_destroy", "_release", "_close",
	)

	for _, typ := range types {
		for _, suffix := range suffixes {
			funcName := typ + suffix
			isCtor := isConstructorFunc(funcName, typ)
			isDtor := isDestructorFunc(funcName, typ)
			if isCtor && isDtor {
				t.Errorf("DISJOINTNESS VIOLATION: %q is both constructor and destructor for %q",
					funcName, typ)
			}
		}
	}
}

// TestIsConstructorDestructor_Detection verifies correct detection.
func TestIsConstructorDestructor_Detection(t *testing.T) {
	typ := "ALooper"

	ctorNames := []string{typ + "_create", typ + "_new"}
	for _, name := range ctorNames {
		if !isConstructorFunc(name, typ) {
			t.Errorf("isConstructorFunc(%q, %q) = false, want true", name, typ)
		}
	}

	dtorNames := []string{
		typ + "_delete", typ + "_free", typ + "_destroy",
		typ + "_release", typ + "_close",
	}
	for _, name := range dtorNames {
		if !isDestructorFunc(name, typ) {
			t.Errorf("isDestructorFunc(%q, %q) = false, want true", name, typ)
		}
	}
}

// TestCountGoParams_Property verifies countGoParams matches comma-counting.
func TestCountGoParams_Property(t *testing.T) {
	cases := []struct {
		sig  string
		want int
	}{
		{"func()", 0},
		{"func(int)", 1},
		{"func(int, string)", 2},
		{"func(int, string, bool)", 3},
		{"func( )", 0},
		{"hello", 0},
		{"", 0},
	}
	for _, tc := range cases {
		if got := countGoParams(tc.sig); got != tc.want {
			t.Errorf("countGoParams(%q) = %d, want %d", tc.sig, got, tc.want)
		}
	}
}

// TestAutoGoTypeName_StripAndroidPrefix verifies the "A" prefix stripping.
func TestAutoGoTypeName_StripAndroidPrefix(t *testing.T) {
	cases := []struct{ in_, want string }{
		{"AImageReader", "ImageReader"},
		{"ALooper", "Looper"},
		{"ASensor", "Sensor"},
		{"GLenum", "GLenum"}, // No strip — not A+uppercase.
		{"Abcdef", "Abcdef"}, // No strip — A+lowercase b.
		{"A", "A"},           // Too short.
		{"AB", "B"},          // Strip — A+B.
		{"a_type", "a_type"}, // lowercase a, no strip.
	}
	for _, tc := range cases {
		if got := autoGoTypeName(tc.in_); got != tc.want {
			t.Errorf("autoGoTypeName(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestInferEnumType_OnlyInt32 verifies inference only applies to int32 params.
func TestInferEnumType_OnlyInt32(t *testing.T) {
	enums := map[string]bool{"Type": true, "Status": true}

	// Only int32 params get inference.
	if got := inferEnumType("type", "int32", enums); got != "Type" {
		t.Errorf("inferEnumType(type, int32) = %q, want Type", got)
	}

	// Non-int32 types are unchanged.
	if got := inferEnumType("type", "uint32", enums); got != "uint32" {
		t.Errorf("inferEnumType(type, uint32) = %q, want uint32", got)
	}
}

// TestToSnakeCase verifies PascalCase → snake_case conversion.
func TestToSnakeCase_Property(t *testing.T) {
	f := func(s string) bool {
		result := toSnakeCase(s)
		// Result should never contain uppercase letters.
		for _, r := range result {
			if r >= 'A' && r <= 'Z' {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("toSnakeCase produced uppercase: %v", err)
	}
}

// TestToSnakeCase_Vectors verifies specific conversions.
func TestToSnakeCase_Vectors(t *testing.T) {
	cases := []struct{ in_, want string }{
		{"StreamBuilder", "stream_builder"},
		{"Model", "model"},
		{"AUDIO", "audio"},
	}
	for _, tc := range cases {
		if got := toSnakeCase(tc.in_); got != tc.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestFixGoAcronyms verifies "Id" → "ID" normalization.
func TestFixGoAcronyms(t *testing.T) {
	cases := []struct{ in_, want string }{
		{"DeviceId", "DeviceID"},
		{"GetId", "GetID"},
		{"Idle", "Idle"},         // "Id" followed by lowercase — no change.
		{"Identity", "Identity"}, // "Id" followed by lowercase — no change.
		{"Id", "ID"},             // "Id" at end.
		{"CameraIdList", "CameraIDList"},
	}
	for _, tc := range cases {
		if got := fixGoAcronyms(tc.in_); got != tc.want {
			t.Errorf("fixGoAcronyms(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestCapiExportType_ScalarUnchanged verifies scalars are passed through.
func TestCapiExportType_ScalarUnchanged(t *testing.T) {
	scalars := []string{
		"int8", "uint8", "int16", "uint16", "int32", "uint32",
		"int64", "uint64", "float32", "float64", "bool", "int", "uint",
	}
	for _, s := range scalars {
		if got := capiExportType(s); got != s {
			t.Errorf("capiExportType(%q) = %q, want unchanged", s, got)
		}
	}
}

// TestCapiExportType_PointerPreservation verifies pointer prefixes are preserved.
func TestCapiExportType_PointerPreservation(t *testing.T) {
	cases := []struct{ in_, want string }{
		{"*camera_status_t", "*Camera_status_t"},
		{"*int32", "*int32"},
		{"[]int32", "[]int32"},
		{"unsafe.Pointer", "unsafe.Pointer"},
		{"string", "string"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := capiExportType(tc.in_); got != tc.want {
			t.Errorf("capiExportType(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestIsTypeRemapped verifies type remapping detection.
func TestIsTypeRemapped(t *testing.T) {
	typeMap := map[string]string{"foo": "Bar"}

	cases := []struct {
		specType string
		want     bool
	}{
		{"foo", true},
		{"*foo", true},
		{"**foo", true},
		{"[]foo", true},
		{"bar", false},
		{"*bar", false},
	}
	for _, tc := range cases {
		if got := isTypeRemapped(tc.specType, typeMap); got != tc.want {
			t.Errorf("isTypeRemapped(%q) = %v, want %v", tc.specType, got, tc.want)
		}
	}
}

// TestStripAndTitle verifies UPPER_SNAKE_CASE to TitleCase conversion.
func TestStripAndTitle(t *testing.T) {
	cases := []struct {
		name, prefix, want string
	}{
		{"AAUDIO_DIRECTION_OUTPUT", "AAUDIO_DIRECTION_", "Output"},
		{"AAUDIO_DIRECTION_INPUT", "AAUDIO_DIRECTION_", "Input"},
		{"FOO_BAR_BAZ", "FOO_BAR_", "Baz"},
		{"FOO_BAR_BAZ", "", "FOO_BAR_BAZ"},
		{"FOO", "FOO", "FOO"}, // Empty remainder returns original.
	}
	for _, tc := range cases {
		if got := stripAndTitle(tc.name, tc.prefix); got != tc.want {
			t.Errorf("stripAndTitle(%q, %q) = %q, want %q", tc.name, tc.prefix, got, tc.want)
		}
	}
}

// TestDestructorPriority verifies the priority ordering.
func TestDestructorPriority(t *testing.T) {
	// release > delete > free > destroy > close
	priorities := []struct {
		suffix string
		want   int
	}{
		{"_release", 5},
		{"_delete", 4},
		{"_free", 3},
		{"_destroy", 2},
		{"_close", 1},
		{"_other", 0},
	}
	for _, tc := range priorities {
		if got := destructorPriority("Type" + tc.suffix); got != tc.want {
			t.Errorf("destructorPriority(%q) = %d, want %d", "Type"+tc.suffix, got, tc.want)
		}
	}

	// Verify strict ordering.
	for i := 0; i < len(priorities)-1; i++ {
		for j := i + 1; j < len(priorities); j++ {
			a := destructorPriority("T" + priorities[i].suffix)
			b := destructorPriority("T" + priorities[j].suffix)
			if a <= b {
				t.Errorf("destructorPriority(%q) = %d should be > destructorPriority(%q) = %d",
					priorities[i].suffix, a, priorities[j].suffix, b)
			}
		}
	}
}

// TestIsAutoDetectedPure verifies the auto-detection heuristic.
func TestIsAutoDetectedPure(t *testing.T) {
	errorTypes := map[string]bool{"camera_status_t": true}

	pureTypes := []string{"string", "bool", "float32", "float64", "int64", "uint64", "unsafe.Pointer"}
	for _, typ := range pureTypes {
		if !isAutoDetectedPure(typ, errorTypes) {
			t.Errorf("isAutoDetectedPure(%q) = false, want true", typ)
		}
	}

	notPureTypes := []string{"", "void", "int32", "uint32", "camera_status_t", "*ALooper"}
	for _, typ := range notPureTypes {
		if isAutoDetectedPure(typ, errorTypes) {
			t.Errorf("isAutoDetectedPure(%q) = true, want false", typ)
		}
	}

	// *uint8 (pointer to scalar) should be pure.
	if !isAutoDetectedPure("*uint8", errorTypes) {
		t.Error("isAutoDetectedPure(*uint8) = false, want true")
	}
}

// TestDurationParamName verifies duration parameter name normalization.
func TestDurationParamName(t *testing.T) {
	cases := []struct{ in_, want string }{
		{"timeoutNanoseconds", "timeout"},
		{"actualDurationNanos", "actualDuration"},
		{"initialTargetWorkDurationNanos", "initialTargetWorkDuration"},
		{"timeoutMillis", "timeout"},
		{"delayMicroseconds", "delay"},
		{"someOtherParam", "someOtherParam"},
	}
	for _, tc := range cases {
		if got := durationParamName(tc.in_); got != tc.want {
			t.Errorf("durationParamName(%q) = %q, want %q", tc.in_, got, tc.want)
		}
	}
}

// TestCapiArg_MultiLevelPointer verifies the fix for the multi-level pointer bug.
func TestCapiArg_MultiLevelPointer(t *testing.T) {
	p := MergedParam{
		Name:     "x",
		GoType:   "**Foo",
		CapiType: "**Foo",
		Remapped: true,
	}
	got := capiArg(p)
	want := "(**capi.Foo)(x)"
	if got != want {
		t.Errorf("capiArg with **Foo = %q, want %q", got, want)
	}

	// Single pointer should also work.
	p2 := MergedParam{
		Name:     "y",
		GoType:   "*Bar",
		CapiType: "*Bar",
		Remapped: true,
	}
	got2 := capiArg(p2)
	want2 := "(*capi.Bar)(y)"
	if got2 != want2 {
		t.Errorf("capiArg with *Bar = %q, want %q", got2, want2)
	}
}

// TestCapiArgOrFixed_MultiLevelPointer verifies the fix applied to capiArgOrFixed.
func TestCapiArgOrFixed_MultiLevelPointer(t *testing.T) {
	fm := FuncMap()
	capiArgOrFixed := fm["capiArgOrFixed"].(func(MergedParam, map[string]string) string)

	p := MergedParam{
		Name:     "x",
		GoType:   "**Foo",
		CapiType: "**Foo",
		Remapped: true,
	}
	got := capiArgOrFixed(p, nil)
	// After fix: should produce (**capi.Foo)(x), not (*capi.*Foo)(x).
	if strings.Contains(got, "*capi.*") {
		t.Errorf("capiArgOrFixed produced invalid Go: %q", got)
	}
	want := "(**capi.Foo)(x)"
	if got != want {
		t.Errorf("capiArgOrFixed = %q, want %q", got, want)
	}
}

// TestCallbackCapiArg_MultiLevelPointer verifies the fix applied to callbackCapiArg.
func TestCallbackCapiArg_MultiLevelPointer(t *testing.T) {
	fm := FuncMap()
	callbackCapiArgFn := fm["callbackCapiArg"].(func(MergedParam, string) string)

	p := MergedParam{
		Name:     "x",
		GoType:   "**Foo",
		CapiType: "**Foo",
		Remapped: true,
	}
	got := callbackCapiArgFn(p, "other")
	if strings.Contains(got, "*capi.*") {
		t.Errorf("callbackCapiArg produced invalid Go: %q", got)
	}
	want := "(**capi.Foo)(x)"
	if got != want {
		t.Errorf("callbackCapiArg = %q, want %q", got, want)
	}
}
