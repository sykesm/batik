// Copyright (c) 2019 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// * Redistributions of source code must retain the above copyright
//   notice, this list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
// * Neither the name of Google Inc. nor the names of its contributors
//   may be used to endorse or promote products derived from this software
//   without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Copyright IBM Corp. All Rights Reserved.
//
// The official implementation all goes out of its way to highlight how the
// "deterministic" marshaling should not be relied on for stability so we'll
// duplicate some of their reflect-based marshaling code to encode messages
// that require stable serialization.
//
// This immplementation is based on the reflection based code from protobuf and
// does not benefit from optimizations from caching the operations required to
// marshal a message.
//
// By default, the deterministic marshaling of a message message is encoded by
// sorting the fields by the proto field number but there are some exceptions.
// One-of fields, in particular, are pushed to the end of the message and
// extensions are brought to the front.
//
// See protobuf/internal/impl/codec_message#makeCodeMethods for details on
// how the coder fields are created and ordered and
// protobuf/internal/impl/encode#marshalAppendPointer for details on how the
// coder fields are referenced after handling extensions.

package protomsg

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MarshalDeterministic performs deterministic marshaling of a protocol buffer
// message. The official implementation emphasizes that we cannot rely on the
// stability of deterministic marshaling across versions and recommends that
// those that require stable, deterministic marshaling implement their own
// encoding.
//
// This method is a slighlty modified clone of the the "slow" encoding
// implemented in the official protobuf package. By copying the implementation
// we can isolate ourselves from any upstream changes and implement any
// additional runtime checks that are necessary. This comes at a cost.
//
// The "slow" encoding relies on reflection and is about an order of magnitude
// slower than than the optimized implementation. This means that we should
// only use this implementation when a persistent, stable encoding is requried.
//
// Additional optimzations can be implemented as necessary. In particualr, we
// can implement custom marshaling for our wire types that does not rely on
// reflection or we can can punt and use the official deterministic marshaling
// implementation for as long as it remains stable. This requires some
// additional test infrascture to validate the constraints hold.
func MarshalDeterministic(m proto.Message) ([]byte, error) {
	if m == nil || !m.ProtoReflect().IsValid() {
		return proto.Marshal(m)
	}

	return marshal(nil, m.ProtoReflect())
}

// marshal is a centralized function that all marshal operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for marshal that do not go through this.
func marshal(b []byte, m protoreflect.Message) ([]byte, error) {
	if m.Descriptor().Syntax() == protoreflect.Proto2 {
		// Explicitly refuse to encode proto2 messages.
		// Support for required fields and groups has been removed.
		return nil, errors.New("protomsg: proto2 syntax is not supported")
	}

	var err error
	b, err = marshalMessageSlow(b, m)
	if err != nil {
		return nil, err
	}
	return b, checkInitialized(m)
}

func marshalMessageSlow(b []byte, m protoreflect.Message) ([]byte, error) {
	if isMessageSet(m.Descriptor()) {
		return b, errors.New("protomsg: no support for message_set_wire_format")
	}

	var err error
	rangeFields(m, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		b, err = marshalField(b, fd, v)
		return err == nil
	})
	if err != nil {
		return b, err
	}
	if len(m.GetUnknown()) != 0 {
		return b, fmt.Errorf("protomsg: refusing to marshal unknown fields with length %d", len(m.GetUnknown()))
	}
	return b, nil
}

// isMessageSet returns whether the message uses the unsupported MessageSet
// wire format.
func isMessageSet(md protoreflect.MessageDescriptor) bool {
	xmd, ok := md.(interface{ IsMessageSet() bool })
	return ok && xmd.IsMessageSet()
}

// rangeFields visits fields in a defined order.
func rangeFields(m protoreflect.Message, f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	var fds []protoreflect.FieldDescriptor
	m.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		fds = append(fds, fd)
		return true
	})
	sort.Slice(fds, func(a, b int) bool {
		return fieldDescLess(fds[a], fds[b])
	})
	for _, fd := range fds {
		if !f(fd, m.Get(fd)) {
			break
		}
	}
}

func marshalField(b []byte, fd protoreflect.FieldDescriptor, value protoreflect.Value) ([]byte, error) {
	switch {
	case fd.IsList():
		return marshalList(b, fd, value.List())
	case fd.IsMap():
		return marshalMap(b, fd, value.Map())
	default:
		b = protowire.AppendTag(b, fd.Number(), wireTypes[fd.Kind()])
		return marshalSingular(b, fd, value)
	}
}

