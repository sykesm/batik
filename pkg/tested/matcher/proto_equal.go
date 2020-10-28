// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package matcher

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
)

// ProtoEqual is a gomega matcher for protobuf messages.
func ProtoEqual(expected proto.Message) types.GomegaMatcher {
	return &protoEqualMatcher{
		expected: expected,
	}
}

type protoEqualMatcher struct {
	expected proto.Message
}

func (e *protoEqualMatcher) Match(actual interface{}) (bool, error) {
	actualMessage, ok := actual.(proto.Message)
	if !ok {
		return false, fmt.Errorf("ProtoEqual expects a proto.Message")
	}

	return proto.Equal(e.expected, actualMessage), nil
}

func (e *protoEqualMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nto proto.Equal\n\t%#v", actual, e.expected)
}

func (e *protoEqualMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nnot to proto.Equal\n\t%#v", actual, e.expected)
}
