// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

module tools

go 1.14

// https://github.com/uber/prototool/issues/559
// https://github.com/fullstorydev/grpcurl/pull/145
replace github.com/fullstorydev/grpcurl => github.com/fullstorydev/grpcurl v1.5.1

require (
	github.com/uber/prototool v1.9.0
	google.golang.org/protobuf v1.25.0
)
