// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"io"
	"testing"

	. "github.com/onsi/gomega"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"github.com/sykesm/batik/pkg/log/pretty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		testName    string
		encoder     zapcore.Encoder
		pretty      bool
		leveler     zapcore.LevelEnabler
		message     string
		expectedOut string
	}{
		{
			testName:    "logs with logfmt",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `ts=.* level=info caller=log/logging_test.go:.* msg=test`,
		},
		{
			testName:    "logs with json",
			encoder:     zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `{"level":"info","ts":.*,"caller":"log/logging_test.go:.*","msg":"test"}`,
		},
		{
			testName:    "logs with color",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      true,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `\x1b.*|\x1b.*INFO.*\x1b.*|.*test.*`,
		},
		{
			testName:    "logs under level",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("warn"),
			message:     "test",
			expectedOut: "^$",
		},
		// {
		// 	testName:    "logs to file",
		// 	level:       "info",
		//  config:      Config{Writer: &os.File{}},
		// 	message:     "test",
		// 	expectedOut: `{"level":"info","time":".*","message":"test"}`,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var w io.Writer
			buf := &bytes.Buffer{}
			w = buf
			if tt.pretty {
				w = &pretty.Writer{Writer: w}
			}
			ws := NewWriteSyncer(w)
			logger := NewLogger(tt.encoder, ws, tt.leveler)

			logger.Info(tt.message)

			gt.Expect(buf.String()).To(MatchRegexp(tt.expectedOut))
		})
	}
}
