# Copyright IBM Corp. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

VERSION ?= "dev"
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%:z)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

GOTOOLS = protoc-gen-go
GOTOOLS_BINDIR = tools/bin

all: gotools batik checks

.PHONY: clean
clean:
	@-rm -rf dist

.PHONY: batik
batik:
	-mkdir -p dist
	@go build \
		-ldflags "\
		-X \"github.com/sykesm/batik/pkg/buildinfo.Version=$(VERSION)\" \
		-X \"github.com/sykesm/batik/pkg/buildinfo.GitCommit=$(GIT_COMMIT)\" \
		-X \"github.com/sykesm/batik/pkg/buildinfo.Timestamp=$(BUILD_TIME)\" \
		" \
		-o dist/$@ \
		github.com/sykesm/batik/cmd/batik

checks: gotools linting unit-test

.PHONY: unit-test unit-tests
unit-test unit-tests:
	@scripts/run-unit-tests

.PHONY: linting
linting:
	@scripts/run-linting

# go tool->path mapping
gotool.protoc-gen-go := github.com/golang/protobuf/protoc-gen-go

gotools: $(patsubst %,$(GOTOOLS_BINDIR)/%, $(GOTOOLS))

$(GOTOOLS_BINDIR)/%: tools/go.sum
	@echo "Building ${gotool.$*} -> $@"
	@cd tools && GOBIN=$(abspath $(GOTOOLS_BINDIR)) go install -tags tools ${gotool.$*}

.PHONY: gotools-clean
gotools-clean:
	@rm -rf $(GOTOOLS_BINDIR)
