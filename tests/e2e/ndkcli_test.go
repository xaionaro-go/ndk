//go:build e2e

// E2E tests for ndkcli commands running on an Android device.
//
// These tests execute ndkcli subcommands as subprocesses and verify
// their output. The ndkcli binary must be pre-built and available
// at the path specified by the NDKCLI_BIN environment variable
// (default: /data/local/tmp/ndkcli).
package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ndkcliBin() string {
	if bin := os.Getenv("NDKCLI_BIN"); bin != "" {
		return bin
	}
	return "/data/local/tmp/ndkcli"
}

func runNdkcli(
	t *testing.T,
	args ...string,
) (string, error) {
	t.Helper()
	cmd := exec.Command(ndkcliBin(), args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func requireNdkcli(
	t *testing.T,
	args ...string,
) string {
	t.Helper()
	out, err := runNdkcli(t, args...)
	require.NoError(t, err, "ndkcli %s failed: %s", strings.Join(args, " "), out)
	return out
}

func TestNdkcli_Help(t *testing.T) {
	out := requireNdkcli(t, "--help")
	assert.Contains(t, out, "ndkcli")
	assert.Contains(t, out, "camera")
	assert.Contains(t, out, "audio")
	assert.Contains(t, out, "sensor")
}

func TestNdkcli_CameraManagerCameraIDList(t *testing.T) {
	out, err := runNdkcli(t, "camera", "manager", "camera-id-list")
	// May fail on emulator without camera permission, but should not crash.
	t.Logf("camera manager camera-id-list: err=%v output=%q", err, out)
}

func TestNdkcli_CameraListDetails(t *testing.T) {
	out, err := runNdkcli(t, "camera", "list-details")
	t.Logf("camera list-details: err=%v output=%q", err, out)
}

func TestNdkcli_AudioStreamBuilderNew(t *testing.T) {
	out := requireNdkcli(t, "audio", "stream-builder", "new")
	assert.Contains(t, out, "created successfully")
}

func TestNdkcli_AudioRecord(t *testing.T) {
	out, err := runNdkcli(t, "audio", "record",
		"--output", "/data/local/tmp/ndkcli_test.pcm",
		"--duration", "1s",
		"--sample-rate", "44100",
		"--channels", "1",
	)
	t.Logf("audio record: err=%v output=%q", err, out)
	// Cleanup.
	os.Remove("/data/local/tmp/ndkcli_test.pcm")
}

func TestNdkcli_SensorRead(t *testing.T) {
	out, err := runNdkcli(t, "sensor", "read", "--type", "1", "--duration", "1s")
	t.Logf("sensor read: err=%v output=%q", err, out)
}

func TestNdkcli_SensorManagerDefaultSensor(t *testing.T) {
	out, err := runNdkcli(t, "sensor", "manager", "default-sensor", "--value", "1")
	t.Logf("sensor manager default-sensor: err=%v output=%q", err, out)
}

func TestNdkcli_ThermalManagerCurrentStatus(t *testing.T) {
	out := requireNdkcli(t, "thermal", "manager", "current-status")
	t.Logf("thermal status: %q", out)
}

func TestNdkcli_ThermalMonitor(t *testing.T) {
	out, err := runNdkcli(t, "thermal", "monitor", "--interval", "500ms", "--duration", "1s")
	t.Logf("thermal monitor: err=%v output=%q", err, out)
}

func TestNdkcli_EGLInfo(t *testing.T) {
	out := requireNdkcli(t, "egl", "info")
	assert.Contains(t, out, "Vendor", "EGL info should contain Vendor")
	t.Logf("egl info:\n%s", out)
}

func TestNdkcli_EGLConfigs(t *testing.T) {
	out, err := runNdkcli(t, "egl", "configs")
	t.Logf("egl configs: err=%v output=%q", err, out)
}

func TestNdkcli_GLES2Info(t *testing.T) {
	out := requireNdkcli(t, "gles2", "info")
	assert.Contains(t, out, "GL_VERSION", "gles2 info should contain GL_VERSION")
	t.Logf("gles2 info:\n%s", out)
}

func TestNdkcli_GLES3Info(t *testing.T) {
	out, err := runNdkcli(t, "gles3", "info")
	t.Logf("gles3 info: err=%v output=%q", err, out)
}

func TestNdkcli_MediaCodecs(t *testing.T) {
	out := requireNdkcli(t, "media", "codecs")
	assert.Contains(t, out, "video/avc", "should probe H.264")
	t.Logf("media codecs:\n%s", out)
}

func TestNdkcli_MediaNewEncoder(t *testing.T) {
	out, err := runNdkcli(t, "media", "new-encoder", "--mime_type", "video/avc")
	t.Logf("media new-encoder: err=%v output=%q", err, out)
}

func TestNdkcli_MediaNewDecoder(t *testing.T) {
	out, err := runNdkcli(t, "media", "new-decoder", "--mime_type", "video/avc")
	t.Logf("media new-decoder: err=%v output=%q", err, out)
}

func TestNdkcli_ConfigShow(t *testing.T) {
	out := requireNdkcli(t, "config", "show")
	assert.Contains(t, out, "Density", "config show should contain Density")
	t.Logf("config show:\n%s", out)
}

func TestNdkcli_ConfigDensity(t *testing.T) {
	out := requireNdkcli(t, "config", "config", "density")
	t.Logf("config density: %q", out)
}

func TestNdkcli_ConfigSdkVersion(t *testing.T) {
	out := requireNdkcli(t, "config", "config", "sdk-version")
	t.Logf("config sdk-version: %q", out)
}

func TestNdkcli_FontMatch(t *testing.T) {
	out, err := runNdkcli(t, "font", "match", "--family", "sans-serif", "--weight", "400")
	t.Logf("font match: err=%v output=%q", err, out)
}

func TestNdkcli_PermissionCheck(t *testing.T) {
	out, err := runNdkcli(t, "permission", "check",
		"--name", "android.permission.CAMERA",
		"--pid", "1000", "--uid", "1000",
	)
	t.Logf("permission check: err=%v output=%q", err, out)
}

func TestNdkcli_TraceIsEnabled(t *testing.T) {
	out := requireNdkcli(t, "trace", "is-enabled")
	t.Logf("trace is-enabled: %q", out)
}

func TestNdkcli_TraceBeginEndSection(t *testing.T) {
	out := requireNdkcli(t, "trace", "begin-section", "--section-name", "test_section")
	t.Logf("trace begin-section: %q", out)

	out = requireNdkcli(t, "trace", "end-section")
	t.Logf("trace end-section: %q", out)
}

func TestNdkcli_TraceSetCounter(t *testing.T) {
	out := requireNdkcli(t, "trace", "set-counter", "--counter-name", "test_counter", "--counter-value", "42")
	t.Logf("trace set-counter: %q", out)
}

func TestNdkcli_LogWrite(t *testing.T) {
	out := requireNdkcli(t, "log", "write", "--tag", "ndktest", "--message", "e2e test", "--priority", "4")
	t.Logf("log write: %q", out)
}

func TestNdkcli_NnapiProbe(t *testing.T) {
	out, err := runNdkcli(t, "nnapi", "probe")
	t.Logf("nnapi probe: err=%v output=%q", err, out)
}

func TestNdkcli_LooperTest(t *testing.T) {
	out, err := runNdkcli(t, "looper", "test")
	t.Logf("looper test: err=%v output=%q", err, out)
}

func TestNdkcli_WindowQuery(t *testing.T) {
	out, err := runNdkcli(t, "window", "query")
	t.Logf("window query: err=%v output=%q", err, out)
}

func TestNdkcli_StorageObb(t *testing.T) {
	out, err := runNdkcli(t, "storage", "obb", "--file", "/nonexistent.obb")
	t.Logf("storage obb: err=%v output=%q", err, out)
}

func TestNdkcli_BinderprocessStartThreadPool(t *testing.T) {
	out, err := runNdkcli(t, "binderprocess", "start-thread-pool")
	t.Logf("binderprocess start-thread-pool: err=%v output=%q", err, out)
}
