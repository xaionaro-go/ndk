package idiomgen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/xaionaro-go/ndk/tools/pkg/specmodel"
)

// capiArg returns the expression to pass a parameter to a capi function.
// Handles need .ptr unwrapping; remapped types need capi.Type() conversion.
func capiArg(p MergedParam) string {
	if p.GoType == "time.Duration" {
		unit := p.DurationUnit
		if unit == "" {
			unit = "ns"
		}
		// All Duration methods return int64; cast to match the capi param type.
		switch unit {
		case "ms":
			return p.CapiType + "(" + p.Name + ".Milliseconds())"
		case "us":
			return p.CapiType + "(" + p.Name + ".Microseconds())"
		case "s":
			return p.CapiType + "(" + p.Name + ".Seconds())"
		default: // "ns"
			return p.CapiType + "(" + p.Name + ".Nanoseconds())"
		}
	}
	if p.IsHandle && p.GoType != "unsafe.Pointer" {
		return p.Name + ".ptr"
	}
	if p.IsString {
		return "&" + p.Name + "Bytes[0]"
	}
	if p.CapiType != "" && (p.CapiType != p.GoType || p.Remapped) {
		if strings.HasPrefix(p.CapiType, "[]") {
			// Slice of a renamed type: reinterpret via unsafe.Pointer.
			// []Int and []capi.EGLint have identical layout.
			return "*(*[]capi." + p.CapiType[2:] + ")(unsafe.Pointer(&" + p.Name + "))"
		}
		if strings.HasPrefix(p.CapiType, "*") {
			// Pointer type: strip all stars, prefix with capi., re-add stars outside.
			// E.g. *Type → (*capi.Type)(name), **Type → (**capi.Type)(name)
			stars := 0
			base := p.CapiType
			for strings.HasPrefix(base, "*") {
				stars++
				base = base[1:]
			}
			return "(" + strings.Repeat("*", stars) + "capi." + base + ")(" + p.Name + ")"
		}
		// Scalar Go types (int32, uint32, etc.) don't need a capi. prefix —
		// they're built-in types, not capi package types.
		if isScalarGoType(p.CapiType) {
			return p.CapiType + "(" + p.Name + ")"
		}
		return "capi." + p.CapiType + "(" + p.Name + ")"
	}
	return p.Name
}

