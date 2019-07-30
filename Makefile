bin/gonews:
	go build -mod vendor -o bin/gonews

.PHONY: dep
dep:
	go mod tidy
	go mod vendor

.PHONY: test
test:
	go test -v -cover -mod vendor ./...

.PHONY: clean
clean:
	rm -rf bin
	go clean -mod vendor
