// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package atexit

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAtExit(t *testing.T) {
	gt := NewGomegaWithT(t)

	atexit := New()

	buf := &bytes.Buffer{}
	atexit.Register(func() { fmt.Fprint(buf, "testHandler1") })
	atexit.Register(func() { fmt.Fprint(buf, "testHandler2") })

	gt.Expect(buf.Len()).To(Equal(0))
	atexit.Exit()
	gt.Expect(buf.String()).To(Equal("testHandler2testHandler1"), "should execute LIFO")

	atexit.Exit()
	gt.Expect(buf.String()).To(Equal("testHandler2testHandler1"), "should only run once")
}
