// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ValidateMessage will verify that a message does not contain unknown fields,
// proto2 messages, or proto2 groups. These checks can be used to complement
// the official deterministic marshaling impelmentation to ensure we don't end
// up with any surprises in the encoded messages.
func ValidateMessage(m proto.Message) error {
	if m != nil {
		return validateMessage(m.ProtoReflect())
	}
	return nil
}

func validateMessage(m protoreflect.Message) error {
	// Explicitly disallow marshaling of messages with unknown fields.
	if len(m.GetUnknown()) != 0 {
		return fmt.Errorf("protomsg: refusing to marshal unknown fields with length %d", len(m.GetUnknown()))
	}

	// Explicitly refuse to encode proto2 messages.
	if m.Descriptor().Syntax() == protoreflect.Proto2 {
		return errors.New("protomsg: proto2 syntax is not supported")
	}

	var err error
	m.Range(func(fd protoreflect.FieldDescriptor, pv protoreflect.Value) bool {
		switch fd.Kind() {
		case protoreflect.MessageKind:
			switch {
			case fd.IsList():
				err = validateList(fd, pv.List())
			case fd.IsMap():
				err = validateMap(fd, pv.Map())
			default:
				err = validateMessage(pv.Message())
			}
			return err == nil

		case protoreflect.GroupKind:
			err = errors.New("protomsg: proto2 groups are not supported")
			return false

		default:
			return true
		}
	})
	return err
}

func validateList(fd protoreflect.FieldDescriptor, list protoreflect.List) error {
	for i, llen := 0, list.Len(); i < llen; i++ {
		if err := validateMessage(list.Get(i).Message()); err != nil {
			return err
		}
	}
	return nil
}

func validateMap(fd protoreflect.FieldDescriptor, mapv protoreflect.Map) error {
	if fd.MapValue().Kind() != protoreflect.MessageKind {
		return nil
	}

	var err error
	mapv.Range(func(_ protoreflect.MapKey, v protoreflect.Value) bool {
		err = validateMessage(v.Message())
		return err == nil
	})
	return err
}
