// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	. "github.com/onsi/gomega"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"github.com/sykesm/batik/pkg/log/pretty"
	"github.com/sykesm/batik/pkg/timeparse"
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
			expectedOut: `"ts=[0-9\.]* level=info logger=test caller=log/logging_test\.go:83 msg=test\\n"`,
		},
		{
			testName:    "logs with json",
			encoder:     zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `"{\\"level\\":\\"info\\",\\"ts\\":[0-9\.]*,\\"logger\\":\\"test\\",\\"caller\\":\\"log/logging_test\.go:83\\",\\"msg\\":\\"test\\"}\\n"`,
		},
		{
			testName:    "logs with color",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      true,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `"\\x1b\[37m[a-z,A-Z]* [0-9]* [0-9:]*\.000000\\x1b\[0m \|\\x1b\[36mINFO\\x1b\[0m\| \\x1b\[34mtest\\x1b\[0m \\x1b\[0mlog/logging_test\.go:83\\x1b\[0m \\x1b\[97mtest\\x1b\[0m \\n"`,
		},
		{
			testName:    "logs under level",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("warn"),
			message:     "test",
			expectedOut: `^""$`,
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
				w = pretty.NewWriter(w, zap.NewProductionEncoderConfig(), timeparse.ParseUnixTime)
			}
			ws := NewWriteSyncer(w)
			logger := NewLogger(tt.encoder, ws, tt.leveler).Named("test")

			logger.Info(tt.message)

			gt.Expect(fmt.Sprintf("%q", buf.String())).To(MatchRegexp(tt.expectedOut))
		})
	}
}
