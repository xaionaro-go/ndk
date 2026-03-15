// Package specgen extracts structured API specs from Go source files produced by c-for-go.
package specgen

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
)

// ParseSources parses Go source files and extracts types, enums, functions,
// and callbacks into a Spec. The files are expected to be c-for-go output
// (pure Go with C.* selector expressions for opaque handle types).
func ParseSources(module, sourcePkg string, filePaths []string) (specmodel.Spec, error) {
	spec := specmodel.Spec{
		Module:        module,
		SourcePackage: sourcePkg,
		Types:         make(map[string]specmodel.TypeDef),
		Enums:         make(map[string][]specmodel.EnumValue),
		Functions:     make(map[string]specmodel.FuncDef),
		Callbacks:     make(map[string]specmodel.CallbackDef),
	}

	fset := token.NewFileSet()
	for _, path := range filePaths {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return spec, fmt.Errorf("parse %s: %w", path, err)
		}
		extractTypes(f, &spec)
		extractConsts(f, &spec)
		extractFunctions(f, &spec)
	}

	return spec, nil
}

// extractTypes processes type declarations: opaque handles (C.* selectors),
// integer typedefs, and callback function types.
func extractTypes(f *ast.File, spec *specmodel.Spec) {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, s := range genDecl.Specs {
			ts := s.(*ast.TypeSpec)
			name := ts.Name.Name

			switch typ := ts.Type.(type) {
			case *ast.SelectorExpr:
				ident, ok := typ.X.(*ast.Ident)
				if !ok {
					break
				}
				switch {
				case ident.Name == "C":
					// type Foo C.Bar → opaque handle (C package selector).
					spec.Types[name] = specmodel.TypeDef{
						Kind:   "opaque_ptr",
						CType:  typ.Sel.Name,
						GoType: "*C." + typ.Sel.Name,
					}
				case ident.Name == "unsafe" && typ.Sel.Name == "Pointer":
					// type Foo unsafe.Pointer → pointer-sized handle (EGL/etc).
					spec.Types[name] = specmodel.TypeDef{
						Kind:   "pointer_handle",
						CType:  name,
						GoType: "unsafe.Pointer",
					}
				}

			case *ast.Ident:
				// type Foo int32 → integer typedef.
				kind := classifyIdent(typ.Name)
				if kind != "" {
					spec.Types[name] = specmodel.TypeDef{
						Kind:   kind,
						CType:  name,
						GoType: typ.Name,
					}
				}

			case *ast.FuncType:
				// type Foo func(...) → callback.
				spec.Callbacks[name] = specmodel.CallbackDef{
					Params:  extractParamList(typ.Params),
					Returns: returnTypeString(typ.Results),
				}
			}
		}
	}
}

// classifyIdent maps Go built-in integer type names to spec kind strings.
func classifyIdent(name string) string {
	switch name {
	case "int32":
		return "typedef_int32"
	case "uint32":
		return "typedef_uint32"
	case "int64":
		return "typedef_int64"
	case "uint64":
		return "typedef_uint64"
	case "int":
		return "typedef_int"
	case "uint":
		return "typedef_uint"
	case "int8":
		return "typedef_int8"
	case "uint8":
		return "typedef_uint8"
	case "int16":
		return "typedef_int16"
	case "uint16":
		return "typedef_uint16"
	case "float32":
		return "typedef_float32"
	case "float64":
		return "typedef_float64"
	default:
		return ""
	}
}

// extractConsts processes const blocks, grouping typed constants into enum values.
// Handles negative values expressed as UnaryExpr with SUB operator.
func extractConsts(f *ast.File, spec *specmodel.Spec) {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		var currentType string
		for _, s := range genDecl.Specs {
			vs := s.(*ast.ValueSpec)

			// Track the current type (iota-style const blocks carry the type forward).
			if vs.Type != nil {
				currentType = typeString(vs.Type)
			}
			if currentType == "" {
				continue
			}

			for i, name := range vs.Names {
				if i >= len(vs.Values) {
					continue
				}
				val, ok := evalConstInt(vs.Values[i])
				if !ok {
					continue
				}
				spec.Enums[currentType] = append(spec.Enums[currentType], specmodel.EnumValue{
					Name:  name.Name,
					Value: val,
				})
			}
		}
	}
}

// evalConstInt evaluates a constant expression to an int64.
// Supports BasicLit (integer) and UnaryExpr with SUB (negative).
func evalConstInt(expr ast.Expr) (int64, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind != token.INT {
			return 0, false
		}
		n, err := strconv.ParseInt(v.Value, 0, 64)
		return n, err == nil

	case *ast.UnaryExpr:
		if v.Op != token.SUB {
			return 0, false
		}
		inner, ok := evalConstInt(v.X)
		if !ok {
			return 0, false
		}
		return -inner, true

	default:
		return 0, false
	}
}

// extractFunctions processes package-level exported function declarations.
// Double-pointer params (**SomeType) are detected as output params.
func extractFunctions(f *ast.File, spec *specmodel.Spec) {
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// Skip methods (have a receiver) and unexported functions.
		if funcDecl.Recv != nil || !funcDecl.Name.IsExported() {
			continue
		}
		// Skip JNI infrastructure, bridge, and c-for-go ref-helper functions
		// (not part of the NDK API).
		name := funcDecl.Name.Name
		if strings.HasPrefix(name, "JNI_") || strings.HasPrefix(name, "Bridge") || strings.HasPrefix(name, "New") {
			continue
		}

		fd := specmodel.FuncDef{
			CName:   funcDecl.Name.Name,
			Params:  extractParamList(funcDecl.Type.Params),
			Returns: returnTypeString(funcDecl.Type.Results),
		}
		spec.Functions[funcDecl.Name.Name] = fd
	}
}

