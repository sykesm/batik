// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"

	cli "github.com/urfave/cli/v2"
)

// These are the types that are supported by the cli package. Please
// implement and test as needed for configuration.
//
//   [ ] BoolFlag
//   [x] DurationFlag
//   [ ] Float64Flag
//   [ ] Float64SliceFlag
//   [ ] GenericFlag
//   [ ] Int64Flag
//   [ ] Int64SliceFlag
//   [ ] IntFlag
//   [ ] IntSliceFlag
//   [ ] PathFlag
//   [x] StringFlag
//   [ ] StringSliceFlag
//   [ ] TimestampFlag
//   [ ] Uint64Flag
//   [x] UintFlag

type DurationFlag struct {
	*cli.DurationFlag
}

func NewDurationFlag(f *cli.DurationFlag) *DurationFlag {
	return &DurationFlag{DurationFlag: f}
}

func (d *DurationFlag) Apply(fs *flag.FlagSet) error {
	if d.DurationFlag.Destination != nil {
		d.DurationFlag.Value = *d.DurationFlag.Destination
	}
	return d.DurationFlag.Apply(fs)
}

type StringFlag struct {
	*cli.StringFlag
}

func NewStringFlag(f *cli.StringFlag) *StringFlag {
	return &StringFlag{StringFlag: f}
}

func (s *StringFlag) Apply(fs *flag.FlagSet) error {
	if s.StringFlag.Destination != nil {
		s.StringFlag.Value = *s.StringFlag.Destination
	}

	return s.StringFlag.Apply(fs)
}

type UintFlag struct {
	*cli.UintFlag
}

func NewUintFlag(f *cli.UintFlag) *UintFlag {
	return &UintFlag{UintFlag: f}
}

func (u *UintFlag) Apply(fs *flag.FlagSet) error {
	if u.UintFlag.Destination != nil {
		u.UintFlag.Value = *u.UintFlag.Destination
	}
	return u.UintFlag.Apply(fs)
}
