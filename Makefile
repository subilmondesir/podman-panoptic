# Panoptic Makefile
# Build automation for cross-platform releases

VERSION := 1.0.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "initial")
DATE := $(shell date -u +%Y-%m-%d)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build binary for current platform
	go build -ldflags="$(LDFLAGS)" -o panoptic ./cmd/panoptic

.PHONY: build-all
build-all: ## Build for all platforms (Linux/Mac amd64/arm64)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/panoptic-linux-amd64 ./cmd/panoptic
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/panoptic-linux-arm64 ./cmd/panoptic
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/panoptic-darwin-amd64 ./cmd/panoptic
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/panoptic-darwin-arm64 ./cmd/panoptic

.PHONY: clean
clean: ## Clean build artifacts
	rm -f panoptic
	rm -rf dist/
	rm -f rapport.html

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: run
run: ## Run without building
	go run ./cmd/panoptic scan

.DEFAULT_GOAL := help