func marshalList(b []byte, fd protoreflect.FieldDescriptor, list protoreflect.List) ([]byte, error) {
	if fd.IsPacked() && list.Len() > 0 {
		b = protowire.AppendTag(b, fd.Number(), protowire.BytesType)
		b, pos := appendSpeculativeLength(b)
		for i, llen := 0, list.Len(); i < llen; i++ {
			var err error
			b, err = marshalSingular(b, fd, list.Get(i))
			if err != nil {
				return b, err
			}
		}
		b = finishSpeculativeLength(b, pos)
		return b, nil
	}

	kind := fd.Kind()
	for i, llen := 0, list.Len(); i < llen; i++ {
		var err error
		b = protowire.AppendTag(b, fd.Number(), wireTypes[kind])
		b, err = marshalSingular(b, fd, list.Get(i))
		if err != nil {
			return b, err
		}
	}
	return b, nil
}

func marshalMap(b []byte, fd protoreflect.FieldDescriptor, mapv protoreflect.Map) ([]byte, error) {
	keyf := fd.MapKey()
	valf := fd.MapValue()
	var err error
	rangeMapSorted(mapv, keyf.Kind(), func(key protoreflect.MapKey, value protoreflect.Value) bool {
		b = protowire.AppendTag(b, fd.Number(), protowire.BytesType)
		var pos int
		b, pos = appendSpeculativeLength(b)

		b, err = marshalField(b, keyf, key.Value())
		if err != nil {
			return false
		}
		b, err = marshalField(b, valf, value)
		if err != nil {
			return false
		}
		b = finishSpeculativeLength(b, pos)
		return true
	})
	return b, err
}

func checkInitialized(m protoreflect.Message) error {
	var err error
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsList():
			if fd.Message() == nil {
				return true
			}
			for i, list := 0, v.List(); i < list.Len() && err == nil; i++ {
				err = checkInitialized(list.Get(i).Message())
			}
		case fd.IsMap():
			if fd.MapValue().Message() == nil {
				return true
			}
			v.Map().Range(func(key protoreflect.MapKey, v protoreflect.Value) bool {
				err = checkInitialized(v.Message())
				return err == nil
			})
		default:
			if fd.Message() == nil {
				return true
			}
			err = checkInitialized(v.Message())
		}
		return err == nil
	})
	return err
}

// fieldDescLess returns true if field a comes before field j in ordered wire marshal
// output. This is a copy of google.golang.org/protobuf/internal/fieldsort#Less.
func fieldDescLess(a, b protoreflect.FieldDescriptor) bool {
	ea := a.IsExtension()
	eb := b.IsExtension()
	oa := a.ContainingOneof()
	ob := b.ContainingOneof()
	switch {
	case ea != eb:
		return ea
	case oa != nil && ob != nil:
		if oa == ob {
			return a.Number() < b.Number()
		}
		return oa.Index() < ob.Index()
	case oa != nil && !oa.IsSynthetic():
		return false
	case ob != nil && !ob.IsSynthetic():
		return true
	default:
		return a.Number() < b.Number()
	}
}

// rangeMapSorted iterates over every map entry in sorted key order,
// calling f for each key and value encountered.
func rangeMapSorted(mapv protoreflect.Map, keyKind protoreflect.Kind, f func(protoreflect.MapKey, protoreflect.Value) bool) {
	var keys []protoreflect.MapKey
	mapv.Range(func(key protoreflect.MapKey, _ protoreflect.Value) bool {
		keys = append(keys, key)
		return true
	})
	sort.Slice(keys, func(i, j int) bool {
		switch keyKind {
		case protoreflect.BoolKind:
			return !keys[i].Bool() && keys[j].Bool()
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return keys[i].Int() < keys[j].Int()
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return keys[i].Uint() < keys[j].Uint()
		case protoreflect.StringKind:
			return keys[i].String() < keys[j].String()
		default:
			panic("invalid kind: " + keyKind.String())
		}
	})
	for _, key := range keys {
		if !f(key, mapv.Get(key)) {
			break
		}
	}
}

