.DEFAULT_GOAL:=help

## === Tasks ===

## Install binary
install:
	go build -o ${GOPATH}/bin/gocontenful main.go

## Build binary
build:
	mkdir -p bin
	go build -o bin/gocontenful main.go

.PHONY: test

## Run tests
test:
	go run ./main.go -exportfile ./test/test-space-export.json ./test/testapi
	go test -count=1 ./...

race:
	go run ./main.go -exportfile ./test/test-space-export.json ./test/testapi
	go test -race ./...

.PHONY: lint
## Run linter
lint:
	golangci-lint run

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	golangci-lint run --fix


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
