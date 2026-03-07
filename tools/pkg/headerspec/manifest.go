package headerspec

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Manifest is the top-level manifest YAML used by capigen/headerspec.
type Manifest struct {
	Generator  GeneratorConfig  `yaml:"GENERATOR"`
	Parser     ParserConfig     `yaml:"PARSER"`
	Translator TranslatorConfig `yaml:"TRANSLATOR"`
}

// GeneratorConfig holds package generation settings.
type GeneratorConfig struct {
	PackageName        string      `yaml:"PackageName"`
	PackageDescription string      `yaml:"PackageDescription"`
	PackageLicense     string      `yaml:"PackageLicense"`
	Includes           []string    `yaml:"Includes"`
	FlagGroups         []FlagGroup `yaml:"FlagGroups"`
}

// FlagGroup is a named set of compiler/linker flags.
type FlagGroup struct {
	Name  string   `yaml:"name"`
	Flags []string `yaml:"flags"`
}

// ParserConfig holds parser-specific settings (used by c-for-go, kept for compatibility).
type ParserConfig struct {
	IncludePaths []string `yaml:"IncludePaths"`
	SourcesPaths []string `yaml:"SourcesPaths"`
}

// TranslatorConfig holds accept/ignore rules and const evaluation settings.
type TranslatorConfig struct {
	ConstRules map[string]string `yaml:"ConstRules"`
	Rules      map[string][]Rule `yaml:"Rules"`
}

// Rule is a single accept/ignore rule with a regex pattern.
type Rule struct {
	Action string `yaml:"action"`
	From   string `yaml:"from"`
	To     string `yaml:"to,omitempty"`
}

// ParseManifest reads and unmarshals a manifest YAML file.
func ParseManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest %s: %w", path, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %s: %w", path, err)
	}

	return &m, nil
}
