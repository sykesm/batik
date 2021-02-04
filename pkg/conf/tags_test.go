// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	. "github.com/onsi/gomega"
)

type Simple struct {
	IngoredInt  int
	Relpath     string `batik:"relpath"`
	IgnoredBool bool
}
type (
	Field       struct{ Simple Simple }
	FieldPtr    struct{ Simple *Simple }
	Embedded    struct{ Simple }
	EmbeddedPtr struct{ *Simple }
	WrongType   struct {
		Relpath int `batik:"relpath"`
	}
)

type Sliced struct {
	Simples []Simple
}

func TestResolveRelpathAddressable(t *testing.T) {
	expected := "source_path/testdata"
	tests := map[string]struct {
		input  interface{}
		result func(interface{}) string
	}{
		"simple ref": {
			input:  &Simple{Relpath: "testdata"},
			result: func(s interface{}) string { return s.(*Simple).Relpath },
		},
		"field ref": {
			input:  &Field{Simple: Simple{Relpath: "testdata"}},
			result: func(s interface{}) string { return s.(*Field).Simple.Relpath },
		},
		"field ptr": {
			input:  FieldPtr{Simple: &Simple{Relpath: "testdata"}},
			result: func(s interface{}) string { return s.(FieldPtr).Simple.Relpath },
		},
		"field ptr ref": {
			input:  &FieldPtr{Simple: &Simple{Relpath: "testdata"}},
			result: func(s interface{}) string { return s.(*FieldPtr).Simple.Relpath },
		},
		"embedded ref": {
			input:  &Embedded{Simple: Simple{Relpath: "testdata"}},
			result: func(s interface{}) string { return s.(*Embedded).Relpath },
		},
		"embedded ptr": {
			input:  EmbeddedPtr{Simple: &Simple{Relpath: "testdata"}},
			result: func(s interface{}) string { return s.(EmbeddedPtr).Relpath },
		},
		"embedded ptr ref": {
			input:  &EmbeddedPtr{Simple: &Simple{Relpath: "testdata"}},
			result: func(s interface{}) string { return s.(*EmbeddedPtr).Relpath },
		},
		"slice field": {
			input:  Sliced{Simples: []Simple{{Relpath: "testdata"}}},
			result: func(s interface{}) string { return s.(Sliced).Simples[0].Relpath },
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tr := &TagResolver{SourcePath: "source_path"}
			err := tr.Resolve(tt.input)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(tt.result(tt.input)).To(Equal(expected))
		})
	}
}

func TestResolveRelpathMultipleSliceEntries(t *testing.T) {
	gt := NewGomegaWithT(t)
	tr := &TagResolver{SourcePath: "source_path"}

	multi := struct {
		Entries []struct {
			Path string `batik:"relpath"`
		}
	}{
		Entries: []struct {
			Path string `batik:"relpath"`
		}{
			{Path: "path/A"},
			{Path: "path/B"},
		},
	}

	err := tr.Resolve(&multi)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(multi.Entries[0].Path).To(Equal("source_path/path/A"))
	gt.Expect(multi.Entries[1].Path).To(Equal("source_path/path/B"))
}

func TestResolveRelpathMultipleFields(t *testing.T) {
	gt := NewGomegaWithT(t)
	tr := &TagResolver{SourcePath: "source_path"}

	multi := struct {
		PathA string `batik:"relpath"`
		PathB string `batik:"relpath"`
	}{
		PathA: "path/A",
		PathB: "path/B",
	}
	err := tr.Resolve(&multi)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(multi.PathA).To(Equal("source_path/path/A"))
	gt.Expect(multi.PathB).To(Equal("source_path/path/B"))
}

func TestResolveRelpathNotAddressable(t *testing.T) {
	tests := map[string]struct {
		input interface{}
	}{
		"simple":   {input: Simple{Relpath: "testdata"}},
		"field":    {input: Field{Simple: Simple{Relpath: "testdata"}}},
		"embedded": {input: Embedded{Simple: Simple{Relpath: "testdata"}}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tr := &TagResolver{}
			err := tr.Resolve(tt.input)
			gt.Expect(err).To(MatchError("field must be addressable"))
		})
	}
}

func TestResolveUnknownDirective(t *testing.T) {
	gt := NewGomegaWithT(t)

	unknownDirective := struct {
		Unknown string `batik:"say-what?"`
	}{}

	tr := &TagResolver{}
	err := tr.Resolve(&unknownDirective)
	gt.Expect(err).To(MatchError("unknown directive: say-what?"))
}

func TestResolveRelpathNil(t *testing.T) {
	gt := NewGomegaWithT(t)
	tr := &TagResolver{}
	err := tr.Resolve(nil)
	gt.Expect(err).NotTo(HaveOccurred())
}

func TestResolveRelpathTypedNil(t *testing.T) {
	gt := NewGomegaWithT(t)
	tr := &TagResolver{}
	var w io.Writer = (*bytes.Buffer)(nil)
	err := tr.Resolve(w)
	gt.Expect(err).NotTo(HaveOccurred())
}

func TestResolveRelpathWrongType(t *testing.T) {
	gt := NewGomegaWithT(t)
	tr := &TagResolver{}
	err := tr.Resolve(&WrongType{})
	gt.Expect(err).To(MatchError("field must be a string"))
}

func TestBasicTypesNoPanic(t *testing.T) {
	type BasicTypes struct {
		Vbool   bool
		Vdata   interface{}
		Vfloat  float64
		Vint    int
		Vint16  int16
		Vint32  int32
		Vint64  int64
		Vint8   int8
		Vstring string
		Vuint   uint
	}

	type BasicPointers struct {
		Vbool   *bool
		Vdata   *interface{}
		Vfloat  *float64
		Vint    *int
		Vint16  *int16
		Vint32  *int32
		Vint64  *int64
		Vint8   *int8
		Vstring *string
		Vuint   *uint
	}

	tests := map[string]struct {
		input interface{}
	}{
		"basic":        {input: BasicTypes{}},
		"basic ref":    {input: &BasicTypes{}},
		"pointers":     {input: BasicPointers{}},
		"pointers ref": {input: &BasicPointers{}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := (&TagResolver{}).Resolve(tt.input)
			NewGomegaWithT(t).Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestNonStructArgs(t *testing.T) {
	tests := []interface{}{
		nil,
		int(1),
		int8(8),
		int16(16),
		int32(32),
		int64(64),
		uint(1),
		uint8(8),
		uint16(8),
		uint32(8),
		uint64(8),
		io.Reader((*bytes.Buffer)(nil)),
		io.Writer(nil),
		"string",
		float32(32.323232),
		float64(64.646464),
		map[string]string{"key": "value"},
		map[string]interface{}{"three": 3},
		make(chan struct{}),
	}
	for _, arg := range tests {
		t.Run(fmt.Sprintf("%T", arg), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tr := &TagResolver{}
			gt.Expect(tr.Resolve(arg)).To(Succeed())
			gt.Expect(tr.Resolve(&arg)).To(Succeed())
		})
	}
}
