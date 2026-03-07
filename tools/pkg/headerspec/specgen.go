package headerspec

import (
	"fmt"
	"os"
	"strings"

	"github.com/xaionaro-go/ndk/tools/pkg/specmodel"
	"gopkg.in/yaml.v3"
)

// GenerateSpec converts filtered declarations into a specmodel.Spec.
func GenerateSpec(
	moduleName string,
	sourcePackage string,
	decls *Declarations,
) *specmodel.Spec {
	spec := &specmodel.Spec{
		Module:        moduleName,
		SourcePackage: sourcePackage,
		Types:         make(map[string]specmodel.TypeDef),
		Enums:         make(map[string][]specmodel.EnumValue),
		Functions:     make(map[string]specmodel.FuncDef),
		Callbacks:     make(map[string]specmodel.CallbackDef),
		Structs:       make(map[string]specmodel.StructDef),
	}

	// Index typedefs by name so we can look them up when resolving
	// parameter types that reference typedefs.
	typedefByName := map[string]*TypedefInfo{}
	for i := range decls.Typedefs {
		td := &decls.Typedefs[i]
		typedefByName[td.Name] = td
	}

	// Process typedefs → types and callbacks.
	for _, td := range decls.Typedefs {
		if td.IsFuncPtr {
			cb := specmodel.CallbackDef{
				Returns: cTypeToGoType(td.FuncReturn),
			}
			for i, p := range td.FuncParams {
				cb.Params = append(cb.Params, specmodel.Param{
					Name: paramName(p.Name, i),
					Type: cTypeToGoType(p.Type),
				})
			}
			spec.Callbacks[td.Name] = cb
			continue
		}

		typeDef := classifyTypedef(&td)
		spec.Types[td.Name] = typeDef
	}

	// Process enums.
	for _, ei := range decls.Enums {
		groupName := ei.TypedefName
		if groupName == "" {
			groupName = ei.Name
		}
		if groupName == "" {
			// Anonymous enum without a typedef: derive a group name from
			// the common prefix of its constant names.
			groupName = deriveEnumGroupName(ei.Constants)
		}
		if groupName == "" {
			continue
		}

		var values []specmodel.EnumValue
		for _, c := range ei.Constants {
			values = append(values, specmodel.EnumValue{
				Name:  c.Name,
				Value: c.Value,
			})
		}
		spec.Enums[groupName] = values
	}

	// Process functions.
	for _, fn := range decls.Functions {
		fd := specmodel.FuncDef{
			CName:   fn.Name,
			Returns: cTypeToGoType(fn.ReturnType),
		}

		for _, p := range fn.Params {
			sp := cParamToSpecParam(p, typedefByName)
			fd.Params = append(fd.Params, sp)
		}

		spec.Functions[fn.Name] = fd
	}

	// Process complete structs.
	for _, si := range decls.Structs {
		if !si.IsComplete {
			continue
		}

		sd := specmodel.StructDef{}
		for _, fi := range si.Fields {
			sf := specmodel.StructField{
				Name: fi.Name,
				Type: cTypeToGoType(fi.Type),
			}

			if fi.IsFuncPtr {
				sf.Type = "func_ptr"
				sf.Returns = cTypeToGoType(fi.FuncReturn)
				for i, p := range fi.FuncParams {
					sf.Params = append(sf.Params, specmodel.Param{
						Name: paramName(p.Name, i),
						Type: cTypeToGoType(p.Type),
					})
				}
			}

			sd.Fields = append(sd.Fields, sf)
		}

		spec.Structs[si.Name] = sd
	}

	// Remove empty maps so they don't appear in YAML output.
	if len(spec.Types) == 0 {
		spec.Types = nil
	}
	if len(spec.Enums) == 0 {
		spec.Enums = nil
	}
	if len(spec.Functions) == 0 {
		spec.Functions = nil
	}
	if len(spec.Callbacks) == 0 {
		spec.Callbacks = nil
	}
	if len(spec.Structs) == 0 {
		spec.Structs = nil
	}

	return spec
}

