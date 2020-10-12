// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log/pretty"
	"github.com/sykesm/batik/pkg/tested"
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
			expectedOut: `"ts=[0-9\.]* level=info logger=test caller=log/logging_test\.go:[0-9]* msg=test\\n"`,
		},
		{
			testName:    "logs with json",
			encoder:     zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `"{\\"level\\":\\"info\\",\\"ts\\":[0-9\.]*,\\"logger\\":\\"test\\",\\"caller\\":\\"log/logging_test\.go:[0-9]*\\",\\"msg\\":\\"test\\"}\\n"`,
		},
		{
			testName:    "logs with color",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      true,
			leveler:     NewLeveler("info"),
			message:     "test",
			expectedOut: `"\\x1b\[37m[[:alpha:]]{3}\s+[0-9]{1,2}\s+[0-9:]{8}\.\d{6}\\x1b\[0m \|\\x1b\[36mINFO\\x1b\[0m\| \\x1b\[34mtest\\x1b\[0m \\x1b\[0mlog/logging_test\.go:[0-9]*\\x1b\[0m \\x1b\[97mtest\\x1b\[0m\\n"`,
		},
		{
			testName:    "logs under level",
			encoder:     zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			pretty:      false,
			leveler:     NewLeveler("warn"),
			message:     "test",
			expectedOut: `^""$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var w io.Writer
			buf := &bytes.Buffer{}
			w = buf
			if tt.pretty {
				w = pretty.NewWriter(w, zap.NewProductionEncoderConfig(), pretty.ParseUnixTime)
			}
			ws := NewWriteSyncer(w)
			logger := NewLogger(tt.encoder, ws, tt.leveler).Named("test")

			logger.Info(tt.message)

			gt.Expect(fmt.Sprintf("%q", buf.String())).To(MatchRegexp(tt.expectedOut))
		})
	}
}

func TestNewWriteSyncer(t *testing.T) {
	gt := NewGomegaWithT(t)

	pr, pw, err := os.Pipe()
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, pw)
	go io.Copy(ioutil.Discard, pr)

	tests := map[string]struct {
		input io.Writer
	}{
		"nil":          {input: nil},
		"file":         {input: pw},
		"write syncer": {input: zapcore.AddSync(&bytes.Buffer{})},
		"naked":        {input: &bytes.Buffer{}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(true).To(BeTrue())

			result := NewWriteSyncer(tt.input)
			gt.Expect(result).NotTo(BeNil())
			gt.Expect(result.Sync()).To(Succeed())
		})
	}
}
