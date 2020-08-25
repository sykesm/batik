// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
)

// A Lookuper is used to lookup the value of a variable from the process
// environment. When a value is found, implementations must return the value
// and true. If a value is not found, implementaions must return the empty
// string and false.
//
// The semantics implied by this interface are consistent with the behavior of
// os.LookupEnv from the standard library.
type Lookuper interface {
	Lookup(name string) (string, bool)
}

// EnvironLookuper uses os.LookupEnv to resolve environment variables.
func EnvironLookuper() Lookuper {
	return &environLookuper{}
}

type environLookuper struct{}

func (*environLookuper) Lookup(name string) (string, bool) {
	return os.LookupEnv(name)
}

// MapLookuper uses the provided map to resolve environment variables.
func MapLookuper(m map[string]string) Lookuper {
	return mapLookuper(m)
}

type mapLookuper map[string]string

func (m mapLookuper) Lookup(name string) (string, bool) {
	val, ok := m[name]
	return val, ok
}
