// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package encoder

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestParseFormat(t *testing.T) {
	var tests = []struct {
		desc       string
		spec       string
		formatters []Formatter
	}{
		{
			desc:       "empty spec",
			spec:       "",
			formatters: []Formatter{},
		},
		{
			desc: "simple verb",
			spec: "%{color}",
			formatters: []Formatter{
				ColorFormatter{},
			},
		},
		{
			desc: "with prefix",
			spec: "prefix %{color}",
			formatters: []Formatter{
				StringFormatter{Value: "prefix "},
				ColorFormatter{},
			},
		},
		{
			desc: "with suffix",
			spec: "%{color} suffix",
			formatters: []Formatter{
				ColorFormatter{},
				StringFormatter{Value: " suffix"}},
		},
		{
			desc: "with prefix and suffix",
			spec: "prefix %{color} suffix",
			formatters: []Formatter{
				StringFormatter{Value: "prefix "},
				ColorFormatter{},
				StringFormatter{Value: " suffix"},
			},
		},
		{
			desc: "with format",
			spec: "%{level:.4s} suffix",
			formatters: []Formatter{
				LevelFormatter{FormatVerb: "%.4s"},
				StringFormatter{Value: " suffix"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(tc.desc), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			formatters, err := ParseFormat(tc.spec)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(formatters).To(Equal(tc.formatters))
		})
	}
}

func TestParseFormatError(t *testing.T) {
	gt := NewGomegaWithT(t)

	_, err := ParseFormat("%{color:bad}")
	gt.Expect(err).To(MatchError("invalid color option: bad"))
}

func TestNewFormatter(t *testing.T) {
	var tests = []struct {
		verb      string
		format    string
		formatter Formatter
		errorMsg  string
	}{
		{verb: "color", format: "", formatter: ColorFormatter{}},
		{verb: "color", format: "bold", formatter: ColorFormatter{Bold: true}},
		{verb: "color", format: "reset", formatter: ColorFormatter{Reset: true}},
		{verb: "color", format: "unknown", errorMsg: "invalid color option: unknown"},
		{verb: "id", format: "", formatter: SequenceFormatter{FormatVerb: "%d"}},
		{verb: "id", format: "04x", formatter: SequenceFormatter{FormatVerb: "%04x"}},
		{verb: "level", format: "", formatter: LevelFormatter{FormatVerb: "%s"}},
		{verb: "level", format: ".4s", formatter: LevelFormatter{FormatVerb: "%.4s"}},
		{verb: "message", format: "", formatter: MessageFormatter{FormatVerb: "%s"}},
		{verb: "message", format: "#30s", formatter: MessageFormatter{FormatVerb: "%#30s"}},
		{verb: "module", format: "", formatter: ModuleFormatter{FormatVerb: "%s"}},
		{verb: "module", format: "ok", formatter: ModuleFormatter{FormatVerb: "%ok"}},
		{verb: "shortfunc", format: "", formatter: ShortFuncFormatter{FormatVerb: "%s"}},
		{verb: "shortfunc", format: "U", formatter: ShortFuncFormatter{FormatVerb: "%U"}},
		{verb: "time", format: "", formatter: TimeFormatter{Layout: "2006-01-02T15:04:05.999Z07:00"}},
		{verb: "time", format: "04:05.999999Z05:00", formatter: TimeFormatter{Layout: "04:05.999999Z05:00"}},
		{verb: "unknown", format: "", errorMsg: "unknown verb: unknown"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			f, err := NewFormatter(tc.verb, tc.format)
			if tc.errorMsg == "" {
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(f).To(Equal(tc.formatter))
			} else {
				gt.Expect(err).To(MatchError(tc.errorMsg))
			}
		})
	}
}

func TestColorFormatter(t *testing.T) {
	var tests = []struct {
		f         ColorFormatter
		level     zapcore.Level
		formatted string
	}{
		{f: ColorFormatter{Reset: true}, level: zapcore.DebugLevel, formatted: ResetColor()},
		{f: ColorFormatter{}, level: zapcore.DebugLevel, formatted: ColorCyan.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.DebugLevel, formatted: ColorCyan.Bold()},
		{f: ColorFormatter{}, level: zapcore.InfoLevel, formatted: ColorBlue.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.InfoLevel, formatted: ColorBlue.Bold()},
		{f: ColorFormatter{}, level: zapcore.WarnLevel, formatted: ColorYellow.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.WarnLevel, formatted: ColorYellow.Bold()},
		{f: ColorFormatter{}, level: zapcore.ErrorLevel, formatted: ColorRed.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.ErrorLevel, formatted: ColorRed.Bold()},
		{f: ColorFormatter{}, level: zapcore.DPanicLevel, formatted: ColorMagenta.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.DPanicLevel, formatted: ColorMagenta.Bold()},
		{f: ColorFormatter{}, level: zapcore.PanicLevel, formatted: ColorMagenta.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.PanicLevel, formatted: ColorMagenta.Bold()},
		{f: ColorFormatter{}, level: zapcore.FatalLevel, formatted: ColorMagenta.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.FatalLevel, formatted: ColorMagenta.Bold()},
		{f: ColorFormatter{}, level: zapcore.Level(99), formatted: ColorNone.Normal()},
		{f: ColorFormatter{Bold: true}, level: zapcore.Level(99), formatted: ColorNone.Normal()},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: tc.level}
			tc.f.Format(buf, entry, nil)
			gt.Expect(buf.String()).To(Equal(tc.formatted))
		})
	}
}

