$(eval CNT_API := $(shell docker ps -f name=payments-api -q | wc -l | awk '{print $1}'))

.PHONY: up
up:
ifeq ($(CNT_API),0)
	docker-compose up --build --no-recreate --detach; sleep 5
endif

.PHONY: bash
bash:
	make up
	docker-compose exec payments-api /bin/sh

.PHONY: down
down:
	docker-compose down --remove-orphans --rmi local

.PHONY: buildall
buildall: buildapi buildcons

.PHONY: buildapi
buildapi: ## Build api app
	go build -o ./bin/server cmd/server/main.go

.PHONY: buildcons
buildcons: ## Build cons app
	go build -o ./bin/cons cmd/consumer/main.go

.PHONY: runapi
runapi: ## Run api app
	go run cmd/server/main.go


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
