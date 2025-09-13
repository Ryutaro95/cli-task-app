# Task CLI Makefile

BINARY_NAME=task-cli
MAIN_PATH=./cmd/task-cli
BUILD_FLAGS=-ldflags="-s -w"

.PHONY: all build test clean run deps help

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  test     - Run all tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Build and run the application" 
	@echo "  deps     - Download dependencies"
	@echo "  help     - Show this help"