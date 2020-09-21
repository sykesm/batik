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
	"github.com/sykesm/batik/pkg/log/pretty/color"
)

type reservedField string

const (
	TimeFormat string = time.StampMicro

	// Reserved logfmt fields
	timeField   reservedField = "ts"
	lvlField                  = "level"
	nameField                 = "logger"
	callerField               = "caller"
	msgField                  = "msg"

	KeyColor              color.Color = color.FgGreen
	ValColor                          = color.FgHiWhite
	LoggerColor                       = color.FgBlue
	TimeLightBgColor                  = color.FgBlack
	TimeDarkBgColor                   = color.FgWhite
	MsgLightBgColor                   = color.FgBlack
	MsgAbsentLightBgColor             = color.FgHiBlack
	MsgDarkBgColor                    = color.FgHiWhite
	MsgAbsentDarkBgColor              = color.FgWhite
	DebugLevelColor                   = color.FgMagenta
	InfoLevelColor                    = color.FgCyan
	WarnLevelColor                    = color.FgYellow
	ErrorLevelColor                   = color.FgRed
	PanicLevelColor                   = color.BgRed
	FatalLevelColor                   = color.BgHiRed
	UnknownLevelColor                 = color.FgMagenta
)

// A pretty.Writer prettifies logfmt lines and writes them to the underlying writer.
type Writer struct {
	io.Writer
}

// Write prettifys input logfmt lines and writes them onto the underlying writer.
// If the lines aren't logfmt, an error is returned.
func (w *Writer) Write(p []byte) (n int, err error) {
	if !bytes.ContainsRune(p, '=') {
		return 0, errors.New("not a logfmt string")
	}

	dec := logfmt.NewDecoder(bytes.NewReader(p))
	dec.ScanRecord()

	parsedTime := parseReservedTimeField(dec, timeField)
	parsedLvl := parseReservedStringField(dec, lvlField)
	parsedName := parseReservedStringField(dec, nameField)
	parsedCaller := parseReservedStringField(dec, callerField)
	parsedMsg := parseReservedStringField(dec, msgField)

	fields := parseUnreservedFields(dec)

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

func parseReservedTimeField(dec *logfmt.Decoder, f reservedField) string {
	dec.ScanKeyval()

	key := dec.Key()
	val := dec.Value()

	if !bytes.Equal(key, []byte(f)) {
		return ""
	}

	if time, ok := tryParseTime(string(val)); ok {
		return TimeDarkBgColor.Sprint(time.Format(TimeFormat))
	}

	return ""
}

func parseReservedStringField(dec *logfmt.Decoder, f reservedField) string {
	dec.ScanKeyval()

	key := string(dec.Key())
	val := string(dec.Value())

	if key != string(f) {
		return ""
	}

	var c color.Color
	switch f {
	case nameField:
		c = LoggerColor
	case lvlField:
		switch val {
		case "debug":
			c = DebugLevelColor
		case "info":
			c = InfoLevelColor
		case "warn", "warning":
			c = WarnLevelColor
		case "error":
			c = ErrorLevelColor
		case "fatal", "panic":
			c = FatalLevelColor
		default:
			c = UnknownLevelColor
		}
		val = strings.ToUpper(val)[:imin(4, len(val))]
	case msgField:
		if val == "" {
			c = MsgAbsentDarkBgColor
			val = "<no msg>"
		} else {
			c = MsgDarkBgColor
		}
	}

	return c.Sprint(val)
}

func parseUnreservedFields(dec *logfmt.Decoder) []string {
	var fields []string

	for dec.ScanKeyval() {
		key := dec.Key()
		val := dec.Value()

		kstr := KeyColor.Sprint(string(key))
		vstr := ValColor.Sprint(string(val))
		fields = append(fields, kstr+"="+vstr)
	}

	return fields
}
