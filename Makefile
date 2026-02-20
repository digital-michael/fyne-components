.PHONY: help test lint fmt vet coverage deps clean tidy build build-examples build-table-demo

# Variables
GO := go
GOPATH := $(shell go env GOPATH)
GOFLAGS := -v
GOTEST := $(GO) test
GOVET := $(GO) vet
GOFMT := $(GO) fmt
GOLINT := $(GOPATH)/bin/golangci-lint
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Default target
help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

all: deps fmt vet lint test ## Run all checks and tests

deps: ## Download and tidy dependencies
	$(GO) mod download
	$(GO) mod tidy

tidy: ## Tidy go.mod and go.sum
	$(GO) mod tidy

fmt: ## Format Go source code
	$(GOFMT) ./...

vet: ## Run go vet on all packages
	$(GOVET) ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@which $(GOLINT) > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	$(GOLINT) run ./...

test: ## Run tests
	$(GOTEST) $(GOFLAGS) ./...

test-verbose: ## Run tests with verbose output
	$(GOTEST) -v ./...

build: build-examples ## Build all examples

build-examples: build-table-demo ## Build all example applications

build-table-demo: ## Build the table demo application
	@echo "Building table-demo..."
	@mkdir -p bin
	$(GO) build -o bin/table-demo ./examples/table-demo
	@echo "Built: bin/table-demo"

coverage: ## Generate test coverage report
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

coverage-text: ## Display test coverage in terminal
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -func=$(COVERAGE_FILE)

clean: ## Clean build artifacts and coverage reports
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -rf bin/ dist/
	$(GO) clean

check: fmt vet lint ## Run formatting, vetting, and linting

ci: deps check test ## Run CI pipeline (deps, checks, and tests)
