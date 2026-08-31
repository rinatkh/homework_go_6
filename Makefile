SHELL := /bin/bash
GO ?= go
GO_PACKAGES := ./...
UNIT_PACKAGES := ./internal/...
INTEGRATION_PACKAGES := ./test/integration/...
BIN_DIR := bin
COVERAGE_FILE ?= coverage.out
COVERAGE_THRESHOLD ?= 80.0
PACKAGE_FILE ?= $(BIN_DIR)/homework_go_6-linux-amd64.tar.gz
CMDS := 01_goroutines 02_channels 03_waitgroup 04_race_mutex 05_close_range 06_select 07_context 08_generics

.PHONY: help deps-check mod-check fmt fmt-check vet test test-unit test-integration test-race coverage coverage-check build package clean run-all ci compile test-goroutines test-channels test-waitgroup test-race-mutex test-close-range test-select test-context test-generics $(addprefix run-,$(CMDS))

help:
	@echo "Available commands:"
	@echo "  make compile          - compile all packages without running tests"
	@echo "  make test-goroutines  - run goroutine tasks"
	@echo "  make test-channels    - run channel tasks"
	@echo "  make test-waitgroup   - run WaitGroup tasks"
	@echo "  make test-race-mutex  - run shared-state tasks"
	@echo "  make test-close-range - run close/range tasks"
	@echo "  make test-select      - run select tasks"
	@echo "  make test-context     - run context tasks"
	@echo "  make test-generics    - run generics tasks"
	@echo "  make test-race        - run all unit tests with race detector"
	@echo "  make ci               - full local CI after solving all tasks"

deps-check:
	$(GO) mod download
	$(GO) mod verify

mod-check:
	$(GO) mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		git diff --exit-code -- go.mod; \
		if [ -f go.sum ]; then git diff --exit-code -- go.sum; fi; \
	else \
		echo "Skipping git diff because this directory is not a git repository"; \
	fi

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './$(BIN_DIR)/*')

fmt-check:
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './$(BIN_DIR)/*'))"; \
	if [ -n "$$files" ]; then echo "Go files are not formatted:"; echo "$$files"; exit 1; fi

vet:
	$(GO) vet $(GO_PACKAGES)

test: test-unit test-integration

test-unit:
	$(GO) test $(UNIT_PACKAGES)

test-integration:
	$(GO) test $(INTEGRATION_PACKAGES)

compile:
	$(GO) test -run '^$$' $(GO_PACKAGES)

test-goroutines:
	$(GO) test ./internal/goroutines/...

test-channels:
	$(GO) test ./internal/channels/...

test-waitgroup:
	$(GO) test ./internal/waitgroup/...

test-race-mutex:
	$(GO) test ./internal/racemutex/...

test-close-range:
	$(GO) test ./internal/closerange/...

test-select:
	$(GO) test ./internal/selectflow/...

test-context:
	$(GO) test ./internal/contextflow/...

test-generics:
	$(GO) test ./internal/generics/...

test-race:
	$(GO) test -race $(UNIT_PACKAGES)

coverage:
	$(GO) test $(UNIT_PACKAGES) -covermode=atomic -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-check: coverage
	@coverage="$$($(GO) tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {gsub("%", "", $$3); print $$3}')"; \
	awk -v coverage="$$coverage" -v threshold="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		if (coverage + 0 < threshold + 0) { printf "coverage %.1f%% is below threshold %.1f%%\n", coverage, threshold; exit 1 } \
		printf "coverage %.1f%% is enough; threshold %.1f%%\n", coverage, threshold; \
	}'

run-all:
	@for cmd in $(CMDS); do \
		echo "== $$cmd =="; \
		$(GO) run ./cmd/$$cmd; \
	done

run-01_goroutines:
	$(GO) run ./cmd/01_goroutines
run-02_channels:
	$(GO) run ./cmd/02_channels
run-03_waitgroup:
	$(GO) run ./cmd/03_waitgroup
run-04_race_mutex:
	$(GO) run ./cmd/04_race_mutex
run-05_close_range:
	$(GO) run ./cmd/05_close_range
run-06_select:
	$(GO) run ./cmd/06_select
run-07_context:
	$(GO) run ./cmd/07_context
run-08_generics:
	$(GO) run ./cmd/08_generics

build:
	@mkdir -p $(BIN_DIR)
	@for cmd in $(CMDS); do $(GO) build -o $(BIN_DIR)/$$cmd ./cmd/$$cmd; done

package: build
	@mkdir -p $(BIN_DIR)
	tar -czf $(PACKAGE_FILE) -C $(BIN_DIR) $(CMDS)

ci: deps-check mod-check fmt-check vet test-unit test-integration test-race coverage-check build package

clean:
	rm -rf $(BIN_DIR) $(COVERAGE_FILE)
