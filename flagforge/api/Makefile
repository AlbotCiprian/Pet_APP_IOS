BINARY?=bin/api
GOFLAGS?=

.PHONY: all build run test fmt lint tidy

all: build

build:
GOFLAGS="$(GOFLAGS)" go build -o $(BINARY) ./cmd/api

run:
go run ./cmd/api

test:
go test ./...

fmt:
go fmt ./...

lint:
golangci-lint run ./...

tidy:
go mod tidy
