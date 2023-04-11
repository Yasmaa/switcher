
all: test

# Tool versions
GOLANGCI_LINT_VERSION?=v1.51.1

.PHONY: test
test:

.PHONY: test-go
test: test-go
test-go:
	go test -v ./...

.PHONY: lint
test: lint
lint: tool/golangci-lint
	tool/golangci-lint run

.PHONY: fix
fix: tool/golangci-lint
	tool/golangci-lint run --fix

.PHONY: test-bench
test: test-bench
test-bench:
	go test -bench=. -v ./...

.PHONY: test-fuzz
test: test-fuzz
test-fuzz:
	go test -fuzz -v ./...

tool/golangci-lint: tool/.golangci-lint.$(GOLANGCI_LINT_VERSION)
	@mkdir -p tool
	GOBIN="$(PWD)/tool" go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

tool/.golangci-lint.$(GOLANGCI_LINT_VERSION):
	@rm -f tool/.golangci-lint.*
	@mkdir -p tool
	touch $@

.PHONY: tool
tool: tool/golangci-lint
