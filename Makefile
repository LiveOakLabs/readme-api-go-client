# Makefile for readme-api-go-client.
# Run 'make help' for a list of targets.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Shows this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## Linting ##
.PHONY: vet
vet: ## Runs 'go vet'
	@go vet ./...

.PHONY: gofumpt
gofumpt: vet ## Check linting with 'gofumpt'
	@go run mvdan.cc/gofumpt -l -d .

.PHONY: lines
lines: ## Check long lines.
	@go run github.com/segmentio/golines -m 120 --dry-run readme/*.go

.PHONY: lines-fix
lines-fix: lines ## Fix long lines
	@go run github.com/segmentio/golines -m 120 -w readme/*.go

.PHONY: golangci-lint
golangci-lint: ## Lint using 'golangci-lint'
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint \
	run --deadline=300s --timeout=300s --out-format checkstyle ./... 2>&1 | tee checkstyle-report.xml

.PHONY: lint
lint: gofumpt vet lines golangci-lint ## Run all linters

## Testing ##
.PHONY: test
test: ## Run unit and race tests with 'go test'
	@go test -count=1 -coverprofile=coverage.txt -covermode count ./readme/...
	@go test -race -short ./readme/...

## Coverage ##
.PHONY: coverage
coverage: test ## Generate a code test coverage report using 'gocover-cobertura'
	@go run github.com/boumenot/gocover-cobertura < coverage.txt > coverage.xml
	@rm -f coverage.txt

.PHONY: coverage-html
coverage-html: test ## Generate an HTML test coverage report
	@go tool cover -html=coverage.txt -o coverage.html
	@rm -f coverage.txt

## Vulnerability checks ##
.PHONY: check-vuln
check-vuln: ## Check for vulnerabilities using 'govulncheck'
	@echo "Checking for vulnerabilities..."
	@go run golang.org/x/vuln/cmd/govulncheck ./...

.PHONY: clean
clean: ## Clean test files
	@rm -f coverage.txt coverage.xml coverage.html checkstyle-report.xml

.PHONY: mocks
mocks: ## Generate mocks
	# @mockery --all --keeptree --output ./mocks --outpkg mocks --case underscore
	@mockery --all
