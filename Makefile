V := @

BIN_DIR := ./bin
GO_NEWS := $(BIN_DIR)/gonews

$(GO_NEWS):
	$(V)go build -mod vendor -ldflags "-linkmode external -extldflags -static" -o $(GO_NEWS) ./cmd/gonews

.PHONY: dep
dep:
	$(V)go mod tidy
	$(V)go mod vendor

.PHONY: test
test:
	$(V)go test -v -cover -mod vendor ./...

.PHONY: clean
clean:
	$(V)rm -rf bin
	$(V)go clean -mod vendor
