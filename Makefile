GO_BIN := go

.PHONY: test-fetch-deps
test-fetch-deps:
	@$(GO_BIN) install honnef.co/go/tools/cmd/staticcheck@2021.1.2
	@$(GO_BIN) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2

.PHONY: test
test: test-unit test-integration test-format test-lint test-security test-go-version

.PHONY: test-lint
test-lint:
	@echo $@
	@golangci-lint run

.PHONY: test-format
test-format:
	@echo $@
	@data=$$(gofmt -l .);\
		 if [ -n "$${data}" ]; then \
			>&2 echo "format is broken:"; \
			>&2 echo "$${data}"; \
			exit 1; \
		 fi

.PHONY: test-security
test-security:
	@echo $@
	@staticcheck ./...

.PHONY: test-go-version
test-go-version:
	@echo $@
	@$(GO_BIN) run ./cmd/assert-version go

.PHONY: test-integration
test-integration:
	@echo $@
	@$(GO_BIN) test ./... -run ^TestIntegration

.PHONY: test-unit
test-unit:
	@echo $@
	@$(GO_BIN) test -short ./...
