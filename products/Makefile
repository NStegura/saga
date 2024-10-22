$(eval CNT_API := $(shell docker ps -f name=products-api -q | wc -l | awk '{print $1}'))

.PHONY: up
up:
ifeq ($(CNT_API),0)
	docker-compose up --build --no-recreate --detach; sleep 5
endif

.PHONY: bash
bash:
	make up
	docker-compose exec products-api /bin/sh

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
	protoc --proto_path=./api --go_out=pkg/api --go_opt=paths=source_relative --go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative api/productsapi.proto

## TESTS

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: ./internal/app/server/iproduct.go \
	   ./internal/app/server/isystem.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@mkdir -p $(MOCKS_DESTINATION)
	@for file in $^; do \
  		input_path=$$file; \
  		if  [ $file="internal/"* ]; then \
  		  output_path=$$(echo "$$input_path" | sed 's/internal\///'); \
  		else \
  		  output_path=$input_path; \
  		fi; \
  		mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$output_path; \
  	done;

.PHONY: test
test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format pkgname -- -coverprofile=cover.out ./...

.PHONY: cover
cover:
	go test -v -coverpkg=./... -coverprofile=cover.out.tmp ./...
	cat cover.out.tmp | \
		grep -v "_easyjson.go" | \
		grep -v "/mocks/" | \
		grep -v "/storage/" | \
		grep -v ".pb.go" \
		> cover.out
	rm cover.out.tmp
	go tool cover -func cover.out

## LINTERS
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w -local github.com/NStegura/saga/products ./cmd
	goimports -w -local github.com/NStegura/saga/products ./internal
	goimports -w -local github.com/NStegura/saga/products ./config

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml --out-format=colored-line-number --sort-results
