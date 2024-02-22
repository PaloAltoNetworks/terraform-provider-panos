default: test

test:
	go test ./...

verbose:
	go test -v ./...

bench:
	go test -v -bench . ./...

fmt:
	find . -type f -name \*.go | xargs gofmt -w

docs:
	go generate

.PHONY: test verbose bench fmt docs
