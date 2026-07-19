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

.PHONY: generate
generate: ## Regenerate mocks
	go tool mockery

.PHONY: gen
gen: ## Generate a CRUD module: make gen name=article
	go run ./cmd/gen $(name)


.PHONY: gen-check
gen-check: ## Verify generator output still compiles
	go run ./cmd/gen gencheck
	go build ./...
	rm -rf internal/gencheck

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

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
