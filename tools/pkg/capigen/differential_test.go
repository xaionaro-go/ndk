package capigen

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// leanOraclePath is the path to the Lean oracle executable built by `lake build`.
const leanOraclePath = "../../../proofs/.lake/build/bin/proofs"

func leanOracleAvailable() bool {
	_, err := os.Stat(leanOraclePath)
	return err == nil
}

// runLeanOracle sends commands to the Lean oracle and returns the outputs.
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

// TestDifferential_ExportName runs the same inputs through Go ExportName and
// Lean exportName, comparing outputs.
func TestDifferential_ExportName(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	inputs := []string{
		"", "_", "__", "___",
		"foo", "Foo", "FOO",
		"_foo", "__bar", "___baz",
		"ALooper", "AImageReader", "camera_status_t",
		"a", "A", "z", "Z", "0start",
	}

	var commands []string
	for _, in_ := range inputs {
		commands = append(commands, fmt.Sprintf("exportName\t%s", in_))
	}

	leanResults := runLeanOracle(t, commands)
	if len(leanResults) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(leanResults))
	}

	for i, in_ := range inputs {
		goResult := ExportName(in_)
		leanResult := leanResults[i]
		if goResult != leanResult {
			t.Errorf("ExportName(%q): Go=%q, Lean=%q", in_, goResult, leanResult)
		}
	}
}

// TestDifferential_StripPointers runs the same inputs through Go and Lean.
func TestDifferential_StripPointers(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	inputs := []string{
		"", "*", "**", "***",
		"ALooper", "*ALooper", "**ALooper", "***int32",
		"int32", "*int32",
	}

	var commands []string
	for _, in_ := range inputs {
		commands = append(commands, fmt.Sprintf("stripPointers\t%s", in_))
	}

	leanResults := runLeanOracle(t, commands)
	if len(leanResults) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(leanResults))
	}

	for i, in_ := range inputs {
		goResult := stripPointers(in_)
		leanResult := leanResults[i]
		if goResult != leanResult {
			t.Errorf("stripPointers(%q): Go=%q, Lean=%q", in_, goResult, leanResult)
		}
	}
}

// TestDifferential_IsScalarGoType compares Go and Lean classification.
func TestDifferential_IsScalarGoType(t *testing.T) {
	if !leanOracleAvailable() {
		t.Skip("Lean oracle not built (run 'cd proofs && lake build')")
	}

	inputs := []string{
		"int8", "uint8", "int16", "uint16", "int32", "uint32",
		"int64", "uint64", "float32", "float64", "bool", "int", "uint",
		"string", "unsafe.Pointer", "ALooper", "", "*int32",
	}

	var commands []string
	for _, in_ := range inputs {
		commands = append(commands, fmt.Sprintf("isScalarGoType\t%s", in_))
	}

	leanResults := runLeanOracle(t, commands)
	if len(leanResults) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(leanResults))
	}

	for i, in_ := range inputs {
		goResult := fmt.Sprintf("%v", isScalarGoType(in_))
		leanResult := leanResults[i]
		if goResult != leanResult {
			t.Errorf("isScalarGoType(%q): Go=%s, Lean=%s", in_, goResult, leanResult)
		}
	}
}
