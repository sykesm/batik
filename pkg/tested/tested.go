// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tested

import (
	"io"
	"io/ioutil"
	"os"
)

type TestingT interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type tHelper interface {
	Helper()
}

// TempDir creates a temporary directory and returns the path to the directory
// and a function that should be called to cleanup. Any errors encountered
// while creating or removing the temporary directory will result in a test
// failure.
//
// The general pattern looks like this:
//
//     temp, cleanup := TempDir(t, "")
//     defer cleanup()
//
// When TB.TestDir and TB.Cleanup arrive in Go 1.15, this will no longer be
// necessary.
func TempDir(t TestingT, dir, pattern string) (tempdir string, cleanup func()) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	td, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	cleanup = func() {
		if err := os.RemoveAll(td); err != nil {
			t.Errorf("failed to remove tempdir: %v", err)
		}
	}
	return td, cleanup
}

// Close closes an io.Closer and will fail the test if Close fails.
func Close(t TestingT, c io.Closer) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if err := c.Close(); err != nil {
		t.Errorf("failed to close: %v", err)
	}
}
