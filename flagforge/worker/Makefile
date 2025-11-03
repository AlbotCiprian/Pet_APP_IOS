BINARY?=bin/worker
GOFLAGS?=

.PHONY: all build run test fmt lint tidy

all: build

build:
GOFLAGS="$(GOFLAGS)" go build -o $(BINARY) ./cmd/worker

run:
go run ./cmd/worker

test:
go test ./...

fmt:
go fmt ./...

lint:
golangci-lint run ./...

tidy:
go mod tidy
