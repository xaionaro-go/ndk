// Package headerspec parses clang JSON AST output and generates spec YAML.
package headerspec

import "encoding/json"

// ASTNode represents a node in clang's JSON AST dump.
type ASTNode struct {
	ID                  string    `json:"id"`
	Kind                string    `json:"kind"`
	Name                string    `json:"name,omitempty"`
	MangledName         string    `json:"mangledName,omitempty"`
	Type                *TypeInfo `json:"type,omitempty"`
	Inner               []ASTNode `json:"inner,omitempty"`
	Loc                 *Location `json:"loc,omitempty"`
	Range               *Range    `json:"range,omitempty"`
	IsImplicit          bool      `json:"isImplicit,omitempty"`
	IsReferenced        bool      `json:"isReferenced,omitempty"`
	CompleteDefinition  bool      `json:"completeDefinition,omitempty"`
	TagUsed             string    `json:"tagUsed,omitempty"`
	FixedUnderlyingType *TypeInfo `json:"fixedUnderlyingType,omitempty"`
	StorageClass        string    `json:"storageClass,omitempty"`
	Value               string    `json:"value,omitempty"`
	CastKind            string    `json:"castKind,omitempty"`
	CC                  string    `json:"cc,omitempty"`
	OwnedTagDecl        *ASTNode  `json:"ownedTagDecl,omitempty"`
	Decl                *ASTNode  `json:"decl,omitempty"`
	ValueCategory       string    `json:"valueCategory,omitempty"`
}

// TypeInfo holds type qualification data from a clang AST node.
type TypeInfo struct {
	QualType          string `json:"qualType"`
	DesugaredQualType string `json:"desugaredQualType,omitempty"`
	TypeAliasDeclID   string `json:"typeAliasDeclId,omitempty"`
}

// Location holds source location information from a clang AST node.
type Location struct {
	File         string    `json:"file,omitempty"`
	Line         int       `json:"line,omitempty"`
	Col          int       `json:"col,omitempty"`
	Offset       int       `json:"offset,omitempty"`
	TokLen       int       `json:"tokLen,omitempty"`
	IncludedFrom *IncluRef `json:"includedFrom,omitempty"`
	SpellingLoc  *Location `json:"spellingLoc,omitempty"`
	ExpansionLoc *Location `json:"expansionLoc,omitempty"`
}

// IncluRef references the file that included this location.
type IncluRef struct {
	File string `json:"file,omitempty"`
}

// Range holds a source range (begin/end locations).
type Range struct {
	Begin *Location `json:"begin,omitempty"`
	End   *Location `json:"end,omitempty"`
}

// ParseClangAST unmarshals clang's -ast-dump=json output into an ASTNode tree.
func ParseClangAST(jsonData []byte) (*ASTNode, error) {
	var root ASTNode
	if err := json.Unmarshal(jsonData, &root); err != nil {
		return nil, err
	}
	return &root, nil
}
