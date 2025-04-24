GO ?= go
GO_DIRS := pkg
GO_SOURCES := $(wildcard *.go)
GO_SOURCES += $(shell find $(GO_DIRS) -type f -name '*.go')
EXECUTABLE ?= bin/lem-in

NPM ?= npm
TS_DIRS := visualizer/src
VITE_CONFIGS := visualizer/tsconfig.json visualizer/vite.config.ts
VITE_SOURCES := visualizer/index.html visualizer/src/main.scss
VITE_SOURCES += $(shell find $(TS_DIRS) -type f -name '*.ts')
VITE_DIST_FILES := visualizer/dist/index.html visualizer/dist/assets/index.css \
	visualizer/dist/assets/index.js visualizer/dist/assets/index.js.map 

.PHONY: all
all: build

.PHONY: help
help:
	@echo "usage: make help|test|build|build-visualizer|watch-visualizer|build-lemin"

.PHONY: test
test:
	@go test  ./...

.PHONY: build
build: build-visualizer build-lemin

.PHONY: build-visualizer
build-visualizer: node_modules $(VITE_DIST_FILES)

.PHONY: build-lemin
build-lemin: $(EXECUTABLE)

.PHONY: watch-visualizer
watch-visualizer: node_modules
	@cd visualizer && npx vite build --watch

# Use --no-save to prevent touch-ing package-lock.json
node_modules: package-lock.json
	$(NPM) install --no-save

$(VITE_DIST_FILES): $(VITE_SOURCES) $(VITE_CONFIGS)
	@cd visualizer && npx vite build

$(EXECUTABLE): $(GO_SOURCES) $(VITE_DIST_FILES)
	$(GO) build -o=$@
