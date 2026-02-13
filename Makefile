.PHONY: fmt test vet lint check

fmt:
	gofmt -w ./internal/licensing ./ilicense ./examples

test:
	GOCACHE=/tmp/go-build-cache go test ./...

vet:
	GOCACHE=/tmp/go-build-cache go vet ./...

lint: fmt vet test

check: lint
