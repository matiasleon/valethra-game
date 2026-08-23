.DEFAULT_GOAL := help

export GOCACHE ?= $(CURDIR)/.cache/go-build

.PHONY: help setup fmt test lint build check docs-check run clean

help:
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "%-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Download and verify Go dependencies
	go mod download
	go mod verify

fmt: ## Format all Go source files
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

test: ## Run all tests with race detection
	go test -race ./...

lint: ## Check formatting, static analysis, and module consistency
	test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './.git/*'))"
	go vet ./...
	go mod tidy -diff

build: ## Build the desktop executable
	mkdir -p bin
	go build -trimpath -o bin/valethra .

docs-check: ## Reject stale product and editor-specific names
	@! rg -n -i 'saturday-chill|cursor|CLAUDE\.md' --glob '*.md' --glob '!docs/archive/**' .

check: lint test build docs-check ## Run the complete local quality gate

run: ## Run the desktop game
	go run .

clean: ## Remove generated build output
	go clean
	rm -f bin/valethra
