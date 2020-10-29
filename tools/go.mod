// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

module tools

go 1.14

// https://github.com/uber/prototool/issues/559
// https://github.com/fullstorydev/grpcurl/pull/145
replace github.com/fullstorydev/grpcurl => github.com/fullstorydev/grpcurl v1.7.0

require (
	github.com/fullstorydev/grpcurl v1.7.0 // indirect
	github.com/gobuffalo/flect v0.2.2 // indirect
	github.com/gofrs/flock v0.8.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0
	github.com/jhump/protoreflect v1.7.0 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.3.0
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/nxadm/tail v1.4.5 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/spf13/cobra v1.1.0 // indirect
	github.com/uber/prototool v1.10.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect
	google.golang.org/genproto v0.0.0-20201015140912-32ed001d685c // indirect
	google.golang.org/grpc v1.33.0 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.0
	google.golang.org/protobuf v1.25.0
)
