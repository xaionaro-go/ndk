package ndk

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPackageName() string {
	smapsName := fmt.Sprint("/proc/", os.Getpid(), "/cmdline")
	b, err := os.ReadFile(smapsName)
	if err != nil {
		return ""
	}
	i := bytes.Index(b, []byte{0})
	if i >= 0 {
		b = b[:i]
	}
	return string(b)
}

func GetLibraryPath() string {
	smapsName := fmt.Sprint("/proc/", os.Getpid(), "/smaps")
	b, err := os.ReadFile(smapsName)
	if err != nil {
		return ""
	}

	for {
		base := []byte("/data/")
		i := bytes.Index(b, base)
		if i < 0 {
			break
		}
		j := bytes.Index(b[i:], []byte("\n"))
		if j < 0 {
			break
		}
		path := string(b[i : i+j])
		if strings.HasSuffix(path, ".so") {
			return filepath.Dir(path)
		}

		b = b[i+j:]
	}
	return ""
}

func FindMatchLibrary(pattern string) []string {
	paths := os.Getenv("LD_LIBRARY_PATH")
	if paths == "" {
		paths = "/system/lib"
	}
	dirs := strings.Split(paths, ":")
	dirs = append([]string{GetLibraryPath()}, dirs...)
	for _, dir := range dirs {
		fns, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			Info("FindMatchLibrary:", err)
		}
		if len(fns) > 0 {
			return fns
		}
	}
	return nil
}
