ANDROID_HOME ?= $(HOME)/Android/Sdk
GOMOBILE     ?= $(shell go env GOPATH)/bin/gomobile
API_LEVEL    ?= 35

AAR := sensorbridge.aar

.PHONY: all clean install-tools

all: $(AAR)

$(AAR): sensorbridge.go
	ANDROID_HOME=$(ANDROID_HOME) $(GOMOBILE) bind \
		-target=android/arm64 \
		-androidapi=$(API_LEVEL) \
		-o $@ .

install-tools:
	go install golang.org/x/mobile/cmd/gomobile@latest
	go install golang.org/x/mobile/cmd/gobind@latest
	$(GOMOBILE) init

clean:
	rm -f $(AAR) sensorbridge-sources.jar
