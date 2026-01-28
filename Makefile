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

.PHONY: mongodb-status
mongodb-status: ## Check local MongoDB status
	@echo "${GREEN}Checking MongoDB status...${NC}"
	@mongosh --eval "db.adminCommand('ping')" --quiet > /dev/null 2>&1 && \
		echo "${GREEN}✓ MongoDB is running${NC}" || \
		echo "${YELLOW}✗ MongoDB is not running${NC}"

.PHONY: mongodb-start
mongodb-start: ## Start local MongoDB
	@echo "${GREEN}Starting MongoDB...${NC}"
	@if command -v brew > /dev/null 2>&1; then \
		brew services start mongodb-community@7.0; \
	elif command -v systemctl > /dev/null 2>&1; then \
		sudo systemctl start mongod; \
	elif command -v sc > /dev/null 2>&1; then \
		sc start MongoDB; \
	else \
		echo "${YELLOW}Please start MongoDB manually${NC}"; \
	fi

.PHONY: mongodb-stop
mongodb-stop: ## Stop local MongoDB
	@echo "${GREEN}Stopping MongoDB...${NC}"
	@if command -v brew > /dev/null 2>&1; then \
		brew services stop mongodb-community@7.0; \
	elif command -v systemctl > /dev/null 2>&1; then \
		sudo systemctl stop mongod; \
	elif command -v sc > /dev/null 2>&1; then \
		sc stop MongoDB; \
	else \
		echo "${YELLOW}Please stop MongoDB manually${NC}"; \
	fi

.PHONY: mongodb-shell
mongodb-shell: ## Open MongoDB shell for test database
	@mongosh auth_service_test

.PHONY: mongodb-clean
mongodb-clean: ## Drop test database
	@echo "${GREEN}Dropping test database...${NC}"
	@mongosh auth_service_test --eval "db.dropDatabase()" --quiet

.PHONY: test-integration
test-integration: ## Run integration tests (requires local MongoDB)
	@echo "${GREEN}Running integration tests...${NC}"
	@echo "${YELLOW}Make sure MongoDB is running! (Run: make mongodb-status)${NC}"
	$(GOTEST) -v -race -cover ./test/integration/...

.PHONY: test-integration-clean
test-integration-clean: ## Run integration tests with clean database
	@echo "${GREEN}Cleaning test database...${NC}"
	@mongosh auth_service_test --eval "db.dropDatabase()" --quiet
	@echo "${GREEN}Running integration tests...${NC}"
	$(GOTEST) -v -race -cover ./test/integration/...

.PHONY: test-all
test-all: test test-integration ## Run all tests (unit + integration)

.PHONY: docker-test-up
docker-test-up: ## Start Docker test MongoDB
	@echo "${GREEN}Starting Docker test MongoDB...${NC}"
	@docker-compose -f docker/docker-compose.test.yml up -d --wait
	@echo "${GREEN}MongoDB ready!${NC}"

.PHONY: docker-test-down
docker-test-down: ## Stop Docker test MongoDB
	@echo "${GREEN}Stopping Docker test MongoDB...${NC}"
	@docker-compose -f docker/docker-compose.test.yml down -v

.PHONY: docker-test-logs
docker-test-logs: ## View Docker MongoDB logs
	@docker-compose -f docker/docker-compose.test.yml logs -f