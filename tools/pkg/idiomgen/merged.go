package idiomgen

import "github.com/xaionaro-go/ndk/tools/pkg/specmodel"

// MergedSpec is the fully resolved model passed to templates.
type MergedSpec struct {
	PackageName   string
	PackageImport string
	PackageDoc    string
	SourcePackage string

	OpaqueTypes           map[string]MergedOpaqueType
	TypeAliases           []MergedTypeAlias
	ErrorEnums            []MergedErrorEnum
	ValueEnums            []MergedValueEnum
	Methods               []MergedMethod
	FreeFunctions         []MergedFreeFunction
	Callbacks             []MergedCallback
	APILevels             map[string]int
	CallbackStructs       []MergedCallbackStruct
	CallbackStructAliases []MergedTypeAlias // type aliases re-exporting capi callback struct types
	StructAccessors       []MergedStructAccessor
	Lifecycle             *MergedLifecycle
	ExtraBridgeC          string
	ExtraBridgeGo         string
}

// MergedOpaqueType is an opaque C handle with Go lifecycle.
type MergedOpaqueType struct {
	GoName                    string
	CapiType                  string        // Name in capi package
	Constructor               string        // capi function name
	ConstructorReturnsPointer bool          // true if constructor returns pointer directly (not out-param)
	ConstructorParams         []MergedParam // additional params for constructor (excluding out-param)
	Destructor                string        // capi function name
	DestructorReturnsError    bool          // true if destructor returns an error type
	Pattern                   string        // builder, ref_counted, singleton
	Interfaces                []string
}

// MergedErrorEnum is an enum that implements the error interface.
type MergedErrorEnum struct {
	GoName       string // Original spec type name (used for capi reference)
	Prefix       string // Error string prefix
	SuccessValue string
	Values       []MergedEnumValue
}

// MergedValueEnum is a non-error enum with stripped prefixes.
type MergedValueEnum struct {
	GoName       string
	SpecName     string // Original spec type name
	BaseType     string // Go base type (e.g., "int32", "uint32"); defaults to "int32"
	StripPrefix  string
	StringMethod bool
	Values       []MergedEnumValue
}

// MergedEnumValue is one constant in a merged enum.
type MergedEnumValue struct {
	GoName   string
	SpecName string // Original C-style constant name
	Value    int64
	ValueStr string // Formatted value for template rendering (e.g., "0x80000000" for unsigned types)
}

// MergedMethod is a function resolved to a method on an opaque type.
type MergedMethod struct {
	GoName              string
	CName               string
	ReceiverType        string // Go name of the receiver opaque type
	Params              []MergedParam
	Returns             string
	Chain               bool
	Pure                bool
	ReturnsNew          string
	ReturnsNewDirect    bool // true if returns_new function returns pointer directly (not out-param)
	ReturnsFrames       bool
	BufGoType           string // Go type override for buffer param (e.g., "[]byte")
	TimeoutUnit         string // "ns" (default), "ms", "us", "s" — for time.Duration conversion
	APILevel            int
	FixedParams         map[string]string     // param name → literal Go value
	CallbackParam       string                // name of callback struct param (triggers bridge call)
	CallbackStruct      string                // spec name of the callback struct
	CallbackGoType      string                // Go name of the callback struct (e.g., "DeviceStateCallbacks")
	ReturnsListAccessor *MergedStructAccessor // if set, method returns a list via struct accessor
	CustomCall          *MergedCustomCall     // if set, method uses custom capi call args
	OutputParams        []MergedOutputParam   // C output params converted to Go return values
	ReturnsBool         bool                  // true when the C function returns bool (not a numeric status code)
}

// MergedParam is a resolved parameter.
type MergedParam struct {
	Name         string
	GoType       string
	CapiType     string // Original spec type (before typeMap resolution)
	IsHandle     bool   // True if param is an opaque handle needing .ptr unwrap
	IsString     bool   // True if param is *byte (C string) needing Go string conversion
	Remapped     bool   // True if the type was remapped through typeMap (needs cast even if names match)
	Direction    string
	DurationUnit string // "ns", "ms", "us", "s" — set when GoType is time.Duration
}

// MergedFreeFunction is a package-level function without a receiver (e.g., GL/EGL procedural APIs).
type MergedFreeFunction struct {
	GoName         string
	CName          string
	Params         []MergedParam
	Returns        string
	CapiReturns    string // Original spec return type (before typeMap resolution)
	IsHandleReturn bool   // True if return is an opaque handle type
	APILevel       int
	OutputParams   []MergedOutputParam // C output params converted to Go return values
	ReturnsBool    bool                // true when the C function returns bool (not a numeric status code)
	ReturnsNew     string              // Go type name for constructor-style functions
}

