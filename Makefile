files := $(shell find . -name "*.go" | grep -v vendor)

GOLINT_VERSION := v0.0.0-20241112194109-818c5a804067
GOIMPORTS_VERSION := v0.24.0
STATICCHECK_VERSION := 2024.1.1

bootstrap:
	go install -v golang.org/x/lint/golint@$(GOLINT_VERSION)
	go install -v golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	go install -v honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)

lint:
	golint -set_exit_status
	staticcheck ./...

fmt:
	goimports -w $(files)

test:
	go test -v ./...

release: fmt test
	./scripts/release.sh
