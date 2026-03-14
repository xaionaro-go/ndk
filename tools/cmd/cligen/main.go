// Command cligen generates cobra CLI subcommands for all exported
// functions and methods in the NDK idiomatic packages. It parses
// Go source with go/ast and emits one _gen.go file per package
// into cmd/ndkcli/.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const goModule = "github.com/xaionaro-go/ndk"

var aliasMap = map[string]string{
	"log":   "ndklog",
	"sync":  "ndksync",
	"image": "ndkimage",
}

var fileSuffixMap = map[string]string{
	"log":  "log_cmd",
	"sync": "sync_cmd",
}

type pkgInfo struct {
	name       string
	dir        string
	alias      string
	types      map[string]*typeInfo
	funcs      []funcInfo
	importPath string
}

type typeInfo struct {
	name        string
	constructor *funcInfo
	methods     []funcInfo
	hasClose    bool
}

type funcInfo struct {
	name    string
	recv    string
	params  []paramInfo
	returns []returnInfo
}

type paramInfo struct {
	name   string
	goType string
}

type returnInfo struct {
	goType string
}

func main() {
	projectRoot := flag.String("root", ".", "Project root directory")
	outputDir := flag.String("output", "cmd/ndkcli", "Output directory")
	flag.Parse()

	pkgDirs := discoverPackages(*projectRoot)
	for _, dir := range pkgDirs {
		pkg, err := parsePackage(filepath.Join(*projectRoot, dir), dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", dir, err)
			continue
		}
		if pkg == nil || (len(pkg.funcs) == 0 && len(pkg.types) == 0) {
			fmt.Fprintf(os.Stderr, "skip %s: no exported API\n", dir)
			continue
		}

		code := generateFile(pkg)
		if code == "" {
			fmt.Fprintf(os.Stderr, "skip %s: no CLI-compatible API\n", dir)
			continue
		}

		formatted, err := format.Source([]byte(code))
		if err != nil {
			fmt.Fprintf(os.Stderr, "format error for %s: %v\nraw:\n%s\n", dir, err, code)
			continue
		}

		base := dir
		if s, ok := fileSuffixMap[dir]; ok {
			base = s
		}
		outPath := filepath.Join(*projectRoot, *outputDir, base+"_gen.go")
		if err := os.WriteFile(outPath, formatted, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
			continue
		}
		fmt.Printf("generated %s (%d commands)\n", outPath, strings.Count(code, "RunE:"))
	}
}

func discoverPackages(root string) []string {
	skip := map[string]bool{
		"capi": true, "spec": true, "templates": true, "tools": true,
		"tests": true, "examples": true, "internal": true, "docs": true,
		"discussion": true, "jni": true, "cmd": true, "display": true,
	}
	entries, _ := os.ReadDir(root)
	var dirs []string
	for _, e := range entries {
		if !e.IsDir() || skip[e.Name()] || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		goFiles, _ := filepath.Glob(filepath.Join(root, e.Name(), "*.go"))
		if len(goFiles) > 0 {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	return dirs
}

func parsePackage(dir, dirName string) (*pkgInfo, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, err
	}

	for pkgName, pkg := range pkgs {
		info := &pkgInfo{
			name:       pkgName,
			dir:        dirName,
			alias:      aliasMap[dirName],
			types:      map[string]*typeInfo{},
			importPath: goModule + "/" + dirName,
		}

		for _, file := range pkg.Files {
			extractFromFile(file, info)
		}

		// Identify constructors.
		var remaining []funcInfo
		for i := range info.funcs {
			fn := &info.funcs[i]
			typeName := constructorFor(fn)
			if typeName != "" {
				if ti, ok := info.types[typeName]; ok {
					fnCopy := *fn
					ti.constructor = &fnCopy
					continue
				}
			}
			remaining = append(remaining, *fn)
		}
		info.funcs = remaining

		// Check Close().
		for _, ti := range info.types {
			for _, m := range ti.methods {
				if m.name == "Close" {
					ti.hasClose = true
				}
			}
		}

		return info, nil
	}
	return nil, nil
}

func extractFromFile(file *ast.File, info *pkgInfo) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || !ts.Name.IsExported() {
					continue
				}
				if _, ok := ts.Type.(*ast.StructType); ok {
					if _, exists := info.types[ts.Name.Name]; !exists {
						info.types[ts.Name.Name] = &typeInfo{name: ts.Name.Name}
					}
				}
			}
		case *ast.FuncDecl:
			if !d.Name.IsExported() {
				continue
			}
			fn := extractFunc(d)
			if fn == nil {
				continue
			}
			if fn.recv != "" {
				ti, ok := info.types[fn.recv]
				if !ok {
					ti = &typeInfo{name: fn.recv}
					info.types[fn.recv] = ti
				}
				ti.methods = append(ti.methods, *fn)
			} else {
				info.funcs = append(info.funcs, *fn)
			}
		}
	}
}

