// Package c2ffi parses c2ffi JSON output and converts it to specmodel.Spec.
package c2ffi

// Declaration is one top-level entry in c2ffi's JSON output array.
// c2ffi outputs a flat JSON array of tagged objects — functions, typedefs,
// enums, structs, etc.
type Declaration struct {
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	NS       int    `json:"ns"`
	ID       int    `json:"id"`
	Location string `json:"location"`

	// Function fields.
	Variadic   bool        `json:"variadic"`
	Inline     bool        `json:"inline"`
	Parameters []Parameter `json:"parameters"`
	ReturnType *TypeRef    `json:"return-type"`

	// Typedef field.
	Type *TypeRef `json:"type"`

	// Enum and struct fields share the same JSON key "fields".
	// Enum fields have Name+Value; struct fields have Name+Type+BitOffset.
	Fields []Field `json:"fields"`

	// Struct-specific.
	BitSize      int `json:"bit-size"`
	BitAlignment int `json:"bit-alignment"`
}

// Parameter is a function parameter in c2ffi output.
type Parameter struct {
	Tag  string  `json:"tag"`
	Name string  `json:"name"`
	Type TypeRef `json:"type"`
}

// TypeRef describes a type in c2ffi's JSON. Types are recursive:
// a pointer is {"tag": ":pointer", "type": {...}}.
type TypeRef struct {
	Tag          string   `json:"tag"`
	Type         *TypeRef `json:"type"`
	BitSize      int      `json:"bit-size"`
	BitAlignment int      `json:"bit-alignment"`

	// Struct/enum inline (for typedef targets like ":enum").
	Name   string  `json:"name"`
	ID     int     `json:"id"`
	Fields []Field `json:"fields"`

	// Array.
	Size int `json:"size"`
}

// Field is used for both enum constants and struct fields.
// c2ffi reuses the "fields" JSON key for both.
type Field struct {
	Tag          string   `json:"tag"`
	Name         string   `json:"name"`
	Value        uint64   `json:"value"`         // enum constant value
	BitOffset    int      `json:"bit-offset"`    // struct field offset
	BitSize      int      `json:"bit-size"`      // struct field size
	BitAlignment int      `json:"bit-alignment"` // struct field alignment
	Type         *TypeRef `json:"type"`          // struct field type
}
