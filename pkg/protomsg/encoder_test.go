// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/protomsg/internal/testprotos/test2"
	"github.com/sykesm/batik/pkg/protomsg/internal/testprotos/test3"
)

// If the encoding of a message using the official protocol buffer
// implemenatation does not match one produced by the custom encoding, that
// indicates an incompatible change has been made in the protocol buffer
// implementation. The custom encoding should not be changed to match the
// official implementation unless it can be done in a backwards compatible
// manner.

func TestEquality(t *testing.T) {
	gt := NewGomegaWithT(t)
	sr := &transaction.StateReference{
		Txid:        []byte("abakjfjdkjfkdjkajdfkdjfkdjkdjkjakfjd"),
		OutputIndex: 9999,
	}

	protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(sr)
	gt.Expect(err).NotTo(HaveOccurred())

	customBytes, err := protomsg.MarshalDeterministic(sr)
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(customBytes).To(Equal(protoBytes))
}

func TestEqualityGenerated(t *testing.T) {
	gt := NewGomegaWithT(t)
	m := makeProto()

	protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
	gt.Expect(err).NotTo(HaveOccurred())

	customBytes, err := protomsg.MarshalDeterministic(m)
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(customBytes).To(Equal(protoBytes))
}

func TestProto2Error(t *testing.T) {
	gt := NewGomegaWithT(t)
	m := &test2.TestMessage{}

	_, err := protomsg.MarshalDeterministic(m)
	gt.Expect(err).To(MatchError("protomsg: proto2 syntax is not supported"))
}

func TestDelegateToProto(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		var m *test3.TestAllTypes
		gt.Expect(m.ProtoReflect().IsValid()).To(BeFalse())

		res, err := protomsg.MarshalDeterministic(m)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(res).To(BeNil())
	})

	t.Run("Nil", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		res, err := protomsg.MarshalDeterministic(nil)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(res).To(BeNil())
	})
}

func TestUnknownFields(t *testing.T) {
	gt := NewGomegaWithT(t)

	m := &test3.TestAllTypes{}
	m.ProtoReflect().SetUnknown(protoreflect.RawFields("raw-fields"))

	_, err := protomsg.MarshalDeterministic(m)
	gt.Expect(err).To(MatchError("protomsg: refusing to marshal unknown fields with length 10"))
}

func TestSortedMap(t *testing.T) {
	tests := []struct {
		desc string
		m    *test3.TestAllTypes
	}{
		{"bool", &test3.TestAllTypes{MapBoolBool: map[bool]bool{true: true, false: false}}},
		{"int64", &test3.TestAllTypes{MapInt64Int64: map[int64]int64{1: 1, 2: 2}}},
		{"uint64", &test3.TestAllTypes{MapUint64Uint64: map[uint64]uint64{99: 99, 88: 88}}},
		{"string", &test3.TestAllTypes{MapStringString: map[string]string{"a": "Alice", "b": "Bob"}}},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(tt.m)
			gt.Expect(err).NotTo(HaveOccurred())

			customBytes, err := protomsg.MarshalDeterministic(tt.m)
			gt.Expect(err).NotTo(HaveOccurred())

			gt.Expect(customBytes).To(Equal(protoBytes))
		})
	}
}

func TestInvalidString(t *testing.T) {
	gt := NewGomegaWithT(t)

	m := &test3.TestAllTypes{SingularString: string([]byte{200, 200, 200})}
	_, err := protomsg.MarshalDeterministic(m)
	gt.Expect(err).To(MatchError("protomsg: field proto.test3.TestAllTypes.singular_string contains invalid UTF-8"))
}

func BenchmarkProtoMarshal(b *testing.B) {
	m := makeProto()
	for n := 0; n < b.N; n++ {
		out, err := proto.Marshal(m)
		if err != nil {
			b.Fatalf("marshal failed: %v", err)
		}
		if len(out) == 0 {
			b.Fatalf("output is empty")
		}
	}
}

func BenchmarkProtoDeterministicMarshal(b *testing.B) {
	m := makeProto()
	for n := 0; n < b.N; n++ {
		out, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
		if err != nil {
			b.Fatalf("marshal failed: %v", err)
		}
		if len(out) == 0 {
			b.Fatalf("output is empty")
		}
	}
}

func BenchmarkCustomDeterministicMarshal(b *testing.B) {
	m := makeProto()
	for n := 0; n < b.N; n++ {
		out, err := protomsg.MarshalDeterministic(m)
		if err != nil {
			b.Fatalf("marshal failed: %v", err)
		}
		if len(out) == 0 {
			b.Fatalf("output is empty")
		}
	}
}

const maxRecurseLevel = 3

func makeProto() *test3.TestAllTypes {
	m := &test3.TestAllTypes{}
	fillMessage(m.ProtoReflect(), 0)
	return m
}

func fillMessage(m protoreflect.Message, level int) {
	if level > maxRecurseLevel {
		return
	}

	fieldDescs := m.Descriptor().Fields()
	for i := 0; i < fieldDescs.Len(); i++ {
		fd := fieldDescs.Get(i)
		switch {
		case fd.IsList():
			setList(m.Mutable(fd).List(), fd, level)
		case fd.IsMap():
			setMap(m.Mutable(fd).Map(), fd, level)
		default:
			setScalarField(m, fd, level)
		}
	}
}

func setScalarField(m protoreflect.Message, fd protoreflect.FieldDescriptor, level int) {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		val := m.NewField(fd)
		fillMessage(val.Message(), level+1)
		m.Set(fd, val)
	default:
		m.Set(fd, scalarField(fd.Kind()))
	}
}

func scalarField(kind protoreflect.Kind) protoreflect.Value {
	switch kind {
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(true)

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(1 << 30)

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(1 << 30)

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(1 << 30)

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(1 << 30)

	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(3.14159265)

	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(3.14159265)

	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes([]byte("hello world"))

	case protoreflect.StringKind:
		return protoreflect.ValueOfString("hello world")

	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(42)
	}

	panic(fmt.Sprintf("FieldDescriptor.Kind %v is not valid", kind))
}

func setList(list protoreflect.List, fd protoreflect.FieldDescriptor, level int) {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for i := 0; i < 10; i++ {
			val := list.NewElement()
			fillMessage(val.Message(), level+1)
			list.Append(val)
		}
	default:
		for i := 0; i < 100; i++ {
			list.Append(scalarField(fd.Kind()))
		}
	}
}

func setMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor, level int) {
	fields := fd.Message().Fields()
	keyDesc := fields.ByNumber(1)
	valDesc := fields.ByNumber(2)

	pkey := scalarField(keyDesc.Kind())
	switch kind := valDesc.Kind(); kind {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		val := mmap.NewValue()
		fillMessage(val.Message(), level+1)
		mmap.Set(pkey.MapKey(), val)
	default:
		mmap.Set(pkey.MapKey(), scalarField(kind))
	}
}