func extractFunc(d *ast.FuncDecl) *funcInfo {
	fn := &funcInfo{name: d.Name.Name}
	if d.Recv != nil && len(d.Recv.List) > 0 {
		fn.recv = recvTypeName(d.Recv.List[0].Type)
		if fn.recv == "" {
			return nil
		}
	}
	if d.Type.Params != nil {
		for _, field := range d.Type.Params.List {
			goType := typeString(field.Type)
			names := fieldNames(field)
			for _, name := range names {
				fn.params = append(fn.params, paramInfo{name: name, goType: goType})
			}
		}
	}
	if d.Type.Results != nil {
		for _, field := range d.Type.Results.List {
			fn.returns = append(fn.returns, returnInfo{goType: typeString(field.Type)})
		}
	}
	return fn
}

func recvTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			return id.Name
		}
	case *ast.Ident:
		return t.Name
	}
	return ""
}

func typeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + typeString(t.Elt)
		}
		return "[...]" + typeString(t.Elt)
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.MapType:
		return "map"
	case *ast.FuncType:
		return "func"
	case *ast.Ellipsis:
		return "..." + typeString(t.Elt)
	default:
		return "unknown"
	}
}

func fieldNames(field *ast.Field) []string {
	if len(field.Names) == 0 {
		return []string{"value"}
	}
	var names []string
	for _, n := range field.Names {
		names = append(names, n.Name)
	}
	return names
}

// isCLIType returns the cobra flag getter name, or "" if unsupported.
func isCLIType(goType string) string {
	switch goType {
	case "string":
		return "GetString"
	case "int32":
		return "GetInt32"
	case "int64":
		return "GetInt64"
	case "int":
		return "GetInt"
	case "uint32":
		return "GetUint32"
	case "uint64":
		return "GetUint64"
	case "uint16":
		return "GetInt32" // promote
	case "float32":
		return "GetFloat32"
	case "float64":
		return "GetFloat64"
	case "bool":
		return "GetBool"
	}
	return ""
}

// flagRegister returns the cobra flag registration method name.
func flagRegister(goType string) string {
	switch goType {
	case "string":
		return "String"
	case "int32":
		return "Int32"
	case "int64":
		return "Int64"
	case "int":
		return "Int"
	case "uint32":
		return "Uint32"
	case "uint64":
		return "Uint64"
	case "uint16":
		return "Int32"
	case "float32":
		return "Float32"
	case "float64":
		return "Float64"
	case "bool":
		return "Bool"
	}
	return ""
}

// flagDefault returns the zero-value literal for the type.
func flagDefault(goType string) string {
	switch goType {
	case "string":
		return `""`
	case "bool":
		return "false"
	default:
		return "0"
	}
}

func isCLICompatible(fn funcInfo) bool {
	for _, p := range fn.params {
		if isCLIType(p.goType) == "" {
			return false
		}
	}
	return true
}

