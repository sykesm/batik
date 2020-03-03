# Copyright IBM Corp. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

VERSION ?= "dev"
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%:z)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

all: batik

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
