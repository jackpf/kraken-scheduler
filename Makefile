APP_NAME = kraken-scheduler
SRC_DIR = ./src/main
TARGET_DIR = ./target
BUILD_VARS ?=
POSTFIX ?=

.PHONY: clean
clean:
	rm -rf target

.PHONY: test
test:
	gofmt -l $(SRC_DIR)
	go test -v $(SRC_DIR)/...

.PHONY: build
build:
	mkdir -p $(TARGET_DIR)
	$(BUILD_VARS) go build -o $(TARGET_DIR)/$(APP_NAME)$(POSTFIX) $(SRC_DIR)

.PHONY: cross-build
cross-build:
	POSTFIX=-windows-amd64.exe BUILD_VARS="GOOS=windows GOARCH=amd64" $(MAKE) build
	POSTFIX=-macos-amd64 BUILD_VARS="GOOS=darwin GOARCH=amd64" $(MAKE) build
	POSTFIX=-macos-arm64 BUILD_VARS="GOOS=darwin GOARCH=arm64" $(MAKE) build
	POSTFIX=-linux-amd64 BUILD_VARS="GOOS=linux GOARCH=amd64" $(MAKE) build
	POSTFIX=-linux-arm BUILD_VARS="GOOS=linux GOARCH=arm" $(MAKE) build
	POSTFIX=-linux-arm64 BUILD_VARS="GOOS=linux GOARCH=arm64" $(MAKE) build

BUILDX=docker buildx create --use; docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v6,linux/arm/v7

.PHONY: package
package:
	$(BUILDX) -t jackpfarrelly/kraken-scheduler:latest .

.PHONY: release
release:
	$(BUILDX) -t jackpfarrelly/kraken-scheduler:latest --push .