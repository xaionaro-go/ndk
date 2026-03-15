// Command ndk-build produces a signed Android APK from a Go package
// using the Android SDK/NDK toolchain.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/AndroidGoLab/ndk/tools/pkg/apkbuild"
)

func main() {
	var (
		output   = flag.String("o", "", "Output APK path (default: <lib-name>.apk)")
		manifest = flag.String("manifest", "AndroidManifest.xml", "Path to AndroidManifest.xml")
		arch     = flag.String("arch", "arm64", "Target architectures, comma-separated (arm64,x86_64,arm,x86)")
		apiLevel = flag.Int("api", 35, "Android API level")
		keystore = flag.String("keystore", "", "Keystore for signing (default: ~/.android/debug.keystore)")
		libName  = flag.String("name", "", "Library name (default: derived from package)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go-ndk-build [flags] <go-package>\n\n")
		fmt.Fprintf(os.Stderr, "Build a signed Android APK from a Go package.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment:\n")
		fmt.Fprintf(os.Stderr, "  ANDROID_HOME       Android SDK path (default: ~/Android/Sdk)\n")
		fmt.Fprintf(os.Stderr, "  ANDROID_NDK_HOME   Android NDK path (default: latest in $ANDROID_HOME/ndk/)\n")
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	goPackage := flag.Arg(0)

	var archs []string
	for _, a := range strings.Split(*arch, ",") {
		a = strings.TrimSpace(a)
		if a != "" {
			archs = append(archs, a)
		}
	}

	cfg := apkbuild.Config{
		GoPackage:     goPackage,
		OutputPath:    *output,
		ManifestPath:  *manifest,
		Architectures: archs,
		APILevel:      *apiLevel,
		KeystorePath:  *keystore,
		LibName:       *libName,
	}

	if err := apkbuild.Build(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputPath := cfg.OutputPath
	if outputPath == "" {
		name := cfg.LibName
		if name == "" {
			name = goPackage
		}
		outputPath = name + ".apk"
	}
	fmt.Printf("APK: %s\n", outputPath)
}
