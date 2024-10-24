## TESTS

.PHONY: test
test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format pkgname -- -coverprofile=cover.out ./...

.PHONY: cover
cover:
	go test -v -coverpkg=./... -coverprofile=cover.out.tmp ./...
	cat cover.out.tmp | \
		grep -v ".pb.go" | \
		grep -v "/cmd/" \
		> cover.out
	rm cover.out.tmp
	go tool cover -func cover.out


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
