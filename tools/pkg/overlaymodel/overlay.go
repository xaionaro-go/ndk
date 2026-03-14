// Package overlaymodel defines the YAML overlay schema for hand-written semantic annotations.
package overlaymodel

// Overlay is the top-level structure for one module's overlay file.
type Overlay struct {
	Module          string                           `yaml:"module"`
	Package         PackageOverlay                   `yaml:"package"`
	Types           map[string]TypeOverlay           `yaml:"types,omitempty"`
	Functions       map[string]FuncOverlay           `yaml:"functions,omitempty"`
	APILevels       map[string]int                   `yaml:"api_levels,omitempty"`
	CallbackStructs map[string]CallbackStructOverlay `yaml:"callback_structs,omitempty"`
	StructAccessors map[string]StructAccessorOverlay `yaml:"struct_accessors,omitempty"`
	Lifecycle       *LifecycleOverlay                `yaml:"lifecycle,omitempty"`
	ExtraEnums      map[string][]ExtraEnumValue      `yaml:"extra_enums,omitempty"`
	ExtraBridgeC    string                           `yaml:"extra_bridge_c,omitempty"`
	ExtraBridgeGo   string                           `yaml:"extra_bridge_go,omitempty"`
}

// ExtraEnumValue defines an additional enum constant not present in the
// auto-generated spec (e.g. from extension headers not parsed by specgen).
type ExtraEnumValue struct {
	Name  string `yaml:"name"`
	Value int64  `yaml:"value"`
}

// PackageOverlay configures the idiomatic Go package.
type PackageOverlay struct {
	GoName   string `yaml:"go_name"`
	GoImport string `yaml:"go_import"`
	Doc      string `yaml:"doc"`
}

// TypeOverlay provides idiomatic annotations for a type.
type TypeOverlay struct {
	GoName       string   `yaml:"go_name,omitempty"`
	GoError      bool     `yaml:"go_error,omitempty"`
	SuccessValue string   `yaml:"success_value,omitempty"`
	ErrorPrefix  string   `yaml:"error_prefix,omitempty"`
	StripPrefix  string   `yaml:"strip_prefix,omitempty"`
	StringMethod bool     `yaml:"string_method,omitempty"`
	Constructor  string   `yaml:"constructor,omitempty"`
	Destructor   string   `yaml:"destructor,omitempty"`
	Pattern      string   `yaml:"pattern,omitempty"` // builder, ref_counted, singleton
	Interfaces   []string `yaml:"interfaces,omitempty"`
	Transparent  bool     `yaml:"transparent,omitempty"` // type alias instead of struct wrapper
	EnumSource   string   `yaml:"enum_source,omitempty"` // spec enum key when typedef name differs from enum group name
}

// FuncOverlay provides idiomatic annotations for a function.
type FuncOverlay struct {
	Receiver          string             `yaml:"receiver,omitempty"`
	GoName            string             `yaml:"go_name,omitempty"`
	Skip              bool               `yaml:"skip,omitempty"`
	Chain             bool               `yaml:"chain,omitempty"`
	Pure              bool               `yaml:"pure,omitempty"`
	ReturnsNew        string             `yaml:"returns_new,omitempty"`
	CallbackParam     string             `yaml:"callback_param,omitempty"`
	UserdataParam     string             `yaml:"userdata_param,omitempty"`
	GoCallbackType    string             `yaml:"go_callback_type,omitempty"`
	GoCallbackSig     string             `yaml:"go_callback_sig,omitempty"`
	BufParam          string             `yaml:"buf_param,omitempty"`
	BufFramesParam    string             `yaml:"buf_frames_param,omitempty"`
	BufGoType         string             `yaml:"buf_go_type,omitempty"` // Go type for buffer param (e.g., "[]byte")
	TimeoutParam      string             `yaml:"timeout_param,omitempty"`
	TimeoutUnit       string             `yaml:"timeout_unit,omitempty"` // "ns" (default), "ms", "us", "s"
	ReturnsFrames     bool               `yaml:"returns_frames,omitempty"`
	FixedParams       map[string]string  `yaml:"fixed_params,omitempty"`
	ReturnsListStruct string             `yaml:"returns_list_struct,omitempty"` // spec name of list struct for iteration
	CustomCall        *CustomCallOverlay `yaml:"custom_call,omitempty"`
	BridgeParams      []BridgeParam      `yaml:"bridge_params,omitempty"`  // params for bridge functions not in spec
	BridgeReturns     string             `yaml:"bridge_returns,omitempty"` // return type for bridge functions not in spec
}

// BridgeParam is a parameter for a bridge function defined only in the overlay.
type BridgeParam struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// CallbackStructOverlay configures Go bindings for a C callback struct.
type CallbackStructOverlay struct {
	GoName       string                          `yaml:"go_name"`
	ContextField string                          `yaml:"context_field"`
	Fields       map[string]CallbackFieldOverlay `yaml:"fields"`
}

// CallbackFieldOverlay configures one field of a callback struct.
type CallbackFieldOverlay struct {
	GoName       string `yaml:"go_name"`
	GoSignature  string `yaml:"go_signature"`
	CallbackType string `yaml:"callback_type,omitempty"`
	GoParamCount int    `yaml:"go_param_count,omitempty"`
}

// StructAccessorOverlay configures Go accessors for a C struct with list-like fields.
type StructAccessorOverlay struct {
	CountField string `yaml:"count_field"`
	ItemsField string `yaml:"items_field"`
	ItemType   string `yaml:"item_type"`
	DeleteFunc string `yaml:"delete_func"`
}

// CustomCallOverlay specifies a method with custom capi call arguments.
type CustomCallOverlay struct {
	Params []CustomCallParam `yaml:"params"`
	Args   string            `yaml:"args"` // literal capi args after receiver (e.g., "nil, 1, &req.ptr, nil")
}

// CustomCallParam is a param in a custom call method's Go signature.
type CustomCallParam struct {
	Name   string `yaml:"name"`
	GoType string `yaml:"go_type"`
}

// LifecycleOverlay configures NativeActivity lifecycle callback generation.
type LifecycleOverlay struct {
	EntryPoint        string `yaml:"entry_point"`
	ActivityType      string `yaml:"activity_type"`
	CallbacksAccessor string `yaml:"callbacks_accessor"`
	CallbackStruct    string `yaml:"callback_struct"`
	EntryCallback     string `yaml:"entry_callback"` // Go field name for entry point callback (e.g., "OnCreate")
}
