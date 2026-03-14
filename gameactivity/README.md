# GameActivity

This package provides Go bindings for the [AGDK GameActivity](https://developer.android.com/games/agdk/game-activity) C library.

GameActivity is an alternative to NativeActivity that provides:
- Better input event handling with motion event batching
- Text input support via GameTextInput
- Modern lifecycle management built on AppCompatActivity

## Prerequisites

GameActivity headers are not part of the standard Android NDK. To use this package:

1. Download the `games-activity` AAR from Maven Central:
   ```
   https://maven.google.com/web/index.html#androidx.games:games-activity
   ```

2. Extract the AAR (it's a ZIP file) and copy the C headers from
   `prefab/modules/game-activity/include/` to `gameactivity/include/`.

3. Copy the C source files from `prefab/modules/game-activity/src/` to
   `gameactivity/csrc/`.

Alternatively, install the AGDK via Android Studio's SDK Manager.

## Build tags

The Go source files use `//go:build ignore` until the AGDK headers are
vendored into the `include/` and `csrc/` subdirectories. After vendoring,
remove the build constraint to enable compilation.
