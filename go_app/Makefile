.PHONY: build test clean lint

# Build the project
build:
	go build ./...

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	go clean
	rm -rf bin/

# Run linters
lint:
	go vet ./...
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate API documentation
docs:
	go generate ./...

# Run all checks
check: lint test build
