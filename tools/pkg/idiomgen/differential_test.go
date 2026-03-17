package idiomgen

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const leanOraclePath = "../../../proofs/.lake/build/bin/proofs"

func leanOracleAvailable() bool {
	_, err := os.Stat(leanOraclePath)
	return err == nil
}

func runLeanOracle(t *testing.T, commands []string) []string {
	t.Helper()
	input := strings.Join(commands, "\n") + "\n"
	cmd := exec.Command(leanOraclePath)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("lean oracle failed: %v", err)
	}
	var results []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}
	return results
}

// TestDifferential_CommonPrefix compares Go and Lean commonPrefix.
func TestDifferential_CommonPrefix(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	type pair struct{ a, b string }
	inputs := []pair{
		{"", ""},
		{"", "hello"},
		{"hello", ""},
		{"hello", "help"},
		{"abc", "abd"},
		{"abc", "xyz"},
		{"test", "test"},
		{"ab", "abcde"},
		{"AAUDIO_DIRECTION_OUTPUT", "AAUDIO_DIRECTION_INPUT"},
		{"ASensor_create", "ASensor_delete"},
	}

	var commands []string
	for _, p := range inputs {
		commands = append(commands, fmt.Sprintf("commonPrefix\t%s\t%s", p.a, p.b))
	}

	leanResults := runLeanOracle(t, commands)
	if len(leanResults) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(leanResults))
	}

	for i, p := range inputs {
		goResult := commonPrefix(p.a, p.b)
		leanResult := leanResults[i]
		if goResult != leanResult {
			t.Errorf("commonPrefix(%q, %q): Go=%q, Lean=%q", p.a, p.b, goResult, leanResult)
		}
	}
}

// commonPrefix is the Go implementation matching c2ffi.commonPrefix.
func commonPrefix(a, b string) string {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:n]
}

// TestDifferential_IsConstructorDestructorFunc compares Go and Lean classification.
func TestDifferential_IsConstructorDestructorFunc(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	type testCase struct {
		funcName, specTypeName string
	}
	inputs := []testCase{
		{"ASensor_create", "ASensor"},
		{"ASensor_new", "ASensor"},
		{"ASensor_delete", "ASensor"},
		{"ASensor_free", "ASensor"},
		{"ASensor_destroy", "ASensor"},
		{"ASensor_release", "ASensor"},
		{"ASensor_close", "ASensor"},
		{"ASensor_getName", "ASensor"},
		{"ALooper_create", "ALooper"},
		{"unrelated", "ALooper"},
	}

	var ctorCommands, dtorCommands []string
	for _, tc := range inputs {
		ctorCommands = append(ctorCommands, fmt.Sprintf("isConstructorFunc\t%s\t%s", tc.funcName, tc.specTypeName))
		dtorCommands = append(dtorCommands, fmt.Sprintf("isDestructorFunc\t%s\t%s", tc.funcName, tc.specTypeName))
	}

	ctorResults := runLeanOracle(t, ctorCommands)
	dtorResults := runLeanOracle(t, dtorCommands)

	for i, tc := range inputs {
		goCtorResult := fmt.Sprintf("%v", isConstructorFunc(tc.funcName, tc.specTypeName))
		goDtorResult := fmt.Sprintf("%v", isDestructorFunc(tc.funcName, tc.specTypeName))

		if goCtorResult != ctorResults[i] {
			t.Errorf("isConstructorFunc(%q, %q): Go=%s, Lean=%s", tc.funcName, tc.specTypeName, goCtorResult, ctorResults[i])
		}
		if goDtorResult != dtorResults[i] {
			t.Errorf("isDestructorFunc(%q, %q): Go=%s, Lean=%s", tc.funcName, tc.specTypeName, goDtorResult, dtorResults[i])
		}
	}
}

// TestDifferential_CountGoParams compares Go and Lean countGoParams.
func TestDifferential_CountGoParams(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	inputs := []string{
		"func()",
		"func(int)",
		"func(int, string)",
		"func(int, string, bool)",
		"func( )",
		"hello",
		"",
	}

	var commands []string
	for _, sig := range inputs {
		commands = append(commands, fmt.Sprintf("countGoParams\t%s", sig))
	}

	leanResults := runLeanOracle(t, commands)
	if len(leanResults) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(leanResults))
	}

	for i, sig := range inputs {
		goResult := fmt.Sprintf("%d", countGoParams(sig))
		leanResult := leanResults[i]
		if goResult != leanResult {
			t.Errorf("countGoParams(%q): Go=%s, Lean=%s", sig, goResult, leanResult)
		}
	}
}

// TestDifferential_ResolveType compares Go and Lean resolveType.
func TestDifferential_ResolveType(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	type testCase struct {
		specType string
		mapStr   string // comma-separated key=value for Lean
		goMap    map[string]string
	}
	inputs := []testCase{
		{"foo", "foo=Bar", map[string]string{"foo": "Bar"}},
		{"*foo", "foo=Bar", map[string]string{"foo": "Bar"}},
		{"**foo", "foo=Bar", map[string]string{"foo": "Bar"}},
		{"[]foo", "foo=Bar", map[string]string{"foo": "Bar"}},
		{"unknown", "foo=Bar", map[string]string{"foo": "Bar"}},
		{"*foo", "*foo=Direct,foo=Prefix", map[string]string{"*foo": "Direct", "foo": "Prefix"}},
		{"baz", "", nil},
	}

	var commands []string
	for _, tc := range inputs {
		commands = append(commands, fmt.Sprintf("resolveType\t%s\t%s", tc.specType, tc.mapStr))
	}

	leanResults := runLeanOracle(t, commands)
	if len(leanResults) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(leanResults))
	}

	for i, tc := range inputs {
		goResult := resolveType(tc.specType, tc.goMap)
		leanResult := leanResults[i]
		if goResult != leanResult {
			t.Errorf("resolveType(%q, %v): Go=%q, Lean=%q", tc.specType, tc.goMap, goResult, leanResult)
		}
	}
}
