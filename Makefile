# gokit developer harness.
# Run `make check` before every commit / PR. CI runs the same targets.

GO        ?= go
PKGS      ?= ./...
GOFILES    = $(shell find . -name '*.go' -not -path './vendor/*')

.DEFAULT_GOAL := check

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## Format all Go files in place
	gofmt -w $(GOFILES)

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofmt-clean
	@out="$$(gofmt -l $(GOFILES))"; \
	if [ -n "$$out" ]; then echo "Not gofmt-clean:"; echo "$$out"; exit 1; fi

.PHONY: tidy
tidy: ## Sync go.mod / go.sum
	$(GO) mod tidy

.PHONY: vet
vet: ## Run go vet
	$(GO) vet $(PKGS)

.PHONY: lint
lint: ## Run golangci-lint if installed (skips gracefully otherwise)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run $(PKGS); \
	else \
		echo "golangci-lint not installed; skipping (go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)"; \
	fi

.PHONY: test
test: ## Run tests with race detector and coverage
	$(GO) test -race -cover $(PKGS)

.PHONY: cover
cover: ## Write and open an HTML coverage report
	$(GO) test -coverprofile=coverage.out $(PKGS)
	$(GO) tool cover -html=coverage.out

.PHONY: build
build: ## Compile all packages
	$(GO) build $(PKGS)

.PHONY: check
check: fmt-check vet lint test ## Full pre-commit gate: format + vet + lint + test

.PHONY: clean
clean: ## Remove build/coverage artifacts
	rm -f coverage.out
	$(GO) clean
