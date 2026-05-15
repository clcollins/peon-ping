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

PEON_DIR     ?= $(HOME)/.claude/hooks/peon-ping
PACKS_DIR    ?= $(PEON_DIR)/packs
PACKS_REPO   ?= https://github.com/PeonPing/og-packs.git
PACKS_REF    ?= v1.1.0

.PHONY: install
install: build ## Build and install peon binary to ~/.local/bin
	mkdir -p $(HOME)/.local/bin
	cp bin/peon $(HOME)/.local/bin/peon

SOUNDPACK ?=

.PHONY: install-pack
install-pack: ## Install a sound pack from upstream (SOUNDPACK=peon)
	@if [ -z "$(SOUNDPACK)" ]; then \
		echo "Usage: make install-pack SOUNDPACK=<name>"; \
		echo ""; \
		echo "Available packs (from $(PACKS_REPO) @ $(PACKS_REF)):"; \
		git ls-remote --refs --tags $(PACKS_REPO) 2>/dev/null | head -1 >/dev/null; \
		git clone --depth=1 --branch=$(PACKS_REF) $(PACKS_REPO) /tmp/og-packs-list 2>/dev/null; \
		ls -1 /tmp/og-packs-list/ | grep -v -E '^\.' | grep -v LICENSE | grep -v README; \
		rm -rf /tmp/og-packs-list; \
		exit 1; \
	fi
	@echo "Installing pack: $(SOUNDPACK)"
	git clone --depth=1 --branch=$(PACKS_REF) $(PACKS_REPO) /tmp/og-packs-dl 2>/dev/null
	@if [ ! -d "/tmp/og-packs-dl/$(SOUNDPACK)" ]; then \
		echo "Error: pack '$(SOUNDPACK)' not found in $(PACKS_REPO) @ $(PACKS_REF)"; \
		rm -rf /tmp/og-packs-dl; \
		exit 1; \
	fi
	mkdir -p $(PACKS_DIR)/$(SOUNDPACK)
	cp -r /tmp/og-packs-dl/$(SOUNDPACK)/* $(PACKS_DIR)/$(SOUNDPACK)/
	@if [ -f "$(PACKS_DIR)/$(SOUNDPACK)/openpeon.json" ]; then \
		mv $(PACKS_DIR)/$(SOUNDPACK)/openpeon.json $(PACKS_DIR)/$(SOUNDPACK)/manifest.json; \
	fi
	rm -rf /tmp/og-packs-dl
	@echo "Installed $(SOUNDPACK) to $(PACKS_DIR)/$(SOUNDPACK)"

.PHONY: install-default-packs
install-default-packs: ## Install all default sound packs (peon, peasant, sc_kerrigan, sc_battlecruiser, glados)
	git clone --depth=1 --branch=$(PACKS_REF) $(PACKS_REPO) /tmp/og-packs-dl 2>/dev/null
	@for pack in peon peasant sc_kerrigan sc_battlecruiser glados; do \
		echo "Installing pack: $$pack"; \
		mkdir -p $(PACKS_DIR)/$$pack; \
		cp -r /tmp/og-packs-dl/$$pack/* $(PACKS_DIR)/$$pack/; \
		if [ -f "$(PACKS_DIR)/$$pack/openpeon.json" ]; then \
			mv $(PACKS_DIR)/$$pack/openpeon.json $(PACKS_DIR)/$$pack/manifest.json; \
		fi; \
	done
	rm -rf /tmp/og-packs-dl
	@echo "Installed 5 default packs to $(PACKS_DIR)"

.PHONY: list-packs
list-packs: ## List installed sound packs
	@if [ -d "$(PACKS_DIR)" ]; then \
		for pack in $(PACKS_DIR)/*/; do \
			if [ -f "$$pack/manifest.json" ]; then \
				basename "$$pack"; \
			fi; \
		done; \
	else \
		echo "No packs installed ($(PACKS_DIR) not found)"; \
	fi

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
