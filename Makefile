GOPATH	?= $(shell go env GOPATH)
CURDIR	= $(shell go list -f '{{.Dir}}' ./...)
FILES	:= $$(find $(CURDIR) -name "*.go")

.PHONY: release upx simulation mockgen test check

default: check

release: check test simulation
	echo "Compiling for Darwin and Linux"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/c.linux
	go build -ldflags="-s -w" -o release/c.darwin

upx: release
	upx release/c.linux
	upx release/c.darwin

simulation:
	ginkgo -v -tags simulation

mockgen:
	mockgen -destination=mock_rsync.go -package=main -source=rsync.go

test:
	@go test -v ./... | sed /PASS/s//$(shell printf "\033[32mPASS\033[0m")/ | sed /FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/

check: vet fmtcheck spellcheck goword staticcheck lint gosec checksucc

vet:
	@echo "vet"
	@go vet -all

fmtcheck:
	@echo "fmtcheck"
	@command -v goimports > /dev/null 2>&1 || GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	@CHANGES="$$(goimports -d $(CURDIR))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run goimports -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi
	@# Annoyingly, goimports does not support the simplify flag.
	@CHANGES="$$(gofmt -s -d $(CURDIR))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run gofmt -s -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi

spellcheck:
	@echo "spellcheck"
	@command -v misspell > /dev/null 2>&1 || GO111MODULE=off go get github.com/client9/misspell/cmd/misspell
	@misspell -locale="US" -error -source="text" **/*

goword:
	@echo "goword"
	@command -v goword > /dev/null 2>&1 || GO111MODULE=off go get github.com/chzchzchz/goword
	@goword $(FILES) 2>&1

staticcheck:
	@echo "staticcheck"
	@command -v staticcheck > /dev/null 2>&1 || GO111MODULE=off go get honnef.co/go/tools/cmd/staticcheck
	@staticcheck -checks="all" -tests $(CURDIR)

lint:
	@echo "lint"
	@command -v golangci-lint > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.27.0
	@golangci-lint run -v --disable-all --deadline=3m \
		--enable=misspell \
	  	--enable=ineffassign \
		--enable=errcheck \
	  	$$($(CURDIR))

gosec:
	@echo "gosec"
	@command -v gosec > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOPATH)/bin v2.3.0
	@gosec -exclude=G204 $(CURDIR)

checksucc:
	@echo "check successfully!"