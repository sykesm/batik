// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package matcher

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
)

// EqualProto is a gomega matcher for protobuf messages.
func EqualProto(expected proto.Message) types.GomegaMatcher {
	return &equalProtoMatcher{
		expected: expected,
	}
}

type equalProtoMatcher struct {
	expected proto.Message
}

func (e *equalProtoMatcher) Match(actual interface{}) (bool, error) {
	actualMessage, ok := actual.(proto.Message)
	if !ok {
		return false, fmt.Errorf("EqualsProto expects a proto.Message")
	}

	return proto.Equal(e.expected, actualMessage), nil
}

func (e *equalProtoMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#vto proto.Equal\n\t%#v", actual, e.expected)
}

func (e *equalProtoMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#vnot to proto.Equal\n\t%#v", actual, e.expected)
}
