.DEFAULT_GOAL:=help
-include .makerc

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%:
	@:

## === Tasks ===

.PHONY: doc
## Run tests
doc:
	@open "http://localhost:6060/pkg/github.com/foomo/contentful/"
	@godoc -http=localhost:6060 -play

.PHONY: install
## Install binary
install:
	@go build -o ${GOPATH}/bin/gocontenful main.go

.PHONY: build
## Build binary
build:
	@mkdir -p bin
	@go build -o bin/gocontenful main.go

.PHONY: test
## Run tests
test: testapi
	@go test -p 1 -coverprofile=coverage.out -race -json ./... | gotestfmt

.PHONY: test
## Run tests
testapi:
	@go run ./main.go -exportfile ./test/test-space-export.json ./test/testapi

## Test & view coverage
cover: test
	@go tool cover -html=coverage.out -o coverage.html; open coverage.html

.PHONY: lint
## Run linter
lint: testapi
	@golangci-lint run

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	@golangci-lint run --fix

.PHONY: tidy
## Run go mod tidy
tidy:
	@go mod tidy

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@go list -u -m -json all | go-mod-outdated -update -direct

## === Utils ===

## Show help text
help:
	@awk '{ \
			if ($$0 ~ /^.PHONY: [a-zA-Z\-\_0-9]+$$/) { \
				helpCommand = substr($$0, index($$0, ":") + 2); \
				if (helpMessage) { \
					printf "\033[36m%-23s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^[a-zA-Z\-\_0-9.]+:/) { \
				helpCommand = substr($$0, 0, index($$0, ":")); \
				if (helpMessage) { \
					printf "\033[36m%-23s\033[0m %s\n", \
						helpCommand, helpMessage"\n"; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^##/) { \
				if (helpMessage) { \
					helpMessage = helpMessage"\n                        "substr($$0, 3); \
				} else { \
					helpMessage = substr($$0, 3); \
				} \
			} else { \
				if (helpMessage) { \
					print "\n                        "helpMessage"\n" \
				} \
				helpMessage = ""; \
			} \
		}' \
		$(MAKEFILE_LIST)
