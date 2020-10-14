# Copyright IBM Corp. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

VERSION ?= "dev"
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%:z)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

GOTOOLS_BINDIR = tools/bin

export PATH := $(GOTOOLS_BINDIR):$(PATH)

all: gotools batik checks

.PHONY: clean
clean:
	@-rm -rf dist

.PHONY: batik
batik:
	@-mkdir -p dist
	@go build \
		-ldflags "\
		-X \"github.com/sykesm/batik/pkg/buildinfo.Version=$(VERSION)\" \
		-X \"github.com/sykesm/batik/pkg/buildinfo.GitCommit=$(GIT_COMMIT)\" \
		-X \"github.com/sykesm/batik/pkg/buildinfo.Timestamp=$(BUILD_TIME)\" \
		" \
		-o dist/$@ \
		github.com/sykesm/batik/cmd/batik

checks: gotools linting unit-test integration-test

.PHONY: unit-test unit-tests
unit-test unit-tests:
	@scripts/run-unit-tests

.PHONY: linting
linting:
	@scripts/run-linting

.PHONY: integration-test integration-tests
integration-test integration-tests:
	@scripts/run-integration-tests

# go tool->path mapping
gotool.ginkgo := github.com/onsi/ginkgo/ginkgo
gotool.protoc-gen-go := google.golang.org/protobuf/cmd/protoc-gen-go
gotool.protoc-gen-go-grpc := google.golang.org/grpc/cmd/protoc-gen-go-grpc
gotool.protoc-gen-grpc-gateway := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
gotool.prototool := github.com/uber/prototool/cmd/prototool
GOTOOLS = ginkgo protoc-gen-go protoc-gen-go-grpc protoc-gen-grpc-gateway prototool

.PHONY: gotools
gotools: $(patsubst %,$(GOTOOLS_BINDIR)/%, $(GOTOOLS))

$(GOTOOLS_BINDIR)/%: tools/go.mod tools/go.sum
	@echo "Building ${gotool.$*} -> $@"
	@cd tools && GOBIN=$(abspath $(GOTOOLS_BINDIR)) go install -tags tools ${gotool.$*}

.PHONY: gotools-clean
gotools-clean:
	@rm -rf $(GOTOOLS_BINDIR)