func TestLevelFormatter(t *testing.T) {
	var tests = []struct {
		level     zapcore.Level
		formatted string
	}{
		{level: zapcore.DebugLevel, formatted: "DEBUG"},
		{level: zapcore.InfoLevel, formatted: "INFO"},
		{level: zapcore.WarnLevel, formatted: "WARN"},
		{level: zapcore.ErrorLevel, formatted: "ERROR"},
		{level: zapcore.DPanicLevel, formatted: "DPANIC"},
		{level: zapcore.PanicLevel, formatted: "PANIC"},
		{level: zapcore.FatalLevel, formatted: "FATAL"},
		{level: zapcore.Level(99), formatted: "LEVEL(99)"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: tc.level}
			LevelFormatter{FormatVerb: "%s"}.Format(buf, entry, nil)
			gt.Expect(buf.String()).To(Equal(tc.formatted))
		})
	}
}

func TestMessageFormatter(t *testing.T) {
	gt := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	entry := zapcore.Entry{Message: "some message text \n\n"}
	f := MessageFormatter{FormatVerb: "%s"}
	f.Format(buf, entry, nil)
	gt.Expect(buf.String()).To(Equal("some message text "))
}

func TestModuleFormatter(t *testing.T) {
	gt := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	entry := zapcore.Entry{LoggerName: "logger/name"}
	f := ModuleFormatter{FormatVerb: "%s"}
	f.Format(buf, entry, nil)
	gt.Expect(buf.String()).To(Equal("logger/name"))
}

func TestSequenceFormatter(t *testing.T) {
	gt := NewGomegaWithT(t)

	mutex := &sync.Mutex{}
	results := map[string]struct{}{}

	ready := &sync.WaitGroup{}
	ready.Add(100)

	finished := &sync.WaitGroup{}
	finished.Add(100)

	SetSequence(0)
	for i := 1; i <= 100; i++ {
		go func(i int) {
			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: zapcore.DebugLevel}
			f := SequenceFormatter{FormatVerb: "%d"}
			ready.Done() // setup complete
			ready.Wait() // wait for all go routines to be ready

			f.Format(buf, entry, nil) // format concurrently

			mutex.Lock()
			results[buf.String()] = struct{}{}
			mutex.Unlock()

			finished.Done()
		}(i)
	}

	finished.Wait()
	for i := 1; i <= 100; i++ {
		gt.Expect(results).To(HaveKey(strconv.Itoa(i)))
	}
}

func TestShortFuncFormatter(t *testing.T) {
	gt := NewGomegaWithT(t)

	callerpc, _, _, ok := runtime.Caller(0)
	gt.Expect(ok).To(BeTrue())
	buf := &bytes.Buffer{}
	entry := zapcore.Entry{Caller: zapcore.EntryCaller{PC: callerpc}}
	ShortFuncFormatter{FormatVerb: "%s"}.Format(buf, entry, nil)
	gt.Expect(buf.String()).To(Equal("TestShortFuncFormatter"))

	buf = &bytes.Buffer{}
	entry = zapcore.Entry{Caller: zapcore.EntryCaller{PC: 0}}
	ShortFuncFormatter{FormatVerb: "%s"}.Format(buf, entry, nil)
	gt.Expect(buf.String()).To(Equal("(unknown)"))
}

func TestTimeFormatter(t *testing.T) {
	gt := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	entry := zapcore.Entry{Time: time.Date(1975, time.August, 15, 12, 0, 0, 333, time.UTC)}
	f := TimeFormatter{Layout: time.RFC3339Nano}
	f.Format(buf, entry, nil)
	gt.Expect(buf.String()).To(Equal("1975-08-15T12:00:00.000000333Z"))
}

func TestMultiFormatter(t *testing.T) {
	entry := zapcore.Entry{
		Message: "message",
		Level:   zapcore.InfoLevel,
	}
	fields := []zapcore.Field{
		zap.String("name", "value"),
	}

	var tests = []struct {
		desc     string
		initial  []Formatter
		update   []Formatter
		expected string
	}{
		{
			desc:     "no formatters",
			initial:  nil,
			update:   nil,
			expected: "",
		},
		{
			desc:     "initial formatters",
			initial:  []Formatter{StringFormatter{Value: "string1"}},
			update:   nil,
			expected: "string1",
		},
		{
			desc:    "set to formatters",
			initial: []Formatter{StringFormatter{Value: "string1"}},
			update: []Formatter{
				StringFormatter{Value: "string1"},
				StringFormatter{Value: "-"},
				StringFormatter{Value: "string2"},
			},
			expected: "string1-string2",
		},
		{
			desc:     "set to empty",
			initial:  []Formatter{StringFormatter{Value: "string1"}},
			update:   []Formatter{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			mf := NewMultiFormatter(tc.initial...)
			if tc.update != nil {
				mf.SetFormatters(tc.update)
			}

			buf := &bytes.Buffer{}
			mf.Format(buf, entry, fields)
			gt.Expect(buf.String()).To(Equal(tc.expected))
		})
	}
}
