package headerspec

import (
	"strconv"
	"strings"
)

// Declarations holds all extracted C declarations from the AST.
type Declarations struct {
	Functions []FuncDecl
	Typedefs  []TypedefInfo
	Enums     []EnumInfo
	Structs   []StructInfo
}

// FuncDecl represents an extracted C function declaration.
type FuncDecl struct {
	Name       string
	ReturnType string
	Params     []ParamInfo
}

// ParamInfo describes one function/callback parameter.
type ParamInfo struct {
	Name string
	Type string
}

// TypedefInfo describes a typedef declaration.
type TypedefInfo struct {
	Name               string
	UnderlyingType     string
	DesugaredType      string // desugared qualType from clang (resolves typedef chains)
	IsOpaqueStruct     bool
	IsEnumTypedef      bool // true if this typedefs an enum
	IsFuncPtr          bool
	FuncParams         []ParamInfo
	FuncReturn         string
}

// EnumInfo describes an enum declaration with its constants.
type EnumInfo struct {
	Name        string
	TypedefName string
	FixedType   string
	Constants   []EnumConstant
}

// EnumConstant is one value in an enum.
type EnumConstant struct {
	Name  string
	Value int64
}

// StructInfo describes a struct declaration.
type StructInfo struct {
	Name       string
	IsComplete bool
	Fields     []FieldInfo
}

// FieldInfo describes one field of a struct.
type FieldInfo struct {
	Name       string
	Type       string
	IsFuncPtr  bool
	FuncParams []ParamInfo
	FuncReturn string
}

// ExtractDeclarations walks the root AST node and extracts all declarations
// from the specified target header files.
//
// Clang only emits loc.file for the FIRST node from each file; subsequent
// nodes only have line/col. We track the "current file" while iterating
// through root.Inner to correctly filter by target headers.
func ExtractDeclarations(
	root *ASTNode,
	targetHeaders []string,
) *Declarations {
	decls := &Declarations{}
	headerSet := buildHeaderSet(targetHeaders)

	// Track anonymous enums by their AST ID so we can link them to
	// wrapping TypedefDecl nodes via ownedTagDecl.
	enumByID := map[string]*EnumInfo{}

	currentFile := ""
	for i := range root.Inner {
		node := &root.Inner[i]

		if file := resolveFile(node.Loc); file != "" {
			currentFile = file
		}

		if !matchesTargetHeader(currentFile, headerSet) {
			continue
		}

		// Skip compiler-implicit declarations.
		if node.IsImplicit {
			continue
		}

		switch node.Kind {
		case "FunctionDecl":
			if fn := extractFunction(node); fn != nil {
				decls.Functions = append(decls.Functions, *fn)
			}

		case "TypedefDecl":
			td := extractTypedef(node)
			if td == nil {
				continue
			}

			// Check if this typedef wraps an anonymous enum.
			ownedID := findOwnedTagDeclID(node)
			if ownedID != "" {
				if ei, ok := enumByID[ownedID]; ok {
					ei.TypedefName = td.Name
				}
			}

			decls.Typedefs = append(decls.Typedefs, *td)

		case "EnumDecl":
			ei := extractEnum(node)
			if ei == nil {
				continue
			}
			decls.Enums = append(decls.Enums, *ei)
			enumByID[node.ID] = &decls.Enums[len(decls.Enums)-1]

		case "RecordDecl":
			si := extractStruct(node)
			if si == nil {
				continue
			}
			decls.Structs = append(decls.Structs, *si)
		}
	}

	return decls
}

// buildHeaderSet normalizes target header paths into a set of suffixes
// for matching against absolute paths from clang.
func buildHeaderSet(headers []string) map[string]bool {
	set := make(map[string]bool, len(headers))
	for _, h := range headers {
		set[h] = true
	}
	return set
}

// matchesTargetHeader checks if the current file path ends with any of
// the target header suffixes (e.g., "android/looper.h").
func matchesTargetHeader(filePath string, headers map[string]bool) bool {
	if filePath == "" {
		return false
	}
	for h := range headers {
		if strings.HasSuffix(filePath, "/"+h) || filePath == h {
			return true
		}
	}
	return false
}

// resolveFile extracts the file path from a node's location. Clang may
// place the file in loc.file, loc.spellingLoc.file, or loc.expansionLoc.file.
func resolveFile(loc *Location) string {
	if loc == nil {
		return ""
	}
	if loc.File != "" {
		return loc.File
	}
	if loc.SpellingLoc != nil && loc.SpellingLoc.File != "" {
		return loc.SpellingLoc.File
	}
	if loc.ExpansionLoc != nil && loc.ExpansionLoc.File != "" {
		return loc.ExpansionLoc.File
	}
	return ""
}

