package idiomgen

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/AndroidGoLab/ndk/tools/pkg/capigen"
	"github.com/AndroidGoLab/ndk/tools/pkg/overlaymodel"
	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
)

// capiExportName returns the exported Go name that capigen generates for a C symbol.
func capiExportName(cName string) string {
	return capigen.ExportName(cName)
}

// capiExportType applies capiExportName to the base type name of a spec type,
// handling pointer and slice prefixes. Scalar Go types are left unchanged.
// E.g. "*camera_status_t" → "*Camera_status_t", "*int32" → "*int32".
func capiExportType(specType string) string {
	if specType == "" || specType == "unsafe.Pointer" || specType == "string" || specType == "bool" {
		return specType
	}
	if isScalarGoType(specType) {
		return specType
	}
	if strings.HasPrefix(specType, "[]") {
		elem := specType[2:]
		if isScalarGoType(elem) {
			return specType
		}
		return "[]" + capiExportName(elem)
	}
	stars := ""
	base := specType
	for strings.HasPrefix(base, "*") {
		stars += "*"
		base = base[1:]
	}
	if base == "" || base == "unsafe.Pointer" || base == "string" || base == "bool" || isScalarGoType(base) {
		return specType
	}
	return stars + capiExportName(base)
}

// exportName capitalizes the first letter of a string to make it an exported Go name.
func exportName(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// isScalarGoType returns true for Go scalar types.
func isScalarGoType(t string) bool {
	switch t {
	case "int8", "uint8", "int16", "uint16", "int32", "uint32",
		"int64", "uint64", "float32", "float64", "bool", "int", "uint":
		return true
	}
	return false
}

// inferEnumType checks whether an int32 param's name matches a value enum
// GoName in this module. If so, returns the enum GoName; otherwise returns
// goType unchanged. This replaces bare int32 with typed enums like sensor.Type.
func inferEnumType(
	specParamName string,
	goType string,
	enumGoNames map[string]bool,
) string {
	if goType != "int32" || specParamName == "" {
		return goType
	}

	candidate := strings.ToUpper(specParamName[:1]) + specParamName[1:]
	candidate = fixGoAcronyms(candidate)
	if enumGoNames[candidate] {
		return candidate
	}

	return goType
}

// autoGoTypeName derives an idiomatic Go type name from a C/spec type name.
// It strips the common Android NDK "A" prefix when followed by an uppercase
// letter: "AImageReader" → "ImageReader", "ALooper" → "Looper".
// Names without that pattern are kept as-is: "GLenum" stays "GLenum".
func autoGoTypeName(specName string) string {
	// Strip leading "A" if followed by uppercase (Android NDK convention).
	if len(specName) >= 2 && specName[0] == 'A' && unicode.IsUpper(rune(specName[1])) {
		return specName[1:]
	}
	return specName
}

// isDestructorFunc returns true if funcName matches common destructor patterns
// for the given spec type name: TypeName_delete, TypeName_free,
// TypeName_destroy, TypeName_release, TypeName_close.
//
// Proven disjoint with isConstructorFunc for all type names
// (proofs/Proofs/ConstructorDestructor.lean).
// Differential tested against Lean oracle.
func isDestructorFunc(funcName, specTypeName string) bool {
	for _, suffix := range []string{"_delete", "_free", "_destroy", "_release", "_close"} {
		if funcName == specTypeName+suffix {
			return true
		}
	}
	return false
}

// destructorPriority returns a priority value for destructor suffix ordering.
// Higher values indicate preferred destructor names.
func destructorPriority(funcName string) int {
	switch {
	case strings.HasSuffix(funcName, "_release"):
		return 5
	case strings.HasSuffix(funcName, "_delete"):
		return 4
	case strings.HasSuffix(funcName, "_free"):
		return 3
	case strings.HasSuffix(funcName, "_destroy"):
		return 2
	case strings.HasSuffix(funcName, "_close"):
		return 1
	default:
		return 0
	}
}

// isConstructorFunc returns true if funcName matches common constructor patterns
// for the given spec type name: TypeName_create, TypeName_new.
// Also matches ModuleName_createTypeName patterns (e.g. AAudio_createStreamBuilder).
//
// Proven disjoint with isDestructorFunc for all type names
// (proofs/Proofs/ConstructorDestructor.lean).
// Differential tested against Lean oracle.
func isConstructorFunc(funcName, specTypeName string) bool {
	for _, suffix := range []string{"_create", "_new"} {
		if funcName == specTypeName+suffix {
			return true
		}
	}
	return false
}

// autoDetectReceiver tries to find the best opaque type receiver for a function
// based on the function name prefix and first parameter type.
// Returns the Go name of the receiver type, or "" if no receiver can be detected.
func autoDetectReceiver(
	funcName string,
	fd specmodel.FuncDef,
	opaqueSpecNames map[string]bool,
	specToGoName map[string]string,
) string {
	// Strategy 1: Match function name prefix against known opaque type names.
	// E.g. "ASensor_getName" matches "ASensor" if ASensor is an opaque type.
	// Try longest match first to handle types like "ASensorEventQueue" vs "ASensor".
	type candidate struct {
		specName string
		goName   string
	}
	var candidates []candidate
	for specName := range opaqueSpecNames {
		prefix := specName + "_"
		if strings.HasPrefix(funcName, prefix) {
			candidates = append(candidates, candidate{specName, specToGoName[specName]})
		}
	}
	// Sort by spec name length descending (longest prefix first).
	sort.Slice(candidates, func(i, j int) bool {
		return len(candidates[i].specName) > len(candidates[j].specName)
	})

	if len(candidates) > 0 {
		// Verify the first param is a pointer to this type.
		best := candidates[0]
		if len(fd.Params) > 0 {
			firstParamBase := strings.TrimPrefix(fd.Params[0].Type, "*")
			if firstParamBase == best.specName && fd.Params[0].Direction != "out" {
				return best.goName
			}
		}
	}

	// Strategy 2: If the first param is a pointer to a known opaque type,
	// use that as receiver.
	if len(fd.Params) > 0 && strings.HasPrefix(fd.Params[0].Type, "*") && fd.Params[0].Direction != "out" {
		firstParamBase := strings.TrimPrefix(fd.Params[0].Type, "*")
		if opaqueSpecNames[firstParamBase] {
			return specToGoName[firstParamBase]
		}
	}

	return ""
}

// autoFuncGoName derives a Go function/method name from a C function name.
// For methods, it strips the receiver's spec type prefix and capitalizes.
// For free functions, it strips a common module prefix if present.
func autoFuncGoName(funcName, receiverSpecName string) string {
	if receiverSpecName != "" {
		prefix := receiverSpecName + "_"
		if strings.HasPrefix(funcName, prefix) {
			suffix := funcName[len(prefix):]
			if len(suffix) > 0 {
				return fixGoAcronyms(strings.ToUpper(suffix[:1]) + suffix[1:])
			}
		}
	}
	// For free functions or unmatched methods, capitalize the function name.
	return fixGoAcronyms(strings.ToUpper(funcName[:1]) + funcName[1:])
}

// isAutoGeneratable returns true if a function can be safely auto-generated
// without overlay annotations. Functions with out-params, non-scalar return
// types (except known opaque types and void), or parameters with types that
// can't be auto-resolved are not safe to auto-generate.
func isAutoGeneratable(fd specmodel.FuncDef, opaqueSpecNames map[string]bool) bool {
	// Functions with out-params need output_params overlay annotations.
	if hasOutParams(fd) {
		return false
	}

	// Functions returning pointers to unknown types (not opaque, not scalar)
	// can't be auto-generated because the template can't wrap them correctly.
	if strings.HasPrefix(fd.Returns, "*") {
		base := strings.TrimPrefix(fd.Returns, "*")
		if !opaqueSpecNames[base] && !isScalarGoType(base) {
			return false
		}
	}

	// Functions returning non-scalar, non-void, non-pointer types (like structs)
	// can't be auto-generated.
	ret := fd.Returns
	if ret != "" && ret != "void" && ret != "string" && ret != "bool" &&
		ret != "unsafe.Pointer" && !isScalarGoType(ret) &&
		!strings.HasPrefix(ret, "*") {
		return false
	}

	// Functions with callback-type parameters (non-pointer, non-scalar types)
	// can't be auto-generated without overlay annotations.
	for _, p := range fd.Params {
		base := p.Type
		for strings.HasPrefix(base, "*") {
			base = base[1:]
		}
		if base == "unsafe.Pointer" || base == "string" || base == "bool" ||
			base == "byte" || isScalarGoType(base) || base == "" {
			continue
		}
		// Opaque pointer params are handled (as handles or receivers).
		if strings.HasPrefix(p.Type, "*") && opaqueSpecNames[base] {
			continue
		}
		// Non-pointer, non-scalar type — likely a callback func type or
		// struct value that the template can't handle.
		if !strings.HasPrefix(p.Type, "*") {
			return false
		}
	}

	return true
}

// hasOutParams returns true if any parameter has direction "out".
func hasOutParams(fd specmodel.FuncDef) bool {
	for _, p := range fd.Params {
		if p.Direction == "out" {
			return true
		}
	}
	return false
}

// isAutoDetectedPure returns true if a function's return type indicates it
// should be treated as a "pure" getter (return value directly) rather than
// an error-returning function. This is used for auto-generated methods
// where no overlay specifies the pattern.
//
// Pure: string, bool, float*, uint64, int64, unsafe.Pointer, *scalar.
// Not pure: int32 (might be an error code), error types, void,
// pointers to opaque types (need struct wrapping via ReturnsNew).
func isAutoDetectedPure(returns string, errorTypes map[string]bool) bool {
	switch {
	case returns == "" || returns == "void":
		return false
	case errorTypes[returns]:
		return false
	case returns == "string" || returns == "bool":
		return true
	case returns == "float32" || returns == "float64":
		return true
	case returns == "int64" || returns == "uint64":
		return true
	case returns == "unsafe.Pointer":
		return true
	case strings.HasPrefix(returns, "*"):
		// Pointers to scalar types (e.g., *uint8) are safe for pure.
		base := strings.TrimPrefix(returns, "*")
		return isScalarGoType(base)
	default:
		// int32, uint32 — might be error codes, treat as error-returning.
		return false
	}
}

// Merge combines a specmodel.Spec and overlaymodel.Overlay into a fully
// resolved MergedSpec ready for template rendering.
func Merge(spec specmodel.Spec, overlay overlaymodel.Overlay) MergedSpec {
	m := MergedSpec{
		PackageName:   overlay.Package.GoName,
		PackageImport: overlay.Package.GoImport,
		PackageDoc:    overlay.Package.Doc,
		SourcePackage: spec.SourcePackage,
		OpaqueTypes:   make(map[string]MergedOpaqueType),
		APILevels:     overlay.APILevels,
	}
	if m.APILevels == nil {
		m.APILevels = make(map[string]int)
	}

	// Build error type set for determining destructor return types.
	errorTypes := make(map[string]bool)
	for eName, tov := range overlay.Types {
		if tov.GoError {
			errorTypes[eName] = true
		}
	}

	// typeMap resolves spec type names to Go type names.
	typeMap := make(map[string]string)

	// Build set of value enum GoNames for param-name-based enum inference.
	// Constructed early so constructor params can also benefit.
	// Includes both spec enums and overlay extra_enums (which may define
	// enum types that don't exist in the generated spec).
	enumGoNames := make(map[string]bool)
	allEnumNames := make(map[string]bool)
	for name := range spec.Enums {
		allEnumNames[name] = true
	}
	for name := range overlay.ExtraEnums {
		allEnumNames[name] = true
	}
	for enumName := range allEnumNames {
		tov := overlay.Types[enumName]
		if tov.GoError {
			continue
		}
		goName := tov.GoName
		if goName == "" {
			goName = autoGoTypeName(enumName)
		}
		enumGoNames[goName] = true
	}

	// Build set of types that are callback structs or struct accessors (not opaque handles).
	bridgeStructTypes := make(map[string]bool)
	for specName := range overlay.CallbackStructs {
		bridgeStructTypes[specName] = true
	}
	for specName := range overlay.StructAccessors {
		bridgeStructTypes[specName] = true
	}

	// specToGoName maps spec type names to Go type names for all opaque types.
	// Built during opaque type processing, used for auto-receiver detection.
	specToGoName := make(map[string]string)

	// Build opaque spec name set for handle detection.
	// Auto-include ALL opaque_ptr types (not just those with overlays).
	// Exclude transparent types — they are type aliases, not struct wrappers.
	opaqueSpecNames := make(map[string]bool)
	for specName, td := range spec.Types {
		if td.Kind != "opaque_ptr" {
			continue
		}
		if bridgeStructTypes[specName] {
			continue
		}
		tov := overlay.Types[specName]
		if tov.Transparent {
			continue
		}
		opaqueSpecNames[specName] = true
	}

	// Value struct types also need to be in opaqueSpecNames so that
	// capiArg calls .cptr() on value struct params in merged functions.
	// Build a separate set for value struct detection in output params.
	valueStructSpecNames := make(map[string]bool)
	for specName, tov := range overlay.Types {
		if tov.ValueStruct {
			opaqueSpecNames[specName] = true
			valueStructSpecNames[specName] = true
		}
	}

	// Auto-detect destructors and constructors from spec function names.
	// These are used when the overlay doesn't specify them explicitly.
	// Functions with explicit overlay entries are NOT auto-claimed, because
	// the overlay takes precedence over auto-detection.
	// Functions above BaseAPILevel are also skipped, because the generated
	// constructor/destructor would go into the base file without build tags.
	autoDestructors := make(map[string]string)  // spec type name → destructor func name
	autoConstructors := make(map[string]string) // spec type name → constructor func name
	for funcName, fd := range spec.Functions {
		// Skip functions that have explicit overlay entries — the overlay
		// author decided what this function should be.
		if _, hasOverlay := overlay.Functions[funcName]; hasOverlay {
			continue
		}
		// Skip higher-API-level functions — their constructor/destructor code
		// would be generated without build tags, causing compilation failures.
		if overlay.APILevels != nil && overlay.APILevels[funcName] > capigen.BaseAPILevel {
			continue
		}
		for specName := range opaqueSpecNames {
			if isDestructorFunc(funcName, specName) {
				// Prefer _release > _delete > _free > _destroy > _close for
				// auto-detection. Higher-priority suffixes overwrite lower ones.
				existing := autoDestructors[specName]
				if existing == "" || destructorPriority(funcName) > destructorPriority(existing) {
					autoDestructors[specName] = funcName
				}
			}
			if isConstructorFunc(funcName, specName) {
				if fd.Returns == "" || fd.Returns == "void" {
					// Constructor that doesn't return anything — check for out-param.
					hasOutParam := false
					for _, p := range fd.Params {
						if p.Direction == "out" {
							hasOutParam = true
							break
						}
					}
					if !hasOutParam {
						continue
					}
				}
				autoConstructors[specName] = funcName
			}
		}
	}

	// Build opaque types and populate typeMap for pointer types.
	// Process ALL opaque_ptr types, not just those with overlay entries.
	for _, specName := range sortedKeys(spec.Types) {
		td := spec.Types[specName]
		if td.Kind != "opaque_ptr" {
			continue
		}
		if bridgeStructTypes[specName] {
			continue
		}
		tov, hasOverlay := overlay.Types[specName]

		// Determine Go name: overlay takes precedence, then auto-naming.
		goName := ""
		switch {
		case hasOverlay && tov.GoName != "":
			goName = tov.GoName
		case hasOverlay:
			goName = specName
		default:
			goName = autoGoTypeName(specName)
		}

		specToGoName[specName] = goName

		// Transparent types become simple type aliases (no struct wrapper).
		// Used for void* typedefs like EGLDisplay that are value types in Go.
		if hasOverlay && tov.Transparent {
			m.TypeAliases = append(m.TypeAliases, MergedTypeAlias{
				GoName:   goName,
				CapiType: capiExportName(specName),
			})
			typeMap[specName] = goName
			typeMap["*"+specName] = "*" + goName
			continue
		}

		// Value struct types are handled separately — they get exported Go
		// fields instead of an opaque pointer wrapper.
		if hasOverlay && tov.ValueStruct {
			typeMap[specName] = goName
			typeMap["*"+specName] = "*" + goName
			continue
		}

		// Resolve destructor: overlay takes precedence, then auto-detected.
		destructor := ""
		if hasOverlay && tov.Destructor != "" {
			destructor = tov.Destructor
		} else if d, ok := autoDestructors[specName]; ok {
			destructor = d
		}

		// Resolve constructor: overlay takes precedence, then auto-detected.
		constructor := ""
		if hasOverlay && tov.Constructor != "" {
			constructor = tov.Constructor
		} else if c, ok := autoConstructors[specName]; ok {
			constructor = c
		}

		destructorReturnsError := false
		if destructor != "" {
			if fd, ok := spec.Functions[destructor]; ok {
				destructorReturnsError = errorTypes[fd.Returns]
			}
		}
		constructorReturnsPointer := false
		var constructorParams []MergedParam
		if constructor != "" {
			if fd, ok := spec.Functions[constructor]; ok {
				constructorReturnsPointer = strings.HasPrefix(fd.Returns, "*")
				ctorParamIdx := 0
				for _, p := range fd.Params {
					if p.Direction == "out" {
						continue
					}
					goType := resolveType(p.Type, typeMap)
					goType = inferEnumType(p.Name, goType, enumGoNames)
					remapped := goType != resolveType(p.Type, typeMap)
					isString := p.Type == "*byte"
					if isString {
						goType = "string"
					}
					baseType := strings.TrimPrefix(p.Type, "*")
					isHandle := strings.HasPrefix(p.Type, "*") && opaqueSpecNames[baseType]
					if isHandle {
						goType = "unsafe.Pointer"
					}
					capiType := capiExportType(p.Type)
					if isFixedArrayGoType(goType) {
						goType = "*" + goType
						capiType = "*" + capiType
					}
					constructorParams = append(constructorParams, MergedParam{
						Name:     safeGoParamName(p.Name, ctorParamIdx),
						GoType:   goType,
						CapiType: capiType,
						Remapped: remapped,
						IsString: isString,
						IsHandle: isHandle,
					})
					ctorParamIdx++
				}
			}
		}
		exportedConstructor := ""
		if constructor != "" {
			exportedConstructor = capiExportName(constructor)
		}
		exportedDestructor := ""
		if destructor != "" {
			exportedDestructor = capiExportName(destructor)
		}

		pattern := ""
		var interfaces []string
		if hasOverlay {
			pattern = tov.Pattern
			interfaces = tov.Interfaces
		}

		m.OpaqueTypes[goName] = MergedOpaqueType{
			GoName:                    goName,
			CapiType:                  capiExportName(specName),
			Constructor:               exportedConstructor,
			ConstructorReturnsPointer: constructorReturnsPointer,
			ConstructorParams:         constructorParams,
			Destructor:                exportedDestructor,
			DestructorReturnsError:    destructorReturnsError,
			Pattern:                   pattern,
			Interfaces:                interfaces,
		}
		typeMap["*"+specName] = "*" + goName
	}

	// Pre-processing: link overlay enum sources to spec enums.
	// When an overlay type has enum_source, copy the spec enum values
	// to spec.Enums[overlayKey] so the existing merge logic finds them.
	// This handles NDK's pattern of separate `typedef int32_t foo_t;`
	// and `enum { FOO_X, FOO_Y };` where the typedef and enum have
	// different names.
	if spec.Enums == nil {
		spec.Enums = make(map[string][]specmodel.EnumValue)
	}
	for typeName, tov := range overlay.Types {
		if tov.EnumSource == "" {
			continue
		}
		if vals, ok := spec.Enums[tov.EnumSource]; ok {
			spec.Enums[typeName] = vals
		}
	}

	// Merge extra enum values from the overlay into the spec enums.
	// When the spec has macros, auto-resolve values from headers instead
	// of relying on manually specified values in the overlay.
	// This must happen before the type alias loop so that extra_enums-only
	// types (like MetadataTag) are recognized as enums, not type aliases.
	consumedMacros := make(map[string]bool)
	for enumName, extras := range overlay.ExtraEnums {
		for _, ev := range extras {
			value := ev.Value
			if macroVal, ok := spec.Macros[ev.Name]; ok {
				value = macroVal
				consumedMacros[ev.Name] = true
			}
			spec.Enums[enumName] = append(spec.Enums[enumName], specmodel.EnumValue{
				Name:  ev.Name,
				Value: value,
			})
		}
	}

	// Also consume macros that already appear in spec enums (from real C enums).
	for _, vals := range spec.Enums {
		for _, v := range vals {
			consumedMacros[v.Name] = true
		}
	}

	// Store unconsumed macros for generation as untyped constants.
	m.UntypedMacros = make(map[string]int64)
	for name, value := range spec.Macros {
		if !consumedMacros[name] {
			m.UntypedMacros[name] = value
		}
	}

	// Build type aliases for typedef and pointer_handle types.
	for _, specName := range sortedKeys(spec.Types) {
		td := spec.Types[specName]
		if td.Kind == "opaque_ptr" {
			continue // Already handled above as OpaqueTypes.
		}
		if !isTypedefKind(td.Kind) {
			continue
		}
		// Skip function pointer typedefs and other non-Go types
		// (capigen skips these too).
		baseType := strings.TrimPrefix(td.Kind, "typedef_")
		if strings.Contains(baseType, ":") {
			continue
		}
		// Skip if this type is classified as an enum (has spec enum values and
		// either has an overlay entry or will be auto-generated as a value enum).
		if _, isEnum := spec.Enums[specName]; isEnum {
			continue
		}
		goName := specName
		if tov, ok := overlay.Types[specName]; ok && tov.GoName != "" {
			goName = tov.GoName
		}
		m.TypeAliases = append(m.TypeAliases, MergedTypeAlias{
			GoName:   goName,
			CapiType: capiExportName(specName),
		})
		typeMap[specName] = goName
		typeMap["*"+specName] = "*" + goName
	}

	// Classify enums: error vs. value.
	// Process ALL enums in the spec, not just those with overlay entries.
	// Overlay entries customize naming; auto-generation uses defaults.
	enumNames := sortedKeys(spec.Enums)
	for _, enumName := range enumNames {
		tov, hasOverlay := overlay.Types[enumName]
		values := spec.Enums[enumName]

		if hasOverlay && tov.GoError {
			m.ErrorEnums = append(m.ErrorEnums, mergeErrorEnum(enumName, tov, values))
			typeMap[enumName] = "Error"
			typeMap["*"+enumName] = "*Error"
		} else {
			goName := ""
			switch {
			case hasOverlay && tov.GoName != "":
				goName = tov.GoName
			case hasOverlay:
				goName = enumName
			default:
				goName = autoGoTypeName(enumName)
			}
			// Use overlay settings if available, otherwise empty TypeOverlay for defaults.
			enumOverlay := tov
			baseType := kindToBaseType(spec.Types[enumName].Kind)
			m.ValueEnums = append(m.ValueEnums, mergeValueEnum(enumName, goName, baseType, enumOverlay, values))
			typeMap[enumName] = goName
			typeMap["*"+enumName] = "*" + goName
		}
	}

	// Build value structs from overlay annotations + spec struct definitions.
	for _, specName := range sortedKeys(overlay.Types) {
		tov := overlay.Types[specName]
		if !tov.ValueStruct {
			continue
		}
		sd, ok := spec.Structs[specName]
		if !ok {
			continue
		}
		goName := tov.GoName
		if goName == "" {
			goName = specName
		}
		capiType := capiExportName(specName)
		var fields []MergedValueStructField
		for _, sf := range sd.Fields {
			fields = append(fields, MergedValueStructField{
				GoName:   exportName(sf.Name),
				CName:    sf.Name,
				GoType:   sf.Type,
				CapiType: sf.Type,
			})
		}
		m.ValueStructs = append(m.ValueStructs, MergedValueStruct{
			GoName:   goName,
			CapiType: capiType,
			Fields:   fields,
		})
	}

	// Build callback annotation lookup by cross-referencing function overlays
	// with spec function parameters to find which callback type each overlay targets.
	cbAnnotations := make(map[string]callbackAnnotation)
	for fname, fov := range overlay.Functions {
		if fov.CallbackParam == "" || fov.GoCallbackType == "" {
			continue
		}
		if fd, ok := spec.Functions[fname]; ok {
			for _, p := range fd.Params {
				if p.Name == fov.CallbackParam {
					cbAnnotations[p.Type] = callbackAnnotation{
						goType: fov.GoCallbackType,
						goSig:  fov.GoCallbackSig,
					}
					typeMap[p.Type] = fov.GoCallbackType
					break
				}
			}
		}
	}

	// Build receiver capi type lookup for GoName auto-derivation.
	receiverCapiTypes := make(map[string]string) // Go name → capi type
	for _, ot := range m.OpaqueTypes {
		receiverCapiTypes[ot.GoName] = ot.CapiType
	}

	// Auto-generate type aliases for types used in function params/returns
	// that aren't in the typeMap yet. This covers cross-module types (off_t,
	// ANativeWindow, etc.) and callback types without explicit overlay entries.
	autoAliasTypes := collectUnresolvedFuncTypes(spec, overlay, typeMap, opaqueSpecNames, specToGoName)
	for _, typeName := range sortedKeys(autoAliasTypes) {
		// If this cross-module type has an overlay entry with a go_name,
		// create a handle struct instead of a simple type alias.
		// This handles cases like ANativeWindow in surfacetexture where the
		// overlay wants go_name: NativeWindow to create a proper wrapper struct.
		if tov, ok := overlay.Types[typeName]; ok && tov.GoName != "" {
			goName := tov.GoName
			if tov.Transparent {
				// Transparent types become simple type aliases, not struct wrappers.
				m.TypeAliases = append(m.TypeAliases, MergedTypeAlias{
					GoName:   goName,
					CapiType: capiExportName(typeName),
				})
				typeMap[typeName] = goName
				typeMap["*"+typeName] = "*" + goName
				continue
			}
			m.OpaqueTypes[goName] = MergedOpaqueType{
				GoName:   goName,
				CapiType: capiExportName(typeName),
			}
			typeMap["*"+typeName] = "*" + goName
			opaqueSpecNames[typeName] = true
			continue
		}
		goName := capiExportName(typeName)
		m.TypeAliases = append(m.TypeAliases, MergedTypeAlias{
			GoName:   goName,
			CapiType: capiExportName(typeName),
		})
		typeMap[typeName] = goName
		typeMap["*"+typeName] = "*" + goName
	}

	// Build a set of functions that are already claimed as constructors or
	// destructors, so they are not also emitted as regular methods.
	// Also claim ALL destructor-like functions (even non-winners) to prevent
	// them from generating methods that conflict with the Close() destructor.
	claimedFuncs := make(map[string]bool)
	for specName := range opaqueSpecNames {
		tov := overlay.Types[specName]
		destructor := ""
		switch {
		case tov.Destructor != "":
			destructor = tov.Destructor
		default:
			destructor = autoDestructors[specName]
		}
		if destructor != "" {
			claimedFuncs[destructor] = true
		}
		constructor := ""
		switch {
		case tov.Constructor != "":
			constructor = tov.Constructor
		default:
			constructor = autoConstructors[specName]
		}
		if constructor != "" {
			claimedFuncs[constructor] = true
		}

		// Claim ALL destructor-like and constructor-like functions for this
		// type, not just the chosen one. This prevents auto-generating methods
		// like Close() that would conflict with the destructor's Close().
		// But do NOT claim functions that have explicit overlay entries —
		// the overlay author controls those.
		for funcName := range spec.Functions {
			if _, hasOverlay := overlay.Functions[funcName]; hasOverlay {
				continue
			}
			if isDestructorFunc(funcName, specName) || isConstructorFunc(funcName, specName) {
				claimedFuncs[funcName] = true
			}
		}
	}

	// Resolve functions → methods or free functions.
	// Process ALL spec functions, not just those with overlay entries.
	// Overlay entries customize behavior; auto-generation uses defaults.
	funcNames := sortedKeys(spec.Functions)
	for _, funcName := range funcNames {
		fd := spec.Functions[funcName]
		fov, hasOverlay := overlay.Functions[funcName]

		// If explicitly skipped in overlay, skip.
		if hasOverlay && fov.Skip {
			continue
		}

		// Skip functions claimed as constructors/destructors.
		if claimedFuncs[funcName] {
			continue
		}

		if hasOverlay {
			// Overlay-directed processing (original behavior).
			switch {
			case fov.Receiver != "":
				m.Methods = append(m.Methods, mergeMethod(funcName, fd, fov, overlay.APILevels, typeMap, enumGoNames, receiverCapiTypes, opaqueSpecNames, overlay.CallbackStructs, overlay.StructAccessors, valueStructSpecNames))
			case fov.GoName != "":
				m.FreeFunctions = append(m.FreeFunctions, mergeFreeFunction(funcName, fd, fov, overlay.APILevels, typeMap, enumGoNames, opaqueSpecNames, valueStructSpecNames))
			}
		} else {
			// Auto-generate: detect receiver and derive Go name.
			receiver := autoDetectReceiver(funcName, fd, opaqueSpecNames, specToGoName)
			switch {
			case receiver != "":
				// Skip auto-generation for functions that would produce
				// invalid code without overlay annotations.
				if !isAutoGeneratable(fd, opaqueSpecNames) {
					break
				}

				// Find the spec type name for this receiver to derive the method name.
				receiverSpecName := ""
				for sn, gn := range specToGoName {
					if gn == receiver {
						receiverSpecName = sn
						break
					}
				}
				goName := autoFuncGoName(funcName, receiverSpecName)
				autoFov := overlaymodel.FuncOverlay{
					Receiver: receiver,
					GoName:   goName,
				}
				// Auto-detect Pure for functions that return non-error types.
				autoFov.Pure = isAutoDetectedPure(fd.Returns, errorTypes)

				// For methods returning a pointer to a known opaque type
				// in this module, use ReturnsNew instead of Pure, because
				// the template needs to wrap the capi pointer in a Go struct.
				returnBase := strings.TrimPrefix(fd.Returns, "*")
				if strings.HasPrefix(fd.Returns, "*") && opaqueSpecNames[returnBase] {
					autoFov.Pure = false
					autoFov.ReturnsNew = specToGoName[returnBase]
				}

				m.Methods = append(m.Methods, mergeMethod(funcName, fd, autoFov, overlay.APILevels, typeMap, enumGoNames, receiverCapiTypes, opaqueSpecNames, overlay.CallbackStructs, overlay.StructAccessors, valueStructSpecNames))

			default:
				// Skip auto-generation for free functions that would produce
				// invalid code without overlay annotations.
				if !isAutoGeneratable(fd, opaqueSpecNames) {
					break
				}

				// Package-level free function.
				goName := autoFuncGoName(funcName, "")
				autoFov := overlaymodel.FuncOverlay{
					GoName: goName,
				}
				m.FreeFunctions = append(m.FreeFunctions, mergeFreeFunction(funcName, fd, autoFov, overlay.APILevels, typeMap, enumGoNames, opaqueSpecNames, valueStructSpecNames))
			}
		}
	}

	// Resolve bridge functions defined only in the overlay (with bridge_params).
	for _, funcName := range sortedKeys(overlay.Functions) {
		if _, inSpec := spec.Functions[funcName]; inSpec {
			continue
		}
		fov := overlay.Functions[funcName]
		if fov.Skip || len(fov.BridgeParams) == 0 {
			continue
		}
		var params []specmodel.Param
		for _, bp := range fov.BridgeParams {
			params = append(params, specmodel.Param{Name: bp.Name, Type: bp.Type})
		}
		fd := specmodel.FuncDef{
			CName:   funcName,
			Params:  params,
			Returns: fov.BridgeReturns,
		}
		switch {
		case fov.Receiver != "":
			m.Methods = append(m.Methods, mergeMethod(funcName, fd, fov, overlay.APILevels, typeMap, enumGoNames, receiverCapiTypes, opaqueSpecNames, overlay.CallbackStructs, overlay.StructAccessors, valueStructSpecNames))
		case fov.GoName != "":
			m.FreeFunctions = append(m.FreeFunctions, mergeFreeFunction(funcName, fd, fov, overlay.APILevels, typeMap, enumGoNames, opaqueSpecNames, valueStructSpecNames))
		}
	}

	// Resolve callbacks.
	cbNames := sortedKeys(spec.Callbacks)
	for _, cbName := range cbNames {
		cbd := spec.Callbacks[cbName]
		m.Callbacks = append(m.Callbacks, mergeCallback(cbName, cbd, cbAnnotations))
	}

	// Generate type aliases for callback structs so idiomatic packages re-export them.
	for _, specName := range sortedKeys(overlay.CallbackStructs) {
		csov := overlay.CallbackStructs[specName]
		m.CallbackStructAliases = append(m.CallbackStructAliases, MergedTypeAlias{
			GoName:   csov.GoName,
			CapiType: csov.GoName, // bridge generates Go struct with GoName in the capi package
		})
	}

	// Merge callback struct overlays with spec struct definitions.
	for _, specName := range sortedKeys(overlay.CallbackStructs) {
		csov := overlay.CallbackStructs[specName]
		sd, hasSpec := spec.Structs[specName]
		mcs := MergedCallbackStruct{
			SpecName:     specName,
			GoName:       csov.GoName,
			ContextField: csov.ContextField,
		}

		if hasSpec {
			for _, field := range sd.Fields {
				if field.Name == csov.ContextField {
					continue
				}
				fov, hasFieldOverlay := csov.Fields[field.Name]
				if !hasFieldOverlay {
					continue
				}
				mcs.Fields = append(mcs.Fields, MergedCallbackField{
					CName:        field.Name,
					GoName:       fov.GoName,
					GoSignature:  fov.GoSignature,
					GoParamCount: countGoParams(fov.GoSignature),
					Params:       field.Params,
				})
			}
		} else {
			for _, cName := range sortedKeys(csov.Fields) {
				fov := csov.Fields[cName]
				goParamCount := fov.GoParamCount
				if goParamCount == 0 {
					goParamCount = countGoParams(fov.GoSignature)
				}
				// Look up callback type params from spec.
				var params []specmodel.Param
				if fov.CallbackType != "" {
					if cbDef, ok := spec.Callbacks[fov.CallbackType]; ok {
						params = cbDef.Params
					}
				}
				mcs.Fields = append(mcs.Fields, MergedCallbackField{
					CName:        cName,
					GoName:       fov.GoName,
					GoSignature:  fov.GoSignature,
					GoParamCount: goParamCount,
					Params:       params,
				})
			}
		}
		m.CallbackStructs = append(m.CallbackStructs, mcs)
	}

	// Merge struct accessor overlays.
	for _, specName := range sortedKeys(overlay.StructAccessors) {
		saov := overlay.StructAccessors[specName]
		m.StructAccessors = append(m.StructAccessors, MergedStructAccessor{
			SpecName:   capiExportName(specName),
			CountField: saov.CountField,
			ItemsField: saov.ItemsField,
			ItemType:   saov.ItemType,
			DeleteFunc: capiExportName(saov.DeleteFunc),
		})
	}

	// Merge lifecycle overlay.
	if overlay.Lifecycle != nil {
		ml := &MergedLifecycle{
			EntryPoint:        overlay.Lifecycle.EntryPoint,
			ActivityType:      overlay.Lifecycle.ActivityType,
			CallbacksAccessor: overlay.Lifecycle.CallbacksAccessor,
			CallbackStruct:    overlay.Lifecycle.CallbackStruct,
			EntryCallback:     overlay.Lifecycle.EntryCallback,
		}

		// Resolve Go activity type name.
		if tov, ok := overlay.Types[overlay.Lifecycle.ActivityType]; ok && tov.GoName != "" {
			ml.GoActivityType = tov.GoName
		} else {
			ml.GoActivityType = overlay.Lifecycle.ActivityType
		}

		// Extract lifecycle fields from spec struct definition.
		if sd, ok := spec.Structs[overlay.Lifecycle.CallbackStruct]; ok {
			for _, field := range sd.Fields {
				if field.Type != "func_ptr" {
					continue
				}
				goName := fixGoAcronyms(strings.ToUpper(field.Name[:1]) + field.Name[1:])
				ml.Fields = append(ml.Fields, MergedCallbackField{
					CName:   field.Name,
					GoName:  goName,
					Params:  field.Params,
					Returns: field.Returns,
				})
			}
		}

		m.Lifecycle = ml
	}

	m.ExtraBridgeC = overlay.ExtraBridgeC
	m.ExtraBridgeGo = overlay.ExtraBridgeGo

	return m
}

func mergeErrorEnum(specName string, tov overlaymodel.TypeOverlay, values []specmodel.EnumValue) MergedErrorEnum {
	ee := MergedErrorEnum{
		GoName:       specName,
		Prefix:       tov.ErrorPrefix,
		SuccessValue: tov.SuccessValue,
	}
	for _, v := range values {
		goName := v.Name
		if tov.StripPrefix != "" && v.Name != tov.SuccessValue {
			goName = "Err" + stripAndTitle(v.Name, tov.StripPrefix)
		}
		ee.Values = append(ee.Values, MergedEnumValue{
			GoName:   goName,
			SpecName: v.Name,
			Value:    v.Value,
		})
	}
	return ee
}

func mergeValueEnum(specName, goName, baseType string, tov overlaymodel.TypeOverlay, values []specmodel.EnumValue) MergedValueEnum {
	if baseType == "" {
		baseType = "int32"
	}

	// If any enum value exceeds the int32 range, promote to uint64.
	if baseType == "int32" {
		for _, v := range values {
			if v.Value > math.MaxInt32 || v.Value < math.MinInt32 {
				baseType = "uint64"
				break
			}
		}
	}

	ve := MergedValueEnum{
		GoName:       goName,
		SpecName:     specName,
		BaseType:     baseType,
		StripPrefix:  tov.StripPrefix,
		StringMethod: tov.StringMethod,
	}
	isUnsigned := baseType == "uint32" || baseType == "uint64"
	seen := make(map[string]bool)
	for _, v := range values {
		// Skip _MAX_ENUM sentinel values (present in all Vulkan enums).
		if strings.HasSuffix(v.Name, "_MAX_ENUM") {
			continue
		}
		stripped := stripAndTitle(v.Name, tov.StripPrefix)
		// Skip duplicate Go names (e.g., Vulkan aliases that produce the same name).
		if seen[stripped] {
			continue
		}
		seen[stripped] = true

		// For unsigned base types, render negative int64 values as hex
		// so that Go interprets them as their unsigned equivalent.
		var valueStr string
		if isUnsigned && v.Value < 0 {
			valueStr = fmt.Sprintf("0x%X", uint64(v.Value))
		}

		ve.Values = append(ve.Values, MergedEnumValue{
			GoName:   stripped,
			SpecName: v.Name,
			Value:    v.Value,
			ValueStr: valueStr,
		})
	}
	return ve
}

// kindToBaseType converts a spec type kind to a Go base type string.
func kindToBaseType(kind string) string {
	switch kind {
	case "typedef_uint32":
		return "uint32"
	case "typedef_int64":
		return "int64"
	case "typedef_uint64":
		return "uint64"
	default:
		return "int32"
	}
}

// stripAndTitle removes the given prefix from name and converts the
// remainder to TitleCase. For example: "AAUDIO_DIRECTION_OUTPUT" with
// prefix "AAUDIO_DIRECTION_" becomes "Output".
func stripAndTitle(name, prefix string) string {
	if prefix == "" {
		return name
	}
	s := strings.TrimPrefix(name, prefix)
	if s == "" {
		return name
	}
	return toTitleCase(s)
}

// toTitleCase converts an UPPER_SNAKE_CASE string to TitleCase,
// normalizing common acronyms (e.g. "DEVICE_ID" → "DeviceID", not "DeviceId").
func toTitleCase(s string) string {
	parts := strings.Split(strings.ToLower(s), "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return fixGoAcronyms(b.String())
}

func mergeMethod(funcName string, fd specmodel.FuncDef, fov overlaymodel.FuncOverlay, apiLevels map[string]int, typeMap map[string]string, enumGoNames map[string]bool, receiverCapiTypes map[string]string, opaqueSpecNames map[string]bool, callbackStructOverlays map[string]overlaymodel.CallbackStructOverlay, structAccessors map[string]overlaymodel.StructAccessorOverlay, valueStructSpecNames map[string]bool) MergedMethod {
	var params []MergedParam
	receiverFound := false
	paramIdx := 0
	for _, p := range fd.Params {
		if !receiverFound && isReceiverParam(p) {
			receiverFound = true
			continue
		}
		goType := resolveType(p.Type, typeMap)
		// Overlay explicit override takes priority over heuristic inference.
		if override, ok := fov.ParamTypes[p.Name]; ok {
			goType = override
		} else {
			goType = inferEnumType(p.Name, goType, enumGoNames)
		}
		baseType := strings.TrimPrefix(p.Type, "*")
		isHandle := strings.HasPrefix(p.Type, "*") && opaqueSpecNames[baseType]
		isString := p.Type == "*byte"
		remapped := isTypeRemapped(p.Type, typeMap) || goType != resolveType(p.Type, typeMap)
		if isString {
			goType = "string"
		}
		capiType := capiExportType(p.Type)
		if isFixedArrayGoType(goType) {
			goType = "*" + goType
			capiType = "*" + capiType
		}
		params = append(params, MergedParam{
			Name:      safeGoParamName(p.Name, paramIdx),
			GoType:    goType,
			CapiType:  capiType,
			IsHandle:  isHandle,
			IsString:  isString,
			Remapped:  remapped,
			Direction: p.Direction,
		})
		paramIdx++
	}

	// If BufGoType is set, override the GoType of the buffer param.
	if fov.BufGoType != "" && fov.BufParam != "" {
		for i := range params {
			if params[i].Name == safeGoName(fov.BufParam) {
				params[i].GoType = fov.BufGoType
				break
			}
		}
	}

	// If TimeoutParam is set, override the GoType to time.Duration.
	if fov.TimeoutParam != "" {
		unit := fov.TimeoutUnit
		if unit == "" {
			unit = "ns"
		}
		for i := range params {
			if params[i].Name == safeGoName(fov.TimeoutParam) {
				params[i].GoType = "time.Duration"
				params[i].DurationUnit = unit
				params[i].Name = durationParamName(fov.TimeoutParam)
				break
			}
		}
	}

	// Auto-derive GoName if not specified in overlay.
	goName := fov.GoName
	if goName == "" {
		goName = deriveGoName(funcName, fov.Receiver, receiverCapiTypes)
	}

	// Resolve return type through typeMap.
	returns := resolveType(fd.Returns, typeMap)

	returnsNewDirect := false
	if fov.ReturnsNew != "" {
		returnsNewDirect = strings.HasPrefix(fd.Returns, "*")
	}
	mm := MergedMethod{
		GoName:           goName,
		CName:            capiExportName(funcName),
		ReceiverType:     fov.Receiver,
		Params:           params,
		Returns:          returns,
		Chain:            fov.Chain,
		Pure:             fov.Pure,
		ReturnsNew:       fov.ReturnsNew,
		ReturnsNewDirect: returnsNewDirect,
		ReturnsFrames:    fov.ReturnsFrames,
		BufGoType:        fov.BufGoType,
		TimeoutUnit:      fov.TimeoutUnit,
		FixedParams:      fov.FixedParams,
		CallbackParam:    fov.CallbackParam,
	}

	// If CallbackParam is set, find the matching callback struct and its Go name.
	if fov.CallbackParam != "" {
		for _, p := range fd.Params {
			if p.Name == fov.CallbackParam {
				specName := strings.TrimPrefix(p.Type, "*")
				mm.CallbackStruct = capiExportName(specName)
				if csov, ok := callbackStructOverlays[specName]; ok {
					mm.CallbackGoType = csov.GoName
				}
				break
			}
		}
	}

	// If ReturnsListStruct is set, look up the struct accessor.
	if fov.ReturnsListStruct != "" {
		if saov, ok := structAccessors[fov.ReturnsListStruct]; ok {
			mm.ReturnsListAccessor = &MergedStructAccessor{
				SpecName:   capiExportName(fov.ReturnsListStruct),
				CountField: saov.CountField,
				ItemsField: saov.ItemsField,
				ItemType:   saov.ItemType,
				DeleteFunc: capiExportName(saov.DeleteFunc),
			}
		}
	}

	// If CustomCall is set, populate the merged custom call.
	if fov.CustomCall != nil {
		mcc := &MergedCustomCall{
			Args: fov.CustomCall.Args,
		}
		for _, p := range fov.CustomCall.Params {
			mcc.Params = append(mcc.Params, MergedCustomCallParam{
				Name:   p.Name,
				GoType: p.GoType,
			})
		}
		mm.CustomCall = mcc
	}

	// Resolve output_params: convert C output pointer params to Go return values.
	if len(fov.OutputParams) > 0 {
		mm.OutputParams = resolveOutputParams(fov.OutputParams, fd, opaqueSpecNames, valueStructSpecNames)
		mm.ReturnsBool = fd.Returns == "bool"
		// Mark output params with direction "out" so skipAndFilterOut removes them.
		for i := range mm.Params {
			if _, isOutput := fov.OutputParams[mm.Params[i].Name]; isOutput {
				mm.Params[i].Direction = "out"
			}
		}
	}

	if apiLevels != nil {
		mm.APILevel = apiLevels[funcName]
	}
	return mm
}

// deriveGoName auto-derives a Go method name from a C function name by
// durationParamName derives a Go parameter name from a C duration/timeout param name.
// "timeoutNanoseconds" → "timeout", "actualDurationNanos" → "actualDuration",
// "initialTargetWorkDurationNanos" → "initialTargetWorkDuration".
func durationParamName(cName string) string {
	for _, suffix := range []string{"Nanoseconds", "Nanos", "Millis", "Microseconds", "Micros"} {
		if strings.HasSuffix(cName, suffix) {
			trimmed := strings.TrimSuffix(cName, suffix)
			if trimmed == "" || trimmed == "timeout" {
				return "timeout"
			}
			return trimmed
		}
	}
	return cName
}

// stripping the receiver's capi type prefix and capitalizing the first letter.
// "AAudioStreamBuilder_setChannelCount" with receiver "StreamBuilder" → "SetChannelCount".
func deriveGoName(funcName, receiver string, receiverCapiTypes map[string]string) string {
	if capiType, ok := receiverCapiTypes[receiver]; ok {
		prefix := capiType + "_"
		if strings.HasPrefix(funcName, prefix) {
			suffix := funcName[len(prefix):]
			if len(suffix) > 0 {
				return fixGoAcronyms(strings.ToUpper(suffix[:1]) + suffix[1:])
			}
		}
	}
	return fixGoAcronyms(funcName)
}

func mergeFreeFunction(funcName string, fd specmodel.FuncDef, fov overlaymodel.FuncOverlay, apiLevels map[string]int, typeMap map[string]string, enumGoNames map[string]bool, opaqueSpecNames map[string]bool, valueStructSpecNames map[string]bool) MergedFreeFunction {
	var params []MergedParam
	for i, p := range fd.Params {
		goType := resolveType(p.Type, typeMap)
		if override, ok := fov.ParamTypes[p.Name]; ok {
			goType = override
		} else {
			goType = inferEnumType(p.Name, goType, enumGoNames)
		}
		baseType := strings.TrimPrefix(p.Type, "*")
		isHandle := strings.HasPrefix(p.Type, "*") && opaqueSpecNames[baseType]
		isString := p.Type == "*byte"
		remapped := isTypeRemapped(p.Type, typeMap) || goType != resolveType(p.Type, typeMap)
		if isString {
			goType = "string"
		}
		capiType := capiExportType(p.Type)
		if isFixedArrayGoType(goType) {
			goType = "*" + goType
			capiType = "*" + capiType
		}
		params = append(params, MergedParam{
			Name:      safeGoParamName(p.Name, i),
			GoType:    goType,
			CapiType:  capiType,
			IsHandle:  isHandle,
			IsString:  isString,
			Remapped:  remapped,
			Direction: p.Direction,
		})
	}
	// If TimeoutParam is set, override the GoType to time.Duration.
	if fov.TimeoutParam != "" {
		unit := fov.TimeoutUnit
		if unit == "" {
			unit = "ns"
		}
		for i := range params {
			if params[i].Name == safeGoName(fov.TimeoutParam) {
				params[i].GoType = "time.Duration"
				params[i].Name = durationParamName(fov.TimeoutParam)
				params[i].DurationUnit = unit
				break
			}
		}
	}

	returns := resolveType(fd.Returns, typeMap)
	isHandleReturn := strings.HasPrefix(fd.Returns, "*") && opaqueSpecNames[strings.TrimPrefix(fd.Returns, "*")]
	ff := MergedFreeFunction{
		GoName:         fov.GoName,
		CName:          capiExportName(funcName),
		Params:         params,
		Returns:        returns,
		CapiReturns:    fd.Returns,
		IsHandleReturn: isHandleReturn,
		ReturnsNew:     fov.ReturnsNew,
	}

	// Resolve output_params: convert C output pointer params to Go return values.
	if len(fov.OutputParams) > 0 {
		ff.OutputParams = resolveOutputParams(fov.OutputParams, fd, opaqueSpecNames, valueStructSpecNames)
		ff.ReturnsBool = fd.Returns == "bool"
		// Mark output params with direction "out" so they are filtered from visible params.
		for i := range ff.Params {
			if _, isOutput := fov.OutputParams[ff.Params[i].Name]; isOutput {
				ff.Params[i].Direction = "out"
			}
		}
	}

	if apiLevels != nil {
		ff.APILevel = apiLevels[funcName]
	}
	return ff
}

// resolveOutputParams builds MergedOutputParam entries from the overlay's output_params map.
// It matches param names against the function's spec params to determine capi types.
func resolveOutputParams(
	outputParams map[string]string,
	fd specmodel.FuncDef,
	opaqueSpecNames map[string]bool,
	valueStructSpecNames map[string]bool,
) []MergedOutputParam {
	// Process output params in the order they appear in the function signature.
	var result []MergedOutputParam
	for _, p := range fd.Params {
		goType, isOutput := outputParams[p.Name]
		if !isOutput {
			continue
		}

		// Strip one level of pointer from the C param type to get the local var type.
		// **AImageReader -> *AImageReader (opaque handle pointer)
		// **uint8 -> *uint8 (scalar pointer)
		// *int32 -> int32 (scalar value)
		// *AHardwareBuffer_Desc -> AHardwareBuffer_Desc (value struct)
		localType := strings.TrimPrefix(p.Type, "*")

		// Determine the base type name (without pointer stars).
		specBase := localType
		for strings.HasPrefix(specBase, "*") {
			specBase = specBase[1:]
		}

		// Check if this is a value struct (before checking opaque handles,
		// since value struct names are also in opaqueSpecNames).
		isValueStruct := valueStructSpecNames[specBase]

		// Determine if the Go type is an opaque handle wrapper.
		// Value structs are NOT opaque handles even though they appear in opaqueSpecNames.
		isHandle := !isValueStruct && opaqueSpecNames[specBase]

		// Build the capi type for the local variable declaration.
		// Opaque types need capi. prefix (e.g., *capi.AImageReader).
		// Value structs need capi. prefix without pointer (e.g., capi.AHardwareBuffer_Desc).
		// Scalar types use the Go type directly (e.g., *uint8, int32).
		var capiType string
		if isHandle || isValueStruct || !isScalarGoType(specBase) {
			// Reconstruct with capi. prefix on the base type.
			stars := strings.TrimSuffix(localType, specBase)
			capiType = stars + "capi." + capiExportName(specBase)
		} else {
			capiType = localType
		}

		result = append(result, MergedOutputParam{
			CParamName:    p.Name,
			GoType:        goType,
			CapiType:      capiType,
			IsHandle:      isHandle,
			IsValueStruct: isValueStruct,
		})
	}
	return result
}

// isTypeRemapped returns true if the spec type was resolved through typeMap.
func isTypeRemapped(specType string, typeMap map[string]string) bool {
	if _, ok := typeMap[specType]; ok {
		return true
	}
	for _, prefix := range []string{"[]", "**", "*"} {
		if strings.HasPrefix(specType, prefix) {
			base := specType[len(prefix):]
			if _, ok := typeMap[base]; ok {
				return true
			}
		}
	}
	return false
}

// resolveType maps a spec type to its Go equivalent using the typeMap.
// Handles composite types like []EGLint → []Int and **EGLint → **Int.
//
// Proven properties (proofs/Proofs/ResolveType.lean):
//   - Identity for empty typeMap
//   - Direct lookup takes precedence over prefix-based
//   - Preserves [], **, * prefix structure
//
// Differential tested against Lean oracle.
func resolveType(specType string, typeMap map[string]string) string {
	if goType, ok := typeMap[specType]; ok {
		return goType
	}
	// Strip slice/pointer prefixes and resolve the base type.
	for _, prefix := range []string{"[]", "**", "*"} {
		if strings.HasPrefix(specType, prefix) {
			base := specType[len(prefix):]
			if goType, ok := typeMap[base]; ok {
				return prefix + goType
			}
		}
	}
	return specType
}

// isFixedArrayGoType returns true if t is a fixed-size array type like "[16]float32".
func isFixedArrayGoType(t string) bool {
	return len(t) > 2 && t[0] == '[' && t[1] != ']'
}

// isReceiverParam returns true if the parameter looks like a method receiver
// (a pointer to an opaque type, not an output parameter).
func isReceiverParam(p specmodel.Param) bool {
	return strings.HasPrefix(p.Type, "*") && p.Direction != "out"
}

// isTypedefKind returns true for all typedef and pointer_handle kind strings.
func isTypedefKind(kind string) bool {
	return strings.HasPrefix(kind, "typedef_") || kind == "pointer_handle"
}

// countGoParams counts the number of parameters in a Go function signature string.
// "func()" → 0, "func(int)" → 1, "func(int, string)" → 2.
//
// Verified by concrete cases (proofs/Proofs/CountGoParams.lean).
// Differential tested against Lean oracle.
func countGoParams(sig string) int {
	// Extract content between ( and ).
	start := strings.Index(sig, "(")
	end := strings.LastIndex(sig, ")")
	if start < 0 || end < 0 || end <= start+1 {
		return 0
	}
	inner := strings.TrimSpace(sig[start+1 : end])
	if inner == "" {
		return 0
	}
	return strings.Count(inner, ",") + 1
}

// callbackAnnotation holds Go callback type info resolved from function overlays.
type callbackAnnotation struct {
	goType string
	goSig  string
}

func mergeCallback(cbName string, cbd specmodel.CallbackDef, annotations map[string]callbackAnnotation) MergedCallback {
	mc := MergedCallback{
		SpecName: cbName,
		Returns:  cbd.Returns,
	}

	if ann, ok := annotations[cbName]; ok {
		mc.GoCallbackType = ann.goType
		mc.GoCallbackSig = ann.goSig
	}

	for _, p := range cbd.Params {
		mc.Params = append(mc.Params, MergedParam{
			Name:      safeGoName(p.Name),
			GoType:    p.Type,
			Direction: p.Direction,
		})
	}
	return mc
}

// collectUnresolvedFuncTypes finds type names used in function params/returns
// that are not yet in the typeMap and are not Go built-in types.
// These need auto-generated type aliases. Scans ALL spec functions (not just
// those with overlay entries) since auto-generation now processes all functions.
func collectUnresolvedFuncTypes(
	spec specmodel.Spec,
	overlay overlaymodel.Overlay,
	typeMap map[string]string,
	opaqueSpecNames map[string]bool,
	specToGoName map[string]string,
) map[string]bool {
	result := make(map[string]bool)

	checkType := func(t string) {
		base := t
		for strings.HasPrefix(base, "*") {
			base = base[1:]
		}
		base = strings.TrimPrefix(base, "[]")
		// Skip array types like [16]float32.
		if len(base) > 0 && base[0] == '[' {
			return
		}
		if base == "" || base == "unsafe.Pointer" || base == "string" ||
			base == "bool" || base == "byte" || base == "func_ptr" {
			return
		}
		if isScalarGoType(base) {
			return
		}
		if typeMap[base] != "" || typeMap["*"+base] != "" {
			return
		}
		// Skip types that are already registered as opaque types.
		if opaqueSpecNames[base] {
			return
		}
		result[base] = true
	}

	// Scan ALL spec functions, not just overlay-listed ones.
	// Skip functions that are explicitly skipped in the overlay.
	for funcName, fd := range spec.Functions {
		if fov, ok := overlay.Functions[funcName]; ok && fov.Skip {
			continue
		}
		for _, p := range fd.Params {
			checkType(p.Type)
		}
		checkType(fd.Returns)
	}

	// Also scan constructor params of opaque types — these reference
	// cross-module types (e.g. ANativeWindow in camera constructors)
	// that need auto-generated type aliases.
	for _, tov := range overlay.Types {
		if tov.Constructor == "" {
			continue
		}
		if fd, ok := spec.Functions[tov.Constructor]; ok {
			for _, p := range fd.Params {
				if p.Direction == "out" {
					continue
				}
				checkType(p.Type)
			}
		}
	}

	return result
}

// sortedKeys returns the sorted keys of a map. It works with any map type
// that has string keys via type parameters.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
