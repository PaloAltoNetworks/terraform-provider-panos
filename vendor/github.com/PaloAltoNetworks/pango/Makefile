default: test

test:
	go test ./...

verbose:
	go test -v ./...

bench:
	go test -v -bench . ./...

fmt:
	find . -type f -name \*.go | xargs gofmt -w

.PHONY: test verbose bench fmt
