// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"strconv"
	"time"
)

var formats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.UnixDate,
	time.RubyDate,
	time.ANSIC,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func parseTimeFloat64(value float64) time.Time {
	v := int64(value)
	switch {
	case v > 1e18:
	case v > 1e15:
		v *= 1e3
	case v > 1e12:
		v *= 1e6
	default:
		return time.Unix(v, 0)
	}

	return time.Unix(v/1e9, v%1e9)
}

// tries to parse time using a couple of formats before giving up
func tryParseTime(value interface{}) (time.Time, bool) {
	var t time.Time
	var err error
	switch value.(type) {
	case string:
		if t, err := strconv.ParseFloat(value.(string), 64); err == nil {
			return parseTimeFloat64(t), true
		}

		for _, layout := range formats {
			t, err = time.Parse(layout, value.(string))
			if err == nil {
				return t, true
			}
		}
	}
	return t, false
}
