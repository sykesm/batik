// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package encoder

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func TestEncodeEntry(t *testing.T) {
	startTime := time.Now()
	var tests = []struct {
		name     string
		spec     string
		fields   []zapcore.Field
		expected string
	}{
		{name: "empty spec and nil fields", spec: "", fields: nil, expected: "\n"},
		{name: "empty spec with fields", spec: "", fields: []zapcore.Field{zap.String("key", "value")}, expected: "key=value\n"},
		{name: "simple spec and nil fields", spec: "simple-string", expected: "simple-string\n"},
		{name: "simple spec and empty fields", spec: "simple-string", fields: []zapcore.Field{}, expected: "simple-string\n"},
		{name: "simple spec with fields", spec: "simple-string", fields: []zapcore.Field{zap.String("key", "value")}, expected: "simple-string key=value\n"},
		{name: "duration", spec: "", fields: []zapcore.Field{zap.Duration("duration", time.Second)}, expected: "duration=1s\n"},
		{name: "time", spec: "", fields: []zapcore.Field{zap.Time("time", startTime)}, expected: fmt.Sprintf("time=%s\n", startTime.Format("2006-01-02T15:04:05.999Z07:00"))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			formatters, err := ParseFormat(tc.spec)
			gt.Expect(err).NotTo(HaveOccurred())

			enc := NewFormatEncoder(formatters...)

			pc, file, l, ok := runtime.Caller(0)
			line, err := enc.EncodeEntry(
				zapcore.Entry{
					// The entry information should be completely omitted
					Level:      zapcore.InfoLevel,
					Time:       startTime,
					LoggerName: "logger-name",
					Message:    "message",
					Caller:     zapcore.NewEntryCaller(pc, file, l, ok),
					Stack:      "stack",
				},
				tc.fields,
			)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(line.String()).To(Equal(tc.expected))
		})
	}
}

type brokenEncoder struct{ zapcore.Encoder }

func (b *brokenEncoder) EncodeEntry(zapcore.Entry, []zapcore.Field) (*buffer.Buffer, error) {
	return nil, errors.New("broken encoder")
}

func TestEncodeFieldsFailed(t *testing.T) {
	gt := NewGomegaWithT(t)

	enc := NewFormatEncoder()
	enc.Encoder = &brokenEncoder{}

	_, err := enc.EncodeEntry(zapcore.Entry{}, nil)
	gt.Expect(err).To(MatchError("broken encoder"))
}

func TestFormatEncoderClone(t *testing.T) {
	gt := NewGomegaWithT(t)

	enc := NewFormatEncoder()
	cloned := enc.Clone()
	gt.Expect(cloned).To(Equal(enc))
}