func shouldSkip(fn funcInfo) bool {
	switch fn.name {
	case "Pointer", "Close":
		return true
	}
	return strings.HasSuffix(fn.name, "FromPointer")
}

func constructorFor(fn *funcInfo) string {
	if strings.HasSuffix(fn.name, "FromPointer") {
		return ""
	}
	if fn.name == "GetInstance" {
		if len(fn.returns) > 0 {
			return strings.TrimPrefix(fn.returns[0].goType, "*")
		}
		return ""
	}
	if !strings.HasPrefix(fn.name, "New") {
		return ""
	}
	typeName := strings.TrimPrefix(fn.name, "New")

	// Verify the function actually returns *TypeName or (*TypeName, error).
	if len(fn.returns) == 0 {
		return ""
	}
	firstRet := fn.returns[0].goType
	if firstRet != "*"+typeName {
		return ""
	}
	return typeName
}

func generateFile(pkg *pkgInfo) string {
	var buf bytes.Buffer
	pkgPrefix := sanitizeVarPrefix(pkg.dir)
	pkgRef := pkg.name
	if pkg.alias != "" {
		pkgRef = pkg.alias
	}

	type leafCmd struct {
		parentVar string
		varName   string
		use       string
		short     string
		body      string   // RunE body
		flagRegs  []string // flag registration lines
	}
	type groupCmd struct {
		parentVar string
		varName   string
		use       string
		short     string
	}

	var leafs []leafCmd
	var groups []groupCmd
	usedVarNames := map[string]bool{}

	// Package-level functions.
	for _, fn := range pkg.funcs {
		if shouldSkip(fn) || !isCLICompatible(fn) {
			continue
		}
		varName := uniqueVar(pkgPrefix+camelCase(fn.name)+"Cmd", usedVarNames)
		body, flagRegs := genFuncBody(pkgRef, fn, varName)
		leafs = append(leafs, leafCmd{
			parentVar: pkgPrefix + "Cmd",
			varName:   varName,
			use:       kebabCase(fn.name),
			short:     pkg.name + "." + fn.name + "()",
			body:      body,
			flagRegs:  flagRegs,
		})
	}

	// Types with methods.
	typeNames := sortedTypeNames(pkg.types)
	for _, typeName := range typeNames {
		ti := pkg.types[typeName]
		typeVar := uniqueVar(pkgPrefix+typeName+"Cmd", usedVarNames)
		var typeLeafs []leafCmd

		// Constructor.
		if ti.constructor != nil && isCLICompatible(*ti.constructor) {
			varName := uniqueVar(pkgPrefix+typeName+"NewCmd", usedVarNames)
			body, flagRegs := genConstructorBody(pkgRef, ti, varName)
			typeLeafs = append(typeLeafs, leafCmd{
				parentVar: typeVar,
				varName:   varName,
				use:       "new",
				short:     "Create " + typeName,
				body:      body,
				flagRegs:  flagRegs,
			})
		}

		// Methods.
		for _, m := range ti.methods {
			if shouldSkip(m) || !isCLICompatible(m) {
				continue
			}
			varName := uniqueVar(pkgPrefix+typeName+camelCase(m.name)+"Cmd", usedVarNames)
			body, flagRegs := genMethodBody(pkgRef, ti, m, varName)
			typeLeafs = append(typeLeafs, leafCmd{
				parentVar: typeVar,
				varName:   varName,
				use:       kebabCase(m.name),
				short:     typeName + "." + m.name + "()",
				body:      body,
				flagRegs:  flagRegs,
			})
		}

		if len(typeLeafs) == 0 {
			continue
		}

		groups = append(groups, groupCmd{
			parentVar: pkgPrefix + "Cmd",
			varName:   typeVar,
			use:       kebabCase(typeName),
			short:     typeName + " operations",
		})
		leafs = append(leafs, typeLeafs...)
	}

	if len(leafs) == 0 {
		return ""
	}

	// Emit file.
	buf.WriteString("// Code generated by cligen. DO NOT EDIT.\n\n")
	buf.WriteString("package main\n\n")
	// Check if any command actually uses the package.
	needsPkgImport := false
	for _, l := range leafs {
		if !strings.Contains(l.body, "requires external context") {
			needsPkgImport = true
			break
		}
	}

	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\t\"github.com/spf13/cobra\"\n")
	if needsPkgImport {
		if pkg.alias != "" {
			fmt.Fprintf(&buf, "\t%s %q\n", pkg.alias, pkg.importPath)
		} else {
			fmt.Fprintf(&buf, "\t%q\n", pkg.importPath)
		}
	}
	buf.WriteString(")\n\n")

	// Top-level package command.
	fmt.Fprintf(&buf, "var %sCmd = &cobra.Command{\n", pkgPrefix)
	fmt.Fprintf(&buf, "\tUse:   %q,\n", pkg.dir)
	fmt.Fprintf(&buf, "\tShort: %q,\n", pkg.name+" NDK module")
	fmt.Fprintf(&buf, "}\n\n")

	// Group commands.
	for _, g := range groups {
		fmt.Fprintf(&buf, "var %s = &cobra.Command{\n", g.varName)
		fmt.Fprintf(&buf, "\tUse:   %q,\n", g.use)
		fmt.Fprintf(&buf, "\tShort: %q,\n", g.short)
		fmt.Fprintf(&buf, "}\n\n")
	}

	// Leaf commands.
	for _, l := range leafs {
		fmt.Fprintf(&buf, "var %s = &cobra.Command{\n", l.varName)
		fmt.Fprintf(&buf, "\tUse:   %q,\n", l.use)
		fmt.Fprintf(&buf, "\tShort: %q,\n", l.short)
		fmt.Fprintf(&buf, "\tRunE: func(cmd *cobra.Command, args []string) error {\n")
		buf.WriteString(l.body)
		fmt.Fprintf(&buf, "\t},\n")
		fmt.Fprintf(&buf, "}\n\n")
	}

	// init().
	buf.WriteString("func init() {\n")

	// Flag registrations.
	for _, l := range leafs {
		for _, reg := range l.flagRegs {
			buf.WriteString(reg)
		}
	}

	// AddCommand.
	for _, g := range groups {
		fmt.Fprintf(&buf, "\t%s.AddCommand(%s)\n", g.parentVar, g.varName)
	}
	for _, l := range leafs {
		fmt.Fprintf(&buf, "\t%s.AddCommand(%s)\n", l.parentVar, l.varName)
	}
	fmt.Fprintf(&buf, "\trootCmd.AddCommand(%sCmd)\n", pkgPrefix)
	buf.WriteString("}\n")

	return buf.String()
}

