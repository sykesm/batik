// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-logfmt/logfmt"
	"github.com/sykesm/batik/pkg/log/pretty/color"
)

// supportedTimeFields enumerates supported timestamp field names
var supportedTimeFields = []string{"time", "ts", "@timestamp", "timestamp"}

// supportedMessageFields enumarates supported Message field names
var supportedMessageFields = []string{"message", "msg"}

// supportedLevelFields enumarates supported level field names
var supportedLevelFields = []string{"level", "lvl", "loglevel", "severity"}

const (
	TimeFormat string = time.StampMicro

	KeyColor              color.Color = color.FgGreen
	ValColor                          = color.FgHiWhite
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
	var (
		finalLevel string
		finalTime  time.Time
		message    string

		fields map[string]string
	)

	fields = make(map[string]string)

	if !bytes.ContainsRune(p, '=') {
		return 0, errors.New("not a logfmt string")
	}

	dec := logfmt.NewDecoder(bytes.NewReader(p))
	for dec.ScanRecord() {
	next_kv:
		for dec.ScanKeyval() {
			key := dec.Key()
			val := dec.Value()
			// process time
			if finalTime.IsZero() {
				foundTime := checkEachUntilFound(supportedTimeFields, func(field string) bool {
					if !bytes.Equal(key, []byte(field)) {
						return false
					}
					time, ok := tryParseTime(string(val))
					if ok {
						finalTime = time
					}
					return ok
				})
				if foundTime {
					continue next_kv
				}
			}

			// process message
			if len(message) == 0 {
				foundMessage := checkEachUntilFound(supportedMessageFields, func(field string) bool {
					if !bytes.Equal(key, []byte(field)) {
						return false
					}
					message = string(val)
					return true
				})
				if foundMessage {
					continue next_kv
				}
			}

			// process log level
			if len(finalLevel) == 0 {
				foundLevel := checkEachUntilFound(supportedLevelFields, func(field string) bool {
					if !bytes.Equal(key, []byte(field)) {
						return false
					}
					finalLevel = string(val)
					return true
				})
				if foundLevel {
					continue next_kv
				}
			}

			// process all other key/value fields
			fields[string(key)] = string(val)
		}
	}

	if dec.Err() != nil {
		return 0, dec.Err()
	}

	buf := bytes.NewBuffer(nil)
	out := tabwriter.NewWriter(buf, 0, 1, 0, '\t', 0)

	msgColor := MsgDarkBgColor
	msgAbsentColor := MsgAbsentDarkBgColor

	var msg string
	if message == "" {
		msg = msgAbsentColor.Sprint("<no msg>")
	} else {
		msg = msgColor.Sprint(message)
	}

	lvl := strings.ToUpper(finalLevel)[:imin(4, len(finalLevel))]
	var level string
	switch finalLevel {
	case "debug":
		level = DebugLevelColor.Sprint(lvl)
	case "info":
		level = InfoLevelColor.Sprint(lvl)
	case "warn", "warning":
		level = WarnLevelColor.Sprint(lvl)
	case "error":
		level = ErrorLevelColor.Sprint(lvl)
	case "fatal", "panic":
		level = FatalLevelColor.Sprint(lvl)
	default:
		level = UnknownLevelColor.Sprint(lvl)
	}

	timeColor := TimeDarkBgColor
	kv := make([]string, 0, len(fields))
	for k, v := range fields {
		kstr := KeyColor.Sprint(k)
		vstr := ValColor.Sprint(v)
		kv = append(kv, kstr+"="+vstr)
	}

	sort.Strings(kv)
	_, _ = fmt.Fprintf(out, "%s |%s| %s\t %s\n",
		timeColor.Sprint(finalTime.Format(TimeFormat)),
		level,
		msg,
		strings.Join(kv, "\t "),
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

// checkEachUntilFound searches a field list for a specific field based on
// the found func.
func checkEachUntilFound(fieldList []string, found func(string) bool) bool {
	for _, field := range fieldList {
		if found(field) {
			return true
		}
	}
	return false
}
