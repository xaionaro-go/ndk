package specmodel

// StructDef describes a C struct's fields extracted from NDK headers.
type StructDef struct {
	Fields []StructField `yaml:"fields"`
}
