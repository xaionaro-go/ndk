export PATH := $(HOME)/go/bin:$(PATH)

ANDROID_HOME ?= $(HOME)/Android/Sdk
NDK_PATH ?= $(or $(ANDROID_NDK_HOME),$(shell ls -d $(ANDROID_HOME)/ndk/* 2>/dev/null | sort -V | tail -1))

# Full module list for NDK production builds.
MODULES  := aaudio camera sensor gles2 gles3 egl vulkan media \
            input nativewindow hardwarebuffer binder thermal \
            performancehint neuralnetworks trace logging font \
            imagedecoder midi multinetwork sync choreographer \
            configuration asset looper nativeactivity surfacecontrol \
            sharedmem permission bitmap storage surfacetexture

# Modules with testdata fixtures (no NDK required).
FIXTURE_MODULES := $(notdir $(wildcard tools/pkg/specgen/testdata/*/))

NDK_SYSROOT := $(NDK_PATH)/toolchains/llvm/prebuilt/linux-x86_64/sysroot/usr/include
C2FFI_BIN   ?= c2ffi

.PHONY: all capi specs idiomatic clean regen fixtures test lint check-examples e2e e2e-build e2e-examples e2e-examples-test e2e-audio ndkcli ndkcli-commands

all: specs capi idiomatic

# Stage 1+2: Generate specs from C headers via c2ffi (requires NDK + c2ffi)
specs:
	@for m in $(MODULES); do \
		manifest="capi/manifests/$$m.yaml"; \
		[ -f "$$manifest" ] || continue; \
		echo "specgen $$m (c2ffi)"; \
		go run ./tools/cmd/specgen \
			-manifest "$$manifest" \
			-ndk-include "$(NDK_SYSROOT)" \
			-c2ffi-bin "$(C2FFI_BIN)" \
			-out "spec/generated/$$m.yaml"; \
	done

# Stage 2: Generate raw CGo binding packages from specs + manifests
capi:
	@for m in $(MODULES); do \
		manifest="capi/manifests/$$m.yaml"; \
		spec="spec/generated/$$m.yaml"; \
		[ -f "$$manifest" ] || continue; \
		[ -f "$$spec" ] || continue; \
		echo "capigen $$m"; \
		go run ./tools/cmd/capigen \
			-spec "$$spec" \
			-manifest "$$manifest" \
			-out "capi/$$m/"; \
	done

# Stage 2 (fixture mode): Extract specs from testdata fixtures
specs-fixtures:
	@for m in $(FIXTURE_MODULES); do \
		case $$m in simple|edgecases) continue;; esac; \
		echo "specgen $$m (fixture)"; \
		go run ./tools/cmd/specgen \
			-module $$m \
			-pkg tools/pkg/specgen/testdata/$$m \
			-out spec/generated/$$m.yaml; \
	done

