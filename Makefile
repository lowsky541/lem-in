
GO ?= go
BIN_DIR := bin
DIST_DIR := dist

GO_DIRS := core util
GO_SOURCES := $(shell find $(GO_DIRS) -type f -name '*.go')

# WASM target
WASM_MAIN := wasm/main.go
WASM_OUT := $(DIST_DIR)/lem-in.wasm

# CLI target
CLI_MAIN := cli/main.go
CLI_OUT := $(BIN_DIR)/lem-in

# ----------------------
# Default target
# ----------------------
.PHONY: all
all: build

# ----------------------
# Test
# ----------------------
.PHONY: test
test:
	$(GO) test -v ./tests/...

# ----------------------
# CLI Build (default)
# ----------------------
.PHONY: build
build: $(CLI_OUT)

$(CLI_OUT): $(GO_SOURCES) $(CLI_MAIN) | $(BIN_DIR)
	$(GO) build -o $@ $(CLI_MAIN)

# ----------------------
# WASM build
# ----------------------
.PHONY: build-wasm
build-wasm: $(WASM_OUT)

$(WASM_OUT): $(GO_SOURCES) $(WASM_MAIN) | $(BIN_DIR)
	GOOS=js GOARCH=wasm $(GO) build -o $@ $(WASM_MAIN)

# ----------------------
# Clean
# ----------------------
.PHONY: clean
clean:
	rm -rf $(CLI_OUT) $(WASM_OUT)

# Create bin directory if it doesn't exist
$(BIN_DIR):
	mkdir -p $(BIN_DIR)
