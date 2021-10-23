.PHONY: lint
lint:
	${GOPATH}/bin/golangci-lint run

.PHONY: test
test:
	go test -v -race -timeout 30s ./internal/...

.PHONY: pipeline
pipeline:
	make lint && make test

.DEFAULT_GOAL := pipeline
