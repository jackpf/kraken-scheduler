APP_NAME = kraken-schedule
SRC_DIR = ./src/main
TARGET_DIR = ./target

test:
	go test -v $(SRC_DIR)/...

build: test
	mkdir -p $(TARGET_DIR)
	go build -o $(TARGET_DIR)/$(APP_NAME) $(SRC_DIR)

install: build
	cp ./target/kraken-schedule /usr/local/bin