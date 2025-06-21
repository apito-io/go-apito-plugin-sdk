.PHONY: help test build-example clean tidy fmt lint

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run tests"
	@echo "  build-example - Build the example plugin"
	@echo "  clean         - Clean build artifacts"
	@echo "  tidy          - Run go mod tidy"
	@echo "  fmt           - Format Go code"
	@echo "  lint          - Run linter"

# Run tests
test:
	go test -v ./...

# Build the example plugin
build-example:
	cd examples/simple-plugin && go build -o simple-plugin main.go

# Clean build artifacts
clean:
	rm -f examples/simple-plugin/simple-plugin
	go clean -cache

# Tidy dependencies
tidy:
	go mod tidy
	cd examples/simple-plugin && go mod tidy

# Format code
fmt:
	go fmt ./...
	cd examples/simple-plugin && go fmt ./...

# Run linter (requires golangci-lint to be installed)
lint:
	golangci-lint run
	cd examples/simple-plugin && golangci-lint run 