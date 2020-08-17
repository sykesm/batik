// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
)

func newLookupError(key string) error {
	return &lookupError{
		key: key,
	}
}

type lookupError struct {
	key string
}

func (e *lookupError) Error() string {
	return fmt.Sprintf("$%s is not defined", e.key)
}

type Lookuper interface {
	Lookup(key string) (string, error)
}

type EnvMap map[string]string

var _ Lookuper = (*EnvMap)(nil)

func (m EnvMap) Lookup(key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", newLookupError(key)
	}

	return v, nil
}

type OsEnv struct{}

var _ Lookuper = (*OsEnv)(nil)

func (o OsEnv) Lookup(key string) (string, error) {
	v, exists := os.LookupEnv(key)
	if !exists {
		return "", newLookupError(key)
	}

	return v, nil
}
