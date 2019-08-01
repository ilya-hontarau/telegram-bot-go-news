V := @

BIN_DIR := ./bin
GO_NEWS := $(BIN_DIR)/gonews

$(GO_NEWS):
	$(V)go build -mod vendor -ldflags "-linkmode external -extldflags -static" -o $(GO_NEWS)

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
