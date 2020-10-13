// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"

	cli "github.com/urfave/cli/v2"
)

// These flags exist to work around an interaction between the cli package, our
// configuration patterns, and the behavior of the Go standard library flag package.
//
// The app creates an instance of a configuration struct that is shared across
// all of the subcommands. The flags used to configure elements of the config
// structure are obtained from this configuration structure during
// instantiation. Later, when the app starts running, the configuration file is
// read and any elements explicitly set in the configuration file are updated
// in the runtime configuration. Finally, when the command executes, the flags
// obtained during instantiation are added to flag sets prior to running the
// command actions. When the flags are added to the flag sets, the "default
// value" established at instantiation time is pushed back to the configuration
// object. This results in any config file overrides being overwritten by the
// initial defaults.
//
// To work around this, we override the Apply method on the cli Flag
// implementations to update the "Value" field from the "Destination" before
// adding the flag to the flag set. This causes the standard library to push
// the current value back to the destination instead of the stale value
// populated during instantiation.

// These are the types that are supported by the cli package. Please implement
// and test as needed for configuration.
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
