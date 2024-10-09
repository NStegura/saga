.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml --out-format=colored-line-number --sort-results
