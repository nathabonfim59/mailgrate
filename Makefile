# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "-X github.com/natbonfim/mailgrate/internal/cmd.Version=$(VERSION) -X github.com/natbonfim/mailgrate/internal/cmd.BuildTime=$(BUILD_TIME)"
BINARY_NAME = mailgrate
BUILD_DIR = build

.PHONY: all run build clean test lint build-static

all: build

# Run the application
run:
	go run cmd/mailgrade/main.go

# Build the application
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/mailgrade/main.go

# Build a static binary using musl
build-static:
	@echo "Building static binary with musl $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-static cmd/mailgrade/main.go
	@echo "Static binary built as $(BUILD_DIR)/$(BINARY_NAME)-static"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

