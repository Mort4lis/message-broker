MAKEFLAGS = --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

BIN_NAME := message-broker
LINTER_VERSION := v1.62.2

.PHONY: lint.install
lint.install:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin ${LINTER_VERSION}

.PHONY: lint
lint:
	golangci-lint --version
	golangci-lint linters
	golangci-lint run -v

.PHONY: lint.fix
lint.fix:
	golangci-lint run --fix

.PHONY: test
test:
	go test -covermode=count -coverprofile=overalls.coverprofile -p 2 -count=1 ./...
	go tool cover -func=overalls.coverprofile
	go tool cover -html=overalls.coverprofile

.PHONY: build
build:
	go build -o ${BIN_NAME} cmd/message-broker

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	rm -f ${BIN_NAME} *.coverprofile coverage.*