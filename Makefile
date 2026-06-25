BINARY_NAME := caching-proxy
BUILD_DIR   := bin
CMD_PATH    := ./cmd/caching-proxy

.PHONY: all build run test cover vet fmt tidy clean ci

all: build

## Build the binary into bin/
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

## Build and run the proxy (pass args with: make run ARGS="--port 3000 --origin http://dummyjson.com")
run: build
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

## Run the test suite with race detector + coverage
test:
	go test ./... -v -race -cover

## Run tests and write an HTML coverage report to coverage.html
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

## Static analysis
vet:
	go vet ./...

## Format all source files
fmt:
	gofmt -l -w .

## Tidy go.mod/go.sum
tidy:
	go mod tidy

## Remove build artifacts
clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

## Run everything CI runs, locally
ci: fmt vet test build