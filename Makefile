APP_VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(APP_VERSION)

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: init
init: ## Copy the example config
	cp -n config/config.example.yml config/config.yml || true

.PHONY: tidy
tidy: ## Tidy go.mod
	go mod tidy

.PHONY: wire
wire: ## Regenerate compile-time dependency injection
	go tool wire generate ./...

.PHONY: wire-check
wire-check: ## Verify generated dependency injection is current
	go tool wire check ./...

.PHONY: generate
generate: wire ## Regenerate dependency injection and mocks
	go tool mockery

.PHONY: gen
gen: ## Generate a CRUD module: make gen name=article
	go run ./cmd/gen $(name)


.PHONY: gen-check
gen-check: ## Verify generator output still compiles
	go run ./cmd/gen gencheck
	@set -e; \
	  cp internal/app/wire.go internal/app/wire.go.gencheck; \
	  restore() { \
	    if [ -f internal/app/wire.go.gencheck ]; then mv -f internal/app/wire.go.gencheck internal/app/wire.go; fi; \
	    rm -f internal/app/wire.go.tmp; \
	    rm -rf internal/gencheck; \
	    go tool wire generate ./... >/dev/null; \
	  }; \
	  trap 'status=$$?; trap - EXIT INT TERM; restore; exit $$status' EXIT INT TERM; \
	  awk '1; /"github.com\/libtnb\/fiber-skeleton\/internal\/order"/ { print "\t\"github.com/libtnb/fiber-skeleton/internal/gencheck\"" }' internal/app/wire.go > internal/app/wire.go.tmp; \
	  mv internal/app/wire.go.tmp internal/app/wire.go; \
	  awk '1; /^[[:space:]]*order\.Module,$$/ { print "\t\tgencheck.Module," }' internal/app/wire.go > internal/app/wire.go.tmp; \
	  mv internal/app/wire.go.tmp internal/app/wire.go; \
	  go tool wire generate ./...; \
	  go build ./...

.PHONY: lint
lint: ## Run golangci-lint
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --timeout=30m ./...

.PHONY: test
test: ## Run tests with race detector and coverage
	go test -race -coverprofile=coverage.out ./...

.PHONY: build
build: ## Build app and cli binaries into bin/
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "$(LDFLAGS)" -o bin/app ./cmd/app
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "$(LDFLAGS)" -o bin/cli ./cmd/cli

.PHONY: run
run: ## Run the HTTP server
	go run -ldflags "-X main.version=$(APP_VERSION)" ./cmd/app

.PHONY: dev
dev: ## Run with hot reload (requires air)
	air
