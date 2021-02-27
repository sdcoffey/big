files := $(shell find . -name "*.go" | grep -v vendor)

bootstrap:
	go install -v golang.org/x/lint/golint@latest
	go install -v golang.org/x/tools/...@latest
	go install -v honnef.co/go/tools/cmd/staticcheck@latest

lint:
	golint -set_exit_status
	staticcheck github.com/sdcoffey/big

clean:
	goimports -w $(files)

test: clean
	go test -v

release: clean test
	./scripts/release.sh