// Package apkbuild provides functionality for building signed Android APKs
// from Go packages using the Android SDK/NDK toolchain.
package apkbuild

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for building an APK.
type Config struct {
	// GoPackage is the Go package to build (e.g., "./examples/camera/display").
	GoPackage string

	// OutputPath is the output APK file path.
	OutputPath string

	// ManifestPath is the path to AndroidManifest.xml.
	ManifestPath string

	// Architectures to build for (e.g., ["arm64", "x86_64"]).
	Architectures []string

	// APILevel is the Android API level (e.g., 35).
	APILevel int

	// KeystorePath is the path to the signing keystore.
	KeystorePath string

	// KeystorePass is the keystore password.
	KeystorePass string

	// LibName is the name for the shared library (without lib prefix and .so suffix).
	LibName string

	// AndroidHome is the path to the Android SDK (default: $ANDROID_HOME or ~/Android/Sdk).
	AndroidHome string

	// NDKPath is the path to the Android NDK (default: latest in AndroidHome/ndk/).
	NDKPath string
}

func (cfg *Config) validate() error {
	if cfg.GoPackage == "" {
		return fmt.Errorf("go package is required")
	}
	if cfg.ManifestPath == "" {
		return fmt.Errorf("manifest path is required")
	}
	return nil
}

func (cfg *Config) setDefaults() {
	if cfg.AndroidHome == "" {
		cfg.AndroidHome = os.Getenv("ANDROID_HOME")
		if cfg.AndroidHome == "" {
			home, _ := os.UserHomeDir()
			cfg.AndroidHome = filepath.Join(home, "Android", "Sdk")
		}
	}

	if cfg.NDKPath == "" {
		cfg.NDKPath = os.Getenv("ANDROID_NDK_HOME")
		if cfg.NDKPath == "" {
			cfg.NDKPath = findLatestDir(filepath.Join(cfg.AndroidHome, "ndk"))
		}
	}

	if cfg.APILevel == 0 {
		cfg.APILevel = 35
	}

	if len(cfg.Architectures) == 0 {
		cfg.Architectures = []string{"arm64"}
	}

	if cfg.LibName == "" {
		cfg.LibName = filepath.Base(cfg.GoPackage)
	}

	if cfg.KeystorePath == "" {
		home, _ := os.UserHomeDir()
		cfg.KeystorePath = filepath.Join(home, ".android", "debug.keystore")
	}

	if cfg.KeystorePass == "" {
		cfg.KeystorePass = "android"
	}

	if cfg.OutputPath == "" {
		cfg.OutputPath = cfg.LibName + ".apk"
	}
}

func (cfg *Config) findBuildTools() (string, error) {
	btDir := filepath.Join(cfg.AndroidHome, "build-tools")
	latest := findLatestDir(btDir)
	if latest == "" {
		return "", fmt.Errorf("no build-tools found in %s", btDir)
	}
	return latest, nil
}
