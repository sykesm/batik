// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		testName    string
		level       string
		logPath     interface{}
		message     string
		expectedOut string
	}{
		{
			testName:    "logs at level",
			level:       "info",
			logPath:     &bytes.Buffer{},
			message:     "test",
			expectedOut: `{"level":"info","time":".*","message":"test"}`,
		},
		{
			testName:    "logs under level",
			level:       "warn",
			logPath:     &bytes.Buffer{},
			message:     "test",
			expectedOut: "^$",
		},
		{
			testName:    "logs to file",
			level:       "info",
			logPath:     "out.txt",
			message:     "test",
			expectedOut: `{"level":"info","time":".*","message":"test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			logger, err := NewLogger(tt.level, tt.logPath)
			gt.Expect(err).NotTo(HaveOccurred())

			logger.Info().Msg(tt.message)

			switch tt.logPath.(type) {
			case io.Writer:
				gt.Expect(tt.logPath).To(MatchRegexp(tt.expectedOut))
			case string:
				defer os.Remove(tt.logPath.(string))
				gt.Expect(tt.logPath).To(BeAnExistingFile())
				bytes, err := ioutil.ReadFile(tt.logPath.(string))
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(bytes).To(MatchRegexp(tt.expectedOut))
			}
		})
	}

	t.Run("logs with no level", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		buf := gbytes.NewBuffer()
		logger, err := NewLogger("warn", buf)
		gt.Expect(err).NotTo(HaveOccurred())

		logger.Log().Msg("test")

		gt.Expect(buf).To(gbytes.Say(`{"time":".*","message":"test"}`))
	})
}

func TestNewLogger_Failures(t *testing.T) {
	tests := []struct {
		testName    string
		level       string
		logPath     interface{}
		expectedErr string
	}{
		{
			testName:    "invalid level",
			level:       "invalid",
			logPath:     gbytes.NewBuffer(),
			expectedErr: "Unknown Level String: 'invalid', defaulting to NoLevel",
		},
		{
			testName:    "invalid path",
			level:       "info",
			logPath:     "/",
			expectedErr: "open /: is a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			logger, err := NewLogger(tt.level, tt.logPath)
			gt.Expect(err).To(MatchError(MatchRegexp(tt.expectedErr)))
			gt.Expect(logger).To(BeNil())
		})
	}
}
