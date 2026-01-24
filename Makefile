# Variables
SERVICE_NAME=auth-service
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=$(GO) fmt

# Colors for output (make output pretty!)
GREEN=\033[0;32m
NC=\033[0m # No Color
YELLOW=\033[1;33m

.PHONY: help
help: ## Show this help message
	@echo '${YELLOW}Available commands:${NC}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}'

.PHONY: test
test: ## Run all tests
	@echo "${GREEN}Running all tests...${NC}"
	CGO_ENABLED=0 $(GOTEST) -v -cover ./...

.PHONY: test-unit
test-unit: ## Run only unit tests
	@echo "${GREEN}Running unit tests...${NC}"
	CGO_ENABLED=0 $(GOTEST) -v -cover ./test/unit/...

.PHONY: test-domain
test-domain: ## Run domain layer tests
	@echo "${GREEN}Running domain tests...${NC}"
	CGO_ENABLED=0 $(GOTEST) -v -cover ./test/unit/domain/...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "${GREEN}Running tests with coverage...${NC}"
	CGO_ENABLED=0 $(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report generated: coverage.html${NC}"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "${GREEN}Running tests (verbose)...${NC}"
	CGO_ENABLED=0 $(GOTEST) -v -cover -count=1 ./...

.PHONY: benchmark
benchmark: ## Run benchmark tests
	@echo "${GREEN}Running benchmarks...${NC}"
	$(GOTEST) -bench=. -benchmem -run=^$ ./...

.PHONY: test-watch
test-watch: ## Run tests in watch mode (requires: go install github.com/cespare/reflex@latest)
	@echo "${GREEN}Running tests in watch mode...${NC}"
	reflex -r '\.go$$' -s -- make test-unit

.PHONY: fmt
fmt: ## Format code
	@echo "${GREEN}Formatting code...${NC}"
	$(GOFMT) ./...

.PHONY: vet
vet: ## Run go vet
	@echo "${GREEN}Running go vet...${NC}"
	$(GOVET) ./...

.PHONY: lint
lint: vet ## Run linters (requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "${GREEN}Running linters...${NC}"
	golangci-lint run ./...

.PHONY: clean
clean: ## Clean test cache and generated files
	@echo "${GREEN}Cleaning...${NC}"
	$(GO) clean -testcache
	rm -f coverage.out coverage.html

.PHONY: deps
deps: ## Download dependencies
	@echo "${GREEN}Downloading dependencies...${NC}"
	$(GO) mod download
	$(GO) mod tidy

.PHONY: check
check: fmt vet test ## Run all checks (fmt, vet, test)
	@echo "${GREEN}All checks passed!${NC}"