// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"fmt"
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

var DefaultOptions = &Options{
	SortLongest: true,
	LightBg:     false,
	TimeFormat:  time.Stamp,

	KeyColor:              color.FgGreen,
	ValColor:              color.FgHiWhite,
	TimeLightBgColor:      color.FgBlack,
	TimeDarkBgColor:       color.FgWhite,
	MsgLightBgColor:       color.FgBlack,
	MsgAbsentLightBgColor: color.FgHiBlack,
	MsgDarkBgColor:        color.FgHiWhite,
	MsgAbsentDarkBgColor:  color.FgWhite,
	DebugLevelColor:       color.FgMagenta,
	InfoLevelColor:        color.FgCyan,
	WarnLevelColor:        color.FgYellow,
	ErrorLevelColor:       color.FgRed,
	PanicLevelColor:       color.BgRed,
	FatalLevelColor:       color.BgHiRed,
	UnknownLevelColor:     color.FgMagenta,
}

type Options struct {
	SortLongest bool
	LightBg     bool
	TimeFormat  string

	KeyColor              color.Color
	ValColor              color.Color
	TimeLightBgColor      color.Color
	TimeDarkBgColor       color.Color
	MsgLightBgColor       color.Color
	MsgAbsentLightBgColor color.Color
	MsgDarkBgColor        color.Color
	MsgAbsentDarkBgColor  color.Color
	DebugLevelColor       color.Color
	InfoLevelColor        color.Color
	WarnLevelColor        color.Color
	ErrorLevelColor       color.Color
	PanicLevelColor       color.Color
	FatalLevelColor       color.Color
	UnknownLevelColor     color.Color
}

// LogfmtEntry stores information about a logfmt log line and can transform the
// output to a color formatted human readable log line.
type LogfmtEntry struct {
	buf *bytes.Buffer
	out *tabwriter.Writer

	Opts *Options

	level   string
	time    time.Time
	message string
	fields  map[string]string
}

// UnmarshalLogfmt attempts to unmarshal a logfmt byte slice to a LogfmtEntry
// for prettification. If it cannot unmarshal the logfmt line, it will return nil.
func UnmarshalLogfmt(data []byte) *LogfmtEntry {
	buf := bytes.NewBuffer(nil)
	out := tabwriter.NewWriter(buf, 0, 1, 0, '\t', 0)
	entry := &LogfmtEntry{
		buf:  buf,
		out:  out,
		Opts: DefaultOptions,
	}

	if !bytes.ContainsRune(data, '=') {
		return nil
	}

	dec := logfmt.NewDecoder(bytes.NewReader(data))
	for dec.ScanRecord() {
	next_kv:
		for dec.ScanKeyval() {
			key := dec.Key()
			val := dec.Value()
			// process time
			if entry.time.IsZero() {
				foundTime := checkEachUntilFound(supportedTimeFields, func(field string) bool {
					if !bytes.Equal(key, []byte(field)) {
						return false
					}
					time, ok := tryParseTime(string(val))
					if ok {
						entry.time = time
					}
					return ok
				})
				if foundTime {
					continue next_kv
				}
			}

			// process message
			if len(entry.message) == 0 {
				foundMessage := checkEachUntilFound(supportedMessageFields, func(field string) bool {
					if !bytes.Equal(key, []byte(field)) {
						return false
					}
					entry.message = string(val)
					return true
				})
				if foundMessage {
					continue next_kv
				}
			}

			// process log level
			if len(entry.level) == 0 {
				foundLevel := checkEachUntilFound(supportedLevelFields, func(field string) bool {
					if !bytes.Equal(key, []byte(field)) {
						return false
					}
					entry.level = string(val)
					return true
				})
				if foundLevel {
					continue next_kv
				}
			}

			// process all other key/value fields
			entry.setField(key, val)
		}
	}

	if dec.Err() != nil {
		return nil
	}

	return entry
}

// Prettify the output in a logrus like fashion.
func (e *LogfmtEntry) Prettify() []byte {
	var (
		msgColor       color.Color
		msgAbsentColor color.Color
	)
	if e.Opts.LightBg {
		msgColor = e.Opts.MsgLightBgColor
		msgAbsentColor = e.Opts.MsgAbsentLightBgColor
	} else {
		msgColor = e.Opts.MsgDarkBgColor
		msgAbsentColor = e.Opts.MsgAbsentDarkBgColor
	}

	var msg string
	if e.message == "" {
		msg = msgAbsentColor.Sprint("<no msg>")
	} else {
		msg = msgColor.Sprint(e.message)
	}

	lvl := strings.ToUpper(e.level)[:imin(4, len(e.level))]
	var level string
	switch e.level {
	case "debug":
		level = e.Opts.DebugLevelColor.Sprint(lvl)
	case "info":
		level = e.Opts.InfoLevelColor.Sprint(lvl)
	case "warn", "warning":
		level = e.Opts.WarnLevelColor.Sprint(lvl)
	case "error":
		level = e.Opts.ErrorLevelColor.Sprint(lvl)
	case "fatal", "panic":
		level = e.Opts.FatalLevelColor.Sprint(lvl)
	default:
		level = e.Opts.UnknownLevelColor.Sprint(lvl)
	}

	var timeColor color.Color
	if e.Opts.LightBg {
		timeColor = e.Opts.TimeLightBgColor
	} else {
		timeColor = e.Opts.TimeDarkBgColor
	}
	_, _ = fmt.Fprintf(e.out, "%s |%s| %s\t %s\n",
		timeColor.Sprint(e.time.Format(e.Opts.TimeFormat)),
		level,
		msg,
		strings.Join(e.joinKVs("="), "\t "),
	)

	_ = e.out.Flush()

	return e.buf.Bytes()
}

func (e *LogfmtEntry) setField(key, val []byte) {
	if e.fields == nil {
		e.fields = make(map[string]string)
	}
	e.fields[string(key)] = string(val)
}

func (e *LogfmtEntry) joinKVs(sep string) []string {

	kv := make([]string, 0, len(e.fields))
	for k, v := range e.fields {
		kstr := e.Opts.KeyColor.Sprint(k)
		vstr := e.Opts.ValColor.Sprint(v)
		kv = append(kv, kstr+sep+vstr)
	}

	sort.Strings(kv)

	if e.Opts.SortLongest {
		sort.Stable(byLongest(kv))
	}

	return kv
}

type byLongest []string

func (s byLongest) Len() int           { return len(s) }
func (s byLongest) Less(i, j int) bool { return len(s[i]) < len(s[j]) }
func (s byLongest) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

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