var wireTypes = map[protoreflect.Kind]protowire.Type{
	protoreflect.BoolKind:     protowire.VarintType,
	protoreflect.EnumKind:     protowire.VarintType,
	protoreflect.Int32Kind:    protowire.VarintType,
	protoreflect.Sint32Kind:   protowire.VarintType,
	protoreflect.Uint32Kind:   protowire.VarintType,
	protoreflect.Int64Kind:    protowire.VarintType,
	protoreflect.Sint64Kind:   protowire.VarintType,
	protoreflect.Uint64Kind:   protowire.VarintType,
	protoreflect.Sfixed32Kind: protowire.Fixed32Type,
	protoreflect.Fixed32Kind:  protowire.Fixed32Type,
	protoreflect.FloatKind:    protowire.Fixed32Type,
	protoreflect.Sfixed64Kind: protowire.Fixed64Type,
	protoreflect.Fixed64Kind:  protowire.Fixed64Type,
	protoreflect.DoubleKind:   protowire.Fixed64Type,
	protoreflect.StringKind:   protowire.BytesType,
	protoreflect.BytesKind:    protowire.BytesType,
	protoreflect.MessageKind:  protowire.BytesType,
}

func marshalSingular(b []byte, fd protoreflect.FieldDescriptor, v protoreflect.Value) ([]byte, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		b = protowire.AppendVarint(b, protowire.EncodeBool(v.Bool()))
	case protoreflect.EnumKind:
		b = protowire.AppendVarint(b, uint64(v.Enum()))
	case protoreflect.Int32Kind:
		b = protowire.AppendVarint(b, uint64(int32(v.Int())))
	case protoreflect.Sint32Kind:
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(int32(v.Int()))))
	case protoreflect.Uint32Kind:
		b = protowire.AppendVarint(b, uint64(uint32(v.Uint())))
	case protoreflect.Int64Kind:
		b = protowire.AppendVarint(b, uint64(v.Int()))
	case protoreflect.Sint64Kind:
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(v.Int()))
	case protoreflect.Uint64Kind:
		b = protowire.AppendVarint(b, v.Uint())
	case protoreflect.Sfixed32Kind:
		b = protowire.AppendFixed32(b, uint32(v.Int()))
	case protoreflect.Fixed32Kind:
		b = protowire.AppendFixed32(b, uint32(v.Uint()))
	case protoreflect.FloatKind:
		b = protowire.AppendFixed32(b, math.Float32bits(float32(v.Float())))
	case protoreflect.Sfixed64Kind:
		b = protowire.AppendFixed64(b, uint64(v.Int()))
	case protoreflect.Fixed64Kind:
		b = protowire.AppendFixed64(b, v.Uint())
	case protoreflect.DoubleKind:
		b = protowire.AppendFixed64(b, math.Float64bits(v.Float()))
	case protoreflect.StringKind:
		if !utf8.ValidString(v.String()) {
			return b, fmt.Errorf("protomsg: field %v contains invalid UTF-8", fd.FullName())
		}
		b = protowire.AppendString(b, v.String())
	case protoreflect.BytesKind:
		b = protowire.AppendBytes(b, v.Bytes())
	case protoreflect.MessageKind:
		var pos int
		var err error
		b, pos = appendSpeculativeLength(b)
		b, err = marshal(b, v.Message())
		if err != nil {
			return b, err
		}
		b = finishSpeculativeLength(b, pos)
	default:
		return b, fmt.Errorf("protomsg: invalid kind %v", fd.Kind())
	}
	return b, nil
}

// When encoding length-prefixed fields, we speculatively set aside some number of bytes
// for the length, encode the data, and then encode the length (shifting the data if necessary
// to make room).
const speculativeLength = 1

func appendSpeculativeLength(b []byte) ([]byte, int) {
	pos := len(b)
	b = append(b, "\x00\x00\x00\x00"[:speculativeLength]...)
	return b, pos
}

func finishSpeculativeLength(b []byte, pos int) []byte {
	mlen := len(b) - pos - speculativeLength
	msiz := protowire.SizeVarint(uint64(mlen))
	if msiz != speculativeLength {
		for i := 0; i < msiz-speculativeLength; i++ {
			b = append(b, 0)
		}
		copy(b[pos+msiz:], b[pos+speculativeLength:])
		b = b[:pos+msiz+mlen]
	}
	protowire.AppendVarint(b[:pos], uint64(mlen))
	return b
}
