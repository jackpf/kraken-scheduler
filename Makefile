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

.PHONY: package
package:
	docker build -t jackpfarrelly/kraken-scheduler:latest .

.PHONY: release
release:
	docker push jackpfarrelly/kraken-scheduler:latest