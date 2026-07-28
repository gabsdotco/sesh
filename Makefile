BINARY_NAME := sesh
BUILD_DIR   := build
INSTALL_PATH := /usr/local/bin
GO          := go

.PHONY: all build clean install uninstall test test-cover fmt vet lint tidy release dev run deps help

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sesh
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "Clean complete"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete. Run '$(BINARY_NAME)' to use."

uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_PATH)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstall complete"

test:
	@which gotestsum > /dev/null 2>&1 && gotestsum --format testname ./... || $(GO) test -v ./...

test-cover:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	@echo ""
	@echo "Coverage by package:"
	go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "Generating HTML report: coverage.html"
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in a browser to see detailed coverage."

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

vet:
	@echo "Running go vet..."
	$(GO) vet ./...

lint: fmt vet
	@echo "Linting complete"

dev: build
	@echo "Running in development mode..."
	$(BUILD_DIR)/$(BINARY_NAME)

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

tidy:
	@echo "Tidying go modules..."
	$(GO) mod tidy
	@echo "Done. Check git diff to see if go.mod or go.sum changed."

deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

.DEFAULT_GOAL := help

release: clean test
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	
	@echo "  darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 ./cmd/sesh
	
	@echo "  darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 ./cmd/sesh
	
	@echo "  linux/amd64..."
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 ./cmd/sesh
	
	@echo "  linux/arm64..."
	@GOOS=linux GOARCH=arm64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 ./cmd/sesh
	
	@echo "Release binaries built in $(BUILD_DIR)/release/"
	@ls -lh $(BUILD_DIR)/release/

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Remove build artifacts and coverage files"
	@echo "  install      - Build and install to $(INSTALL_PATH)"
	@echo "  uninstall    - Remove from $(INSTALL_PATH)"
	@echo "  test         - Run tests (uses gotestsum if available)"
	@echo "  test-cover   - Run tests with coverage report + HTML"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  lint         - Run fmt and vet"
	@echo "  tidy         - Run go mod tidy"
	@echo "  release      - Build cross-platform binaries for release"
	@echo "  dev          - Build and run"
	@echo "  run          - Build and run (alias)"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  help         - Show this help"