// extractFunction builds a FuncDecl from a FunctionDecl AST node.
func extractFunction(node *ASTNode) *FuncDecl {
	if node.Name == "" {
		return nil
	}

	retType, paramTypes := parseFunctionQualType(node.Type)

	fd := &FuncDecl{
		Name:       node.Name,
		ReturnType: retType,
	}

	// Extract parameters from ParmVarDecl inner nodes.
	paramIdx := 0
	for i := range node.Inner {
		child := &node.Inner[i]
		if child.Kind != "ParmVarDecl" {
			continue
		}

		p := ParamInfo{Name: child.Name}
		if child.Type != nil {
			p.Type = child.Type.QualType
		} else if paramIdx < len(paramTypes) {
			p.Type = paramTypes[paramIdx]
		}
		fd.Params = append(fd.Params, p)
		paramIdx++
	}

	return fd
}

// parseFunctionQualType splits a C function qualType like "int (int, void *)"
// into its return type and parameter types.
func parseFunctionQualType(ti *TypeInfo) (string, []string) {
	if ti == nil {
		return "void", nil
	}
	qt := ti.QualType

	// Find the opening paren that starts the parameter list.
	// For "int (*)(int, int)" style, we need to be more careful, but for
	// FunctionDecl qualTypes the format is "ReturnType (Param1, Param2, ...)"
	parenIdx := strings.Index(qt, "(")
	if parenIdx < 0 {
		return qt, nil
	}

	retType := strings.TrimSpace(qt[:parenIdx])
	paramsPart := qt[parenIdx+1 : len(qt)-1]
	if paramsPart == "" || paramsPart == "void" {
		return retType, nil
	}

	params := splitCParams(paramsPart)
	return retType, params
}