// extractParamList converts an ast.FieldList into a slice of Param.
// Detects double-pointer params (**T) as output direction.
func extractParamList(fields *ast.FieldList) []specmodel.Param {
	if fields == nil {
		return nil
	}
	var params []specmodel.Param
	for _, field := range fields.List {
		ts := typeString(field.Type)
		dir := ""
		if isDoublePointer(field.Type) {
			dir = "out"
		}

		if len(field.Names) == 0 {
			// Unnamed parameter.
			params = append(params, specmodel.Param{
				Name:      "",
				Type:      ts,
				Direction: dir,
			})
			continue
		}
		for _, name := range field.Names {
			params = append(params, specmodel.Param{
				Name:      name.Name,
				Type:      ts,
				Direction: dir,
			})
		}
	}
	return params
}

// returnTypeString extracts the return type string from a function's result list.
// Returns empty string for no return or multiple returns.
func returnTypeString(results *ast.FieldList) string {
	if results == nil || len(results.List) == 0 {
		return ""
	}
	if len(results.List) == 1 {
		return typeString(results.List[0].Type)
	}
	return ""
}

// isDoublePointer returns true if expr is **T (StarExpr wrapping StarExpr).
func isDoublePointer(expr ast.Expr) bool {
	outer, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	_, ok = outer.X.(*ast.StarExpr)
	return ok
}

// typeString converts an AST type expression to its string representation.
func typeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + typeString(t.Elt)
		}
		return "[...]" + typeString(t.Elt)
	case *ast.Ellipsis:
		return "..." + typeString(t.Elt)
	case *ast.MapType:
		return "map[" + typeString(t.Key) + "]" + typeString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// MergeFunctions adds function definitions from parsed C headers into the spec.
// Only adds functions that are not already present (Go AST extraction takes priority).
func MergeFunctions(
	spec *specmodel.Spec,
	funcs map[string]specmodel.FuncDef,
) {
	if spec.Functions == nil {
		spec.Functions = make(map[string]specmodel.FuncDef)
	}

	for name, fd := range funcs {
		if _, exists := spec.Functions[name]; exists {
			continue
		}
		spec.Functions[name] = fd
	}
}

// MergeStructs adds struct definitions from parsed C headers into the spec.
// Only includes structs whose type name exists in spec.Types.
func MergeStructs(
	spec *specmodel.Spec,
	structs map[string]specmodel.StructDef,
) {
	if spec.Structs == nil {
		spec.Structs = make(map[string]specmodel.StructDef)
	}

	for name, sd := range structs {
		spec.Structs[name] = sd
	}
}

var cgoIncludeRe = regexp.MustCompile(`#include\s+[<"]([^>"]+)[>"]`)

// ExtractIncludeDirsFromGoFiles scans Go source files for CGo preamble
// #include <path> directives and returns the unique set of parent directories
// resolved under the given sysroot include path.
func ExtractIncludeDirsFromGoFiles(
	goFiles []string,
	sysrootInclude string,
) ([]string, error) {
	seen := make(map[string]struct{})
	var dirs []string

	for _, goFile := range goFiles {
		includes, err := extractCGoIncludes(goFile)
		if err != nil {
			return nil, fmt.Errorf("extract includes from %s: %w", goFile, err)
		}

		for _, inc := range includes {
			dir := filepath.Dir(inc)
			fullDir := filepath.Join(sysrootInclude, dir)

			if _, ok := seen[fullDir]; ok {
				continue
			}
			seen[fullDir] = struct{}{}

			if info, err := os.Stat(fullDir); err == nil && info.IsDir() {
				dirs = append(dirs, fullDir)
			}
		}
	}

	return dirs, nil
}

// extractCGoIncludes reads a Go source file and returns all #include <path>
// directives found in CGo comment preambles.
func extractCGoIncludes(goFile string) ([]string, error) {
	f, err := os.Open(goFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var includes []string
	scanner := bufio.NewScanner(f)
	inCgo := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Detect CGo preamble blocks: lines starting with // #include or
		// following a // #cgo directive, or lines with #include inside /* */ blocks.
		if strings.HasPrefix(trimmed, "// #") || strings.HasPrefix(trimmed, "//#") {
			inCgo = true
			text := strings.TrimPrefix(trimmed, "//")
			text = strings.TrimSpace(text)
			if m := cgoIncludeRe.FindStringSubmatch(text); m != nil {
				includes = append(includes, m[1])
			}
			continue
		}

		// Also handle multi-line /* */ CGo comments.
		if strings.HasPrefix(trimmed, "/*") {
			inCgo = true
		}
		if inCgo && strings.Contains(trimmed, "*/") {
			inCgo = false
		}
		if inCgo {
			if m := cgoIncludeRe.FindStringSubmatch(line); m != nil {
				includes = append(includes, m[1])
			}
			continue
		}

		// CGo preamble ends at a non-comment line.
		if inCgo && !strings.HasPrefix(trimmed, "//") {
			inCgo = false
		}
	}

	return includes, scanner.Err()
}
