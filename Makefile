.PHONY: all build test clean install build-all lint

all: build

BINARY=sdeconvert
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=${VERSION}"
GO_DIR=wanderer-sde

build:
	cd $(GO_DIR) && go build ${LDFLAGS} -o ../bin/${BINARY} ./cmd/sdeconvert

test:
	cd $(GO_DIR) && go test -v ./...

test-coverage:
	cd $(GO_DIR) && go test -coverprofile=coverage.out ./...
	cd $(GO_DIR) && go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ output/ $(GO_DIR)/coverage.out $(GO_DIR)/coverage.html

install:
	cd $(GO_DIR) && go install ${LDFLAGS} ./cmd/sdeconvert

# Cross-compilation
build-all:
	cd $(GO_DIR) && GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ../bin/${BINARY}-linux-amd64 ./cmd/sdeconvert
	cd $(GO_DIR) && GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ../bin/${BINARY}-darwin-amd64 ./cmd/sdeconvert
	cd $(GO_DIR) && GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ../bin/${BINARY}-darwin-arm64 ./cmd/sdeconvert
	cd $(GO_DIR) && GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ../bin/${BINARY}-windows-amd64.exe ./cmd/sdeconvert

# Development helpers
fmt:
	cd $(GO_DIR) && go fmt ./...

vet:
	cd $(GO_DIR) && go vet ./...

lint: fmt vet

tidy:
	cd $(GO_DIR) && go mod tidy
