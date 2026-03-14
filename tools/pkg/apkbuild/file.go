package apkbuild

import "os"

// findLatestDir returns the path to the last directory entry (sorted by name)
// within the given parent directory. This works for version-numbered directories
// like build-tools/35.0.0 or ndk/27.2.12479018.
func findLatestDir(parent string) string {
	entries, err := os.ReadDir(parent)
	if err != nil || len(entries) == 0 {
		return ""
	}

	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].IsDir() {
			return parent + "/" + entries[i].Name()
		}
	}
	return ""
}

// copyFile reads src and writes its contents to dst with mode 0o644.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
