.DEFAULT_GOAL:=help
-include .makerc

# --- Config -----------------------------------------------------------------

# Newline hack for error output
define br


endef

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%: .mise .husky
	@:

.PHONY: .mise
# Install dependencies
.mise: msg := $(br)$(br)Please ensure you have 'mise' installed and activated!$(br)$(br)$$ brew update$(br)$$ brew install mise$(br)$(br)See the documentation: https://mise.jdx.dev/getting-started.html$(br)$(br)
.mise:
ifeq (, $(shell command -v mise))
	$(error ${msg})
endif
	@mise install

.PHONY: .husky
# Configure git hooks for husky
.husky:
	@git config core.hooksPath .husky

### Tasks

.PHONY: check
## Run tidy, lint & tests
check: tidy lint test

.PHONY: tidy
## Run go mod tidy
tidy:
	@go mod tidy

.PHONY: lint
## Run linter
lint:
	@golangci-lint run

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	@golangci-lint run --fix

.PHONY: testapi
## Run api tests
testapi:
	@go run ./main.go -exportfile ./test/test-space-export.json ./test/testapi

.PHONY: test
## Run tests
test: testapi
	@echo "〉go test"
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out -tags=safe -race ./...

.PHONY: coverage
## Test & view coverage
coverage: test
	@go tool cover -html=coverage.out -o coverage.html; open coverage.html

.PHONY: build
## Build binary
build:
	@echo "〉building bin/gocontenful"
	@go build -o bin/gocontenful main.go

.PHONY: install
## Install binary
install:
	@echo "〉installing ${GOPATH}/bin/gocontentful"
	@go build -o ${GOPATH}/bin/gocontentful main.go

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@echo "〉go mod outdated"
	@go list -u -m -json all | go-mod-outdated -update -direct

### Utils

.PHONY: docs
## Open go docs
docs:
	@echo "〉starting go docs"
	@go doc -http

.PHONY: help
## Show help text
help:
	@echo "GoContentful - A Contentful Entry-Reference Mapper for Go\n"
	@echo "Usage:\n  make [task]"
	@awk '{ \
		if($$0 ~ /^### /){ \
			if(help) printf "%-23s %s\n\n", cmd, help; help=""; \
			printf "\n%s:\n", substr($$0,5); \
		} else if($$0 ~ /^[a-zA-Z0-9._-]+:/){ \
			cmd = substr($$0, 1, index($$0, ":")-1); \
			if(help) printf "  %-23s %s\n", cmd, help; help=""; \
		} else if($$0 ~ /^##/){ \
			help = help ? help "\n                        " substr($$0,3) : substr($$0,3); \
		} else if(help){ \
			print "\n                        " help "\n"; help=""; \
		} \
	}' $(MAKEFILE_LIST)


