# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint

# Build parameters
BIN_DIR=bin
PROJECTCTL_BIN=$(BIN_DIR)/projectctl
PROJECTS_BIN=$(BIN_DIR)/projects

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X github.com/MiguelRodo/projects/pkg/version.Version=$(VERSION) \
                  -X github.com/MiguelRodo/projects/pkg/version.GitCommit=$(COMMIT) \
                  -X github.com/MiguelRodo/projects/pkg/version.BuildDate=$(DATE)"

.PHONY: all build clean test coverage lint fmt vet install help

all: fmt vet lint test build

## build: Compile projectctl and projects binaries into bin/
build:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(PROJECTCTL_BIN) ./cmd/projectctl
	$(GOBUILD) $(LDFLAGS) -o $(PROJECTS_BIN) ./cmd/projects

## test: Run unit tests with race detection and coverage
test:
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

## coverage: Display coverage report
coverage: test
	$(GOCMD) tool cover -func=coverage.txt

## lint: Run golangci-lint
lint:
	$(GOLINT) run ./...

## fmt: Format Go source code
fmt:
	$(GOFMT) ./...

## vet: Run go vet
vet:
	$(GOVET) ./...

## install: Install binaries to GOPATH/bin
install:
	$(GOCMD) install $(LDFLAGS) ./cmd/projectctl
	$(GOCMD) install $(LDFLAGS) ./cmd/projects

## clean: Remove build artifacts and test reports
clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR) coverage.txt coverage.html

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'
