V := @

BIN_DIR := ./bin
GO_NEWS := $(BIN_DIR)/gonews

bin/gonews:
	$(V)go build -mod vendor -o $(GO_NEWS)

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