func genFuncBody(
	pkgRef string,
	fn funcInfo,
	cmdVar string,
) (string, []string) {
	var body bytes.Buffer
	var flagRegs []string

	// Read flags.
	for _, p := range fn.params {
		getter := isCLIType(p.goType)
		flagName := kebabCase(p.name)
		regMethod := flagRegister(p.goType)
		defVal := flagDefault(p.goType)
		flagRegs = append(flagRegs, fmt.Sprintf("\t%s.Flags().%s(%q, %s, %q)\n",
			cmdVar, regMethod, flagName, defVal, p.name))

		if p.goType == "uint16" {
			fmt.Fprintf(&body, "\t\t%sRaw, _ := cmd.Flags().%s(%q)\n", p.name, getter, flagName)
			fmt.Fprintf(&body, "\t\t%s := uint16(%sRaw)\n", p.name, p.name)
		} else {
			fmt.Fprintf(&body, "\t\t%s, _ := cmd.Flags().%s(%q)\n", p.name, getter, flagName)
		}
	}

	callArgs := paramNames(fn.params)
	hasError := len(fn.returns) > 0 && fn.returns[len(fn.returns)-1].goType == "error"
	hasResult := len(fn.returns) > 0 && fn.returns[0].goType != "error"

	writeCall(&body, pkgRef+"."+fn.name, callArgs, hasResult, hasError, false)

	return body.String(), flagRegs
}

