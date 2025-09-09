.PHONY: test build clean install run coverage

BINARY_NAME=task-cli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Test commands
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build commands
build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} .

build-all:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe .

# Development commands
run: build
	./bin/${BINARY_NAME}

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

install: build
	cp bin/${BINARY_NAME} ${GOPATH}/bin/

# Linting and formatting
fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

# Development workflow
dev: fmt vet test build