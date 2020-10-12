// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-logfmt/logfmt"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log/pretty/color"
)

const (
	timeFormat string = time.StampMicro

	keyColor color.Color = color.FgGreen
	valColor             = color.FgHiWhite

	nameColor      = color.FgBlue
	timeColor      = color.FgWhite
	callerColor    = color.None
	msgColor       = color.FgHiWhite
	msgAbsentColor = color.FgWhite

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
	io.Writer

	// The zapcore.EncoderConfig used by the logger for encoding logfmt lines.
	encoderConfig zapcore.EncoderConfig
	// Maps a reserved key to a function controlling how to colorize output for values
	// associated with that key.
	colorFuncs map[string]func(f string) string
	// Function for parsing time values in a logfmt line.
	parseTime TimeParser
}

// A TimeParser is used to convert encoded time stamps to a time.Time.
type TimeParser func(string) (time.Time, error)

func NewWriter(w io.Writer, e zapcore.EncoderConfig, parseTime TimeParser) *Writer {
	colorFuncs := map[string]func(f string) string{
		e.NameKey: func(v string) string { return nameColor.Sprint(v) },
		e.LevelKey: func(v string) string {
			var c color.Color
			switch v {
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
			v = strings.ToUpper(v)[:imin(4, len(v))]
			return c.Sprint(v)
		},
		e.MessageKey: func(v string) string {
			if v == "" {
				return msgAbsentColor.Sprint("<no msg>")
			}
			return msgColor.Sprint(v)
		},
		e.CallerKey: func(v string) string { return callerColor.Sprint(v) },
		e.TimeKey:   func(v string) string { return timeColor.Sprint(v) },
	}

	return &Writer{
		Writer:        w,
		encoderConfig: e,
		colorFuncs:    colorFuncs,
		parseTime:     parseTime,
	}
}

// Write prettifys input logfmt lines and writes them onto the underlying writer.
// If the lines aren't logfmt, an error is returned.
// TODO: If the loglines are not in proper order (time, level, name, caller, message, fields),
// or if a reserved field is missing (such as if the logger is unnamed), the Write will miss
// subsequent fields. We should probably add workarounds to this behavior.
func (w *Writer) Write(p []byte) (n int, err error) {
	if !bytes.ContainsRune(p, '=') {
		return 0, errors.New("not a logfmt string")
	}

	dec := logfmt.NewDecoder(bytes.NewReader(p))
	dec.ScanRecord()

	parsedTime := w.parseReservedTimeField(dec, w.encoderConfig.TimeKey)
	parsedLvl := w.parseReservedStringField(dec, w.encoderConfig.LevelKey)
	parsedName := w.parseReservedStringField(dec, w.encoderConfig.NameKey)
	parsedCaller := w.parseReservedStringField(dec, w.encoderConfig.CallerKey)
	parsedMsg := w.parseReservedStringField(dec, w.encoderConfig.MessageKey)

	fields := w.parseUnreservedFields(dec)

	if dec.Err() != nil {
		return 0, dec.Err()
	}

	buf := bytes.NewBuffer(nil)
	out := tabwriter.NewWriter(buf, 0, 1, 0, '\t', 0)

	_, _ = fmt.Fprintf(out, "%s |%s| %s\t %s\t %s\t %s\n",
		parsedTime,
		parsedLvl,
		parsedName,
		parsedCaller,
		parsedMsg,
		strings.Join(fields, "\t "),
	)

	_ = out.Flush()

	return w.Writer.Write(buf.Bytes())
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (w *Writer) parseReservedTimeField(dec *logfmt.Decoder, f string) string {
	dec.ScanKeyval()

	key := string(dec.Key())
	val := string(dec.Value())

	if key != f {
		return ""
	}

	if t, err := w.parseTime(val); err == nil {
		if c, ok := w.colorFuncs[key]; ok {
			return c(t.Format(timeFormat))
		}
	}

	return ""
}

func (w *Writer) parseReservedStringField(dec *logfmt.Decoder, f string) string {
	dec.ScanKeyval()

	key := string(dec.Key())
	val := string(dec.Value())

	if key != f {
		return ""
	}

	if c, ok := w.colorFuncs[key]; ok {
		return c(val)
	}

	return ""
}

func (w *Writer) parseUnreservedFields(dec *logfmt.Decoder) []string {
	var fields []string

	for dec.ScanKeyval() {
		key := string(dec.Key())
		val := string(dec.Value())
		fields = append(fields, keyColor.Sprint(key)+"="+valColor.Sprint(val))
	}

	return fields
}
