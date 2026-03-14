package apkbuild

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Build produces a signed APK from a Go package.
//
// The pipeline mirrors what the Makefile does:
//  1. Cross-compile Go → .so (go build -buildmode=c-shared)
//  2. Package with aapt (aapt package -f -M manifest -I platform.jar -F base.apk)
//  3. Add native libs (cd buildDir && zip -r base.apk lib/)
//  4. Align (zipalign -f -p 4)
//  5. Sign (apksigner sign --ks keystore)
func Build(cfg Config) error {
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	cfg.setDefaults()

	buildTools, err := cfg.findBuildTools()
	if err != nil {
		return err
	}

	platformJar := filepath.Join(
		cfg.AndroidHome, "platforms",
		fmt.Sprintf("android-%d", cfg.APILevel), "android.jar",
	)
	if _, err := os.Stat(platformJar); err != nil {
		return fmt.Errorf(
			"platform jar not found: %s (install Android SDK platform %d)",
			platformJar, cfg.APILevel,
		)
	}

	buildDir, err := os.MkdirTemp("", "ndk-build-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(buildDir) }()

	for _, arch := range cfg.Architectures {
		if err := crossCompile(&cfg, buildDir, arch); err != nil {
			return fmt.Errorf("cross-compile %s: %w", arch, err)
		}
	}

	baseAPK := filepath.Join(buildDir, "base.apk")
	if err := runCmd(
		filepath.Join(buildTools, "aapt"),
		"package", "-f",
		"-M", cfg.ManifestPath,
		"-I", platformJar,
		"-F", baseAPK,
	); err != nil {
		return fmt.Errorf("aapt package: %w", err)
	}

	if err := runCmdInDir(buildDir, "zip", "-r", "base.apk", "lib/"); err != nil {
		return fmt.Errorf("adding native libs: %w", err)
	}

	alignedAPK := filepath.Join(buildDir, "aligned.apk")
	if err := runCmd(
		filepath.Join(buildTools, "zipalign"),
		"-f", "-p", "4", baseAPK, alignedAPK,
	); err != nil {
		return fmt.Errorf("zipalign: %w", err)
	}

	if err := ensureKeystore(&cfg); err != nil {
		return fmt.Errorf("keystore: %w", err)
	}

	if err := runCmd(
		filepath.Join(buildTools, "apksigner"),
		"sign",
		"--ks", cfg.KeystorePath,
		"--ks-pass", "pass:"+cfg.KeystorePass,
		"--key-pass", "pass:"+cfg.KeystorePass,
		alignedAPK,
	); err != nil {
		return fmt.Errorf("apksigner: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.OutputPath), 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	if err := copyFile(alignedAPK, cfg.OutputPath); err != nil {
		return fmt.Errorf("copying APK: %w", err)
	}

	return nil
}

func crossCompile(cfg *Config, buildDir, arch string) error {
	goarch, abi, ccPrefix := archMapping(arch, cfg.APILevel)
	if goarch == "" {
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	cc := filepath.Join(
		cfg.NDKPath, "toolchains", "llvm", "prebuilt",
		"linux-x86_64", "bin", ccPrefix+"-clang",
	)

	libDir := filepath.Join(buildDir, "lib", abi)
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		return err
	}

	soPath := filepath.Join(libDir, "lib"+cfg.LibName+".so")

	cmd := exec.Command(
		"go", "build", "-buildmode=c-shared",
		"-o", soPath, cfg.GoPackage,
	)
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"GOOS=android",
		"GOARCH="+goarch,
		"CC="+cc,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Remove the generated .h file that go build -buildmode=c-shared produces.
	_ = os.Remove(strings.TrimSuffix(soPath, ".so") + ".h")
	return nil
}

func ensureKeystore(cfg *Config) error {
	if _, err := os.Stat(cfg.KeystorePath); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(cfg.KeystorePath), 0o755); err != nil {
		return err
	}

	return runCmd("keytool",
		"-genkeypair", "-v",
		"-keystore", cfg.KeystorePath,
		"-storepass", cfg.KeystorePass,
		"-alias", "androiddebugkey",
		"-keypass", cfg.KeystorePass,
		"-keyalg", "RSA",
		"-keysize", "2048",
		"-validity", "10000",
		"-dname", "CN=Debug,O=Debug,C=US",
	)
}
