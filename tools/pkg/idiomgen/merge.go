package idiomgen

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/xaionaro-go/ndk/tools/pkg/capigen"
	"github.com/xaionaro-go/ndk/tools/pkg/overlaymodel"
	"github.com/xaionaro-go/ndk/tools/pkg/specmodel"
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

// isScalarGoType returns true for Go scalar types.
func isScalarGoType(t string) bool {
	switch t {
	case "int8", "uint8", "int16", "uint16", "int32", "uint32",
		"int64", "uint64", "float32", "float64", "bool", "int", "uint":
		return true
	}
	return false
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

	// Build set of types that are callback structs or struct accessors (not opaque handles).
	bridgeStructTypes := make(map[string]bool)
	for specName := range overlay.CallbackStructs {
		bridgeStructTypes[specName] = true
	}
	for specName := range overlay.StructAccessors {
		bridgeStructTypes[specName] = true
	}

	// Build opaque spec name set for handle detection.
	// Exclude transparent types — they are type aliases, not struct wrappers.
	opaqueSpecNames := make(map[string]bool)
	for specName, td := range spec.Types {
		if td.Kind == "opaque_ptr" {
			tov, hasOverlay := overlay.Types[specName]
			if !hasOverlay {
				continue // No overlay entry → no struct wrapper → not a handle.
			}
			if tov.Transparent {
				continue
			}
			opaqueSpecNames[specName] = true
		}
	}

	// Build opaque types and populate typeMap for pointer types.
	for specName, td := range spec.Types {
		if td.Kind != "opaque_ptr" {
			continue
		}
		if bridgeStructTypes[specName] {
			continue
		}
		tov, hasOverlay := overlay.Types[specName]
		if !hasOverlay {
			continue
		}
		goName := tov.GoName
		if goName == "" {
			goName = specName
		}

		// Transparent types become simple type aliases (no struct wrapper).
		// Used for void* typedefs like EGLDisplay that are value types in Go.
		if tov.Transparent {
			m.TypeAliases = append(m.TypeAliases, MergedTypeAlias{
				GoName:   goName,
				CapiType: capiExportName(specName),
			})
			typeMap[specName] = goName
			typeMap["*"+specName] = "*" + goName
			continue
		}

		destructorReturnsError := false
		if tov.Destructor != "" {
			if fd, ok := spec.Functions[tov.Destructor]; ok {
				destructorReturnsError = errorTypes[fd.Returns]
			}
		}
		constructorReturnsPointer := false
		var constructorParams []MergedParam
		if tov.Constructor != "" {
			if fd, ok := spec.Functions[tov.Constructor]; ok {
				constructorReturnsPointer = strings.HasPrefix(fd.Returns, "*")
				for _, p := range fd.Params {
					if p.Direction == "out" {
						continue
					}
					goType := resolveType(p.Type, typeMap)
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
						Name:     safeGoName(p.Name),
						GoType:   goType,
						CapiType: capiType,
						IsString: isString,
						IsHandle: isHandle,
					})
				}
			}
		}
		exportedConstructor := ""
		if tov.Constructor != "" {
			exportedConstructor = capiExportName(tov.Constructor)
		}
		exportedDestructor := ""
		if tov.Destructor != "" {
			exportedDestructor = capiExportName(tov.Destructor)
		}
		m.OpaqueTypes[goName] = MergedOpaqueType{
			GoName:                    goName,
			CapiType:                  capiExportName(specName),
			Constructor:               exportedConstructor,
			ConstructorReturnsPointer: constructorReturnsPointer,
			ConstructorParams:         constructorParams,
			Destructor:                exportedDestructor,
			DestructorReturnsError:    destructorReturnsError,
			Pattern:                   tov.Pattern,
			Interfaces:                tov.Interfaces,
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
	// This must happen before the type alias loop so that extra_enums-only
	// types (like MetadataTag) are recognized as enums, not type aliases.
	for enumName, extras := range overlay.ExtraEnums {
		for _, ev := range extras {
			spec.Enums[enumName] = append(spec.Enums[enumName], specmodel.EnumValue{
				Name:  ev.Name,
				Value: ev.Value,
			})
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
		// Skip if this type is classified as an enum in the overlay.
		if _, isEnum := spec.Enums[specName]; isEnum {
			if _, hasOverlay := overlay.Types[specName]; hasOverlay {
				continue
			}
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
	// Only include enums that have an overlay entry.
	enumNames := sortedKeys(spec.Enums)
	for _, enumName := range enumNames {
		tov, hasOverlay := overlay.Types[enumName]
		if !hasOverlay {
			continue
		}
		values := spec.Enums[enumName]

		if tov.GoError {
			m.ErrorEnums = append(m.ErrorEnums, mergeErrorEnum(enumName, tov, values))
			typeMap[enumName] = "Error"
			typeMap["*"+enumName] = "*Error"
		} else {
			goName := tov.GoName
			if goName == "" {
				goName = enumName
			}
			baseType := kindToBaseType(spec.Types[enumName].Kind)
			m.ValueEnums = append(m.ValueEnums, mergeValueEnum(enumName, goName, baseType, tov, values))
			typeMap[enumName] = goName
			typeMap["*"+enumName] = "*" + goName
		}
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
	autoAliasTypes := collectUnresolvedFuncTypes(spec, overlay, typeMap)
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

	// Resolve functions → methods or free functions.
	funcNames := sortedKeys(spec.Functions)
	for _, funcName := range funcNames {
		fd := spec.Functions[funcName]
		fov, ok := overlay.Functions[funcName]
		if !ok {
			continue
		}
		if fov.Skip {
			continue
		}
		switch {
		case fov.Receiver != "":
			m.Methods = append(m.Methods, mergeMethod(funcName, fd, fov, overlay.APILevels, typeMap, receiverCapiTypes, opaqueSpecNames, overlay.CallbackStructs, overlay.StructAccessors))
		case fov.GoName != "":
			m.FreeFunctions = append(m.FreeFunctions, mergeFreeFunction(funcName, fd, fov, overlay.APILevels, typeMap, opaqueSpecNames))
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
			m.Methods = append(m.Methods, mergeMethod(funcName, fd, fov, overlay.APILevels, typeMap, receiverCapiTypes, opaqueSpecNames, overlay.CallbackStructs, overlay.StructAccessors))
		case fov.GoName != "":
			m.FreeFunctions = append(m.FreeFunctions, mergeFreeFunction(funcName, fd, fov, overlay.APILevels, typeMap, opaqueSpecNames))
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
				fov, hasOverlay := csov.Fields[field.Name]
				if !hasOverlay {
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

func mergeMethod(funcName string, fd specmodel.FuncDef, fov overlaymodel.FuncOverlay, apiLevels map[string]int, typeMap map[string]string, receiverCapiTypes map[string]string, opaqueSpecNames map[string]bool, callbackStructOverlays map[string]overlaymodel.CallbackStructOverlay, structAccessors map[string]overlaymodel.StructAccessorOverlay) MergedMethod {
	var params []MergedParam
	receiverFound := false
	for _, p := range fd.Params {
		if !receiverFound && isReceiverParam(p) {
			receiverFound = true
			continue
		}
		goType := resolveType(p.Type, typeMap)
		baseType := strings.TrimPrefix(p.Type, "*")
		isHandle := strings.HasPrefix(p.Type, "*") && opaqueSpecNames[baseType]
		isString := p.Type == "*byte"
		remapped := isTypeRemapped(p.Type, typeMap)
		if isString {
			goType = "string"
		}
		capiType := capiExportType(p.Type)
		if isFixedArrayGoType(goType) {
			goType = "*" + goType
			capiType = "*" + capiType
		}
		params = append(params, MergedParam{
			Name:      safeGoName(p.Name),
			GoType:    goType,
			CapiType:  capiType,
			IsHandle:  isHandle,
			IsString:  isString,
			Remapped:  remapped,
			Direction: p.Direction,
		})
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
		mm.OutputParams = resolveOutputParams(fov.OutputParams, fd, opaqueSpecNames)
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

func mergeFreeFunction(funcName string, fd specmodel.FuncDef, fov overlaymodel.FuncOverlay, apiLevels map[string]int, typeMap map[string]string, opaqueSpecNames map[string]bool) MergedFreeFunction {
	var params []MergedParam
	for _, p := range fd.Params {
		goType := resolveType(p.Type, typeMap)
		baseType := strings.TrimPrefix(p.Type, "*")
		isHandle := strings.HasPrefix(p.Type, "*") && opaqueSpecNames[baseType]
		isString := p.Type == "*byte"
		remapped := isTypeRemapped(p.Type, typeMap)
		if isString {
			goType = "string"
		}
		capiType := capiExportType(p.Type)
		if isFixedArrayGoType(goType) {
			goType = "*" + goType
			capiType = "*" + capiType
		}
		params = append(params, MergedParam{
			Name:      safeGoName(p.Name),
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
		ff.OutputParams = resolveOutputParams(fov.OutputParams, fd, opaqueSpecNames)
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
		localType := strings.TrimPrefix(p.Type, "*")

		// Determine the base type name (without pointer stars).
		specBase := localType
		for strings.HasPrefix(specBase, "*") {
			specBase = specBase[1:]
		}

		// Determine if the Go type is an opaque handle wrapper.
		isHandle := opaqueSpecNames[specBase]

		// Build the capi type for the local variable declaration.
		// Opaque types need capi. prefix (e.g., *capi.AImageReader).
		// Scalar types use the Go type directly (e.g., *uint8, int32).
		var capiType string
		if isHandle || !isScalarGoType(specBase) {
			// Reconstruct with capi. prefix on the base type.
			stars := strings.TrimSuffix(localType, specBase)
			capiType = stars + "capi." + capiExportName(specBase)
		} else {
			capiType = localType
		}

		result = append(result, MergedOutputParam{
			CParamName: p.Name,
			GoType:     goType,
			CapiType:   capiType,
			IsHandle:   isHandle,
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
// (for functions that have overlay entries) that are not yet in the typeMap
// and are not Go built-in types. These need auto-generated type aliases.
func collectUnresolvedFuncTypes(
	spec specmodel.Spec,
	overlay overlaymodel.Overlay,
	typeMap map[string]string,
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
		result[base] = true
	}

	for funcName, fov := range overlay.Functions {
		if fov.Skip {
			continue
		}
		fd, ok := spec.Functions[funcName]
		if !ok {
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
