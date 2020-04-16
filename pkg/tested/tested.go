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
			t.Errorf("remove tempdir: %v", err)
		}
	}
	return td, cleanup
}

func Close(t TestingT, c io.Closer) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if err := c.Close(); err != nil {
		t.Errorf("close: %v", err)
	}
}
