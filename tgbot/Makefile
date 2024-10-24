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


## TESTS

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: ./internal/services/business/iordercli.go \
	   ./internal/services/business/ipaymentcli.go \
	   ./internal/services/business/iproductcli.go
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