// classifyTypedef determines the specmodel.TypeDef for a non-function-pointer typedef.
func classifyTypedef(td *TypedefInfo) specmodel.TypeDef {
	ut := td.UnderlyingType

	// typedef struct X X → opaque_ptr
	if td.IsOpaqueStruct {
		return specmodel.TypeDef{
			Kind:   "opaque_ptr",
			CType:  td.Name,
			GoType: "*C." + td.Name,
		}
	}

	// typedef void * → pointer_handle
	if ut == "void *" {
		return specmodel.TypeDef{
			Kind:   "pointer_handle",
			CType:  td.Name,
			GoType: "unsafe.Pointer",
		}
	}

	// typedef enum X → typedef_int32 (C enums are int-sized by default).
	if strings.HasPrefix(ut, "enum ") {
		return specmodel.TypeDef{
			Kind:   "typedef_int32",
			CType:  td.Name,
			GoType: "int32",
		}
	}

	// typedef to integer types → typedef_intN / typedef_uintN
	goType := cTypeToGoType(ut)
	kind := goTypeToKind(goType)

	return specmodel.TypeDef{
		Kind:   kind,
		CType:  td.Name,
		GoType: goType,
	}
}

// goTypeToKind converts a Go type string to a spec kind string.
func goTypeToKind(goType string) string {
	switch goType {
	case "int8":
		return "typedef_int8"
	case "uint8":
		return "typedef_uint8"
	case "int16":
		return "typedef_int16"
	case "uint16":
		return "typedef_uint16"
	case "int32":
		return "typedef_int32"
	case "uint32":
		return "typedef_uint32"
	case "int64":
		return "typedef_int64"
	case "uint64":
		return "typedef_uint64"
	case "bool":
		return "typedef_bool"
	case "float32":
		return "typedef_float32"
	case "float64":
		return "typedef_float64"
	default:
		return "typedef_" + goType
	}
}

// cParamToSpecParam converts a C parameter to a spec Param, applying type
// conversions and detecting output parameters (double pointers).
func cParamToSpecParam(
	p ParamInfo,
	typedefByName map[string]*TypedefInfo,
) specmodel.Param {
	sp := specmodel.Param{
		Name: p.Name,
	}

	cType := stripNullabilityAnnotations(p.Type)

	// Detect const qualifier.
	if strings.HasPrefix(cType, "const ") {
		sp.Const = true
	}

	// Handle special pointer types before generic pointer logic.
	switch cType {
	case "void *", "const void *":
		sp.Type = "unsafe.Pointer"
		return sp
	case "void **":
		sp.Type = "[]unsafe.Pointer"
		return sp
	case "const char *", "const char *const", "char *":
		sp.Type = "string"
		return sp
	}

	// Detect output parameters: double pointer "X **"
	if strings.HasSuffix(cType, " **") || strings.HasSuffix(cType, "**") {
		base := strings.TrimSpace(strings.TrimSuffix(cType, "**"))
		base = strings.TrimPrefix(base, "const ")
		sp.Type = "*" + cTypeToGoType(base+" *")
		sp.Direction = "out"
		return sp
	}

	// Single pointer to known type: "int *" → "[]int32" (pointer-to-scalar = out/array).
	if strings.HasSuffix(cType, " *") {
		base := strings.TrimSpace(strings.TrimSuffix(cType, " *"))
		base = strings.TrimPrefix(base, "const ")

		// Check if base is a typedef that we know.
		if td, ok := typedefByName[base]; ok {
			if td.IsOpaqueStruct {
				sp.Type = "*" + base
				return sp
			}
		}

		// Pointer to scalar types become slices (matching specgen behavior).
		goBase := cTypeToGoType(base)
		if isScalarGoType(goBase) {
			sp.Type = "[]" + goBase
			return sp
		}

		sp.Type = "*" + goBase
		return sp
	}

	sp.Type = cTypeToGoType(cType)
	return sp
}

// isScalarGoType returns true for Go scalar/primitive types.
func isScalarGoType(t string) bool {
	switch t {
	case "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64",
		"float32", "float64", "bool":
		return true
	}
	return false
}

// stripNullabilityAnnotations removes Clang nullability annotations from a type string.
func stripNullabilityAnnotations(s string) string {
	s = strings.ReplaceAll(s, " _Nonnull", "")
	s = strings.ReplaceAll(s, " _Nullable", "")
	s = strings.ReplaceAll(s, " _Null_unspecified", "")
	return strings.TrimSpace(s)
}

