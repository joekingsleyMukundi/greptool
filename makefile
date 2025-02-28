APP_NAME = grepper
GO = go
GO_FILES = $(shell find ./grepper -name '*.go')
GO_VERSION = $(shell $(GO) version)
all: build

build:
	$(GO) build -o ./$(APP_NAME) ./grepper  # Compile Go code inside the 'grepper' directory
windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o ./$(APP_NAME).exe ./grepper

darwin:
	GOOS=darwin GOARCH=amd64 $(GO) build -o ./$(APP_NAME)_mac_intel ./grepper
	GOOS=darwin GOARCH=arm64 $(GO) build -o ./$(APP_NAME)_mac_m1 ./grepper
linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o ./$(APP_NAME)_linux_amd64 ./grepper
clean:
	rm -f ./$(APP_NAME) ./$(APP_NAME).exe ./$(APP_NAME)_mac_intel ./$(APP_NAME)_mac_m1 ./$(APP_NAME)_linux_amd64
install:
	sudo mv ./$(APP_NAME)_linux_amd64 /usr/local/bin/grepper
confirm:
	sudo chmod +x /usr/local/bin/grepper

.PHONY: all build clean install windows darwin linux confirm