// MergedTypeAlias is a type alias re-exporting a capi type.
type MergedTypeAlias struct {
	GoName   string // Name in idiomatic package
	CapiType string // Name in capi package
}

// PerTypeData holds the data for rendering a single opaque type's file.
type PerTypeData struct {
	PackageName   string
	SourcePackage string
	Type          MergedOpaqueType
	Methods       []MergedMethod
	FreeFunctions []MergedFreeFunction        // factory functions that return this type
	OpaqueTypes   map[string]MergedOpaqueType // needed by lookupCapiType in templates
}

// MergedCallback is a resolved callback type.
type MergedCallback struct {
	SpecName       string
	GoCallbackType string
	GoCallbackSig  string
	Params         []MergedParam
	Returns        string
}

// MergedCallbackStruct is a C callback struct with Go type mapping.
type MergedCallbackStruct struct {
	SpecName     string // C struct name (e.g., "ACameraDevice_StateCallbacks")
	GoName       string // Go struct name (e.g., "DeviceStateCallbacks")
	ContextField string // Name of the void* context field
	Fields       []MergedCallbackField
}

// MergedCallbackField is one callback function pointer in a struct.
type MergedCallbackField struct {
	CName        string            // C field name (e.g., "onDisconnected")
	GoName       string            // Go field name (e.g., "OnDisconnected")
	GoSignature  string            // Go function signature (e.g., "func()")
	GoParamCount int               // number of params in Go callback signature
	Params       []specmodel.Param // C params from spec (for generating //export)
	Returns      string            // non-void C return type (e.g., "void*")
}

// MergedStructAccessor describes a C struct with list-like fields.
type MergedStructAccessor struct {
	SpecName   string // C struct name (e.g., "ACameraIdList")
	CountField string // Field name for count (e.g., "numCameras")
	ItemsField string // Field name for items (e.g., "cameraIds")
	ItemType   string // Go type of each item (e.g., "string")
	DeleteFunc string // capi function to free the list
}

// MergedOutputParam describes a C output parameter converted to a Go return value.
type MergedOutputParam struct {
	CParamName string // Original C param name (e.g., "reader")
	GoType     string // Go idiomatic type (e.g., "*ImageReader")
	CapiType   string // Capi pointer type to declare (e.g., "*capi.AImageReader")
	IsHandle   bool   // True if the Go type is an opaque handle wrapper
}

// MergedCustomCall holds custom call information for methods with non-standard capi patterns.
type MergedCustomCall struct {
	Params []MergedCustomCallParam // Go visible params
	Args   string                  // literal capi args after receiver
}

// MergedCustomCallParam is a param in a custom call method's Go signature.
type MergedCustomCallParam struct {
	Name   string
	GoType string
}

// MergedLifecycle describes a NativeActivity lifecycle pattern.
type MergedLifecycle struct {
	EntryPoint        string                // Exported symbol (e.g., "ANativeActivity_onCreate")
	ActivityType      string                // C type name (e.g., "ANativeActivity")
	GoActivityType    string                // Go type name (e.g., "Activity")
	CallbacksAccessor string                // C expression (e.g., "activity->callbacks")
	CallbackStruct    string                // C struct name (e.g., "ANativeActivityCallbacks")
	EntryCallback     string                // Go field name for entry point callback (e.g., "OnCreate")
	Fields            []MergedCallbackField // Lifecycle event callbacks
}

// PerValueEnumData holds data for rendering a single value enum file.
type PerValueEnumData struct {
	PackageName string
	Enum        MergedValueEnum
}

// PerTypeAliasData holds data for rendering a single type alias file.
type PerTypeAliasData struct {
	PackageName   string
	SourcePackage string
	Alias         MergedTypeAlias
}

// PerCallbackData holds data for rendering a single callback type file.
type PerCallbackData struct {
	PackageName string
	Callback    MergedCallback
}

// BridgeData holds data for rendering bridge templates into a capi package.
type BridgeData struct {
	PackageName     string
	CIncludes       []string
	CallbackStructs []MergedCallbackStruct
	StructAccessors []MergedStructAccessor
	Lifecycle       *MergedLifecycle
	ExtraBridgeC    string
	ExtraBridgeGo   string
}
