// Check common Android permissions for the current process.
//
// Queries the Android Permission Manager for four commonly requested
// permissions -- INTERNET, CAMERA, RECORD_AUDIO, and ACCESS_FINE_LOCATION --
// using the calling process's own PID and UID. For each permission the
// program prints whether it is currently granted or denied.
//
// The result tells you what the system would allow right now; a "denied"
// result means the app either never declared the permission in its manifest
// or the user has not yet granted it at runtime.
//
// This program must run on an Android device.
package main

import (
	"fmt"
	"syscall"

	"github.com/xaionaro-go/ndk/permission"
)

func main() {
	pid := permission.Pid_t(syscall.Getpid())
	uid := permission.Uid_t(syscall.Getuid())

	fmt.Printf("Checking permissions for PID %d, UID %d\n\n", pid, uid)

	permissions := []struct {
		label string
		name  string
	}{
		{"INTERNET", "android.permission.INTERNET"},
		{"CAMERA", "android.permission.CAMERA"},
		{"RECORD_AUDIO", "android.permission.RECORD_AUDIO"},
		{"ACCESS_FINE_LOCATION", "android.permission.ACCESS_FINE_LOCATION"},
	}

	for _, p := range permissions {
		var result int32
		status := permission.CheckPermission(p.name, pid, uid, &result)
		if status != 0 {
			fmt.Printf("  %-24s  error (status %d)\n", p.label+":", status)
			continue
		}

		switch permission.Result(result) {
		case permission.Granted:
			fmt.Printf("  %-24s  granted\n", p.label+":")
		case permission.Denied:
			fmt.Printf("  %-24s  denied\n", p.label+":")
		default:
			fmt.Printf("  %-24s  unknown result (%d)\n", p.label+":", result)
		}
	}
}
