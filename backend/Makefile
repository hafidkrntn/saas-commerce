# ==================== App Commands ====================

# Setup runs automatically on first `make run`, `make build`, or `make test`.
# It configures git to use the project's hooks directory directly.
setup:
	@git config core.hooksPath scripts/hooks
	@chmod +x scripts/hooks/*
	@echo "✅ Git hooks configured (scripts/hooks)"

run: setup
	go run ./cmd/api/

dev: setup
	air

build: setup
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd/api/

test: setup
	@echo "Running tests with gotestsum..."
	@gots=$$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...); \
	nots=$$(go list -f '{{if not (or .TestGoFiles .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./...); \
	gotestsum --format=testname -- -cover $$gots; status=$$?; \
	if [ -n "$$nots" ]; then gotestsum --format=testname -- $$nots; fi; \
	exit $$status

lint:
	golangci-lint run ./...

wire:
	@GOTOOLCHAIN=go1.25.7 go run github.com/google/wire/cmd/wire@latest gen ./cmd/api/

# ==================== Migration Commands (goose) ====================
# Requires: go install github.com/pressly/goose/v3/cmd/goose@latest
# DB_DSN must be set in .env or environment

GOOSE_CMD=goose

migrate-up:
	$(GOOSE_CMD) up

migrate-down:
	$(GOOSE_CMD) down

migrate-status:
	$(GOOSE_CMD) status

migrate-create:
	@if [ -z "$(name)" ]; then echo "❌ Usage: make migrate-create name=<migration_name>"; exit 1; fi
	$(GOOSE_CMD) create $(name) sql

migrate-redo:
	$(GOOSE_CMD) redo

migrate-reset:
	$(GOOSE_CMD) reset

# ==================== Module Scaffold ====================
# Usage: make new-module name=product [tests=true]
# Generates interfaces, mockgen mocks, and (optionally) unit tests.
# Default: tests=false (no test files generated).

new-module:
	@if [ -z "$(name)" ]; then echo "❌ Usage: make new-module name=<module_name> [tests=true]"; exit 1; fi
	@./scripts/new_module.sh "$(name)" "$(tests)"
	@echo "✅ Module '$(name)' created at internal/$(name)/"

# ==================== Docker ====================

docker-build:
	docker build -t $${APP_NAME:-my-service}:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

.PHONY: run build test lint wire migrate-up migrate-down migrate-status migrate-create migrate-redo migrate-reset new-module docker-build docker-up docker-down
