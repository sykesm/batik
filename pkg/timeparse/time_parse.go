// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package timeparse

import (
	"errors"
	"strconv"
	"time"
)

// ParseUnixTime parses a unix epoch passed as either a string or float64.
// It will error on any other type of value or if the string cannot be parsed
// to a float64.
func ParseUnixTime(value interface{}) (time.Time, error) {
	var (
		f   float64
		err error
	)

	switch value.(type) {
	case string:
		f, err = strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return time.Time{}, err
		}
	case float64:
		f = value.(float64)
	default:
		return time.Time{}, errors.New("unix time is not string or float64")
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
