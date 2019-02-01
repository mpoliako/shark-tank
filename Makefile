GO_PACKAGES=$(shell go list ./... | grep -v /vendor/)
SHELL=/usr/bin/env bash -o pipefail
BINARIES = $(shell ls cmd)
BIN_PREFIX = linux_amd64/
BIN_FOLDER = bin
ifeq ($(shell go env GOOS),linux)
BIN_PREFIX =
endif

.PHONY: vendor
vendor:
	@dep status -v
	@dep ensure -v -vendor-only

.PHONY: build
build:
	@rm -rf $(BIN_FOLDER)
	@for pkg in $(GO_PACKAGES); do \
		GOOS=linux CGO_ENABLED=0 go install -tags netgo --ldflags '-extldflags "-static" -X "main.version=$(v)"' $$pkg || exit 1; \
	done; \
	mkdir -p $(BIN_FOLDER)
	@for cmd in $(BINARIES); do \
		cp $(GOPATH)/bin/$(BIN_PREFIX)$$cmd $(BIN_FOLDER) || exit 1; \
	done; \
