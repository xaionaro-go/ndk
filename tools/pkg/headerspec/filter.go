package headerspec

import "regexp"

// ApplyRules filters declarations using accept/ignore rules from the manifest.
//
// Rules are evaluated in order. For each declaration name:
//  1. Check all ignore rules first; if any match, the declaration is excluded.
//  2. Check all accept rules; if any match, the declaration is included.
//  3. If no rule matches, the declaration is excluded (default deny).
func ApplyRules(decls *Declarations, rules []Rule) *Declarations {
	if len(rules) == 0 {
		return decls
	}

	compiled := compileRules(rules)

	out := &Declarations{}

	for i := range decls.Functions {
		if shouldAccept(decls.Functions[i].Name, compiled) {
			out.Functions = append(out.Functions, decls.Functions[i])
		}
	}

	for i := range decls.Typedefs {
		if shouldAccept(decls.Typedefs[i].Name, compiled) {
			out.Typedefs = append(out.Typedefs, decls.Typedefs[i])
		}
	}

	for i := range decls.Enums {
		e := &decls.Enums[i]
		name := e.TypedefName
		if name == "" {
			name = e.Name
		}

		// For enums, also check if any of their constants match.
		if shouldAccept(name, compiled) || anyConstantAccepted(e.Constants, compiled) {
			out.Enums = append(out.Enums, *e)
		}
	}

	for i := range decls.Structs {
		if shouldAccept(decls.Structs[i].Name, compiled) {
			out.Structs = append(out.Structs, decls.Structs[i])
		}
	}

	return out
}

type compiledRule struct {
	action string
	re     *regexp.Regexp
}

func compileRules(rules []Rule) []compiledRule {
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile(r.From)
		if err != nil {
			continue
		}
		compiled = append(compiled, compiledRule{
			action: r.Action,
			re:     re,
		})
	}
	return compiled
}

// shouldAccept checks whether a given name is accepted by the rules.
// Ignore rules are checked first (any match rejects). Then accept rules
// are checked (any match accepts). No match means reject.
func shouldAccept(name string, rules []compiledRule) bool {
	if name == "" {
		return false
	}

	// First pass: check if any ignore rule matches.
	for _, r := range rules {
		if r.action == "ignore" && r.re.MatchString(name) {
			return false
		}
	}

	// Second pass: check if any accept rule matches.
	for _, r := range rules {
		if r.action == "accept" && r.re.MatchString(name) {
			return true
		}
	}

	return false
}

// anyConstantAccepted checks whether any enum constant name matches an accept
// rule (and no ignore rule).
func anyConstantAccepted(
	constants []EnumConstant,
	rules []compiledRule,
) bool {
	for _, c := range constants {
		if shouldAccept(c.Name, rules) {
			return true
		}
	}
	return false
}
