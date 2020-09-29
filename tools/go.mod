// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

module tools

go 1.14

// https://github.com/uber/prototool/issues/559
// https://github.com/fullstorydev/grpcurl/pull/145
replace github.com/fullstorydev/grpcurl => github.com/fullstorydev/grpcurl v1.7.0

replace go.uber.org/zap => go.uber.org/zap v1.14.1

require (
	github.com/grpc-ecosystem/grpc-gateway v1.15.0
	github.com/onsi/ginkgo v1.14.0
	github.com/uber/prototool v1.10.0
	google.golang.org/protobuf v1.25.0
)
