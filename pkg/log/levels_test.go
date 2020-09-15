// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

func TestNameToLevel(t *testing.T) {
	var tests = []struct {
		names       []string
		level       zapcore.Level
		expectedErr string
	}{
		{names: []string{"PAYLOAD", "payload"}, level: payloadLevel},
		{names: []string{"DEBUG", "debug"}, level: zapcore.DebugLevel},
		{names: []string{"INFO", "info"}, level: zapcore.InfoLevel},
		{names: []string{"WARNING", "warning", "WARN", "warn"}, level: zapcore.WarnLevel},
		{names: []string{"ERROR", "error"}, level: zapcore.ErrorLevel},
		{names: []string{"DPANIC", "dpanic"}, level: zapcore.DPanicLevel},
		{names: []string{"PANIC", "panic"}, level: zapcore.PanicLevel},
		{names: []string{"FATAL", "fatal"}, level: zapcore.FatalLevel},
		{names: []string{"invalid"}, level: zapcore.InfoLevel},
		{names: []string{""}, level: zapcore.InfoLevel},
	}

	for _, tc := range tests {
		for _, name := range tc.names {
			t.Run(name, func(t *testing.T) {
				gt := NewGomegaWithT(t)
				level := NameToLevel(name)
				gt.Expect(level).To(Equal(tc.level))
			})
		}
	}
}
