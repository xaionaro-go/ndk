//go:build e2e

// E2E tests that exercise NDK idiomatic Go bindings on an Android device.
//
// Each test creates and tears down its own resources. Tests tolerate
// emulator limitations (e.g. absent camera, limited sensors) by
// skipping rather than failing.
//
// Build the test binary:
//
//	CGO_ENABLED=1 GOOS=android GOARCH=amd64 \
//	  CC=$NDK/.../x86_64-linux-android35-clang \
//	  go test -c -tags e2e -o tests/e2e/examples.test ./tests/e2e/
//
// Run on device:
//
//	adb push tests/e2e/examples.test /data/local/tmp/
//	adb shell /data/local/tmp/examples.test -test.v
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndroidGoLab/ndk/audio"
	"github.com/AndroidGoLab/ndk/camera"
	"github.com/AndroidGoLab/ndk/config"
	"github.com/AndroidGoLab/ndk/egl"
	"github.com/AndroidGoLab/ndk/log"
	"github.com/AndroidGoLab/ndk/media"
	"github.com/AndroidGoLab/ndk/sensor"
	"github.com/AndroidGoLab/ndk/thermal"
	"github.com/AndroidGoLab/ndk/trace"
)

func TestExample_AudioStreamInfo(t *testing.T) {
	builder, err := audio.NewStreamBuilder()
	require.NoError(t, err, "NewStreamBuilder")
	defer builder.Close()

	builder.
		SetSampleRate(44100).
		SetChannelCount(1).
		SetFormat(audio.PcmI16).
		SetDirection(audio.Output)

	stream, err := builder.Open()
	if err != nil {
		t.Skipf("cannot open audio stream on this device: %v", err)
	}
	defer stream.Close()

	rate := stream.SampleRate()
	assert.Greater(t, rate, int32(0), "SampleRate should be positive")
	t.Logf("SampleRate: %d Hz", rate)

	chans := stream.ChannelCount()
	assert.Greater(t, chans, int32(0), "ChannelCount should be positive")
	t.Logf("ChannelCount: %d", chans)

	burst := stream.FramesPerBurst()
	assert.Greater(t, burst, int32(0), "FramesPerBurst should be positive")
	t.Logf("FramesPerBurst: %d", burst)
}

func TestExample_CameraList(t *testing.T) {
	mgr := camera.NewManager()
	require.NotNil(t, mgr, "NewManager returned nil")
	defer mgr.Close()

	ids, err := mgr.CameraIDList()
	// On emulator without camera the list may be empty but no error.
	assert.NoError(t, err, "CameraIDList should not return an error")
	t.Logf("camera IDs: %v (count=%d)", ids, len(ids))
}

func TestExample_SensorList(t *testing.T) {
	mgr := sensor.GetInstance()
	require.NotNil(t, mgr, "GetInstance returned nil")

	accel := mgr.DefaultSensor(sensor.Accelerometer)
	require.NotNil(t, accel, "DefaultSensor(Accelerometer) returned nil")

	name := accel.Name()
	if name == "" {
		t.Skip("accelerometer not available on this emulator")
	}

	assert.NotEmpty(t, name, "sensor name should be non-empty")
	t.Logf("accelerometer: name=%q vendor=%q type=%s resolution=%.6f minDelay=%d",
		name, accel.Vendor(), sensor.Type(accel.Type()), accel.Resolution(), accel.MinDelay())
}

func TestExample_ThermalStatus(t *testing.T) {
	mgr := thermal.NewManager()
	require.NotNil(t, mgr, "NewManager returned nil")
	defer mgr.Close()

	// CurrentStatus should not panic; on a cool emulator it returns StatusNone.
	status := mgr.CurrentStatus()
	assert.GreaterOrEqual(t, int32(status), int32(thermal.StatusError),
		"status should be a valid ThermalStatus value")
	t.Logf("thermal status: %s (%d)", status, status)
}

func TestExample_EGLInfo(t *testing.T) {
	var defaultDisplay egl.EGLNativeDisplayType
	dpy := egl.GetDisplay(defaultDisplay)
	require.NotNil(t, dpy, "GetDisplay returned EGL_NO_DISPLAY")

	var major, minor egl.Int
	ret := egl.Initialize(dpy, &major, &minor)
	require.NotEqual(t, egl.Boolean(0), ret,
		"Initialize failed, error=0x%X", egl.GetError())
	defer egl.Terminate(dpy)

	t.Logf("EGL %d.%d", major, minor)

	vendor := egl.QueryString(dpy, egl.EGL_VENDOR)
	assert.NotEmpty(t, vendor, "EGL vendor string should be non-empty")
	t.Logf("EGL vendor: %s", vendor)
}

func TestExample_MediaCodecs(t *testing.T) {
	t.Run("Encoder", func(t *testing.T) {
		enc := media.NewEncoder("video/avc")
		require.NotNil(t, enc, "NewEncoder returned nil")
		assert.NotNil(t, enc.Pointer(), "encoder Pointer should be non-nil")
		t.Log("video/avc encoder created")
		enc.Close()
	})

	t.Run("Decoder", func(t *testing.T) {
		dec := media.NewDecoder("video/avc")
		require.NotNil(t, dec, "NewDecoder returned nil")
		assert.NotNil(t, dec.Pointer(), "decoder Pointer should be non-nil")
		t.Log("video/avc decoder created")
		dec.Close()
	})
}

func TestExample_ConfigShow(t *testing.T) {
	cfg := config.NewConfig()
	require.NotNil(t, cfg, "NewConfig returned nil")
	defer cfg.Close()

	density := cfg.Density()
	// An empty config returns 0 for density (no activity context), so
	// only verify the call succeeds without panic.
	t.Logf("density: %d", density)

	sdk := cfg.SdkVersion()
	// Empty config returns 0 for sdk version, but the call must not crash.
	t.Logf("SdkVersion: %d", sdk)

	t.Logf("ScreenSize: %d  Orientation: %d  ScreenDp: %dx%d",
		cfg.ScreenSize(), cfg.Orientation(),
		cfg.ScreenWidthDp(), cfg.ScreenHeightDp())
}

func TestExample_TraceEnabled(t *testing.T) {
	// Just verify the call does not panic.
	enabled := trace.IsEnabled()
	t.Logf("trace.IsEnabled: %v", enabled)
}

func TestExample_LogWrite(t *testing.T) {
	ret := log.Write(int32(log.Info), "ndktest", "E2E test message")
	assert.GreaterOrEqual(t, ret, int32(0),
		"log.Write should return >= 0 (got %d)", ret)
	t.Logf("log.Write returned: %d", ret)
}
