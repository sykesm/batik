// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-logfmt/logfmt"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log/pretty/color"
)

const (
	timeFormat string = time.StampMicro

	keyColor color.Color = color.FgGreen
	valColor             = color.FgHiWhite

	nameColor   = color.FgBlue
	timeColor   = color.FgWhite
	callerColor = color.None
	msgColor    = color.FgHiWhite

	debugLevelColor   = color.FgMagenta
	infoLevelColor    = color.FgCyan
	warnLevelColor    = color.FgYellow
	errorLevelColor   = color.FgRed
	panicLevelColor   = color.BgRed
	fatalLevelColor   = color.BgHiRed
	unknownLevelColor = color.FgMagenta
)

// A pretty.Writer prettifies logfmt lines and writes them to the underlying writer.
type Writer struct {
	// writer is the outpt stream we write pretty lines to.
	writer io.Writer
	// The zapcore.EncoderConfig used by the logger for encoding logfmt lines.
	encoderConfig zapcore.EncoderConfig
	// Function for parsing time values in a logfmt line.
	parseTime TimeParser
}

// A TimeParser is used to convert encoded time stamps to a time.Time.
type TimeParser func(string) (time.Time, error)

func NewWriter(w io.Writer, e zapcore.EncoderConfig, parseTime TimeParser) *Writer {
	return &Writer{
		writer:        w,
		encoderConfig: e,
		parseTime:     parseTime,
	}
}

type keyValuePair struct {
	key, value string
}

type header struct {
	time   string
	level  string
	name   string
	caller string
	msg    string
	parsed []keyValuePair
}

// Write prettifys input logfmt lines and writes them to the underlying writer.
// If the lines don't adhere to the logfmt format, an error is returned.
func (w *Writer) Write(p []byte) (int, error) {
	if !bytes.ContainsRune(p, '=') {
		return 0, errors.New("not a logfmt string")
	}

	reserved := map[string]struct{}{
		w.encoderConfig.TimeKey:    {},
		w.encoderConfig.LevelKey:   {},
		w.encoderConfig.NameKey:    {},
		w.encoderConfig.CallerKey:  {},
		w.encoderConfig.MessageKey: {},
	}

	var h header
	var parsed []keyValuePair
	dec := logfmt.NewDecoder(bytes.NewReader(p))
	dec.ScanRecord()

	// Scan until all reserved fields are read or until the end of the
	// record. The goal is to reduce memory pressure from parsing the
	// entire record in one go.
	for len(reserved) != 0 && dec.ScanKeyval() {
		kv := keyValuePair{
			key:   string(dec.Key()),
			value: string(dec.Value()),
		}
		if _, ok := reserved[kv.key]; ok {
			delete(reserved, kv.key)
			switch kv.key {
			case w.encoderConfig.TimeKey:
				h.time = w.formatTime(kv.value)
			case w.encoderConfig.LevelKey:
				h.level = formatLevel(kv.value)
			case w.encoderConfig.NameKey:
				h.name = formatName(kv.value)
			case w.encoderConfig.CallerKey:
				h.caller = formatCaller(kv.value)
			case w.encoderConfig.MessageKey:
				h.msg = formatMessage(kv.value)
			default:
				panic("unexpected reserved key: " + kv.key)
			}
			continue
		}
		parsed = append(parsed, kv)
	}

	// Print common header
	if _, err := fmt.Fprintf(w.writer,
		"%s |%s| %s %s %s",
		h.time, h.level, h.name, h.caller, h.msg,
	); err != nil {
		return 0, err
	}

	// Print any non-header fields that have already been parsed
	for i := range parsed {
		if _, err := fmt.Fprintf(w.writer, " %s=%s",
			keyColor.Sprint(parsed[i].key),
			valColor.Sprint(quoteString(parsed[i].value)),
		); err != nil {
			return 0, err
		}
	}
	parsed = nil

	// Print all remaining fields
	for dec.ScanKeyval() {
		if _, err := fmt.Fprintf(w.writer, " %s=%s",
			keyColor.Sprint(dec.Key()),
			valColor.Sprint(quoteString(string(dec.Value()))),
		); err != nil {
			return 0, err
		}
	}

	// Write a newline
	if _, err := fmt.Fprintf(w.writer, "\n"); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (w *Writer) formatTime(value string) string {
	t, _ := w.parseTime(value)
	return timeColor.Sprint(t.Format(timeFormat))
}

func formatLevel(value string) string {
	var c color.Color
	switch strings.ToLower(value) {
	case "debug":
		c = debugLevelColor
	case "info":
		c = infoLevelColor
	case "warn", "warning":
		c = warnLevelColor
	case "error":
		c = errorLevelColor
	case "fatal", "panic":
		c = fatalLevelColor
	default:
		c = unknownLevelColor
	}
	if len(value) < 4 {
		value += "    "
	}
	value = strings.ToUpper(value)[:4]
	return c.Sprint(value)
}

func formatName(value string) string {
	return nameColor.Sprint(value)
}

func formatCaller(value string) string {
	return callerColor.Sprint(value)
}

func formatMessage(value string) string {
	return msgColor.Sprint(value)
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func poolBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

func needsQuotedValueRune(r rune) bool {
	return r <= ' ' || r == '=' || r == '"' || r == utf8.RuneError
}

func quoteString(s string) string {
	const hex = "0123456789abcdef"

	if strings.IndexFunc(s, needsQuotedValueRune) == -1 {
		return s
	}

	buf := getBuffer()
	defer poolBuffer(buf)

	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \n, \r, and \t.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.String()
}