func genConstructorBody(
	pkgRef string,
	ti *typeInfo,
	cmdVar string,
) (string, []string) {
	var body bytes.Buffer
	var flagRegs []string
	fn := ti.constructor

	for _, p := range fn.params {
		getter := isCLIType(p.goType)
		flagName := kebabCase(p.name)
		regMethod := flagRegister(p.goType)
		defVal := flagDefault(p.goType)
		flagRegs = append(flagRegs, fmt.Sprintf("\t%s.Flags().%s(%q, %s, %q)\n",
			cmdVar, regMethod, flagName, defVal, p.name))

		if p.goType == "uint16" {
			fmt.Fprintf(&body, "\t\t%sRaw, _ := cmd.Flags().%s(%q)\n", p.name, getter, flagName)
			fmt.Fprintf(&body, "\t\t%s := uint16(%sRaw)\n", p.name, p.name)
		} else {
			fmt.Fprintf(&body, "\t\t%s, _ := cmd.Flags().%s(%q)\n", p.name, getter, flagName)
		}
	}

	callArgs := paramNames(fn.params)
	hasError := len(fn.returns) > 1 && fn.returns[len(fn.returns)-1].goType == "error"

	if hasError {
		objVar := "obj"
		if !ti.hasClose {
			objVar = "_"
		}
		fmt.Fprintf(&body, "\t\t%s, err := %s.%s(%s)\n", objVar, pkgRef, fn.name, callArgs)
		body.WriteString("\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n")
	} else {
		if ti.hasClose {
			fmt.Fprintf(&body, "\t\tobj := %s.%s(%s)\n", pkgRef, fn.name, callArgs)
		} else {
			fmt.Fprintf(&body, "\t\t_ = %s.%s(%s)\n", pkgRef, fn.name, callArgs)
		}
	}
	if ti.hasClose {
		body.WriteString("\t\tdefer obj.Close()\n")
	}
	body.WriteString("\t\tfmt.Println(\"created successfully\")\n")
	body.WriteString("\t\treturn nil\n")

	return body.String(), flagRegs
}

