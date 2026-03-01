.PHONY: all build test lint generate-schema lint-schema clean

# Default target
all: build test lint

# Build all packages
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Generate JSON schemas from Go types
generate-schema: build
	go run ./tools/generate

# Lint generated schemas with schemago (if installed)
lint-schema: generate-schema
	@if command -v schemago >/dev/null 2>&1; then \
		echo "Linting schemas with schemago..."; \
		schemago lint schema/design-system.schema.json; \
	else \
		echo "schemago not installed, skipping schema lint"; \
	fi

# Clean generated files
clean:
	rm -rf schema/*.schema.json
	rm -rf schema/foundations/*.schema.json
	rm -rf schema/components/*.schema.json

# Format code
fmt:
	gofmt -w .
	cd sdk/go && gofmt -w .

# Tidy modules
tidy:
	go mod tidy
	cd sdk/go && go mod tidy

# Install dependencies
deps:
	go mod download
	cd sdk/go && go mod download

# Run schema generation and verify output
verify-schema: generate-schema
	@echo "Verifying schema files exist..."
	@test -f schema/design-system.schema.json || (echo "Missing design-system.schema.json" && exit 1)
	@test -f schema/meta.schema.json || (echo "Missing meta.schema.json" && exit 1)
	@test -f schema/foundations/foundations.schema.json || (echo "Missing foundations.schema.json" && exit 1)
	@test -f schema/components/component.schema.json || (echo "Missing component.schema.json" && exit 1)
	@echo "All schema files present!"
