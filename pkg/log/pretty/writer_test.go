// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
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

	t.Run("NonLogfmt", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		buf := &bytes.Buffer{}
		w := NewWriter(buf, encoderConfig, timeParser)
		_, err := w.Write([]byte("this is not a logfmt line"))
		gt.Expect(err).To(MatchError("not a logfmt string"))
	})

	t.Run("Logfmt", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		buf := &bytes.Buffer{}
		w := NewWriter(buf, encoderConfig, timeParser)

		testLine := `ts=1600356328.141956 level=info logger=batik caller=app/caller.go:99 msg="the message"`
		expectedLine := fmt.Sprintf("\x1b[37m%s\x1b[0m |\x1b[36mINFO\x1b[0m| \x1b[34mbatik\x1b[0m \x1b[0mapp/caller.go:99\x1b[0m \x1b[97mthe message\x1b[0m \n", ts.Format(time.StampMicro))

		n, err := w.Write([]byte(testLine))
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(n).To(Equal(110))
		gt.Expect(buf.String()).To(Equal(expectedLine))
	})
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
