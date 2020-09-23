// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"strconv"
	"time"
)

// ParseUnixTime parses a unix epoch passed as a string.
// It will error if the string cannot be parsed to a float64.
func ParseUnixTime(value string) (time.Time, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return time.Time{}, err
	}

	v := int64(f)

	switch {
	case v > 1e18:
	case v > 1e15:
		v *= 1e3
	case v > 1e12:
		v *= 1e6
	default:
		return time.Unix(v, 0), nil
	}

	return time.Unix(v/1e9, v%1e9), nil
}
