## LINTERS
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w -local github.com/NStegura/saga/tgbot ./cmd
	goimports -w -local github.com/NStegura/saga/tgbot ./internal
	goimports -w -local github.com/NStegura/saga/tgbot ./config

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml --out-format=colored-line-number --sort-results
