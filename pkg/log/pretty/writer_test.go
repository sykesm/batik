// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log/pretty/color"
)

func TestWrite(t *testing.T) {
	ts := time.Now()
	encoderConfig := zapcore.EncoderConfig{
		NameKey:    "logger",
		LevelKey:   "level",
		MessageKey: "msg",
		CallerKey:  "caller",
		TimeKey:    "ts",
	}
	timeParser := func(s string) (time.Time, error) {
		return ts, nil
	}

	tests := map[string]struct {
		input      string
		expected   string
		errMatcher types.GomegaMatcher
	}{
		"not logfmt": {
			input:      "this is not a message in logfmt format",
			errMatcher: MatchError("not a logfmt string"),
		},
		"logfmt": {
			input: `ts=1600356328.141956 level=info logger=batik caller=app/caller.go:99 msg="the message"`,
			expected: "\x1b[37m" + ts.Format(time.StampMicro) + "\x1b[0m " +
				"|\x1b[36mINFO\x1b[0m| " +
				"\x1b[34mbatik\x1b[0m " +
				"\x1b[0mapp/caller.go:99\x1b[0m " +
				"\x1b[97mthe message\x1b[0m\n",
		},
		"with fields": {
			input: `ts=1600356328.141956 level=info logger=batik caller=app/caller.go:99 msg="the message" key1=value1 key2="[value two]"`,
			expected: "\x1b[37m" + ts.Format(time.StampMicro) + "\x1b[0m " +
				"|\x1b[36mINFO\x1b[0m| " +
				"\x1b[34mbatik\x1b[0m " +
				"\x1b[0mapp/caller.go:99\x1b[0m " +
				"\x1b[97mthe message\x1b[0m " +
				"\x1b[32mkey1\x1b[0m=\x1b[97mvalue1\x1b[0m " +
				"\x1b[32mkey2\x1b[0m=\x1b[97m[value two]\x1b[0m\n",
		},
		"interposed fields in header": {
			input: `ts=1600356328.141956 key1=value1 level=info logger=batik caller=app/caller.go:99 msg="the message" key2="[value two]"`,
			expected: "\x1b[37m" + ts.Format(time.StampMicro) + "\x1b[0m " +
				"|\x1b[36mINFO\x1b[0m| " +
				"\x1b[34mbatik\x1b[0m " +
				"\x1b[0mapp/caller.go:99\x1b[0m " +
				"\x1b[97mthe message\x1b[0m " +
				"\x1b[32mkey1\x1b[0m=\x1b[97mvalue1\x1b[0m " +
				"\x1b[32mkey2\x1b[0m=\x1b[97m[value two]\x1b[0m\n",
		},
		"missing header fields": {
			input:    `ts=1600356328.141956 level=info`,
			expected: "\x1b[37m" + ts.Format(time.StampMicro) + "\x1b[0m |\x1b[36mINFO\x1b[0m|   \n",
		},
		"missing header keys with fields": {
			input: `ts=1600356328.141956 level=info key1=value1`,
			expected: "\x1b[37m" + ts.Format(time.StampMicro) + "\x1b[0m " +
				"|\x1b[36mINFO\x1b[0m|    " +
				"\x1b[32mkey1\x1b[0m=\x1b[97mvalue1\x1b[0m\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			buf := &bytes.Buffer{}
			w := NewWriter(buf, encoderConfig, timeParser)

			n, err := w.Write([]byte(tt.input))
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(n).To(Equal(len(tt.input)), "write should return 0 <= n <= len(input)")
			gt.Expect(buf.String()).To(Equal(tt.expected))
		})
	}
}

type failingWriter struct {
	failAfter int
}

func (f *failingWriter) Write(p []byte) (int, error) {
	if f.failAfter == 0 {
		return 0, errors.New("boom")
	}
	f.failAfter--
	return len(p), nil
}

func TestWriteErrors(t *testing.T) {
	tests := map[string]struct {
		fw         *failingWriter
		errMatcher types.GomegaMatcher
	}{
		"header":                        {fw: &failingWriter{0}, errMatcher: MatchError("boom")},
		"fields before header complete": {fw: &failingWriter{1}, errMatcher: MatchError("boom")},
		"fields":                        {fw: &failingWriter{2}, errMatcher: MatchError("boom")},
		"newline":                       {fw: &failingWriter{3}, errMatcher: MatchError("boom")},
		"success":                       {fw: &failingWriter{4}, errMatcher: BeNil()},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			w := NewWriter(tt.fw, zap.NewProductionEncoderConfig(), func(string) (time.Time, error) { return time.Now(), nil })

			input := `ts=1600356328.141956 key1=value1 level=info logger=batik caller=app/caller.go:99 msg="the message" key2="[value two]"`
			_, err := w.Write([]byte(input))
			gt.Expect(err).To(tt.errMatcher)
		})
	}
}

func TestLevelColors(t *testing.T) {
	tests := []struct {
		level string
		color color.Color
		text  string
	}{
		{level: "debug", color: color.FgMagenta, text: "DEBU"},
		{level: "DEBUG", color: color.FgMagenta, text: "DEBU"},
		{level: "info", color: color.FgCyan, text: "INFO"},
		{level: "warn", color: color.FgYellow, text: "WARN"},
		{level: "warning", color: color.FgYellow, text: "WARN"},
		{level: "wArnInG", color: color.FgYellow, text: "WARN"},
		{level: "error", color: color.FgRed, text: "ERRO"},
		{level: "fatal", color: color.BgHiRed, text: "FATA"},
		{level: "panic", color: color.BgHiRed, text: "PANI"},
		{level: "default", color: color.FgMagenta, text: "DEFA"},
		{level: "z", color: color.FgMagenta, text: "Z   "},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			buf := &bytes.Buffer{}
			w := NewWriter(
				buf,
				zapcore.EncoderConfig{LevelKey: "level"},
				func(string) (time.Time, error) { return time.Now(), nil },
			)

			w.Write([]byte("time=123.456 level=" + tt.level + "\n"))
			gt.Expect(buf.String()).To(ContainSubstring("|" + tt.color.Sprint(tt.text) + "|"))
		})
	}
}
