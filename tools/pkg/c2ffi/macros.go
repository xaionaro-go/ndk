package c2ffi

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// macroRe matches "#define NAME VALUE" where VALUE is a hex or decimal integer
// literal, optionally with C integer suffixes (u, l, ul, ull, etc.).
var macroRe = regexp.MustCompile(
	`^#define\s+([A-Z][A-Z0-9_]+)\s+(0[xX][0-9a-fA-F]+|[0-9]+)[uUlL]*\s*$`,
)

// ExtractMacros parses C header files for #define integer constants.
// Only macros from headers matching targetHeaders (suffix match) are included.
func ExtractMacros(
	ndkSysroot string,
	targetHeaders []string,
) (map[string]int64, error) {
	macros := make(map[string]int64)

	for _, header := range targetHeaders {
		path := ndkSysroot + "/" + header
		if err := extractMacrosFromFile(path, macros); err != nil {
			return nil, fmt.Errorf("extracting macros from %s: %w", header, err)
		}
	}

	return macros, nil
}

func extractMacrosFromFile(
	path string,
	macros map[string]int64,
) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#define") {
			continue
		}

		m := macroRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		name := m[1]
		valueStr := m[2]

		// Skip version/feature-test macros (e.g. EGL_VERSION_1_0, GL_ES_VERSION_2_0).
		if strings.Contains(name, "_VERSION_") {
			continue
		}
		// Skip prototype guard macros (e.g. EGL_EGL_PROTOTYPES, GL_GLES_PROTOTYPES).
		if strings.HasSuffix(name, "_PROTOTYPES") {
			continue
		}

		value, err := strconv.ParseInt(valueStr, 0, 64)
		if err != nil {
			// Handle unsigned values that overflow int64
			// (e.g. GL_TIMEOUT_IGNORED = 0xFFFFFFFFFFFFFFFF).
			uval, uerr := strconv.ParseUint(valueStr, 0, 64)
			if uerr != nil {
				continue
			}
			value = int64(uval)
		}

		macros[name] = value
	}

	return scanner.Err()
}
