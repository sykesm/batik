# Copyright IBM Corp. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

VERSION ?= "dev"
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%:z)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

all: batik checks

.PHONY: batik
batik:
	-mkdir -p dist
	go build \
		-ldflags "\
		-X \"github.com/sykesm/batik/pkg/buildinfo.Version=$(VERSION)\" \
		-X \"github.com/sykesm/batik/pkg/buildinfo.GitCommit=$(GIT_COMMIT)\" \
		-X \"github.com/sykesm/batik/pkg/buildinfo.Timestamp=$(BUILD_TIME)\" \
		" \
		-o dist/$@ \
		github.com/sykesm/batik/cmd/batik

checks: linting unit-test

.PHONY: unit-test unit-tests
unit-tests: unit-test
unit-test:
	scripts/run-unit-tests

.PHONY: lintint
linting:
	scripts/run-linting
