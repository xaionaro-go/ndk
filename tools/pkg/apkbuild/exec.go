package apkbuild

import (
	"os"
	"os/exec"
)

// runCmd executes a command at the given path with the provided arguments,
// forwarding stdout/stderr to the current process.
func runCmd(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCmdInDir executes a command by name (resolved via PATH) in the given
// working directory, forwarding stdout/stderr to the current process.
func runCmdInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
