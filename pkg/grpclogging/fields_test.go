// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpclogging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-logfmt/logfmt"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/grpclogging/internal/testprotos/echo"
)

func TestProtoMessage(t *testing.T) {
	gt := NewGomegaWithT(t)

	field := ProtoMessage("message", &echo.Message{
		Message:  "I am the Walrus!",
		Sequence: 42,
	})
	gt.Expect(field.Key).To(Equal("message"))
	gt.Expect(field.Type).To(Equal(zapcore.ReflectType))

	jm, ok := field.Interface.(json.Marshaler)
	gt.Expect(ok).To(BeTrue())

	encoded, err := jm.MarshalJSON()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(encoded).To(MatchJSON(`{"message":"I am the Walrus!", "sequence":42}`))
}

func TestProtoMessageAny(t *testing.T) {
	gt := NewGomegaWithT(t)

	field := ProtoMessage("any-message", "some-string")
	gt.Expect(field.Key).To(Equal("any-message"))
	gt.Expect(field.Type).To(Equal(zapcore.StringType))
	gt.Expect(field.String).To(Equal("some-string"))
}

func TestProtoMessageEncoding(t *testing.T) {
	gt := NewGomegaWithT(t)

	encoder := zaplogfmt.NewEncoder(zapcore.EncoderConfig{})
	buf, err := encoder.EncodeEntry(
		zapcore.Entry{},
		[]zapcore.Field{
			ProtoMessage("proto", &echo.Message{Message: "I am the Walrus!", Sequence: 42}),
		},
	)
	gt.Expect(err).NotTo(HaveOccurred())

	dec := logfmt.NewDecoder(bytes.NewBuffer(buf.Bytes()))
	gt.Expect(dec.ScanRecord()).To(BeTrue())
	gt.Expect(dec.ScanKeyval()).To(BeTrue())
	gt.Expect(string(dec.Key())).To(Equal("proto"))
	gt.Expect(dec.Value()).To(MatchJSON(`{"message":"I am the Walrus!", "sequence": 42}`))
}

func TestError(t *testing.T) {
	gt := NewGomegaWithT(t)

	err := errors.New("error message")
	_, ok := err.(fmt.Formatter)
	gt.Expect(ok).To(BeTrue(), "should be an error that implements fmt.Formatter")

	field := Error(err)
	gt.Expect(field.Type).To(Equal(zapcore.ErrorType))
	gt.Expect(field.Integer).NotTo(Equal(err))

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})
	buf, err := encoder.EncodeEntry(zapcore.Entry{}, []zapcore.Field{field})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(buf).To(MatchJSON(`{"error": "error message"}`))
}

func TestErrorWithNil(t *testing.T) {
	gt := NewGomegaWithT(t)
	field := Error(nil)
	gt.Expect(field.Type).To(Equal(zapcore.SkipType))
}
