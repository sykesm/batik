// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

var buf *bytes.Buffer

func TestRegisterExitHandler(t *testing.T) {
	gt := NewGomegaWithT(t)

	clearExitHandlers()
	RegisterExitHandler(testHandler1)
	gt.Expect(handlers).To(HaveLen(1))
}

func TestExit(t *testing.T) {
	gt := NewGomegaWithT(t)

	clearExitHandlers()
	RegisterExitHandler(testHandler1)
	RegisterExitHandler(testHandler2)

	buf = &bytes.Buffer{}
	Exit()
	gt.Expect(buf).To(MatchRegexp("testHandler2testHandler1"))

}

func clearExitHandlers() {
	hLock.Lock()
	defer hLock.Unlock()

	handlers = []func(){}
}

func testHandler1() {
	fmt.Fprintf(buf, "testHandler1")
}
func testHandler2() {
	fmt.Fprintf(buf, "testHandler2")
}
