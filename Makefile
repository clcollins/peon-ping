# PEON-PING — Claude Code sound notifications (personal fork)

CONTAINER_SUBSYS ?= podman
VERSION          ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

GO               := go
GOFMT            := gofmt
GOBIN            ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT    ?= $(shell command -v golangci-lint 2>/dev/null || echo "$(GOBIN)/golangci-lint")
CHECKMAKE        ?= $(shell command -v checkmake 2>/dev/null || echo "$(GOBIN)/checkmake")

# ── Build ──────────────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build the peon binary
	$(GO) build -ldflags="-X main.version=$(VERSION)" -o bin/peon ./cmd/peon/...

# ── Test ───────────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all Go tests with race detector
	$(GO) test -race -count=1 -timeout=5m ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	$(GO) test -race -count=1 -v -timeout=5m ./...

.PHONY: cover
cover: ## Run tests and emit coverage profile
	$(GO) test -race -count=1 -timeout=5m \
		-coverprofile=coverage.out -covermode=atomic ./...

# ── Code Quality ───────────────────────────────────────────────────────────────

.PHONY: fmt
fmt: ## Check formatting (gofmt)
	@echo "Checking gofmt..."
	@diff=$$($(GOFMT) -l .); if [ -n "$$diff" ]; then \
		echo "The following files are not formatted:"; echo "$$diff"; exit 1; \
	fi

.PHONY: fmt-fix
fmt-fix: ## Apply gofmt formatting
	$(GOFMT) -w .

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

MARKDOWNLINT_VERSION ?= 0.20.0
MARKDOWNLINT ?= $(shell command -v markdownlint-cli2 2>/dev/null \
	|| echo "npx --yes markdownlint-cli2@$(MARKDOWNLINT_VERSION)")

.PHONY: markdown-lint
markdown-lint: ## Lint all markdown files
	$(MARKDOWNLINT) "docs/**/*.md" "*.md"

.PHONY: makefile-lint
makefile-lint: $(CHECKMAKE) ## Lint this Makefile
	$(CHECKMAKE) Makefile

# ── Documentation ─────────────────────────────────────────────────────────────

.PHONY: docs-check
docs-check: ## Verify docs/ contains plan documents
	@count=$$(find docs -name '*.md' 2>/dev/null | wc -l); \
	[ "$$count" -gt 0 ] || \
		{ echo "ERROR: No plan documents in docs/."; exit 1; }; \
	echo "OK: $$count plan document(s) found."

# ── Install ───────────────────────────────────────────────────────────────────

.PHONY: install
install: build ## Build and install peon binary to ~/.local/bin
	mkdir -p $(HOME)/.local/bin
	cp bin/peon $(HOME)/.local/bin/peon

# ── Local CI ──────────────────────────────────────────────────────────────────

.PHONY: ci
ci: fmt vet lint test build docs-check ## Run all CI checks locally

CI_IMAGE ?= peon-ping-ci:local

.PHONY: ci-build
ci-build: ## Build the CI container image
	$(CONTAINER_SUBSYS) build -f test/Containerfile.ci -t $(CI_IMAGE) test/

.PHONY: ci-all
ci-all: ci-build ## Build CI container and run all checks inside it
	$(CONTAINER_SUBSYS) run --rm --userns=keep-id \
		-v $$(pwd):/src:Z -w /src $(CI_IMAGE) make ci

# ── Dependencies / Tools ──────────────────────────────────────────────────────

.PHONY: tidy
tidy: ## Run go mod tidy
	$(GO) mod tidy

.PHONY: tidy-check
tidy-check: ## Verify go.mod and go.sum are tidy
	$(GO) mod tidy
	git diff --exit-code go.mod
	@if [ -f go.sum ]; then git diff --exit-code go.sum; fi

$(GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3

$(CHECKMAKE):
	$(GO) install github.com/mrtazz/checkmake/cmd/checkmake@v0.2.2

# ── Clean ─────────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin/ coverage.out

# ── Help ──────────────────────────────────────────────────────────────────────

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' | sort
