// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/protomsg/internal/testprotos/test2"
	"github.com/sykesm/batik/pkg/protomsg/internal/testprotos/test3"
)

func TestValidateSuccess(t *testing.T) {
	gt := NewGomegaWithT(t)
	err := protomsg.ValidateMessage(&test3.TestAllTypes{
		SingularString:         "a string",
		SingularNestedMessage:  &test3.TestAllTypes_NestedMessage{},
		RepeatedNestedMessage:  []*test3.TestAllTypes_NestedMessage{{}},
		MapStringString:        map[string]string{"key": "value"},
		MapStringNestedMessage: map[string]*test3.TestAllTypes_NestedMessage{"key": {}},
	})
	gt.Expect(err).NotTo(HaveOccurred())
}

func TestValidateNil(t *testing.T) {
	gt := NewGomegaWithT(t)
	err := protomsg.ValidateMessage(nil)
	gt.Expect(err).NotTo(HaveOccurred())
}

func TestValidateSyntaxProto2(t *testing.T) {
	gt := NewGomegaWithT(t)
	err := protomsg.ValidateMessage(&test2.TestMessage{})
	gt.Expect(err).To(MatchError("protomsg: proto2 syntax is not supported"))
}

func TestValidateUnknownFields(t *testing.T) {
	nested := &test3.TestAllTypes_NestedMessage{}
	nested.ProtoReflect().SetUnknown(protoreflect.RawFields("raw-fields"))

	tests := []struct {
		desc string
		f    func(*test3.TestAllTypes)
	}{
		{"top", func(m *test3.TestAllTypes) { m.ProtoReflect().SetUnknown(protoreflect.RawFields("raw-fields")) }},
		{"nested", func(m *test3.TestAllTypes) { m.SingularNestedMessage = nested }},
		{"list", func(m *test3.TestAllTypes) { m.RepeatedNestedMessage = []*test3.TestAllTypes_NestedMessage{nested} }},
		{"map", func(m *test3.TestAllTypes) {
			m.MapStringNestedMessage = map[string]*test3.TestAllTypes_NestedMessage{"key": nested}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			m := &test3.TestAllTypes{}

			tt.f(m)
			err := protomsg.ValidateMessage(m)
			gt.Expect(err).To(MatchError("protomsg: refusing to marshal unknown fields with length 10"))
		})
	}
}
