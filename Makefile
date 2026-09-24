.DEFAULT_GOAL := help

export GOCACHE ?= $(CURDIR)/.cache/go-build

.PHONY: help setup fmt test lint build check docs-check run clean v2-setup v2-test v2-build v2-check v2-run godot-v2-import godot-v2-test godot-v2-smoke godot-v2-check godot-v2-run

help:
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "%-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

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

v2-setup: ## Install the Three.js V2 dependencies from the lockfile
	npm --prefix prototype-web ci

v2-test: ## Run the deterministic Three.js V2 rule tests
	npm --prefix prototype-web run test

v2-build: ## Type-check and build the Three.js V2 prototype
	npm --prefix prototype-web run build

v2-check: v2-test v2-build ## Run all Three.js V2 prototype checks

v2-run: ## Run the Three.js V2 prototype at http://127.0.0.1:4173
	npm --prefix prototype-web run dev -- --port 4173

godot-v2-import: ## Import the previous Godot V2 assets and resources
	godot --headless --path prototype-v2 --import

godot-v2-test: godot-v2-import ## Run the previous Godot V2 rule tests
	godot --headless --path prototype-v2 --script res://tests/test_encounter_rules.gd
	godot --headless --path prototype-v2 --script res://tests/test_scene_flow.gd

godot-v2-smoke: ## Load and run the previous Godot V2 scene headlessly
	godot --headless --path prototype-v2 --quit-after 2

godot-v2-check: godot-v2-test godot-v2-smoke ## Run all previous Godot V2 prototype checks

godot-v2-run: ## Run the previous Godot V2 prototype
	godot --path prototype-v2

clean: ## Remove generated build output
	go clean
	rm -f bin/valethra