# Stage 3: Generate idiomatic Go from specs + overlays
idiomatic:
	@for overlay in spec/overlays/*.yaml; do \
		m=$$(basename "$$overlay" .yaml); \
		spec="spec/generated/$$m.yaml"; \
		[ -f "$$spec" ] || continue; \
		goname=$$(grep 'go_name:' "$$overlay" | head -1 | awk '{print $$2}'); \
		[ -z "$$goname" ] && goname="$$m"; \
		echo "idiomgen $$m -> $$goname/"; \
		go run ./tools/cmd/idiomgen \
			-spec "$$spec" \
			-overlay "$$overlay" \
			-templates templates/ \
			-capi-dir "capi/$$m/" \
			-out "$$goname/"; \
	done

# Generate everything from fixtures (no NDK required)
fixtures: specs-fixtures idiomatic

# Regenerate everything from scratch (requires NDK + c2ffi)
regen: clean specs capi idiomatic

test:
	go test $$(go list ./... | grep -v -E '/(capi|tests|examples)/|/ndk/[a-z][a-z0-9]*$$') -count=1

lint:
	golangci-lint run ./tools/...

# Cross-compile all examples for Android arm64 to catch compile errors (requires NDK)
check-examples:
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
		CC=$(NDK_PATH)/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android35-clang \
		go build ./examples/...

# Cross-compile E2E test binary for Android x86_64 (requires NDK)
e2e-build:
	CGO_ENABLED=1 GOOS=android GOARCH=amd64 \
		CC=$(NDK_PATH)/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android35-clang \
		go build -o tests/e2e/e2e_test ./tests/e2e

# Run full E2E test on Android emulator (requires SDK + NDK + KVM)
e2e: e2e-build
	ANDROID_HOME=$(dir $(NDK_PATH)) ./tests/e2e/run.sh

# Run all examples on Android emulator (requires running emulator + NDK)
e2e-examples:
	./tests/e2e/run-examples.sh

# Run examples E2E Go test on Android emulator (requires running emulator + NDK)
e2e-examples-test:
	./tests/e2e/run-examples-test.sh

# Run audio recording E2E test (requires running emulator with audio + NDK)
e2e-audio:
	./tests/e2e/run-audio-e2e.sh

# Build ndkcli for Android (requires NDK)
ndkcli:
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
		CC=$(NDK_PATH)/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android$(APK_API)-clang \
		go build -o ndkcli ./cmd/ndkcli/

# Print all ndkcli subcommands (extracted from source, no binary needed)
ndkcli-commands:
	@go run ./tools/cmd/ndkcli-commands

clean:
	@for m in $(MODULES); do rm -rf "capi/$$m/"; done
	rm -rf spec/generated/

# --- APK Build (examples) ---

APK_API      ?= 35
APK_ARCH     ?= arm64

ifeq ($(APK_ARCH),arm64)
  APK_ABI        := arm64-v8a
  APK_GOARCH     := arm64
  APK_CC_PREFIX  := aarch64-linux-android$(APK_API)
else ifeq ($(APK_ARCH),x86_64)
  APK_ABI        := x86_64
  APK_GOARCH     := amd64
  APK_CC_PREFIX  := x86_64-linux-android$(APK_API)
endif

APK_CC         := $(NDK_PATH)/toolchains/llvm/prebuilt/linux-x86_64/bin/$(APK_CC_PREFIX)-clang
APK_BUILD_TOOLS := $(shell ls -d $(ANDROID_HOME)/build-tools/* 2>/dev/null | sort -V | tail -1)
APK_PLATFORM   := $(ANDROID_HOME)/platforms/android-$(APK_API)/android.jar
APK_KEYSTORE   := $(HOME)/.android/debug.keystore

$(APK_KEYSTORE):
	mkdir -p $(dir $@)
	keytool -genkeypair -v -keystore $@ -storepass android \
		-alias androiddebugkey -keypass android -keyalg RSA \
		-keysize 2048 -validity 10000 \
		-dname "CN=Debug,O=Debug,C=US"

.PHONY: apk-displaycamera

apk-displaycamera: $(APK_KEYSTORE)
	$(eval DIR := examples/camera/display)
	$(eval BUILD := $(DIR)/_build)
	@rm -rf $(BUILD)
	@mkdir -p $(BUILD)/lib/$(APK_ABI)
	CGO_ENABLED=1 GOOS=android GOARCH=$(APK_GOARCH) CC=$(APK_CC) \
		go build -buildmode=c-shared \
		-o $(BUILD)/lib/$(APK_ABI)/libdisplaycamera.so \
		./$(DIR)
	@rm -f $(BUILD)/lib/$(APK_ABI)/libdisplaycamera.h
	$(APK_BUILD_TOOLS)/aapt package -f \
		-M $(DIR)/AndroidManifest.xml \
		-I $(APK_PLATFORM) \
		-F $(BUILD)/base.apk
	cd $(BUILD) && zip -r base.apk lib/
	$(APK_BUILD_TOOLS)/zipalign -f -p 4 \
		$(BUILD)/base.apk $(BUILD)/aligned.apk
	$(APK_BUILD_TOOLS)/apksigner sign \
		--ks $(APK_KEYSTORE) --ks-pass pass:android \
		--key-pass pass:android \
		$(BUILD)/aligned.apk
	mv $(BUILD)/aligned.apk $(DIR)/displaycamera.apk
	@rm -rf $(BUILD)
	@echo "APK: $(DIR)/displaycamera.apk"
