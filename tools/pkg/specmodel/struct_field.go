package specmodel

// StructField is one field in a C struct.
type StructField struct {
	Name    string  `yaml:"name"`
	Type    string  `yaml:"type"`
	Params  []Param `yaml:"params,omitempty"`
	Returns string  `yaml:"returns,omitempty"` // non-void return type for func_ptr fields
}
