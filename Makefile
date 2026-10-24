.PHONY: test race vet build format binding process artifact verify

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

vet:
	go vet ./...

build:
	go build ./...

format:
	test -z "$$(gofmt -l cmd internal)"

binding:
	./scripts/check-published-control.sh

process:
	./scripts/test-process.sh

artifact:
	./scripts/test-reproducible-artifact.sh

verify: test race vet build format binding process artifact
	git diff --check
