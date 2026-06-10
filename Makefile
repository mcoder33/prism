# Local mirror of .github/workflows/ci.yml — `make ci` runs everything CI runs.

GOLANGCI_LINT_VERSION := v2.11.4

.PHONY: help build install test cover vet fmt lint vuln mod-check ci

help: ## list targets
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-12s %s\n", $$1, $$2}'

build: ## compile all packages
	go build ./...

install: ## build + put `prism` on PATH
	go install .

test: ## unit tests with race detector
	go test -race ./...

cover: ## tests with coverage summary
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out | tail -1

vet: ## go vet
	go vet ./...

fmt: ## gofmt all sources in place
	gofmt -w .

lint: ## golangci-lint (auto-installs pinned version if missing)
	@command -v golangci-lint >/dev/null || \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	golangci-lint run --timeout 5m

vuln: ## govulncheck (auto-installs if missing)
	@command -v govulncheck >/dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

mod-check: ## go.mod/go.sum drift check
	go mod tidy
	git diff --exit-code go.mod go.sum

ci: test vet lint vuln mod-check ## everything CI runs
