package apkbuild

import "fmt"

// archMapping returns the GOARCH, Android ABI name, and NDK CC prefix
// for a given architecture alias.
func archMapping(
	arch string,
	apiLevel int,
) (goarch, abi, ccPrefix string) {
	switch arch {
	case "arm64":
		return "arm64", "arm64-v8a", fmt.Sprintf("aarch64-linux-android%d", apiLevel)
	case "x86_64":
		return "amd64", "x86_64", fmt.Sprintf("x86_64-linux-android%d", apiLevel)
	case "arm":
		return "arm", "armeabi-v7a", fmt.Sprintf("armv7a-linux-androideabi%d", apiLevel)
	case "x86":
		return "386", "x86", fmt.Sprintf("i686-linux-android%d", apiLevel)
	default:
		return "", "", ""
	}
}
