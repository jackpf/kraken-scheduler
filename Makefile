APP_NAME = kraken-scheduler
SRC_DIR = ./src/main
TARGET_DIR = ./target
BUILD_VARS ?=
POSTFIX ?=

clean:
	rm -rf target

test:
	go test -v $(SRC_DIR)/...

build:
	mkdir -p $(TARGET_DIR)
	$(BUILD_VARS) go build -o $(TARGET_DIR)/$(APP_NAME)$(POSTFIX) $(SRC_DIR)

cross-build:
	POSTFIX=-windows-amd64.exe BUILD_VARS="GOOS=windows GOARCH=amd64" $(MAKE) build
	POSTFIX=-macos-amd64 BUILD_VARS="GOOS=darwin GOARCH=amd64" $(MAKE) build
	POSTFIX=-macos-arm64 BUILD_VARS="GOOS=darwin GOARCH=arm64" $(MAKE) build
	POSTFIX=-linux-amd64 BUILD_VARS="GOOS=linux GOARCH=amd64" $(MAKE) build

install: build
	cp ./target/kraken-scheduler /usr/local/bin