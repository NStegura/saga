.PHONY: protoc
protoc:
	protoc --proto_path=./api --go_out=pkg/api --go_opt=paths=source_relative --go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative api/paymentsapi.proto


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