// FuncMap returns template helper functions for code generation templates.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"safeGoName": safeGoName,
		"unexport":   unexport,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"trimPrefix": strings.TrimPrefix,
		"skip": func(params []MergedParam, n int) []MergedParam {
			if n >= len(params) {
				return nil
			}
			return params[n:]
		},
		// skipAndFilterOut filters out params with direction "out".
		// The receiver is already stripped by the merger, so no index-based skipping.
		"skipAndFilterOut": func(params []MergedParam) []MergedParam {
			var out []MergedParam
			for _, p := range params {
				if p.Direction == "out" {
					continue
				}
				out = append(out, p)
			}
			return out
		},
		"lookupCapiType": func(types map[string]MergedOpaqueType, goName string) string {
			if t, ok := types[goName]; ok {
				return t.CapiType
			}
			return goName
		},
		// sortedOpaqueTypes returns opaque types sorted by Go name for deterministic output.
		"sortedOpaqueTypes": func(m map[string]MergedOpaqueType) []MergedOpaqueType {
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			result := make([]MergedOpaqueType, len(keys))
			for i, k := range keys {
				result[i] = m[k]
			}
			return result
		},
		// anyStringMethod returns true if any value enum has StringMethod enabled.
		"anyStringMethod": func(enums []MergedValueEnum) bool {
			for _, e := range enums {
				if e.StringMethod {
					return true
				}
			}
			return false
		},
		// anyUnsafeParam returns true if any method has an unsafe.Pointer parameter.
		"anyUnsafeParam": func(methods []MergedMethod) bool {
			for _, m := range methods {
				for _, p := range m.Params {
					if strings.Contains(p.GoType, "unsafe.Pointer") {
						return true
					}
				}
			}
			return false
		},
		// anyDurationParam returns true if any method or free function has a time.Duration parameter.
		"anyDurationParam": func(methods []MergedMethod, fns []MergedFreeFunction) bool {
			for _, m := range methods {
				for _, p := range m.Params {
					if p.GoType == "time.Duration" {
						return true
					}
				}
			}
			for _, f := range fns {
				for _, p := range f.Params {
					if p.GoType == "time.Duration" {
						return true
					}
				}
			}
			return false
		},
		// anyDurationFreeFunc returns true if any free function has a time.Duration parameter.
		"anyDurationFreeFunc": func(fns []MergedFreeFunction) bool {
			for _, f := range fns {
				for _, p := range f.Params {
					if p.GoType == "time.Duration" {
						return true
					}
				}
			}
			return false
		},
		// anyUnsafeFreeFunc returns true if any free function has an unsafe.Pointer parameter or return,
		// or if any parameter needs a slice type reinterpret cast.
		"anyUnsafeFreeFunc": func(fns []MergedFreeFunction) bool {
			for _, f := range fns {
				if strings.Contains(f.Returns, "unsafe.Pointer") {
					return true
				}
				for _, p := range f.Params {
					if strings.Contains(p.GoType, "unsafe.Pointer") {
						return true
					}
					// Slice of a renamed type needs unsafe reinterpret.
					if p.CapiType != "" && p.CapiType != p.GoType && strings.HasPrefix(p.CapiType, "[]") {
						return true
					}
				}
			}
			return false
		},
		// capiArg returns the expression to pass a parameter to a capi function.
		// Handles need .ptr unwrapping; remapped types need capi.Type() conversion.
		"capiArg": capiArg,
		// capiArgBuf returns the expression for a buffer parameter, converting
		// typed slices (e.g., []byte) to unsafe.Pointer for capi calls.
		"capiArgBuf": func(p MergedParam, bufGoType string) string {
			if strings.HasPrefix(p.GoType, "[]") {
				return "unsafe.Pointer(&" + p.Name + "[0])"
			}
			return capiArg(p)
		},
		// trimStar removes a leading "*" from a string.
		"trimStar": func(s string) string {
			return strings.TrimPrefix(s, "*")
		},
		"toSnakeCase": toSnakeCase,
		// skipSpecParams skips the first n parameters from a specmodel.Param slice.
		"skipSpecParams": func(params []specmodel.Param, n int) []specmodel.Param {
			if n >= len(params) {
				return nil
			}
			return params[n:]
		},
		// cTypeToGo maps a C/spec type string to a Go CGo type for use in //export functions.
		"cTypeToGo": func(cType string) string {
			if cType == "void*" || cType == "unsafe.Pointer" {
				return "unsafe.Pointer"
			}
			// Go integer types → C integer types.
			goToCGoMap := map[string]string{
				"int8": "C.int8_t", "uint8": "C.uint8_t",
				"int16": "C.int16_t", "uint16": "C.uint16_t",
				"int32": "C.int32_t", "uint32": "C.uint32_t",
				"int64": "C.int64_t", "uint64": "C.uint64_t",
				"float32": "C.float", "float64": "C.double",
				"int": "C.int",
			}
			if cgoType, ok := goToCGoMap[cType]; ok {
				return cgoType
			}
			if strings.HasPrefix(cType, "*") {
				base := cType[1:]
				if cgoBase, ok := goToCGoMap[base]; ok {
					return "*" + cgoBase
				}
				return "*C." + base
			}
			return "C." + cType
		},
		// toCType converts a spec type (Go-style) to C declaration syntax.
		// "*ACameraDevice" becomes "ACameraDevice*", "void*" stays "void*",
		// "int" stays "int", "**ACameraDevice" becomes "ACameraDevice**".
		"toCType": toCType,
		// filterFixedParams returns params whose names are not in fixedParams.
		"filterFixedParams": func(params []MergedParam, fixedParams map[string]string) []MergedParam {
			if len(fixedParams) == 0 {
				return params
			}
			var out []MergedParam
			for _, p := range params {
				if _, fixed := fixedParams[p.Name]; !fixed {
					out = append(out, p)
				}
			}
			return out
		},
		// capiArgOrFixed returns the fixed literal value if the param is fixed,
		// otherwise returns the normal capiArg expression.
		"capiArgOrFixed": func(p MergedParam, fixedParams map[string]string) string {
			if v, ok := fixedParams[p.Name]; ok {
				return v
			}
			if p.IsHandle {
				return p.Name + ".ptr"
			}
			if p.IsString {
				return "&" + p.Name + "Bytes[0]"
			}
			if p.CapiType != "" && (p.CapiType != p.GoType || p.Remapped) {
				if strings.HasPrefix(p.CapiType, "*") {
					return "(*capi." + p.CapiType[1:] + ")(" + p.Name + ")"
				}
				if isScalarGoType(p.CapiType) {
					return p.CapiType + "(" + p.Name + ")"
				}
				return "capi." + p.CapiType + "(" + p.Name + ")"
			}
			return p.Name
		},
		// callbackVisibleParams returns params visible in the Go signature for callback methods.
		// Excludes direction:out params and the callback param itself. Converts *byte to string.
		"callbackVisibleParams": func(params []MergedParam, callbackParam string) []MergedParam {
			var out []MergedParam
			for _, p := range params {
				if p.Direction == "out" || p.Name == callbackParam {
					continue
				}
				out = append(out, p)
			}
			return out
		},
		// callbackCapiArg returns the capi expression for a param in a callback method call.
		"callbackCapiArg": func(p MergedParam, callbackParam string) string {
			if p.Name == callbackParam {
				return "&cbsC"
			}
			if p.IsString {
				return "&" + p.Name + "Bytes[0]"
			}
			if p.IsHandle {
				return p.Name + ".ptr"
			}
			if p.CapiType != "" && (p.CapiType != p.GoType || p.Remapped) {
				if strings.HasPrefix(p.CapiType, "*") {
					return "(*capi." + p.CapiType[1:] + ")(" + p.Name + ")"
				}
				if isScalarGoType(p.CapiType) {
					return p.CapiType + "(" + p.Name + ")"
				}
				return "capi." + p.CapiType + "(" + p.Name + ")"
			}
			return p.Name
		},
		// filterOutputParams returns params whose names are not in the output params set.
		"filterOutputParams": func(params []MergedParam, outputParams []MergedOutputParam) []MergedParam {
			if len(outputParams) == 0 {
				return params
			}
			opNames := make(map[string]bool, len(outputParams))
			for _, op := range outputParams {
				opNames[op.CParamName] = true
			}
			var out []MergedParam
			for _, p := range params {
				if !opNames[p.Name] {
					out = append(out, p)
				}
			}
			return out
		},
		// trimStarPrefix removes the leading "*" from a type string (e.g., "*ImageReader" -> "ImageReader").
		"trimStarPrefix": func(s string) string {
			return strings.TrimPrefix(s, "*")
		},
		// zeroValue returns the Go zero value literal for a type.
		// Pointer and handle types return "nil", numeric types return "0".
		"zeroValue": func(op MergedOutputParam) string {
			if strings.HasPrefix(op.GoType, "*") || op.IsHandle {
				return "nil"
			}
			return "0"
		},
		// stringParams returns only params that need string→byte conversion.
		"stringParams": func(params []MergedParam) []MergedParam {
			var out []MergedParam
			for _, p := range params {
				if p.IsString {
					out = append(out, p)
				}
			}
			return out
		},
		// add returns a + b. Used for index arithmetic in templates.
		"add": func(a, b int) int { return a + b },
		// lifecycleGoParamType converts a C param type to a Go type for lifecycle callbacks.
		// Pointer types → unsafe.Pointer, int → int, size_t → uintptr, etc.
		"lifecycleGoParamType": func(cType string) string {
			switch cType {
			case "int", "int32_t":
				return "int32"
			case "int64_t":
				return "int64"
			case "size_t":
				return "uintptr"
			case "void*", "unsafe.Pointer":
				return "unsafe.Pointer"
			default:
				if strings.HasPrefix(cType, "*") {
					return "unsafe.Pointer"
				}
				return cType
			}
		},
		// lifecycleConvertParam returns Go expression to convert a C param for lifecycle callback dispatch.
		"lifecycleConvertParam": func(p specmodel.Param) string {
			switch p.Type {
			case "int", "int32_t", "int32":
				return "int32(" + p.Name + ")"
			case "int64_t", "int64":
				return "int64(" + p.Name + ")"
			case "uint32_t", "uint32":
				return "uint32(" + p.Name + ")"
			case "uint64_t", "uint64":
				return "uint64(" + p.Name + ")"
			case "size_t":
				return "uintptr(" + p.Name + ")"
			default:
				return "unsafe.Pointer(" + p.Name + ")"
			}
		},
		// toCParamDecl converts a specmodel.Param to a C declaration string (e.g., "const ARect* rect").
		"toCParamDecl": func(p specmodel.Param) string {
			cType := toCType(p.Type)
			if p.Const {
				return "const " + cType + " " + p.Name
			}
			return cType + " " + p.Name
		},
		// tailSpecParams returns the last n params from a specmodel.Param slice.
		"tailSpecParams": func(params []specmodel.Param, n int) []specmodel.Param {
			if n >= len(params) {
				return params
			}
			return params[len(params)-n:]
		},
		// goConvertSpecParam converts a C param to a Go expression for callback dispatch.
		"goConvertSpecParam": func(p specmodel.Param) string {
			switch p.Type {
			case "int":
				return "int(" + p.Name + ")"
			case "int32_t", "int32":
				return "int(" + p.Name + ")"
			case "int64_t", "int64":
				return "int64(" + p.Name + ")"
			case "uint32_t", "uint32":
				return "uint32(" + p.Name + ")"
			case "void*", "unsafe.Pointer":
				return p.Name
			default:
				if strings.HasPrefix(p.Type, "*") {
					return "unsafe.Pointer(" + p.Name + ")"
				}
				return p.Name
			}
		},
	}
}

