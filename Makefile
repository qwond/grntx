.DEFAULT_GOAL := help

# include env file if exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Semver from tag
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Commit hash for builds identification
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Check if build is dirty
GIT_DIRTY := $(shell git diff --quiet || echo "-dirty")

# Full version string
VERSION := $(GIT_TAG)-$(GIT_COMMIT)$(GIT_DIRTY)

.PHONY: version
version:
	@echo "Version: $(VERSION)"

test: ## Runs tests
	go test ./...