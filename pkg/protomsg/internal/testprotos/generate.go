// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testprotos

//go:generate protoc -I . --go_out=. --go_opt=paths=source_relative test2/test.proto
//go:generate protoc --experimental_allow_proto3_optional -I . --go_out=. --go_opt=paths=source_relative test3/test.proto test3/test_import.proto