// cTypeToGoType converts a C type string to a Go type string.
func cTypeToGoType(cType string) string {
	cType = strings.TrimSpace(cType)
	cType = stripNullabilityAnnotations(cType)

	// Handle pointer types first.
	if cType == "void *" || cType == "const void *" {
		return "unsafe.Pointer"
	}
	if cType == "const char *" || cType == "const char *const" || cType == "const char *const *" {
		return "string"
	}
	if cType == "char *" {
		return "string"
	}

	// Strip pointer suffix and recurse.
	if strings.HasSuffix(cType, " *") {
		base := strings.TrimSpace(strings.TrimSuffix(cType, " *"))
		return "*" + cTypeToGoType(base)
	}
	if strings.HasSuffix(cType, "*") && !strings.Contains(cType, "(") {
		base := strings.TrimSpace(strings.TrimSuffix(cType, "*"))
		return "*" + cTypeToGoType(base)
	}

	// Strip "const " prefix for base type lookups.
	base := strings.TrimPrefix(cType, "const ")

	// Strip "struct " prefix.
	base = strings.TrimPrefix(base, "struct ")

	// Strip "enum " prefix.
	base = strings.TrimPrefix(base, "enum ")

	switch base {
	case "void":
		return ""
	case "int":
		return "int32"
	case "unsigned int", "unsigned":
		return "uint32"
	case "long":
		return "int64"
	case "unsigned long":
		return "uint64"
	case "long long":
		return "int64"
	case "unsigned long long":
		return "uint64"
	case "short":
		return "int16"
	case "unsigned short":
		return "uint16"
	case "char", "signed char":
		return "int8"
	case "unsigned char":
		return "uint8"
	case "int8_t":
		return "int8"
	case "uint8_t":
		return "uint8"
	case "int16_t":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "int32_t":
		return "int32"
	case "uint32_t":
		return "uint32"
	case "int64_t":
		return "int64"
	case "uint64_t":
		return "uint64"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "size_t":
		return "uint64"
	case "ssize_t":
		return "int64"
	case "_Bool", "bool":
		return "bool"
	}

	// If it's a known NDK typedef name (like ALooper), return as-is.
	return base
}

// paramName returns a parameter name, using a generated placeholder if empty.
func paramName(name string, index int) string {
	if name != "" {
		return name
	}
	return fmt.Sprintf("p%d", index)
}

// deriveEnumGroupName computes a group name from the common prefix of
// enum constant names. For example, constants ["ALOOPER_POLL_WAKE",
// "ALOOPER_POLL_CALLBACK"] yield "ALOOPER_POLL".
//
// For single-constant enums, the prefix up to the second-to-last
// underscore is used (e.g., "ALOOPER_PREPARE_ALLOW_NON_CALLBACKS"
// yields "ALOOPER_PREPARE").
func deriveEnumGroupName(constants []EnumConstant) string {
	if len(constants) == 0 {
		return ""
	}
	if len(constants) == 1 {
		// Single constant: find a meaningful prefix by dropping the
		// last two underscore-delimited segments.
		name := constants[0].Name
		parts := strings.Split(name, "_")
		if len(parts) > 2 {
			return strings.Join(parts[:len(parts)-2], "_")
		}
		if idx := strings.LastIndex(name, "_"); idx > 0 {
			return name[:idx]
		}
		return name
	}

	prefix := constants[0].Name
	for _, c := range constants[1:] {
		prefix = commonPrefix(prefix, c.Name)
	}

	// Trim trailing underscore.
	prefix = strings.TrimRight(prefix, "_")

	if prefix == "" {
		return ""
	}

	return prefix
}

// commonPrefix returns the common prefix of two strings.
func commonPrefix(a, b string) string {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:n]
}

// WriteSpecYAML marshals a Spec to YAML and writes it to the given path.
func WriteSpecYAML(spec *specmodel.Spec, path string) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("marshalling spec: %w", err)
	}

	header := "# Code generated by headerspec. DO NOT EDIT.\n\n"
	content := header + string(data)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}