// splitCParams splits a comma-separated C parameter list, respecting nested
// parentheses (for function pointer parameters).
func splitCParams(s string) []string {
	var params []string
	depth := 0
	start := 0
	for i, c := range s {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				params = append(params, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	if rest := strings.TrimSpace(s[start:]); rest != "" {
		params = append(params, rest)
	}
	return params
}

// extractTypedef builds a TypedefInfo from a TypedefDecl AST node.
func extractTypedef(node *ASTNode) *TypedefInfo {
	if node.Name == "" || node.Type == nil {
		return nil
	}

	td := &TypedefInfo{
		Name:           node.Name,
		UnderlyingType: node.Type.QualType,
		DesugaredType:  node.Type.DesugaredQualType,
	}

	qt := node.Type.QualType
	desugared := node.Type.DesugaredQualType

	// Detect opaque struct: "struct X" typedef with no complete definition.
	// Also check the desugared type for typedef-to-typedef chains that
	// ultimately resolve to a struct.
	if strings.HasPrefix(qt, "struct ") || strings.HasPrefix(desugared, "struct ") {
		td.IsOpaqueStruct = true
	}

	// Detect enum typedef: "enum X" qualType or desugared type.
	if strings.HasPrefix(qt, "enum ") {
		td.IsEnumTypedef = true
	}

	// Detect function pointer: qualType contains "(*)" pattern.
	if isFuncPtrQualType(qt) {
		td.IsFuncPtr = true
		td.FuncReturn, td.FuncParams = parseFuncPtrFromAST(node)
	}

	return td
}

// isFuncPtrQualType checks whether a qualType string represents a function pointer.
func isFuncPtrQualType(qt string) bool {
	return strings.Contains(qt, "(*)")
}

// parseFuncPtrFromAST extracts function pointer return type and parameter
// types from the inner AST nodes of a TypedefDecl.
//
// The inner node chain is:
//
//	PointerType → ParenType → FunctionProtoType
//
// FunctionProtoType.Inner contains: [returnType, param1Type, param2Type, ...]
// Note: parameter NAMES are not available in the AST for function pointer
// typedefs; only types.
func parseFuncPtrFromAST(node *ASTNode) (string, []ParamInfo) {
	fpt := findFunctionProtoType(node.Inner)
	if fpt == nil {
		// Fall back to parsing the qualType string.
		return parseFuncPtrQualType(node.Type.QualType)
	}

	if len(fpt.Inner) == 0 {
		return "void", nil
	}

	// First inner node is the return type.
	retType := "void"
	if fpt.Inner[0].Type != nil {
		retType = fpt.Inner[0].Type.QualType
	}

	// Remaining inner nodes are parameter types.
	var params []ParamInfo
	for i := 1; i < len(fpt.Inner); i++ {
		child := &fpt.Inner[i]
		paramType := ""
		if child.Type != nil {
			paramType = child.Type.QualType
		}
		params = append(params, ParamInfo{
			Name: "", // Names aren't available in the AST for func ptr typedefs.
			Type: paramType,
		})
	}

	return retType, params
}

// findFunctionProtoType recursively searches inner nodes for a FunctionProtoType.
func findFunctionProtoType(nodes []ASTNode) *ASTNode {
	for i := range nodes {
		n := &nodes[i]
		if n.Kind == "FunctionProtoType" {
			return n
		}
		if found := findFunctionProtoType(n.Inner); found != nil {
			return found
		}
	}
	return nil
}

// parseFuncPtrQualType is a fallback that parses "RetType (*)(P1, P2, ...)"
// from the qualType string.
func parseFuncPtrQualType(qt string) (string, []ParamInfo) {
	// Format: "RetType (*)(P1, P2, ...)"
	starIdx := strings.Index(qt, "(*)")
	if starIdx < 0 {
		return "void", nil
	}

	retType := strings.TrimSpace(qt[:starIdx])

	// Extract params between the last "(" and ")"
	rest := qt[starIdx+3:]
	if len(rest) < 2 || rest[0] != '(' {
		return retType, nil
	}
	rest = rest[1 : len(rest)-1] // strip outer parens

	if rest == "" || rest == "void" {
		return retType, nil
	}

	parts := splitCParams(rest)
	var params []ParamInfo
	for _, p := range parts {
		params = append(params, ParamInfo{Type: strings.TrimSpace(p)})
	}
	return retType, params
}

// findOwnedTagDeclID looks for an ownedTagDecl in the inner nodes of a
// TypedefDecl, which links a typedef to its anonymous enum or struct.
func findOwnedTagDeclID(node *ASTNode) string {
	for i := range node.Inner {
		child := &node.Inner[i]
		if child.OwnedTagDecl != nil {
			return child.OwnedTagDecl.ID
		}
		if id := findOwnedTagDeclID(child); id != "" {
			return id
		}
	}
	return ""
}

// extractEnum builds an EnumInfo from an EnumDecl AST node.
func extractEnum(node *ASTNode) *EnumInfo {
	ei := &EnumInfo{
		Name: node.Name,
	}

	if node.FixedUnderlyingType != nil {
		ei.FixedType = node.FixedUnderlyingType.QualType
	}

	for i := range node.Inner {
		child := &node.Inner[i]
		if child.Kind != "EnumConstantDecl" {
			continue
		}

		ec := EnumConstant{Name: child.Name}
		ec.Value = extractEnumConstantValue(child)
		ei.Constants = append(ei.Constants, ec)
	}

	if len(ei.Constants) == 0 {
		return nil
	}

	return ei
}

// extractEnumConstantValue extracts the integer value from an EnumConstantDecl.
// The value is typically in a ConstantExpr node, which may be a direct child
// or nested inside an ImplicitCastExpr (for typed enums like "enum X : int8_t").
func extractEnumConstantValue(node *ASTNode) int64 {
	if v, ok := findConstantExprValue(node.Inner); ok {
		return v
	}
	return 0
}

// findConstantExprValue recursively searches for a ConstantExpr with a value.
func findConstantExprValue(nodes []ASTNode) (int64, bool) {
	for i := range nodes {
		child := &nodes[i]
		if child.Kind == "ConstantExpr" && child.Value != "" {
			v, err := strconv.ParseInt(child.Value, 10, 64)
			if err == nil {
				return v, true
			}
		}
		if v, ok := findConstantExprValue(child.Inner); ok {
			return v, ok
		}
	}
	return 0, false
}

// extractStruct builds a StructInfo from a RecordDecl AST node.
func extractStruct(node *ASTNode) *StructInfo {
	if node.Name == "" || node.TagUsed != "struct" {
		return nil
	}

	si := &StructInfo{
		Name:       node.Name,
		IsComplete: node.CompleteDefinition,
	}

	if !si.IsComplete {
		return si
	}

	for i := range node.Inner {
		child := &node.Inner[i]
		if child.Kind != "FieldDecl" {
			continue
		}
		fi := extractField(child)
		si.Fields = append(si.Fields, fi)
	}

	return si
}

// extractField builds a FieldInfo from a FieldDecl AST node.
func extractField(node *ASTNode) FieldInfo {
	fi := FieldInfo{
		Name: node.Name,
	}

	if node.Type != nil {
		fi.Type = node.Type.QualType

		// Use desugared type for function pointer detection when available.
		desugared := node.Type.DesugaredQualType
		if desugared == "" {
			desugared = fi.Type
		}

		if isFuncPtrQualType(desugared) {
			fi.IsFuncPtr = true
			fi.FuncReturn, fi.FuncParams = parseFuncPtrQualType(desugared)
		}
	}

	return fi
}
