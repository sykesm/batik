// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"github.com/pkg/errors"
)

const (
	// ErrHalt indicates transaction submission should be halted
	ErrHalt = errorString("halt processing")
)

// errorString is a converstion type for constant errors.
type errorString string

func (e errorString) Error() string { return string(e) }

// fatalError represents a fatal error encountered during transaction processing.
type fatalError struct {
	cause error
	kind  errorString
}

func (f *fatalError) Error() string {
	return f.kind.Error() + ": " + f.cause.Error()
}

func (f *fatalError) Unwrap() error {
	return f.kind
}

func newHaltError(err error, msg string, args ...interface{}) error {
	fe := &fatalError{kind: ErrHalt, cause: err}
	return errors.WithMessagef(fe, msg, args...)
}