// unexport lowercases the first rune of a string.
func unexport(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// toCType converts a spec type (Go-style) to C declaration syntax.
// "*ACameraDevice" becomes "ACameraDevice*", "void*" stays "void*".
func toCType(specType string) string {
	if specType == "void*" || specType == "unsafe.Pointer" {
		return "void*"
	}
	stars := 0
	base := specType
	for strings.HasPrefix(base, "*") {
		stars++
		base = base[1:]
	}
	// Convert Go integer types to C types.
	goToCMap := map[string]string{
		"int8": "int8_t", "uint8": "uint8_t",
		"int16": "int16_t", "uint16": "uint16_t",
		"int32": "int32_t", "uint32": "uint32_t",
		"int64": "int64_t", "uint64": "uint64_t",
		"float32": "float", "float64": "double",
		"int": "int",
	}
	if cType, ok := goToCMap[base]; ok {
		base = cType
	}
	if stars == 0 {
		return base
	}
	return base + strings.Repeat("*", stars)
}

// goReservedWords is the set of Go keywords and predeclared identifiers that
// cannot be used as parameter names.
var goReservedWords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

// safeGoName returns a Go-safe version of a name by prefixing reserved words
// with _ and normalizing common acronyms (e.g. "deviceId" → "deviceID").
func safeGoName(name string) string {
	if goReservedWords[name] {
		return "_" + name
	}
	return fixGoAcronyms(name)
}

// safeGoParamName returns a Go-safe parameter name, synthesizing positional
// names for unnamed C parameters: index 0 → "arg0", index 1 → "arg1", etc.
func safeGoParamName(name string, index int) string {
	if name == "" {
		return fmt.Sprintf("arg%d", index)
	}
	return safeGoName(name)
}

// acronymRe matches "Id" when it appears as a complete word segment
// (at end of string or followed by a non-lowercase letter).
// This avoids false positives like "Idle", "Identify", "Identity".
var acronymRe = regexp.MustCompile(`Id([^a-z]|$)`)

// fixGoAcronyms normalizes common acronyms in Go identifiers to follow
// Go conventions (e.g. "Id" → "ID").
func fixGoAcronyms(s string) string {
	return acronymRe.ReplaceAllString(s, "ID$1")
}

// toSnakeCase converts PascalCase or UPPER_CASE to snake_case.
// Examples: StreamBuilder → stream_builder, Model → model, AUDIO → audio,
// LOOPER_POLL → looper_poll, IMAGE_FORMATS → image_formats.
func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			// Insert underscore only when transitioning from a lowercase
			// letter or digit to an uppercase letter. Consecutive uppercase
			// letters (e.g. "AUDIO") are treated as a single word.
			if i > 0 {
				prev := rune(s[i-1])
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					b.WriteByte('_')
				}
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
