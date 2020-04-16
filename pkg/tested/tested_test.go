// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tested

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type tb struct {
	testing.TB
	fatalfMsg string
	errorfMsg string
}

func (tb *tb) Fatalf(format string, args ...interface{}) {
	tb.fatalfMsg = fmt.Sprintf(format, args...)
}

func (tb *tb) Errorf(format string, args ...interface{}) {
	tb.errorfMsg = fmt.Sprintf(format, args...)
}

func TestTempDir(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "tempdir")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func(td string) {
		if err := os.RemoveAll(td); err != nil {
			t.Errorf("remove tempdir: %v", err)
		}
	}(tempdir)

	t.Run("Success", func(t *testing.T) {
		dir, cleanup := TempDir(t, tempdir, "success")
		if filepath.Dir(dir) != tempdir {
			t.Errorf("dir: expected %s but got %s", tempdir, filepath.Dir(dir))
		}
		fi, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("stat: %v", err)
		}
		if !fi.IsDir() {
			t.Fatalf("stat: expected %s to be a directory", dir)
		}

		cleanup()
		_, err = os.Stat(dir)
		if !os.IsNotExist(errors.Unwrap(err)) {
			t.Fatalf("cleanup: expected %s to be removed", dir)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		ro := filepath.Join(tempdir, "ro")
		if err := os.Mkdir(ro, 0110); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		tt := &tb{TB: t}
		TempDir(tt, ro, "fail")
		if tt.fatalfMsg == "" {
			t.Fatalf("expected failure to fail test")
		}
	})
}

type closeFunc func() error

func (c closeFunc) Close() error {
	return c()
}

func TestClose(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var closed bool
		success := func() error {
			closed = true
			return nil
		}

		Close(t, closeFunc(success))
		if closed != true {
			t.Fatalf("expected close to be called")
		}
	})

	t.Run("Fail", func(t *testing.T) {
		fail := func() error {
			return errors.New("boom!")
		}

		tt := &tb{TB: t}
		Close(tt, closeFunc(fail))
		if tt.errorfMsg == "" {
			t.Fatalf("expected close to fail the test")
		}
	})
}
