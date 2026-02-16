GO ?= go
GO_DIRS := core
GO_SOURCES := $(shell find $(GO_DIRS) -type f -name '*.go')

BIN_DIR := bin

.PHONY: all
all: build

.PHONY: test
test:
	@$(GO) test -v  ./tests/...

.PHONY: build
build: bin/lem-in

.PHONY: build-wasm
build-wasm: bin/lem-in.wasm

bin/lem-in.wasm: $(GO_SOURCES) wasm/main.go
	GOOS=js GOARCH=wasm $(GO) build -o=$@ wasm/main.go

bin/lem-in: $(GO_SOURCES) cli/main.go
	$(GO) build -o=$@ cli/main.go
