BINARY_NAME        := scylla
CONFIG_PACKAGE     := github.com/luganodes/slashing-observer/config
VERSION            := $(shell git describe --tags --always --dirty | sed 's/^v//')
BUILD_FLAGS        := -ldflags "\
	-X $(CONFIG_PACKAGE).EXTERNAL_VERSION=$(VERSION) \
	-X $(CONFIG_PACKAGE).EXTERNAL_APP_NAME=$(BINARY_NAME)"



.PHONY: all build clean

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME)
	@echo "Build complete: ./$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	@if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
	@echo "Clean complete."

