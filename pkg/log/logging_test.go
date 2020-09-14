// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		testName    string
		config      Config
		message     string
		expectedOut string
	}{
		{
			testName:    "logs with logfmt",
			config:      Config{LogSpec: "info", Writer: &bytes.Buffer{}, Format: "logfmt"},
			message:     "test",
			expectedOut: `ts=.* level=info caller=log/logging_test.go:.* msg=test`,
		},
		{
			testName:    "logs with json",
			config:      Config{LogSpec: "info", Writer: &bytes.Buffer{}, Format: "json"},
			message:     "test",
			expectedOut: `{"level":"info","ts":.*,"caller":"log/logging_test.go:.*","msg":"test"}`,
		},
		{
			testName:    "logs under level",
			config:      Config{LogSpec: "warn", Writer: &bytes.Buffer{}},
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

			logger, err := NewLogger(tt.config)
			gt.Expect(err).NotTo(HaveOccurred())

			logger.Info(tt.message)

			// switch tt.config.Writer.(type) {
			// case io.Writer:
			gt.Expect(tt.config.Writer.(*bytes.Buffer).String()).To(MatchRegexp(tt.expectedOut))
			// case string:
			// 	defer os.Remove(tt.logPath.(string))
			// 	gt.Expect(tt.logPath).To(BeAnExistingFile())
			// 	bytes, err := ioutil.ReadFile(tt.logPath.(string))
			// 	gt.Expect(err).NotTo(HaveOccurred())
			// 	gt.Expect(bytes).To(MatchRegexp(tt.expectedOut))
			// }
		})
	}
}

func TestNewLogger_Failures(t *testing.T) {
	tests := []struct {
		testName    string
		config      Config
		expectedErr string
	}{
		{
			testName:    "invalid level",
			config:      Config{LogSpec: "invalid", Writer: &bytes.Buffer{}},
			expectedErr: "invalid log level: invalid",
		},
		// {
		// 	testName:    "invalid path",
		// 	config:      Config{LogSpec: "info", Writer: &os.File{}},
		// 	expectedErr: "open /: is a directory",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			logger, err := NewLogger(tt.config)
			gt.Expect(err).To(MatchError(tt.expectedErr))
			gt.Expect(logger).To(BeNil())
		})
	}
}