func genMethodBody(
	pkgRef string,
	ti *typeInfo,
	m funcInfo,
	cmdVar string,
) (string, []string) {
	var body bytes.Buffer
	var flagRegs []string

	// Constructor params as flags too.
	if ti.constructor != nil {
		for _, p := range ti.constructor.params {
			getter := isCLIType(p.goType)
			flagName := "ctor-" + kebabCase(p.name)
			regMethod := flagRegister(p.goType)
			defVal := flagDefault(p.goType)
			flagRegs = append(flagRegs, fmt.Sprintf("\t%s.Flags().%s(%q, %s, \"constructor: %s\")\n",
				cmdVar, regMethod, flagName, defVal, p.name))

			if p.goType == "uint16" {
				fmt.Fprintf(&body, "\t\t%sRaw, _ := cmd.Flags().%s(%q)\n", "ctor"+camelCase(p.name), getter, flagName)
				fmt.Fprintf(&body, "\t\t%s := uint16(%sRaw)\n", "ctor"+camelCase(p.name), "ctor"+camelCase(p.name))
			} else {
				fmt.Fprintf(&body, "\t\t%s, _ := cmd.Flags().%s(%q)\n", "ctor"+camelCase(p.name), getter, flagName)
			}
		}
	}

	// Method params.
	for _, p := range m.params {
		getter := isCLIType(p.goType)
		flagName := kebabCase(p.name)
		regMethod := flagRegister(p.goType)
		defVal := flagDefault(p.goType)
		flagRegs = append(flagRegs, fmt.Sprintf("\t%s.Flags().%s(%q, %s, %q)\n",
			cmdVar, regMethod, flagName, defVal, p.name))

		if p.goType == "uint16" {
			fmt.Fprintf(&body, "\t\t%sRaw, _ := cmd.Flags().%s(%q)\n", p.name, getter, flagName)
			fmt.Fprintf(&body, "\t\t%s := uint16(%sRaw)\n", p.name, p.name)
		} else {
			fmt.Fprintf(&body, "\t\t%s, _ := cmd.Flags().%s(%q)\n", p.name, getter, flagName)
		}
	}

	// Construct receiver.
	if ti.constructor == nil {
		// No constructor available — can't read flags either.
		var stubBody bytes.Buffer
		stubBody.WriteString("\t\tfmt.Println(\"requires external context (NativeActivity, JNI, etc.)\")\n")
		stubBody.WriteString("\t\treturn nil\n")
		return stubBody.String(), nil
	}

	ctorArgs := ctorParamNames(ti.constructor.params)
	ctorHasError := len(ti.constructor.returns) > 1 && ti.constructor.returns[len(ti.constructor.returns)-1].goType == "error"

	if ctorHasError {
		fmt.Fprintf(&body, "\t\tobj, err := %s.%s(%s)\n", pkgRef, ti.constructor.name, ctorArgs)
		body.WriteString("\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n")
	} else {
		fmt.Fprintf(&body, "\t\tobj := %s.%s(%s)\n", pkgRef, ti.constructor.name, ctorArgs)
	}
	if ti.hasClose {
		body.WriteString("\t\tdefer obj.Close()\n")
	}

	// Call method.
	callArgs := paramNames(m.params)
	hasError := len(m.returns) > 0 && m.returns[len(m.returns)-1].goType == "error"
	hasResult := len(m.returns) > 0 && m.returns[0].goType != "error"

	writeCall(&body, "obj."+m.name, callArgs, hasResult, hasError, ctorHasError)

	return body.String(), flagRegs
}

func writeCall(
	buf *bytes.Buffer,
	callExpr string,
	callArgs string,
	hasResult bool,
	hasError bool,
	alreadyHasErr bool,
) {
	errAssign := "err :="
	if alreadyHasErr {
		errAssign = "err ="
	}

	switch {
	case hasResult && hasError:
		// Always use := since result is new. err gets re-declared (valid in Go).
		fmt.Fprintf(buf, "\t\tresult, err := %s(%s)\n", callExpr, callArgs)
		buf.WriteString("\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n")
		buf.WriteString("\t\tfmt.Println(result)\n")
	case hasResult:
		fmt.Fprintf(buf, "\t\tresult := %s(%s)\n", callExpr, callArgs)
		buf.WriteString("\t\tfmt.Println(result)\n")
	case hasError:
		fmt.Fprintf(buf, "\t\t%s %s(%s)\n", errAssign, callExpr, callArgs)
		buf.WriteString("\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n")
		buf.WriteString("\t\tfmt.Println(\"ok\")\n")
	default:
		fmt.Fprintf(buf, "\t\t%s(%s)\n", callExpr, callArgs)
		buf.WriteString("\t\tfmt.Println(\"ok\")\n")
	}
	buf.WriteString("\t\treturn nil\n")
}

func paramNames(params []paramInfo) string {
	var names []string
	for _, p := range params {
		names = append(names, p.name)
	}
	return strings.Join(names, ", ")
}

func ctorParamNames(params []paramInfo) string {
	var names []string
	for _, p := range params {
		names = append(names, "ctor"+camelCase(p.name))
	}
	return strings.Join(names, ", ")
}

func kebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			prev := rune(s[i-1])
			if unicode.IsLower(prev) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) {
				result = append(result, '-')
			}
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func camelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func uniqueVar(name string, used map[string]bool) string {
	if !used[name] {
		used[name] = true
		return name
	}
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s%d", name, i)
		if !used[candidate] {
			used[candidate] = true
			return candidate
		}
	}
}

func sanitizeVarPrefix(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

func sortedTypeNames(types map[string]*typeInfo) []string {
	var names []string
	for n := range types {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
