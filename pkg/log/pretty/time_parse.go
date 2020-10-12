// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ParseUnixTime parses a unix epoch passed as a string. It returns an error
// when the string cannot be parsed as a float representing the number of
// seconds elapsed since the Unix epoch.
func ParseUnixTime(value string) (time.Time, error) {
	// Split value into values representing seconds and nanoseconds
	pieces := strings.Split(value, ".")
	if len(pieces) > 2 {
		return time.Time{}, errors.Errorf("ParseUnixTime: invalid syntax: %q", value)
	}
	for len(pieces) < 2 {
		pieces = append(pieces, "0")
	}
	for i := range pieces {
		if len(pieces[i]) == 0 {
			pieces[i] = "0"
		}
	}

	// Treat the first part of the value as a uint64 containing seconds and treat
	// the second part as a fractional float by prepending with 0. The nanos will
	// be scaled by 1e9 to construct the original Unix time.
	secs, err := strconv.ParseUint(pieces[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	nanos, err := strconv.ParseFloat("0."+pieces[1], 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(secs), int64(nanos*1e9)), nil
}
