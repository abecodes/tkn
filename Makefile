.DEFAULT_GOAL := help
# To avoid includes breaking the makefile_list, a copy is created
MAKE_FILE := $(lastword $(MAKEFILE_LIST))

# ====================
# Env variables
include .makerc

# ====================
# Create help output
help:
	@grep $(if $(filter Darwin,$(shell uname -s)), -E, -P) '^[a-zA-Z_-]+:.*?## .*$$' $(MAKE_FILE) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help
# ====================
# Formatting and linging

fmt: ## Format code with [gofumpt](https://github.com/mvdan/gofumpt)
	@gofumpt -d -w .
.PHONY: fmt

lint: fmt ## Lint code
	@revive -config revive.toml -formatter stylish -exclude ./vendor/... ./...
.PHONY: lint

vet: fmt ## Check parameters and assignments
	@go vet ./...
.PHONY: vet

shadow: fmt ## Check for shadowed variables
	@shadow ./...
.PHONY: vet

static: fmt ## Check common statics
	@staticcheck -f stylish ./...
.PHONY: vet

prepare: lint vet shadow static ## Execute all format and lint the code
.PHONY: prepare

# ====================
# Project

clean: ## Remove previously build binaries
	go clean
.PHONY: clean

install: ## Install required packages and vendor them
	@go mod tidy
	@go mod vendor
.PHONY: install

test: prepare ## Run all tests
	@go test ./... -count=1
.PHONY: test

# ====================
# Build

echo:
	@echo 'Building $(VERSION)...'
.PHONY: echo

prebuild:
	@$(MAKE) clean
	@$(MAKE) echo
	@$(MAKE) prepare -j$(shell sysctl -n hw.physicalcpu)
.PHONY: prebuild

build: ## Build the binary
	@$(MAKE) prebuild
	go build -o tkn cmd/main.go
.PHONY: build

