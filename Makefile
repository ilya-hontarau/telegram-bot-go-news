V := @

BIN_DIR := ./bin
GO_NEWS := $(BIN_DIR)/gonews

$(GO_NEWS):
	$(V)CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o $(GO_NEWS) ./cmd/gonews

.PHONY: dep
dep:
	$(V)go mod tidy
	$(V)go mod vendor

.PHONY: test
test:
	$(V)go test -v -mod vendor -coverprofile=cover.out -race ./...

.PHONY: lint
lint:
	$(V)golangci-lint run --config .golangci.local.yml

.PHONY: clean
clean:
	$(V)rm -rf bin
	$(V)go clean -mod vendor
