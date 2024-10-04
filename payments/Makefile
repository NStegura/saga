## LINTERS
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w -local github.com/NStegura/metrics ./cmd
	goimports -w -local github.com/NStegura/metrics ./internal
	goimports -w -local github.com/NStegura/metrics ./config

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml --out-format=colored-line-number --sort-results
