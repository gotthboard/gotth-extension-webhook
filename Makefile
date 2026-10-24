.PHONY: test race fuzz vet build format modules binding scripts closeout process artifact verify

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

fuzz:
	go test ./internal/webhook -run=^$$ -fuzz=FuzzDecodeConfig -fuzztime=2s
	go test ./internal/webhook -run=^$$ -fuzz=FuzzCanonicalAlert -fuzztime=2s

vet:
	go vet ./...

build:
	go build ./...

format:
	test -z "$$(gofmt -l cmd internal)"

modules:
	go mod verify

binding:
	./scripts/check-published-control.sh

scripts:
	sh -n scripts/*.sh
	python3 -c "import ast,pathlib; ast.parse(pathlib.Path('scripts/release-metadata.py').read_text())"
	python3 -c "import tomllib; tomllib.load(open('workflow.toml','rb'))"

closeout:
	./scripts/test-v1-release-closeout.sh

process:
	./scripts/test-process.sh

artifact:
	./scripts/test-reproducible-artifact.sh

verify: test race fuzz vet build format modules binding scripts closeout process artifact
	git diff --check
