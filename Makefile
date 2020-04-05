# Lint (https://github.com/golangci/golangci-lint)
LINTER_OPTIONS ?= run# Arguments to golangci-lint
LINTER_BINARY ?= golangci-lint# Name of the binary of golangci-lint
LINTER_VERSION ?= 1.24.0# Version of golangci-lint to use in CI

lint:
ifneq (1,$(shell $(LINTER_BINARY) version 2>&1 | grep -c $(LINTER_VERSION)))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(LINTER_VERSION)
endif
	$(LINTER_BINARY) $(LINTER_OPTIONS)
.PHONY: lint
